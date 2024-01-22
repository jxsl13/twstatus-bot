

-- name: GetPlayerCountNotificationRequest :many
SELECT
	guild_id,
	channel_id,
	message_id,
	user_id,
	threshold
FROM player_count_notification_requests
WHERE guild_id = $1
AND channel_id = $2
AND message_id = $3
AND user_id = $4
LIMIT 1;


-- name: SetPlayerCountNotificationRequest :exec
INSERT INTO player_count_notification_requests (
	guild_id,
	channel_id,
	message_id,
	user_id,
	threshold
) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (guild_id, channel_id, message_id, user_id)
DO UPDATE SET threshold = $5;


-- name: RemovePlayerCountNotificationRequests :exec
DELETE FROM player_count_notification_requests;


-- name: RemovePlayerCountNotificationRequest :exec
DELETE FROM player_count_notification_requests
WHERE guild_id = $1
AND channel_id = $2
AND message_id = $3
AND user_id = $4
AND threshold = $5;
