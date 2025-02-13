package model

import "github.com/google/uuid"

type EmployeeInventory struct {
	Id         uuid.UUID
	EmployeeId uuid.UUID
	ItemId     uuid.UUID
	Amount     int
}
