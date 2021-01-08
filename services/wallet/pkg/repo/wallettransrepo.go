package repo

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/nazarov-pro/stock-exchange/services/wallet/pkg/domain"
)

type walletTransactionRepo struct {
	db     *sql.DB
	logger log.Logger
}

//NewWalletTransactionRepo a new instance of wallet transaction repo
func NewWalletTransactionRepo(db *sql.DB, logger log.Logger) domain.WalletTransactionRepository {
	return &walletTransactionRepo{db: db, logger: logger}
}

func (repo *walletTransactionRepo) FindByID(ctx context.Context, ID uint64) (*domain.WalletTransaction, error) {
	sqlQuery := `SELECT "WT"."ID", "WT"."WALLET_ID", "WT"."REFERENCE_ID", "WT"."TYPE", "WT"."STATUS", "WT"."VERSION", "WT"."AMOUNT", "WT"."CURRENCY_CODE", "WT"."CREATED_DATE", COALESCE("WT"."LAST_UPDATE_DATE", 0) FROM "WALLET"."WALLET_TRANSACTIONS" "WT" WHERE "WT"."ID"=$1`
	sqlRow := repo.db.QueryRowContext(ctx, sqlQuery, ID)
	var walletTrans domain.WalletTransaction
	err := sqlRow.Scan(&walletTrans.ID, &walletTrans.WalletID, &walletTrans.ReferenceID, &walletTrans.Type, &walletTrans.Status, &walletTrans.Version, &walletTrans.Amount, &walletTrans.CurrencyCode, &walletTrans.CreatedDate, &walletTrans.LastUpdateDate)
	switch err {
	case sql.ErrNoRows:
		return nil, domain.ErrWalletNotFount
	case nil:
		return &walletTrans, nil
	default:
		return nil, err
	}
}

func (repo *walletTransactionRepo) FindByWalletID(ctx context.Context, walletID uint64) (*[]domain.WalletTransaction, error) {
	sqlQuery := `SELECT "WT"."ID", "WT"."WALLET_ID", "WT"."REFERENCE_ID", "WT"."TYPE", "WT"."STATUS", "WT"."VERSION", "WT"."AMOUNT", "WT"."CURRENCY_CODE", "WT"."CREATED_DATE", COALESCE("WT"."LAST_UPDATE_DATE", 0) FROM "WALLET"."WALLET_TRANSACTIONS" "WT" WHERE "WT"."WALLET_ID"=$1`
	sqlRows, err := repo.db.QueryContext(ctx, sqlQuery, walletID)
	if err != nil {
		return nil, err
	}
	defer sqlRows.Close()
	walletTransactions := make([]domain.WalletTransaction, 0)
	for sqlRows.Next() {
		var walletTrans domain.WalletTransaction
		err := sqlRows.Scan(&walletTrans.ID, &walletTrans.WalletID, &walletTrans.ReferenceID, &walletTrans.Type, &walletTrans.Status, &walletTrans.Version, &walletTrans.Amount, &walletTrans.CurrencyCode, &walletTrans.CreatedDate, &walletTrans.LastUpdateDate)
		if err != nil {
			return nil, err
		}
		walletTransactions = append(walletTransactions, walletTrans)
	}
	err = sqlRows.Err()

	switch err {
	case sql.ErrNoRows:
		return nil, domain.ErrWalletNotFount
	case nil:
		return &walletTransactions, nil
	default:
		return nil, err
	}
}

func (repo *walletTransactionRepo) FindByAccountID(ctx context.Context, accountID uint64) (*[]domain.WalletTransaction, error) {
	sqlQuery := `SELECT "WT"."ID", "WT"."WALLET_ID", "WT"."REFERENCE_ID", "WT"."TYPE", "WT"."STATUS", "WT"."VERSION", "WT"."AMOUNT", "WT"."CURRENCY_CODE", "WT"."CREATED_DATE", COALESCE("WT"."LAST_UPDATE_DATE", 0) FROM "WALLET"."WALLET_TRANSACTIONS" "WT" JOIN "WALLET"."WALLETS" "W" ON "WT"."WALLET_ID"="W"."ID" WHERE "W"."ACCOUNT_ID"=$1`
	sqlRows, err := repo.db.QueryContext(ctx, sqlQuery, accountID)
	if err != nil {
		return nil, err
	}
	defer sqlRows.Close()
	walletTransactions := make([]domain.WalletTransaction, 0)
	for sqlRows.Next() {
		var walletTrans domain.WalletTransaction
		err := sqlRows.Scan(&walletTrans.ID, &walletTrans.WalletID, &walletTrans.ReferenceID, &walletTrans.Type, &walletTrans.Status, &walletTrans.Version, &walletTrans.Amount, &walletTrans.CurrencyCode, &walletTrans.CreatedDate, &walletTrans.LastUpdateDate)
		if err != nil {
			return nil, err
		}
		walletTransactions = append(walletTransactions, walletTrans)
	}
	err = sqlRows.Err()

	switch err {
	case sql.ErrNoRows:
		return nil, domain.ErrWalletNotFount
	case nil:
		return &walletTransactions, nil
	default:
		return nil, err
	}
}

func (repo *walletTransactionRepo) FindByIDAndAccountID(ctx context.Context, ID uint64, accountID uint64) (*domain.WalletTransaction, error) {
	sqlQuery := `SELECT "WT"."ID", "WT"."WALLET_ID", "WT"."REFERENCE_ID", "WT"."TYPE", "WT"."STATUS", "WT"."VERSION", "WT"."AMOUNT", "WT"."CURRENCY_CODE", "WT"."CREATED_DATE", COALESCE("WT"."LAST_UPDATE_DATE", 0) FROM "WALLET"."WALLET_TRANSACTIONS" "WT" JOIN "WALLET"."WALLETS" "W" ON "WT"."WALLET_ID"="W"."ID" WHERE "WT"."ID"=$1 AND "W"."ACCOUNT_ID"=$2`
	sqlRow := repo.db.QueryRowContext(ctx, sqlQuery, ID, accountID)
	var walletTrans domain.WalletTransaction
	err := sqlRow.Scan(&walletTrans.ID, &walletTrans.WalletID, &walletTrans.ReferenceID, &walletTrans.Type, &walletTrans.Status, &walletTrans.Version, &walletTrans.Amount, &walletTrans.CurrencyCode, &walletTrans.CreatedDate, &walletTrans.LastUpdateDate)
	switch err {
	case sql.ErrNoRows:
		return nil, domain.ErrWalletNotFount
	case nil:
		return &walletTrans, nil
	default:
		return nil, err
	}
}

func (repo *walletTransactionRepo) Save(ctx context.Context, walletTransaction *domain.WalletTransaction) error {
	tx := ctx.Value("tx").(*sql.Tx)
	id, err := repo.generateID(ctx)
	if err != nil {
		return err
	}
	walletTransaction.ID = id

	sqlQuery := `INSERT INTO "WALLET"."WALLET_TRANSACTIONS"("ID", "WALLET_ID", "REFERENCE_ID", "TYPE", "STATUS", "VERSION", "AMOUNT", "CURRENCY_CODE", "CREATED_DATE") values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	result, err := tx.ExecContext(ctx, sqlQuery, walletTransaction.ID, walletTransaction.WalletID, walletTransaction.ReferenceID, walletTransaction.Type, walletTransaction.Status, walletTransaction.Version, walletTransaction.Amount, walletTransaction.CurrencyCode, walletTransaction.CreatedDate)
	if err != nil {
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows > 0 {
		return nil
	}
	return domain.ErrNoAffectedRows
}

func (repo *walletTransactionRepo) UpdateStatus(ctx context.Context, walletTransaction *domain.WalletTransaction, newStatus domain.WalletTransactionStatus, newVersion string, updateDate int64) error {
	tx := ctx.Value("tx").(*sql.Tx)
	sqlQuery := `UPDATE "WALLET"."WALLET_TRANSACTIONS" SET "STATUS"=$1, "VERSION"=$2, "LAST_UPDATE_DATE"=$3 WHERE "ID"=$4 AND "WALLET_ID"=$5 AND "VERSION"=$6`
	result, err := tx.ExecContext(ctx, sqlQuery, newStatus, newVersion, updateDate, walletTransaction.ID, walletTransaction.WalletID, walletTransaction.Version)
	if err != nil {
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows > 0 {
		walletTransaction.Status = newStatus
		walletTransaction.Version = newVersion
		walletTransaction.LastUpdateDate = updateDate
		return nil
	}
	return domain.ErrNoAffectedRows
}

func (repo *walletTransactionRepo) generateID(ctx context.Context) (uint64, error) {
	sqlQuery := `SELECT nextval('"WALLET"."WALLET_TRANSACTIONS_ID_SEQ"')`
	sqlRow := repo.db.QueryRowContext(ctx, sqlQuery)
	var id uint64
	err := sqlRow.Scan(&id)
	return id, err
}
