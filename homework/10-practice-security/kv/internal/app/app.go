package app

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"store/config"
	"store/internal/handler"
	"store/internal/repo/redis"
	"store/internal/server"
	"store/internal/service"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config new failed: %w", err)
	}

	var public, port string
	flag.StringVar(&public, "public", "signature.pub", "verbose output")
	flag.StringVar(&port, "port", "8081", "verbose output")
	flag.Parse()

	cfg.PubPath = public

	redis, err := redis.New(cfg)
	if err != nil {
		return fmt.Errorf("redis new failed: %w", err)
	}
	defer redis.Close()

	service := service.New(redis, cfg)
	handler, err := handler.New(service, cfg)
	if err != nil {
		return fmt.Errorf("handler new failed: %w", err)
	}

	server := &server.Server{}

	go func() {
		if err := server.Run(fmt.Sprintf("0.0.0.0:%s", port), handler.InitRouters(), cfg); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf(fmt.Sprintf("server run failed: %v", err))
			return
		}
	}()

	if err := server.ShutDown(); err != nil {
		return fmt.Errorf("server shut down failed: %w", err)
	}
	return nil
}
