package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nazarov-pro/stock-exchange/services/account/domain"
	"github.com/nazarov-pro/stock-exchange/services/account/internal/config"
	accountsvc "github.com/nazarov-pro/stock-exchange/services/account/internal/impl"
	"github.com/nazarov-pro/stock-exchange/services/account/internal/repo"
	"github.com/nazarov-pro/stock-exchange/services/account/internal/transport"
	httptransport "github.com/nazarov-pro/stock-exchange/services/account/internal/http"
)

func main() {
	config := config.Config
	httpAddr := fmt.Sprintf(
		"%s:%d", 
		config.GetString("server.hostname"), config.GetInt("server.port"),
	)
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.GetString("db.host"), config.GetInt("db.port"), 
		config.GetString("db.user"), config.GetString("db.password"), 
		config.GetString("db.database"),
	)
	appName := config.GetString("app.name")

	// initialize our structured logger for the service
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

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	var db *sql.DB
	{
		var err error
		// Connect to the "ordersdb" database
		db, err = sql.Open("postgres", psqlInfo)

		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	// Create Account Service
	var svc domain.Service
	{
		repository, err := repo.New(db, logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		svc = accountsvc.New(repository, logger)
	}

	var endpoints transport.Endpoints
	{
		endpoints = transport.MakeEndpoints(&svc)
	}

	var h http.Handler
	{
		ocTracing := kitoc.HTTPServerTrace()
		serverOptions := []kithttp.ServerOption{ocTracing}
		h = httptransport.NewService(endpoints, serverOptions, logger)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", httpAddr)
		server := &http.Server{
			Addr:    httpAddr,
			Handler: h,
		}
		errs <- server.ListenAndServe()
	}()

	level.Error(logger).Log("exit", <-errs)
}
