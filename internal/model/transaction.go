package model

import "github.com/google/uuid"

type Transaction struct {
	Id           uuid.UUID
	FromEmployee uuid.UUID
	ToEmployee   uuid.UUID
	Amount       int
}
