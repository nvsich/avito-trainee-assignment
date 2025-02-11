package model

import "github.com/google/uuid"

type Employee struct {
	Id           uuid.UUID
	Login        string
	Balance      int
	PasswordHash string
}
