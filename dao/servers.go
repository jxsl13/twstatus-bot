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

	stmt, err := tx.PrepareContext(ctx, `
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
	score_kind,
	clients
)
VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?);`)
	if err != nil {
		return fmt.Errorf("failed to prepare servers statement: %w", err)
	}

	for _, server := range servers {
		_, err = stmt.ExecContext(ctx,
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
			string(server.ClientsJSON()),
		)
		if err != nil {
			return fmt.Errorf("failed to insert server %s: %w", server.Address, err)
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
