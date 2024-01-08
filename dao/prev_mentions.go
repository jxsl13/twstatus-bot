package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
)

// TODO: continue here
func ChangedMessageMentions(
	ctx context.Context,
	tx *sql.Tx,
	currentMentions model.MessageMentions,
) (
	messageMentions model.MessageMentions,
	err error,
) {

	// removed mentions
	// changed mentions
	// unchanged mentions

	return messageMentions, nil
}

func ListPrevMessageMentions(ctx context.Context, conn Conn) (messageMentions model.MessageMentions, err error) {
	rows, err := conn.QueryContext(ctx, `
SELECT
	guild_id,
	channel_id,
	message_id,
	user_id
FROM prev_message_mentions
ORDER BY guild_id ASC, channel_id ASC, message_id ASC, user_id ASC;`)
	if err != nil {
		return nil, fmt.Errorf("failed to query previous channel notifications: %w", err)
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	messageMentions = make(model.MessageMentions, 64)
	for rows.Next() {
		var mt model.MessageTarget
		var userID discord.UserID
		err = rows.Scan(
			&mt.GuildID,
			&mt.ChannelID,
			&mt.MessageID,
			&userID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan previous channel notification: %w", err)
		}
		messageMentions[mt] = append(messageMentions[mt], userID)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate previous channel notifications: %w", err)
	}

	return messageMentions, nil
}

func RemoveMessageMentions(ctx context.Context, conn Conn, messageMentions model.MessageMentions) (err error) {
	stmt, err := conn.PrepareContext(ctx, `
DELETE FROM prev_message_mentions
WHERE guild_id = ?
AND channel_id = ?
AND message_id = ?;`)
	if err != nil {
		return fmt.Errorf("failed to prepare delete previous message mention statement: %w", err)
	}
	defer func() {
		err = errors.Join(err, stmt.Close())
	}()

	for mt := range messageMentions {
		_, err = stmt.ExecContext(ctx, mt.GuildID, mt.ChannelID, mt.MessageID)
		if err != nil {
			return fmt.Errorf("failed to delete previous message mention: %w", err)
		}
	}

	return nil
}
