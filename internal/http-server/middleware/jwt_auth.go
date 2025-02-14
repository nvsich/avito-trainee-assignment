package middleware

import (
	"avito-shop/internal/http-server/dto/response"
	"avito-shop/internal/lib/logger/sl"
	"avito-shop/internal/service"
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const UserContextKey contextKey = "user"

func NewJwtAuth(log *slog.Logger, signKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(slog.String("component", "middleware/jwt_auth"))

		fn := func(w http.ResponseWriter, r *http.Request) {
			requestId := middleware.GetReqID(r.Context())
			const requestIdKey = "request_id"

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Error("missing Authorization header", slog.String(requestIdKey, requestId))

				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, &response.ErrorResponse{Errors: "missing auth header"})
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &service.TokenClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					log.Error("unexpected signing method",
						slog.String(requestIdKey, requestId),
						slog.String("method", token.Method.Alg()),
					)

					return nil, errors.New("unexpected signing method")
				}

				return []byte(signKey), nil
			})

			if err != nil || !token.Valid {
				log.Error("invalid token", slog.String(requestIdKey, requestId), sl.Err(err))

				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, &response.ErrorResponse{Errors: "invalid token"})
				return
			}

			if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
				log.Error("expired token", slog.String(requestIdKey, requestId))

				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, &response.ErrorResponse{Errors: "token expired"})
				return
			}

			log.Info("successful authentication",
				slog.String("user_id", claims.Username),
				slog.String(requestIdKey, requestId),
			)

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
