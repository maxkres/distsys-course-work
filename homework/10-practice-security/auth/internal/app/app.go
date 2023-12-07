package app

import (
	"auth/config"
	"auth/internal/handler"
	"auth/internal/handler/grpc"
	"auth/internal/repo/redis"
	"auth/internal/server"
	"auth/internal/service"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config new failed: %w", err)
	}

	var private, public, port string
	flag.StringVar(&private, "private", "signature.pem", "verbose output")
	flag.StringVar(&public, "public", "signature.pub", "verbose output")
	flag.StringVar(&port, "port", "8080", "verbose output")
	flag.Parse()

	cfg.PubPath = public
	cfg.PemPath = private

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

	grpcServer := grpc.New(service, cfg)
	go func() {
		if err := grpcServer.Run(); err != nil {
			log.Println(fmt.Sprintf("grpc server run failed: %v", err))
			return
		}
	}()

	if err := server.ShutDown(); err != nil {
		return fmt.Errorf("server shut down failed: %w", err)
	}
	return nil
}
