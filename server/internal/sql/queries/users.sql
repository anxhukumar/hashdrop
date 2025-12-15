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

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?;
