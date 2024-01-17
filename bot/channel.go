package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/jxsl13/twstatus-bot/dao"
	"github.com/jxsl13/twstatus-bot/model"
)

func (b *Bot) listChannels(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
	q, closer, err := b.ConnQueries(ctx)
	if err != nil {
		return errorResponse(err)
	}
	defer closer()

	guildId := data.Event.GuildID
	channels, err := dao.ListChannels(ctx, q, guildId)
	if err != nil {
		return errorResponse(err)
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(channels.StatusString()),
		Flags:   discord.EphemeralMessage,
	}
}

func (b *Bot) addChannel(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
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

	channel := model.Channel{
		GuildID: data.Event.GuildID,
		ID:      optionalChannelID(data),
		Running: false,
	}
	err = dao.AddChannel(ctx, q, channel)
	if err != nil {
		return errorResponse(err)
	}

	msg := fmt.Sprintf("Added channel: %s", channel)
	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
		Flags:   discord.EphemeralMessage,
	}
}

func (b *Bot) removeChannel(ctx context.Context, data cmdroute.CommandData) (resp *api.InteractionResponseData) {
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

	var (
		guildID   = data.Event.GuildID
		channelID = optionalChannelID(data)
	)

	channel, err := dao.GetChannel(ctx, q, guildID, channelID)
	if err != nil {
		return errorResponse(err)
	}

	trackings, err := dao.ListTrackingsByChannelID(ctx, q, guildID, channelID)
	if err != nil {
		return errorResponse(err)
	}

	msgIDs := make([]discord.MessageID, 0, len(trackings))
	for _, t := range trackings {
		msgIDs = append(msgIDs, t.MessageID)
	}

	delErr := b.state.DeleteMessages(channelID, msgIDs, "channel was removed")
	if delErr != nil {
		log.Printf("failed to delete messages: %v", delErr)
	}

	err = dao.RemoveChannel(
		ctx,
		q,
		guildID,
		channelID,
	)
	if err != nil {
		return errorResponse(err)
	}

	msg := fmt.Sprintf("removed channel %s", channel)
	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
		Flags:   discord.EphemeralMessage,
	}
}
