package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrNotEnoughCoins         = errors.New("not enough coins")
	ErrNegativeTransferAmount = errors.New("negative transfer amount")
	ErrReceiverNotFound       = errors.New("receiver not found")
	ErrSenderNotFound         = errors.New("sender not found")
	ErrTransferToSameEmployee = errors.New("transfer to same employee")

	ErrEmployeeNotFound = errors.New("employee not found")
	ErrItemNotFound     = errors.New("item not found")
)
