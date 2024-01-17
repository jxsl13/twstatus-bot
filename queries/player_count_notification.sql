

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
WHERE guild_id = $1
AND channel_id = $2
AND message_id = $3
ORDER BY user_id ASC;


-- name: GetPlayerCountNotification :many
SELECT
	guild_id,
	channel_id,
	message_id,
	user_id,
	threshold
FROM player_count_notifications
WHERE guild_id = $1
AND channel_id = $2
AND message_id = $3
AND user_id = $4
LIMIT 1;


-- name: SetPlayerCountNotification :exec
INSERT INTO player_count_notifications (
	guild_id,
	channel_id,
	message_id,
	user_id,
	threshold
) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (guild_id, channel_id, message_id, user_id)
DO UPDATE SET threshold = $5;


-- name: RemovePlayerCountNotifications :exec
DELETE FROM player_count_notifications;


-- name: RemovePlayerCountNotification :exec
DELETE FROM player_count_notifications
WHERE guild_id = $1
AND channel_id = $2
AND message_id = $3
AND user_id = $4
AND threshold = $5;
