package handlers

import (
	"avito-shop/internal/http-server/dto"
	rep "avito-shop/internal/http-server/dto/response"
	"avito-shop/internal/http-server/handlers"
	mw "avito-shop/internal/http-server/middleware"
	"avito-shop/internal/model"
	"avito-shop/internal/service"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type mockInfoService struct {
	mock.Mock
}

func (m *mockInfoService) Get(ctx context.Context, username string) (*model.EmployeeInfo, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*model.EmployeeInfo), args.Error(1)
}

func TestNewInfoHandlerFunc(t *testing.T) {
	validUsername := "valid-user"
	validEmployeeId := uuid.New()

	validEmployeeInfo := &model.EmployeeInfo{
		Coins: 123,
		Inventory: []model.InventoryItem{
			{"test-item", 1},
			{"test-item-1", 2},
		},
		CoinHistory: model.CoinHistory{
			Received: []model.CoinTransaction{
				{"from-test-user", 111},
			},
			Sent: []model.CoinTransaction{
				{"to-test-user", 222},
			},
		},
	}

	tests := []struct {
		name           string
		setup          func(*mockInfoService) *http.Request
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "successful retrieval of employee info",
			setup: func(mockInfo *mockInfoService) *http.Request {
				mockInfo.On("Get", mock.Anything, validUsername).
					Return(validEmployeeInfo, nil)

				req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
				claims := &service.TokenClaims{
					Username:   validUsername,
					EmployeeId: validEmployeeId,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				}
				ctx := context.WithValue(req.Context(), mw.UserContextKey, claims)
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   dto.ToInfoResponse(*validEmployeeInfo),
		},
		{
			name: "employee not found",
			setup: func(mockInfo *mockInfoService) *http.Request {
				mockInfo.On("Get", mock.Anything, validUsername).
					Return(&model.EmployeeInfo{}, service.ErrEmployeeNotFound)

				req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
				claims := &service.TokenClaims{
					Username:   validUsername,
					EmployeeId: validEmployeeId,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				}
				ctx := context.WithValue(req.Context(), mw.UserContextKey, claims)
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   rep.ErrorResponse{Errors: "employee not found"},
		},
		{
			name: "internal server error",
			setup: func(mockInfo *mockInfoService) *http.Request {
				mockInfo.On("Get", mock.Anything, validUsername).
					Return(&model.EmployeeInfo{}, errors.New("internal error"))

				req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
				claims := &service.TokenClaims{
					Username:   validUsername,
					EmployeeId: validEmployeeId,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				}
				ctx := context.WithValue(req.Context(), mw.UserContextKey, claims)
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   rep.ErrorResponse{Errors: "internal server error"},
		},
		{
			name: "missing JWT token in context",
			setup: func(mockInfo *mockInfoService) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/api/info", nil)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   rep.ErrorResponse{Errors: "internal server error"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
			mockInfoService := new(mockInfoService)

			req := tc.setup(mockInfoService)

			handler := handlers.NewInfoHandlerFunc(logger, mockInfoService)

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Result().StatusCode)

			expectedResp, err := json.Marshal(tc.expectedBody)
			assert.NoError(t, err)
			assert.JSONEq(t, string(expectedResp), w.Body.String())

			mockInfoService.AssertExpectations(t)
		})
	}
}
