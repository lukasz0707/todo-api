package utils

import "github.com/gofiber/fiber/v2"

func ErrorResponse(c *fiber.Ctx, status int, message string) error {

	return c.Status(status).JSON(&fiber.Map{
		"error": message,
	})

}
