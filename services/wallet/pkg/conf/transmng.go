package conf

import (
	"database/sql"

	"github.com/go-kit/kit/log"
)

type TransactionManager interface {
	Begin() (*sql.Tx, error)
}

type transactionManager struct {
	db     *sql.DB
	logger log.Logger
}

func NewTransactionManager(db *sql.DB, logger log.Logger) TransactionManager {
	return &transactionManager{db: db, logger: logger}
}

func (trMng transactionManager) Begin() (*sql.Tx, error) {
	return trMng.db.Begin()
}
