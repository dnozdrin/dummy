package options

type Config struct {
	LogLevel        string
	Port            int
	HealthCheckPort int
	PrometheusPort  int
	SQLDB           SQLDBConfig
	CacheAddr       string
}

type SQLDBConfig struct {
	Host         string
	Port         int
	User         string
	Pass         string
	DBName       string
	MaxOpenConns int
}
