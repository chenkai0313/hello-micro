package app

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

var Config appConfig

type appConfig struct {
	MysqlDbDns string `mapstructure:"mysql_db_dns"`
	Server     Server `mapstructure:"server"`
	Jaeger     Jaeger `mapstructure:"jaeger"`
	Redis      Redis  `mapstructure:"redis"`
}

type Server struct {
	Version string `mapstructure:"version"`
	Name    string `mapstructure:"name"`
}

type Redis struct {
	HostName string `mapstructure:"hostname"`
	DB       int    `mapstructure:"database"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type Jaeger struct {
	Addr string `mapstructure:"addr"`
}

func LoadConfig(configPaths ...string) {
	v := viper.New()

	if os.Getenv("ENV") == "prod" {
		v.SetConfigName("prod")
		v.SetConfigType("yaml")
	} else {
		v.SetConfigName("test")
		v.SetConfigType("yaml")
	}
	v.SetDefault("mysql_db_dns", "127.0.0.1")

	v.SetDefault("redis.hostname", "127.0.0.1")
	v.SetDefault("redis.database", 1)
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")

	v.SetDefault("server.name", "test")
	v.SetDefault("server.version", "latest")

	v.SetDefault("jaeger.addr", "")

	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	if err := v.ReadInConfig(); err != nil {
		log.Panic(fmt.Errorf("config error failed to read the configuration file: %s", err))
	}
	if err := v.Unmarshal(&Config); err != nil {
		log.Panic("config error", err)
	}
}
