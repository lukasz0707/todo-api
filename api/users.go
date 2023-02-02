package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/lukasz0707/todo-api/db/sqlc"
	"github.com/lukasz0707/todo-api/utils"
)

type createUserRequest struct {
	Username string `json:"username" validate:"required,alphanumunicode,min=5,max=30"`
	Password string `json:"password" validate:"required,min=8,alphanumunicode"`
	Email    string `json:"email" validate:"required,email"`
}

// type createUserResponse struct {
// 	Username string `json:"username"`
// 	FullName string `json:"full_name"`
// 	Email    string `json:"email"`
// }

func (server *Server) createUser(c *fiber.Ctx) error {
	var req createUserRequest
	err := c.BodyParser(&req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "cannot parse json")
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
	return c.JSON(user)
}

type getUserRequest struct {
	ID *int64 `validate:"required,min=0"`
}

type getUserResponse struct {
	ID                int64              `json:"id"`
	Username          string             `json:"username"`
	Email             string             `json:"email"`
	PasswordChangedAt pgtype.Timestamptz `json:"password_changed_at"`
	CreatedAt         pgtype.Timestamptz `json:"created_at"`
}

func (server *Server) getUser(c *fiber.Ctx) error {
	n, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	req := getUserRequest{ID: &n}
	if err := utils.Validate(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	user, err := server.store.GetUser(c.Context(), *req.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	resp := getUserResponse{
		ID:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	return c.JSON(resp)

}
