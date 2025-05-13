package db

import (
	"database/sql"

	"gitlab.com/digineat/go-broker-test/internal/model"
)

func FetchUnprocessedTrades(db *sql.DB) ([]model.TradeQ, error) {
	rows, err := db.Query(`
		SELECT id, account, volume, open, close, side
		FROM trades_q
		WHERE processed = 0
		LIMIT 100
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []model.TradeQ
	for rows.Next() {
		var t model.TradeQ
		if err := rows.Scan(&t.ID, &t.Account, &t.Volume, &t.Open, &t.Close, &t.Side); err != nil {
			return nil, err
		}
		trades = append(trades, t)
	}
	return trades, nil
}

func UpdateAccountStats(db *sql.DB, acc string, profit float64) error {
	_, err := db.Exec(`
		INSERT INTO account_stats (account, trades, profit)
		VALUES (?, 1, ?)
		ON CONFLICT(account) DO UPDATE SET
			trades = trades + 1,
			profit = profit + excluded.profit
	`, acc, profit)
	return err
}

func MarkTradeProcessed(db *sql.DB, id int64) error {
	_, err := db.Exec(`UPDATE trades_q SET processed = 1 WHERE id = ?`, id)
	return err
}
