-- name: ListTrackedServers :many
SELECT
	c.guild_id,
	c.channel_id,
	t.message_id,
	ts.timestamp,
	ts.address,
	ts.protocols,
	ts.name,
	ts.gametype,
	ts.passworded,
	ts.map,
	ts.map_sha256sum,
	ts.map_size,
	ts.version,
	ts.max_clients,
	ts.max_players,
	ts.score_kind
FROM channels c
JOIN tracking t ON c.channel_id = t.channel_id
JOIN active_servers ts ON t.address = ts.address
WHERE c.running = 1
ORDER BY c.guild_id ASC, c.channel_id ASC;


-- name: DeleteActiveServers :exec
DELETE FROM active_servers;

-- name: InsertActiveServers :exec
INSERT INTO active_servers (
	timestamp,
	address,
	protocols,
	name,
	gametype,
	passworded,
	map,
	map_sha256sum,
	map_size,
	version,
	max_clients,
	max_players,
	score_kind
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);


-- name: ListTrackedServerClients :many
SELECT
	c.guild_id,
	c.channel_id,
	t.message_id,
	tsc.name,
	tsc.clan,
	tsc.country_id,
	(CASE WHEN tsc.score = -9999 THEN 9223372036854775807 ELSE tsc.score END) as score,
	tsc.is_player,
	tsc.team,
	f.abbr,
	(CASE WHEN fm.emoji NOT NULL THEN fm.emoji ELSE f.emoji END) as flag_emoji
FROM channels c
JOIN tracking t ON c.channel_id = t.channel_id
JOIN active_server_clients tsc ON t.address = tsc.address
JOIN flags f ON tsc.country_id = f.flag_id
LEFT JOIN flag_mappings fm ON
	(
		t.channel_id = fm.channel_id AND
		tsc.country_id = fm.flag_id
	)
WHERE c.running = 1
ORDER BY
    c.guild_id ASC,
    c.channel_id ASC,
    t._rowid_,
    score DESC,
    tsc.name ASC;


-- name: DeleteActiveServerClients :exec
DELETE FROM active_server_clients;


-- name: InsertActiveServerClients :exec
INSERT INTO active_server_clients (
	address,
	name,
	clan,
	country_id,
	score,
	is_player,
	team
) VALUES (?,?,?,?,?,?,?);


-- name: ExistsServer :one
SELECT
	address
FROM active_servers
WHERE address = ?
LIMIT 1;


