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
		Currency:  req.Currency,
		Balance:   req.InitialBalance,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	err := s.repo.Create(ctx, wallet)

	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *Service) GetWallet(ctx context.Context, id int64) (*Wallet, error) {
	if id <= 0 {
		return nil, ErrInvalidId
	}

	wallet, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *Service) Transfer(ctx context.Context, req TransferRequest) error {
	if req.FromWalletId <= 0 || req.ToWalletId <= 0 {
		return ErrInvalidId
	}

	if req.FromWalletId == req.ToWalletId {
		return ErrSameWalletTransfer
	}

	if req.Amount <= 0 {
		return ErrInvalidAmount
	}

	err := s.repo.TransferTx(ctx, req.FromWalletId, req.ToWalletId,req.Amount, req.Description)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error) {
	if filter.WalletID <= 0 {
		return nil, ErrInvalidId
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	}else if filter.Limit > 100 {
		filter.Limit = 100
	}

	return s.repo.GetTransaction(ctx, filter)
}