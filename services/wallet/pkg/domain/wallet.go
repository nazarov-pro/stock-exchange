package domain

import (
	"context"

	"github.com/nazarov-pro/stock-exchange/services/wallet/pkg/domain/api"
)

//WalletStatus - defines status of the wallet
type WalletStatus uint8

const (
	//WsUnknown - unknown status
	WsUnknown WalletStatus = iota
	//WsActivated - wallet is ready to use
	WsActivated
	//WsSuspended - temporary blocked (because of inactivity of the user)
	WsSuspended
	//WsBlocked - user blocked wallet manually
	WsBlocked
	//WsDeleted - user deleted account or wallet
	WsDeleted
)

// Currency currency ISO 4217
// https://www.iso.org/iso-4217-currency-codes.html
type Currency struct {
	Code     string // identifier - ex: AZN
	Name     string
	Number   uint8
	MinUnits uint8
}

//Wallet defines wallet of specified account
type Wallet struct {
	ID             uint64
	AccountID      uint64
	Balance        float64
	CurrencyCode   string
	Status         WalletStatus
	Version        string
	CreatedDate    int64
	LastUpdateDate int64
}

// WalletTransactionStatus - describes status of wallet transaction
type WalletTransactionStatus int8

const (
	// WtsUnknown - unknown status
	WtsUnknown WalletTransactionStatus = iota
	// WtsRegistered - transaction just gestired
	WtsRegistered
	// WtsBlocked - amount blocked
	WtsBlocked
	// WtsSuccessful - Operation was successfully completed
	WtsSuccessful
	// WtsReversingFailed - Reversing operation was failed
	WtsReversingFailed
	// WtsInSufficientFunds - insufficient funds
	WtsInSufficientFunds
	// WtsFailed - Operation was failed
	WtsFailed

	// WtsSystemError - system error occurred
	WtsSystemError
)

// WalletTransactionType - transaction type
type WalletTransactionType int8

const (
	//WttUnknown - unknown transaction
	WttUnknown WalletTransactionType = iota
	//WttDebit - debit transaction
	WttDebit
	//WttCredit - credit transaction
	WttCredit
)

//WalletTransaction defines wallet transaction
type WalletTransaction struct {
	ID             uint64
	WalletID       uint64
	ReferenceID    string // uuid or other identifier input
	Type           WalletTransactionType
	Status         WalletTransactionStatus
	Amount         float64
	CurrencyCode   string
	CreatedDate    int64
	LastUpdateDate int64
}

// WalletOperationType - transaction type
type WalletOperationType int8

const (
	//WohtUnknown - unknown transaction
	WohtUnknown WalletOperationType = iota
	//WohtCreation - wallet creation
	WohtCreation
	//WohtDebit - debit transaction
	WohtDebit
	//WohtCredit - credit transaction
	WohtCredit
	//WohtStatus - status change transaction
	WohtStatus
)

//WalletOperation stores all operation
type WalletOperation struct {
	ID            uint64
	WalletID      uint64
	TransactionID uint64 // transaction id if credit or debit request
	Type          WalletOperationType
	Values        map[string]string
	CreatedDate   int64
}

//WalletRepository - wallet reposittory
type WalletRepository interface {
	FindByID(ctx context.Context, ID uint64) (Wallet, error)
	FindByAccountID(ctx context.Context, accountID uint64) ([]Wallet, error)
	FindByIDandAccountID(ctx context.Context, ID uint64, accountID uint64) (Wallet, error)

	Save(ctx context.Context, wallet *Wallet) error
	UpdateBalance(ctx context.Context, wallet *Wallet, newVersion string, newBalance int64) error
	UpdateStatus(ctx context.Context, wallet *Wallet, newVersion string, newStatus WalletStatus) error
}

//WalletTransactionRepository - wallet transaction reposittory
type WalletTransactionRepository interface {
	FindByID(ctx context.Context, ID uint64) (WalletTransaction, error)
	FindByWalletID(ctx context.Context, walletID uint64) ([]WalletTransaction, error)
	FindByAccountID(ctx context.Context, accountID uint64) ([]WalletTransaction, error)
	FindByIDandAccountID(ctx context.Context, ID uint64, AccountID uint64) (WalletTransaction, error)

	Save(ctx context.Context, walletTransaction *WalletTransaction) error
	UpdateStatus(ctx context.Context, ID uint64, newStatus WalletTransactionStatus) error
}

//WalletOperationRepository - wallet operation element
type WalletOperationRepository interface {
	Save(ctx context.Context, walletOperation *WalletOperation) error
}

// WalletService - all wallet related stuff in here
type WalletService interface {
	GetWalletsByAccountId(ctx context.Context, accountID uint64) ([]api.WalletResponse, error)
	GetWalletTransactionsByAccountID(ctx context.Context, accountID uint64) ([]api.WalletTransactionResponse, error)
	GetWalletByIdAndAccountId(ctx context.Context, walletID uint64, accountID uint64) (api.WalletResponse, error)

	Create(ctx context.Context, req api.WalletCreationRequest) (api.WalletResponse, error)
	Credit(ctx context.Context, req api.WalletCreditRequest) (api.WalletResponse, error)
	Debit(ctx context.Context, req api.WalletDebitRequest) (api.WalletResponse, error)
}
