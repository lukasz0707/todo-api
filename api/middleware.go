package api

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lukasz0707/todo-api/token"
	"github.com/lukasz0707/todo-api/utils"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorizationHeader := c.Get(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "authorization header is not provided")
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "invalid authorization header format")
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, fmt.Sprintf("unsuported authorization type %s", authorizationType))
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		}

		if payload.TokenType != "access_token" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "wrong token type")
		}

		c.Locals(authorizationPayloadKey, payload)
		return c.Next()
	}
}

func authAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload, ok := c.Locals("authorization_payload").(*token.Payload)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Locals authorization_payload error")
		}

		if payload.Role != "admin" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "you don't have permissions to access that data")
		}
		return c.Next()
	}
}
