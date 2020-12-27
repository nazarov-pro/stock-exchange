package domain

import (
	"github.com/nazarov-pro/stock-exchange/pkg/container"
)

//SingleWalletResponse - returns a wallet
type SingleWalletResponse struct {
	*container.APIResponse
	item WalletResponse
}

//MultiWalletsResponse - returns multiple wallets
type MultiWalletsResponse struct {
	*container.APIResponse
	items *[]WalletResponse
}

// WalletResponse - wallet response
type WalletResponse struct {
	ID           uint64
	Balance      float64
	CurrencyCode string
	Status       WalletStatus
	Version      string
	CreatedDate  int64
}

//WalletTransactionResponse - wallet transaction response
type WalletTransactionResponse struct {
}

//WalletCreationRequest - wallet creation request
type WalletCreationRequest struct {
}

//WalletDebitRequest - wallet debit request
type WalletDebitRequest struct {
}

//WalletCreditRequest - wallet credit request
type WalletCreditRequest struct {
}
