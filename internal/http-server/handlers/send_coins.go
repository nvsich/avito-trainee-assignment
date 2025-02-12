package handlers

import (
	req "avito-shop/internal/http-server/dto/request"
	resp "avito-shop/internal/http-server/dto/response"
	mw "avito-shop/internal/http-server/middleware"
	"avito-shop/internal/lib/logger/sl"
	"avito-shop/internal/service"
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Transfer interface {
	SendCoins(ctx context.Context, from string, to string, amount int) error
}

func NewSendCoinsHandlerFunc(log *slog.Logger, transferService Transfer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.NewSendCoinsHandlerFunc"

		log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var request req.SendCoinRequest

		err := render.DecodeJSON(r.Body, &request)
		if err != nil {
			log.Error("failed to parse request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ErrorResponse{Errors: "failed to parse request"})
			return
		}

		if err = validator.New().Struct(request); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("failed to validate request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ErrorResponse{Errors: "failed to validate request"})
			return
		}

		log.Debug("validation passed")

		claims, ok := r.Context().Value(mw.UserContextKey).(*service.TokenClaims)
		if !ok || claims == nil {
			log.Error("failed to get claims from context")

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.ErrorResponse{Errors: "internal server error"})
			return
		}

		err = transferService.SendCoins(r.Context(), claims.Login, request.ToUser, request.Amount)
		if err != nil {
			if errors.Is(err, service.ErrNotEnoughCoins) {
				log.Info("failed to send coins", sl.Err(err))

				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.ErrorResponse{Errors: "not enough coins to send"})
				return
			}

			if errors.Is(err, service.ErrNegativeTransferAmount) {
				log.Info("failed to send coins", sl.Err(err))

				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.ErrorResponse{Errors: "negative amount"})
				return
			}

			if errors.Is(err, service.ErrReceiverNotFound) {
				log.Info("failed to send coins", sl.Err(err))

				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.ErrorResponse{Errors: "receiver not found"})
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
