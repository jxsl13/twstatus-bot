package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/tern/v2/migrate"
)

func New(ctx context.Context, host string, port int, database, username, password string, options ...Option) (*DB, error) {
	opts := &Options{
		ssl:          false,
		migrationsFs: nil,
		versionTable: "migration_versions",
	}

	for _, option := range options {
		option(opts)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		username,
		password,
		host,
		port,
		database,
		opts.SSL(),
	)

	if opts.migrationsFs == nil {
		// no migrations
		db, err := sql.Open("pgx", dsn)
		if err != nil {
			return nil, err
		}

		return &DB{
			DB: db,
		}, nil
	}

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	m, err := migrate.NewMigratorEx(
		ctx, conn,
		opts.versionTable,
		&migrate.MigratorOptions{
			DisableTx: false,
		},
	)
	if err != nil {
		return nil, err
	}
	err = m.LoadMigrations(opts.migrationsFs)
	if err != nil {
		return nil, err
	}

	m.OnStart = func(i int32, s1, s2, s3 string) {
		fmt.Printf("migrating database: version %d %s %s %s\n", i, s1, s2, s3)
	}

	err = m.Migrate(ctx)
	if err != nil {
		return nil, err
	}

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
