package wallet

import (
	"context"
	"errors"
	"time"
)

var ErrWalletNotFound = errors.New("wallet not found")
var ErrInvalidId = errors.New("invalid wallet id")
var ErrInsufficientBalance = errors.New("insufficient balance")
var ErrInvalidAmount = errors.New("amount must be greater than 0")
var ErrSameWalletTransfer = errors.New("cannot transfer money to the same wallet")

type Wallet struct {
	ID int64 `json:"id"`
	OwnerName string `json:"owner_name"`
	Currency string `json:"currency"`
	Balance float64 `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateWalletRequest struct {
	OwnerName string `json:"owner_name"`
	Currency string `json:"currency"`
	InitialBalance float64 `json:"initial_balance"`
}

type TransferRequest struct {
	FromWalletId int64 `json:"from_wallet_id"`
	ToWalletId int64 `json:"to_wallet_id"`
	Amount float64 `json:"amount"`
	Description string `json:"description"`
}

type Repository interface {
	Create(ctx context.Context, wallet *Wallet) error
	GetByID(ctx context.Context, id int64) (*Wallet, error)
	TransferTx(ctx context.Context, fromId int64, toId int64, amount float64, description string) error
}