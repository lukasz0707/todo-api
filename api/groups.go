package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	db "github.com/lukasz0707/todo-api/db/sqlc"
	"github.com/lukasz0707/todo-api/token"
	"github.com/lukasz0707/todo-api/utils"
)

type CreateGroupRequest struct {
	GroupName string `json:"group_name" validate:"required"`
}

func (server *Server) createGroup(c *fiber.Ctx) error {
	var req CreateGroupRequest
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
	result, err := server.store.CreateGroupTx(c.Context(), db.CreateGroupTxParams{
		UserID:    payload.UserID,
		GroupName: req.GroupName,
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(result)

}
