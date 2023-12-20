package dao

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
)

func GetChannel(ctx context.Context, conn Conn, guildId discord.GuildID, channelID discord.ChannelID) (model.Channel, error) {
	rows, err := conn.QueryContext(ctx, `
SELECT message_id, running
FROM channels
WHERE guild_id = ?
AND channel_id = ?
LIMIT 1;`,
		int64(guildId),
		int64(channelID),
	)
	if err != nil {
		return model.Channel{}, fmt.Errorf("failed to query channel: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return model.Channel{}, fmt.Errorf("%w: channel %d", ErrNotFound, channelID)
	}
	err = rows.Err()
	if err != nil {
		return model.Channel{}, fmt.Errorf("failed to iterate over channels: %w", err)
	}

	var channel model.Channel
	err = rows.Scan(
		&channel.MessageID,
		&channel.Running,
	)
	if err != nil {
		return model.Channel{}, fmt.Errorf("failed to scan channel: %w", err)
	}

	channel.GuildID = guildId
	channel.ID = channelID

	return channel, nil
}

func ListChannels(ctx context.Context, conn Conn, guildID discord.GuildID) (channels model.Channels, err error) {
	rows, err := conn.QueryContext(ctx, `
SELECT channel_id, message_id, running 
FROM channels 
WHERE guild_id = ?
ORDER BY channel_id ASC;`,
		int64(guildID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query channels: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var channel model.Channel
		err = rows.Scan(
			&channel.ID,
			&channel.MessageID,
			&channel.Running,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan channel: %w", err)
		}
		channel.GuildID = guildID
		channels = append(channels, channel)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to iterate over channels: %w", err)
	}

	return channels, nil
}

func AddChannel(ctx context.Context, tx *sql.Tx, channel model.Channel) (err error) {
	_, err = tx.ExecContext(ctx, `
INSERT INTO channels (guild_id, channel_id, message_id, running) 
VALUES (?, ?, ?, ?);`,
		channel.GuildID,
		channel.ID,
		channel.MessageID,
		channel.Running,
	)

	if err != nil {
		if IsUniqueConstraintErr(err) {
			return fmt.Errorf("%w: channel %d", ErrAlreadyExists, channel.ID)
		}
		return fmt.Errorf("failed to insert channel %d: %w", channel.ID, err)
	}

	err = insertDefaultFlags(ctx, tx, channel.ID)
	if err != nil {
		return fmt.Errorf("failed to insert default flags: %w", err)
	}

	return err
}

func insertDefaultFlags(ctx context.Context, tx *sql.Tx, channelID discord.ChannelID) (err error) {
	stmt, err := tx.PrepareContext(ctx, `
INSERT INTO flags (flag_id, channel_id, abbr, symbol) 
VALUES (?, ?, ?, ?);`)
	if err != nil {
		return fmt.Errorf("failed to add default flags: failed to prepare statement: %w", err)
	}

	for _, id := range flagKeys {
		vals := flags[id]
		abbr, flag := vals[0], vals[1]
		_, err = stmt.ExecContext(ctx, id, int64(channelID), abbr, flag)
		if err != nil {
			return fmt.Errorf("failed to insert default flag %s: %w", abbr, err)
		}
	}
	return nil
}

func RemoveChannel(ctx context.Context, tx *sql.Tx, guildID discord.GuildID, channelID discord.ChannelID) (channel model.Channel, err error) {

	channel, err = GetChannel(ctx, tx, guildID, channelID)
	if err != nil {
		return model.Channel{}, err
	}

	_, err = tx.ExecContext(ctx, `
DELETE FROM channels
WHERE guild_id = ?
AND channel_id = ?;`,
		channel.GuildID,
		channel.ID,
	)
	if err != nil {
		return model.Channel{}, fmt.Errorf("failed to delete channel %d: %w", channelID, err)
	}

	return channel, nil
}
