package repo

import (
	"avito-shop/internal/model"
	"avito-shop/internal/repo"
	"avito-shop/internal/repo/pgdb"
	"context"
	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PGInventoryRepoTestSuite struct {
	PGDBTestSuite
	ctx           context.Context
	inventoryRepo *pgdb.PgInventoryRepo
}

func (s *PGInventoryRepoTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.inventoryRepo = pgdb.NewPgInventoryRepo(
		&pgdb.Postgres{
			Pool:    s.pool,
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		},
		trmpgx.DefaultCtxGetter,
	)

	_, err := s.pool.Exec(context.Background(),
		`truncate table employees restart identity cascade;
              truncate table items restart identity cascade;
              truncate table employee_inventory restart identity cascade;`)
	s.Require().NoError(err)
}

func TestPGInventoryRepo(t *testing.T) {
	suite.Run(t, new(PGInventoryRepoTestSuite))
}

func (s *PGInventoryRepoTestSuite) TestSave() {
	employeeId := uuid.New()
	itemId := uuid.New()

	s.insertEmployee(employeeId, "employee")
	s.insertItem(itemId, "item")

	testEmployeeInventory := model.EmployeeInventory{
		Id:         uuid.New(),
		EmployeeId: employeeId,
		ItemId:     itemId,
		Amount:     123,
	}

	s.Run("should save employee inventory", func() {
		err := s.inventoryRepo.Save(s.ctx, &testEmployeeInventory)
		s.Require().NoError(err)

		var saved model.EmployeeInventory
		err = s.pool.QueryRow(s.ctx, "select id, employee_id, item_id, amount from employee_inventory").
			Scan(&saved.Id, &saved.EmployeeId, &saved.ItemId, &saved.Amount)
		s.Require().NoError(err)

		s.Require().Equal(testEmployeeInventory.Id, saved.Id)
		s.Require().Equal(testEmployeeInventory.EmployeeId, saved.EmployeeId)
		s.Require().Equal(testEmployeeInventory.ItemId, saved.ItemId)
		s.Require().Equal(testEmployeeInventory.Amount, saved.Amount)
	})
}

func (s *PGInventoryRepoTestSuite) TestFindAllInventoryItemsByEmployee() {
	employeeId := uuid.New()
	itemId1 := uuid.New()
	itemId2 := uuid.New()
	item1Name := "item1"
	item2Name := "item2"
	item1Amount := 10
	item2Amount := 20

	s.insertEmployee(employeeId, "employee")
	s.insertItem(itemId1, item1Name)
	s.insertItem(itemId2, item2Name)
	s.insertInventory(model.EmployeeInventory{Id: uuid.New(), EmployeeId: employeeId, ItemId: itemId1, Amount: item1Amount})
	s.insertInventory(model.EmployeeInventory{Id: uuid.New(), EmployeeId: employeeId, ItemId: itemId2, Amount: item2Amount})

	s.Run("should find all inventory items for employee", func() {
		inventoryItems, err := s.inventoryRepo.FindAllInventoryItemsByEmployee(s.ctx, employeeId)
		s.Require().NoError(err)
		s.Require().Len(inventoryItems, 2)

		expected := map[string]int{
			item1Name: item1Amount,
			item2Name: item2Amount,
		}

		for _, inventoryItem := range inventoryItems {
			s.Require().Equal(expected[inventoryItem.Type], inventoryItem.Quantity)
		}
	})
}

//FindByEmployeeAndItem(ctx context.Context, employeeId uuid.UUID, itemId uuid.UUID) (*model.EmployeeInventory, error)
//	UpdateById(ctx context.Context, id uuid.UUID, employeeInventory *model.EmployeeInventory) error

func (s *PGInventoryRepoTestSuite) TestFindByEmployeeAndItem() {
	employeeId := uuid.New()
	employeeIdWithoutInventory := uuid.New()
	itemId := uuid.New()

	s.insertEmployee(employeeId, "employee")
	s.insertEmployee(employeeIdWithoutInventory, "employee without inventory")
	s.insertItem(itemId, "item")

	testEmployeeInventory := model.EmployeeInventory{
		Id:         uuid.New(),
		EmployeeId: employeeId,
		ItemId:     itemId,
		Amount:     123,
	}

	s.insertInventory(testEmployeeInventory)

	s.Run("should find all inventory items for employee", func() {
		employeeInventory, err := s.inventoryRepo.FindByEmployeeAndItem(s.ctx, employeeId, itemId)
		s.Require().NoError(err)
		s.Require().Equal(testEmployeeInventory.Id, employeeInventory.Id)
		s.Require().Equal(testEmployeeInventory.EmployeeId, employeeInventory.EmployeeId)
		s.Require().Equal(testEmployeeInventory.Amount, employeeInventory.Amount)
	})

	s.Run("should return ErrImployeeInventory not found", func() {
		employeeInventory, err := s.inventoryRepo.FindByEmployeeAndItem(s.ctx, employeeIdWithoutInventory, itemId)
		s.Require().ErrorIs(err, repo.ErrEmployeeInventoryNotFound)
		s.Require().Nil(employeeInventory)
	})
}

func (s *PGInventoryRepoTestSuite) insertEmployee(employeeId uuid.UUID, username string) {
	_, err := s.pool.Exec(s.ctx,
		"insert into employees (id, username, password_hash, balance) VALUES ($1, $2, 'hash', 1000)",
		employeeId, username)
	s.Require().NoError(err)
}

func (s *PGInventoryRepoTestSuite) insertItem(itemId uuid.UUID, name string) {
	_, err := s.pool.Exec(s.ctx,
		"insert into items (id, name, price) VALUES ($1, $2, 123)",
		itemId, name)
	s.Require().NoError(err)
}

func (s *PGInventoryRepoTestSuite) insertInventory(t model.EmployeeInventory) {
	_, err := s.pool.Exec(s.ctx,
		"insert into employee_inventory (id, employee_id, item_id, amount) VALUES ($1, $2, $3, $4)",
		t.Id, t.EmployeeId, t.ItemId, t.Amount)
	s.Require().NoError(err)
}
