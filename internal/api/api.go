package api

import (
	"messaggio/infra"
	"messaggio/internal/manager"
	"messaggio/pkg/http/middleware"
	"messaggio/pkg/http/request"
	"messaggio/pkg/util/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Server interface {
	Run()
}

type server struct {
	infra      infra.Infra
	app        *fiber.App
	service    manager.ServiceManager
	middleware middleware.Middleware
}

func NewServer(infra infra.Infra) Server {
	return &server{
		infra:      infra,
		app:        fiber.New(),
		service:    manager.NewServiceManager(infra),
		middleware: middleware.NewMiddleware(infra.Config().GetString("secret.key")),
	}
}

// Run starts the server and initializes necessary middleware and handlers.
// It sets up rate limiting based on the configured RPS limit,
// enables CORS middleware, registers application handlers, and API routes.
// It also starts a background service to synchronize algorithm statuses.
// Finally, it logs the start of algorithm synchronization and listens on the configured port.
func (c *server) Run() {
	c.app.Use(c.middleware.RPSLimit(c.infra.Config().GetInt("rps_limit")))

	c.app.Use(c.middleware.CORS())
	c.handlers()
	c.v1()

	log := logger.GetLogger()
	log.Info("Start algorithm sync")

	c.app.Listen(":" + c.infra.Port())
}

// handlers sets up custom route handlers for specific routes on the server.
// It assigns default handlers for handling unknown routes and an index route.
func (c *server) handlers() {
	h := request.DefaultHandler()

	c.app.Use(func(ctx *fiber.Ctx) error {
		if ctx.Route().Path == "*" {
			return h.NoRoute(ctx)
		}
		return ctx.Next()
	})

	c.app.Get("/", func(ctx *fiber.Ctx) error {
		return h.Index(ctx)
	})
}

// v1 configures versioned API endpoints (v1) for client operations.
// It sets up routes for client management operations such as adding, updating, deleting clients,
// and updating algorithm statuses associated with clients.
func (c *server) v1() {
	c.app.Get("/docs/*", swagger.HandlerDefault)
}
