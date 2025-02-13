package pgdb

import (
	"avito-shop/internal/model"
	"context"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
)

type PGTransferRepo struct {
	*Postgres
	getter *trmpgx.CtxGetter
}

func NewPGTransferRepo(p *Postgres, c *trmpgx.CtxGetter) *PGTransferRepo {
	return &PGTransferRepo{p, c}
}

func (r *PGTransferRepo) Save(ctx context.Context, transfer *model.Transfer) error {
	const op = "repo.PGTransferRepo.Save"

	query, args, err := r.Builder.
		Insert("transfers").
		Columns("id, from_employee, to_employee, amount").
		Values(transfer.Id, transfer.FromEmployee, transfer.ToEmployee, transfer.Amount).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.Pool)

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
