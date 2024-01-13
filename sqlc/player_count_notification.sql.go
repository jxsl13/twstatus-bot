// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: player_count_notification.sql

package sqlc

import (
	"context"
)

const getMessageTargetNotifications = `-- name: GetMessageTargetNotifications :many
SELECT
	user_id,
	threshold
FROM player_count_notifications
WHERE guild_id = ?
AND channel_id = ?
AND message_id = ?
ORDER BY user_id ASC
`

type GetMessageTargetNotificationsParams struct {
	GuildID   int64
	ChannelID int64
	MessageID int64
}

type GetMessageTargetNotificationsRow struct {
	UserID    int64
	Threshold int64
}

func (q *Queries) GetMessageTargetNotifications(ctx context.Context, arg GetMessageTargetNotificationsParams) ([]GetMessageTargetNotificationsRow, error) {
	rows, err := q.query(ctx, q.getMessageTargetNotificationsStmt, getMessageTargetNotifications, arg.GuildID, arg.ChannelID, arg.MessageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetMessageTargetNotificationsRow{}
	for rows.Next() {
		var i GetMessageTargetNotificationsRow
		if err := rows.Scan(&i.UserID, &i.Threshold); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPlayerCountNotification = `-- name: GetPlayerCountNotification :many
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
LIMIT 1
`

type GetPlayerCountNotificationParams struct {
	GuildID   int64
	ChannelID int64
	MessageID int64
	UserID    int64
}

func (q *Queries) GetPlayerCountNotification(ctx context.Context, arg GetPlayerCountNotificationParams) ([]PlayerCountNotification, error) {
	rows, err := q.query(ctx, q.getPlayerCountNotificationStmt, getPlayerCountNotification,
		arg.GuildID,
		arg.ChannelID,
		arg.MessageID,
		arg.UserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PlayerCountNotification{}
	for rows.Next() {
		var i PlayerCountNotification
		if err := rows.Scan(
			&i.GuildID,
			&i.ChannelID,
			&i.MessageID,
			&i.UserID,
			&i.Threshold,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPlayerCountNotifications = `-- name: ListPlayerCountNotifications :many
SELECT
	guild_id,
	channel_id,
	message_id,
	user_id,
	threshold
FROM player_count_notifications
ORDER BY
    guild_id ASC,
    channel_id ASC,
    message_id ASC,
    user_id ASC
`

func (q *Queries) ListPlayerCountNotifications(ctx context.Context) ([]PlayerCountNotification, error) {
	rows, err := q.query(ctx, q.listPlayerCountNotificationsStmt, listPlayerCountNotifications)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PlayerCountNotification{}
	for rows.Next() {
		var i PlayerCountNotification
		if err := rows.Scan(
			&i.GuildID,
			&i.ChannelID,
			&i.MessageID,
			&i.UserID,
			&i.Threshold,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removePlayerCountNotification = `-- name: RemovePlayerCountNotification :exec
DELETE FROM player_count_notifications
WHERE guild_id = ?
AND channel_id = ?
AND message_id = ?
AND user_id = ?
AND threshold = ?
`

type RemovePlayerCountNotificationParams struct {
	GuildID   int64
	ChannelID int64
	MessageID int64
	UserID    int64
	Threshold int64
}

func (q *Queries) RemovePlayerCountNotification(ctx context.Context, arg RemovePlayerCountNotificationParams) error {
	_, err := q.exec(ctx, q.removePlayerCountNotificationStmt, removePlayerCountNotification,
		arg.GuildID,
		arg.ChannelID,
		arg.MessageID,
		arg.UserID,
		arg.Threshold,
	)
	return err
}

const removePlayerCountNotifications = `-- name: RemovePlayerCountNotifications :exec
DELETE FROM player_count_notifications
`

func (q *Queries) RemovePlayerCountNotifications(ctx context.Context) error {
	_, err := q.exec(ctx, q.removePlayerCountNotificationsStmt, removePlayerCountNotifications)
	return err
}

const setPlayerCountNotification = `-- name: SetPlayerCountNotification :exec
REPLACE INTO player_count_notifications (
	guild_id,
	channel_id,
	message_id,
	user_id,
	threshold
) VALUES (?, ?, ?, ?, ?)
`

type SetPlayerCountNotificationParams struct {
	GuildID   int64
	ChannelID int64
	MessageID int64
	UserID    int64
	Threshold int64
}

func (q *Queries) SetPlayerCountNotification(ctx context.Context, arg SetPlayerCountNotificationParams) error {
	_, err := q.exec(ctx, q.setPlayerCountNotificationStmt, setPlayerCountNotification,
		arg.GuildID,
		arg.ChannelID,
		arg.MessageID,
		arg.UserID,
		arg.Threshold,
	)
	return err
}
