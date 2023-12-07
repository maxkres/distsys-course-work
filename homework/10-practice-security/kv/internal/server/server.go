package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"store/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(host string, handler http.Handler, cfg *config.Config) error {
	s.httpServer = &http.Server{
		Addr:    host,
		Handler: handler,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) ShutDown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shuttig down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("shut down failed: %w", err)
	}

	log.Println("Server exiting.")
	return nil
}
