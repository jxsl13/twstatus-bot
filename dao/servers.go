package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"

	"github.com/jxsl13/twstatus-bot/model"
)

var activeServersSQL = fmt.Sprintf(`
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
	ts.score_kind,
	tsc.name,
	tsc.clan,
	tsc.country_id,
	(CASE WHEN tsc.score = -9999 THEN %d ELSE tsc.score END) as score,
	tsc.is_player,
	f.abbr,
	(CASE WHEN fm.emoji NOT NULL THEN fm.emoji ELSE f.emoji END) as flag_emoji
FROM channels c
JOIN tracking t ON c.channel_id = t.channel_id
JOIN tw_servers ts ON t.address = ts.address
JOIN tw_server_clients tsc ON ts.address = tsc.address
JOIN flags f ON tsc.country_id = f.flag_id
LEFT JOIN flag_mappings fm ON
	(
		t.guild_id = fm.guild_id AND
		t.channel_id = fm.channel_id AND
		tsc.country_id = fm.flag_id
	)
WHERE c.running = 1
ORDER BY c.guild_id, c.channel_id, t._rowid_, score DESC, tsc.name ASC`, math.MaxInt)

func ActiveServers(ctx context.Context, conn Conn) (servers map[model.Target]model.ServerStatus, err error) {
	rows, err := conn.QueryContext(ctx, activeServersSQL) // stay architecture independent
	if err != nil {
		return nil, fmt.Errorf("failed to query active servers: %w", err)
	}
	defer rows.Close()

	servers = make(map[model.Target]model.ServerStatus)
	for rows.Next() {
		var (
			target    model.Target
			server    model.ServerStatus
			client    model.ClientStatus
			protocols []byte
		)

		err = rows.Scan(
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

			&client.Name,
			&client.Clan,
			&client.Country,
			&client.Score,
			&client.IsPlayer,
			&client.FlagAbbr,
			&client.FlagEmoji,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan server status: %w", err)
		}

		s := servers[target]
		s.Address = server.Address
		err = json.Unmarshal(protocols, &s.Protocols)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal protocols json: %w", err)
		}
		s.Name = server.Name
		s.Gametype = server.Gametype
		s.Passworded = server.Passworded
		s.Map = server.Map
		s.MapSha256Sum = server.MapSha256Sum
		s.MapSize = server.MapSize
		s.Version = server.Version
		s.MaxClients = server.MaxClients
		s.MaxPlayers = server.MaxPlayers
		s.ScoreKind = server.ScoreKind
		s.Clients = append(s.Clients, client)

		servers[target] = s
	}
	return servers, nil
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
