package app

import (
	"avito-shop/internal/config"
	"avito-shop/internal/http-server/handlers"
	mw "avito-shop/internal/http-server/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

func setupRouter(cfg *config.Config, log *slog.Logger, services *serviceProvider) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mw.NewLogger(log))

	router.Post("/api/auth", handlers.NewAuthHandlerFunc(log, services.AuthService))
	router.Group(func(router chi.Router) {
		router.Use(mw.NewJwtAuth(log, cfg.JWT.SignKey))
		router.Post("/api/sendCoin", handlers.NewSendCoinsHandlerFunc(log, services.TransferService))
		router.Get("/api/buy/{item}", handlers.NewBuyItemHandlerFunc(log, services.BuyItemService))
		router.Get("/api/info", handlers.NewInfoHandlerFunc(log, services.InfoService))
	})

	return router
}
