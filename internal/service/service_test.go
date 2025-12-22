package service_test

import (
	"bank-app/internal/model"
	"bank-app/internal/service"
	"fmt"
	"testing"
)

type mockStorage struct {
	mockStorageAccs map[string]*model.Account

	FailedLoad   bool
	FailedUpdate bool

	UpdateBalanceCounter int
}

func NewMockStorage() *mockStorage {
	return &mockStorage{
		mockStorageAccs: make(map[string]*model.Account),
	}
}

func (ms *mockStorage) Load(accountID string) (*model.Account, error) {
	if ms.FailedLoad {
		return nil, fmt.Errorf("load failed by flag")
	}

	if acc, exist := ms.mockStorageAccs[accountID]; exist {
		return acc, nil
	}

	return nil, fmt.Errorf("account %s not found", accountID)
}

func (ms *mockStorage) UpdateBalance(accountID string, amount float64) error {
	ms.UpdateBalanceCounter++

	if ms.FailedUpdate {
		return fmt.Errorf("update failed by flag")
	}

	if acc, exist := ms.mockStorageAccs[accountID]; exist {
		acc.Balance = amount
		return nil
	}

	return fmt.Errorf("account %s not found", accountID)
}

func (ms *mockStorage) Save(acc model.Account) error {
	ms.mockStorageAccs[acc.ID] = &acc
	return nil
}

func (ms *mockStorage) GetBalance(accountID string) float64 {
	return ms.mockStorageAccs[accountID].Balance
}

func TestService_Deposit(t *testing.T) {
	tests := []struct {
		name               string
		accountID          string
		amount             float64
		setupMock          func(*mockStorage)
		wantErr            bool
		wantBalance        float64
		wantBalanceCounter int
		msgErr             string
	}{
		{
			name:      "без ошибок",
			accountID: "acc1",
			amount:    100.0,
			setupMock: func(ms *mockStorage) {
				ms.Save(model.Account{ID: "acc1", Owner: "anton", Balance: 1})
			},
			wantErr:            false,
			wantBalance:        101.0,
			wantBalanceCounter: 1,
		}, {
			name:      "empty ID",
			accountID: "",
			amount:    100,
			setupMock: func(ms *mockStorage) {},
			wantErr:   true,
			msgErr:    "empty ID field",
		}, {
			name:      "amount <=0",
			accountID: "acc1",
			amount:    -5,
			setupMock: func(ms *mockStorage) {
				ms.Save(model.Account{ID: "acc1", Owner: "anton", Balance: 2})
			},
			wantErr: true,
			msgErr:  "amount should be greater than zero",
		}, {
			name:      "account not found",
			accountID: "XxX",
			amount:    10,
			setupMock: func(ms *mockStorage) {},
			wantErr:   true,
			msgErr:    "account XxX not found",
		}, {
			name:      "load flag err",
			accountID: "acc1",
			amount:    100,
			setupMock: func(ms *mockStorage) {
				ms.Reset()
				ms.Save(model.Account{ID: "acc1", Owner: "anton", Balance: 100})
				ms.FailedLoad = true
			},
			wantErr: true,
			msgErr:  "load failed by flag",
		}, {
			name:      "update balance flag err",
			accountID: "acc1",
			amount:    100,
			setupMock: func(ms *mockStorage) {
				ms.Reset()
				ms.Save(model.Account{ID: "acc1", Owner: "anton", Balance: 100})
				ms.FailedUpdate = true
			},
			wantErr:            true,
			wantBalanceCounter: 1,
			msgErr:             "update failed by flag",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockStorage()
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewService(mockRepo)

			err := svc.Deposit(tt.accountID, tt.amount)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Deposit() expected err, got nil")
				} else if tt.msgErr != "" && err.Error() != tt.msgErr {
					t.Errorf("Deposit() error = %v, want %v", err.Error(), tt.msgErr)
				}
			} else {
				if err != nil {
					t.Errorf("Deposit() unexpected err: %v", err)
				}
			}

			if !tt.wantErr && err == nil {
				acc, _ := mockRepo.Load(tt.accountID)
				if acc.Balance != tt.wantBalance {
					t.Errorf("Deposit() balance = %f, exepcted = %f",
						acc.Balance, tt.wantBalance)
				}
			}
			if mockRepo.UpdateBalanceCounter != tt.wantBalanceCounter {
				t.Errorf("Deposit() UpdateBalanceCounter have: %v , want: %v",
					mockRepo.UpdateBalanceCounter, tt.wantBalanceCounter)
			}
		})
	}
}

func TestService_Withdraw(t *testing.T) {
	tests := []struct {
		name               string
		accountID          string
		amount             float64
		setupMock          func(ms *mockStorage)
		wantBalance        float64
		wantErr            bool
		wantBalanceCounter int
		msgErr             string
	}{
		{
			name:      "good try",
			accountID: "acc1",
			amount:    10,
			setupMock: func(ms *mockStorage) {
				ms.Save(model.Account{ID: "acc1", Owner: "anton", Balance: 100})
			},
			wantBalance:        90,
			wantErr:            false,
			wantBalanceCounter: 1,
		}, {
			name:      "empty ID",
			accountID: "",
			amount:    0,
			wantErr:   true,
			msgErr:    "empty ID field",
		}, {
			name:      "amount should be greater than zero",
			accountID: "cxc",
			wantErr:   true,
			msgErr:    "amount should be greater than zero",
		}, {
			name:      "account not found",
			accountID: "XXX",
			amount:    10,
			setupMock: func(ms *mockStorage) {
				ms.Save(model.Account{ID: "acc1", Owner: "anton", Balance: 100})
			},
			wantErr: true,
			msgErr:  "account XXX not found",
		}, {
			name:      "cannot withdraw amount greater than balance",
			accountID: "acc1",
			amount:    100,
			setupMock: func(ms *mockStorage) {
				ms.Save(model.Account{ID: "acc1", Owner: "anton", Balance: 99})
			},
			wantErr:            true,
			msgErr:             "cannot withdraw amount greater than balance",
			wantBalanceCounter: 0,
		}, {
			name:      "update flag error",
			accountID: "acc1",
			amount:    100,
			setupMock: func(ms *mockStorage) {
				ms.Reset()
				ms.Save(model.Account{ID: "acc1", Owner: "anton", Balance: 100})
				ms.FailedUpdate = true
			},
			wantErr:            true,
			msgErr:             "update failed by flag",
			wantBalanceCounter: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockStorage()
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewService(mockRepo)

			err := svc.Withdraw(tt.accountID, tt.amount)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Withdraw() expected err, got nil")
				} else if tt.msgErr != "" && err.Error() != tt.msgErr {
					t.Errorf("Withdraw() have err: %v, want: %v", err, tt.msgErr)
				}
			} else {
				if err != nil {	
					t.Errorf("Withdraw() unexpected error: %v", err)
				}
			}

			if tt.wantBalanceCounter != mockRepo.UpdateBalanceCounter {
				t.Errorf("want conunt: %d, have %d", tt.wantBalanceCounter, mockRepo.UpdateBalanceCounter)
			}

			if !tt.wantErr && err == nil {
				acc, _ := mockRepo.Load(tt.accountID)
				if acc.Balance != tt.wantBalance {
					t.Errorf("balance %f - amount %f != %f wantBalance",
						acc.Balance, tt.amount, tt.wantBalance)
				}
			}
		})
	}
}

func TestService_CheckBalance(t *testing.T) {
	tests := []struct {
		name        string
		accountID   string
		wantBalance float64
		setupMock   func(ms *mockStorage)
		wantErr     bool
		msgErr      string
	}{
		{
			name:        "good try",
			accountID:   "abc123",
			wantBalance: 100,
			setupMock: func(ms *mockStorage) {
				ms.Save(model.Account{ID: "abc123", Owner: "anton", Balance: 100})
			},
			wantErr: false,
		}, {
			name:      "empty ID field",
			accountID: "",
			wantErr:   true,
			msgErr:    "empty ID field",
		}, {
			name:      "flag load err",
			accountID: "abc123",
			setupMock: func(ms *mockStorage) {
				ms.Reset()
				ms.Save(model.Account{ID: "abc123", Owner: "anton", Balance: 100})
				ms.FailedLoad = true
			},
			wantErr: true,
			msgErr:  "load failed by flag",
		}, {
			name:      "account not found",
			accountID: "non-existent",
			wantErr:   true,
			msgErr:    "account non-existent not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockStorage()
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewService(mockRepo)

			gotBalance, err := svc.CheckBalance(tt.accountID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CheckBalance want: %v, have: nil", tt.msgErr)
				} else if tt.msgErr != "" && err.Error() != tt.msgErr {
					t.Errorf("CheckBalance want: %v, have: %v", tt.msgErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("CheckBalance() unexpected err: %v", err)
				}
				if gotBalance != tt.wantBalance {
					t.Errorf("CheckBalance() balance = %f, want %f", gotBalance, tt.wantBalance)
				}
			}
		})
	}
}

func (ms *mockStorage) Reset() {
	ms.FailedLoad = false
	ms.FailedUpdate = false
	ms.UpdateBalanceCounter = 0
	ms.mockStorageAccs = make(map[string]*model.Account)
}
