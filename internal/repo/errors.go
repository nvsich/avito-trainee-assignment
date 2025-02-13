package repo

import "errors"

var (
	ErrEmployeeExists   = errors.New("employee already exists")
	ErrEmployeeNotFound = errors.New("employee not found")

	ErrItemNotFound = errors.New("item not found")

	ErrNoChange = errors.New("no change")
)
