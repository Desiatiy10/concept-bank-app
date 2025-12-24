package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	DepositTx  TransactionType = "deposit"
	WithdrawTx TransactionType = "withdraw"
)

type Transaction struct {
	ID        string          `json:"id"`
	AccountID string          `json:"account_id"`
	Type      TransactionType `json:"type"`
	Amount    float64         `json:"amount"`
	CreatedAt time.Time       `json:"created_at"`
}

func generateTransationID() string {
	return uuid.New().String()
}

func NewDepositTransaction(accountID string, amount float64) Transaction {
	return Transaction{
		ID:        generateTransationID(),
		AccountID: accountID,
		Type:      DepositTx,
		Amount:    amount,
		CreatedAt: time.Now(),
	}
}

func NewWithdrawTransaction(accountID string, amount float64) Transaction {
	return Transaction{
		ID:        generateTransationID(),
		AccountID: accountID,
		Type:      WithdrawTx,
		Amount:    amount,
		CreatedAt: time.Now(),
	}
}
