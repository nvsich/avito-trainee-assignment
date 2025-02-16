package handlers

import (
	"avito-shop/internal/http-server/dto"
	"avito-shop/internal/lib/logger/sl"
	"avito-shop/internal/model"
	"avito-shop/internal/service"
	"context"
	"errors"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Info interface {
	Get(ctx context.Context, username string) (*model.EmployeeInfo, error)
}

func NewInfoHandlerFunc(log *slog.Logger, infoService Info) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.NewInfoHandlerFunc"
		log = setupLogger(log, op, r)

		claims, ok := getClaimsFromContext(r, log)
		if !ok {
			renderError(w, r, http.StatusInternalServerError, "internal server error")
			return
		}

		employeeInfo, err := infoService.Get(r.Context(), claims.Username)
		if err != nil {
			handleInfoError(w, r, log, err)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, dto.ToInfoResponse(*employeeInfo))
	}
}

func handleInfoError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	if errors.Is(err, service.ErrEmployeeNotFound) {
		log.Error("Employee not found", sl.Err(err))
		renderError(w, r, http.StatusUnauthorized, "employee not found")
		return
	}

	log.Error("Info retrieval failed", sl.Err(err))
	renderError(w, r, http.StatusInternalServerError, "internal server error")
}
