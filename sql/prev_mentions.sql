
-- name: ListPreviousMessageMentions :many
SELECT
	guild_id,
	channel_id,
	message_id,
	user_id
FROM prev_message_mentions
ORDER BY guild_id ASC, channel_id ASC, message_id ASC, user_id ASC;


-- name: RemoveMessageMentions :exec
DELETE FROM prev_message_mentions
WHERE guild_id = ?
AND channel_id = ?
AND message_id = ?;