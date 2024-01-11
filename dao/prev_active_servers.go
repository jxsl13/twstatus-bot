package dao

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
)

func PrevActiveServers(ctx context.Context, conn Conn) (servers map[model.MessageTarget]model.ServerStatus, err error) {
	servers, err = prevActiveServers(ctx, conn)
	if err != nil {
		return nil, err
	}

	clients, err := prevActiveClients(ctx, conn, servers)
	if err != nil {
		return nil, err
	}

	for target := range servers {
		server := servers[target]
		for _, client := range clients[target] {
			server.AddClientStatus(client)
		}
		servers[target] = server
	}

	return servers, nil
}

func prevActiveServers(ctx context.Context, conn Conn) (servers map[model.MessageTarget]model.ServerStatus, err error) {
	rows, err := conn.QueryContext(ctx, `
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
ORDER BY guild_id ASC, channel_id ASC, message_id ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to query previous active servers: %w", err)
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	servers = make(map[model.MessageTarget]model.ServerStatus, 2048)
	for rows.Next() {
		var (
			target    model.MessageTarget
			server    model.ServerStatus
			protocols []byte
			timestamp int64
		)
		err = rows.Scan(
			&target.MessageID,
			&target.GuildID,
			&target.ChannelID,

			&timestamp,
			&server.Address,
			&protocols,
			&server.Name,
			&server.Gametype,
			&server.Passworded,
			&server.Map,
			&server.MapSha256Sum,
			&server.MapSize,
			&server.Version,
			&server.MaxClients,
			&server.MaxPlayers,
			&server.ScoreKind,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan previous server status: %w", err)
		}
		err = server.ProtocolsFromJSON(protocols)
		if err != nil {
			return nil, err
		}
		server.Timestamp = time.UnixMilli(timestamp)
		servers[target] = server
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over previous server status: %w", err)
	}

	return servers, nil
}

func prevActiveClients(
	ctx context.Context,
	conn Conn,
	servers map[model.MessageTarget]model.ServerStatus,
) (
	_ map[model.MessageTarget]model.ClientStatusList,
	err error,
) {
	if len(servers) == 0 {
		return map[model.MessageTarget]model.ClientStatusList{}, nil
	}

	args := make([]any, 0, len(servers))
	for target := range servers {
		args = append(args, target.MessageID)
	}

	rows, err := conn.QueryContext(ctx,
		fmt.Sprintf(`
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
WHERE message_id IN (%s)
ORDER BY id ASC`,
			strings.Repeat("?,", len(args)-1)+"?"),
		args...,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query previous active clients: %w", err)
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	result := make(map[model.MessageTarget]model.ClientStatusList)
	for rows.Next() {
		var (
			target model.MessageTarget
			client model.ClientStatus
		)
		err = rows.Scan(
			&target.MessageID,
			&target.GuildID,
			&target.ChannelID,

			&client.Name,
			&client.Clan,
			&client.Team,
			&client.Country,
			&client.Score,
			&client.IsPlayer,
			&client.FlagAbbr,
			&client.FlagEmoji,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan previous client status: %w", err)
		}
		result[target] = append(result[target], client)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over previous client status: %w", err)
	}
	return result, nil
}

func addPrevActiveServers(
	ctx context.Context,
	conn Conn,
	servers map[model.MessageTarget]model.ServerStatus,
) (err error) {
	stmt, err := conn.PrepareContext(ctx, `
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
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() {
		err = errors.Join(err, stmt.Close())
	}()
	for target, server := range servers {
		_, err = stmt.ExecContext(ctx,
			target.MessageID,
			target.GuildID,
			target.ChannelID,

			server.Timestamp.UnixMilli(),
			server.Address,
			string(server.ProtocolsJSON()),
			server.Name,
			server.Gametype,
			server.Passworded,
			server.Map,
			server.MapSha256Sum,
			server.MapSize,
			server.Version,
			server.MaxClients,
			server.MaxPlayers,
			server.ScoreKind,
		)
		if err != nil {
			return fmt.Errorf("failed to insert previous server status: %#v -> %#v: %w", target, server, err)
		}
	}
	return nil
}

func removePrevActiveServers(ctx context.Context, conn Conn, messageIds []discord.MessageID) (err error) {
	if len(messageIds) == 0 {
		return nil
	}

	args := make([]any, 0, len(messageIds))
	for _, id := range messageIds {
		args = append(args, id)
	}

	_, err = conn.ExecContext(
		ctx,
		fmt.Sprintf(
			`DELETE FROM prev_active_servers WHERE message_id IN (%s);`,
			strings.Repeat("?,", len(messageIds)-1)+"?",
		),
		args...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete previous active servers: %w", err)
	}
	return nil
}

func addPrevActiveClients(ctx context.Context, conn Conn, servers map[model.MessageTarget]model.ServerStatus) (err error) {
	stmt, err := conn.PrepareContext(ctx, `
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
) VALUES (?,?,?,?,?,?,?,?,?,?,?);`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() {
		err = errors.Join(err, stmt.Close())
	}()

	for target, server := range servers {
		for _, client := range server.Clients {
			_, err = stmt.ExecContext(ctx,
				target.MessageID,
				target.GuildID,
				target.ChannelID,

				client.Name,
				client.Clan,
				client.Team,
				client.Country,
				client.Score,
				client.IsPlayer,
				client.FlagAbbr,
				client.FlagEmoji,
			)
			if err != nil {
				return fmt.Errorf("failed to insert previous client status: %#v -> %#v: %w", target, client, err)
			}
		}
	}

	return nil
}

func removePrevActiveClients(ctx context.Context, conn Conn, messageIds []discord.MessageID) (err error) {
	if len(messageIds) == 0 {
		return nil
	}

	args := make([]any, 0, len(messageIds))
	for _, id := range messageIds {
		args = append(args, id)
	}

	_, err = conn.ExecContext(
		ctx,
		fmt.Sprintf(
			`DELETE FROM prev_active_server_clients WHERE message_id IN (%s);`,
			strings.Repeat("?,", len(messageIds)-1)+"?",
		),
		args...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete previous active clients: %w", err)
	}
	return nil
}
