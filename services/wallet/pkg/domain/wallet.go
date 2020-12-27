package domain

import (
	"context"
	"errors"
)

var (
	//ErrWalletNotFount - wallet not found
	ErrWalletNotFount = errors.New("Not any wallet found")
	//ErrWalletTransNotFount - wallet not found
	ErrWalletTransNotFount = errors.New("Not any wallet transactions found")
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
	Version        string
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
	FindByID(ctx context.Context, ID uint64) (*Wallet, error)
	FindByAccountID(ctx context.Context, accountID uint64) (*[]Wallet, error)
	FindByIDAndAccountID(ctx context.Context, ID uint64, accountID uint64) (*Wallet, error)

	Save(ctx context.Context, wallet *Wallet) error
	UpdateBalance(ctx context.Context, wallet *Wallet, newVersion string, newBalance float64, updateDate int64) error
	UpdateStatus(ctx context.Context, wallet *Wallet, newVersion string, newStatus WalletStatus, updateDate int64) error
}

//WalletTransactionRepository - wallet transaction reposittory
type WalletTransactionRepository interface {
	FindByID(ctx context.Context, ID uint64) (*WalletTransaction, error)
	FindByWalletID(ctx context.Context, walletID uint64) (*[]WalletTransaction, error)
	FindByAccountID(ctx context.Context, accountID uint64) (*[]WalletTransaction, error)
	FindByIDAndAccountID(ctx context.Context, ID uint64, walletID uint64) (*WalletTransaction, error)

	Save(ctx context.Context, walletTransaction *WalletTransaction) error
	UpdateStatus(ctx context.Context, walletTransaction *WalletTransaction, newStatus WalletTransactionStatus, newVersion string, updateDate int64) error
}

//WalletOperationRepository - wallet operation element
type WalletOperationRepository interface {
	Save(ctx context.Context, walletOperation *WalletOperation) error
}

// WalletService - all wallet related stuff in here
type WalletService interface {
	GetWalletsByAccountID(ctx context.Context, accountID uint64) (*[]WalletResponse, error)
	GetWalletTransactionsByAccountID(ctx context.Context, accountID uint64) (*[]WalletTransactionResponse, error)
	GetWalletByIDAndAccountID(ctx context.Context, walletID uint64, accountID uint64) (*WalletResponse, error)

	Create(ctx context.Context, req *WalletCreationRequest) (*WalletResponse, error)
	Credit(ctx context.Context, req *WalletCreditRequest) (*WalletResponse, error)
	Debit(ctx context.Context, req *WalletDebitRequest) (*WalletResponse, error)
}
