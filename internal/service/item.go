package service

import (
	"avito-shop/internal/model"
	"avito-shop/internal/repo"
	"context"
	"errors"
	"fmt"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
)

type ItemService struct {
	trManager     *manager.Manager
	itemRepo      ItemRepo
	employeeRepo  EmployeeRepo
	inventoryRepo InventoryRepo
}

func NewItemService(
	trManager *manager.Manager,
	itemRepo ItemRepo,
	employeeRepo EmployeeRepo,
	inventoryRepo InventoryRepo,
) *ItemService {
	return &ItemService{
		trManager:     trManager,
		itemRepo:      itemRepo,
		employeeRepo:  employeeRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *ItemService) BuyItem(ctx context.Context, itemName string, login string) error {
	const op = "service.ItemService.BuyItem"

	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		employee, err := s.employeeRepo.FindByLogin(ctx, login)
		if err != nil {
			if errors.Is(err, repo.ErrEmployeeNotFound) {
				return ErrEmployeeNotFound
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		item, err := s.itemRepo.FindByName(ctx, itemName)
		if err != nil {
			if errors.Is(err, repo.ErrItemNotFound) {
				return ErrItemNotFound
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		if item.Price > employee.Balance {
			return ErrNotEnoughCoins
		}

		employee.Balance = employee.Balance - item.Price

		if err = s.employeeRepo.UpdateByLogin(ctx, employee.Login, employee); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if err = s.inventoryRepo.Save(ctx, &model.EmployeeInventory{
			Id:         uuid.New(),
			EmployeeId: employee.Id,
			ItemId:     item.Id,
		}); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	return err
}
