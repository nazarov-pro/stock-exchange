package svc

import (
	"strconv"
	"context"

	"github.com/go-kit/kit/log"
	"github.com/nazarov-pro/stock-exchange/pkg/util/gen"
	"github.com/nazarov-pro/stock-exchange/pkg/util/time"
	"github.com/nazarov-pro/stock-exchange/services/wallet/pkg/conf"
	"github.com/nazarov-pro/stock-exchange/services/wallet/pkg/domain"
)

type walletSvc struct {
	walletRepo          domain.WalletRepository
	walletTransRepo     domain.WalletTransactionRepository
	walletOperationRepo domain.WalletOperationRepository
	transManager        conf.TransactionManager
	logger              log.Logger
}

//NewWalletSvc a new instance of wallet service
func NewWalletSvc(
	walletRepo domain.WalletRepository,
	walletTransRepo domain.WalletTransactionRepository,
	walletOperationRepo domain.WalletOperationRepository,
	transManager conf.TransactionManager,
	logger log.Logger,
) domain.WalletService {
	return &walletSvc{
		walletRepo:          walletRepo,
		walletTransRepo:     walletTransRepo,
		walletOperationRepo: walletOperationRepo,
		transManager:        transManager,
		logger:              logger,
	}
}

func (svc *walletSvc) GetWalletsByAccountID(ctx context.Context, accountID uint64) (*[]domain.WalletResponse, error) {
	wallets, err := svc.walletRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		svc.logger.Log("err", err)
		return nil, err
	}
	responses := make([]domain.WalletResponse, len(*wallets), cap(*wallets))
	for i, wallet := range *wallets {
		responses[i] = domain.WalletResponse{ID: wallet.ID, Balance: wallet.Balance, CurrencyCode: wallet.CurrencyCode, Status: wallet.Status, CreatedDate: wallet.CreatedDate}
	}
	return &responses, nil
}

func (svc *walletSvc) GetWalletTransactionsByAccountID(ctx context.Context, accountID uint64) (*[]domain.WalletTransactionResponse, error) {
	transactions, err := svc.walletTransRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		svc.logger.Log("err", err)
		return nil, err
	}
	responses := make([]domain.WalletTransactionResponse, len(*transactions), cap(*transactions))
	for i, transaction := range *transactions {
		responses[i] = domain.WalletTransactionResponse{ID: transaction.ID, WalletID: transaction.WalletID,
			ReferenceID: transaction.ReferenceID, Type: transaction.Type, Status: transaction.Status,
			Amount: transaction.Amount, CurrencyCode: transaction.CurrencyCode,
			CreatedDate: transaction.CreatedDate, LastUpdateDate: transaction.LastUpdateDate}
	}
	return &responses, nil
}

func (svc *walletSvc) GetWalletByIDAndAccountID(ctx context.Context, walletID uint64, accountID uint64) (*domain.WalletResponse, error) {
	wallet, err := svc.walletRepo.FindByIDAndAccountID(ctx, walletID, accountID)
	if err != nil {
		return nil, err
	}
	response := domain.WalletResponse{ID: wallet.ID, Balance: wallet.Balance, CurrencyCode: wallet.CurrencyCode, Status: wallet.Status, CreatedDate: wallet.CreatedDate}
	return &response, nil
}

func (svc *walletSvc) Create(ctx context.Context, req *domain.WalletCreationRequest) (*domain.WalletResponse, error) {
	tx, err := svc.transManager.Begin()
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, "tx", tx)
	wallet := &domain.Wallet{
		AccountID:    req.AccountID,
		CurrencyCode: req.CurrencyCode,
		Balance:      0.0,
		Version:      gen.NewUUID(),
		CreatedDate:  time.Epoch(),
		Status:       domain.WsActivated,
	}

	err = svc.walletRepo.Save(ctx, wallet)
	if err != nil {
		svc.logger.Log("err", err)
		svc.logger.Log("msg", "transaction rollback", "err", tx.Rollback())
		return nil, err
	}

	walletOperation := domain.WalletOperation{Type: domain.WotCreation, CreatedDate: time.Epoch(), WalletID: wallet.ID}

	err = svc.walletOperationRepo.Save(ctx, &walletOperation)
	if err != nil {
		svc.logger.Log("err", err)
		svc.logger.Log("msg", "transaction rollback", "err", tx.Rollback())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		svc.logger.Log("err", err)
		return nil, err
	}

	response := domain.WalletResponse{ID: wallet.ID, Balance: wallet.Balance, CurrencyCode: wallet.CurrencyCode, Status: wallet.Status, CreatedDate: wallet.CreatedDate}
	return &response, nil
}

func (svc *walletSvc) Credit(ctx context.Context, req *domain.WalletCreditRequest) (*domain.WalletResponse, error) {
	tx, err := svc.transManager.Begin()
	if err != nil {
		return nil, err
	}

	wallet, err := svc.walletRepo.FindByIDAndAccountID(ctx, req.WalletID, req.AccountID)
	if err != nil {
		svc.logger.Log("err", err)
		return nil, err
	}

	if wallet.CurrencyCode != req.CurrencyCode {
		return nil, nil
	}

	if req.Amount <= 0.0 {
		return nil, nil
	}
	oldVersion := wallet.Version
	oldBalance := wallet.Balance

	newBalance := wallet.Balance + req.Amount
	newVersion := gen.NewUUID()
	updatedAt := time.Epoch()

	ctx = context.WithValue(ctx, "tx", tx)
	err = svc.walletRepo.UpdateBalance(ctx, wallet, newVersion, newBalance, updatedAt)
	if err != nil {
		svc.logger.Log("err", err)
		svc.logger.Log("msg", "transaction rollback", "err", tx.Rollback())
		return nil, err
	}

	transaction := &domain.WalletTransaction{
		WalletID:     wallet.ID,
		ReferenceID:  req.ReferenceID,
		Type:         domain.WttCredit,
		Status:       domain.WtsSuccessful,
		Amount:       req.Amount,
		CurrencyCode: req.CurrencyCode,
		Version:      newVersion,
		CreatedDate:  updatedAt,
	}

	err = svc.walletTransRepo.Save(ctx, transaction)
	if err != nil {
		svc.logger.Log("err", err)
		svc.logger.Log("msg", "transaction rollback", "err", tx.Rollback())
		return nil, err
	}

	values := make(map[string]string, 4)
	values["new_balance"] = strconv.FormatInt(newBalance, 10)
	values["old_balance"] = strconv.FormatInt(oldBalance, 10)
	values["new_version"] = newVersion
	values["old_version"] = oldVersion
	walletOperation := domain.WalletOperation{
		WalletID:      wallet.ID,
		TransactionID: transaction.ID,
		Type:          domain.WotCredit,
		CreatedDate:   updatedAt,
		Values:        values,
	}

	err = svc.walletOperationRepo.Save(ctx, &walletOperation)
	if err != nil {
		svc.logger.Log("err", err)
		svc.logger.Log("msg", "transaction rollback", "err", tx.Rollback())
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		svc.logger.Log("err", err)
		return nil, err
	}

	response := domain.WalletResponse{ID: wallet.ID, Balance: wallet.Balance, CurrencyCode: wallet.CurrencyCode, Status: wallet.Status, CreatedDate: wallet.CreatedDate}
	return &response, nil
}

func (svc *walletSvc) Debit(ctx context.Context, req *domain.WalletDebitRequest) (*domain.WalletResponse, error) {
	tx, err := svc.transManager.Begin()
	if err != nil {
		return nil, err
	}

	wallet, err := svc.walletRepo.FindByIDAndAccountID(ctx, req.WalletID, req.AccountID)
	if err != nil {
		svc.logger.Log("err", err)
		return nil, err
	}

	if wallet.CurrencyCode != req.CurrencyCode {
		return nil, nil
	}

	if wallet.Balance < req.Amount {
		return nil, nil
	}
	oldVersion := wallet.Version
	oldBalance := wallet.Balance

	newBalance := wallet.Balance - req.Amount
	newVersion := gen.NewUUID()
	updatedAt := time.Epoch()

	ctx = context.WithValue(ctx, "tx", tx)
	err = svc.walletRepo.UpdateBalance(ctx, wallet, newVersion, newBalance, updatedAt)
	if err != nil {
		svc.logger.Log("err", err)
		svc.logger.Log("msg", "transaction rollback", "err", tx.Rollback())
		return nil, err
	}

	transaction := &domain.WalletTransaction{
		WalletID:     wallet.ID,
		ReferenceID:  req.ReferenceID,
		Type:         domain.WttDebit,
		Status:       domain.WtsSuccessful,
		Amount:       req.Amount,
		CurrencyCode: req.CurrencyCode,
		Version:      newVersion,
		CreatedDate:  updatedAt,
	}

	err = svc.walletTransRepo.Save(ctx, transaction)
	if err != nil {
		svc.logger.Log("err", err)
		svc.logger.Log("msg", "transaction rollback", "err", tx.Rollback())
		return nil, err
	}

	values := make(map[string]string, 4)
	values["new_balance"] = strconv.FormatInt(newBalance, 10)
	values["old_balance"] = strconv.FormatInt(oldBalance, 10)
	values["new_version"] = newVersion
	values["old_version"] = oldVersion
	walletOperation := domain.WalletOperation{
		WalletID:      wallet.ID,
		TransactionID: transaction.ID,
		Type:          domain.WotDebit,
		CreatedDate:   updatedAt,
		Values:        values,
	}

	err = svc.walletOperationRepo.Save(ctx, &walletOperation)
	if err != nil {
		svc.logger.Log("err", err)
		svc.logger.Log("msg", "transaction rollback", "err", tx.Rollback())
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		svc.logger.Log("err", err)
		return nil, err
	}

	response := domain.WalletResponse{ID: wallet.ID, Balance: wallet.Balance, CurrencyCode: wallet.CurrencyCode, Status: wallet.Status, CreatedDate: wallet.CreatedDate}
	return &response, nil
}
