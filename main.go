package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lukasz0707/todo-api/api"
	db "github.com/lukasz0707/todo-api/db/sqlc"
	"github.com/lukasz0707/todo-api/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	pool, err := sql.Open("pgx", config.DBSource)
	if err != nil {
		log.Fatal("unable to use data source name", err)
	}
	defer pool.Close()

	store := db.NewStore(pool)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	log.Fatal(server.Router.Listen(config.HTTPServerAddress))
}
