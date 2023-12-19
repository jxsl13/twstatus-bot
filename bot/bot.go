package bot

import (
	"context"
	"database/sql"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

type Bot struct {
	state *state.State
	db    *sql.DB
}

// New requires a discord bot token and returns a Bot instance.
// A bot token starts with Nj... and can be obtained from the discord developer portal.
func New(token string, db *sql.DB) (*Bot, error) {
	s := state.New("Bot " + token)
	bot := &Bot{
		state: s,
		db:    db,
	}

	s.AddHandler(func(*gateway.ReadyEvent) {
		me, _ := s.Me()
		log.Println("connected to the gateway as", me.Tag())
	})

	r := cmdroute.NewRouter()

	r.AddFunc("ping", bot.Ping)
	r.AddFunc("register", bot.registerServer)

	s.AddInteractionHandler(r)
	s.AddIntents(
		gateway.IntentGuilds,
	)

	err := cmdroute.OverwriteCommands(s, []api.CreateCommandData{
		{
			Name:        "ping",
			Description: "Ping!",
		},
		{
			Name:        "register",
			Description: "Register a server",
			Options: []discord.CommandOption{
				&discord.StringOption{
					OptionName:  "address",
					Description: "ipv4:port or [ipv6]:port",
					MinLength:   option.NewInt(16),
					MaxLength:   option.NewInt(64),
					Required:    true,
				},
			},
			/*DefaultMemberPermissions: discord.NewPermissions(
				discord.PermissionAdministrator,
			),*/
		},
	})
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

func errorResponse(err error) *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Content:         option.NewNullableString("**Error:** " + err.Error()),
		Flags:           discord.EphemeralMessage,
		AllowedMentions: &api.AllowedMentions{ /* none */ },
	}
}
