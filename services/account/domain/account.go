package domain

import (
	"context"
	"errors"
)

var (
	// ErrExistedUsername - defines username has already been taken
	ErrExistedUsername = errors.New("username has already been taken")
	// ErrExistedEmail - defines email has already been taken
	ErrExistedEmail = errors.New("email has already been taken")
	// ErrAccountNotFound couldnt find any relevant account
	ErrAccountNotFound = errors.New("account is not found")
)

// Status represent account's current status
type Status uint8

const (
	// Registered account just registered
	Registered Status = iota
	// Activated account verified
	Activated
	// Deactivated account deactivated (reached max. of inactivity), verification pending
	Deactivated
	// Suspended account suspended temporarily
	Suspended
	// Deleted account deleted permanently
	Deleted
)

// Account holds account based functions such as registration, resetting password, activating account
type Account struct {
	ID             int64
	Username       string
	Email          string
	Password       string
	Status         Status
	ActivationCode string
	CreatedDate    int64
	LastUpdateDate int64
}

// Repository account crud operations
type Repository interface {
	//GetById(ctx context.Context, id int64) (Account, error)

	FindByUsernameOrEmail(ctx context.Context, username string, email string) (*Account, error)

	Create(ctx context.Context, acc *Account) error

	UpdateStatus(ctx context.Context, email string, activationCode string, fromStatus Status, toStatus Status) error
}

// Service - Registration, Resetting password and activating account
type Service interface {
	Register(ctx context.Context, req* RegisterAccountRequest) (*Account, error)

	Activate(ctx context.Context, req* ActivateAccountRequest) error
	// Reset password and activate/deactivate/suspend/delete 
}
