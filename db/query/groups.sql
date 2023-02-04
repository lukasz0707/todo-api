-- name: Creategroup :one
INSERT INTO groups (
  group_name,
  owner_id
) VALUES (
    $1, $2
)RETURNING *;