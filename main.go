package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lukasz0707/todo-api/api"
	db "github.com/lukasz0707/todo-api/db/sqlc"
)

func main() {
	pool, err := sql.Open("pgx", "postgresql://root:secret@localhost:5432/todo?sslmode=disable")
	if err != nil {
		log.Fatal("unable to use data source name", err)
	}
	defer pool.Close()

	store := db.NewStore(pool)
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	log.Fatal(server.Router.Listen(":3000"))
}
