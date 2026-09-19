package wallet

import (
	"context"
	"errors"
	"time"
)

var ErrWalletNotFound = errors.New("wallet not found")
var ErrInvalidId = errors.New("invalid wallet id") 

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

type Repository interface {
	Create(ctx context.Context, wallet *Wallet) error
	GetByID(ctx context.Context, id int64) (*Wallet, error)
}