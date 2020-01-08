package options

import (
	"github.com/akhripko/dummy/kafka/producer"
	"github.com/akhripko/dummy/providers/grpc/hellosrv"
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
	KafkaTopic      KafkaTopic
	KafkaProducer   producer.Config
}

type SQLDBConfig struct {
	Host         string
	Port         int
	User         string
	Pass         string
	DBName       string
	MaxOpenConns int
}

type KafkaTopic struct {
	Hello string
}
