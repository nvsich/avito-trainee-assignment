package service

import (
	"avito-shop/internal/model"
	"context"
)

type EmployeeRepo interface {
	Save(ctx context.Context, employee *model.Employee) error
	FindByLogin(ctx context.Context, login string) (*model.Employee, error)
}

type TransferRepo interface {
	Save(ctx context.Context, transfer *model.Transfer) error
}

type TxProvider interface {
	RunInTransaction(txFunc func(adapters Adapters) error) error
}

type Adapters struct {
	EmployeeRepo  EmployeeRepo
	TransfersRepo TransferRepo
}
