package dao

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

func PrevActiveServers(ctx context.Context, q *sqlc.Queries) (servers map[model.MessageTarget]model.ServerStatus, err error) {
	servers, err = prevActiveServers(ctx, q)
	if err != nil {
		return nil, err
	}

	// enrich with clients
	servers, err = prevActiveClients(ctx, q, servers)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func prevActiveServers(ctx context.Context, q *sqlc.Queries) (servers map[model.MessageTarget]model.ServerStatus, err error) {
	pas, err := q.ListPrevActiveServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous active servers: %w", err)
	}

	servers = make(map[model.MessageTarget]model.ServerStatus, len(pas))
	for _, s := range pas {
		var (
			target = model.MessageTarget{
				ChannelTarget: model.ChannelTarget{
					GuildID:   discord.GuildID(s.GuildID),
					ChannelID: discord.ChannelID(s.ChannelID),
				},
				MessageID: discord.MessageID(s.MessageID),
			}
			server = model.ServerStatus{
				Timestamp:    s.Timestamp.Time,
				Address:      s.Address,
				Name:         s.Name,
				Gametype:     s.Gametype,
				Passworded:   s.Passworded,
				Map:          s.Map,
				MapSha256Sum: s.MapSha256sum,
				MapSize:      s.MapSize,
				Version:      s.Version,
				MaxClients:   s.MaxClients,
				MaxPlayers:   s.MaxPlayers,
				ScoreKind:    s.ScoreKind,
			}
		)

		err = server.ProtocolsFromJSON(s.Protocols)
		if err != nil {
			return nil, err
		}
		servers[target] = server
	}

	return servers, nil
}

func prevActiveClients(
	ctx context.Context,
	q *sqlc.Queries,
	servers map[model.MessageTarget]model.ServerStatus,
) (
	_ map[model.MessageTarget]model.ServerStatus,
	err error,
) {
	if len(servers) == 0 {
		return map[model.MessageTarget]model.ServerStatus{}, nil
	}

	for t, s := range servers {
		var (
			target model.MessageTarget
			client model.ClientStatus
		)
		rows, err := q.GetPrevActiveServerClients(ctx, int64(t.MessageID))
		if err != nil {
			return nil, fmt.Errorf("failed to get previous active clients: %w", err)
		}
		for _, row := range rows {
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
				Team:      row.Team,
				Country:   row.CountryID,
				Score:     row.Score,
				IsPlayer:  row.IsPlayer,
				FlagAbbr:  row.FlagAbbr,
				FlagEmoji: row.FlagEmoji,
			}
			s.AddClientStatus(client)
			servers[target] = s
		}
	}
	return servers, nil
}

func addPrevActiveServers(
	ctx context.Context,
	q *sqlc.Queries,
	servers map[model.MessageTarget]model.ServerStatus,
) (err error) {

	for t, s := range servers {
		err = q.AddPrevActiveServer(ctx, sqlc.AddPrevActiveServerParams{
			MessageID: int64(t.MessageID),
			GuildID:   int64(t.GuildID),
			ChannelID: int64(t.ChannelID),
			Timestamp: pgtype.Timestamptz{
				Time:  s.Timestamp,
				Valid: true,
			},
			Address:      s.Address,
			Protocols:    s.ProtocolsJSON(),
			Name:         s.Name,
			Gametype:     s.Gametype,
			Passworded:   s.Passworded,
			Map:          s.Map,
			MapSha256sum: s.MapSha256Sum,
			MapSize:      s.MapSize,
			Version:      s.Version,
			MaxClients:   s.MaxClients,
			MaxPlayers:   s.MaxPlayers,
			ScoreKind:    s.ScoreKind,
		})
		if err != nil {
			return fmt.Errorf("failed to insert previous server status: %#v -> %#v: %w", t, s, err)
		}
	}

	return nil
}

func removePrevActiveServers(ctx context.Context, q *sqlc.Queries, messageIds []discord.MessageID) (err error) {
	if len(messageIds) == 0 {
		return nil
	}

	for _, id := range messageIds {
		err = q.RemovePrevActiveServer(ctx, int64(id))
		if err != nil {
			return fmt.Errorf("failed to delete previous active server: %w", err)
		}
	}
	return nil
}

func addPrevActiveClients(ctx context.Context, q *sqlc.Queries, servers map[model.MessageTarget]model.ServerStatus) (err error) {

	for target, server := range servers {
		for _, client := range server.Clients {
			err = q.AddPrevActiveServerClient(ctx, sqlc.AddPrevActiveServerClientParams{
				MessageID: int64(target.MessageID),
				GuildID:   int64(target.GuildID),
				ChannelID: int64(target.ChannelID),
				Name:      client.Name,
				Clan:      client.Clan,
				Team:      client.Team,
				CountryID: client.Country,
				Score:     client.Score,
				IsPlayer:  client.IsPlayer,
				FlagAbbr:  client.FlagAbbr,
				FlagEmoji: client.FlagEmoji,
			})
			if err != nil {
				return fmt.Errorf("failed to insert previous server client: %#v -> %#v: %w", target, client, err)
			}
		}
	}

	return nil
}

func removePrevActiveClients(ctx context.Context, q *sqlc.Queries, messageIds []discord.MessageID) (err error) {
	if len(messageIds) == 0 {
		return nil
	}

	for _, id := range messageIds {
		err = q.RemovePrevActiveServerClient(ctx, int64(id))
		if err != nil {
			return err
		}
	}
	return nil
}
