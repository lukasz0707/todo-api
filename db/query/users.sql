-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  email
) VALUES (
    $1, $2, $3
) RETURNING id, username, email, password_changed_at, created_at;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;