package dao

import (
	"context"
	"fmt"

	"github.com/jxsl13/twstatus-bot/model"
)

func (dao *DAO) ListFlags(ctx context.Context) (_ []model.Flag, err error) {
	fs, err := dao.q.ListFlags(ctx)
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

func (dao *DAO) GetFlag(ctx context.Context, flagId int16) (_ model.Flag, err error) {
	fs, err := dao.q.GetFlag(ctx, flagId)
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

func (dao *DAO) GetFlagByAbbr(ctx context.Context, abbr string) (_ model.Flag, err error) {
	fs, err := dao.q.GetFlagByAbbr(ctx, abbr)
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
