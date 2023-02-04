package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	db "github.com/lukasz0707/todo-api/db/sqlc"
)

type Server struct {
	store  db.Store
	Router *fiber.App
}

func NewServer(store db.Store) (*Server, error) {
	server := &Server{
		store: store,
	}
	app := server.setupRouter()
	server.Router = app

	return server, nil
}

func (server *Server) setupRouter() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/metrics", monitor.New(monitor.Config{Title: "TodoApi Metrics Page"}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Everything allright :)")
	})
	app.Get("/users/:id", server.getUser)
	app.Post("/users", server.createUser)

	return app
}
