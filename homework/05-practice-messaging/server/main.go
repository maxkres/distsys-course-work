package main

import (
	"errors"
	"log"
	"net/http"
	"server/config"
	"server/internal/handlers"
	"server/internal/rabbit"
	"server/internal/server"
	"server/internal/service"
	"time"
)

func main() {
	time.Sleep(10 * time.Second)
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config new failed: %v", err)
	}

	delivery := make(chan string)

	rb, err := rabbit.New(cfg, delivery)
	if err != nil {
		log.Fatalf("rabbit new failed: %v", err)
	}

	rb.Receive()

	svc := service.New(rb)
	handler := handlers.New(svc, delivery)

	srv := &server.Server{}

	go func() {
		if err := srv.Run(handler.InitRouters(), cfg); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server run failed: %v", err)
			return
		}
	}()

	if err := srv.WaitForShutDown(); err != nil {
		log.Fatalf("server shut down failed: %v", err)
	}
}
