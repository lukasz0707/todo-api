package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ErrorResponse(c *fiber.Ctx, status int, message string) error {

	return c.Status(status).JSON(&fiber.Map{
		"error": message,
	})

}

func ErrorWrapper(message string, err error) string {
	wrappedError := fmt.Errorf("%s: %w", message, err)
	return wrappedError.Error()
}
