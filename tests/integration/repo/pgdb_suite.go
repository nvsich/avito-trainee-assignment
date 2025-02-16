package repo

import (
	"avito-shop/tests/setup"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PGDBTestSuite struct {
	suite.Suite
	pool    *pgxpool.Pool
	cleanup func()
}

func (s *PGDBTestSuite) SetupSuite() {
	s.pool, s.cleanup = setup.TestPostgres(s.T())
}

func (s *PGDBTestSuite) TearDownSuite() {
	s.cleanup()
}

func TestPGDBTestSuite(t *testing.T) {
	suite.Run(t, new(PGDBTestSuite))
}
