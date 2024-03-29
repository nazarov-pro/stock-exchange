package http

import (
	"context"
	"encoding/json"

	"net/http"

	// "net/http/pprof"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nazarov-pro/stock-exchange/pkg/container"
	"github.com/nazarov-pro/stock-exchange/services/wallet/pkg/domain"
	"github.com/nazarov-pro/stock-exchange/services/wallet/pkg/transport"
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

	r.Methods("GET").Path("/wallets").Handler(kithttp.NewServer(
		svcEndpoints.GetWalletsOfAccount,
		decodeActivateAccountRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/wallets").Handler(kithttp.NewServer(
		svcEndpoints.CreateWallet,
		decodeCreateWallet,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/wallets/credit").Handler(kithttp.NewServer(
		svcEndpoints.CreditWallet,
		decodeCreditWallet,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/wallets/debit").Handler(kithttp.NewServer(
		svcEndpoints.DebitWallet,
		decodeDebitWallet,
		encodeResponse,
		options...,
	))

	// Add the pprof routes
	// r.HandleFunc("/debug/pprof/", pprof.Index)
	// r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	// r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// r.Handle("/debug/pprof/block", pprof.Handler("block"))
	// r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	// r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	// r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	return r
}

func decodeActivateAccountRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodeCreditWallet(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req domain.WalletCreditRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return &req, nil
}

func decodeDebitWallet(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req domain.WalletDebitRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return &req, nil
}

func decodeCreateWallet(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req domain.WalletCreationRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return &req, nil
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
