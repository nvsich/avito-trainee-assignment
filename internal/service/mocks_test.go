package service

import (
	"avito-shop/internal/model"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockEmployeeRepo struct {
	mock.Mock
}

func (m *mockEmployeeRepo) Save(ctx context.Context, employee *model.Employee) error {
	args := m.Called(ctx, employee)
	return args.Error(0)
}

func (m *mockEmployeeRepo) FindByUsername(ctx context.Context, username string) (*model.Employee, error) {
	args := m.Called(ctx, username)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Employee), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockEmployeeRepo) UpdateByUsername(ctx context.Context, username string, employee *model.Employee) error {
	args := m.Called(ctx, username, employee)
	return args.Error(0)
}

type mockTransferRepo struct {
	mock.Mock
}

func (m *mockTransferRepo) Save(ctx context.Context, transfer *model.Transfer) error {
	args := m.Called(ctx, transfer)
	return args.Error(0)
}

func (m *mockTransferRepo) FindAllForReceiverGroupedBySenders(
	ctx context.Context, receiverId uuid.UUID) ([]model.CoinTransaction, error) {
	args := m.Called(ctx, receiverId)
	if args.Get(0) != nil {
		return args.Get(0).([]model.CoinTransaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockTransferRepo) FindAllForSenderGroupedByReceivers(
	ctx context.Context, senderId uuid.UUID) ([]model.CoinTransaction, error) {
	args := m.Called(ctx, senderId)
	if args.Get(0) != nil {
		return args.Get(0).([]model.CoinTransaction), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockItemRepo struct {
	mock.Mock
}

func (m *mockItemRepo) FindByName(ctx context.Context, itemName string) (*model.Item, error) {
	args := m.Called(ctx, itemName)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Item), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockItemRepo) FindById(ctx context.Context, itemId uuid.UUID) (*model.Item, error) {
	args := m.Called(ctx, itemId)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Item), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockInventoryRepo struct {
	mock.Mock
}

func (m *mockInventoryRepo) Save(ctx context.Context, inventory *model.EmployeeInventory) error {
	args := m.Called(ctx, inventory)
	return args.Error(0)
}

func (m *mockInventoryRepo) FindAllInventoryItemsByEmployee(
	ctx context.Context, employeeId uuid.UUID) ([]model.InventoryItem, error) {
	args := m.Called(ctx, employeeId)
	if args.Get(0) != nil {
		return args.Get(0).([]model.InventoryItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockInventoryRepo) FindByEmployeeAndItem(
	ctx context.Context, employeeId uuid.UUID, itemId uuid.UUID) (*model.EmployeeInventory, error) {
	args := m.Called(ctx, employeeId, itemId)
	if args.Get(0) != nil {
		return args.Get(0).(*model.EmployeeInventory), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockInventoryRepo) UpdateById(
	ctx context.Context, id uuid.UUID, employeeInventory *model.EmployeeInventory) error {
	args := m.Called(ctx, id, employeeInventory)
	return args.Error(0)
}

type mockTransactionManager struct {
}

func (m *mockTransactionManager) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
