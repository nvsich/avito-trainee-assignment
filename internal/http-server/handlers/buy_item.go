package handlers

import (
	resp "avito-shop/internal/http-server/dto/response"
	mw "avito-shop/internal/http-server/middleware"
	"avito-shop/internal/lib/logger/sl"
	"avito-shop/internal/service"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

const item = "item"

type BuyItem interface {
	Buy(ctx context.Context, itemName string, username string) error
}

func NewBuyItemHandlerFunc(log *slog.Logger, buyItemService BuyItem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.NewBuyItemHandlerFunc"

		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		itemName := chi.URLParam(r, item)

		if itemName == "" {
			log.Error("empty item name")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ErrorResponse{Errors: "empty item name"})
			return
		}

		claims, ok := r.Context().Value(mw.UserContextKey).(*service.TokenClaims)
		if !ok || claims == nil {
			log.Error("failed to get claims from context")

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.ErrorResponse{Errors: "internal server error"})
			return
		}

		err := buyItemService.Buy(r.Context(), itemName, claims.Username)

		if err != nil {
			if errors.Is(err, service.ErrEmployeeNotFound) {
				log.Info("failed to buy item", sl.Err(err))

				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, resp.ErrorResponse{Errors: "employee not found"})
				return
			}

			if errors.Is(err, service.ErrNotEnoughCoins) {
				log.Info("failed to buy item", sl.Err(err))

				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.ErrorResponse{Errors: "not enough coins"})
				return
			}

			if errors.Is(err, service.ErrItemNotFound) {
				log.Info("failed to buy item", sl.Err(err))

				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.ErrorResponse{Errors: "item not found"})
				return
			}

			log.Error("failed to send coins", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.ErrorResponse{Errors: "internal server error"})
			return
		}

		render.Status(r, http.StatusOK)
	}
}
