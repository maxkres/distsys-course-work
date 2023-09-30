package main

import (
	"log"
	"net"
	"net/http"
	"server/internal/config"
	"server/internal/handler"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config new failed: %v", err)
	}

	h := handler.New(cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Controller)

	middleware := h.Host(mux)

	host := net.JoinHostPort(cfg.Host, cfg.Port)

	log.Printf("Server started on host: %v\n", host)
	err = http.ListenAndServe(host, middleware)
	if err != nil {
		log.Fatalf("listen and server failed: %v", err)
	}
}
