package main

import (
	"log"

	"store/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("app run failed: %v", err)
	}
}
