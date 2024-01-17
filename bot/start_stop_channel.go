package bot

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/jxsl13/twstatus-bot/dao"
)

// optional channel id parameter
func (b *Bot) startChannel(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
	q, closer, err := b.TxQueries(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer func() {
		err = closer(err)
		if err != nil {
			resp = errorResponse(err)
		}
	}()

	channel, err := dao.StartChannel(ctx,
		q,
		data.Event.GuildID,
		optionalChannelID(data),
	)
	if err != nil {
		return errorResponse(err)
	}

	msg := fmt.Sprintf("Started channel: %s", channel)
	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
		Flags:   discord.EphemeralMessage,
	}
}

// optional channel id parameter
func (b *Bot) stopChannel(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
	q, closer, err := b.TxQueries(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer func() {
		err = closer(err)
		if err != nil {
			resp = errorResponse(err)
		}
	}()

	channel, err := dao.StopChannel(ctx,
		q,
		data.Event.GuildID,
		optionalChannelID(data),
	)
	if err != nil {
		return errorResponse(err)
	}

	msg := fmt.Sprintf("Stopped channel: %s", channel)
	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
		Flags:   discord.EphemeralMessage,
	}
}
