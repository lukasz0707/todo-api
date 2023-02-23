package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	db "github.com/lukasz0707/todo-api/db/sqlc"
	"github.com/lukasz0707/todo-api/token"
	"github.com/lukasz0707/todo-api/utils"
)

type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	Router     *fiber.App
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	app := server.setupRouter()
	server.Router = app

	return server, nil
}

func (server *Server) setupRouter() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{AllowOrigins: "http://localhost:5173", AllowCredentials: true}))
	app.Get("/v1/", func(c *fiber.Ctx) error {
		return c.SendString("All good :)")
	})

	app.Post("/v1/users", server.createUser)
	app.Post("/v1/users/login", server.loginUser)
	app.Post("/v1/tokens/renew_access", server.renewAccessToken)

	appAuth := app.Group("/v1", authMiddleware(server.tokenMaker))
	appAuth.Get("/users/:id", server.getUserByID)
	appAuth.Post("/group", server.createGroup)
	appAuth.Post("/todo", server.createTodo)
	appAuth.Get("/todo_by_user_id", server.GetTodosByUserID)
	appAuth.Get("/todo_by_group_id", server.GetTodosByGroupID)

	appAdmin := app.Group("/v1/admin", authMiddleware(server.tokenMaker), authAdmin())
	appAdmin.Get("/metrics", monitor.New(monitor.Config{Title: "TodoApi Metrics Page"}))

	return app
}
