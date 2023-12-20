package bot

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/jxsl13/twstatus-bot/dao"
	"github.com/jxsl13/twstatus-bot/model"
)

func (b *Bot) listChannels(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	guildId := data.Event.GuildID
	channels, err := dao.ListChannels(ctx, b.db, guildId)
	if err != nil {
		return errorResponse(err)
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(channels.String()),
	}
}

func (b *Bot) addChannel(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
	var (
		guildId   = data.Event.GuildID
		channelId = data.Event.ChannelID
	)

	tx, closer, err := b.Tx(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer func() {
		err = closer(err)
		if err != nil {
			resp = errorResponse(err)
		}
	}()

	channel := model.Channel{
		GuildID: guildId,
		ID:      channelId,
		Running: 0,
	}
	err = dao.AddChannel(ctx, tx, channel)
	if err != nil {
		return errorResponse(err)
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("added channel: %s", channel)),
		Flags:   discord.EphemeralMessage,
	}
}

func (b *Bot) removeChannel(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
	var (
		guildId   = data.Event.GuildID
		channelId = data.Event.ChannelID
	)
	tx, closer, err := b.Tx(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer func() {
		err = closer(err)
		if err != nil {
			resp = errorResponse(err)
		}
	}()

	channel, err := dao.RemoveChannel(ctx, tx, guildId, channelId)
	if err != nil {
		return errorResponse(err)
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("removed channel %s", channel)),
		Flags:   discord.EphemeralMessage,
	}
}
