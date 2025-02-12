package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrNotEnoughCoins         = errors.New("not enough coins")
	ErrNegativeTransferAmount = errors.New("negative transfer amount")
	ErrReceiverNotFound       = errors.New("receiver not found")
)
