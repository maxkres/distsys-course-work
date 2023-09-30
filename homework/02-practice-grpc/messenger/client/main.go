package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/distsys-course/2021/02-practice-grpc/grpc/mes_grpc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChatMessage struct {
	Author   string    `json:"author"`
	Text     string    `json:"text"`
	SendTime time.Time `json:"sendTime"`
}

type MessengerClient struct {
	pendingMessages []ChatMessage
	pendingMutex    sync.Mutex
	grpcClient      mes_grpc.MessengerServerClient
}

func NewMessengerClient(serverAddr string) (*MessengerClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(serverAddr, opts...)

	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	client := mes_grpc.NewMessengerServerClient(conn)

	return &MessengerClient{
		grpcClient: client,
	}, nil
}

func (c *MessengerClient) ReadMessages() error {
	stream, err := c.grpcClient.ReadMessages(context.Background(), &mes_grpc.Empty{})
	if err != nil {
		return fmt.Errorf("read message failed: %w", err)
	}

	for {
		message, err := stream.Recv()

		if err == io.EOF {
			continue
		}
		if err != nil {
			return fmt.Errorf("recv failed: %w", err)
		}

		chatMsg := ChatMessage{
			Author:   message.Author,
			Text:     message.Text,
			SendTime: message.SendTime.AsTime(),
		}
		c.pendingMutex.Lock()
		c.pendingMessages = append(c.pendingMessages, chatMsg)
		c.pendingMutex.Unlock()
	}
}

func (c *MessengerClient) GetPending() (messages []ChatMessage) {
	c.pendingMutex.Lock()
	result := c.pendingMessages
	c.pendingMessages = nil
	c.pendingMutex.Unlock()
	log.Printf("%+v\n", result)
	return result
}

type MessageResponse struct {
	SendTime *string `json:"sendTime"`
	Error    *string `json:"error"`
}

func main() {
	time.Sleep(1 * time.Second)
	r := gin.Default()
	serverAddr := os.Getenv("MESSENGER_SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = "localhost:51075"
		fmt.Println("Missing MESSENGER_SERVER_ADDR variable, using default value: " + serverAddr)
	}

	client, err := NewMessengerClient(serverAddr)
	if err != nil {
		log.Fatalf("new messenger client failed: %v", err)
	}

	r.POST("/getAndFlushMessages", func(c *gin.Context) {
		c.JSON(http.StatusOK, client.GetPending())
	})

	r.POST("/sendMessage", func(c *gin.Context) {
		msg := mes_grpc.ChatMessage{}
		if err := c.BindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Errorf("bind json failed: %w", err).Error(),
			})
			return
		}

		sendTime := time.Now()
		msg.SendTime = timestamppb.New(sendTime)
		_, err := client.grpcClient.SendMessage(c, &msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Errorf("send message failed: %w", err).Error(),
			})
			return
		}
		sendTimeStr := sendTime.Format(time.RFC3339Nano)
		c.JSON(http.StatusOK, MessageResponse{SendTime: &sendTimeStr})
		return
	})

	go func() {
		err := client.ReadMessages()
		if err != nil {
			log.Printf("read messages failed: %v", err)
		}
	}()

	addr := os.Getenv("MESSENGER_HTTP_PORT")
	if addr == "" {
		addr = "0.0.0.0:8080"
		fmt.Println("Missing MESSENGER_HTTP_PORT variable, using default value: 8080")
	} else {
		addr = "0.0.0.0:" + addr
	}
	if err := r.Run(addr); err != nil {
		panic(err)
	}
}
