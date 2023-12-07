package main

import (
	"log"
	"sync"
	"time"
	"worker/config"
	"worker/internal/rabbit"
	"worker/internal/service"
)

func main() {
	time.Sleep(15 * time.Second)
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config new failed: %v", err)
	}

	delivery := make(chan string)

	rb, err := rabbit.New(cfg, delivery)
	if err != nil {
		log.Fatalf("rabbit new failed: %v", err)
	}

	svc := service.New(delivery, rb)

	var wg sync.WaitGroup
	wg.Add(2)

	rb.Receive()
	svc.StoreImage()

	wg.Wait()
}
