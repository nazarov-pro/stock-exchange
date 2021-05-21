package repo

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/nazarov-pro/stock-exchange/services/wallet/pkg/domain"
)

type walletRepo struct {
	db     *sql.DB
	logger log.Logger
}

//NewWalletRepo a new instance of wallet repo
func NewWalletRepo(db *sql.DB, logger log.Logger) domain.WalletRepository {
	return &walletRepo{db: db, logger: logger}
}

func (repo *walletRepo) FindByID(ctx context.Context, ID uint64) (*domain.Wallet, error) {
	sqlQuery := `SELECT "W"."ID", "W"."ACCOUNT_ID", "W"."VERSION", "W"."BALANCE", "W"."CURRENCY_CODE","W"."STATUS", "W"."CREATED_DATE", COALESCE("W"."LAST_UPDATE_DATE", 0) FROM "WALLET"."WALLETS" "W" WHERE "W"."ID"=$1`
	sqlRow := repo.db.QueryRowContext(ctx, sqlQuery, ID)
	var wallet domain.Wallet
	err := sqlRow.Scan(&wallet.ID, &wallet.AccountID, &wallet.Version, &wallet.Balance, &wallet.CurrencyCode, &wallet.Status, &wallet.CreatedDate, &wallet.LastUpdateDate)
	switch err {
	case sql.ErrNoRows:
		return nil, domain.ErrWalletNotFount
	case nil:
		return &wallet, nil
	default:
		return nil, err
	}
}

func (repo *walletRepo) FindByAccountID(ctx context.Context, accountID uint64) (*[]domain.Wallet, error) {
	sqlQuery := `SELECT "W"."ID", "W"."ACCOUNT_ID", "W"."VERSION", "W"."BALANCE", "W"."CURRENCY_CODE","W"."STATUS", "W"."CREATED_DATE", COALESCE("W"."LAST_UPDATE_DATE", 0) FROM "WALLET"."WALLETS" "W" WHERE "W"."ACCOUNT_ID"=$1`
	sqlRows, err := repo.db.QueryContext(ctx, sqlQuery, accountID)
	if err != nil {
		return nil, err
	}
	defer sqlRows.Close()
	wallets := make([]domain.Wallet, 0)
	for sqlRows.Next() {
		var wallet domain.Wallet
		err = sqlRows.Scan(&wallet.ID, &wallet.AccountID, &wallet.Version, &wallet.Balance, &wallet.CurrencyCode, &wallet.Status, &wallet.CreatedDate, &wallet.LastUpdateDate)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}
	err = sqlRows.Err()

	switch err {
	case sql.ErrNoRows:
		return nil, domain.ErrWalletNotFount
	case nil:
		return &wallets, nil
	default:
		return nil, err
	}
}

func (repo *walletRepo) FindByIDAndAccountID(ctx context.Context, ID uint64, accountID uint64) (*domain.Wallet, error) {
	sqlQuery := `SELECT "W"."ID", "W"."ACCOUNT_ID", "W"."VERSION", "W"."BALANCE", "W"."CURRENCY_CODE","W"."STATUS", "W"."CREATED_DATE", COALESCE("W"."LAST_UPDATE_DATE", 0) FROM "WALLET"."WALLETS" "W" WHERE "W"."ID"=$1 AND "W"."ACCOUNT_ID"=$2`
	sqlRow := repo.db.QueryRowContext(ctx, sqlQuery, ID, accountID)
	var wallet domain.Wallet
	err := sqlRow.Scan(&wallet.ID, &wallet.AccountID, &wallet.Version, &wallet.Balance, &wallet.CurrencyCode, &wallet.Status, &wallet.CreatedDate, &wallet.LastUpdateDate)
	switch err {
	case sql.ErrNoRows:
		return nil, domain.ErrWalletNotFount
	case nil:
		return &wallet, nil
	default:
		return nil, err
	}
}

func (repo *walletRepo) Save(ctx context.Context, wallet *domain.Wallet) error {
	tx := ctx.Value("tx").(*sql.Tx)
	id, err := repo.generateID(ctx)
	if err != nil {
		return err
	}
	wallet.ID = id

	sqlQuery := `INSERT INTO "WALLET"."WALLETS"("ID", "ACCOUNT_ID", "BALANCE", "CURRENCY_CODE", "STATUS", "VERSION", "CREATED_DATE") values ($1, $2, $3, $4, $5, $6, $7)`
	result, err := tx.ExecContext(ctx, sqlQuery, wallet.ID, wallet.AccountID, wallet.Balance, wallet.CurrencyCode, wallet.Status, wallet.Version, wallet.CreatedDate)
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

func (repo *walletRepo) UpdateBalance(ctx context.Context, wallet *domain.Wallet, newVersion string, newBalance int64, updateDate int64) error {
	tx := ctx.Value("tx").(*sql.Tx)
	sqlQuery := `UPDATE "WALLET"."WALLETS" SET "BALANCE"=$1, "VERSION"=$2, "LAST_UPDATE_DATE"=$3 WHERE "ID"=$4 AND "ACCOUNT_ID"=$5 AND "VERSION"=$6`
	result, err := tx.ExecContext(ctx, sqlQuery, newBalance, newVersion, updateDate, wallet.ID, wallet.AccountID, wallet.Version)
	if err != nil {
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows > 0 {
		wallet.Balance = newBalance
		wallet.LastUpdateDate = updateDate
		wallet.Version = newVersion
		return nil
	}
	return domain.ErrNoAffectedRows
}

func (repo *walletRepo) UpdateStatus(ctx context.Context, wallet *domain.Wallet, newVersion string, newStatus domain.WalletStatus, updateDate int64) error {
	tx := ctx.Value("tx").(*sql.Tx)
	sqlQuery := `UPDATE "WALLET"."WALLETS" SET "STATUS"=$1, "VERSION"=$2, "LAST_UPDATE_DATE"=$3 WHERE "ID"=$4 AND "ACCOUNT_ID"=$5 AND "VERSION"=$6`
	result, err := tx.ExecContext(ctx, sqlQuery, newStatus, newVersion, updateDate, wallet.ID, wallet.AccountID, wallet.Version)
	if err != nil {
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows > 0 {
		wallet.Status = newStatus
		wallet.LastUpdateDate = updateDate
		wallet.Version = newVersion
		return nil
	}
	return domain.ErrNoAffectedRows

}

func (repo *walletRepo) generateID(ctx context.Context) (uint64, error) {
	sqlQuery := `SELECT nextval('"WALLET"."WALLETS_ID_SEQ"')`
	sqlRow := repo.db.QueryRowContext(ctx, sqlQuery)
	var id uint64
	err := sqlRow.Scan(&id)
	return id, err
}
