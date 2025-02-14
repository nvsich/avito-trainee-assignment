package service

import (
	"avito-shop/internal/model"
	"avito-shop/internal/repo"
	"context"
	"errors"
	"fmt"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const newEmployeeInitialBalance = 1000

type TokenClaims struct {
	jwt.StandardClaims
	EmployeeId uuid.UUID
	Username   string
}

type AuthService struct {
	employeeRepo EmployeeRepo
	signKey      string
	tokenTTL     time.Duration
	trManager    *manager.Manager
}

func NewAuthService(
	trManager *manager.Manager, employeeRepo EmployeeRepo, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		employeeRepo: employeeRepo,
		signKey:      signKey,
		tokenTTL:     tokenTTL,
		trManager:    trManager,
	}
}

func (s *AuthService) Authorize(ctx context.Context, username string, password string) (string, error) {
	const op = "service.AuthService.Authenticate"

	var token string

	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		employee, err := s.employeeRepo.FindByUsername(ctx, username)

		if err != nil {
			if errors.Is(err, repo.ErrEmployeeNotFound) {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				if err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}

				newEmployee := &model.Employee{
					Id:           uuid.New(),
					Username:     username,
					PasswordHash: string(hashedPassword),
					Balance:      newEmployeeInitialBalance,
				}

				err = s.employeeRepo.Save(ctx, newEmployee)
				if err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}

				token, err = s.generateJWT(newEmployee.Id, newEmployee.Username)
				if err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}
				return nil
			}

			return fmt.Errorf("%s: %w", op, err)
		}

		if err = bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(password)); err != nil {
			return ErrInvalidCredentials
		}

		token, err = s.generateJWT(employee.Id, employee.Username)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})

	return token, err
}

func (s *AuthService) generateJWT(id uuid.UUID, username string) (string, error) {
	expirationTime := time.Now().Add(s.tokenTTL)
	claims := &TokenClaims{
		Username:   username,
		EmployeeId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.signKey))
}
