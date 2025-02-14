package app

import (
	"avito-shop/internal/config"
	"avito-shop/internal/http-server/handlers"
	mw "avito-shop/internal/http-server/middleware"
	"avito-shop/internal/lib/logger/sl"
	"avito-shop/internal/repo/pgdb"
	"avito-shop/internal/service"
	"context"
	"errors"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TODO: decompose for several steps

func Run(configPath string) {
	cfg := config.MustLoad(configPath)

	log := mustSetupLogger(cfg.Log.Level)

	pg, err := pgdb.New(cfg.PG.URL, cfg.MaxPoolSize)
	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
		os.Exit(1)
	}
	defer pg.Close()

	trManager := manager.Must(trmpgx.NewDefaultFactory(pg.Pool))
	pgEmployeeRepo := pgdb.NewPGEmployeeRepo(pg, trmpgx.DefaultCtxGetter)
	pgTransferRepo := pgdb.NewPGTransferRepo(pg, trmpgx.DefaultCtxGetter)
	pgItemRepo := pgdb.NewPGItemRepo(pg, trmpgx.DefaultCtxGetter)
	pgInventoryRepo := pgdb.NewPgInventoryRepo(pg, trmpgx.DefaultCtxGetter)

	authService := service.NewAuthService(pgEmployeeRepo, cfg.JWT.SignKey, cfg.JWT.TokenTTL)
	transferService := service.NewTransferService(trManager, pgEmployeeRepo, pgTransferRepo)
	buyItemService := service.NewItemService(trManager, pgItemRepo, pgEmployeeRepo, pgInventoryRepo)
	infoService := service.NewInfoService(trManager, pgEmployeeRepo, pgInventoryRepo, pgTransferRepo, pgItemRepo)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mw.NewLogger(log))

	router.Post("/api/auth", handlers.NewAuthHandlerFunc(log, authService))
	router.Group(func(router chi.Router) {
		router.Use(mw.NewJwtAuth(log, cfg.JWT.SignKey))
		router.Post("/api/sendCoin", handlers.NewSendCoinsHandlerFunc(log, transferService))
		router.Get("/api/buy/{item}", handlers.NewBuyItemHandlerFunc(log, buyItemService))
		router.Get("/api/info", handlers.NewInfoHandlerFunc(log, infoService))
	})

	server := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	go func() {
		log.Info("starting server", slog.String("addr", server.Addr))

		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown server", sl.Err(err))
	}

	select {
	case <-ctx.Done():
		log.Info("context done, shutting down server...")
	}

	log.Info("server stopped")
}
