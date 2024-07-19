package api

import (
	"context"
	"messaggio/infra"
	"messaggio/internal/api/handlers"
	"messaggio/internal/usecase"
	"messaggio/pkg/http/middleware"
	"messaggio/pkg/http/request"
	"messaggio/pkg/util/logger"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type Server interface {
	Run(ctx context.Context) error
}

type server struct {
	infra      infra.Infra
	app        *fiber.App
	service    usecase.ServiceUseCase
	middleware middleware.Middleware
}

func NewServer(infra infra.Infra) Server {
	return &server{
		infra: infra,
		app: fiber.New(fiber.Config{
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  30 * time.Second,
			Concurrency:  1000000,
		}),
		service:    usecase.NewServiceUseCase(infra),
		middleware: middleware.NewMiddleware(infra.Config().GetString("secret.key")),
	}
}

func (s *server) Run(ctx context.Context) error {
	// Setup middleware and routes
	s.handlers()
	s.v1()

	// Run the server in a goroutine
	go func() {
		if err := s.app.Listen(s.infra.Port(), fiber.ListenConfig{EnablePrefork: true}); err != nil {
			if err.Error() != "http: Server closed" {
				logger.GetLogger().Fatalf("Server error: %v", err)
			}
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	// Graceful shutdown with a timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Graceful shutdown with a timeout
	err := s.app.Shutdown()
	if err != nil {
		logger.GetLogger().Errorf("Server shutdown error: %v", err)
	}

	// Wait for a while to ensure all connections are closed
	select {
	case <-shutdownCtx.Done():
		logger.GetLogger().Info("Server stopped gracefully")
		return nil
	case <-time.After(5 * time.Second):
		logger.GetLogger().Error("Server shutdown timed out")
		return nil
	}
}

func (s *server) handlers() {
	h := request.DefaultHandler()

	s.app.Use(recover.New()) // Включаем middleware для восстановления после паники

	s.app.Use(func(ctx fiber.Ctx) error {
		if ctx.Route().Path == "*" {
			return h.NoRoute(ctx)
		}
		return ctx.Next()
	})

	s.app.Get("/", func(ctx fiber.Ctx) error {
		return h.Index(ctx)
	})
}

func (s *server) v1() {
	messageHandler := handlers.NewMessageHandler(s.service.MessageService(), s.infra.KafkaProducer())

	api := s.app.Group("/api")
	{
		message := api.Group("/message")
		{
			message.Post("/", messageHandler.CreateMessage)
			message.Get("/stat", messageHandler.GetStatistics)
		}
	}
}
