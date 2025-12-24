package main

import (
	"bank-app/internal/model"
	"bank-app/internal/service"
	"bank-app/internal/storage"
	"fmt"
)

func main() {
	acc := model.Account{
		ID:      "123",
		Owner:   "Stas",
		Balance: 0,
	}

	repo := storage.NewFileStorage(
		"data/accounts.json",
		"data/transactions.json")

	svc := service.NewService(repo)

	if err := svc.Deposit(acc.ID, 10); err != nil {
		fmt.Printf("Ошибка Deposit: %v\n", err)
	}
	if err := svc.Withdraw(acc.ID, 1); err != nil {
		fmt.Printf("Ошибка Withdraw: %v\n", err)
	}
}
