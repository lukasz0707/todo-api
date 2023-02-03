-- name: CreateUsersTodos :one
INSERT INTO users_todos (
  user_id,
  todos_id,
  has_permissions
) VALUES (
    $1, $2, $3
)RETURNING *;