package bot

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

type Bot struct {
	state *state.State
}

func New(token string) (*Bot, error) {
	s := state.New("Bot " + token)
	bot := &Bot{
		state: s,
	}

	r := cmdroute.NewRouter()

	r.AddFunc("ping", bot.Ping)

	s.AddInteractionHandler(r)
	s.AddIntents(
		gateway.IntentGuilds | gateway.IntentGuildMessages | gateway.IntentGuildMessageReactions | gateway.IntentGuildMessageTyping | gateway.IntentDirectMessages | gateway.IntentDirectMessageReactions | gateway.IntentDirectMessageTyping | gateway.IntentMessageContent,
	)

	err := cmdroute.OverwriteCommands(s, []api.CreateCommandData{
		{
			Name:        "ping",
			Description: "Ping!",
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
