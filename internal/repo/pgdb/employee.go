package pgdb

import (
	"avito-shop/internal/model"
	"avito-shop/internal/repo"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PGEmployeeRepo struct {
	*Postgres
}

func NewPGEmployeeRepo(p *Postgres) *PGEmployeeRepo {
	return &PGEmployeeRepo{p}
}

func (r *PGEmployeeRepo) Save(ctx context.Context, employee *model.Employee) error {
	const op = "repo.pgdb.Save"

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
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return repo.ErrEmployeeExists
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PGEmployeeRepo) FindByLogin(ctx context.Context, login string) (*model.Employee, error) {
	const op = "repo.pgdb.FindByLogin"

	query, args, err := r.Builder.
		Select("id, login, password_hash, balance").
		From("employees").
		Where("login = ?", login).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var employee model.Employee
	err = r.Pool.QueryRow(ctx, query, args...).
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
