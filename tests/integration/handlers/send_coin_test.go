package handlers

import (
	"avito-shop/internal/http-server/dto/request"
	rep "avito-shop/internal/http-server/dto/response"
	"avito-shop/internal/http-server/handlers"
	mw "avito-shop/internal/http-server/middleware"
	"avito-shop/internal/service"
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

type mockTransferService struct {
	mock.Mock
}

func (m *mockTransferService) SendCoins(ctx context.Context, from string, to string, amount int) error {
	args := m.Called(ctx, from, to, amount)
	return args.Error(0)
}

func TestNewSendCoinsHandlerFunc(t *testing.T) {
	validSender := "sender-user"
	validReceiver := "receiver-user"
	validAmount := 100
	validEmployeeId := uuid.New()

	tests := []struct {
		name           string
		setup          func(*mockTransferService) *http.Request
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "successful transfer",
			setup: func(mockService *mockTransferService) *http.Request {
				mockService.On("SendCoins", mock.Anything, validSender, validReceiver, validAmount).
					Return(nil)

				requestBody := request.SendCoinRequest{ToUser: validReceiver, Amount: validAmount}
				jsonBody, _ := json.Marshal(requestBody)

				req := httptest.NewRequest(http.MethodPost, "/api/send-coins", strings.NewReader(string(jsonBody)))
				req.Header.Set("Content-Type", "application/json")

				claims := &service.TokenClaims{
					Username:   validSender,
					EmployeeId: validEmployeeId,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				}
				ctx := context.WithValue(req.Context(), mw.UserContextKey, claims)
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "not enough coins",
			setup: func(mockService *mockTransferService) *http.Request {
				mockService.On("SendCoins", mock.Anything, validSender, validReceiver, validAmount).
					Return(service.ErrNotEnoughCoins)

				requestBody := request.SendCoinRequest{ToUser: validReceiver, Amount: validAmount}
				jsonBody, _ := json.Marshal(requestBody)

				req := httptest.NewRequest(http.MethodPost, "/api/send-coins", strings.NewReader(string(jsonBody)))
				req.Header.Set("Content-Type", "application/json")

				claims := &service.TokenClaims{
					Username:   validSender,
					EmployeeId: validEmployeeId,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				}
				ctx := context.WithValue(req.Context(), mw.UserContextKey, claims)
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   rep.ErrorResponse{Errors: "not enough coins to send"},
		},
		{
			name: "missing JWT token",
			setup: func(mockService *mockTransferService) *http.Request {
				requestBody := request.SendCoinRequest{ToUser: validReceiver, Amount: validAmount}
				jsonBody, _ := json.Marshal(requestBody)
				return httptest.NewRequest(http.MethodPost, "/api/send-coins", strings.NewReader(string(jsonBody)))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   rep.ErrorResponse{Errors: "internal server error"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
			mockTransferService := new(mockTransferService)

			req := tc.setup(mockTransferService)

			handler := handlers.NewSendCoinsHandlerFunc(logger, mockTransferService)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Result().StatusCode)

			if tc.expectedBody != nil {
				expectedResp, err := json.Marshal(tc.expectedBody)
				assert.NoError(t, err)
				assert.JSONEq(t, string(expectedResp), w.Body.String())
			}

			mockTransferService.AssertExpectations(t)
		})
	}
}
