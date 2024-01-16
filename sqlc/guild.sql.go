// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: guild.sql

package sqlc

import (
	"context"
)

const addGuild = `-- name: AddGuild :exec
INSERT INTO guilds (
    guild_id,
    description
) VALUES ($1, $2)
`

type AddGuildParams struct {
	GuildID     int64  `db:"guild_id"`
	Description string `db:"description"`
}

func (q *Queries) AddGuild(ctx context.Context, arg AddGuildParams) error {
	_, err := q.db.Exec(ctx, addGuild, arg.GuildID, arg.Description)
	return err
}

const getGuild = `-- name: GetGuild :many
SELECT guild_id, description
FROM guilds
WHERE guild_id = $1
LIMIT 1
`

func (q *Queries) GetGuild(ctx context.Context, guildID int64) ([]Guild, error) {
	rows, err := q.db.Query(ctx, getGuild, guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Guild{}
	for rows.Next() {
		var i Guild
		if err := rows.Scan(&i.GuildID, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listGuilds = `-- name: ListGuilds :many
SELECT guild_id, description
FROM guilds
ORDER BY guild_id ASC
`

func (q *Queries) ListGuilds(ctx context.Context) ([]Guild, error) {
	rows, err := q.db.Query(ctx, listGuilds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Guild{}
	for rows.Next() {
		var i Guild
		if err := rows.Scan(&i.GuildID, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeGuild = `-- name: RemoveGuild :exec
DELETE FROM guilds WHERE guild_id = $1
`

func (q *Queries) RemoveGuild(ctx context.Context, guildID int64) error {
	_, err := q.db.Exec(ctx, removeGuild, guildID)
	return err
}
