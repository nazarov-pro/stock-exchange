package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/oklog/run"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nazarov-pro/stock-exchange/services/account/pkg/conf"
	httptransport "github.com/nazarov-pro/stock-exchange/services/account/pkg/http"
	accountsvc "github.com/nazarov-pro/stock-exchange/services/account/pkg/impl"
	"github.com/nazarov-pro/stock-exchange/services/account/pkg/repo"
	"github.com/nazarov-pro/stock-exchange/services/account/pkg/transport"
)

var startTime = time.Now()

func main() {
	var (
		config   = conf.Config
		httpAddr = fmt.Sprintf(
			"%s:%d",
			config.GetString("server.hostname"), config.GetInt("server.port"),
		)
		appName = config.GetString("app.name")
	)

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

	db, err := conf.ConnectDb()
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(-1)
	}

	// Service Initalization
	repository, err := repo.New(db, logger)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(-1)
	}
	svc := accountsvc.New(repository, logger)

	endpoints := transport.MakeEndpoints(&svc)

	var h http.Handler
	{
		ocTracing := kitoc.HTTPServerTrace()
		serverOptions := []kithttp.ServerOption{ocTracing}
		h = httptransport.NewService(endpoints, serverOptions, logger)
	}

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

	//HTTP Server Runner
	{
		ln, _ := net.Listen("tcp", httpAddr)
		level.Info(logger).Log("transport", "tcp", "addr", httpAddr)
		g.Add(func() error {
			return http.Serve(ln, h)
		}, func(error) {

			ln.Close()
		})
	}

	defer os.Remove(appName)
	pid := os.Getpid()
	ioutil.WriteFile(appName, []byte(fmt.Sprint(pid)), 0644)

	level.Info(logger).Log("msg", "service started", "startuptime", time.Since(startTime))
	defer level.Info(logger).Log("msg", "service ended")

	level.Error(logger).Log("err", g.Run())
}
