package storage_test

import (
	"bank-app/internal/model"
	"bank-app/internal/storage"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStorage_SaveNewAccount(t *testing.T) {
	tests := []struct {
		name string

		accountID       string
		initialAccounts []model.Account
		accountToSave   model.Account

		wantErr     bool
		corruptFile bool
	}{
		{
			name:            "account success",
			initialAccounts: []model.Account{},
			accountToSave:   model.Account{ID: "123abc", Owner: "Anton", Balance: 0},
		}, {
			name:            "allready exist",
			initialAccounts: []model.Account{{ID: "abc123", Owner: "Anton", Balance: 100}},
			accountToSave:   model.Account{ID: "abc123", Owner: "Anton", Balance: 100},
			wantErr:         true,
		}, {
			name:        "corrupted json file",
			corruptFile: true,
			accountToSave: model.Account{
				ID: "2", Owner: "Anton", Balance: 10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {})

		dir := t.TempDir()

		accPath := filepath.Join(dir, "accounts.json")
		txPath := filepath.Join(dir, "transactions.json")

		if tt.corruptFile {
			require.NoError(t, os.WriteFile(accPath, []byte("{not json file}"), 0644))
		} else if len(tt.initialAccounts) > 0 {
			writeJSON(t, accPath, tt.initialAccounts)
		}

		storage := storage.NewFileStorage(accPath, txPath)

		err := storage.SaveNewAccount(tt.accountToSave)

		if tt.wantErr {
			require.Error(t, err)
			return
		}

		accs := readAccounts(t, accPath)
		require.Len(t, accs, len(tt.initialAccounts)+1)

		require.NoError(t, err)

	}
}

func TestFileStorage_LoadAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		initialAccounts []model.Account
		accountID       string
		wantAccount     model.Account
		wantErr         bool
	}{
		{
			name: "load succuss",
			initialAccounts: []model.Account{
				{ID: "abc123", Owner: "Anton", Balance: 100},
			},
			accountID:   "abc123",
			wantAccount: model.Account{ID: "abc123", Owner: "Anton", Balance: 100},
		}, {
			name:      "account not found",
			accountID: "ABC!@#",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			accPath := filepath.Join(dir, "accounts.json")
			txPath := filepath.Join(dir, "transactions.json")

			if len(tt.initialAccounts) > 0 {
				writeJSON(t, accPath, tt.initialAccounts)
			}

			storage := storage.NewFileStorage(accPath, txPath)

			acc, err := storage.LoadAccount(tt.accountID)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantAccount, *acc)
		})
	}

}

func TestFileStorage_ApplyTransaction(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string

		initialAccounts    []model.Account
		initialTransaction []model.Transaction

		accountID string
		amount    float64
		tx        model.Transaction

		wantErr     bool
		wantBalance float64
		wantTxCount int
	}{
		{
			name: "deposit success",
			initialAccounts: []model.Account{
				{ID: "abc123", Owner: "Anton", Balance: 100},
			},
			accountID:   "abc123",
			amount:      50,
			tx:          model.NewDepositTransaction("abc123", 50),
			wantBalance: 150,
			wantTxCount: 1,
		}, {
			name: "withdraw success",
			initialAccounts: []model.Account{
				{ID: "abc123", Owner: "Anton", Balance: 100}},
			accountID:   "abc123",
			amount:      -50,
			tx:          model.NewWithdrawTransaction("abc123", 50),
			wantBalance: 50,
			wantTxCount: 1,
		}, {
			name: "withdraw insufficient funds",
			initialAccounts: []model.Account{
				{ID: "abc123", Owner: "Anton", Balance: 100}},
			accountID: "abc 123",
			amount:    -150,
			wantErr:   true,
		}, {
			name: "account not found",
			initialAccounts: []model.Account{
				{ID: "abc123", Owner: "Anton", Balance: 100}},
			accountID: "",
			amount:    100,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			dir := t.TempDir()

			accPath := filepath.Join(dir, "accounts.json")
			txPath := filepath.Join(dir, "transactions.json")

			fs := storage.NewFileStorage(accPath, txPath)

			if len(tt.initialAccounts) > 0 {
				writeJSON(t, accPath, tt.initialAccounts)
			}

			err := fs.ApplyTransaction(tt.accountID, tt.amount, tt.tx)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			accounts := readAccounts(t, accPath)
			require.Len(t, accounts, 1)
			assert.Equal(t, tt.wantBalance, accounts[0].Balance)

			txs := readTransactions(t, txPath)
			assert.Len(t, txs, tt.wantTxCount)
		})
	}
}

func writeJSON(t *testing.T, path string, v any) {
	t.Helper()

	data, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)

	err = os.WriteFile(path, data, 0644)
	require.NoError(t, err)
}

func readAccounts(t *testing.T, path string) []model.Account {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var accs []model.Account

	require.NoError(t, json.Unmarshal(data, &accs))

	return accs
}

func readTransactions(t *testing.T, path string) []model.Transaction {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var txs []model.Transaction

	require.NoError(t, json.Unmarshal(data, &txs))

	return txs
}
