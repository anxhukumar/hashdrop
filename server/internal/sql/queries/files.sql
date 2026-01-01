-- name: CreatePendingFile :exec
INSERT INTO files (id, user_id, file_name, mime_type, s3_key, created_at, updated_at)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    datetime('now'),
    datetime('now')
);

-- name: GetS3KeyFromFileID :one
SELECT s3_key
FROM files
WHERE id = ? AND user_id = ?;

-- name: UpdateUploadedFile :exec
UPDATE files
SET 
    plaintext_hash = ?,
    plaintext_size_bytes = ?,
    encrypted_size_bytes = ?,
    key_management_mode = ?,
    passphrase_salt = ?,
    status = ?,
    updated_at = datetime('now')
WHERE id = ? AND user_id = ? AND status='pending';

-- name: UpdateFailedFile :exec
UPDATE files
SET
    status = ?,
    updated_at = datetime('now')
WHERE id = ? AND user_id = ? AND status='pending';