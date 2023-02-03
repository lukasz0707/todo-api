-- name: CreateTodo :one
INSERT INTO todos (
  todo_name,
  group_name,
  deadline
) VALUES (
    $1, $2, $3
) RETURNING *;