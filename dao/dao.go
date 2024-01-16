package dao

import (
	"context"
	"errors"

	"github.com/jxsl13/twstatus-bot/db"
	"modernc.org/sqlite"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

const (
	UniqueConstraintViolation     = 2067
	PrimaryKeyConstraintViolation = 1555
)

func IsPrimaryKeyConstraintErr(err error) bool {
	serr, ok := err.(*sqlite.Error)
	if ok {
		return serr.Code() == PrimaryKeyConstraintViolation
	}
	return false
}

func IsUniqueConstraintErr(err error) bool {
	serr, ok := err.(*sqlite.Error)
	if ok {
		return serr.Code() == UniqueConstraintViolation
	}
	return false
}

func InitDatabase(ctx context.Context, db *db.DB) (err error) {

	return nil
}
