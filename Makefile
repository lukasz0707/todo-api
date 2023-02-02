postgres:
	docker run --name postgres15todo -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it postgres15todo  createdb --username=root --owner=root todo

dropdb:
	docker exec -it postgres15todo dropdb todo

migratecreate:
	migrate create -ext sql -dir db/migration/ -seq init_schema

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/todo?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/todo?sslmode=disable" -verbose down

sqlc:
	sqlc --experimental generate

.PHONY: postgres createdb dropdb migratecreate migrateup migratedown sqlc