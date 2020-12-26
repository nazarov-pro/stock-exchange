package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"

	_ "github.com/lib/pq"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/pkg/conf"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/pkg/kafka"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/pkg/repo"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/pkg/svc"
)

var startTime = time.Now()

func main() {
	config := conf.Config
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

	db, err := conf.ConnectDb()
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(0)
	}

	repo := repo.New(db, logger)
	svc := svc.New(repo, logger)
	emailConsumer := kafka.NewEmailConsumer(svc)
	ctx := context.Background()

	var g run.Group
	// Signal Catcher
	{
		sigChan := make(chan os.Signal)
		g.Add(func() error {
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			return fmt.Errorf("%v", <-sigChan)
		}, func(error) {
			close(sigChan)
		})
	}

	// Email Consumer
	{
		ctx, cancel := context.WithCancel(ctx)
		g.Add(func() error {
			return emailConsumer.Consume(ctx)
		}, func(error) {
			cancel()
		})
	}

	// Email Consumer 2
	{
		ctx, cancel := context.WithCancel(ctx)
		g.Add(func() error {
			return emailConsumer.Consume(ctx)
		}, func(error) {
			cancel()
		})
	}

	pid := os.Getpid()
	ioutil.WriteFile(appName, []byte(fmt.Sprint(pid)), 0644)

	level.Info(logger).Log("msg", "service started", "startuptime", time.Since(startTime))
	if err = g.Run(); err != nil {
		level.Error(logger).Log("err", err)
	}

	level.Info(logger).Log("msg", "service ended")
	os.Remove(appName)
}
