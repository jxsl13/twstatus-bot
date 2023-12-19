package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
)

func GetChannel(ctx context.Context, conn Conn, guildId discord.GuildID, channelID discord.ChannelID) (model.Channel, error) {
	rows, err := conn.QueryContext(ctx, `SELECT message_id, running FROM channel WHERE guild_id = ? AND channel_id = ?;`,
		uint64(guildId),
		uint64(channelID),
	)
	if err != nil {
		return model.Channel{}, fmt.Errorf("failed to query channel: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return model.Channel{}, fmt.Errorf("%w: channel %d", ErrNotFound, channelID)
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
	rows, err := conn.QueryContext(ctx, `SELECT channel_id, message_id, running FROM channel WHERE guild_id = ?;`,
		uint64(guildID),
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

	return channels, nil
}

func AddChannel(ctx context.Context, conn Conn, channel model.Channel) (err error) {
	_, err = conn.ExecContext(ctx, `INSERT INTO channel (guild_id, channel_id, message_id, running) VALUES (?, ?, ?, ?);`,
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

	err = insertDefaultFlags(ctx, conn, channel.ID)
	if err != nil {
		return fmt.Errorf("failed to insert default flags: %w", err)
	}

	return err
}

func insertDefaultFlags(ctx context.Context, conn Conn, channelID discord.ChannelID) (err error) {
	for id, vals := range flags {
		abbr, flag := vals[0], vals[1]
		_, err = conn.ExecContext(ctx, `INSERT INTO flag (flag_id, channel_id, abbr, symbol) VALUES (?, ?, ?, ?);`,
			id,
			uint64(channelID),
			abbr,
			flag,
		)
		if err != nil {
			return fmt.Errorf("failed to insert default flag %s: %w", abbr, err)
		}
	}
	return nil
}

func RemoveChannel(ctx context.Context, db *sql.DB, guildID discord.GuildID, channelID discord.ChannelID) (channel model.Channel, err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return model.Channel{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}
	}()

	channel, err = GetChannel(ctx, tx, guildID, channelID)
	if err != nil {
		return model.Channel{}, err
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM channel WHERE guild_id = ? AND channel_id = ?`,
		channel.GuildID,
		channel.ID,
	)
	if err != nil {
		return model.Channel{}, fmt.Errorf("failed to delete channel %d: %w", channelID, err)
	}

	return channel, nil
}
