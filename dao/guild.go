package dao

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

func ListGuilds(ctx context.Context, q *sqlc.Queries) (guilds model.Guilds, err error) {
	gs, err := q.ListGuilds(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query guilds: %w", err)
	}
	guilds = make(model.Guilds, 0, len(gs))
	for _, g := range gs {
		guilds = append(guilds, model.Guild{
			ID:          discord.GuildID(g.GuildID),
			Description: g.Description,
		})
	}
	return guilds, nil
}

func AddGuild(ctx context.Context, q *sqlc.Queries, guild model.Guild) (err error) {
	err = q.AddGuild(ctx, guild.ToSQLC())
	if err != nil {
		if IsUniqueConstraintErr(err) {
			return fmt.Errorf("%w: guild %d", ErrAlreadyExists, guild.ID)
		}
		return fmt.Errorf("failed to insert guild %d: %w", guild.ID, err)
	}
	return nil
}

func GetGuild(ctx context.Context, q *sqlc.Queries, guildID discord.GuildID) (guild model.Guild, err error) {
	gs, err := q.GetGuild(ctx, int64(guildID))
	if err != nil {
		return model.Guild{}, fmt.Errorf("failed to query guild: %w", err)
	}
	if len(gs) == 0 {
		return model.Guild{}, fmt.Errorf("%w: guild %d", ErrNotFound, guildID)
	}
	g := gs[0]
	return model.Guild{
		ID:          discord.GuildID(g.GuildID),
		Description: g.Description,
	}, nil
}

func RemoveGuild(ctx context.Context, q *sqlc.Queries, guildID discord.GuildID) (guild model.Guild, err error) {

	guild, err = GetGuild(ctx, q, guildID)
	if err != nil {
		return model.Guild{}, fmt.Errorf("failed to get guild: %w", err)
	}

	err = q.RemoveGuild(ctx, int64(guildID))
	if err != nil {
		return model.Guild{}, fmt.Errorf("failed to remove guild: %w", err)
	}
	return guild, err
}
