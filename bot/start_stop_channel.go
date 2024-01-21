package bot

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

// optional channel id parameter
func (b *Bot) startChannel(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
	dao, closer, err := b.TxDAO(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer func() {
		err = closer(err)
		if err != nil {
			resp = errorResponse(err)
		}
	}()

	channel, err := dao.StartChannel(
		ctx,
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
	dao, closer, err := b.TxDAO(ctx)
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
