package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/distsys-course/2021/02-practice-grpc/grpc/mes_grpc"
	"google.golang.org/grpc"
)

type ChatServer struct {
	mes_grpc.UnsafeMessengerServerServer

	messages   []*mes_grpc.ChatMessage
	mu         sync.Mutex
	subscribed map[mes_grpc.MessengerServer_ReadMessagesServer]struct{}
}

func NewMessengerServer(serverAddr string) error {
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("listen failed: %w", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	chat := &ChatServer{
		subscribed: make(map[mes_grpc.MessengerServer_ReadMessagesServer]struct{}),
	}
	mes_grpc.RegisterMessengerServerServer(grpcServer, chat)

	fmt.Printf("gRPC server started\nHost: %v\n", serverAddr)
	err = grpcServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("serve failed: %w", err)
	}

	return nil
}

func (s *ChatServer) SendMessage(ctx context.Context, message *mes_grpc.ChatMessage) (*mes_grpc.Time, error) {
	var wg sync.WaitGroup

	sendTime := timestamppb.New(time.Now())

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.mu.Lock()
		defer s.mu.Unlock()
		if message.SendTime == nil {
			message.SendTime = sendTime
		}
		s.messages = append(s.messages, message)
		fmt.Printf("Received message from %s: %s\n", message.Author, message.Text)

		for client := range s.subscribed {
			log.Println(message)
			err := client.Send(message)
			if err != nil {
				log.Printf("Failed to send message to client: %v", err)
			}
		}
	}()

	wg.Wait()
	return &mes_grpc.Time{SendTime: sendTime}, nil
}

func (s *ChatServer) ReadMessages(empty *mes_grpc.Empty, stream mes_grpc.MessengerServer_ReadMessagesServer) error {
	go func() {
		s.mu.Lock()
		s.subscribed[stream] = struct{}{}
		s.mu.Unlock()

	}()
	select {}
}

func main() {
	serverAddr := os.Getenv("MESSENGER_SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":51075"
		fmt.Println("Missing MESSENGER_SERVER_ADDR variable, using default value: " + serverAddr)
	}

	err := NewMessengerServer(serverAddr)
	if err != nil {
		log.Fatalf("new messenger server failed: %v", err)
	}
}
