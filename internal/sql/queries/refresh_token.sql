-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    ?,
    datetime('now'),
    datetime('now'),
    ?,
    ?
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE token = ?;

-- name: GetUserFromRefreshToken :one
SELECT users.*
FROM users
INNER JOIN refresh_tokens
ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = ?
    AND refresh_tokens.revoked_at IS NULL
    AND datetime('now') < refresh_tokens.expires_at;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = datetime('now'), revoked_at = datetime('now')
WHERE token = ?;