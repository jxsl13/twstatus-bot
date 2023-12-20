package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

var commandList = []api.CreateCommandData{
	{
		Name:        "list-guilds",
		Description: "List all guilds that are allowed to use this bot",
	},
	{
		Name:        "add-guild",
		Description: "Allow a guild to use this bot",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "description",
				Description: "A description for this guild.",
				MinLength:   option.NewInt(4),
				MaxLength:   option.NewInt(256),
				Required:    true,
			},
			&discord.StringOption{
				OptionName:  "id",
				Description: "The guild id of the guild you want to add.",
				MinLength:   option.NewInt(1),
				Required:    false,
			},
		},
	},
	{
		Name:        "remove-guild",
		Description: "Remove a guild from the allowed guilds",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "id",
				Description: "The guild id of the guild you want to remove.",
				MinLength:   option.NewInt(1),
				Required:    false,
			},
		},
	},
	{
		Name:        "add-channel",
		Description: "Add a channel to the allowed channels",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "id",
				Description: "The channel id of the channel you want to add.",
				MinLength:   option.NewInt(1),
				MaxLength:   option.NewInt(64),
				Required:    false,
			},
		},
	},
	{
		Name:        "remove-channel",
		Description: "Remove a channel from the allowed channels",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "id",
				Description: "The channel id of the channel you want to remove.",
				MinLength:   option.NewInt(1),
				MaxLength:   option.NewInt(64),
				Required:    false,
			},
		},
	},
	{
		Name:        "list-channels",
		Description: "List all channels of the current guild that are registered for this bot",
	},
	{
		Name:        "update-servers",
		Description: "Update the server list",
	},
}

type Bot struct {
	state       *state.State
	db          *sql.DB
	superAdmins []discord.UserID
}

// New requires a discord bot token and returns a Bot instance.
// A bot token starts with Nj... and can be obtained from the discord developer portal.
func New(ctx context.Context, token string, db *sql.DB, superAdmins []discord.UserID, guildID discord.GuildID) (*Bot, error) {
	s := state.New("Bot " + token)
	app, err := s.CurrentApplication()
	if err != nil {
		return nil, fmt.Errorf("failed to get current application: %w", err)
	}

	bot := &Bot{
		state:       s,
		db:          db,
		superAdmins: superAdmins,
	}

	s.AddHandler(func(*gateway.ReadyEvent) {
		me, _ := s.Me()
		log.Println("connected to the gateway as", me.Tag())
		src, dst, err := bot.updateServers(ctx)
		if err != nil {
			log.Printf("failed to initialize server list: %v", err)
		} else {
			log.Printf("initialized server list with %d source and %d target servers", src, dst)
		}
	})

	r := cmdroute.NewRouter()

	r.AddFunc("list-guilds", bot.listGuilds)
	r.AddFunc("add-guild", bot.addGuild)
	r.AddFunc("remove-guild", bot.removeGuild)
	r.AddFunc("list-channels", bot.listChannels)
	r.AddFunc("add-channel", bot.addChannel)
	r.AddFunc("remove-channel", bot.removeChannel)
	r.AddFunc("update-servers", bot.updateServerList)

	s.AddInteractionHandler(r)
	s.AddIntents(
		gateway.IntentGuilds,
	)

	if guildID != discord.NullGuildID {
		_, err = s.BulkOverwriteGuildCommands(app.ID, guildID, commandList)
		if err != nil {
			return nil, err
		}
	} else {
		err = cmdroute.OverwriteCommands(s, commandList)
		if err != nil {
			return nil, err
		}
	}

	return bot, nil
}

func (b *Bot) Connect(ctx context.Context) error {
	return b.state.Connect(ctx)
}

func (b *Bot) Close() error {
	return b.state.Close()
}

func (b *Bot) IsSuperAdmin(userID discord.UserID) bool {
	for _, admin := range b.superAdmins {
		if admin == userID {
			return true
		}
	}
	return false
}

func ErrAccessForbidden() *api.InteractionResponseData {
	return errorResponse(fmt.Errorf("access forbidden"))
}

func errorResponse(err error) *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Content:         option.NewNullableString("**Error:** " + err.Error()),
		Flags:           discord.EphemeralMessage,
		AllowedMentions: &api.AllowedMentions{ /* none */ },
	}
}

// Tx returns a transaction and a closer function.
func (b *Bot) Tx(ctx context.Context) (*sql.Tx, func(error) error, error) {
	closer := func(err error) error { return err }
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, closer, err
	}

	return tx, func(err error) error {
		if err != nil {
			return errors.Join(err, tx.Rollback())
		}
		return tx.Commit()
	}, nil
}
