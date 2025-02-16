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

type PgInventoryRepo struct {
	*Postgres
	getter *trmpgx.CtxGetter
}

func NewPgInventoryRepo(p *Postgres, c *trmpgx.CtxGetter) *PgInventoryRepo {
	return &PgInventoryRepo{p, c}
}

func (r *PgInventoryRepo) Save(ctx context.Context, inventory *model.EmployeeInventory) error {
	const op = "repo.pgdb.PGInventoryRepo.Save"

	query, args, err := r.Builder.
		Insert("employee_inventory").
		Columns("id, employee_id, item_id, amount").
		Values(inventory.Id, inventory.EmployeeId, inventory.ItemId, inventory.Amount).
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

func (r *PgInventoryRepo) FindAllInventoryItemsByEmployee(
	ctx context.Context, employeeId uuid.UUID) ([]model.InventoryItem, error) {
	const op = "repo.pgdb.PGInventoryRepo.FindAllInventoryItemsByEmployee"

	query, args, err := r.Builder.
		Select("items.name, sum(employee_inventory.amount) as amount").
		From("employee_inventory").
		LeftJoin("items on items.id = employee_inventory.item_id").
		Where("employee_inventory.employee_id = ?", employeeId).
		GroupBy("items.name").
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

	var items []model.InventoryItem
	for rows.Next() {
		var item model.InventoryItem
		err = rows.Scan(&item.Type, &item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return items, nil
}

func (r *PgInventoryRepo) FindByEmployeeAndItem(
	ctx context.Context, employeeId uuid.UUID, itemId uuid.UUID) (*model.EmployeeInventory, error) {
	const op = "repo.pgdb.PGInventoryRepo.FindByEmployeeAndItem"

	query, args, err := r.Builder.
		Select("id, employee_id, item_id, amount").
		From("employee_inventory").
		Where("employee_id = ? AND item_id = ?", employeeId, itemId).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	conn := r.getter.DefaultTrOrDB(ctx, r.Pool)

	var inventory model.EmployeeInventory
	err = conn.QueryRow(ctx, query, args...).Scan(
		&inventory.Id,
		&inventory.EmployeeId,
		&inventory.ItemId,
		&inventory.Amount,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrEmployeeInventoryNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &inventory, nil
}

func (r *PgInventoryRepo) UpdateById(
	ctx context.Context, id uuid.UUID, employeeInventory *model.EmployeeInventory) error {
	const op = "repo.pgdb.PGInventoryRepo.UpdateById"

	query, args, err := r.Builder.
		Update("employee_inventory").
		Set("employee_id", employeeInventory.EmployeeId).
		Set("item_id", employeeInventory.ItemId).
		Set("amount", employeeInventory.Amount).
		Where("id = ?", id).
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
