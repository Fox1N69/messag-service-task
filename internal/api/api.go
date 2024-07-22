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

type IPBlockConfig struct {
	Enable          bool
	RateLimitWindow time.Duration
	MaxRequests     int
}
type Server interface {
	Run(ctx context.Context) error
}

type server struct {
	infra         infra.Infra
	app           *fiber.App
	service       usecase.ServiceUseCase
	middleware    middleware.Middleware
	log           logger.Logger
	ipBlockConfig IPBlockConfig
}

// NewServer creates a new instance of the server with the given infrastructure configuration.
//
// This function initializes the server with the provided infrastructure, sets up
// the Fiber application with default configuration, and prepares middleware and
// service use cases required for handling HTTP requests.
//
// Parameters:
//   - infra: An instance of infra.Infra that provides configuration and dependencies
//
// Returns:
//   - Server: A new instance of a Server interface, which is implemented by the server struct
func NewServer(infra infra.Infra) Server {
	logger := logger.GetLogger()
	infraBlockConfig := infra.Config().Sub("ip_block_config")

	ipBlockConfig := IPBlockConfig{
		Enable:          infraBlockConfig.GetBool("enable"),
		MaxRequests:     infraBlockConfig.GetInt("max_requests"),
		RateLimitWindow: time.Duration(infraBlockConfig.GetInt("rate_limit_window_seconds")) * time.Second,
	}

	return &server{
		infra: infra,
		app: fiber.New(fiber.Config{
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  30 * time.Second,
			Concurrency:  10000,
		}),
		service:       usecase.NewServiceUseCase(infra),
		middleware:    middleware.NewMiddleware(infra.Config().GetString("secret.key")),
		log:           logger,
		ipBlockConfig: ipBlockConfig,
	}
}

// Run starts the HTTP server and listens for incoming requests.
//
// It sets up the necessary middleware, routes, and handlers for the server.
// The server is started in a separate goroutine to allow the main goroutine
// to handle other tasks, such as waiting for a shutdown signal.
//
// When the context is canceled, indicating a shutdown request, the method
// performs a graceful shutdown of the server. It attempts to shut down the
// server gracefully by waiting up to 5 seconds for active connections to
// close. If the shutdown does not complete within the timeout period, the
// method logs an error and returns.
//
// The method returns an error if there is an issue with starting the server
// or if the server fails to shut down gracefully.
//
// Parameters:
//   - ctx: The context used to signal when the server should stop
//
// Returns:
//   - error: An error value indicating if an error occurred during server
//     startup or shutdown
func (s *server) Run(ctx context.Context) error {
	// Setup middleware and routes
	s.middlewares()
	s.handlers()
	s.routes()

	// Run the server in a goroutine
	go func() {
		if err := s.app.Listen(s.infra.Port(), fiber.ListenConfig{EnablePrefork: true}); err != nil {
			if err.Error() != "http: Server closed" {
				s.log.Fatalf("Server error: %v", err)
			}
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Graceful shutdown with a timeout
	err := s.app.Shutdown()
	if err != nil {
		s.log.Errorf("Server shutdown error: %v", err)
	}

	select {
	case <-shutdownCtx.Done():
		s.log.Info("Server stopped gracefully")
		return nil
	case <-time.After(5 * time.Second):
		s.log.Error("Server shutdown timed out")
		return nil
	}
}

// middlewares sets up the middleware for the server.
//
// It adds middleware components to the Fiber application, including recovery
// middleware for recovering from panics and rate limiting middleware for
// controlling request rates.
//
// This method is called during server startup to ensure that the middleware
// is properly configured before handling any requests.
func (s *server) middlewares() {
	if s.ipBlockConfig.Enable {
		s.app.Use(s.middleware.IPBlock(s.infra.RedisClient(), middleware.IPBlockConfig(s.ipBlockConfig)))
	}

	s.app.Use(recover.New())
}

// handlers sets up the request handlers for the server.
//
// It configures routes and assigns handlers to various endpoints. This includes
// setting up a default handler for routes that are not matched by other handlers
// and a handler for the root path.
//
// This method is called during server startup to ensure that all request handlers
// are correctly registered and ready to process incoming requests.
func (s *server) handlers() {
	h := request.DefaultHandler()

	s.app.Use(func(ctx fiber.Ctx) error {
		if ctx.Route().Path == "*" {
			return h.NoRoute(ctx)
		}
		return ctx.Next()
	})

	s.app.Get("/", h.Index)
}

// routes sets up the API routes for the server.
//
// This method configures the routes for the API version 1 by creating route groups
// and assigning handlers to the endpoints. It sets up a route group under the `/api`
// path and further organizes the `/message` routes. Specifically, it registers
// handlers for creating messages and retrieving message statistics.
//
// It uses the `messageHandler` created with the message service and Kafka producer
// to handle requests for the `/message` endpoints. This method ensures that all
// routes related to message operations are correctly configured and ready to
// process incoming API requests.
//
// This method is called during server startup to set up the routing for the
// application.
func (s *server) routes() {
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
