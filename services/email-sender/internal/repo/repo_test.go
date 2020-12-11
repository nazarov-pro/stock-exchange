package repo

import (
	"database/sql"
	"fmt"
	"testing"
	"os"
	"context"

	_ "github.com/lib/pq"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/nazarov-pro/stock-exchange/services/email-sender/domain"
)

func TestInsert(t *testing.T) {
	msg := domain.Message{
		Content:     "cnt",
		Subject:     "sbj",
		Sender:      "sender",
		Status:      domain.Sent,
		CreatedDate: 0,
		Recipients:  []string{"1", "2"},
	}
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "secret", "postgres",
	)
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "account",
			"ts", log.DefaultTimestampUTC,
			"clr", log.DefaultCaller,
		)
	}

	db, _ := sql.Open("postgres", psqlInfo)
	repo, _ := New(db, logger)
	err := repo.Insert(context.Background(), &msg)

	if err != nil {
		t.Fatalf("Something went wrong when saving email msg, err: %v", err)
	}
}
