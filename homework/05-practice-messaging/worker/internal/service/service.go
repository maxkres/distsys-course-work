package service

import (
	"context"
	"fmt"
	"hash/fnv"
	"os"
	"strings"
)

type Service struct {
	deliver chan string
	mb      MessageBrokerInterface
}

type MessageBrokerInterface interface {
	Publish(ctx context.Context, body string) error
}

func New(del chan string, mb MessageBrokerInterface) Service {
	return Service{
		deliver: del,
		mb:      mb,
	}
}

func (svc Service) StoreImage() {
	go func() {
		for message := range svc.deliver {
			str := strings.Split(message, ":")
			hash := fnv.New64a()
			hash.Write([]byte(str[1]))
			hashValue := hash.Sum64()
			result := fmt.Sprint(hashValue)

			file, _ := os.Create("data/" + str[0])

			file.Write([]byte(result))

			svc.mb.Publish(context.Background(), str[0])
		}
	}()
}
