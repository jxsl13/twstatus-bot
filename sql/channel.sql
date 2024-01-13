
-- name: GetChannel :many
SELECT running
FROM channels
WHERE guild_id = ?
AND channel_id = ?
LIMIT 1;


-- name: ListGuildChannels :many
SELECT channel_id, running
FROM channels
WHERE guild_id = ?
ORDER BY channel_id ASC;

-- name: AddGuildChannel :exec
INSERT INTO channels (channel_id, guild_id, running)
VALUES (?, ?, ?);

-- name: RemoveGuildChannel :exec
DELETE FROM channels
WHERE guild_id = ?
AND channel_id = ?;


-- name: StartChannel :exec
UPDATE channels
SET running = 1
WHERE guild_id = ?
AND channel_id = ?;

-- name: StopChannel :exec
UPDATE channels
SET running = 0
WHERE guild_id = ?
AND channel_id = ?;