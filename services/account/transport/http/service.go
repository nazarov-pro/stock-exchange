package http

import (
	"context"
	"encoding/json"

	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nazarov-pro/stock-exchange/services/account/transport"
	"github.com/nazarov-pro/stock-exchange/pkg/container"
	"github.com/nazarov-pro/stock-exchange/services/account"
)

// NewService wires Go kit endpoints to the HTTP transport.
func NewService(
	svcEndpoints transport.Endpoints, options []kithttp.ServerOption, logger log.Logger,
) http.Handler {
	// set-up router and initialize http endpoints
	var (
		r            = mux.NewRouter()
		errorLogger  = kithttp.ServerErrorLogger(logger)
		errorEncoder = kithttp.ServerErrorEncoder(encodeErrorResponse)
	)
	options = append(options, errorLogger, errorEncoder)
	//options := []kithttp.ServerOption{
	//	kithttp.ServerErrorLogger(logger),
	//	kithttp.ServerErrorEncoder(encodeError),
	//}
	// HTTP Post - /orders
	r.Methods("POST").Path("/accounts").Handler(kithttp.NewServer(
		svcEndpoints.RegisterAccount,
		decodeRegisterRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/accounts/activate").Handler(kithttp.NewServer(
		svcEndpoints.ActivateAccount,
		decodeActivateAccountRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeActivateAccountRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req account.ActivateAccountRequest
	vals := r.URL.Query()
	req.Email = vals.Get("email")
	req.ActivationCode = vals.Get("activationCode")
	return req, nil
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req account.RegisterAccountRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeErrorResponse(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(
		container.APIResponse{
			Status:     "Operation was failed",
			Successful: false,
			Error: &container.APIError{
				Code:        "001",
				Description: err.Error(),
			},
		},
	)
}

func codeFrom(err error) int {
	switch err {
	// case order.ErrOrderNotFound:
	// 	return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}