package db

import (
	"database/sql"
	"errors"

	"gitlab.com/digineat/go-broker-test/internal/model"
)

func InsertTrade(db *sql.DB, t model.Trade) error {
	_, err := db.Exec(`
		INSERT INTO trades_q (account, symbol, volume, open, close, side, processed)
		VALUES (?, ?, ?, ?, ?, ?, 0)
	`, t.Account, t.Symbol, t.Volume, t.Open, t.Close, t.Side)
	return err
}

// GetAccountStats возвращает статистику по аккаунту
func GetAccountStats(db *sql.DB, acc string) (*model.AccountStats, error) {
	var stats model.AccountStats
	err := db.QueryRow(`
		SELECT account, trades, profit
		FROM account_stats
		WHERE account = ?
	`, acc).Scan(&stats.Account, &stats.Trades, &stats.Profit)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func CheckHealth(db *sql.DB) error {
	return db.Ping()
}
