package dao

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
)

func StartChannel(ctx context.Context, tx *sql.Tx, guildID discord.GuildID, channelID discord.ChannelID) (c model.Channel, err error) {
	channel, err := GetChannel(ctx, tx, guildID, channelID)
	if err != nil {
		return c, err
	}

	if channel.Running {
		return c, fmt.Errorf("channel %s is already active", channel)
	}

	_, err = tx.ExecContext(ctx, `
UPDATE channels
SET running = 1
WHERE guild_id = ?
AND channel_id = ?;`, guildID, channelID)
	if err != nil {
		return c, fmt.Errorf("failed to start channel %s: %w", channel, err)
	}
	return channel, nil
}

func StopChannel(ctx context.Context, tx *sql.Tx, guildID discord.GuildID, channelID discord.ChannelID) (c model.Channel, err error) {
	channel, err := GetChannel(ctx, tx, guildID, channelID)
	if err != nil {
		return c, err
	}

	if !channel.Running {
		return c, fmt.Errorf("channel %s is already inactive", channel)
	}

	_, err = tx.ExecContext(ctx, `
UPDATE channels
SET running = 0
WHERE guild_id = ?
AND channel_id = ?;`, guildID, channelID)
	if err != nil {
		return c, fmt.Errorf("failed to stop channel %s: %w", channel, err)
	}
	return channel, nil
}
