package options

import (
	"github.com/akhripko/dummy/remote/grpc/hellosrv"
	"github.com/akhripko/dummy/storage/postgres"
)

type Config struct {
	LogLevel        string
	HTTPPort        int
	GraphqlPort     int
	GRPCPort        int
	HealthCheckPort int
	PrometheusPort  int
	SQLDB           postgres.Config
	CacheAddr       string
	HelloSrvConf    hellosrv.Config
}

type SQLDBConfig struct {
	Host         string
	Port         int
	User         string
	Pass         string
	DBName       string
	MaxOpenConns int
}
