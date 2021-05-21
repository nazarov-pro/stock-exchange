package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/nazarov-pro/stock-exchange/pkg/container"
	"github.com/nazarov-pro/stock-exchange/services/wallet/pkg/domain"
)

// Endpoints holds all Go kit endpoints for the Account service.
type Endpoints struct {
	GetWalletsOfAccount endpoint.Endpoint
	CreateWallet        endpoint.Endpoint
	CreditWallet        endpoint.Endpoint
	DebitWallet         endpoint.Endpoint
}

// MakeEndpoints initializes all go kit endpoints for account service
func MakeEndpoints(svc *domain.WalletService) Endpoints {
	return Endpoints{
		GetWalletsOfAccount: makeGetWalletsOfAccountEndpoint(svc),
		CreateWallet:        makeCreateWalletEndpoint(svc),
		CreditWallet:        makeCreditWalletEndpoint(svc),
		DebitWallet:         makeDebitWalletEndpoint(svc),
	}
}

func makeGetWalletsOfAccountEndpoint(svc *domain.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		response, err := (*svc).GetWalletsByAccountID(ctx, 1)
		if err == nil {
			return container.APIMultiItemsResponse{
				Metadata: &container.APIResponse{
					Status:     "success",
					Successful: true,
				},
				Items: response,
			}, nil
		}

		return container.APIMultiItemsResponse{
			Metadata: &container.APIResponse{
				Status:     "success",
				Successful: true,
			},
		}, nil
	}
}

func makeCreateWalletEndpoint(svc *domain.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		response, err := (*svc).Create(ctx, request.(*domain.WalletCreationRequest))
		if err == nil {
			return container.APISingleItemResponse{
				Metadata: &container.APIResponse{
					Status:     "success",
					Successful: true,
				},
				Item: response,
			}, nil
		}

		return container.APISingleItemResponse{
			Metadata: &container.APIResponse{
				Status:     "success",
				Successful: true,
			},
		}, nil
	}
}

func makeCreditWalletEndpoint(svc *domain.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		response, err := (*svc).Credit(ctx, request.(*domain.WalletCreditRequest))
		if err == nil {
			return container.APISingleItemResponse{
				Metadata: &container.APIResponse{
					Status:     "success",
					Successful: true,
				},
				Item: response,
			}, nil
		}

		return container.APISingleItemResponse{
			Metadata: &container.APIResponse{
				Status:     "success",
				Successful: true,
			},
		}, nil
	}
}

func makeDebitWalletEndpoint(svc *domain.WalletService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		response, err := (*svc).Debit(ctx, request.(*domain.WalletDebitRequest))
		if err == nil {
			return container.APISingleItemResponse{
				Metadata: &container.APIResponse{
					Status:     "success",
					Successful: true,
				},
				Item: response,
			}, nil
		}

		return container.APISingleItemResponse{
			Metadata: &container.APIResponse{
				Status:     "success",
				Successful: true,
			},
		}, nil
	}
}
