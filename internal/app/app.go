package app

import (
	"avito-shop/internal/config"
	"avito-shop/internal/lib/logger/sl"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(envPath string) {
	cfg := config.MustLoad(envPath)
	log := mustSetupLogger(cfg.Log.Level)
	pg, trManager := mustSetupDatabase(cfg, log)
	defer pg.Close()

	services := newServiceProvider(cfg, pg, trManager)
	router := setupRouter(cfg, log, services)
	server := setupServer(cfg, router)

	// TODO: move running server and shutdown to different method

	go func() {
		log.Info("starting server", slog.String("addr", server.Addr))

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown server", sl.Err(err))
	}

	select {
	case <-ctx.Done():
		log.Info("context done, shutting down server...")
	}

	log.Info("server stopped")
}
