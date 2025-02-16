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

func TestTransferService_SendCoins(t *testing.T) {
	mockTrManager := new(mockTransactionManager)

	tests := []struct {
		name          string
		setup         func(*mockEmployeeRepo, *mockTransferRepo)
		expectedError error
	}{
		{
			name: "successful transfer",
			setup: func(mer *mockEmployeeRepo, mtr *mockTransferRepo) {
				sender := &model.Employee{Id: uuid.New(), Username: "sender", Balance: 1000}
				receiver := &model.Employee{Id: uuid.New(), Username: "receiver", Balance: 500}

				mer.On("FindByUsername", mock.Anything, "sender").
					Return(sender, nil)
				mer.On("FindByUsername", mock.Anything, "receiver").
					Return(receiver, nil)
				mer.On("UpdateByUsername", mock.Anything, "sender", mock.Anything).
					Return(nil)
				mer.On("UpdateByUsername", mock.Anything, "receiver", mock.Anything).
					Return(nil)
				mtr.On("Save", mock.Anything, mock.Anything).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "sender not found",
			setup: func(mer *mockEmployeeRepo, mtr *mockTransferRepo) {
				mer.On("FindByUsername", mock.Anything, "sender").
					Return(nil, repo.ErrEmployeeNotFound)
			},
			expectedError: ErrSenderNotFound,
		},
		{
			name: "receiver not found",
			setup: func(mer *mockEmployeeRepo, mtr *mockTransferRepo) {
				sender := &model.Employee{Id: uuid.New(), Username: "sender", Balance: 1000}

				mer.On("FindByUsername", mock.Anything, "sender").
					Return(sender, nil)
				mer.On("FindByUsername", mock.Anything, "receiver").
					Return(nil, repo.ErrEmployeeNotFound)
			},
			expectedError: ErrReceiverNotFound,
		},
		{
			name: "not enough balance",
			setup: func(mer *mockEmployeeRepo, mtr *mockTransferRepo) {
				sender := &model.Employee{Id: uuid.New(), Username: "sender", Balance: 100}

				mer.On("FindByUsername", mock.Anything, "sender").
					Return(sender, nil)
			},
			expectedError: ErrNotEnoughCoins,
		},
		{
			name: "error saving transfer",
			setup: func(mer *mockEmployeeRepo, mtr *mockTransferRepo) {
				sender := &model.Employee{Id: uuid.New(), Username: "sender", Balance: 1000}
				receiver := &model.Employee{Id: uuid.New(), Username: "receiver", Balance: 500}

				mer.On("FindByUsername", mock.Anything, "sender").
					Return(sender, nil)
				mer.On("FindByUsername", mock.Anything, "receiver").
					Return(receiver, nil)
				mer.On("UpdateByUsername", mock.Anything, "sender", mock.Anything).
					Return(nil)
				mer.On("UpdateByUsername", mock.Anything, "receiver", mock.Anything).
					Return(nil)
				mtr.On("Save", mock.Anything, mock.Anything).
					Return(errors.New("save error"))
			},
			expectedError: errors.New("save error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockEmployeeRepo := new(mockEmployeeRepo)
			mockTransferRepo := new(mockTransferRepo)
			transferService := NewTransferService(mockTrManager, mockEmployeeRepo, mockTransferRepo)

			tc.setup(mockEmployeeRepo, mockTransferRepo)

			err := transferService.SendCoins(context.Background(), "sender", "receiver", 200)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockEmployeeRepo.AssertExpectations(t)
			mockTransferRepo.AssertExpectations(t)
		})
	}
}
