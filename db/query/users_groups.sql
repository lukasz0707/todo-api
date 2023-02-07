-- name: AssignUserToGroup :one
INSERT INTO users_groups (
  user_id,
  group_id
) VALUES (
    $1, $2
)RETURNING *;