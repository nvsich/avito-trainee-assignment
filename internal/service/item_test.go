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

func TestItemService_Buy(t *testing.T) {
	mockTrManager := new(mockTransactionManager)

	tests := []struct {
		name          string
		setup         func(*mockItemRepo, *mockEmployeeRepo, *mockInventoryRepo)
		expectedError error
	}{
		{
			name: "successful item purchase",
			setup: func(mir *mockItemRepo, mer *mockEmployeeRepo, minr *mockInventoryRepo) {
				employeeID := uuid.New()
				employee := &model.Employee{Id: employeeID, Username: "test_user", Balance: 1000}
				item := &model.Item{Id: uuid.New(), Name: "Item1", Price: 500}
				employeeInventory := &model.EmployeeInventory{
					Id: uuid.New(), EmployeeId: employeeID, ItemId: item.Id, Amount: 1}

				mer.On("FindByUsername", mock.Anything, "test_user").
					Return(employee, nil)
				mir.On("FindByName", mock.Anything, "Item1").
					Return(item, nil)
				minr.On("FindByEmployeeAndItem", mock.Anything, employeeID, item.Id).
					Return(employeeInventory, nil)
				mer.On("UpdateByUsername", mock.Anything, "test_user", mock.Anything).
					Return(nil)
				minr.On("UpdateById", mock.Anything, employeeInventory.Id, mock.Anything).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "employee not found",
			setup: func(mir *mockItemRepo, mer *mockEmployeeRepo, minr *mockInventoryRepo) {
				mer.On("FindByUsername", mock.Anything, "test_user").
					Return(nil, repo.ErrEmployeeNotFound)
			},
			expectedError: ErrEmployeeNotFound,
		},
		{
			name: "item not found",
			setup: func(mir *mockItemRepo, mer *mockEmployeeRepo, minr *mockInventoryRepo) {
				employee := &model.Employee{Id: uuid.New(), Username: "test_user", Balance: 1000}

				mer.On("FindByUsername", mock.Anything, "test_user").
					Return(employee, nil)
				mir.On("FindByName", mock.Anything, "Item1").
					Return(nil, repo.ErrItemNotFound)
			},
			expectedError: ErrItemNotFound,
		},
		{
			name: "not enough balance",
			setup: func(mir *mockItemRepo, mer *mockEmployeeRepo, minr *mockInventoryRepo) {
				employee := &model.Employee{Id: uuid.New(), Username: "test_user", Balance: 100}
				item := &model.Item{Id: uuid.New(), Name: "Item1", Price: 500}

				mer.On("FindByUsername", mock.Anything, "test_user").
					Return(employee, nil)
				mir.On("FindByName", mock.Anything, "Item1").
					Return(item, nil)
			},
			expectedError: ErrNotEnoughCoins,
		},
		{
			name: "error updating balance",
			setup: func(mir *mockItemRepo, mer *mockEmployeeRepo, minr *mockInventoryRepo) {
				employee := &model.Employee{Id: uuid.New(), Username: "test_user", Balance: 1000}
				item := &model.Item{Id: uuid.New(), Name: "Item1", Price: 500}

				mer.On("FindByUsername", mock.Anything, "test_user").
					Return(employee, nil)
				mir.On("FindByName", mock.Anything, "Item1").
					Return(item, nil)
				mer.On("UpdateByUsername", mock.Anything, "test_user", mock.Anything).
					Return(errors.New("update error"))
			},
			expectedError: errors.New("update error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockItemRepo := new(mockItemRepo)
			mockEmployeeRepo := new(mockEmployeeRepo)
			mockInventoryRepo := new(mockInventoryRepo)
			itemService := NewItemService(mockTrManager, mockItemRepo, mockEmployeeRepo, mockInventoryRepo)

			tc.setup(mockItemRepo, mockEmployeeRepo, mockInventoryRepo)

			err := itemService.Buy(context.Background(), "Item1", "test_user")

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockItemRepo.AssertExpectations(t)
			mockEmployeeRepo.AssertExpectations(t)
			mockInventoryRepo.AssertExpectations(t)
		})
	}
}
