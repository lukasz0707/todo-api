-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
    $1, $2, $3, $4
) RETURNING id, username, full_name, email, password_changed_at, created_at;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;