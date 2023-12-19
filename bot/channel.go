package bot

import (
	"context"
	"errors"
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
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return errorResponse(err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}

		if err != nil {
			resp = errorResponse(err)
		}
	}()

	var (
		guildId   = data.Event.GuildID
		channelId = data.Event.ChannelID
	)

	channel, err := dao.GetChannel(ctx, tx, guildId, channelId)
	if err != nil && !errors.Is(err, dao.ErrNotFound) {
		return errorResponse(err)
	} else if err == nil {
		// delete previous message
		_ = b.state.DeleteMessage(channel.ID, channel.MessageID, "removed message due to re-registration of Teeworlds server status tracking.")
	}

	// else - new registration

	msg, err := b.state.SendMessage(channelId, "initial message")
	if err != nil {
		return errorResponse(err)
	}
	defer func() {
		if err != nil {
			_ = b.state.DeleteMessage(channelId, msg.ID, "removed message due to failed registration of Teeworlds server status tracking.")
		}
	}()

	err = dao.AddChannel(ctx, tx, model.Channel{
		GuildID:   guildId,
		ID:        channelId,
		MessageID: msg.ID,
		Running:   0,
	})
	if err != nil {
		return errorResponse(err)
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString("added channel"),
		Flags:   discord.EphemeralMessage,
	}
}

func (b *Bot) removeChannel(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	guildId := data.Event.GuildID
	channelId := data.Event.ChannelID

	channel, err := dao.RemoveChannel(ctx, b.db, guildId, channelId)
	if err != nil {
		return errorResponse(err)
	}

	_ = b.state.DeleteMessage(channel.ID, channel.MessageID, "removed message due to deregistration of Teeworlds server status tracking.")

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("removed channel %d", channel.ID)),
		Flags:   discord.EphemeralMessage,
	}
}
