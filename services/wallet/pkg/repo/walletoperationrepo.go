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
	sqlQuery := `INSERT INTO "WALLET"."WALLET_OPERATIONS"("ID", "WALLET_ID", "TRANSACTION_ID", "TYPE", "CREATED_DATE") values (nextval('"WALLET"."WALLET_OPERATIONS_ID_SEQ"'), $1, $2, $3, $4) RETURNING "ID"`
	sqlRow := repo.db.QueryRowContext(ctx, sqlQuery, walletOperation.WalletID, walletOperation.TransactionID, walletOperation.Type, walletOperation.CreatedDate)
	err := sqlRow.Scan(&walletOperation.ID)
	return err
}
