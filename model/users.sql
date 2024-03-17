-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUsers :many
SELECT * FROM users ORDER BY id;

-- name: IsEmailTaken :one
SELECT 1 FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (
  first_name,
  last_name,
  email,
  password,
  verified,
  verification_token,
  avatar
) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;
