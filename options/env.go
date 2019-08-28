package options

import (
	"github.com/spf13/viper"
)

func ReadEnv() *Config {
	viper.AutomaticEnv()

	viper.SetEnvPrefix("APP")

	viper.SetDefault("LOG_LEVEL", "DEBUG")
	viper.SetDefault("HTTP_PORT", 8080)
	viper.SetDefault("HEALTH_CHECK_PORT", 8888)
	viper.SetDefault("PROMETHEUS_PORT", 9090)

	viper.SetDefault("SQLDB_HOST", "localhost")
	viper.SetDefault("SQLDB_PORT", 5432)
	viper.SetDefault("SQLDB_USER", "")
	viper.SetDefault("SQLDB_PASS", "")
	viper.SetDefault("SQLDB_DB_NAME", "db")
	viper.SetDefault("SQLDB_MAX_OPEN_CONNS", 10)

	viper.SetDefault("CACHE_ADDR", ":6379")

	return &Config{
		LogLevel:        viper.GetString("LOG_LEVEL"),
		Port:            viper.GetInt("HTTP_PORT"),
		HealthCheckPort: viper.GetInt("HEALTH_CHECK_PORT"),
		PrometheusPort:  viper.GetInt("PROMETHEUS_PORT"),
		SQLDB: SQLDBConfig{
			Host:         viper.GetString("SQLDB_HOST"),
			Port:         viper.GetInt("SQLDB_PORT"),
			User:         viper.GetString("SQLDB_USER"),
			Pass:         viper.GetString("SQLDB_PASS"),
			DBName:       viper.GetString("SQLDB_DB_NAME"),
			MaxOpenConns: viper.GetInt("SQLDB_MAX_OPEN_CONNS"),
		},
		CacheAddr: viper.GetString("CACHE_ADDR"),
	}
}
