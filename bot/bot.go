package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/jxsl13/twstatus-bot/db"
)

const (
	channelOptionName = "channel"
)

var ownerCommandList = []api.CreateCommandData{
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
				MaxLength:   option.NewInt(20),
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
				MaxLength:   option.NewInt(20),
				Required:    false,
			},
		},
	},
	{
		Name:        "update-servers",
		Description: "Update the server list",
	},
	{
		Name:        "update-messages",
		Description: "Update all discord messages with server status lists",
	},
}

var userCommandList = []api.CreateCommandData{
	{
		Name:           "add-channel",
		Description:    "Add a channel to the allowed channels",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to add.",
				Required:    false,
			},
		},
	},
	{
		Name:           "remove-channel",
		Description:    "Remove a channel from the allowed channels",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to remove.",
				Required:    false,
			},
		},
	},
	{
		Name:           "list-channels",
		Description:    "List all channels of the current guild that are registered for this bot",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
	},
	{
		Name:           "list-flag-mappings",
		Description:    "List all flag mappings for the current or given channel",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to list the flag mappings for.",
				Required:    false,
			},
		},
	},
	{
		Name:           "add-flag-mapping",
		Description:    "Add a flag mapping for the current channel",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "abbr",
				Description: "The abbreviation of the flag you want to add a different emoji for.",
				Required:    true,
				MinLength:   option.NewInt(2),
				MaxLength:   option.NewInt(7), // len("default")
			},
			&discord.StringOption{
				OptionName:  "emoji",
				Description: "The emoji you want to use for this flag (any text).",
				Required:    true,
				MinLength:   option.NewInt(1), // :X:
				MaxLength:   option.NewInt(256),
			},
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to add a flag mapping for.",
				Required:    false,
			},
		},
	},
	{
		Name:           "remove-flag-mapping",
		Description:    "Remove a flag mapping for the current or provided channel",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "abbr",
				Description: "The abbreviation of the flag you want to remove a mapping for.",
				Required:    true,
				MinLength:   option.NewInt(2),
				MaxLength:   option.NewInt(7), // len("default")
			},
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to remove a flag mapping for.",
				Required:    false,
			},
		},
	},
	{
		Name:           "list-flags",
		Description:    "show all known flags",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
	},
	{
		Name:           "add-tracking",
		Description:    "Add tracking of a Teeworlds server for the current or given channel",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "address",
				Description: "One or a list of comma separated server addresses that you want to track.",
				Required:    true,
				MinLength:   option.NewInt(9),
			},
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to track the server for.",
				Required:    false,
			},
		},
	},
	{
		Name:           "start",
		Description:    "Start the bot for  the given channel",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to start the bot for.",
				Required:    false,
			},
		},
	},
	{
		Name:           "stop",
		Description:    "Stop the bot for the given channel",
		NoDMPermission: true,
		DefaultMemberPermissions: discord.NewPermissions(
			discord.PermissionAdministrator,
		),
		Options: []discord.CommandOption{
			&discord.ChannelOption{
				OptionName:  channelOptionName,
				Description: "The channel id of the channel you want to stop the bot for.",
				Required:    false,
			},
		},
	},
}

type Bot struct {
	ctx         context.Context
	state       *state.State
	db          *db.DB
	superAdmins []discord.UserID
	useEmbeds   bool
}

// New requires a discord bot token and returns a Bot instance.
// A bot token starts with Nj... and can be obtained from the discord developer portal.
func New(
	ctx context.Context,
	token string,
	db *db.DB,
	superAdmins []discord.UserID,
	guildID discord.GuildID,
	pollingInterval time.Duration,
	legacyMessageFormat bool,
) (*Bot, error) {
	s := state.New("Bot " + token)
	app, err := s.CurrentApplication()
	if err != nil {
		return nil, fmt.Errorf("failed to get current application: %w", err)
	}

	bot := &Bot{
		ctx:         ctx,
		state:       s,
		db:          db,
		superAdmins: superAdmins,
		useEmbeds:   !legacyMessageFormat,
	}

	s.AddIntents(
		gateway.IntentGuilds | gateway.IntentGuildMessages | gateway.IntentGuildMessageReactions,
	)

	s.AddHandler(func(*gateway.ReadyEvent) {
		me, _ := s.Me()
		log.Println("connected to the gateway as", me.Tag())
		src, dst, err := bot.updateServers(ctx)
		if err != nil {
			log.Printf("failed to initialize server list: %v", err)
		} else {
			log.Printf("initialized server list with %d source and %d target servers", src, dst)
		}

		// start polling
		go bot.async(pollingInterval)
	})

	// requires guild message intents
	s.AddHandler(bot.handleMessageDeletion)
	s.AddHandler(bot.handleAddGuild)
	s.AddHandler(bot.handleRemoveGuild)

	r := cmdroute.NewRouter()

	// bot owner commands
	r.AddFunc("list-guilds", bot.listGuilds)
	r.AddFunc("add-guild", bot.addGuildCommand)
	r.AddFunc("remove-guild", bot.removeGuildCommand)
	r.AddFunc("update-servers", bot.updateServerListCommand)
	r.AddFunc("update-messages", bot.updateDiscordMessagesCommand)

	// user commands
	r.AddFunc("list-channels", bot.listChannels)
	r.AddFunc("add-channel", bot.addChannel)
	r.AddFunc("remove-channel", bot.removeChannel)
	r.AddFunc("list-flags", bot.listFlags)
	r.AddFunc("add-flag-mapping", bot.addFlagMapping)
	r.AddFunc("list-flag-mappings", bot.listFlagMappings)
	r.AddFunc("remove-flag-mapping", bot.removeFlagMapping)
	r.AddFunc("add-tracking", bot.addTracking)
	r.AddFunc("start", bot.startChannel)
	r.AddFunc("stop", bot.stopChannel)

	s.AddInteractionHandler(r)

	_, err = s.BulkOverwriteGuildCommands(app.ID, guildID, ownerCommandList)
	if err != nil {
		return nil, err
	}

	// update user facing commands
	err = cmdroute.OverwriteCommands(s, userCommandList)
	if err != nil {
		return nil, err
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
