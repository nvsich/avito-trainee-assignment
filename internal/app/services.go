package app

import (
	"avito-shop/internal/config"
	"avito-shop/internal/repo/pgdb"
	"avito-shop/internal/service"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type serviceProvider struct {
	AuthService     *service.AuthService
	TransferService *service.TransferService
	BuyItemService  *service.ItemService
	InfoService     *service.InfoService
}

func newServiceProvider(cfg *config.Config, pg *pgdb.Postgres, trManager *manager.Manager) *serviceProvider {
	pgEmployeeRepo := pgdb.NewPGEmployeeRepo(pg, trmpgx.DefaultCtxGetter)
	pgTransferRepo := pgdb.NewPGTransferRepo(pg, trmpgx.DefaultCtxGetter)
	pgItemRepo := pgdb.NewPGItemRepo(pg, trmpgx.DefaultCtxGetter)
	pgInventoryRepo := pgdb.NewPgInventoryRepo(pg, trmpgx.DefaultCtxGetter)

	return &serviceProvider{
		AuthService:     service.NewAuthService(pgEmployeeRepo, cfg.JWT.SignKey, cfg.JWT.TokenTTL),
		TransferService: service.NewTransferService(trManager, pgEmployeeRepo, pgTransferRepo),
		BuyItemService:  service.NewItemService(trManager, pgItemRepo, pgEmployeeRepo, pgInventoryRepo),
		InfoService:     service.NewInfoService(trManager, pgEmployeeRepo, pgInventoryRepo, pgTransferRepo, pgItemRepo),
	}
}
