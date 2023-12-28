package dao

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
)

func GetFlagMapping(ctx context.Context, conn Conn, guildId discord.GuildID, channelId discord.ChannelID, flagId int) (model.FlagMapping, error) {
	rows, err := conn.QueryContext(ctx, `
SELECT
	m.flag_id,
	m.emoji,
	f.abbr
FROM flag_mappings m
JOIN flags f ON m.flag_id = f.flag_id
WHERE m.guild_id = ?
AND m.channel_id = ?
AND m.flag_id = ? LIMIT 1;`, guildId, channelId, flagId)
	if err != nil {
		return model.FlagMapping{}, fmt.Errorf("failed to query flag mapping: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return model.FlagMapping{}, fmt.Errorf("%w: flag mapping", ErrNotFound)
	}
	err = rows.Err()
	if err != nil {
		return model.FlagMapping{}, fmt.Errorf("failed to iterate over flag mapping: %w", err)
	}

	var mapping model.FlagMapping = model.FlagMapping{
		GuildID:   guildId,
		ChannelID: channelId,
	}
	err = rows.Scan(
		&mapping.FlagID,
		&mapping.Emoji,
		&mapping.Abbr,
	)
	if err != nil {
		return model.FlagMapping{}, fmt.Errorf("failed to scan flag mapping: %w", err)
	}
	return mapping, nil
}

func AddFlagMapping(ctx context.Context, conn Conn, mapping model.FlagMapping) error {
	_, err := conn.ExecContext(ctx, `
REPLACE INTO flag_mappings (guild_id, channel_id, flag_id, emoji)
VALUES (?, ?, ?, ?);`,
		mapping.GuildID,
		mapping.ChannelID,
		mapping.FlagID,
		mapping.Emoji,
	)
	if err != nil {
		return fmt.Errorf("failed to insert flag mapping: %w", err)
	}
	return nil
}

func ListFlagMappings(ctx context.Context, conn Conn, guildId discord.GuildID, channelId discord.ChannelID) (model.FlagMappings, error) {
	rows, err := conn.QueryContext(ctx, `
SELECT
    m.flag_id,
	m.emoji,
	f.abbr
FROM flag_mappings m
JOIN flags f ON m.flag_id = f.flag_id
WHERE m.guild_id = ?
AND m.channel_id = ?;`, guildId, channelId)
	if err != nil {
		return nil, fmt.Errorf("failed to query flag mappings: %w", err)
	}
	defer rows.Close()
	result := model.FlagMappings{}
	for rows.Next() {
		var mapping model.FlagMapping = model.FlagMapping{
			GuildID:   guildId,
			ChannelID: channelId,
		}

		err = rows.Scan(
			&mapping.FlagID,
			&mapping.Emoji,
			&mapping.Abbr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan flag mapping: %w", err)
		}
		result = append(result, mapping)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over flag mappings: %w", err)
	}
	return result, nil
}

func RemoveFlagMapping(ctx context.Context, tx *sql.Tx, guildId discord.GuildID, channelId discord.ChannelID, abbr string) error {

	flag, err := GetFlagByAbbr(ctx, tx, abbr)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
DELETE FROM flag_mappings
WHERE guild_id = ?
AND channel_id = ?
AND flag_id = ?;`, guildId, channelId, flag.ID)
	if err != nil {
		return fmt.Errorf("failed to remove flag mapping: %w", err)
	}
	return nil
}
