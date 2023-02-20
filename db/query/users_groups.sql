-- name: AssignUserToGroup :one
INSERT INTO users_groups (
  user_id,
  group_id
) VALUES (
    $1, $2
)RETURNING *;
-- name: AssignOwnerToGroup :one
INSERT INTO users_groups (
  user_id,
  group_id,
  role
) VALUES (
    $1, $2, 'owner'
)RETURNING *;

-- name: SelectFromUsersGroups :one
SELECT * FROM users_groups WHERE user_id = $1 AND group_id = $2;