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

	if err := repo.Save(acc); err == nil {
		fmt.Printf("Сохранен аккаунт: %v\n", acc.ID)
	} else {
		fmt.Printf("Ошибка: %v\n", err)
	}
	if err := svc.Deposit(acc.ID, 100); err == nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Printf("Баланс на аккаунте %v изменен. +%d\n", acc.ID, 100)
	}
	if err := svc.Withdraw(acc.ID, 50); err == nil {
		fmt.Printf("Баланс на аккаунте %v изменен. -%d\n", acc.ID, 50)
	} else {
		fmt.Printf("Ошибка: %v\n", err)
	}
}
