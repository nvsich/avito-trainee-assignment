package handlers

import (
	"avito-shop/internal/http-server/dto"
	resp "avito-shop/internal/http-server/dto/response"
	mw "avito-shop/internal/http-server/middleware"
	"avito-shop/internal/lib/logger/sl"
	"avito-shop/internal/model"
	"avito-shop/internal/service"
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
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

		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		claims, ok := r.Context().Value(mw.UserContextKey).(*service.TokenClaims)
		if !ok || claims == nil {
			log.Error("failed to get claims from context")

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.ErrorResponse{Errors: "internal server error"})
			return
		}

		employeeInfo, err := infoService.Get(r.Context(), claims.Username)

		if err != nil {
			if errors.Is(err, service.ErrEmployeeNotFound) {
				log.Error("failed to get employee info", sl.Err(err))

				// TODO: подумать над логикой тут и в других хэндлерах
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.ErrorResponse{Errors: "employee not found"})
				return
			}

			log.Error("failed to get employee info", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.ErrorResponse{Errors: "internal server error"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, dto.ToInfoResponse(*employeeInfo))
	}
}
