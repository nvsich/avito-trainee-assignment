package handlers

import (
	"avito-shop/internal/http-server/dto/request"
	"avito-shop/internal/http-server/dto/response"
	"avito-shop/internal/http-server/handlers"
	"avito-shop/internal/service"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Authorize(ctx context.Context, username, password string) (string, error) {
	args := m.Called(ctx, username, password)
	return args.String(0), args.Error(1)
}

func setupRouter(log *slog.Logger, authService *mockAuthService) http.Handler {
	r := chi.NewRouter()
	r.Post("/api/auth", handlers.NewAuthHandlerFunc(log, authService))
	return r
}

func TestAuthHandler(t *testing.T) {
	const (
		validUser     = "valid-user"
		validPassword = "valid-password"
		validToken    = "valid-token"
		wrongPassword = "wrong-password"
	)
	tests := []struct {
		name           string
		setup          func(*mockAuthService) ([]byte, error)
		requestBody    any
		mockAuthReturn string
		mockAuthError  error
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "successful authentication",
			setup: func(mockAuth *mockAuthService) ([]byte, error) {
				reqBody, err := json.Marshal(request.AuthRequest{Username: validUser, Password: validPassword})
				if err != nil {
					return nil, err
				}
				mockAuth.On("Authorize", mock.Anything, validUser, validPassword).Return(validToken, nil)
				return reqBody, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   response.AuthResponse{Token: validToken},
		},
		{
			name: "invalid password",
			setup: func(mockAuth *mockAuthService) ([]byte, error) {
				reqBody, err := json.Marshal(request.AuthRequest{Username: validUser, Password: wrongPassword})
				if err != nil {
					return nil, err
				}
				mockAuth.On("Authorize", mock.Anything, validUser, wrongPassword).Return("", service.ErrInvalidCredentials)
				return reqBody, nil
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   response.ErrorResponse{Errors: "invalid credentials"},
		},
		{
			name: "service error",
			setup: func(mockAuth *mockAuthService) ([]byte, error) {
				reqBody, err := json.Marshal(request.AuthRequest{Username: validUser, Password: validPassword})
				if err != nil {
					return nil, err
				}
				mockAuth.On("Authorize", mock.Anything, validUser, validPassword).Return("", errors.New("internal error"))
				return reqBody, nil
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   response.ErrorResponse{Errors: "internal error"},
		},
		{
			name: "invalid JSON",
			setup: func(mockAuth *mockAuthService) ([]byte, error) {
				return []byte("{invalid_json}"), nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   response.ErrorResponse{Errors: "failed to parse request"},
		},
		{
			name: "empty request body",
			setup: func(mockAuth *mockAuthService) ([]byte, error) {
				return nil, nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   response.ErrorResponse{Errors: "failed to parse request"},
		},
		{
			name: "missing username",
			setup: func(mockAuth *mockAuthService) ([]byte, error) {
				reqBody, err := json.Marshal(request.AuthRequest{Password: validPassword})
				if err != nil {
					return nil, err
				}
				return reqBody, nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   response.ErrorResponse{Errors: "invalid request body"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
			mockAuthService := new(mockAuthService)

			reqBody, err := tc.setup(mockAuthService)
			assert.NoError(t, err)

			r := setupRouter(logger, mockAuthService)

			req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, "application/json", w.Result().Header.Get("Content-Type"))

			expectedResp, err := json.Marshal(tc.expectedBody)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatus, w.Result().StatusCode)
			assert.JSONEq(t, string(expectedResp), w.Body.String())

			mockAuthService.AssertExpectations(t)
		})
	}
}
