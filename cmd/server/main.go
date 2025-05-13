package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	dbTrade "gitlab.com/digineat/go-broker-test/internal/db"
	"gitlab.com/digineat/go-broker-test/internal/model"
)

func main() {
	// Command line flags
	dbPath := flag.String("db", "data.db", "path to SQLite database")
	listenAddr := flag.String("listen", "8080", "HTTP server listen address")
	flag.Parse()

	// Initialize database connection
	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := dbTrade.InitSchema(db); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize HTTP server
	mux := http.NewServeMux()

	// POST /trades endpoint
	mux.HandleFunc("POST /trades", func(w http.ResponseWriter, r *http.Request) {
		var trade model.Trade
		if err := json.NewDecoder(r.Body).Decode(&trade); err != nil {
			log.Printf("Server Status: %+v", http.StatusBadRequest)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if trade.Account == "" || trade.Symbol == "" || trade.Volume <= 0 || (trade.Side != "buy" && trade.Side != "sell") {
			log.Printf("Server Status: %+v", http.StatusBadRequest)
			http.Error(w, "Invalid trade data", http.StatusBadRequest)
			return
		}

		log.Printf("Received trade: %+v", trade)

		if err := dbTrade.InsertTrade(db, trade); err != nil {
			log.Printf("Failed to insert trade: %v", err)
			log.Printf("Server Status: %+v", http.StatusInternalServerError)

			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		log.Printf("Server Status: %+v", http.StatusOK)

		w.WriteHeader(http.StatusOK)
	})

	// GET /stats/{acc} endpoint
	mux.HandleFunc("GET /stats/", func(w http.ResponseWriter, r *http.Request) {
		acc := strings.TrimPrefix(r.URL.Path, "/stats/")
		if acc == "" {
			log.Printf("Server Status: %+v", http.StatusBadRequest)
			http.Error(w, "Missing account", http.StatusBadRequest)
			return
		}

		stats, err := dbTrade.GetAccountStats(db, acc)
		if err != nil {
			log.Printf("Server Status: %+v", http.StatusInternalServerError)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if stats == nil {
			log.Printf("Server Status: %+v", http.StatusNotFound)
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})

	// GET /healthz endpoint
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := dbTrade.CheckHealth(db); err != nil {
			log.Printf("Server Status: %+v", http.StatusInternalServerError)

			http.Error(w, "Database connection failed", http.StatusInternalServerError)
			return
		}

		log.Printf("Server Status: %+v", http.StatusOK)
		w.WriteHeader(http.StatusOK)
	})

	// Start server
	serverAddr := fmt.Sprintf(":%s", *listenAddr)
	log.Printf("Starting server on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
