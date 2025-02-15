package handlers

import (
	req "avito-shop/internal/http-server/dto/request"
	resp "avito-shop/internal/http-server/dto/response"
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

type Auth interface {
	Authorize(ctx context.Context, username string, password string) (string, error)
}

func NewAuthHandlerFunc(log *slog.Logger, authService Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.NewAuthHandlerFunc"

		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var request req.AuthRequest

		err := render.DecodeJSON(r.Body, &request)
		if err != nil {
			log.Error("failed to parse request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ErrorResponse{Errors: "failed to parse request"})
			return
		}

		log.Debug("request body successfully decoded")

		if err = validator.New().Struct(request); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("invalid request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ErrorResponse{Errors: "invalid request body"})
			return
		}

		log.Debug("validation passed", slog.String("username", request.Username))

		token, err := authService.Authorize(r.Context(), request.Username, request.Password)
		if err != nil {
			if errors.Is(err, service.ErrInvalidCredentials) {
				log.Info("invalid login attempt", slog.String("username", request.Username))

				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, resp.ErrorResponse{Errors: "invalid credentials"})
				return
			}

			log.Error("failed to authorize", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.ErrorResponse{Errors: "internal error"})
			return
		}

		log.Info("user successfully authenticated", slog.String("username", request.Username))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp.AuthResponse{Token: token})
	}
}
