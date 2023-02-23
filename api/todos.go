package api

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	db "github.com/lukasz0707/todo-api/db/sqlc"
	"github.com/lukasz0707/todo-api/token"
	"github.com/lukasz0707/todo-api/utils"
)

type CreateTodoRequest struct {
	TodoName string    `json:"todo_name" validate:"required"`
	GroupID  int64     `json:"group_id" validate:"required"`
	Deadline time.Time `json:"deadline"`
}

func (server *Server) createTodo(c *fiber.Ctx) error {
	var req CreateTodoRequest
	err := c.BodyParser(&req)
	if err != nil {
		fmt.Println(err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "cannot parse json")
	}
	if err := utils.Validate(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	payload, ok := c.Locals("authorization_payload").(*token.Payload)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Locals authorization_payload error")
	}

	arg := db.CreateTodoTxParams{
		UserID:   payload.UserID,
		TodoName: req.TodoName,
		GroupID:  req.GroupID,
		Deadline: req.Deadline,
	}

	result, err := server.store.CreateTodoTx(c.Context(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "You don't belong to that group")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(result)
}

func (server *Server) GetTodosByUserID(c *fiber.Ctx) error {

	payload, ok := c.Locals("authorization_payload").(*token.Payload)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Locals authorization_payload error")
	}

	result, err := server.store.GetTodosByUserID(c.Context(), payload.UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(result)
}

type GetTodosByGroupIDRequest struct {
	GroupID int64 `json:"group_id"`
}

func (server *Server) GetTodosByGroupID(c *fiber.Ctx) error {
	var req GetTodosByGroupIDRequest
	err := c.BodyParser(&req)
	if err != nil {
		fmt.Println(err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "cannot parse json")
	}
	if err := utils.Validate(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	payload, ok := c.Locals("authorization_payload").(*token.Payload)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Locals authorization_payload error")
	}

	result, err := server.store.GetTodosByGroupID(c.Context(), db.GetTodosByGroupIDParams{
		UserID:  payload.UserID,
		GroupID: req.GroupID,
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(result)
}
