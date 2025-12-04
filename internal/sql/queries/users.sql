-- name: CreateNewUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    ?,
    datetime('now'),
    datetime('now'),
    ?,
    ?
)
RETURNING *;