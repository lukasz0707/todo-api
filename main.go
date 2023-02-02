package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lukasz0707/todo-api/api"
	db "github.com/lukasz0707/todo-api/db/sqlc"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), "postgresql://root:secret@localhost:5432/todo?sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	store := db.New(dbpool)
	server, _ := api.NewServer(store)

	log.Fatal(server.Router.Listen(":3000"))
}
