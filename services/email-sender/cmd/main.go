package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/lib/pq"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/internal/repo"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/internal/impl"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/internal/kafka/consumer"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/internal/config"
)

func main() {
	config := config.Config
	
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.GetString("db.host"), config.GetInt("db.port"), 
		config.GetString("db.user"), config.GetString("db.password"), 
		config.GetString("db.database"),
	)
	appName := config.GetString("app.name")
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", appName,
			"ts", log.DefaultTimestampUTC,
			"clr", log.DefaultCaller,
		)
	}

	closeFunc := func () {
		level.Info(logger).Log("msg", "stopped")
		time.Sleep(5 * time.Second)
	}

	defer closeFunc()

	level.Info(logger).Log("msg", "started")
	db, _ := sql.Open("postgres", psqlInfo)
	repo, _ := repo.New(db, logger)
	svc := impl.New(repo, logger)

	consumer.ConsumeEmails(svc)
}
