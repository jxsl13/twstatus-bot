

-- name: GetPlayerCountNotificationMessages :many


SELECT
	t.guild_id,
	t.channel_id,
	pcr.message_id AS req_message_id,
	COALESCE(pcm.message_id, 0)::bigint AS prev_message_id,
	pcr.user_id,
	MIN(pcr.threshold)::smallint AS threshold,
	MAX(COALESCE(np.num_players, 0))::smallint AS num_players
FROM channels c
JOIN tracking t ON c.channel_id = t.channel_id
LEFT JOIN (
	SELECT ac.address, count(*) AS num_players
	FROM active_server_clients ac
	WHERE ac.address = ANY($1::TEXT[])
	GROUP BY ac.address
    ORDER BY ac.address
) np ON np.address = t.address
JOIN player_count_notification_requests pcr
ON (
	t.guild_id = pcr.guild_id AND
	t.channel_id = pcr.channel_id AND
	t.message_id = pcr.message_id AND
	num_players >= pcr.threshold
)
LEFT JOIN player_count_notification_messages pcm
ON (t.channel_id = pcm.channel_id)
WHERE c.running = TRUE
GROUP BY
	t.guild_id,
	t.channel_id,
	pcm.message_id,
	pcr.message_id,
	pcr.user_id,
	num_players
ORDER BY t.guild_id, t.channel_id, pcm.message_id, num_players, pcr.user_id;



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