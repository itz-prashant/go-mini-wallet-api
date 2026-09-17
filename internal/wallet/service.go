package wallet

import (
	"context"
	"errors"
	"strings"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

var ErrEmptyOwnerName = errors.New("owner name can not be empty")
var ErrNegativeBalance = errors.New("initial balance can not be negative")

func (s *Service) CreateWallet(ctx context.Context, req CreateWalletRequest) (*Wallet, error) {
	if strings.TrimSpace(req.OwnerName) == "" {
		return nil, ErrEmptyOwnerName
	}

	if req.InitialBalance < 0 {
		return nil, ErrNegativeBalance
	}

	if req.Currency == "" {
		req.Currency = "INR"
	}

	currentTime := time.Now().UTC()

	wallet := &Wallet{
		OwnerName: req.OwnerName,
		Currency: req.Currency,
		Balance: req.InitialBalance,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	err := s.repo.Create(ctx, wallet)

	if err != nil {
		return nil, err
	}

	return wallet, nil
}