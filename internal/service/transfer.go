package service

import (
	"avito-shop/internal/model"
	"avito-shop/internal/repo"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type TransferService struct {
	trManager    TransactionManager
	employeeRepo EmployeeRepo
	transferRepo TransferRepo
}

func NewTransferService(
	trManager TransactionManager,
	employeeRepo EmployeeRepo,
	transferRepo TransferRepo,
) *TransferService {
	return &TransferService{
		trManager:    trManager,
		employeeRepo: employeeRepo,
		transferRepo: transferRepo,
	}
}
func (s *TransferService) SendCoins(ctx context.Context, fromUsername string, toUsername string, amount int) error {
	const op = "service.TransferService.SendCoins"

	if fromUsername == toUsername {
		return ErrTransferToSameEmployee
	}

	if amount < 0 {
		return ErrNegativeTransferAmount
	}

	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		fromEmployee, err := s.employeeRepo.FindByUsername(ctx, fromUsername)
		if err != nil {
			if errors.Is(err, repo.ErrEmployeeNotFound) {
				return ErrSenderNotFound
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		if fromEmployee.Balance < amount {
			return ErrNotEnoughCoins
		}

		toEmployee, err := s.employeeRepo.FindByUsername(ctx, toUsername)
		if err != nil {
			if errors.Is(err, repo.ErrEmployeeNotFound) {
				return ErrReceiverNotFound
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		fromEmployee.Balance -= amount
		toEmployee.Balance += amount

		err = s.employeeRepo.UpdateByUsername(ctx, fromUsername, fromEmployee)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		err = s.employeeRepo.UpdateByUsername(ctx, toUsername, toEmployee)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = s.transferRepo.Save(ctx, &model.Transfer{
			Id:           uuid.New(),
			FromEmployee: fromEmployee.Id,
			ToEmployee:   toEmployee.Id,
			Amount:       amount,
		})
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	return err
}
