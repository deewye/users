-- name: InsertUser :exec
INSERT INTO users (email, name, birthday)
VALUES ($1, $2, $3);

-- name: GetUserByID :one
SELECT id, email, name, birthday, created_at, updated_at
FROM users
WHERE id = $1 LIMIT 1;