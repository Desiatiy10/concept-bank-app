package storage

import (
	"bank-app/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type Storage interface {
	Save(account model.Account) error
	Load(id string) (*model.Account, error)
	UpdateBalance(accountID string, newBalance float64) error

	SaveTransaction(model.Transaction) error
}

type FileStorage struct {
	accountFilePath     string
	TransactionFilePath string
	mu                  sync.Mutex
}

func NewFileStorage(accPath string, TxPath string) *FileStorage {
	return &FileStorage{
		accountFilePath:     accPath,
		TransactionFilePath: TxPath,
	}
}

func (fs *FileStorage) Save(acc model.Account) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	accounts := []model.Account{}

	data, err := os.ReadFile(fs.accountFilePath)
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &accounts); err != nil {
			return fmt.Errorf("corrupted json file: %w", err)
		}
		for _, a := range accounts {
			if acc.ID == a.ID {
				return fmt.Errorf("account with ID %s allready exist", acc.ID)
			}
		}
	}

	accounts = append(accounts, acc)

	newData, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshal accounts to file: %w", err)
	}

	err = os.WriteFile(fs.accountFilePath, newData, 0644)
	if err != nil {
		return fmt.Errorf("error open file: %w", err)
	}

	return nil
}

func (fs *FileStorage) Load(id string) (*model.Account, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	data, err := os.ReadFile(fs.accountFilePath)

	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("account %s not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("file read error: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("account %s not found", id)
	}

	accs := []model.Account{}

	if err := json.Unmarshal(data, &accs); err != nil {
		return nil, fmt.Errorf("corrupted json file %w", err)
	}

	for _, v := range accs {
		if id == v.ID {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("account not found")
}

func (fs *FileStorage) UpdateBalance(accountID string, newBalance float64) error {
	if accountID == "" {
		return fmt.Errorf("accountID cannon be empty")
	}
	if newBalance < 0 {
		return fmt.Errorf("amount must be 0 or high")
	}

	fs.mu.Lock()
	defer fs.mu.Unlock()

	data, err := os.ReadFile(fs.accountFilePath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	accs := []model.Account{}
	if err := json.Unmarshal(data, &accs); err != nil {
		return fmt.Errorf("corrupted file error: %w", err)
	}

	updated := false
	for i := range accs {
		if accs[i].ID == accountID {
			accs[i].Balance = newBalance
			updated = true
			break
		}
	}
	if !updated {
		return fmt.Errorf("account %s not found", accountID)
	}

	newData, err := json.MarshalIndent(accs, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	return os.WriteFile(fs.accountFilePath, newData, 0644)
}

func (fs *FileStorage) SaveTransaction(transaction model.Transaction) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	txs := []model.Transaction{}

	data, err := os.ReadFile(fs.TransactionFilePath)
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &txs); err != nil {
			return fmt.Errorf("failed to unmarshal transactions to slice: %v", err)
		}
	}

	txs = append(txs, transaction)

	newData, err := json.MarshalIndent(txs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %v", err)
	}

	return os.WriteFile(fs.TransactionFilePath, newData, 0644)
}
