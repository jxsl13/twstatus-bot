
-- name: ListFlagMappings :many
SELECT
    m.flag_id,
	m.emoji,
	f.abbr
FROM flag_mappings m
JOIN flags f ON m.flag_id = f.flag_id
WHERE m.guild_id = ?
AND m.channel_id = ?
ORDER BY f.abbr ASC;

-- name: AddFlagMapping :exec
REPLACE INTO flag_mappings (
    guild_id,
    channel_id,
    flag_id, emoji
) VALUES (?, ?, ?, ?);

-- name: GetFlagMapping :one
SELECT
	m.flag_id,
	m.emoji,
	f.abbr
FROM flag_mappings m
JOIN flags f ON m.flag_id = f.flag_id
WHERE m.guild_id = ?
AND m.channel_id = ?
AND m.flag_id = ?
LIMIT 1;

-- name: RemoveFlagMapping :exec
DELETE FROM flag_mappings
WHERE guild_id = ?
AND channel_id = ?
AND flag_id = ?;

