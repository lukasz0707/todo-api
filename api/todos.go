package api

// import (
// 	"fmt"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	db "github.com/lukasz0707/todo-api/db/sqlc"
// 	"github.com/lukasz0707/todo-api/utils"
// )

// type CreateTodoRequest struct {
// 	UserID    int64     `json:"user_id" validate:"required"`
// 	TodoName  string    `json:"todo_name" validate:"required"`
// 	GroupName string    `json:"group_name" validate:"required"`
// 	Deadline time.Time `json:"deadline"`
// }

// func (server *Server) createTodo(c *fiber.Ctx) error {
// 	var req CreateTodoRequest
// 	err := c.BodyParser(&req)
// 	if err != nil {
// 		fmt.Println(err)
// 		return utils.ErrorResponse(c, fiber.StatusBadRequest, "cannot parse json")
// 	}
// 	if err := utils.Validate(req); err != nil {
// 		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
// 	}

// 	arg := db.CreateTodoTxParams{
// 		UserID:    req.UserID,
// 		TodoName:  req.TodoName,
// 		GroupName: req.GroupName,
// 		Deadline:  req.Deadline,
// 	}

// 	result, err := server.store.CreateTodoTx(c.Context(), arg)
// 	if err != nil {
// 		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
// 	}
// 	return c.JSON(result)
// }
