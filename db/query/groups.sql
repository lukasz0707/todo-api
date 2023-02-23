-- name: CreateGroup :one
INSERT INTO groups (
  group_name
) VALUES (
    $1
)RETURNING *;

-- name: GetGroupsByUserID :many
SELECT * FROM groups WHERE id IN (SELECT group_id FROM users_groups WHERE user_id = $1);