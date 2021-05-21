package repo

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/nazarov-pro/stock-exchange/services/wallet/pkg/domain"
)

type walletOperationRepo struct {
	db     *sql.DB
	logger log.Logger
}

//NewWalletOperationRepo a new instance of wallet transaction repo
func NewWalletOperationRepo(db *sql.DB, logger log.Logger) domain.WalletOperationRepository {
	return &walletOperationRepo{db: db, logger: logger}
}

func (repo *walletOperationRepo) Save(ctx context.Context, walletOperation *domain.WalletOperation) error {
	tx := ctx.Value("tx").(*sql.Tx)
	id, err := repo.generateID(ctx)
	if err != nil {
		return err
	}
	walletOperation.ID = id

	var result sql.Result
	if walletOperation.TransactionID == 0 {
		sqlQuery := `INSERT INTO "WALLET"."WALLET_OPERATIONS"("ID", "WALLET_ID", "TYPE", "CREATED_DATE") values ($1, $2, $3, $4)`
		result, err = tx.ExecContext(ctx, sqlQuery, walletOperation.ID, walletOperation.WalletID, walletOperation.Type, walletOperation.CreatedDate)
	} else {
		sqlQuery := `INSERT INTO "WALLET"."WALLET_OPERATIONS"("ID", "WALLET_ID", "TRANSACTION_ID", "TYPE", "CREATED_DATE") values ($1, $2, $3, $4, $5)`
		result, err = tx.ExecContext(ctx, sqlQuery, walletOperation.ID, walletOperation.WalletID, walletOperation.TransactionID, walletOperation.Type, walletOperation.CreatedDate)
	}

	if err != nil {
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows > 0 {
		return repo.saveOperationValues(ctx, walletOperation)
	}
	return domain.ErrNoAffectedRows
}

func (repo *walletOperationRepo) saveOperationValues(ctx context.Context, walletOperation *domain.WalletOperation) error {
	tx := ctx.Value("tx").(*sql.Tx)
	sqlQuery := `INSERT INTO "WALLET"."WALLET_OPERATION_VALUES"("WALLET_OPERATION_ID", "KEY", "VALUE") values ($1, $2, $3)`
	for key, value := range walletOperation.Values {
		_, err := tx.ExecContext(ctx, sqlQuery, walletOperation.ID, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *walletOperationRepo) generateID(ctx context.Context) (uint64, error) {
	sqlQuery := `SELECT nextval('"WALLET"."WALLET_OPERATIONS_ID_SEQ"')`
	sqlRow := repo.db.QueryRowContext(ctx, sqlQuery)
	var id uint64
	err := sqlRow.Scan(&id)
	return id, err
}
