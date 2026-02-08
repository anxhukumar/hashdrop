-- name: CreateNewUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, verified)
VALUES (
    ?,
    datetime('now'),
    datetime('now'),
    ?,
    ?,
    0
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetVerifiedUserByEmail :one
SELECT *
FROM users
WHERE email = ? AND verified = 1;

-- name: GetUnverifiedUserByEmail :one
SELECT *
FROM users
WHERE email = ? AND verified = 0;

-- name: DeleteUserById :exec
DELETE
FROM users
WHERE id = ?;

-- name: MarkUserVerifiedByEmail :exec
UPDATE users
SET verified = 1,
    updated_at = datetime('now')
WHERE email = ? AND verified = 0;