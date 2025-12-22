package model

type Account struct {
	ID      string
	Owner   string
	Balance float64
}

func NewAccount(id string, owner string, balance float64) *Account {
	return &Account{
		ID:      id,
		Owner:   owner,
		Balance: balance,
	}
}
