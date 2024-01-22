

-- name: GetPlayerCountNotificationMessages :many
SELECT
	t.guild_id,
	t.channel_id,
	COALESCE(pcm.message_id, 0)::bigint AS prev_message_id,
	pcn.user_id,
	MAX(COALESCE(np.num_players, 0))::smallint AS num_players
FROM channels c
JOIN tracking t ON c.channel_id = t.channel_id
LEFT JOIN (
	SELECT ac.address, count(*) AS num_players
	FROM active_server_clients ac
	GROUP BY ac.address
    ORDER BY ac.address
) np ON np.address = t.address
JOIN player_count_notifications pcn
ON (
	t.guild_id = pcn.guild_id AND
	t.channel_id = pcn.channel_id AND
	t.message_id = pcn.message_id AND
	num_players >= pcn.threshold
)
LEFT JOIN player_count_notification_messages pcm
ON (t.channel_id = pcm.channel_id)
WHERE c.running = TRUE
GROUP BY t.guild_id, t.channel_id, pcm.message_id, pcn.user_id, num_players
ORDER BY t.guild_id, t.channel_id, pcm.message_id, num_players, pcn.user_id;


-- name: AddPlayerCountNotificationMessage :exec
INSERT INTO player_count_notification_messages (channel_id, message_id)
VALUES ($1, $2)
ON CONFLICT (channel_id)
DO UPDATE SET
    message_id = EXCLUDED.message_id;

-- name: RemovePlayerCountNotificationMessage :exec
DELETE FROM player_count_notification_messages
WHERE channel_id = $1
AND message_id = $2;