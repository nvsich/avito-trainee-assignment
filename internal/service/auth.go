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
	Username   string
}

type AuthService struct {
	employeeRepo EmployeeRepo
	signKey      string
	tokenTTL     time.Duration
	trManager    TransactionManager
}

func NewAuthService(
	trManager TransactionManager, employeeRepo EmployeeRepo, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		employeeRepo: employeeRepo,
		signKey:      signKey,
		tokenTTL:     tokenTTL,
		trManager:    trManager,
	}
}

func (s *AuthService) Authorize(ctx context.Context, username string, password string) (string, error) {
	const op = "service.AuthService.Authorize"

	var token string
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		employee, err := s.getOrCreateEmployee(ctx, username, password)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if err = s.verifyPassword(employee.PasswordHash, password); err != nil {
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

func (s *AuthService) getOrCreateEmployee(ctx context.Context, username, password string) (*model.Employee, error) {
	employee, err := s.employeeRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repo.ErrEmployeeNotFound) {
			return s.createNewEmployee(ctx, username, password)
		}
		return nil, err
	}
	return employee, nil
}

func (s *AuthService) createNewEmployee(ctx context.Context, username, password string) (*model.Employee, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newEmployee := &model.Employee{
		Id:           uuid.New(),
		Username:     username,
		PasswordHash: string(hashedPassword),
		Balance:      newEmployeeInitialBalance,
	}

	if err = s.employeeRepo.Save(ctx, newEmployee); err != nil {
		return nil, err
	}

	return newEmployee, nil
}

func (s *AuthService) verifyPassword(storedHash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
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
