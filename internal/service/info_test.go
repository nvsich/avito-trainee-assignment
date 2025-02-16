package service

import (
	"avito-shop/internal/model"
	"avito-shop/internal/repo"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestInfoService_Get(t *testing.T) {
	mockTrManager := new(mockTransactionManager)

	tests := []struct {
		name          string
		setup         func(*mockEmployeeRepo, *mockInventoryRepo, *mockTransferRepo)
		expectedError error
	}{
		{
			name: "successful info retrieval",
			setup: func(mer *mockEmployeeRepo, mir *mockInventoryRepo, mtr *mockTransferRepo) {
				employeeID := uuid.New()
				employee := &model.Employee{Id: employeeID, Username: "test_user", Balance: 1000}
				inventoryItems := []model.InventoryItem{{Type: "Item1", Quantity: 1}, {Type: "Item2", Quantity: 2}}
				coinHistorySent := []model.CoinTransaction{{User: "receiver", Amount: 50}}
				coinHistoryReceived := []model.CoinTransaction{{User: "sender", Amount: 30}}

				mer.On("FindByUsername", mock.Anything, "test_user").
					Return(employee, nil)
				mir.On("FindAllInventoryItemsByEmployee", mock.Anything, employeeID).
					Return(inventoryItems, nil)
				mtr.On("FindAllForSenderGroupedByReceivers", mock.Anything, employeeID).
					Return(coinHistorySent, nil)
				mtr.On("FindAllForReceiverGroupedBySenders", mock.Anything, employeeID).
					Return(coinHistoryReceived, nil)
			},
			expectedError: nil,
		},
		{
			name: "employee not found",
			setup: func(mer *mockEmployeeRepo, mir *mockInventoryRepo, mtr *mockTransferRepo) {
				mer.On("FindByUsername", mock.Anything, "test_user").
					Return(nil, repo.ErrEmployeeNotFound)
			},
			expectedError: ErrEmployeeNotFound,
		},
		{
			name: "inventory retrieval error",
			setup: func(mer *mockEmployeeRepo, mir *mockInventoryRepo, mtr *mockTransferRepo) {
				employeeID := uuid.New()
				employee := &model.Employee{Id: employeeID, Username: "test_user", Balance: 1000}

				mer.On("FindByUsername", mock.Anything, "test_user").
					Return(employee, nil)
				mir.On("FindAllInventoryItemsByEmployee", mock.Anything, employeeID).
					Return(nil, errors.New("inventory error"))
			},
			expectedError: errors.New("inventory error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockEmployeeRepo := new(mockEmployeeRepo)
			mockInventoryRepo := new(mockInventoryRepo)
			mockTransferRepo := new(mockTransferRepo)
			infoService := NewInfoService(
				mockTrManager, mockEmployeeRepo, mockInventoryRepo, mockTransferRepo, nil)

			tc.setup(mockEmployeeRepo, mockInventoryRepo, mockTransferRepo)

			info, err := infoService.Get(context.Background(), "test_user")

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, info)
			}

			mockEmployeeRepo.AssertExpectations(t)
			mockInventoryRepo.AssertExpectations(t)
			mockTransferRepo.AssertExpectations(t)
		})
	}
}
