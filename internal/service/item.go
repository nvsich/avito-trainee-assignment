package service

import (
	"avito-shop/internal/model"
	"avito-shop/internal/repo"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type ItemService struct {
	trManager     TransactionManager
	itemRepo      ItemRepo
	employeeRepo  EmployeeRepo
	inventoryRepo InventoryRepo
}

func NewItemService(
	trManager TransactionManager,
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

func (s *ItemService) Buy(ctx context.Context, itemName string, username string) error {
	const op = "service.ItemService.BuyItem"

	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		employee, err := s.employeeRepo.FindByUsername(ctx, username)
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

		employee.Balance -= item.Price

		if err = s.employeeRepo.UpdateByUsername(ctx, employee.Username, employee); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		employeeInventory, err := s.inventoryRepo.FindByEmployeeAndItem(ctx, employee.Id, item.Id)

		if err != nil {
			if errors.Is(err, repo.ErrEmployeeInventoryNotFound) {
				if err = s.inventoryRepo.Save(ctx, &model.EmployeeInventory{
					Id:         uuid.New(),
					EmployeeId: employee.Id,
					ItemId:     item.Id,
					Amount:     1,
				}); err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}
				return nil
			}

			return fmt.Errorf("%s: %w", op, err)
		}

		employeeInventory.Amount++

		if err = s.inventoryRepo.UpdateById(ctx, employeeInventory.Id, employeeInventory); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	return err
}
