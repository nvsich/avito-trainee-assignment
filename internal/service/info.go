package service

import (
	"avito-shop/internal/model"
	"avito-shop/internal/repo"
	"context"
	"errors"
	"fmt"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type InfoService struct {
	trManager     *manager.Manager
	employeeRepo  EmployeeRepo
	inventoryRepo InventoryRepo
	transferRepo  TransferRepo
	itemRepo      ItemRepo
}

func NewInfoService(
	trManager *manager.Manager,
	employeeRepo EmployeeRepo,
	inventoryRepo InventoryRepo,
	transferRepo TransferRepo,
	itemRepo ItemRepo,
) *InfoService {
	return &InfoService{
		trManager:     trManager,
		employeeRepo:  employeeRepo,
		inventoryRepo: inventoryRepo,
		transferRepo:  transferRepo,
		itemRepo:      itemRepo,
	}
}

func (s *InfoService) Get(ctx context.Context, username string) (*model.EmployeeInfo, error) {
	const op = "service.InfoService.Get"
	var employeeInfo model.EmployeeInfo

	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		employee, err := s.employeeRepo.FindByUsername(ctx, username)
		if err != nil {
			if errors.Is(err, repo.ErrEmployeeNotFound) {
				return ErrEmployeeNotFound
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		inventoryItems, err := s.inventoryRepo.FindAllInventoryItemsByEmployee(ctx, employee.Id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		transfersAsSender, err := s.transferRepo.FindAllForSenderGroupedByReceivers(ctx, employee.Id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		transfersAsReceiver, err := s.transferRepo.FindAllForReceiverGroupedBySenders(ctx, employee.Id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		employeeInfo = model.EmployeeInfo{
			Coins:     employee.Balance,
			Inventory: inventoryItems,
			CoinHistory: model.CoinHistory{
				Sent:     transfersAsSender,
				Received: transfersAsReceiver,
			},
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &employeeInfo, nil
}
