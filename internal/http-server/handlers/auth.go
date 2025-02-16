package handlers

import (
	req "avito-shop/internal/http-server/dto/request"
	resp "avito-shop/internal/http-server/dto/response"
	"avito-shop/internal/lib/logger/sl"
	"avito-shop/internal/service"
	"context"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Auth interface {
	Authorize(ctx context.Context, username string, password string) (string, error)
}

func NewAuthHandlerFunc(log *slog.Logger, authService Auth, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.NewAuthHandlerFunc"
		log = setupLogger(log, op, r)

		var request req.AuthRequest

		if err := render.DecodeJSON(r.Body, &request); err != nil {
			log.Error("failed to parse request", sl.Err(err))
			renderError(w, r, http.StatusBadRequest, "failed to parse request")
			return
		}

		log.Debug("Request body decoded")

		if err := validate.Struct(request); err != nil {
			log.Error("Invalid request", sl.Err(err))
			renderError(w, r, http.StatusBadRequest, "invalid request body")
			return
		}

		log.Debug("Validation passed", slog.String("username", request.Username))

		token, err := authService.Authorize(r.Context(), request.Username, request.Password)
		if err != nil {
			if errors.Is(err, service.ErrInvalidCredentials) {
				log.Info("Invalid login attempt", slog.String("username", request.Username))
				renderError(w, r, http.StatusUnauthorized, "invalid credentials")
				return
			}

			log.Error("Authorization failed", sl.Err(err))
			renderError(w, r, http.StatusInternalServerError, "internal error")
			return
		}

		log.Info("User authenticated", slog.String("username", request.Username))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp.AuthResponse{Token: token})
	}
}
