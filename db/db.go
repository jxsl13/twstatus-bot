package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
)

func New(ctx context.Context, host string, port uint16, database, username, password string, options ...Option) (*DB, error) {
	opts := &Options{
		sslmode:      SSLModeDisable,
		migrationsFs: nil,
		versionTable: "schema_versions",
		connTimeout:  30 * time.Second,
	}

	for _, option := range options {
		option(opts)
	}

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(username, password),
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   database,
		RawQuery: url.Values{
			"sslmode":          []string{opts.sslmode.String()},
			"application_name": []string{"twbot"},
		}.Encode(),
	}

	cfg, err := pgxpool.ParseConfig(dsn.String())
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	db := &DB{
		pool: pool,
	}

	if opts.migrationsFs == nil {
		// no migrations
		return db, nil
	}

	conn, closer, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer closer()

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

	m.OnStart = func(i int32, s1, s2, _ string) {
		log.Printf("migrating database: version %d %s %s\n", i, s1, s2)
	}

	err = m.Migrate(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type DB struct {
	pool *pgxpool.Pool
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) Conn(ctx context.Context) (c *pgx.Conn, closer func(), err error) {
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return nil, nil, err
	}
	return conn.Conn(), conn.Release, nil
}

func (db *DB) Tx(ctx context.Context) (tx pgx.Tx, closer func(error) error, err error) {
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return nil, nil, err
	}
	tx, err = conn.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	return tx, func(e error) error {
		defer conn.Release()
		if e != nil {
			return errors.Join(e, tx.Rollback(ctx))
		} else {
			return tx.Commit(ctx)
		}
	}, nil
}
