package dao

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

func ListAllTrackings(ctx context.Context, q *sqlc.Queries) (trackings model.Trackings, err error) {
	latr, err := q.ListAllTrackings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get trackings: %w", err)
	}

	result := make(model.Trackings, 0, len(latr))
	for _, t := range latr {
		result = append(result, model.Tracking{
			MessageTarget: model.MessageTarget{
				ChannelTarget: model.ChannelTarget{
					GuildID:   discord.GuildID(t.GuildID),
					ChannelID: discord.ChannelID(t.ChannelID),
				},
				MessageID: discord.MessageID(t.MessageID),
			},
			Address: t.Address,
		})
	}

	return result, nil
}

func ListTrackingsByChannelID(ctx context.Context, q *sqlc.Queries, guildID discord.GuildID, channelID discord.ChannelID) (trackings model.Trackings, err error) {
	latr, err := q.ListChannelTrackings(ctx, sqlc.ListChannelTrackingsParams{
		GuildID:   int64(guildID),
		ChannelID: int64(channelID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get trackings: %w", err)
	}
	result := make(model.Trackings, 0, len(latr))
	for _, t := range latr {
		result = append(result, model.Tracking{
			MessageTarget: model.MessageTarget{
				ChannelTarget: model.ChannelTarget{
					GuildID:   guildID,
					ChannelID: channelID,
				},
				MessageID: discord.MessageID(t.MessageID),
			},
			Address: t.Address,
		})
	}
	return result, nil
}

func AddTracking(ctx context.Context, q *sqlc.Queries, tracking model.Tracking) (err error) {
	cs, err := q.GetChannel(ctx, sqlc.GetChannelParams{
		GuildID:   int64(tracking.GuildID),
		ChannelID: int64(tracking.ChannelID),
	})
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
	}
	if len(cs) == 0 {
		return fmt.Errorf("channel %s is not known", tracking.ChannelID)
	}

	found, err := ExistsServer(ctx, q, tracking.Address)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("server %s does not exist", tracking.Address)
	}

	err = q.AddTracking(ctx, sqlc.AddTrackingParams{
		GuildID:   int64(tracking.GuildID),
		ChannelID: int64(tracking.ChannelID),
		Address:   tracking.Address,
		MessageID: int64(tracking.MessageID),
	})
	if err != nil {
		if IsUniqueConstraintErr(err) {
			return fmt.Errorf("%w: tracking %s", ErrAlreadyExists, tracking.Address)
		}
		return fmt.Errorf("failed to insert tracking for %s: %w", tracking.Address, err)
	}
	return nil

}

func RemoveTrackingByMessageID(ctx context.Context, q *sqlc.Queries, guildID discord.GuildID, messageID discord.MessageID) (err error) {
	err = q.RemoveTrackingByMessageId(
		ctx,
		sqlc.RemoveTrackingByMessageIdParams{
			GuildID:   int64(guildID),
			MessageID: int64(messageID),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to remove tracking by message id: %w", err)
	}
	return nil
}
