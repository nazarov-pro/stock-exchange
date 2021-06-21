package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"github.com/nazarov-pro/stock-exchange/services/auth"

	"github.com/go-kit/kit/endpoint"
)

func MakeAuthEndpoint(svc auth.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthRequest)
		token, err := svc.Auth(req.ClientID, req.ClientSecret)
		if err != nil {
			return nil, err
		}
		return AuthResponse{token, ""}, nil
	}
}

func DecodeAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}


type AuthRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type AuthResponse struct {
	Token string `json:"token,omitempty"`
	Err   string `json:"err,omitempty"`
}

func AuthErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusUnauthorized
	msg := err.Error()

	w.WriteHeader(code)
	err = json.NewEncoder(w).Encode(AuthResponse{Token: "", Err: msg})
	if err != nil {
		panic(err)
	}
}