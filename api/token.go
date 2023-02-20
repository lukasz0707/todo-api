package api

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lukasz0707/todo-api/utils"
)

// type renewAccessTokenRequest struct {
// 	RefreshToken string `json:"refresh_token" validate:"required"`
// }

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(c *fiber.Ctx) error {
	// var req renewAccessTokenRequest
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "refresh_token cookie not provided")
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	session, err := server.store.GetSession(c.Context(), refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			if server.config.Environment == "development" {
				fmt.Println(utils.ErrorWrapper("session not found", err))
			}
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "invalid session id")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if session.UserID != refreshPayload.UserID {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "incorrect session user")
	}

	if session.RefreshToken != refreshToken {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "mismatched session token")
	}

	if time.Now().After(session.ExpiresAt) {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "expired session")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.UserID, "access_token", server.config.AccessTokenDuration, refreshPayload.Role)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	return c.JSON(rsp)
}
