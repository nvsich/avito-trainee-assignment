package model

import "github.com/google/uuid"

type Item struct {
	Id    uuid.UUID
	Name  string
	Price int
}
