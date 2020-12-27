package domain

import (
	"encoding/json"
	"testing"

	"github.com/nazarov-pro/stock-exchange/pkg/container"
)

func TestSerializeApiSingleItemToJson(t *testing.T) {
	walletResponse := &WalletResponse{ID: 1, Balance: 23.33, CurrencyCode: "AZN", Status: WsActivated, Version: "123", CreatedDate: 123}
	metadata := &container.APIResponse{Status: "OK", Successful: true}
	apiResponse := &container.APISingleItemResponse{Metadata: metadata, Item: walletResponse}
	bytes, err := json.Marshal(apiResponse)
	if err != nil {
		t.Fatalf("Error occurred when marshalling %v, err: %v", walletResponse, err)
	}
	jsonAsTxt := string(bytes)
	if jsonAsTxt == "" {
		t.Fatalf("Json as text should not be empty")
	}

}

func TestSerializeApiMultiItemsToJson(t *testing.T) {
	walletResponse := WalletResponse{ID: 1, Balance: 23.33, CurrencyCode: "AZN", Status: WsActivated, Version: "123", CreatedDate: 123}
	walletResponse2 := WalletResponse{ID: 2, Balance: 23.33, CurrencyCode: "AZN", Status: WsActivated, Version: "123", CreatedDate: 123}
	walletResponses := make([]WalletResponse, 2)

	walletResponses[0] = walletResponse
	walletResponses[1] = walletResponse2

	metadata := &container.APIResponse{Status: "OK", Successful: true}
	apiResponse := &container.APIMultiItemsResponse{Metadata: metadata, Items: walletResponses}
	bytes, err := json.Marshal(apiResponse)
	if err != nil {
		t.Fatalf("Error occurred when marshalling %v, err: %v", walletResponse, err)
	}
	jsonAsTxt := string(bytes)
	if jsonAsTxt == "" {
		t.Fatalf("Json as text should not be empty")
	}

}
