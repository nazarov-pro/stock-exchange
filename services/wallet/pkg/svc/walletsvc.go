package svc

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/nazarov-pro/stock-exchange/services/wallet/pkg/domain"
)

type walletSvc struct {
	walletRepo          domain.WalletRepository
	walletTransRepo     domain.WalletTransactionRepository
	walletOperationRepo domain.WalletOperationRepository
	logger              log.Logger
}

//NewWalletRepo a new instance of wallet repo
func NewWalletRepo(walletRepo domain.WalletRepository, walletTransRepo domain.WalletTransactionRepository, walletOperationRepo domain.WalletOperationRepository, logger log.Logger) domain.WalletService {
	return &walletSvc{walletRepo: walletRepo, walletTransRepo: walletTransRepo, walletOperationRepo: walletOperationRepo, logger: logger}
}

func (svc *walletSvc) GetWalletsByAccountID(ctx context.Context, accountID uint64) (*[]domain.WalletResponse, error) {
	wallets, err := svc.walletRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	wlts := *wallets
	responses := make([]domain.WalletResponse, len(wlts), cap(wlts))
	for i, wallet := range (wlts) {
		responses[i] = domain.WalletResponse{ID: wallet.ID, Balance: wallet.Balance, CurrencyCode: wallet.CurrencyCode, Status: wallet.Status, Version: wallet.Version, CreatedDate: wallet.CreatedDate}
	}
	return &responses, nil
}

func (svc *walletSvc) GetWalletTransactionsByAccountID(ctx context.Context, accountID uint64) (*[]domain.WalletTransactionResponse, error) {

	return nil, nil
}

func (svc *walletSvc) GetWalletByIDAndAccountID(ctx context.Context, walletID uint64, accountID uint64) (*domain.WalletResponse, error) {
	return nil, nil
}

func (svc *walletSvc) Create(ctx context.Context, req *domain.WalletCreationRequest) (*domain.WalletResponse, error) {
	return nil, nil
}

func (svc *walletSvc) Credit(ctx context.Context, req *domain.WalletCreditRequest) (*domain.WalletResponse, error) {
	return nil, nil
}

func (svc *walletSvc) Debit(ctx context.Context, req *domain.WalletDebitRequest) (*domain.WalletResponse, error) {
	return nil, nil
}
