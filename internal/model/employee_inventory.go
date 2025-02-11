package model

import "github.com/google/uuid"

type EmployeeInventory struct {
	Id         uuid.UUID
	EmployeeId uuid.UUID
	ProductId  uuid.UUID
	Amount     int
}
