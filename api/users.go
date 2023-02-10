package api

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/lukasz0707/todo-api/db/sqlc"
	"github.com/lukasz0707/todo-api/token"
	"github.com/lukasz0707/todo-api/utils"
)

type createUserRequest struct {
	Username  string `json:"username" validate:"required,alphanumunicode,min=5,max=30"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"alphanumunicode,required,min=1,max=75"`
	LastName  string `json:"last_name" validate:"alphanumunicode,required,min=1,max=75"`
	Email     string `json:"email" validate:"required,email"`
}

type userResponse struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:                user.ID,
		Username:          user.Username,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(c *fiber.Ctx) error {
	var req createUserRequest
	err := c.BodyParser(&req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, utils.ErrorWrapper("cannot parse json", err))
	}
	if err := utils.Validate(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(c.Context(), arg)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			switch pqErr.Message {
			case `duplicate key value violates unique constraint "users_username_key"`:
				return utils.ErrorResponse(c, fiber.StatusForbidden, "username already exists")
			case `duplicate key value violates unique constraint "users_email_key"`:
				return utils.ErrorResponse(c, fiber.StatusForbidden, "email already exists")
			}
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := newUserResponse(user)

	return c.JSON(resp)
}

type getUserRequest struct {
	ID *int64 `validate:"required,min=0"`
}

func (server *Server) getUserByID(c *fiber.Ctx) error {
	payload, ok := c.Locals("authorization_payload").(*token.Payload)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Locals authorization_payload error")
	}
	payloadUserID := payload.UserID

	n, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if payloadUserID != n {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "you don't have permissions to access that data")
	}

	req := getUserRequest{ID: &n}
	if err := utils.Validate(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	user, err := server.store.GetUserByID(c.Context(), *req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := newUserResponse(user)

	return c.JSON(resp)

}

type loginUserRequest struct {
	Username string `json:"username" validate:"required,alphanumunicode,min=5,max=30"`
	Password string `json:"password" validate:"required,min=8"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) loginUser(c *fiber.Ctx) error {
	var req loginUserRequest
	err := c.BodyParser(&req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, utils.ErrorWrapper("cannot parse json", err))
	}

	user, err := server.store.GetUserByUsername(c.Context(), req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			if server.config.Environment == "development" {
				fmt.Println(utils.ErrorWrapper("username not found", err))
			}
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "invalid credentials")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "invalid credentials")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.ID, "access_token", server.config.AccessTokenDuration)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.ID, "refresh_token", server.config.RefreshTokenDuration)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	session, err := server.store.CreateSession(c.Context(), db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    string(c.Context().UserAgent()),
		ClientIp:     c.IP(),
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}
	return c.JSON(rsp)
}
