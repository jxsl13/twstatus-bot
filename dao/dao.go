package dao

import "database/sql"

func InitDatabase(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS flags (
		id INTEGER,
		flag TEXT,
		symbol TEXT,
		PRIMARY KEY (id)
	);
	CREATE INDEX IF NOT EXISTS idx_flags_id ON flags (id);
	`)
	return err
}
