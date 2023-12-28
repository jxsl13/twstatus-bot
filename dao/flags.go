package dao

import (
	"context"
	"fmt"

	"github.com/jxsl13/twstatus-bot/model"
)

func GetFlagList(ctx context.Context, conn Conn) ([]model.Flag, error) {
	rows, err := conn.QueryContext(ctx, `SELECT flag_id, abbr, emoji FROM flags ORDER BY abbr;`)
	if err != nil {
		return nil, fmt.Errorf("failed to query flags: %w", err)
	}
	defer rows.Close()

	flags := []model.Flag{}
	for rows.Next() {
		var flag model.Flag
		err = rows.Scan(
			&flag.ID,
			&flag.Abbr,
			&flag.Emoji,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan flag: %w", err)
		}
		flags = append(flags, flag)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over flags: %w", err)
	}
	return flags, nil
}

func GetFlag(ctx context.Context, conn Conn, flagId int) (model.Flag, error) {
	rows, err := conn.QueryContext(ctx, `
SELECT flag_id, abbr, emoji
FROM flags
WHERE flag_id = ? LIMIT 1;`, flagId)
	if err != nil {
		return model.Flag{}, fmt.Errorf("failed to query flag: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return model.Flag{}, fmt.Errorf("%w: flag", ErrNotFound)
	}
	err = rows.Err()
	if err != nil {
		return model.Flag{}, fmt.Errorf("failed to iterate over flag: %w", err)
	}

	var flag model.Flag
	err = rows.Scan(
		&flag.ID,
		&flag.Abbr,
		&flag.Emoji,
	)
	if err != nil {
		return model.Flag{}, fmt.Errorf("failed to scan flag: %w", err)
	}
	return flag, nil
}

func GetFlagByAbbr(ctx context.Context, conn Conn, abbr string) (model.Flag, error) {
	rows, err := conn.QueryContext(ctx, `
SELECT flag_id, abbr, emoji
FROM flags
WHERE abbr = ? LIMIT 1;`, abbr)
	if err != nil {
		return model.Flag{}, fmt.Errorf("failed to query flag: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return model.Flag{}, fmt.Errorf("%w: flag with abbr %s", ErrNotFound, abbr)
	}

	var flag model.Flag
	err = rows.Scan(
		&flag.ID,
		&flag.Abbr,
		&flag.Emoji,
	)
	if err != nil {
		return model.Flag{}, fmt.Errorf("failed to scan flag: %w", err)
	}
	return flag, nil
}
