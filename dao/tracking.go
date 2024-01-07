package dao

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
)

func ListAllTrackings(ctx context.Context, conn Conn) (trackings model.Trackings, err error) {
	rows, err := conn.QueryContext(ctx, `
SELECT guild_id, channel_id, address, message_id
FROM tracking
ORDER BY guild_id ASC, channel_id ASC;`)
	if err != nil {
		return nil, fmt.Errorf("failed to get trackings: %w", err)
	}
	defer rows.Close()

	result := make(model.Trackings, 0, 64)
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
		result = append(result, tracking)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate tracking: %w", err)
	}
	return result, nil
}

func ListTrackingsByChannelID(ctx context.Context, conn Conn, guildID discord.GuildID, channelID discord.ChannelID) (trackings model.Trackings, err error) {
	rows, err := conn.QueryContext(ctx, `
SELECT guild_id, channel_id, address, message_id
FROM tracking
WHERE guild_id = ?
AND channel_id = ?;`, guildID, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trackings for channel %d: %w", channelID, err)
	}
	defer rows.Close()

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

func AddTracking(ctx context.Context, tx *sql.Tx, tracking model.Tracking) (err error) {
	_, err = GetChannel(ctx, tx, tracking.GuildID, tracking.ChannelID)
	if err != nil {
		return err
	}

	found, err := ExistsServer(ctx, tx, tracking.Address)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("server %s does not exist", tracking.Address)
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO tracking (guild_id, channel_id, address, message_id)
VALUES (?, ?, ?, ?);`,
		tracking.GuildID,
		tracking.ChannelID,
		tracking.Address,
		tracking.MessageID,
	)
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
