package repo

import "errors"

var (
	ErrEmployeeExists   = errors.New("employee already exists")
	ErrEmployeeNotFound = errors.New("employee not found")

	ErrItemNotFound              = errors.New("item not found")
	ErrEmployeeInventoryNotFound = errors.New("employee inventory not found")

	ErrInventoryNotFound = errors.New("inventory not found")
)
