package pgdb

import (
	"avito-shop/internal/model"
	"avito-shop/internal/repo"
	"context"
	"errors"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PGItemRepo struct {
	*Postgres
	getter *trmpgx.CtxGetter
}

func NewPGItemRepo(p *Postgres, c *trmpgx.CtxGetter) *PGItemRepo {
	return &PGItemRepo{p, c}
}

func (r *PGItemRepo) FindByName(ctx context.Context, itemName string) (*model.Item, error) {
	const op = "repo.pgdb.PGItemRepo.FindByName"

	query, args, err := r.Builder.
		Select("id, name, price").
		From("items").
		Where("name = ?", itemName).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.Pool)

	var item model.Item
	err = conn.QueryRow(ctx, query, args...).
		Scan(&item.Id, &item.Name, &item.Price)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrItemNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &item, nil
}

func (r *PGItemRepo) FindById(ctx context.Context, itemId uuid.UUID) (*model.Item, error) {
	const op = "repo.pgdb.PGItemRepo.FindById"

	query, args, err := r.Builder.
		Select("id, name, price").
		From("items").
		Where("id = ?", itemId).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.Pool)

	var item model.Item
	err = conn.QueryRow(ctx, query, args...).
		Scan(&item.Id, &item.Name, &item.Price)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrItemNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &item, nil
}
