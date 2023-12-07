package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	RedisDbHost     string `mapstructure:"REDIS_DB_HOST"`
	RedisDbPassword string `mapstructure:"REDIS_DB_PASSWORD"`
	RedisDbName     int    `mapstructure:"REDIS_DB_NAME"`

	GrpcHost string `mapstructure:"GRPC_HOST"`

	PubPath string `mapstructure:"PUB_PATH"`
}

func New() (*Config, error) {
	viper.AddConfigPath("./config")
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}
	return config, nil
}
