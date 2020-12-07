package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/lib/pq"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/repository"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/implementation"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/consumer"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/config"
)

func main() {
	config := config.Config
	
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.GetString("db.host"), config.GetInt("db.port"), 
		config.GetString("db.user"), config.GetString("db.password"), 
		config.GetString("db.database"),
	)
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", config.GetString("app.name"),
			"ts", log.DefaultTimestampUTC,
			"clr", log.DefaultCaller,
		)
	}

	db, _ := sql.Open("postgres", psqlInfo)
	repo, _ := repository.New(db, logger)
	svc := implementation.New(repo, logger)
	consumer.ConsumeEmails(svc)
}
