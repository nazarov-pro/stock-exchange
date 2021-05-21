package domain

// WalletResponse - wallet response
type WalletResponse struct {
	ID           uint64       `json:"id"`
	Balance      int64        `json:"balance"`
	CurrencyCode string       `json:"currencyCode"`
	Status       WalletStatus `json:"status"`
	CreatedDate  int64        `json:"createdDate"`
}

//WalletTransactionResponse - wallet transaction response
type WalletTransactionResponse struct {
	ID             uint64                  `json:"id"`
	WalletID       uint64                   `json:"walletId"`
	ReferenceID    string                  `json:"referenceId"`
	Type           WalletTransactionType   `json:"type"`
	Status         WalletTransactionStatus `json:"status"`
	Amount         int64                  `json:"amount"`
	CurrencyCode   string                  `json:"currencyCode"`
	CreatedDate    int64                   `json:"createdDate"`
	LastUpdateDate int64                   `json:"lastUpdateDate"`
}

//WalletCreationRequest - wallet creation request
type WalletCreationRequest struct {
	AccountID    uint64 `json:"accountId"`
	CurrencyCode string `json:"currencyCode"`
}

//WalletDebitRequest - wallet debit request
type WalletDebitRequest struct {
	ReferenceID  string `json:"referenceId"`
	AccountID    uint64 `json:"accountId"`
	WalletID     uint64 `json:"walletId"`
	Amount       int64  `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}

//WalletCreditRequest - wallet credit request
type WalletCreditRequest struct {
	ReferenceID  string `json:"referenceId"`
	AccountID    uint64 `json:"accountId"`
	WalletID     uint64 `json:"walletId"`
	Amount       int64  `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}
