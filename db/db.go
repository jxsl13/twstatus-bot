package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func New(filePath string) (*sql.DB, error) {
	return sql.Open("sqlite", filePath)
}
