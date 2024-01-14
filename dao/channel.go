package dao

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

func GetChannel(ctx context.Context, q *sqlc.Queries, guildId discord.GuildID, channelID discord.ChannelID) (_ model.Channel, err error) {
	runnings, err := q.GetChannel(ctx, sqlc.GetChannelParams{
		GuildID:   int64(guildId),
		ChannelID: int64(channelID),
	})

	if err != nil {
		return model.Channel{}, fmt.Errorf("failed to query channel: %w", err)
	}

	if len(runnings) == 0 {
		return model.Channel{}, fmt.Errorf("%w: channel %d", ErrNotFound, channelID)
	}

	running := runnings[0]
	return model.Channel{
		GuildID: guildId,
		ID:      channelID,
		Running: running != 0,
	}, nil
}

func ListChannels(ctx context.Context, q *sqlc.Queries, guildID discord.GuildID) (model.Channels, error) {
	channels, err := q.ListGuildChannels(ctx, int64(guildID))
	if err != nil {
		return nil, fmt.Errorf("failed to query channels: %w", err)
	}

	result := make(model.Channels, 0, len(channels))
	for _, c := range channels {
		result = append(result, model.Channel{
			GuildID: guildID,
			ID:      discord.ChannelID(c.ChannelID),
			Running: c.Running != 0,
		})
	}

	return result, nil
}

func AddChannel(ctx context.Context, q *sqlc.Queries, channel model.Channel) (err error) {
	err = q.AddGuildChannel(ctx, sqlc.AddGuildChannelParams{
		GuildID:   int64(channel.GuildID),
		ChannelID: int64(channel.ID),
		Running:   channel.RunningInt64(),
	})
	if err != nil {
		if IsPrimaryKeyConstraintErr(err) {
			return fmt.Errorf("%w: channel %s", ErrAlreadyExists, channel)
		}
		return fmt.Errorf("failed to insert channel %d: %w", channel.ID, err)
	}
	return nil
}

func RemoveChannel(ctx context.Context, q *sqlc.Queries, guildID discord.GuildID, channelID discord.ChannelID) (err error) {
	err = q.RemoveGuildChannel(ctx, sqlc.RemoveGuildChannelParams{
		GuildID:   int64(guildID),
		ChannelID: int64(channelID),
	})
	if err != nil {
		return fmt.Errorf("failed to delete channel %d: %w", channelID, err)
	}
	return nil
}
