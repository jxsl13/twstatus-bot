package dao

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
)

func AddTracking(ctx context.Context, tx *sql.Tx, tracking model.Tracking) (err error) {
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
