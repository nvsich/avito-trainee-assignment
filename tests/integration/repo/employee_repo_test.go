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

type PGEmployeeRepoTestSuite struct {
	PGDBTestSuite
	ctx          context.Context
	employeeRepo *pgdb.PGEmployeeRepo
}

func (s *PGEmployeeRepoTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.employeeRepo = pgdb.NewPGEmployeeRepo(
		&pgdb.Postgres{
			Pool:    s.pool,
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		},
		trmpgx.DefaultCtxGetter,
	)

	_, err := s.pool.Exec(context.Background(), "truncate table employees restart identity cascade")
	s.Require().NoError(err)
}

func TestPGEmployeeRepo(t *testing.T) {
	suite.Run(t, new(PGEmployeeRepoTestSuite))
}

func (s *PGEmployeeRepoTestSuite) TestSave() {
	testEmployee := model.Employee{
		Id:           uuid.New(),
		Username:     "test username",
		PasswordHash: "test passwordHash",
		Balance:      1,
	}

	s.Run("should save employee", func() {
		err := s.employeeRepo.Save(s.ctx, &testEmployee)
		s.Require().NoError(err)

		var savedEmployee model.Employee
		err = s.pool.QueryRow(s.ctx, "select id, username, password_hash, balance from employees where id = $1", testEmployee.Id).
			Scan(&savedEmployee.Id, &savedEmployee.Username, &savedEmployee.PasswordHash, &savedEmployee.Balance)
		s.Require().NoError(err)

		s.Require().Equal(testEmployee.Id, savedEmployee.Id)
		s.Require().Equal(testEmployee.Username, savedEmployee.Username)
		s.Require().Equal(testEmployee.PasswordHash, savedEmployee.PasswordHash)
		s.Require().Equal(testEmployee.Balance, savedEmployee.Balance)
	})

	s.Run("should fail on duplicate username", func() {
		duplicateEmployee := model.Employee{
			Id:           uuid.New(),
			Username:     testEmployee.Username,
			PasswordHash: "another_password_hash",
			Balance:      50,
		}

		err := s.employeeRepo.Save(s.ctx, &duplicateEmployee)
		s.Require().ErrorIs(err, repo.ErrEmployeeExists)
	})
}

func (s *PGEmployeeRepoTestSuite) TestFindByUsername() {
	testEmployee := model.Employee{
		Id:           uuid.New(),
		Username:     "test username",
		PasswordHash: "test passwordHash",
		Balance:      1,
	}

	_, err := s.pool.Exec(s.ctx,
		"insert into employees(id, username, password_hash, balance) values ($1, $2, $3, $4)",
		testEmployee.Id, testEmployee.Username, testEmployee.PasswordHash, testEmployee.Balance)
	s.Require().NoError(err)

	s.Run("should find employee by username", func() {
		employee, err := s.employeeRepo.FindByUsername(s.ctx, testEmployee.Username)
		s.Require().NoError(err)
		s.Require().NotNil(employee)
		s.Require().Equal(testEmployee.Id, employee.Id)
		s.Require().Equal(testEmployee.Username, employee.Username)
		s.Require().Equal(testEmployee.PasswordHash, employee.PasswordHash)
		s.Require().Equal(testEmployee.Balance, employee.Balance)
	})

	s.Run("should not find employee by username", func() {
		employee, err := s.employeeRepo.FindByUsername(s.ctx, "non existing username")
		s.Require().ErrorIs(err, repo.ErrEmployeeNotFound)
		s.Require().Nil(employee)
	})
}

func (s *PGEmployeeRepoTestSuite) TestUpdateByUsername() {
	testEmployee := model.Employee{
		Id:           uuid.New(),
		Username:     "test username",
		PasswordHash: "test passwordHash",
		Balance:      1,
	}

	_, err := s.pool.Exec(s.ctx,
		"insert into employees(id, username, password_hash, balance) values ($1, $2, $3, $4)",
		testEmployee.Id, testEmployee.Username, testEmployee.PasswordHash, testEmployee.Balance)
	s.Require().NoError(err)

	s.Run("should update employee", func() {
		newPasswordHash := "new passwordHash"
		newBalance := 2
		updatedEmployee := &model.Employee{
			Id:           testEmployee.Id,
			Username:     testEmployee.Username,
			PasswordHash: newPasswordHash,
			Balance:      newBalance,
		}

		err = s.employeeRepo.UpdateByUsername(s.ctx, testEmployee.Username, updatedEmployee)
		s.Require().NoError(err)

		var savedEmployee model.Employee
		err = s.pool.QueryRow(s.ctx, "select id, username, password_hash, balance from employees where id = $1", testEmployee.Id).
			Scan(&savedEmployee.Id, &savedEmployee.Username, &savedEmployee.PasswordHash, &savedEmployee.Balance)
		s.Require().NoError(err)

		s.Require().Equal(updatedEmployee.Id, savedEmployee.Id)
		s.Require().Equal(updatedEmployee.Username, savedEmployee.Username)
		s.Require().Equal(updatedEmployee.PasswordHash, newPasswordHash)
		s.Require().Equal(updatedEmployee.Balance, newBalance)
	})

	s.Run("should not fail on not found username", func() {
		newPasswordHash := "new passwordHash"
		newBalance := 2
		updatedEmployee := &model.Employee{
			Id:           testEmployee.Id,
			Username:     testEmployee.Username,
			PasswordHash: newPasswordHash,
			Balance:      newBalance,
		}

		err = s.employeeRepo.UpdateByUsername(s.ctx, "non existing id", updatedEmployee)
		s.Require().NoError(err)
	})
}
