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
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);


-- name: RemovePrevActiveServer :exec
DELETE FROM prev_active_servers
WHERE message_id = $1;


-- name: GetPrevActiveServerClients :many
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
WHERE message_id = $1
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
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);



-- name: RemovePrevActiveServerClient :exec
DELETE FROM prev_active_server_clients
WHERE message_id = $1;

