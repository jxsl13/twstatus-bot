package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/tern/migrate"
)

func New(host string, port int, database, username, password string, ssl bool) (*DB, error) {
	sslMode := "disable"
	if ssl {
		sslMode = "enable"
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		username,
		password,
		host,
		port,
		database,
		sslMode,
	)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &DB{
		DB: db,
	}, nil
}

type DB struct {
	*sql.DB
	sync.Mutex
}

func Migrations(ctx context.Context) {
	migrate.NewMigrator()
}
