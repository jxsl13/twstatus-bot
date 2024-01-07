package dao

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jxsl13/twstatus-bot/model"
)

func GetTargetListNotifications(ctx context.Context, tx *sql.Tx, servers map[model.Target]model.ChangedServerStatus) (map[model.Target]model.ChangedServerStatus, error) {
	stmt, err := tx.PrepareContext(ctx, `
	SELECT
	user_id,
	threshold
FROM player_count_notifications
WHERE guild_id = ?
AND channel_id = ?
AND message_id = ?
ORDER BY user_id ASC;`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare list notifications statement: %w", err)
	}
	defer stmt.Close()

	for t := range servers {
		rows, err := stmt.QueryContext(ctx, t.GuildID, t.ChannelID, t.MessageID)
		if err != nil {
			return nil, fmt.Errorf("failed to query notifications: %w", err)
		}

		for rows.Next() {
			var n model.PlayerCountNotification
			err = rows.Scan(
				&n.UserID,
				&n.Threshold,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to scan notification: %w", err)
			}
			server := servers[t]
			if n.Notify(&server) {
				server.UserNotifications = append(server.UserNotifications, n.UserID)
				servers[t] = server
			}
		}
		err = rows.Err()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate notifications: %w", err)
		}
	}

	return servers, nil
}

func ListAllPlayerCountNotifications(ctx context.Context, conn Conn) (notifications []model.PlayerCountNotification, err error) {
	rows, err := conn.QueryContext(ctx, `
SELECT
	guild_id,
	channel_id,
	message_id,
	user_id,
	threshold
FROM player_count_notifications
ORDER BY guild_id ASC, channel_id ASC, message_id ASC, user_id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n model.PlayerCountNotification
		err = rows.Scan(
			&n.GuildID,
			&n.ChannelID,
			&n.MessageID,
			&n.UserID,
			&n.Threshold,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func GetPlayerCountNotification(ctx context.Context, conn Conn, n model.UserTarget) (notification model.PlayerCountNotification, err error) {

	rows, err := conn.QueryContext(ctx, `
SELECT
	guild_id,
	channel_id,
	message_id,
	user_id,
	threshold
FROM player_count_notifications
WHERE guild_id = ?
AND channel_id = ?
AND message_id = ?
AND user_id = ?
LIMIT 1;`, n.GuildID, n.ChannelID, n.MessageID, n.UserID)
	if err != nil {
		return model.PlayerCountNotification{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return model.PlayerCountNotification{}, fmt.Errorf("no player count notification found for %w", ErrNotFound)
	}
	err = rows.Err()
	if err != nil {
		return model.PlayerCountNotification{}, err
	}

	err = rows.Scan(
		&notification.GuildID,
		&notification.ChannelID,
		&notification.MessageID,
		&notification.UserID,
		&notification.Threshold,
	)
	if err != nil {
		return model.PlayerCountNotification{}, err
	}

	return notification, nil
}
func SetPlayerCountNotifications(ctx context.Context, conn Conn, notifications []model.PlayerCountNotification) (err error) {
	stmt, err := conn.PrepareContext(ctx, `
REPLACE INTO player_count_notifications (
	guild_id,
	channel_id,
	message_id,
	user_id,
	threshold
) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare set player notifications statement: %w", err)
	}
	defer stmt.Close()

	for _, n := range notifications {
		_, err = stmt.ExecContext(ctx,
			n.GuildID,
			n.ChannelID,
			n.MessageID,
			n.UserID,
			n.Threshold,
		)
		if err != nil {
			return fmt.Errorf("failed to set player notification: %w", err)
		}
	}

	return nil
}

func SetPlayerCountNotification(ctx context.Context, conn Conn, n model.PlayerCountNotification) (err error) {
	_, err = conn.ExecContext(ctx, `
REPLACE INTO player_count_notifications (
	guild_id,
	channel_id,
	message_id,
	user_id,
	threshold
) VALUES (?, ?, ?, ?, ?)`,
		n.GuildID,
		n.ChannelID,
		n.MessageID,
		n.UserID,
		n.Threshold,
	)
	return err
}

func RemovePlayerCountNotifications(ctx context.Context, conn Conn) (err error) {
	_, err = conn.ExecContext(ctx, `
DELETE FROM player_count_notifications`)
	return err
}

func RemovePlayerCountNotification(ctx context.Context, conn Conn, n model.PlayerCountNotification) (err error) {
	_, err = conn.ExecContext(ctx, `
DELETE FROM player_count_notifications
WHERE guild_id = ?
AND channel_id = ?
AND message_id = ?
AND user_id = ?
AND threshold = ?;`,
		n.GuildID,
		n.ChannelID,
		n.MessageID,
		n.UserID,
		n.Threshold,
	)
	return err
}
