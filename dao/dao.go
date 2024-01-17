package dao

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

const (
	UniqueConstraintViolation = "23505"
)

func IsUniqueConstraintErr(err error) bool {
	serr, ok := err.(*pgconn.PgError)
	if ok {
		return serr.Code == UniqueConstraintViolation
	}
	return false
}
