package dao

import (
	"context"
	"fmt"

	"github.com/jxsl13/twstatus-bot/model"
	"github.com/jxsl13/twstatus-bot/sqlc"
)

func ListFlags(ctx context.Context, q *sqlc.Queries) (_ []model.Flag, err error) {
	fs, err := q.ListFlags(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query flags: %w", err)
	}
	flags := make([]model.Flag, 0, len(fs))
	for _, flag := range fs {
		flags = append(flags, model.Flag{
			ID:    flag.FlagID,
			Abbr:  flag.Abbr,
			Emoji: flag.Emoji,
		})
	}
	return flags, nil
}

func GetFlag(ctx context.Context, q *sqlc.Queries, flagId int64) (_ model.Flag, err error) {
	fs, err := q.GetFlag(ctx, flagId)
	if err != nil {
		return model.Flag{}, fmt.Errorf("failed to query flag: %w", err)
	}
	if len(fs) == 0 {
		return model.Flag{}, fmt.Errorf("%w: flag with id %d", ErrNotFound, flagId)
	}
	flag := fs[0]

	return model.Flag{
		ID:    flag.FlagID,
		Abbr:  flag.Abbr,
		Emoji: flag.Emoji,
	}, nil
}

func GetFlagByAbbr(ctx context.Context, q *sqlc.Queries, abbr string) (_ model.Flag, err error) {
	fs, err := q.GetFlagByAbbr(ctx, abbr)
	if err != nil {
		return model.Flag{}, fmt.Errorf("failed to query flag: %w", err)
	}
	if len(fs) == 0 {
		return model.Flag{}, fmt.Errorf("%w: flag with abbr %s", ErrNotFound, abbr)
	}
	flag := fs[0]
	return model.Flag{
		ID:    flag.FlagID,
		Abbr:  flag.Abbr,
		Emoji: flag.Emoji,
	}, nil
}
