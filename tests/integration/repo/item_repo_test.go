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

type PGItemRepoTestSuite struct {
	PGDBTestSuite
	ctx      context.Context
	itemRepo *pgdb.PGItemRepo
}

func (s *PGItemRepoTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.itemRepo = pgdb.NewPGItemRepo(
		&pgdb.Postgres{
			Pool:    s.pool,
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		},
		trmpgx.DefaultCtxGetter,
	)

	_, err := s.pool.Exec(context.Background(), "truncate table items restart identity cascade;")
	s.Require().NoError(err)
}

func TestPGItemRepo(t *testing.T) {
	suite.Run(t, new(PGItemRepoTestSuite))
}

func (s *PGItemRepoTestSuite) TestFindByName() {
	testItem := model.Item{
		Id:    uuid.New(),
		Name:  "tests item",
		Price: 100,
	}

	_, err := s.pool.Exec(s.ctx,
		"insert into items(id, name, price) values ($1, $2, $3)",
		testItem.Id, testItem.Name, testItem.Price)
	s.Require().NoError(err)

	s.Run("should find item by name", func() {
		item, err := s.itemRepo.FindByName(s.ctx, testItem.Name)
		s.Require().NoError(err)
		s.Require().NotNil(item)
		s.Require().Equal(testItem.Id, item.Id)
		s.Require().Equal(testItem.Name, item.Name)
		s.Require().Equal(testItem.Price, item.Price)
	})

	s.Run("should not find item by name", func() {
		item, err := s.itemRepo.FindByName(s.ctx, "non existing item")
		s.Require().ErrorIs(err, repo.ErrItemNotFound)
		s.Require().Nil(item)
	})
}

func (s *PGItemRepoTestSuite) TestPGItemRepo_FindById() {
	testItem := model.Item{
		Id:    uuid.New(),
		Name:  "tests item",
		Price: 100,
	}

	_, err := s.pool.Exec(s.ctx,
		"insert into items(id, name, price) values ($1, $2, $3)",
		testItem.Id, testItem.Name, testItem.Price)
	s.Require().NoError(err)

	s.Run("should find item by id", func() {
		item, err := s.itemRepo.FindById(s.ctx, testItem.Id)
		s.Require().NoError(err)
		s.Require().NotNil(item)
		s.Require().Equal(testItem.Id, item.Id)
		s.Require().Equal(testItem.Name, item.Name)
		s.Require().Equal(testItem.Price, item.Price)
	})

	s.Run("should not find item by id", func() {
		nonExistingId := uuid.New()
		item, err := s.itemRepo.FindById(s.ctx, nonExistingId)
		s.Require().ErrorIs(err, repo.ErrItemNotFound)
		s.Require().Nil(item)
	})
}
