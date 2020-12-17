package domain

import (
	"github.com/nazarov-pro/stock-exchange/pkg/container"
)

// RegisterAccountRequest request for the registering account
type RegisterAccountRequest struct {
	Username string `json:"username" validate:"required,username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwd"`
}

// RegisterAccountResponse response after registering account
type RegisterAccountResponse struct {
	*container.APIResponse
}

// ActivateAccountRequest account activation request (not json from query params)
type ActivateAccountRequest struct {
	Email          string `validate:"required,email"`
	ActivationCode string `validate:"required"`
}

// ActivateAccountResponse consist of response metadata
type ActivateAccountResponse struct {
	*container.APIResponse
}
