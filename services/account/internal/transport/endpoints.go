package transport

import (
	"context"
	"unicode"

	"github.com/go-kit/kit/endpoint"
	validate "github.com/go-playground/validator/v10"
	"github.com/nazarov-pro/stock-exchange/pkg/container"
	"github.com/nazarov-pro/stock-exchange/services/account/domain"
)

var (
	validator *validate.Validate
)

func init() {
	validator = validate.New()
	validator.RegisterValidation("username", func(fl validate.FieldLevel) bool {
		size := 0

		for i, c := range fl.Field().String() {
			if i == 0 {
				// the first character must be a letter
				if !unicode.IsLetter(c) {
					return false
				}
			} else {
				if unicode.IsSymbol(c) && c != '_' && c != '-' && c != '.' {
					return false
				}
			}
			size++
		}
		return size >= 6 && size < 25
	})

	validator.RegisterValidation("passwd", func(fl validate.FieldLevel) bool {
		size := 0
		specialCharaters := 0
		digitCharaters := 0
		uppercaseLetter := 0
		lowercaseLetter := 0

		for _, c := range fl.Field().String() {

			if unicode.IsUpper(c) {
				uppercaseLetter++
			} else if unicode.IsLower(c) {
				lowercaseLetter++
			} else if unicode.IsDigit(c) {
				digitCharaters++
			} else {
				specialCharaters++
			}

			size++
		}

		if specialCharaters == 0 || digitCharaters == 0 ||
			uppercaseLetter == 0 || lowercaseLetter == 0 {
			return false
		}
		return size >= 8 && size < 25
	})

}

// Endpoints holds all Go kit endpoints for the Account service.
type Endpoints struct {
	RegisterAccount endpoint.Endpoint
	ActivateAccount endpoint.Endpoint
}

// MakeEndpoints initializes all go kit endpoints for account service
func MakeEndpoints(svc *domain.Service) Endpoints {
	return Endpoints{
		RegisterAccount: makeRegisterAccountEndpoint(svc),
		ActivateAccount: makeActivateAccountEndpoint(svc),
	}
}

func makeActivateAccountEndpoint(svc *domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(domain.ActivateAccountRequest)
		err := validator.Struct(req)
		if err != nil {
			for _, e := range err.(validate.ValidationErrors) {
				return nil, e
			}
		}
		err = (*svc).Activate(ctx, &req)
		if err == nil {
			return domain.ActivateAccountResponse{
				APIResponse: &container.APIResponse{
					Status:     "Your account activated successfully.",
					Successful: true,
				},
			}, nil
		}

		return domain.ActivateAccountResponse{
			APIResponse: &container.APIResponse{
				Status:     "Operation was failed",
				Successful: false,
				Error: &container.APIError{
					Code:        "001",
					Description: err.Error(),
				},
			},
		}, nil
	}
}

func makeRegisterAccountEndpoint(svc *domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(domain.RegisterAccountRequest)
		err := validator.Struct(req)
		if err != nil {
			for _, e := range err.(validate.ValidationErrors) {
				return nil, e
			}
		}

		_, err = (*svc).Register(ctx, &req)
		if err == nil {
			return domain.RegisterAccountResponse{
				APIResponse: &container.APIResponse{
					Status:     "Activation link will be send the email.",
					Successful: true,
				},
			}, nil
		}
		return domain.RegisterAccountResponse{
			APIResponse: &container.APIResponse{
				Status:     "Operation was failed",
				Successful: false,
				Error: &container.APIError{
					Code:        "001",
					Description: err.Error(),
				},
			},
		}, nil
	}
}
