-- name: CreatePendingFile :one
INSERT INTO files (id, user_id, file_name, mime_type, s3_key, created_at, updated_at)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
)
RETURNING *;