
-- name: GetChannel :many
SELECT running
FROM channels
WHERE guild_id = $1
AND channel_id = $2
LIMIT 1;


-- name: ListGuildChannels :many
SELECT channel_id, running
FROM channels
WHERE guild_id = $1
ORDER BY channel_id ASC;

-- name: AddGuildChannel :exec
INSERT INTO channels (channel_id, guild_id, running)
VALUES ($1, $2, $3);

-- name: RemoveGuildChannel :exec
DELETE FROM channels
WHERE guild_id = $1
AND channel_id = $2;


-- name: StartChannel :exec
UPDATE channels
SET running = TRUE
WHERE guild_id = $1
AND channel_id = $2;

-- name: StopChannel :exec
UPDATE channels
SET running = FALSE
WHERE guild_id = $1
AND channel_id = $2;