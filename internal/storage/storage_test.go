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

func TestFileStorage_ApplyTransaction(t *testing.T) {
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
		}, {},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {})

		dir := t.TempDir()

		txPathDir := filepath.Join(dir, "transaction.json")
		accPathDir := filepath.Join(dir, "account.json")

		fs := storage.NewFileStorage(txPathDir, accPathDir)

		if len(tt.initialAccounts) > 0 {
			writeJSON(t, accPathDir, tt.initialAccounts)
		}
		if len(tt.initialTransaction) > 0 {
			writeJSON(t, txPathDir, tt.initialTransaction)
		}

		err := fs.ApplyTransaction(tt.accountID, tt.amount, tt.tx)

		if tt.wantErr {
			require.Error(t, err)
			return
		}

		assert.NoError(t, err)

		accounts := readAccounts(t, accPathDir)
		require.Len(t, accounts, 1)
		txs := readTransactions(t, txPathDir)
		require.Len(t, txs, tt.wantTxCount)
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
