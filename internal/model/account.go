package model

import "errors"

type Account struct {
	ID      string
	Owner   string
	Balance float64
}

var (
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

func NewAccount(id string, owner string, balance float64) *Account {
	return &Account{
		ID:      id,
		Owner:   owner,
		Balance: balance,
	}
}

func (a *Account) Apply(amount float64) error {
	if amount == 0 {
		return ErrInvalidAmount
	}

	if amount < 0 && a.Balance+amount < 0 {
		return ErrInsufficientFunds
	}

	a.Balance += amount
	return nil
}
