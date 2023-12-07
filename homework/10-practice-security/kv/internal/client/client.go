package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"store/config"
	"store/pkg/proto"
)

type Client struct {
	client proto.AuthClient
	conn   *grpc.ClientConn
	cfg    *config.Config
}

func New(cfg *config.Config) (*Client, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(cfg.GrpcHost, opts...)

	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	client := proto.NewAuthClient(conn)

	return &Client{client, conn, cfg}, nil
}

func (u *Client) CheckJWT(ctx context.Context, token string) (string, error) {
	request := &proto.Token{
		Token: token,
	}
	response, err := u.client.CheckJWT(ctx, request)
	if err != nil {
		return "", fmt.Errorf("get jwt failed: %w", err)
	}
	return response.Username, nil
}

func (u *Client) Close() error {
	return u.conn.Close()
}
