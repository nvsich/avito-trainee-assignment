package pgdb

import (
	"avito-shop/internal/model"
	"avito-shop/internal/repo"
	"context"
	"errors"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
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
		Columns("id, login, password_hash, balance").
		Values(employee.Id, employee.Login, employee.PasswordHash, employee.Balance).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PGEmployeeRepo) FindByLogin(ctx context.Context, login string) (*model.Employee, error) {
	const op = "repo.pgdb.PGEmployeeRepo.FindByLogin"

	query, args, err := r.Builder.
		Select("id, login, password_hash, balance").
		From("employees").
		Where("login = ?", login).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.Pool)

	var employee model.Employee
	err = conn.QueryRow(ctx, query, args...).
		Scan(
			&employee.Id,
			&employee.Login,
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

func (r *PGEmployeeRepo) UpdateByLogin(ctx context.Context, login string, employee *model.Employee) error {
	const op = "repo.pgdb.PGEmployeeRepo.UpdateByLogin"

	query, args, err := r.Builder.
		Update("employees").
		Set("password_hash", employee.PasswordHash).
		Set("balance", employee.Balance).
		Where("login = ?", login).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.Pool)

	cmdTag, err := conn.Exec(ctx, query, args...)

	if cmdTag.RowsAffected() == 0 {
		return repo.ErrNoChange
	}

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
