package app

import (
	"avito-shop/internal/config"
	"avito-shop/internal/http-server/handlers"
	mw "avito-shop/internal/http-server/middleware"
	"avito-shop/internal/lib/logger/sl"
	"avito-shop/internal/repo/pgdb"
	"avito-shop/internal/service"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// TODO: decompose for several steps

func Run(configPath string) {
	cfg := config.MustLoad(configPath)

	log := mustSetupLogger(cfg.Log.Level)

	pg, err := pgdb.New(cfg.PG.URL)
	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
		os.Exit(1)
	}
	defer pg.Close()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mw.NewLogger(log))

	pgEmployeeRepo := pgdb.NewPGEmployeeRepo(pg)
	authService := service.NewAuthService(pgEmployeeRepo, cfg.JWT.SignKey, cfg.JWT.TokenTTL)

	router.Post("/api/auth", handlers.NewAuthHandlerFunc(log, authService))
	router.Group(func(router chi.Router) {
		router.Use(mw.NewJwtAuth(log, cfg.JWT.SignKey))

		// router.Post("/api/path...) to protect
	})

	log.Info("starting server", slog.String("port", cfg.HTTP.Port))

	server := &http.Server{
		Addr:         ":" + cfg.HTTP.Port, // TODO: add addr to config
		Handler:      router,
		ReadTimeout:  5 * time.Second,   // TODO: add to config
		WriteTimeout: 10 * time.Second,  // TODO: add to config
		IdleTimeout:  120 * time.Second, // TODO: add to config
	}

	log.Info("server started", slog.String("address", server.Addr))

	if err = server.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

	// TODO: add graceful shutdown

	log.Error("stopped server", sl.Err(server.Shutdown(context.Background())))
}
