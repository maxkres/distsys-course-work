package service

import (
	"context"
)

type Service struct {
	mb MessageBrokerInterface
}

type MessageBrokerInterface interface {
	Publish(ctx context.Context, body string) error
}

func New(mb MessageBrokerInterface) Service {
	return Service{
		mb: mb,
	}
}

func (svc Service) StoreImage(ctx context.Context, url string) error {
	return svc.mb.Publish(ctx, url)
}
