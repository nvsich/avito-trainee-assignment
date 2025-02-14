package pgdb

import (
	"avito-shop/internal/model"
	"context"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
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

func (r *PGTransferRepo) FindAllForReceiverGroupedBySenders(
	ctx context.Context, receiverId uuid.UUID) ([]model.CoinTransaction, error) {
	const op = "repo.PGTransferRepo.FindAllForReceiverGroupedBySenders"

	query, args, err := r.Builder.
		Select("e.login as user, sum(t.amount) as amount").
		From("transfers t").
		Join("employees e on t.from_employee = e.id").
		Where("t.to_employee = ?", receiverId).
		GroupBy("e.login").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.Pool)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var transactions []model.CoinTransaction
	for rows.Next() {
		var t model.CoinTransaction
		err = rows.Scan(&t.User, &t.Amount)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return transactions, nil
}

func (r *PGTransferRepo) FindAllForSenderGroupedByReceivers(
	ctx context.Context, senderId uuid.UUID) ([]model.CoinTransaction, error) {
	const op = "repo.PGTransferRepo.FindAllForReceiverGroupedBySenders"

	query, args, err := r.Builder.
		Select("e.login as user, sum(t.amount) as amount").
		From("transfers t").
		Join("employees e on t.to_employee = e.id").
		Where("t.from_employee = ?", senderId).
		GroupBy("e.login").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.Pool)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var transactions []model.CoinTransaction
	for rows.Next() {
		var t model.CoinTransaction
		err = rows.Scan(&t.User, &t.Amount)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return transactions, nil
}
