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
	SaveNewAccount(account model.Account) error
	LoadAccount(accountID string) (*model.Account, error)
	ApplyTransaction(accountID string, amount float64, tx model.Transaction) error
}

type FileStorage struct {
	accountFilePath     string
	transactionFilePath string
	mu                  sync.Mutex
}

func NewFileStorage(accPath, txPath string) *FileStorage {
	return &FileStorage{
		accountFilePath:     accPath,
		transactionFilePath: txPath,
	}
}

func (fs *FileStorage) SaveNewAccount(acc model.Account) error {
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
		return fmt.Errorf("marshal accounts: %w", err)
	}

	err = os.WriteFile(fs.accountFilePath, newData, 0644)
	if err != nil {
		return fmt.Errorf("write acounts file: %w", err)
	}

	return nil
}

func (fs *FileStorage) LoadAccount(accountID string) (*model.Account, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	data, err := os.ReadFile(fs.accountFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("account %s not found", accountID)
	}
	if err != nil {
		return nil, fmt.Errorf("read account file: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("account %s not found", accountID)
	}

	accs := []model.Account{}

	if err := json.Unmarshal(data, &accs); err != nil {
		return nil, fmt.Errorf("corrupted json file %w", err)
	}

	for i := range accs {
		if accs[i].ID == accountID {
			return &accs[i], nil
		}
	}

	return nil, fmt.Errorf("account %v not found", accountID)
}

func (fs *FileStorage) ApplyTransaction(
	accountID string, amount float64, tx model.Transaction) error {
	if accountID == "" {
		return fmt.Errorf("empty ID field")
	}

	fs.mu.Lock()
	defer fs.mu.Unlock()

	if err := fs.applyTransactionUnsafe(accountID, amount, tx); err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) loadAccountsUnsafe() ([]model.Account, error) {
	dataAccs, err := os.ReadFile(fs.accountFilePath)
	if errors.Is(err, os.ErrNotExist) || len(dataAccs) == 0 {
		return []model.Account{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("file read: %w", err)
	}

	sliceAccs := []model.Account{}

	if err := json.Unmarshal(dataAccs, &sliceAccs); err != nil {
		return nil, fmt.Errorf("corrupted json file %w", err)
	}

	return sliceAccs, nil
}

func (fs *FileStorage) saveTransactionUnsafe(tx model.Transaction) error {
	sliceTransactions := []model.Transaction{}

	dataTxs, err := os.ReadFile(fs.transactionFilePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("read transaction file: %w", err)
		}
	} else if len(dataTxs) > 0 {
		if err := json.Unmarshal(dataTxs, &sliceTransactions); err != nil {
			return fmt.Errorf("unmarshal transactions: %w", err)
		}
	}

	sliceTransactions = append(sliceTransactions, tx)

	newDataTxs, err := json.MarshalIndent(sliceTransactions, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal transactions: %w", err)
	}

	if err := os.WriteFile(fs.transactionFilePath, newDataTxs, 0644); err != nil {
		return fmt.Errorf("write transactions file: %w", err)
	}

	return nil
}

func (fs *FileStorage) applyTransactionUnsafe(
	accountID string, amount float64, tx model.Transaction) error {
	accounts, err := fs.loadAccountsUnsafe()
	if err != nil {
		return err
	}

	var acc *model.Account
	for i := range accounts {
		if accounts[i].ID == accountID {
			acc = &accounts[i]
			break
		}
	}
	if acc == nil {
		return fmt.Errorf("account %s not found", accountID)
	}

	if err := acc.Apply(amount); err != nil {
		return err
	}

	if err := fs.writeAccountsUnsafe(accounts); err != nil {
		return err
	}

	if err := fs.saveTransactionUnsafe(tx); err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) writeAccountsUnsafe(accounts []model.Account) error {
	newData, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal accounts: %w", err)
	}

	if err := os.WriteFile(fs.accountFilePath, newData, 0644); err != nil {
		return fmt.Errorf("write accounts file: %w", err)
	}

	return nil
}
