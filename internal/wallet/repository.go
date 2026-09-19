package wallet

import (
	"context"
	"database/sql"
	"errors"
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

func (s *sqliteRepository) GetByID(ctx context.Context, id int64) (*Wallet, error) {

	query := `SELECT id, owner_name, balance, currency, created_at, updated_at FROM wallets WHERE id = ?`

	row := s.db.QueryRowContext(ctx, query, id)

	var wallet Wallet

	err := row.Scan(&wallet.ID, &wallet.OwnerName, &wallet.Balance, &wallet.Currency, &wallet.CreatedAt, &wallet.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to get wallet by id: %w", err)
	}

	return &wallet, nil
}

func (s *sqliteRepository) TransferTx(ctx context.Context, fromId int64, toId int64, amount float64, description string) error {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	deductQuery := `UPDATE wallets SET balance = balance - ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND balance >= ?`

	res, err := tx.ExecContext(ctx, deductQuery, amount, fromId, amount)

	if err != nil {
		return fmt.Errorf("failed to deduct balance from wallet %d: %w", fromId, err)
	}

	rowAffected, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to check rows affected on sender: %w", err)
	}

	if rowAffected == 0 {
		return ErrInsufficientBalance
	}

	creditQuery := `UPDATE wallets SET balance = balance + ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	res, err = tx.ExecContext(ctx, creditQuery, amount, toId)

	if err != nil {
		return fmt.Errorf("failed to credit balance to wallet %d: %w", toId, err)
	}

	rowAffected, err = res.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to check rows affected on receiver: %w", err)
	}

	if rowAffected == 0 {
		return ErrWalletNotFound
	}

	insertDebitQuery := `INSERT INTO transactions (wallet_id, counterpart_wallet_id, type, amount, description, created_at) VALUES(?,?,'transfer_out', ?, ?, CURRENT_TIMESTAMP)`

	_, err = tx.ExecContext(ctx, insertDebitQuery, fromId, toId, amount, description)

	if err != nil {
		return fmt.Errorf("failed to record debit transer: %w", err)
	}

	insertCreditQuery := `INSERT INTO transactions (wallet_id, counterpart_wallet_id, type, amount, description, created_at) VALUES(?,?,'transfer_in', ?, ?, CURRENT_TIMESTAMP)`

	_, err = tx.ExecContext(ctx, insertCreditQuery, toId, fromId, amount, description)

	if err != nil {
		return fmt.Errorf("failed to record credit transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *sqliteRepository) GetTransaction(ctx context.Context, filter TransactionFilter) ([]Transaction, error) {
	query := `SELECT id, wallet_id, counterpart_wallet_id, type, amount, description, created_at FROM transactions WHERE wallet_id = ?`

	args := []any{filter.WalletID}

	if filter.Type != "" {
		query += " AND type = ?"
		args = append(args, filter.Type)
	}

	offset := (filter.Page - 1) * filter.Limit
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, filter.Limit, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}

	defer rows.Close()

	transactions := make([]Transaction, 0)

	for rows.Next() {
		var t Transaction
		err := rows.Scan(
			&t.ID,
			&t.WalletID,
			&t.CounterpartWalletID,
			&t.Type,
			&t.Amount,
			&t.Description,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("transaction rows iteration error: %w", err)
	}

	return transactions, nil
}
