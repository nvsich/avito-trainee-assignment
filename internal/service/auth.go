package service

import (
	"avito-shop/internal/model"
	"avito-shop/internal/repo"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const newEmployeeInitialBalance = 1000

type TokenClaims struct {
	jwt.StandardClaims
	EmployeeId uuid.UUID
	Login      string
}

type AuthService struct {
	employeeRepo EmployeeRepo
	signKey      string
	tokenTTL     time.Duration
}

func NewAuthService(employeeRepo EmployeeRepo, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		employeeRepo: employeeRepo,
		signKey:      signKey,
		tokenTTL:     tokenTTL,
	}
}

// TODO: подумать насчет хранения []byte(passwordHash) в БД
// TODO: hash password in handler?

func (s *AuthService) Authorize(ctx context.Context, login string, password string) (string, error) {
	const op = "service.AuthService.Authenticate"

	employee, err := s.employeeRepo.FindByLogin(ctx, login)

	if err != nil {
		if errors.Is(err, repo.ErrEmployeeNotFound) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return "", fmt.Errorf("%s: %w", op, err)
			}

			newEmployee := &model.Employee{
				Id:           uuid.New(),
				Login:        login,
				PasswordHash: string(hashedPassword),
				Balance:      newEmployeeInitialBalance,
			}

			err = s.employeeRepo.Save(ctx, newEmployee)
			if err != nil {
				// TODO: handle repo.ErrEmployeeExists (transactions?)
				return "", fmt.Errorf("%s: %w", op, err)
			}

			token, err := s.generateJWT(newEmployee.Id, newEmployee.Login)
			if err != nil {
				return "", fmt.Errorf("%s: %w", op, err)
			}
			return token, nil
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	// TODO: не создавать новый токен а проверять старый?
	return s.generateJWT(employee.Id, employee.Login)
}

func (s *AuthService) generateJWT(id uuid.UUID, login string) (string, error) {
	expirationTime := time.Now().Add(s.tokenTTL)
	claims := &TokenClaims{
		Login:      login,
		EmployeeId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.signKey))
}
