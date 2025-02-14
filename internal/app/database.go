package app

import (
	"avito-shop/internal/config"
	"avito-shop/internal/lib/logger/sl"
	"avito-shop/internal/repo/pgdb"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"log/slog"
)

func mustSetupDatabase(cfg *config.Config, log *slog.Logger) (*pgdb.Postgres, *manager.Manager) {
	pg, err := pgdb.New(cfg.PG.URL, cfg.MaxPoolSize)
	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
		panic(err)
	}

	trManager := manager.Must(trmpgx.NewDefaultFactory(pg.Pool))

	return pg, trManager
}
