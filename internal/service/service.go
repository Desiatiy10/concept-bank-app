package service

import (
	"bank-app/internal/model"
	"bank-app/internal/storage"
	"fmt"
)

type Service interface {
	Deposit(accountID string, amount float64) error
	Withdraw(accountID string, amount float64) error
}

type service struct {
	repo storage.Storage
}

func NewService(repo storage.Storage) Service {
	return &service{
		repo: repo,
	}
}



func (s *service) Deposit(accountID string, amount float64) error {
	if accountID == "" {
		return fmt.Errorf("empty ID field")
	}
	if amount <= 0 {
		return fmt.Errorf("amount should be greater than zero")
	}

	tx := model.NewDepositTransaction(accountID, amount)
	return s.repo.ApplyTransaction(accountID, amount, tx)
}

func (s *service) Withdraw(accountID string, amount float64) error {
	if accountID == "" {
		return fmt.Errorf("empty ID field")
	}
	if amount <= 0 {
		return fmt.Errorf("amount should be greater than zero Withdraw")
	}

	tx := model.NewWithdrawTransaction(accountID, amount)
	return s.repo.ApplyTransaction(accountID, -amount, tx)
}

// func (s *service) GetTransactions(accountID string) ([]model.Transaction, error) {

// 	return
// }
