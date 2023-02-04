-- name: CreateTodo :one
INSERT INTO todos (
  group_id,
  todo_name,
  status,
  deadline
) VALUES (
    $1, $2, $3, $4
) RETURNING *;