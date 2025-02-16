package handlers

import (
	"avito-shop/internal/http-server/dto/response"
	"avito-shop/internal/http-server/handlers"
	mw "avito-shop/internal/http-server/middleware"
	"avito-shop/internal/service"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type mockBuyItemService struct {
	mock.Mock
}

func (m *mockBuyItemService) Buy(ctx context.Context, itemName string, username string) error {
	args := m.Called(ctx, itemName, username)
	return args.Error(0)
}

func TestNewBuyItemHandlerFunc(t *testing.T) {
	validUsername := "valid-user"
	validItemName := "test-item"

	tests := []struct {
		name           string
		setup          func(*mockBuyItemService) *http.Request
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "successful purchase",
			setup: func(mockBuy *mockBuyItemService) *http.Request {
				mockBuy.On("Buy", mock.Anything, validItemName, validUsername).Return(nil)

				r := httptest.NewRequest(http.MethodGet, "/api/buy/"+validItemName, nil)
				rCtx := chi.NewRouteContext()
				rCtx.URLParams.Add("item", validItemName)
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rCtx))
				ctx := context.WithValue(r.Context(), mw.UserContextKey, &service.TokenClaims{
					Username: validUsername,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
					},
				})
				return r.WithContext(ctx)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
		},
		{
			name: "employee not found",
			setup: func(mockBuy *mockBuyItemService) *http.Request {
				mockBuy.On("Buy", mock.Anything, validItemName, validUsername).Return(service.ErrEmployeeNotFound)

				r := httptest.NewRequest(http.MethodGet, "/api/buy/"+validItemName, nil)
				rCtx := chi.NewRouteContext()
				rCtx.URLParams.Add("item", validItemName)
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rCtx))
				ctx := context.WithValue(r.Context(), mw.UserContextKey, &service.TokenClaims{
					Username: validUsername,
				})
				return r.WithContext(ctx)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   response.ErrorResponse{Errors: "employee not found"},
		},
		{
			name: "not enough coins",
			setup: func(mockBuy *mockBuyItemService) *http.Request {
				mockBuy.On("Buy", mock.Anything, validItemName, validUsername).Return(service.ErrNotEnoughCoins)

				r := httptest.NewRequest(http.MethodGet, "/api/buy/"+validItemName, nil)
				rCtx := chi.NewRouteContext()
				rCtx.URLParams.Add("item", validItemName)
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rCtx))
				ctx := context.WithValue(r.Context(), mw.UserContextKey, &service.TokenClaims{
					Username: validUsername,
				})
				return r.WithContext(ctx)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   response.ErrorResponse{Errors: "not enough coins"},
		},
		{
			name: "item not found",
			setup: func(mockBuy *mockBuyItemService) *http.Request {
				mockBuy.On("Buy", mock.Anything, validItemName, validUsername).Return(service.ErrItemNotFound)

				r := httptest.NewRequest(http.MethodGet, "/api/buy/"+validItemName, nil)
				rCtx := chi.NewRouteContext()
				rCtx.URLParams.Add("item", validItemName)
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rCtx))
				ctx := context.WithValue(r.Context(), mw.UserContextKey, &service.TokenClaims{
					Username: validUsername,
				})
				return r.WithContext(ctx)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   response.ErrorResponse{Errors: "item not found"},
		},
		{
			name: "internal server error",
			setup: func(mockBuy *mockBuyItemService) *http.Request {
				mockBuy.On("Buy", mock.Anything, validItemName, validUsername).Return(errors.New("internal error"))

				r := httptest.NewRequest(http.MethodGet, "/api/buy/"+validItemName, nil)
				rCtx := chi.NewRouteContext()
				rCtx.URLParams.Add("item", validItemName)
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rCtx))
				ctx := context.WithValue(r.Context(), mw.UserContextKey, &service.TokenClaims{
					Username: validUsername,
				})
				return r.WithContext(ctx)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   response.ErrorResponse{Errors: "internal server error"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
			mockBuyService := new(mockBuyItemService)

			req := tc.setup(mockBuyService)

			handler := handlers.NewBuyItemHandlerFunc(logger, mockBuyService)

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Result().StatusCode)

			if tc.expectedBody != nil {
				expectedResp, err := json.Marshal(tc.expectedBody)
				assert.NoError(t, err)
				assert.JSONEq(t, string(expectedResp), w.Body.String())
			}

			mockBuyService.AssertExpectations(t)
		})
	}
}
