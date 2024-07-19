package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type Middleware interface {
	CORS() fiber.Handler
	RPSLimit(rps int) fiber.Handler
}

type middleware struct {
	secretKey string
}

func NewMiddleware(secretKey string) Middleware {
	return &middleware{secretKey: secretKey}
}

// CORS returns a middleware handler that adds CORS headers to the response.
//
// It sets the following headers:
//   - Access-Control-Allow-Origin: *
//   - Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
//   - Access-Control-Allow-Headers: Origin, Content-Type, Authorization
//   - Access-Control-Expose-Headers: Content-Length
//   - Access-Control-Allow-Credentials: true
//
// If the incoming request method is OPTIONS, it responds with HTTP status
// 204 (No Content) and aborts further processing.
func (m *middleware) CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
	})
}

// RPSLimit returns a middleware handler that limits the requests per second (RPS).
//
// It uses a rate limiter initialized with the provided RPS value. For each request,
// it logs the time elapsed since the previous request using the rate limiter.
func (m *middleware) RPSLimit(rps int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        rps,
		Expiration: 60 * time.Second,
	})
}
