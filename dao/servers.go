package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jxsl13/twstatus-bot/model"
)

func ActiveServers(ctx context.Context, tx *sql.Tx) (servers map[model.Target]model.ServerStatus, err error) {
	servers, err = activeServers(ctx, tx)
	if err != nil {
		return nil, err
	}

	clients, err := activeClients(ctx, tx)
	if err != nil {
		return nil, err
	}

	for target := range servers {
		server := servers[target]
		server.Clients = clients[target]
		servers[target] = server
	}

	return servers, nil
}

func activeServers(ctx context.Context, conn Conn) (servers map[model.Target]model.ServerStatus, err error) {
	serverRows, err := conn.QueryContext(ctx, `
SELECT
	c.guild_id,
	c.channel_id,
	t.message_id,
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
JOIN tw_servers ts ON t.address = ts.address
WHERE c.running = 1
`)
	if err != nil {
		return nil, fmt.Errorf("failed to query active servers: %w", err)
	}
	defer serverRows.Close()

	servers = make(map[model.Target]model.ServerStatus)
	for serverRows.Next() {
		var (
			target    model.Target
			server    model.ServerStatus
			protocols []byte
		)
		err = serverRows.Scan(
			&target.GuildID,
			&target.ChannelID,
			&target.MessageID,

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
			return nil, fmt.Errorf("failed to scan server status: %w", err)
		}
		err = json.Unmarshal(protocols, &server.Protocols)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal protocols json: %w", err)
		}

		server.Clients = make(model.ClientStatusList, 0, 4)
		servers[target] = server
	}
	return servers, nil
}

func activeClients(ctx context.Context, conn Conn) (map[model.Target]model.ClientStatusList, error) {
	rows, err := conn.QueryContext(ctx, `
SELECT
	c.guild_id,
	c.channel_id,
	t.message_id,
	tsc.name,
	tsc.clan,
	tsc.country_id,
	(CASE WHEN tsc.score = -9999 THEN 9223372036854775807 ELSE tsc.score END) as score,
	tsc.is_player,
	f.abbr,
	(CASE WHEN fm.emoji NOT NULL THEN fm.emoji ELSE f.emoji END) as flag_emoji
FROM channels c
JOIN tracking t ON c.channel_id = t.channel_id
JOIN tw_server_clients tsc ON t.address = tsc.address
JOIN flags f ON tsc.country_id = f.flag_id
LEFT JOIN flag_mappings fm ON
	(
		t.guild_id = fm.guild_id AND
		t.channel_id = fm.channel_id AND
		tsc.country_id = fm.flag_id
	)
WHERE c.running = 1
ORDER BY c.guild_id, c.channel_id, t._rowid_, score DESC, tsc.name ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to query active players: %w", err)
	}
	defer rows.Close()

	result := make(map[model.Target]model.ClientStatusList)
	for rows.Next() {
		var (
			target model.Target
			client model.ClientStatus
		)
		err = rows.Scan(
			&target.GuildID,
			&target.ChannelID,
			&target.MessageID,

			&client.Name,
			&client.Clan,
			&client.Country,
			&client.Score,
			&client.IsPlayer,
			&client.FlagAbbr,
			&client.FlagEmoji,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan client status: %w", err)
		}
		result[target] = append(result[target], client)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over client status: %w", err)
	}

	return result, nil
}

func ListServers(ctx context.Context, conn Conn) (servers []model.Server, err error) {
	rows, err := conn.QueryContext(ctx, `
SELECT
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
	score_kind,
	clients
FROM tw_servers
ORDER BY address ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to query servers: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var server model.Server
		var (
			clientsJSON   []byte
			protocolsJSON []byte
		)

		err = rows.Scan(
			&server.Address,
			&protocolsJSON,
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
			&clientsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan server: %w", err)
		}

		err = json.Unmarshal(clientsJSON, &server.Clients)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal clients json: %w", err)
		}

		err = json.Unmarshal(protocolsJSON, &server.Protocols)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal protocols json: %w", err)
		}

		servers = append(servers, server)
	}

	return servers, nil
}

func SetServers(ctx context.Context, tx *sql.Tx, servers []model.Server) error {
	flags, err := ListFlags(ctx, tx)
	if err != nil {
		return err
	}

	knownFlags := make(map[int]bool)
	for _, flag := range flags {
		knownFlags[flag.ID] = true
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM tw_server_clients`)
	if err != nil {
		return fmt.Errorf("failed to delete server clients: %w", err)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM tw_servers`)
	if err != nil {
		return fmt.Errorf("failed to delete servers: %w", err)
	}

	serverStmt, err := tx.PrepareContext(ctx, `
INSERT INTO tw_servers (
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
)
VALUES (?,?,?,?,?,?,?,?,?,?,?,?);`)
	if err != nil {
		return fmt.Errorf("failed to prepare servers statement: %w", err)
	}

	clientStmt, err := tx.PrepareContext(ctx, `
REPLACE INTO tw_server_clients (
	address,
	name,
	clan,
	country_id,
	score,
	is_player
) VALUES (?,?,?,?,?,?);`)
	if err != nil {
		return fmt.Errorf("failed to prepare clients statement: %w", err)
	}

	var isPlayer int
	for _, server := range servers {
		_, err = serverStmt.ExecContext(ctx,
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
			return fmt.Errorf("failed to insert server %s: %w", server.Address, err)
		}

		for _, client := range server.Clients {
			// skip connecting players
			if client.IsConnecting() {
				continue
			}

			if !knownFlags[client.Country] {
				// set to known
				client.Country = -1
			}

			isPlayer = 0
			if client.IsPlayer {
				isPlayer = 1
			}
			_, err = clientStmt.ExecContext(ctx,
				server.Address,
				client.Name,
				client.Clan,
				client.Country,
				client.Score,
				isPlayer,
			)
			if err != nil {
				return fmt.Errorf("failed to insert client %s for address %s: %w", client.Name, server.Address, err)
			}
		}
	}
	return nil
}

func ExistsServer(ctx context.Context, conn Conn, address string) (found bool, err error) {
	rows, err := conn.QueryContext(ctx, `
SELECT
	address
FROM tw_servers
WHERE address = ?
LIMIT 1;`, address)

	if err != nil {
		return false, fmt.Errorf("failed to query server address: %s: %w", address, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return false, nil
	}
	err = rows.Err()
	if err != nil {
		return false, fmt.Errorf("failed to iterate over server addresses: %w", err)
	}
	return true, nil
}
