package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AccessTokenExp  int    `mapstructure:"ACCESS_TOKEN_EXP"`
	RefreshTokenExp int    `mapstructure:"REFRESH_TOKEN_EXP"`
	Hs256Secret     string `mapstructure:"HS256_SECRET"`

	RedisDbHost     string `mapstructure:"REDIS_DB_HOST"`
	RedisDbPassword string `mapstructure:"REDIS_DB_PASSWORD"`
	RedisDbName     int    `mapstructure:"REDIS_DB_NAME"`

	GrpcHost string `mapstructure:"GRPC_HOST"`

	PemPath string `mapstructure:"PEM_PATH"`
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
