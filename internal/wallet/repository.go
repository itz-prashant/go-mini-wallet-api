package wallet

import (
	"context"
	"database/sql"
	"fmt"
)

type sqliteRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &sqliteRepository{
		db: db,
	}
}

func (s *sqliteRepository) Create(ctx context.Context, wallet *Wallet) error {
	query := `INSERT INTO wallets (owner_name, currency, balance, created_at, updated_at) VALUES(?,?,?,?,?)`

	result, err := s.db.ExecContext(ctx, query, wallet.OwnerName, wallet.Currency, wallet.Balance, wallet.CreatedAt, wallet.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return fmt.Errorf("Failed to get last insert id %w", err)
	}

	wallet.ID = id

	return nil
}
