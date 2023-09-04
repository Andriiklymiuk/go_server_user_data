-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users
LIMIT 10;

-- name: CreateUser :one
INSERT INTO users (name) 
VALUES ($1) 
RETURNING *;

-- name: UpdateUser :one
UPDATE users 
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1;