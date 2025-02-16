package service

import (
	"avito-shop/internal/model"
	"avito-shop/internal/repo"
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestAuthService_Authorize(t *testing.T) {
	t.Parallel()

	mockRepo := new(mockEmployeeRepo)
	mockTrManager := new(mockTransactionManager)
	signKey := "test_key"
	tokenTTL := time.Hour

	authService := NewAuthService(mockTrManager, mockRepo, signKey, tokenTTL)

	existingUserID := uuid.New()
	existingUsername := "existing_user"
	password := "securePassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	newUsername := "new_user"
	newPassword := "newPassword"

	tests := []struct {
		name          string
		setup         func()
		username      string
		password      string
		expectedError error
		expectToken   bool
	}{
		{
			name: "successful auth of existing user",
			setup: func() {
				mockRepo.ExpectedCalls = nil

				mockRepo.On("FindByUsername", mock.Anything, existingUsername).
					Return(&model.Employee{Id: existingUserID, Username: existingUsername, PasswordHash: string(hashedPassword)}, nil)
			},
			username:      existingUsername,
			password:      password,
			expectedError: nil,
			expectToken:   true,
		},
		{
			name: "wrong password",
			setup: func() {
				mockRepo.ExpectedCalls = nil

				mockRepo.On("FindByUsername", mock.Anything, existingUsername).
					Return(&model.Employee{Id: existingUserID, Username: existingUsername, PasswordHash: string(hashedPassword)}, nil)
			},
			username:      existingUsername,
			password:      "wrong-password",
			expectedError: ErrInvalidCredentials,
			expectToken:   false,
		},
		{
			name: "creating new user",
			setup: func() {
				mockRepo.ExpectedCalls = nil

				mockRepo.On("FindByUsername", mock.Anything, newUsername).
					Return(nil, repo.ErrEmployeeNotFound)

				mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.Employee")).
					Return(nil)
			},
			username:      newUsername,
			password:      newPassword,
			expectedError: nil,
			expectToken:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			token, err := authService.Authorize(context.Background(), tc.username, tc.password)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			if tc.expectToken {
				assert.NotEmpty(t, token)

				parsedToken, _ := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(signKey), nil
				})

				assert.NotNil(t, parsedToken)
			} else {
				assert.Empty(t, token)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
