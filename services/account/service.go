package account

import "context"

// Service - Registration, Resetting password and activating account
type Service interface {
	Register(ctx context.Context, req* RegisterAccountRequest) (*Account, error)

	Activate(ctx context.Context, req* ActivateAccountRequest) error
	// Reset password and activate/deactivate/suspend/delete 
}
