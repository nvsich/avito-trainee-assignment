package handlers

import (
	req "avito-shop/internal/http-server/dto/request"
	"avito-shop/internal/lib/logger/sl"
	"avito-shop/internal/service"
	"context"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Transfer interface {
	SendCoins(ctx context.Context, from string, to string, amount int) error
}

func NewSendCoinsHandlerFunc(log *slog.Logger, transferService Transfer, vld *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.NewSendCoinsHandlerFunc"
		log = setupLogger(log, op, r)

		var request req.SendCoinRequest

		if err := render.DecodeJSON(r.Body, &request); err != nil {
			log.Error("Failed to parse request", sl.Err(err))
			renderError(w, r, http.StatusBadRequest, "failed to parse request")
			return
		}

		if err := vld.Struct(request); err != nil {
			log.Error("Invalid request", sl.Err(err))
			renderError(w, r, http.StatusBadRequest, "invalid request body")
			return
		}

		claims, ok := getClaimsFromContext(r, log)
		if !ok {
			renderError(w, r, http.StatusInternalServerError, internalServerError)
			return
		}

		if err := transferService.SendCoins(r.Context(), claims.Username, request.ToUser, request.Amount); err != nil {
			handleTransferError(w, r, log, err)
			return
		}

		render.Status(r, http.StatusOK)
	}
}

func handleTransferError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	var status int
	var message string

	switch {
	case errors.Is(err, service.ErrTransferToSameEmployee):
		status, message = http.StatusBadRequest, "can't send coins to yourself"
	case errors.Is(err, service.ErrNotEnoughCoins):
		status, message = http.StatusBadRequest, "not enough coins to send"
	case errors.Is(err, service.ErrNegativeTransferAmount):
		status, message = http.StatusBadRequest, "negative amount"
	case errors.Is(err, service.ErrReceiverNotFound):
		status, message = http.StatusBadRequest, "receiver not found"
	case errors.Is(err, service.ErrSenderNotFound):
		status, message = http.StatusInternalServerError, internalServerError
	default:
		status, message = http.StatusInternalServerError, internalServerError
		log.Error("Transfer failed", sl.Err(err))
	}

	if status != http.StatusInternalServerError {
		log.Info("Transfer failed", sl.Err(err))
	}

	renderError(w, r, status, message)
}
