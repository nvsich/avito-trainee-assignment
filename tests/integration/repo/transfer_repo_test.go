package repo

import (
	"avito-shop/internal/model"
	"avito-shop/internal/repo/pgdb"
	"context"
	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PGTransferRepoTestSuite struct {
	PGDBTestSuite
	ctx          context.Context
	transferRepo *pgdb.PGTransferRepo
}

func (s *PGTransferRepoTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.transferRepo = pgdb.NewPGTransferRepo(
		&pgdb.Postgres{
			Pool:    s.pool,
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		},
		trmpgx.DefaultCtxGetter,
	)

	_, err := s.pool.Exec(s.ctx,
		`truncate table transfers restart identity cascade;
		      truncate table employees restart identity cascade;`)
	s.Require().NoError(err)
}

func TestPGTransferRepo(t *testing.T) {
	suite.Run(t, new(PGTransferRepoTestSuite))
}

func (s *PGTransferRepoTestSuite) TestSave() {
	senderId := uuid.New()
	receiverId := uuid.New()

	testTransfer := model.Transfer{
		Id:           uuid.New(),
		FromEmployee: senderId,
		ToEmployee:   receiverId,
		Amount:       1,
	}

	s.insertEmployee(senderId, "sender")
	s.insertEmployee(receiverId, "receiver")

	s.Run("should save transfer", func() {
		err := s.transferRepo.Save(s.ctx, &testTransfer)
		s.Require().NoError(err)

		var savedTransfer model.Transfer
		err = s.pool.QueryRow(s.ctx, "select id, from_employee, to_employee, amount from transfers").
			Scan(&savedTransfer.Id, &savedTransfer.FromEmployee, &savedTransfer.ToEmployee, &savedTransfer.Amount)
		s.Require().NoError(err)

		s.Require().Equal(testTransfer.Id, savedTransfer.Id)
		s.Require().Equal(testTransfer.FromEmployee, savedTransfer.FromEmployee)
		s.Require().Equal(testTransfer.ToEmployee, savedTransfer.ToEmployee)
		s.Require().Equal(testTransfer.Amount, savedTransfer.Amount)
	})
}

func (s *PGTransferRepoTestSuite) TestFindAllForReceiverGroupedBySenders() {
	sender1 := uuid.New()
	sender2 := uuid.New()
	receiver := uuid.New()
	sender1Username := "sender1"
	sender2Username := "sender2"

	s.insertEmployee(sender1, sender1Username)
	s.insertEmployee(sender2, sender2Username)
	s.insertEmployee(receiver, "receiver")

	sender1FirstTransfer := 50
	sender1SecondTransfer := 30
	sender2FirstTransfer := 20

	testTransfers := []model.Transfer{
		{Id: uuid.New(), FromEmployee: sender1, ToEmployee: receiver, Amount: sender1FirstTransfer},
		{Id: uuid.New(), FromEmployee: sender2, ToEmployee: receiver, Amount: sender2FirstTransfer},
		{Id: uuid.New(), FromEmployee: sender1, ToEmployee: receiver, Amount: sender1SecondTransfer},
	}

	for _, transfer := range testTransfers {
		_, err := s.pool.Exec(s.ctx,
			"insert into transfers (id, from_employee, to_employee, amount) values ($1, $2, $3, $4)",
			transfer.Id, transfer.FromEmployee, transfer.ToEmployee, transfer.Amount)

		s.Require().NoError(err)
	}

	s.Run("should return grouped transactions for receiver", func() {
		transactions, err := s.transferRepo.FindAllForReceiverGroupedBySenders(s.ctx, receiver)
		s.Require().NoError(err)
		s.Require().Len(transactions, 2)

		expectedAmounts := map[string]int{
			sender1Username: sender1FirstTransfer + sender1SecondTransfer,
			sender2Username: sender2FirstTransfer,
		}

		for _, t := range transactions {
			s.Require().Equal(expectedAmounts[t.User], t.Amount)
		}
	})
}

func (s *PGTransferRepoTestSuite) TestFindAllForSenderGroupedByReceivers() {
	sender := uuid.New()
	receiver1 := uuid.New()
	receiver2 := uuid.New()
	receiver1Username := "receiver1"
	receiver2Username := "receiver2"

	s.insertEmployee(sender, "sender")
	s.insertEmployee(receiver1, receiver1Username)
	s.insertEmployee(receiver2, receiver2Username)

	receiver1FirstTransfer := 50
	receiver1SecondTransfer := 30
	receiver2FirstTransfer := 20

	testTransfers := []model.Transfer{
		{Id: uuid.New(), FromEmployee: sender, ToEmployee: receiver1, Amount: receiver1FirstTransfer},
		{Id: uuid.New(), FromEmployee: sender, ToEmployee: receiver2, Amount: receiver2FirstTransfer},
		{Id: uuid.New(), FromEmployee: sender, ToEmployee: receiver1, Amount: receiver1SecondTransfer},
	}

	for _, transfer := range testTransfers {
		_, err := s.pool.Exec(s.ctx,
			"insert into transfers (id, from_employee, to_employee, amount) VALUES ($1, $2, $3, $4)",
			transfer.Id, transfer.FromEmployee, transfer.ToEmployee, transfer.Amount)

		s.Require().NoError(err)
	}

	s.Run("should return grouped transactions for sender", func() {
		transactions, err := s.transferRepo.FindAllForSenderGroupedByReceivers(s.ctx, sender)
		s.Require().NoError(err)
		s.Require().Len(transactions, 2)

		expectedAmounts := map[string]int{
			receiver1Username: receiver1FirstTransfer + receiver1SecondTransfer,
			receiver2Username: receiver2FirstTransfer,
		}

		for _, t := range transactions {
			s.Require().Equal(expectedAmounts[t.User], t.Amount)
		}
	})
}

func (s *PGTransferRepoTestSuite) insertEmployee(employeeId uuid.UUID, username string) {
	_, err := s.pool.Exec(s.ctx,
		"insert into employees (id, username, password_hash, balance) VALUES ($1, $2, 'hash', 1000)",
		employeeId, username)
	s.Require().NoError(err)
}
