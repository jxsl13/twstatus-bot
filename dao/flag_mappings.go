package dao

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

func ListFlagMappings(
	ctx context.Context,
	q *sqlc.Queries,
	guildId discord.GuildID,
	channelId discord.ChannelID,
) (
	_ model.FlagMappings,
	err error,
) {

	fms, err := q.ListFlagMappings(ctx, sqlc.ListFlagMappingsParams{
		GuildID:   int64(guildId),
		ChannelID: int64(channelId),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query flag mappings: %w", err)
	}
	result := make(model.FlagMappings, 0, len(fms))
	for _, fm := range fms {
		result = append(result, model.FlagMapping{
			GuildID:   guildId,
			ChannelID: channelId,
			FlagID:    fm.FlagID,
			Emoji:     fm.Emoji,
			Abbr:      fm.Abbr,
		})
	}

	return result, nil
}

func AddFlagMapping(ctx context.Context, q *sqlc.Queries, mapping model.FlagMapping) (err error) {
	return q.AddFlagMapping(ctx, sqlc.AddFlagMappingParams{
		GuildID:   int64(mapping.GuildID),
		ChannelID: int64(mapping.ChannelID),
		FlagID:    mapping.FlagID,
		Emoji:     mapping.Emoji,
	})
}

func GetFlagMapping(
	ctx context.Context,
	q *sqlc.Queries,
	guildId discord.GuildID,
	channelId discord.ChannelID,
	flagId int16,
) (
	_ model.FlagMapping,
	err error,
) {

	fm, err := q.GetFlagMapping(ctx, sqlc.GetFlagMappingParams{
		GuildID:   int64(guildId),
		ChannelID: int64(channelId),
		FlagID:    flagId,
	})
	if err != nil {
		return model.FlagMapping{}, fmt.Errorf("failed to query flag mapping: %w", err)
	}
	if len(fm) == 0 {
		return model.FlagMapping{}, fmt.Errorf("%w: flag mapping", ErrNotFound)
	}
	mapping := fm[0]
	return model.FlagMapping{
		GuildID:   guildId,
		ChannelID: channelId,
		FlagID:    mapping.FlagID,
		Emoji:     mapping.Emoji,
		Abbr:      mapping.Abbr,
	}, nil
}

func RemoveFlagMapping(
	ctx context.Context,
	q *sqlc.Queries,
	guildId discord.GuildID,
	channelId discord.ChannelID,
	abbr string,
) (err error) {

	flag, err := q.GetFlagByAbbr(ctx, abbr)
	if err != nil {
		return err
	}
	if len(flag) == 0 {
		return fmt.Errorf("%w: flag %s", ErrNotFound, abbr)
	}
	f := flag[0]
	return q.RemoveFlagMapping(ctx, sqlc.RemoveFlagMappingParams{
		GuildID:   int64(guildId),
		ChannelID: int64(channelId),
		FlagID:    f.FlagID,
	})
}
