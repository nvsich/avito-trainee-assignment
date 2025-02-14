package model

import "github.com/google/uuid"

type Employee struct {
	Id           uuid.UUID
	Username     string
	Balance      int
	PasswordHash string
}
