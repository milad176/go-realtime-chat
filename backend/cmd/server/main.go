package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/milad176/go-realtime-chat/backend/internal/api"
	"github.com/milad176/go-realtime-chat/backend/internal/config"
	"github.com/milad176/go-realtime-chat/backend/internal/db"
	"github.com/milad176/go-realtime-chat/backend/internal/repository"
	"github.com/milad176/go-realtime-chat/backend/internal/ws"
)

func main() {
	log.Println("Starting Go RealTime Chat Backend...")

	cfg := config.LoadConfig()

	pg, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	db.RunMigrations(pg)
	messageRepo := repository.NewMessageRepository(pg)

	hub := ws.NewHub(messageRepo)
	go hub.Run()

	server := api.NewServer(pg, hub)

	log.Printf("HTTP server listening on :%s\n", cfg.ServerPort)

	httpServer := server.NewHTTPServer(cfg.ServerPort)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	shutdown := make(chan os.Signal, 1)

	signal.Notify(
		shutdown,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-shutdown

	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}

	log.Println("Closing PostgreSQL connection pool")
	pg.Close()

	log.Println("Application stopped gracefully")
}
