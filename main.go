package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/akhripko/dummy/srvgql"
	"github.com/akhripko/dummy/srvgrpc"
	"github.com/akhripko/dummy/srvhttp"

	"github.com/akhripko/dummy/cache/redis"
	"github.com/akhripko/dummy/healthcheck"
	"github.com/akhripko/dummy/metrics"
	"github.com/akhripko/dummy/options"
	"github.com/akhripko/dummy/prometheus"
	"github.com/akhripko/dummy/service"
	"github.com/akhripko/dummy/storage/postgres"
	log "github.com/sirupsen/logrus"
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

	// build storage
	storage, err := postgres.New(ctx, config.SQLDB)
	if err != nil {
		log.Error("sql db init error:", err.Error())
		os.Exit(1)
	}
	// build cache
	cache, err := redis.New(ctx, config.CacheAddr)
	if err != nil {
		log.Error("cache init error:", err.Error())
		os.Exit(1)
	}

	//p, err := producer.New(ctx, wg, config.KafkaTopic.Hello, config.KafkaProducer)
	//if err != nil {
	//	log.Error("kafka producer init error:", err.Error())
	//	os.Exit(1)
	//}

	//hellosrvClient, err := hellosrv.New(ctx, config.HelloSrvConf)
	//if err != nil {
	//	log.Error("hellosrv client init error:", err.Error())
	//	os.Exit(1)
	//}

	// build main service
	srv, err := service.New(storage, cache)
	if err != nil {
		log.Error("service init error:", err.Error())
		os.Exit(1)
	}

	http, err := srvhttp.New(config.HTTPPort, srv)
	if err != nil {
		log.Error("http service init error:", err.Error())
		os.Exit(1)
	}

	grpc, err := srvgrpc.New(config.GRPCPort, srv)
	if err != nil {
		log.Error("grpc service init error:", err.Error())
		os.Exit(1)
	}

	gql, err := srvgql.New(config.GraphqlPort, srv)
	if err != nil {
		log.Error("graphql service init error:", err.Error())
		os.Exit(1)
	}

	// build prometheus srv
	prometheusSrv := prometheus.New(config.PrometheusPort)
	// build healthcheck srv
	healthSrv := healthcheck.New(
		config.HealthCheckPort,
		srv.HealthCheck,
		prometheusSrv.HealthCheck,
		http.HealthCheck,
		gql.HealthCheck,
		grpc.HealthCheck,
	)

	// run srv
	http.Run(ctx, wg)
	grpc.Run(ctx, wg)
	gql.Run(ctx, wg)
	healthSrv.Run(ctx, wg)
	prometheusSrv.Run(ctx, wg)

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
