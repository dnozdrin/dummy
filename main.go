package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/akhripko/dummy/healthcheck"
	"github.com/akhripko/dummy/metrics"
	"github.com/akhripko/dummy/options"
	"github.com/akhripko/dummy/prometheus"
	"github.com/akhripko/dummy/service"
	"github.com/akhripko/dummy/storage/cache/redis"
	"github.com/akhripko/dummy/storage/sql/postgres"
)

func main() {
	// read service config from os env
	config := options.ReadEnv()

	// init logger
	initLogger(config)

	log.Info("begin...")
	// register metrics
	metrics.Register()

	// prepare main context
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)
	var wg = &sync.WaitGroup{}

	// build db
	db, err := postgres.New(ctx, postgres.Config(config.SQLDB))
	if err != nil {
		log.Error("sql db init error:", err.Error())
		os.Exit(1)
	}
	// build cache
	ccl, err := redis.New(ctx, config.CacheAddr)
	if err != nil {
		log.Error("cache init error:", err.Error())
		os.Exit(1)
	}

	// build main service
	srv := service.New(config.Port, db, ccl)
	// build prometheus service
	prometheusSrv := prometheus.New(config.PrometheusPort)
	// build healthcheck service
	healthSrv := healthcheck.New(config.HealthCheckPort, srv.HealthCheck, prometheusSrv.HealthCheck)

	// run service
	healthSrv.Run(ctx, wg)
	prometheusSrv.Run(ctx, wg)
	srv.Run(ctx, wg)

	// wait while services work
	wg.Wait()
	log.Info("end")
}

func initLogger(config *options.Config) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)

	switch strings.ToLower(config.LogLevel) {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Error("Got Interrupt signal")
		stop()
	}()
}
