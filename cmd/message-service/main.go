package main

import (
	"context"
	"messaggio/infra"
	"messaggio/internal/api"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Init config
	i := infra.New("config/config.json")
	// Set project mode
	i.SetMode()
	// Get custom logrus logger
	log := i.GetLogger()

	if os.Getenv("FIBER_PREFORK_CHILD") == "" {
		// Connect to database and migration
		psqlClient, err := i.PSQLClient()
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}
		defer psqlClient.Close()
		log.Info("Connected to PSQLClient")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start API server
	server := api.NewServer(i)
	go func() {
		if err := server.Run(ctx); err != nil && err != context.Canceled {
			log.Fatalf("API server error: %v", err)
		}
	}()

	<-ctx.Done()

	if os.Getenv("FIBER_PREFORK_CHILD") == "" {
		log.Info("Shutting down gracefully...")
	}

	stop()
}
