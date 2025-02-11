package repo

import "errors"

var (
	ErrEmployeeExists   = errors.New("employee already exists")
	ErrEmployeeNotFound = errors.New("employee not found")
)
