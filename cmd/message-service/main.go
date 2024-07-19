package main

import (
	"context"
	"messaggio/infra"
	"messaggio/internal/api"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	// Check if the process is a master process
	if os.Getenv("FIBER_PREFORK_CHILD") == "" {
		// Init config
		i := infra.New("config/config.json")
		// Set project mode
		i.SetMode()

		// Get custom logrus logger
		log := i.GetLogger()

		// Connect to database and migration
		psqlClient, err := i.PSQLClient()
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}
		defer psqlClient.Close()
		log.Info("Connected to PSQLClient")

		// Setup signal handling
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		// Start API server
		server := api.NewServer(i)
		go func() {
			if err := server.Run(ctx); err != nil && err != context.Canceled {
				log.Fatalf("API server error: %v", err)
			}
		}()

		// Wait for termination signal
		<-ctx.Done()
		log.Info("Shutting down gracefully...")
		stop()
	} else {
		// Child process: Just run the server
		server := api.NewServer(infra.New("config/config.json"))
		if err := server.Run(context.Background()); err != nil && err != context.Canceled {
			logrus.Fatalf("API server error: %v", err)
		}
	}
}
