

-- name: ListPlayerCountNotifications :many
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
    user_id ASC;


-- name: GetMessageTargetNotifications :many
SELECT
	user_id,
	threshold
FROM player_count_notifications
WHERE guild_id = ?
AND channel_id = ?
AND message_id = ?
ORDER BY user_id ASC;


-- name: GetPlayerCountNotification :many
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
LIMIT 1;


-- name: SetPlayerCountNotification :exec
REPLACE INTO player_count_notifications (
	guild_id,
	channel_id,
	message_id,
	user_id,
	threshold
) VALUES (?, ?, ?, ?, ?);


-- name: RemovePlayerCountNotifications :exec
DELETE FROM player_count_notifications;


-- name: RemovePlayerCountNotification :exec
DELETE FROM player_count_notifications
WHERE guild_id = ?
AND channel_id = ?
AND message_id = ?
AND user_id = ?
AND threshold = ?;
