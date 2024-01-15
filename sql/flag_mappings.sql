
-- name: ListFlagMappings :many
SELECT
    m.flag_id,
	m.emoji,
	f.abbr
FROM flag_mappings m
JOIN flags f ON m.flag_id = f.flag_id
WHERE m.guild_id = $1
AND m.channel_id = $2
ORDER BY f.abbr ASC;

-- name: AddFlagMapping :exec
INSERT INTO flag_mappings (
    guild_id,
    channel_id,
    flag_id,
	emoji
) VALUES ($1, $2, $3, $4)
ON CONFLICT (channel_id, flag_id) DO UPDATE
SET
	guild_id = $1,
    emoji = $4;

-- name: GetFlagMapping :many
SELECT
	m.flag_id,
	m.emoji,
	f.abbr
FROM flag_mappings m
JOIN flags f ON m.flag_id = f.flag_id
WHERE m.guild_id = $1
AND m.channel_id = $2
AND m.flag_id = $3
LIMIT 1;

-- name: RemoveFlagMapping :exec
DELETE FROM flag_mappings
WHERE guild_id = $1
AND channel_id = $2
AND flag_id = $3;

