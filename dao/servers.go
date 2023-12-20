package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jxsl13/twstatus-bot/model"
)

func SetServers(ctx context.Context, tx *sql.Tx, servers []model.Server) error {
	_, err := tx.ExecContext(ctx, `DELETE FROM tw_servers`)
	if err != nil {
		return err
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

			if !KnownFlag(client.Country) {
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
