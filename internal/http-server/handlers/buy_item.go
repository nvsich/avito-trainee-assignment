package handlers

import (
	"avito-shop/internal/lib/logger/sl"
	"avito-shop/internal/service"
	"context"
	"errors"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

const itemParam = "item"

type BuyItem interface {
	Buy(ctx context.Context, itemName string, username string) error
}

func NewBuyItemHandlerFunc(log *slog.Logger, buyItemService BuyItem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.NewBuyItemHandlerFunc"
		log = setupLogger(log, op, r)

		itemName, ok := getURLParam(r, itemParam, log)
		if !ok {
			renderError(w, r, http.StatusBadRequest, "empty item name")
			return
		}

		claims, ok := getClaimsFromContext(r, log)
		if !ok {
			renderError(w, r, http.StatusInternalServerError, "internal server error")
			return
		}

		if err := buyItemService.Buy(r.Context(), itemName, claims.Username); err != nil {
			handleBuyError(w, r, log, err)
			return
		}

		render.Status(r, http.StatusOK)
	}
}

func handleBuyError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	var status int
	var message string

	switch {
	case errors.Is(err, service.ErrEmployeeNotFound):
		status, message = http.StatusUnauthorized, "employee not found"
	case errors.Is(err, service.ErrNotEnoughCoins):
		status, message = http.StatusBadRequest, "not enough coins"
	case errors.Is(err, service.ErrItemNotFound):
		status, message = http.StatusBadRequest, "item not found"
	default:
		status, message = http.StatusInternalServerError, "internal server error"
		log.Error("Buy operation failed", sl.Err(err))
	}

	if status != http.StatusInternalServerError {
		log.Info("Buy operation failed", sl.Err(err))
	}

	renderError(w, r, status, message)
}
