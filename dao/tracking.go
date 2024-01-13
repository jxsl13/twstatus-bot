package dao

import (
	"context"
	"errors"
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

func ListTrackingsByChannelID(ctx context.Context, conn Conn, guildID discord.GuildID, channelID discord.ChannelID) (trackings model.Trackings, err error) {
	rows, err := conn.QueryContext(ctx, `
SELECT guild_id, channel_id, address, message_id
FROM tracking
WHERE guild_id = ?
AND channel_id = ?
ORDER BY message_id ASC;`, guildID, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trackings for channel %d: %w", channelID, err)
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	for rows.Next() {
		var tracking model.Tracking
		err = rows.Scan(
			&tracking.GuildID,
			&tracking.ChannelID,
			&tracking.Address,
			&tracking.MessageID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tracking: %w", err)
		}
		trackings = append(trackings, tracking)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate tracking: %w", err)
	}
	return trackings, nil
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

func RemoveTrackingByMessageID(ctx context.Context, conn Conn, guildID discord.GuildID, messageID discord.MessageID) (err error) {
	_, err = conn.ExecContext(ctx, `
DELETE FROM tracking
WHERE guild_id = ?
AND message_id = ?;`, guildID, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete tracking for message id: %s: %w", messageID, err)
	}
	return nil
}
