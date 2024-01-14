package dao

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

func StartChannel(ctx context.Context, q *sqlc.Queries, guildID discord.GuildID, channelID discord.ChannelID) (c model.Channel, err error) {
	channel, err := GetChannel(ctx, q, guildID, channelID)
	if err != nil {
		return c, err
	}

	if channel.Running {
		return c, fmt.Errorf("channel %s is already active", channel)
	}

	err = q.StartChannel(ctx, sqlc.StartChannelParams{
		GuildID:   int64(guildID),
		ChannelID: int64(channelID),
	})
	if err != nil {
		return c, fmt.Errorf("failed to start channel %s: %w", channel, err)
	}
	return channel, nil
}

func StopChannel(ctx context.Context, q *sqlc.Queries, guildID discord.GuildID, channelID discord.ChannelID) (c model.Channel, err error) {
	channel, err := GetChannel(ctx, q, guildID, channelID)
	if err != nil {
		return c, err
	}

	if !channel.Running {
		return c, fmt.Errorf("channel %s is already inactive", channel)
	}

	err = q.StopChannel(ctx, sqlc.StopChannelParams{
		GuildID:   int64(guildID),
		ChannelID: int64(channelID),
	})
	if err != nil {
		return c, fmt.Errorf("failed to stop channel %s: %w", channel, err)
	}
	return channel, nil
}
