

-- name: ListGuilds :many
SELECT guild_id, description
FROM guilds
ORDER BY guild_id ASC;

-- name: AddGuild :exec
INSERT INTO guilds (
    guild_id,
    description
) VALUES (?, ?);


-- name: GetGuild :many
SELECT guild_id, description
FROM guilds
WHERE guild_id = ?
LIMIT 1;


-- name: RemoveGuild :exec
DELETE FROM guilds WHERE guild_id = ?;