package main

import (
	"database/sql"
	"flag"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	util "gitlab.com/digineat/go-broker-test/internal"
	dbWorker "gitlab.com/digineat/go-broker-test/internal/db"
)

func main() {
	// Command line flags
	dbPath := flag.String("db", "data.db", "path to SQLite database")
	pollInterval := flag.Duration("poll", 100*time.Millisecond, "polling interval")
	flag.Parse()

	// Initialize database connection
	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Printf("Worker started with polling interval: %v", *pollInterval)

	// Main worker loop
	for {
		trades, err := dbWorker.FetchUnprocessedTrades(db)
		if err != nil {
			log.Printf("Error fetching trades: %v", err)
			time.Sleep(*pollInterval)
			continue
		}

		for _, t := range trades {
			lot := 100000.0
			var profit float64
			if t.Side == "buy" {
				profit = util.Round((t.Close-t.Open)*t.Volume*lot, 2)
			} else if t.Side == "sell" {
				profit = util.Round((t.Open-t.Close)*t.Volume*lot, 2)
			}

			if err := dbWorker.UpdateAccountStats(db, t.Account, profit); err != nil {
				log.Printf("Failed to update stats for %s: %v", t.Account, err)
				continue
			}

			if err := dbWorker.MarkTradeProcessed(db, t.ID); err != nil {
				log.Printf("Failed to mark trade %d as processed: %v", t.ID, err)
			}
		}

		// Sleep for the specified interval
		time.Sleep(*pollInterval)
	}
}
