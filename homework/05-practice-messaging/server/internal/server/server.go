package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"server/config"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(handler http.Handler, cfg config.Config) error {
	s.httpServer = &http.Server{
		Addr:    cfg.ServerHost,
		Handler: handler,
	}
	fmt.Println(s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) WaitForShutDown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("shut down failed: %w", err)
	}

	return nil
}
