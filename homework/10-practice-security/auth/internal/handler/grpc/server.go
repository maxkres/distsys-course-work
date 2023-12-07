package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"auth/config"
	"auth/internal/service"
	"auth/pkg/proto"
	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedAuthServer
	listener   net.Listener
	grpcServer *grpc.Server
	service    *service.Service
	cfg        *config.Config
}

func New(s *service.Service, cfg *config.Config) *Server {
	return &Server{service: s, cfg: cfg}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.cfg.GrpcHost)

	if err != nil {
		return fmt.Errorf("listen failed: %w", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	s.listener = listener
	s.grpcServer = grpcServer

	proto.RegisterAuthServer(grpcServer, s)
	err = grpcServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("serve failed: %w", err)
	}

	return nil
}

func (s *Server) CheckJWT(c context.Context, token *proto.Token) (*proto.Response, error) {
	username, err := s.service.Verify(token.Token)
	if err != nil {
		return nil, err
	}
	return &proto.Response{Username: username}, nil
}

func (s *Server) Stop() error {
	log.Println("Shuttig down grpc...")

	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("listener close failed: %w", err)
	}

	s.grpcServer.Stop()
	log.Println("Grpc server exiting.")
	return nil
}
