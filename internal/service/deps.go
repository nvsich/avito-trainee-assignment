package service

import (
	"avito-shop/internal/model"
	"context"
	"github.com/google/uuid"
)

type EmployeeRepo interface {
	Save(ctx context.Context, employee *model.Employee) error
	FindByUsername(ctx context.Context, username string) (*model.Employee, error)
	UpdateByUsername(ctx context.Context, username string, employee *model.Employee) error
}

type TransferRepo interface {
	Save(ctx context.Context, transfer *model.Transfer) error
	FindAllForReceiverGroupedBySenders(ctx context.Context, receiverId uuid.UUID) ([]model.CoinTransaction, error)
	FindAllForSenderGroupedByReceivers(ctx context.Context, senderId uuid.UUID) ([]model.CoinTransaction, error)
}

type ItemRepo interface {
	FindByName(ctx context.Context, itemName string) (*model.Item, error)
	FindById(ctx context.Context, itemId uuid.UUID) (*model.Item, error)
}

type InventoryRepo interface {
	Save(ctx context.Context, inventory *model.EmployeeInventory) error
	FindAllInventoryItemsByEmployee(ctx context.Context, employeeId uuid.UUID) ([]model.InventoryItem, error)
	FindByEmployeeAndItem(ctx context.Context, employeeId uuid.UUID, itemId uuid.UUID) (*model.EmployeeInventory, error)
	UpdateById(ctx context.Context, id uuid.UUID, employeeInventory *model.EmployeeInventory) error
}
