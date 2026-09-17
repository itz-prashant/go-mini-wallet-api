package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/itz-prashant/mini-wallet-api/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on", cfg.StoragePath)

	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS wallets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		owner_name TEXT NOT NULL,
		currency TEXT NOT NULL DEFAULT 'INR',
		balance REAL NOT NULL DEFAULT 0.0 CHECK (BALANCE >= 0.0),
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`);

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		wallet_id INTEGER NOT NULL,
		counterpart_wallet_id INTEGER,
		type TEXT NOT NULL CHECK (type IN ('credit', 'debit', 'transfer_in', 'transfer_out')),
		amount REAL NOT NULL CHECK (amount > 0.0),
		description TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (wallet_id) REFERENCES wallets (id),
		FOREIGN. KEY (counterpart_wallet_id) REFERENCES wallets (id)
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{Db: db}, nil
}