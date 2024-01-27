// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: prev_active_servers.sql

package sqlc

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

const addPrevActiveServer = `-- name: AddPrevActiveServer :exec
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
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
`

type AddPrevActiveServerParams struct {
	MessageID    int64              `db:"message_id"`
	GuildID      int64              `db:"guild_id"`
	ChannelID    int64              `db:"channel_id"`
	Timestamp    pgtype.Timestamptz `db:"timestamp"`
	Address      string             `db:"address"`
	Protocols    json.RawMessage    `db:"protocols"`
	Name         string             `db:"name"`
	Gametype     string             `db:"gametype"`
	Passworded   bool               `db:"passworded"`
	Map          string             `db:"map"`
	MapSha256sum *string            `db:"map_sha256sum"`
	MapSize      *int32             `db:"map_size"`
	Version      string             `db:"version"`
	MaxClients   int16              `db:"max_clients"`
	MaxPlayers   int16              `db:"max_players"`
	ScoreKind    string             `db:"score_kind"`
}

func (q *Queries) AddPrevActiveServer(ctx context.Context, arg AddPrevActiveServerParams) error {
	_, err := q.db.Exec(ctx, addPrevActiveServer,
		arg.MessageID,
		arg.GuildID,
		arg.ChannelID,
		arg.Timestamp,
		arg.Address,
		arg.Protocols,
		arg.Name,
		arg.Gametype,
		arg.Passworded,
		arg.Map,
		arg.MapSha256sum,
		arg.MapSize,
		arg.Version,
		arg.MaxClients,
		arg.MaxPlayers,
		arg.ScoreKind,
	)
	return err
}

const addPrevActiveServerClient = `-- name: AddPrevActiveServerClient :exec
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
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`

type AddPrevActiveServerClientParams struct {
	MessageID int64  `db:"message_id"`
	GuildID   int64  `db:"guild_id"`
	ChannelID int64  `db:"channel_id"`
	Name      string `db:"name"`
	Clan      string `db:"clan"`
	Team      *int16 `db:"team"`
	CountryID int16  `db:"country_id"`
	Score     int32  `db:"score"`
	IsPlayer  bool   `db:"is_player"`
	FlagAbbr  string `db:"flag_abbr"`
	FlagEmoji string `db:"flag_emoji"`
}

func (q *Queries) AddPrevActiveServerClient(ctx context.Context, arg AddPrevActiveServerClientParams) error {
	_, err := q.db.Exec(ctx, addPrevActiveServerClient,
		arg.MessageID,
		arg.GuildID,
		arg.ChannelID,
		arg.Name,
		arg.Clan,
		arg.Team,
		arg.CountryID,
		arg.Score,
		arg.IsPlayer,
		arg.FlagAbbr,
		arg.FlagEmoji,
	)
	return err
}

const getPrevActiveServerClients = `-- name: GetPrevActiveServerClients :many
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
LIMIT 1
`

type GetPrevActiveServerClientsRow struct {
	MessageID int64  `db:"message_id"`
	GuildID   int64  `db:"guild_id"`
	ChannelID int64  `db:"channel_id"`
	Name      string `db:"name"`
	Clan      string `db:"clan"`
	Team      *int16 `db:"team"`
	CountryID int16  `db:"country_id"`
	Score     int32  `db:"score"`
	IsPlayer  bool   `db:"is_player"`
	FlagAbbr  string `db:"flag_abbr"`
	FlagEmoji string `db:"flag_emoji"`
}

func (q *Queries) GetPrevActiveServerClients(ctx context.Context, messageID int64) ([]GetPrevActiveServerClientsRow, error) {
	rows, err := q.db.Query(ctx, getPrevActiveServerClients, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPrevActiveServerClientsRow{}
	for rows.Next() {
		var i GetPrevActiveServerClientsRow
		if err := rows.Scan(
			&i.MessageID,
			&i.GuildID,
			&i.ChannelID,
			&i.Name,
			&i.Clan,
			&i.Team,
			&i.CountryID,
			&i.Score,
			&i.IsPlayer,
			&i.FlagAbbr,
			&i.FlagEmoji,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPrevActiveServers = `-- name: ListPrevActiveServers :many
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
ORDER BY guild_id ASC, channel_id ASC, message_id ASC
`

func (q *Queries) ListPrevActiveServers(ctx context.Context) ([]PrevActiveServer, error) {
	rows, err := q.db.Query(ctx, listPrevActiveServers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PrevActiveServer{}
	for rows.Next() {
		var i PrevActiveServer
		if err := rows.Scan(
			&i.MessageID,
			&i.GuildID,
			&i.ChannelID,
			&i.Timestamp,
			&i.Address,
			&i.Protocols,
			&i.Name,
			&i.Gametype,
			&i.Passworded,
			&i.Map,
			&i.MapSha256sum,
			&i.MapSize,
			&i.Version,
			&i.MaxClients,
			&i.MaxPlayers,
			&i.ScoreKind,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removePrevActiveServer = `-- name: RemovePrevActiveServer :exec
DELETE FROM prev_active_servers
WHERE message_id = $1
`

func (q *Queries) RemovePrevActiveServer(ctx context.Context, messageID int64) error {
	_, err := q.db.Exec(ctx, removePrevActiveServer, messageID)
	return err
}

const removePrevActiveServerClient = `-- name: RemovePrevActiveServerClient :exec
DELETE FROM prev_active_server_clients
WHERE message_id = $1
`

func (q *Queries) RemovePrevActiveServerClient(ctx context.Context, messageID int64) error {
	_, err := q.db.Exec(ctx, removePrevActiveServerClient, messageID)
	return err
}
