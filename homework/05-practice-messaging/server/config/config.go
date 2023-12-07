package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	RabbitQueueName string `env:"RABBIT_QUEUE_NAME"`
	RabbitHost      string `env:"RABBIT_HOST"`
	RabbitPort      string `env:"RABBIT_PORT"`

	ServerHost string `env:"SERVER_HOST"`
}

func New() (Config, error) {
	err := godotenv.Load("./config/.env")
	if err != nil {
		return Config{}, fmt.Errorf("load failed: %w", err)
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse failed: %w", err)
	}
	return cfg, nil
}

func (c *Config) GetRabbitUrl() string {
	return fmt.Sprintf("amqp://guest:guest@%s:%s/", c.RabbitHost, c.RabbitPort)
}
