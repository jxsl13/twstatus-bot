package dao

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jxsl13/twstatus-bot/sqlc"
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

func NewDAO(q *sqlc.Queries) *DAO {
	return &DAO{
		q: q,
	}
}

type DAO struct {
	q *sqlc.Queries
}
