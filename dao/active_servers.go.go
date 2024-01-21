package dao

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

func ChangedServers(ctx context.Context, q *sqlc.Queries) (_ map[model.MessageTarget]model.ChangedServerStatus, err error) {
	previousServers, err := PrevActiveServers(ctx, q)
	if err != nil {
		return nil, err
	}

	currentServers, err := ActiveServers(ctx, q)
	if err != nil {
		return nil, err
	}

	changedServers := make(map[model.MessageTarget]model.ChangedServerStatus, 64)

	// removed servers
	for target := range previousServers {
		if empty, ok := currentServers[target]; !ok {
			changedServers[target] = model.ChangedServerStatus{
				Target:  target,
				Prev:    previousServers[target],
				Curr:    empty,
				Offline: true,
			}
		}
	}

	// to add
	added := make(map[model.MessageTarget]model.ServerStatus, 64)
	for target, server := range currentServers {
		if prev, ok := previousServers[target]; ok {
			// found in prev -> check if changed
			if !server.Equals(prev) {
				changedServers[target] = model.ChangedServerStatus{
					Target: target,
					Prev:   prev,
					Curr:   server,
				}
				added[target] = server
			}
		} else {
			// not found in prev -> new server
			changedServers[target] = model.ChangedServerStatus{
				Target: target,
				Prev:   model.ServerStatus{},
				Curr:   server,
			}
			added[target] = server
		}
	}

	var messageIDs []discord.MessageID
	for target := range changedServers {
		messageIDs = append(messageIDs, target.MessageID)
	}
	err = removePrevActiveServers(ctx, q, messageIDs)
	if err != nil {
		return nil, err
	}

	err = removePrevActiveClients(ctx, q, messageIDs)
	if err != nil {
		return nil, err
	}

	err = addPrevActiveServers(ctx, q, added)
	if err != nil {
		return nil, err
	}

	err = addPrevActiveClients(ctx, q, added)
	if err != nil {
		return nil, err
	}

	changedServers, err = GetTargetListNotifications(ctx, q, changedServers)
	if err != nil {
		return nil, err
	}

	return changedServers, nil
}

func ActiveServers(ctx context.Context, q *sqlc.Queries) (servers map[model.MessageTarget]model.ServerStatus, err error) {

	servers, err = activeServers(ctx, q)
	if err != nil {
		return nil, err
	}

	servers, err = activeClients(ctx, q, servers)
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func activeServers(ctx context.Context, q *sqlc.Queries) (servers map[model.MessageTarget]model.ServerStatus, err error) {
	ltsr, err := q.ListTrackedServers(ctx)
	if err != nil {
		return nil, err
	}
	servers = make(map[model.MessageTarget]model.ServerStatus, len(ltsr))
	for _, row := range ltsr {
		target := model.MessageTarget{
			ChannelTarget: model.ChannelTarget{
				GuildID:   discord.GuildID(row.GuildID),
				ChannelID: discord.ChannelID(row.ChannelID),
			},
			MessageID: discord.MessageID(row.MessageID),
		}
		server := model.ServerStatus{
			Timestamp:    row.Timestamp.Time,
			Address:      row.Address,
			Name:         row.Name,
			Gametype:     row.Gametype,
			Passworded:   row.Passworded,
			Map:          row.Map,
			MapSha256Sum: row.MapSha256sum,
			MapSize:      row.MapSize,
			Version:      row.Version,
			MaxClients:   row.MaxClients,
			MaxPlayers:   row.MaxPlayers,
			ScoreKind:    row.ScoreKind,
		}

		err = server.ProtocolsFromJSON([]byte(row.Protocols))
		if err != nil {
			return nil, err
		}
		servers[target] = server

	}
	return servers, nil
}

func activeClients(ctx context.Context, q *sqlc.Queries, servers map[model.MessageTarget]model.ServerStatus) (_ map[model.MessageTarget]model.ServerStatus, err error) {
	if len(servers) == 0 {
		return map[model.MessageTarget]model.ServerStatus{}, nil
	}

	ltscr, err := q.ListTrackedServerClients(ctx)
	if err != nil {
		return nil, err
	}

	for _, row := range ltscr {
		var (
			target = model.MessageTarget{
				ChannelTarget: model.ChannelTarget{
					GuildID:   discord.GuildID(row.GuildID),
					ChannelID: discord.ChannelID(row.ChannelID),
				},
				MessageID: discord.MessageID(row.MessageID),
			}
			client = model.ClientStatus{
				Name:      row.Name,
				Clan:      row.Clan,
				Country:   row.CountryID,
				Score:     row.Score, // TODO: fix this
				IsPlayer:  row.IsPlayer,
				Team:      row.Team,
				FlagAbbr:  row.Abbr,
				FlagEmoji: row.FlagEmoji, // TODO: fix this
			}
		)

		server := servers[target]
		server.AddClientStatus(client)
		servers[target] = server
	}

	return servers, nil
}

func SetServers(ctx context.Context, q *sqlc.Queries, servers model.ServerList) error {
	flags, err := q.ListFlags(ctx)
	if err != nil {
		return fmt.Errorf("failed to list flags: %w", err)
	}

	knownFlags := make(map[int16]bool)
	for _, flag := range flags {
		knownFlags[flag.FlagID] = true
	}

	err = q.DeleteActiveServerClients(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete active server clients: %w", err)
	}

	err = q.DeleteActiveServers(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete active servers: %w", err)
	}

	ss, cs := servers.ToSQLC(knownFlags)
	i, err := q.InsertActiveServers(ctx, ss)
	if err != nil {
		return fmt.Errorf("failed to insert servers: %w", err)
	}
	if i != int64(len(ss)) {
		return fmt.Errorf("failed to insert all servers: %d/%d", i, len(ss))
	}

	i, err = q.InsertActiveServerClients(ctx, cs)
	if err != nil {
		return fmt.Errorf("failed to insert clients: %w", err)
	}
	if i != int64(len(cs)) {
		return fmt.Errorf("failed to insert all clients: %d/%d", i, len(cs))
	}

	return nil
}

func ExistsServer(ctx context.Context, q *sqlc.Queries, address string) (found bool, err error) {
	addr, err := q.ExistsServer(ctx, address)
	if err != nil {
		return false, err
	}

	if len(addr) == 0 {
		return false, nil
	}

	return addr[0] == address, nil
}
