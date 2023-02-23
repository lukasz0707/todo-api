-- name: CreateTodo :one
INSERT INTO todos (
  group_id,
  todo_name,
  deadline
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetTodosByGroupID :many
SELECT * FROM todos WHERE todos.group_id IN (SELECT users_groups.group_id FROM users_groups WHERE user_id = $1 and users_groups.group_id = $2);

-- name: GetTodosByUserID :many
SELECT * FROM todos WHERE group_id IN (SELECT group_id FROM users_groups WHERE user_id = $1);