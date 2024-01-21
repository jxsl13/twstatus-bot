package dao

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

func (dao *DAO) StartChannel(ctx context.Context, guildID discord.GuildID, channelID discord.ChannelID) (c model.Channel, err error) {
	channel, err := dao.GetChannel(ctx, guildID, channelID)
	if err != nil {
		return c, err
	}

	if channel.Running {
		return c, fmt.Errorf("channel %s is already active", channel)
	}

	err = dao.q.StartChannel(ctx, sqlc.StartChannelParams{
		GuildID:   int64(guildID),
		ChannelID: int64(channelID),
	})
	if err != nil {
		return c, fmt.Errorf("failed to start channel %s: %w", channel, err)
	}
	return channel, nil
}

func (dao *DAO) StopChannel(ctx context.Context, guildID discord.GuildID, channelID discord.ChannelID) (c model.Channel, err error) {
	channel, err := dao.GetChannel(ctx, guildID, channelID)
	if err != nil {
		return c, err
	}

	if !channel.Running {
		return c, fmt.Errorf("channel %s is already inactive", channel)
	}

	err = dao.q.StopChannel(ctx, sqlc.StopChannelParams{
		GuildID:   int64(guildID),
		ChannelID: int64(channelID),
	})
	if err != nil {
		return c, fmt.Errorf("failed to stop channel %s: %w", channel, err)
	}
	return channel, nil
}
