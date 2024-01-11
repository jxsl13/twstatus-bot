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
	guild_id INTEGER PRIMARY KEY,
	description TEXT NOT NULL DEFAULT ''
) STRICT;

CREATE TABLE IF NOT EXISTS channels (
	channel_id INTEGER PRIMARY KEY,
	guild_id INTEGER
		REFERENCES guilds(guild_id)
		ON DELETE CASCADE,
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
		REFERENCES guilds(guild_id)
		ON DELETE CASCADE,
	channel_id INTEGER
		REFERENCES channels(channel_id)
		ON DELETE CASCADE,
	flag_id INTEGER
		REFERENCES flags(flag_id)
		ON DELETE CASCADE,
	emoji TEXT NOT NULL,
	PRIMARY KEY (channel_id, flag_id)
) STRICT;

CREATE TABLE IF NOT EXISTS tracking (
	message_id INTEGER PRIMARY KEY,
	guild_id INTEGER
		REFERENCES guilds(guild_id)
		ON DELETE CASCADE,
	channel_id INTEGER
		REFERENCES channels(channel_id)
		ON DELETE CASCADE,
	address TEXT NOT NULL,
	CONSTRAINT tracking_unique_address UNIQUE (guild_id, channel_id, address)
	CONSTRAINT tracking_unique_message_id UNIQUE (guild_id, channel_id, message_id)
) STRICT;

CREATE TABLE IF NOT EXISTS player_count_notifications (
	guild_id INTEGER
		REFERENCES guilds(guild_id)
		ON DELETE CASCADE,
	channel_id INTEGER
		REFERENCES channels(channel_id)
		ON DELETE CASCADE,
	message_id INTEGER
		REFERENCES tracking(message_id)
		ON DELETE CASCADE
		ON UPDATE CASCADE,
	user_id INTEGER NOT NULL,
	threshold INTEGER NOT NULL
		CHECK( threshold > 0)
) STRICT;

CREATE TABLE IF NOT EXISTS prev_message_mentions (
	guild_id INTEGER
		REFERENCES guilds(guild_id)
		ON DELETE CASCADE,
	channel_id INTEGER
		REFERENCES channels(channel_id)
		ON DELETE CASCADE,
	message_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	PRIMARY KEY (guild_id, channel_id, message_id, user_id)
) STRICT;

CREATE TABLE IF NOT EXISTS active_servers (
	timestamp INTEGER NOT NULL,
	address TEXT PRIMARY KEY,
	protocols TEXT NOT NULL,
	name TEXT NOT NULL,
	gametype TEXT NOT NULL,
	passworded INTEGER NOT NULL DEFAULT 0
		CHECK( passworded IN (0,1)),
	map TEXT NOT NULL,
	map_sha256sum TEXT,
	map_size INTEGER,
	version TEXT NOT NULL,
	max_clients INTEGER NOT NULL,
	max_players INTEGER NOT NULL,
	score_kind TEXT NOT NULL DEFAULT 'points'
		CHECK(score_kind IN ('points','time'))
) STRICT;

CREATE TABLE IF NOT EXISTS active_server_clients (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	message_id INTEGER
		REFERENCES tracking(message_id)
		ON DELETE CASCADE,
	address TEXT
		REFERENCES active_servers(address)
		ON DELETE CASCADE,
	name TEXT NOT NULL,
	clan TEXT NOT NULL,
	country_id INTEGER
		REFERENCES flags(flag_id),
	score INTEGER NOT NULL,
	is_player INTEGER NOT NULL
		CHECK( is_player IN (0,1)),
	team INTEGER
) STRICT;

CREATE TABLE IF NOT EXISTS prev_active_servers (
	message_id INTEGER PRIMARY KEY
		REFERENCES tracking(message_id)
		ON DELETE CASCADE
		ON UPDATE CASCADE,
	guild_id INTEGER
		REFERENCES guilds(guild_id)
		ON DELETE CASCADE,
	channel_id INTEGER
		REFERENCES channels(channel_id)
		ON DELETE CASCADE,
	timestamp INTEGER NOT NULL,
	address TEXT NOT NULL,
	protocols TEXT NOT NULL,
	name TEXT NOT NULL,
	gametype TEXT NOT NULL,
	passworded INTEGER NOT NULL DEFAULT 0
		CHECK( passworded IN (0,1)),
	map TEXT NOT NULL,
	map_sha256sum TEXT,
	map_size INTEGER,
	version TEXT NOT NULL,
	max_clients INTEGER NOT NULL,
	max_players INTEGER NOT NULL,
	score_kind TEXT NOT NULL DEFAULT 'points'
		CHECK(score_kind IN ('points','time'))
) STRICT;

CREATE TABLE IF NOT EXISTS prev_active_server_clients (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	message_id INTEGER
		REFERENCES tracking(message_id)
		ON DELETE CASCADE
		ON UPDATE CASCADE,
	guild_id INTEGER
		REFERENCES guilds(guild_id)
		ON DELETE CASCADE,
	channel_id INTEGER
		REFERENCES channels(channel_id)
		ON DELETE CASCADE,
	name TEXT NOT NULL,
	clan TEXT NOT NULL,
	team INTEGER,
	country_id INTEGER
		REFERENCES flags(flag_id),
	score INTEGER NOT NULL,
	is_player INTEGER NOT NULL
		CHECK( is_player IN (0,1)),
	flag_abbr TEXT NOT NULL,
	flag_emoji TEXT NOT NULL
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
