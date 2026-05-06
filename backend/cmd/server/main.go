package main

import (
	"log"

	"github.com/milad176/go-realtime-chat/backend/internal/api"
	"github.com/milad176/go-realtime-chat/backend/internal/config"
	"github.com/milad176/go-realtime-chat/backend/internal/db"
)

func main() {
	log.Println("Starting Go RealTime Chat Backend...")

	cfg := config.LoadConfig()

	pg, err := db.PostgresConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(pg)

	log.Printf("HTTP server listening on :%s\n", cfg.ServerPort)

	if err := server.Start(cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
