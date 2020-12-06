package repository

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/nazarov-pro/stock-exchange/services/account"
)

// PgRepository represents PostgreSql repository
type PgRepository struct {
	db     *sql.DB
	logger log.Logger
}

// New created a new repository backed by postgressql
func New(db *sql.DB, logger log.Logger) (account.Repository, error) {
	return &PgRepository{
		db:     db,
		logger: log.With(logger, "repo", "postgresql"),
	}, nil
}

// Create - inserting a new account to the database
func (repo *PgRepository) Create(ctx context.Context, acc *account.Account) error {
	sql := `
			INSERT INTO "ACCOUNT"."ACCOUNTS"("ID", "USERNAME", "EMAIL", "PASSWORD", "STATUS", "ACTIVATION_CODE", "CREATED_DATE")
			values (nextval('"ACCOUNT"."ACCOUNT_ID_SEQ"'), $1, $2, $3, $4, $5, $6)
	`
	_, err := repo.db.ExecContext(ctx, sql, acc.Username, acc.Email, acc.Password, acc.Status, acc.ActivationCode, acc.CreatedDate)
	if err != nil {
		return err
	}
	return nil
}

// FindByUsernameOrEmail - Find Account by username or email
func (repo *PgRepository) FindByUsernameOrEmail(ctx context.Context, username string, email string) (*account.Account, error) {
	sqlQuery := `
				SELECT "A"."ID", "A"."USERNAME", "A"."EMAIL", "A"."PASSWORD", 
				"A"."STATUS", "A"."ACTIVATION_CODE" FROM "ACCOUNT"."ACCOUNTS" "A" 
				WHERE "A"."USERNAME"=$1 OR "A"."EMAIL"=$2 limit 1
	`
	sqlRow := repo.db.QueryRowContext(ctx, sqlQuery, username, email)

	var acc account.Account
	err := sqlRow.Scan(&acc.ID, &acc.Username, &acc.Email, &acc.Password,
		&acc.Status, &acc.ActivationCode)

	switch err {
	case sql.ErrNoRows:
		return nil, account.ErrAccountNotFound
	case nil:
		return &acc, nil
	default:
		return nil, err
	}
}

// UpdateStatus updating account's status
func (repo *PgRepository) UpdateStatus(ctx context.Context, email string, activationCode string,
	 fromStatus account.Status, toStatus account.Status) error {
	sqlQuery := `
		UPDATE "ACCOUNT"."ACCOUNTS" SET "STATUS"=$1 WHERE 
		"EMAIL"=$2 AND "ACTIVATION_CODE"=$3 AND "STATUS"=$4
	`
	sqlResult, err := repo.db.ExecContext(ctx, sqlQuery, toStatus, email, activationCode, fromStatus)
	if err != nil {
		return err
	}

	affectedRows, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return account.ErrAccountNotFound
	}
	return nil
}
