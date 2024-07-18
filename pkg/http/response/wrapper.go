package response

import (
	"messaggio/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type Wrapper interface {
	Write(code int, message string) error
	Error(code int, err error) error
}

type wrapper struct {
	c *fiber.Ctx
}

func New(c *fiber.Ctx) Wrapper {
	return &wrapper{c: c}
}

// Write writes a JSON response with the provided HTTP status code and message.
//
// It serializes the response into JSON format using the provided code and message.
func (w *wrapper) Write(code int, message string) error {
	return w.c.Status(code).JSON(domain.Response{Code: code, Message: message})
}

// Error writes a JSON response with the provided HTTP status code and error message.
//
// It serializes the response into JSON format using the provided code and the
// error message extracted from the error object.
func (w *wrapper) Error(code int, err error) error {
	return w.c.Status(code).JSON(domain.Response{Code: code, Message: err.Error()})
}
