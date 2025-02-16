package handlers

import (
	resp "avito-shop/internal/http-server/dto/response"
	mw "avito-shop/internal/http-server/middleware"
	"avito-shop/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

const (
	internalServerError = "internal server error"
)

func renderError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, resp.ErrorResponse{Errors: message})
}

func getClaimsFromContext(r *http.Request, log *slog.Logger) (*service.TokenClaims, bool) {
	claims, ok := r.Context().Value(mw.UserContextKey).(*service.TokenClaims)
	if !ok || claims == nil {
		log.Error("failed to get claims from context")
		return nil, false
	}
	return claims, true
}

func getURLParam(r *http.Request, param string, log *slog.Logger) (string, bool) {
	value := chi.URLParam(r, param)
	if value == "" {
		log.Error("empty parameter", slog.String("param", param))
		return "", false
	}
	return value, true
}

func setupLogger(log *slog.Logger, op string, r *http.Request) *slog.Logger {
	return log.With(
		slog.String("operation", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
}
