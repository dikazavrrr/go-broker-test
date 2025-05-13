package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gitlab.com/digineat/go-broker-test/internal/model"
)

func TestInsertTrade(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer dbConn.Close()

	trade := model.Trade{
		Account: "acc1",
		Symbol:  "AAPL",
		Volume:  10,
		Open:    100.0,
		Close:   110.0,
		Side:    "buy",
	}

	mock.ExpectExec("INSERT INTO trades_q").
		WithArgs(trade.Account, trade.Symbol, trade.Volume, trade.Open, trade.Close, trade.Side).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = InsertTrade(dbConn, trade)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestGetAccountStats(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer dbConn.Close()

	account := "acc1"

	mock.ExpectQuery("SELECT account, trades, profit FROM account_stats WHERE account = ?").
		WithArgs(account).
		WillReturnRows(sqlmock.NewRows([]string{"account", "trades", "profit"}).
			AddRow("acc1", 1, 10.0))

	stats, err := GetAccountStats(dbConn, account)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Log(stats)
	if stats.Account != "acc1" || stats.Trades != 1 || stats.Profit != 10.0 {
		t.Errorf("unexpected result: %+v", stats)
	}
}
