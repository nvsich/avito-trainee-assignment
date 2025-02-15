package pgdb

import (
	"avito-shop/internal/model"
	"avito-shop/internal/repo"
	"context"
	"errors"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PGEmployeeRepo struct {
	*Postgres
	getter *trmpgx.CtxGetter
}

func NewPGEmployeeRepo(p *Postgres, c *trmpgx.CtxGetter) *PGEmployeeRepo {
	return &PGEmployeeRepo{p, c}
}

func (r *PGEmployeeRepo) Save(ctx context.Context, employee *model.Employee) error {
	const op = "repo.pgdb.PGEmployeeRepo.Save"

	query, args, err := r.Builder.
		Insert("employees").
		Columns("id, username, password_hash, balance").
		Values(employee.Id, employee.Username, employee.PasswordHash, employee.Balance).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.Pool)

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return repo.ErrEmployeeExists
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PGEmployeeRepo) FindByUsername(ctx context.Context, username string) (*model.Employee, error) {
	const op = "repo.pgdb.PGEmployeeRepo.FindByUsername"

	query, args, err := r.Builder.
		Select("id, username, password_hash, balance").
		From("employees").
		Where("username = ?", username).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.Pool)

	var employee model.Employee
	err = conn.QueryRow(ctx, query, args...).
		Scan(
			&employee.Id,
			&employee.Username,
			&employee.PasswordHash,
			&employee.Balance,
		)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &employee, nil
}

func (r *PGEmployeeRepo) UpdateByUsername(ctx context.Context, username string, employee *model.Employee) error {
	const op = "repo.pgdb.PGEmployeeRepo.UpdateByUsername"

	query, args, err := r.Builder.
		Update("employees").
		Set("password_hash", employee.PasswordHash).
		Set("balance", employee.Balance).
		Where("username = ?", username).
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
