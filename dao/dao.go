package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jxsl13/twstatus-bot/db"
	"github.com/jxsl13/twstatus-bot/model"
	"modernc.org/sqlite"
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

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Conn interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

func InitDatabase(ctx context.Context, db *db.DB, wal bool) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}
	}()

	stmt := `
PRAGMA strict = ON;
PRAGMA foreign_keys = OFF;
`
	if wal {
		stmt += `
PRAGMA journal_mode = WAL;
`
	}

	stmt += `
CREATE TABLE IF NOT EXISTS guilds (
	guild_id INTEGER NOT NULL PRIMARY KEY,
	description TEXT NOT NULL DEFAULT ""
) STRICT;

CREATE TABLE IF NOT EXISTS channels (
	guild_id INTEGER NOT NULL REFERENCES guilds(guild_id) ON DELETE CASCADE,
	channel_id INTEGER PRIMARY KEY,
	running INTEGER
		CHECK( running IN (0,1))
		NOT NULL DEFAULT 0,
	UNIQUE(guild_id, channel_id)
) STRICT;
CREATE INDEX IF NOT EXISTS channels_id ON channels(channel_id);

CREATE TABLE IF NOT EXISTS flags (
	flag_id INTEGER PRIMARY KEY,
	abbr TEXT NOT NULL UNIQUE,
	emoji TEXT NOT NULL
) STRICT;
CREATE INDEX IF NOT EXISTS flags_id ON flags(flag_id);

CREATE TABLE IF NOT EXISTS flag_mappings (
	guild_id INTEGER
		NOT NULL
		REFERENCES guild(guild_id),
	channel_id INTEGER
		NOT NULL
		REFERENCES channels(channel_id)
			ON DELETE CASCADE,
	flag_id INTEGER NOT NULL
		REFERENCES flags(flag_id)
			ON DELETE CASCADE,
	emoji TEXT NOT NULL,
	PRIMARY KEY (flag_id, channel_id)
) STRICT;

CREATE TABLE IF NOT EXISTS tracking (
	guild_id INTEGER
		NOT NULL
		REFERENCES guild(guild_id),
	channel_id INTEGER
		NOT NULL
		REFERENCES channels(channel_id)
			ON DELETE CASCADE,
	address TEXT NOT NULL,
	message_id INTEGER NOT NULL,
	CONSTRAINT tracking_unique_address UNIQUE (guild_id, channel_id, address)
	CONSTRAINT tracking_unique_message_id UNIQUE (guild_id, channel_id, message_id)
) STRICT;

CREATE TABLE IF NOT EXISTS tw_servers (
	address TEXT PRIMARY KEY,
	protocols TEXT NOT NULL,
	name TEXT NOT NULL,
	gametype TEXT NOT NULL,
	passworded INTEGER
		CHECK( passworded IN (0,1))
		NOT NULL DEFAULT 0,
	map TEXT NOT NULL,
	map_sha256sum TEXT,
	map_size INTEGER,
	version TEXT NOT NULL,
	max_clients INTEGER NOT NULL,
	max_players INTEGER NOT NULL,
	score_kind TEXT
		CHECK(score_kind IN ('points','time'))
		NOT NULL DEFAULT 'points'
) STRICT;

CREATE TABLE IF NOT EXISTS tw_server_clients (
	address TEXT NOT NULL
		REFERENCES tw_servers(address)
			ON DELETE CASCADE,
	name TEXT NOT NULL,
	clan TEXT NOT NULL,
	country_id INTEGER NOT NULL
		REFERENCES flags(flag_id),
	score INTEGER NOT NULL,
	is_player INTEGER NOT NULL
) STRICT;
`
	stmt += `
PRAGMA foreign_key_check; -- validate foreign keys
`
	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	flagStmt, err := tx.PrepareContext(ctx, `
REPLACE INTO flags (flag_id, abbr, emoji)
VALUES (?, ?, ?);`,
	)
	if err != nil {
		return err
	}

	for _, flag := range model.Flags() {
		_, err = flagStmt.ExecContext(ctx, flag.ID, flag.Abbr, flag.Emoji)
		if err != nil {
			return fmt.Errorf("failed to insert flag %s: %w", flag.Abbr, err)
		}
	}

	return nil
}
