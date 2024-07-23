package middleware

import (
	"context"
	"messaggio/pkg/util/logger"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/sirupsen/logrus"
)

type Middleware interface {
	CORS() fiber.Handler
	RateLimit(rps int) fiber.Handler
	IPBlock(redisClient *redis.Client, config IPBlockConfig) fiber.Handler
}

type middleware struct {
	secretKey string
	log       logger.Logger
}

type IPBlockConfig struct {
	Enable          bool
	RateLimitWindow time.Duration
	MaxRequests     int
}

func NewMiddleware(secretKey string) Middleware {
	return &middleware{
		secretKey: secretKey,
		log:       logger.GetLogger(),
	}
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
	//TODO: Настроить CORS
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	})
}

// RPSLimit returns a middleware handler that limits the requests per second (RPS).
//
// It uses a rate limiter initialized with the provided RPS value. For each request,
// it logs the time elapsed since the previous request using the rate limiter.
func (m *middleware) RateLimit(rps int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        rps,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
	})
}

func (m *middleware) IPBlock(redisClient *redis.Client, config IPBlockConfig) fiber.Handler {
	if !config.Enable {
		return func(c fiber.Ctx) error {
			logrus.Info("test")
			return c.Next()
		}
	}

	return func(c fiber.Ctx) error {
		logrus.Info("test")
		ip := c.IP()
		key := "rate_limit:" + ip
		blockedKey := "blocked_ips"
		blockedTTL := 10 * time.Minute

		count, err := redisClient.Incr(context.Background(), key).Result()
		if err != nil {
			m.log.Errorf("error icrementing Redis key: %v", err)
			return c.SendStatus(501)
		}

		if count == 1 {
			redisClient.Expire(context.Background(), key, config.RateLimitWindow)
		}

		if count > int64(config.MaxRequests) {
			_, err := redisClient.SAdd(context.Background(), blockedKey, ip).Result()
			if err != nil {
				m.log.Errorf("error add IP to block list: %v", err)
				return c.SendStatus(501)
			}

			_, err = redisClient.Expire(context.Background(), blockedKey, blockedTTL).Result()
			if err != nil {
				m.log.Errorf("error setting expiration for blocked IP list: %v", err)
				return c.SendStatus(501)
			}

			return c.Status(429).SendString("Rate limit exceeded")
		}

		isBlock, err := redisClient.SIsMember(context.Background(), blockedKey, ip).Result()
		if err != nil {
			m.log.Errorf("error check if IP is block: %v", err)
			return c.SendStatus(501)
		}
		if isBlock {
			return c.Status(fiber.StatusForbidden).SendString("Access denied")
		}

		return c.Next()
	}
}
