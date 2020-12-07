package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/nazarov-pro/stock-exchange/services/email-sender"
)

// PgRepository represents PostgreSql repository
type PgRepository struct {
	db     *sql.DB
	logger log.Logger
}

// New created a new repository backed by postgressql
func New(db *sql.DB, logger log.Logger) (email.Repository, error) {
	return &PgRepository{
		db:     db,
		logger: log.With(logger, "repo", "postgresql"),
	}, nil
}

// Insert inserting to the db
func (repo *PgRepository) Insert(ctx context.Context, msg *email.Message) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	sqlQuery := `
		INSERT INTO "EMAIL"."EMAIL_MESSAGES"("ID", "SUBJECT", "SENDER", "CONTENT", "STATUS", "CREATED_DATE")
		VALUES (nextval('"EMAIL"."EMAIL_MESSAGES_ID_SEQ"'), $1, $2, $3, $4, $5) RETURNING "ID"
	`
	raw := repo.db.QueryRowContext(ctx, sqlQuery, msg.Subject, msg.Sender, msg.Content, msg.Status, msg.CreatedDate)
	err = raw.Scan(&msg.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	sqlQuery = `
	INSERT INTO "EMAIL"."EMAIL_MESSAGE_RECIPIENTS"("EMAIL_MSG_ID", "RECIPIENT")
	VALUES `

	params := make([]interface{}, len(msg.Recipients)+1)
	params[0] = msg.ID

	for i, recepient := range msg.Recipients {
		params[i+1] = recepient
		sqlQuery += fmt.Sprintf("($1, $%d),", i+2)
	}

	sqlQuery = sqlQuery[:len(sqlQuery)-1]
	_, err = repo.db.ExecContext(ctx, sqlQuery, params...)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
