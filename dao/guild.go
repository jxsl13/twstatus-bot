package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
)

func GetGuild(ctx context.Context, conn Conn, guildID discord.GuildID) (guild model.Guild, err error) {
	rows, err := conn.QueryContext(ctx, `SELECT guild_id, description FROM guild WHERE guild_id = ?`,
		uint64(guildID),
	)
	if err != nil {
		return model.Guild{}, fmt.Errorf("failed to query guild: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return model.Guild{}, fmt.Errorf("%w: guild %d", ErrNotFound, guildID)
	}

	err = rows.Scan(
		&guild.ID,
		&guild.Description,
	)
	if err != nil {
		return model.Guild{}, fmt.Errorf("failed to scan guild: %w", err)
	}

	return guild, nil
}

func AddGuild(ctx context.Context, conn Conn, guild model.Guild) (err error) {
	_, err = conn.ExecContext(ctx, `INSERT INTO guild (guild_id, description) VALUES (?, ?)`,
		uint64(guild.ID),
		guild.Description,
	)

	if err != nil {
		if IsUniqueConstraintErr(err) {
			return fmt.Errorf("%w: guild %d", ErrAlreadyExists, guild.ID)
		}
		return fmt.Errorf("failed to insert guild %d: %w", guild.ID, err)
	}
	return nil
}

func ListGuilds(ctx context.Context, conn Conn) (guilds model.Guilds, err error) {
	rows, err := conn.QueryContext(ctx, `SELECT guild_id, description FROM guild`)
	if err != nil {
		return nil, fmt.Errorf("failed to query guilds: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var guild model.Guild
		err = rows.Scan(
			&guild.ID,
			&guild.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan guild: %w", err)
		}
		guilds = append(guilds, guild)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to iterate guilds: %w", err)
	}

	return guilds, nil
}

func RemoveGuild(ctx context.Context, db *sql.DB, guildID discord.GuildID) (guild model.Guild, err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return model.Guild{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}
	}()

	guild, err = GetGuild(ctx, tx, guildID)
	if err != nil {
		return model.Guild{}, fmt.Errorf("failed to get guild: %w", err)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM guild WHERE guild_id = ?`, uint64(guild.ID))
	if err != nil {
		return model.Guild{}, fmt.Errorf("failed to delete guild: %w", err)
	}
	return guild, nil
}
