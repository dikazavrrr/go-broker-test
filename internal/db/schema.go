package db

import (
	"database/sql"
)

func InitSchema(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS trades_q (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	account TEXT NOT NULL,
	symbol TEXT NOT NULL,
	volume REAL NOT NULL,
	open REAL NOT NULL,
	close REAL NOT NULL,
	side TEXT NOT NULL,
	processed BOOLEAN DEFAULT 0
);

CREATE TABLE IF NOT EXISTS account_stats (
	account TEXT PRIMARY KEY,
	trades INTEGER NOT NULL,
	profit REAL NOT NULL
);
`
	_, err := db.Exec(schema)
	return err
}
