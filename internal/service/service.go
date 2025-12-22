package service

import (
	"bank-app/internal/model"
	"bank-app/internal/storage"
	"fmt"
	"time"
)

type Service interface {
	CheckBalance(id string) (float64, error)
	Deposit(accountID string, amount float64) error
	Withdraw(accountID string, amount float64) error

	GetTransactions(accountID string) ([]model.Transaction, error)
}

type service struct {
	repo storage.Storage
}

func NewService(repo storage.Storage) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CheckBalance(accountID string) (float64, error) {
	if accountID == "" {
		return 0, fmt.Errorf("empty ID field")
	}

	acc, err := s.repo.Load(accountID)
	if err != nil {
		return 0, err
	}

	return acc.Balance, err
}

func (s *service) Deposit(accountID string, amount float64) error {
	if accountID == "" {
		return fmt.Errorf("empty ID field")
	}
	if amount <= 0 {
		return fmt.Errorf("amount should be greater than zero")
	}

	acc, err := s.repo.Load(accountID)
	if err != nil {
		return err
	}

	newBalance := acc.Balance + amount

	transaction := model.Transaction{
		ID:                model.GenerateTransationID(),
		AccountID:         accountID,
		Type:              model.DepositTx,
		Amount:            amount,
		CreatedAt:         time.Now(),
	}

	if err := s.repo.SaveTransaction(transaction); err != nil {
		return err
	}

	return s.repo.UpdateBalance(acc.ID, newBalance)
}

func (s *service) Withdraw(accountID string, amount float64) error {
	if accountID == "" {
		return fmt.Errorf("empty ID field")
	}
	if amount <= 0 {
		return fmt.Errorf("amount should be greater than zero")
	}
	acc, err := s.repo.Load(accountID)
	if err != nil {
		return err
	}

	if amount > acc.Balance {
		return fmt.Errorf("cannot withdraw amount greater than balance")
	}

	newBalance := acc.Balance - amount

	transaction := model.Transaction{
		ID:                model.GenerateTransationID(),
		AccountID:         accountID,
		Type:              model.WithdrawTx,
		Amount:            amount,
		CreatedAt:         time.Now(),
	}

	if err := s.repo.SaveTransaction(transaction); err != nil {
		return err
	}

	return s.repo.UpdateBalance(accountID, newBalance)
}

func (s *service) GetTransactions(accountID string) ([]model.Transaction, error) {

	return
}
