-- name: ListPrevActiveServers :many
SELECT
	message_id,
	guild_id,
	channel_id,
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
FROM prev_active_servers
ORDER BY guild_id ASC, channel_id ASC, message_id ASC;

-- name: AddPrevActiveServer :exec
INSERT INTO prev_active_servers (
	message_id,
	guild_id,
	channel_id,
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
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);


-- name: RemovePrevActiveServer :exec
DELETE FROM prev_active_servers
WHERE message_id = ?;


-- name: GetPrevActiveServerClients :one
SELECT
	message_id,
	guild_id,
	channel_id,
	name,
	clan,
	team,
	country_id,
	score,
	is_player,
	flag_abbr,
	flag_emoji
FROM prev_active_server_clients
WHERE message_id = ?
ORDER BY id ASC
LIMIT 1;


-- name: AddPrevActiveServerClient :exec
INSERT INTO prev_active_server_clients (
	message_id,
	guild_id,
	channel_id,
	name,
	clan,
	team,
	country_id,
	score,
	is_player,
	flag_abbr,
	flag_emoji
) VALUES (?,?,?,?,?,?,?,?,?,?,?);



-- name: RemovePrevActiveServerClient :exec
DELETE FROM prev_active_server_clients
WHERE message_id = ?;

