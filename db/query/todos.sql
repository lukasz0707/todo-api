-- name: CreateTodo :one
INSERT INTO todos (
  group_id,
  todo_name,
  deadline
) VALUES (
    $1, $2, $3
) RETURNING *;