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

-- name: GetS3KeyForUploadVerification :one
SELECT s3_key
FROM files
WHERE id = ? AND user_id = ? AND status='pending';

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
    status = 'failed',
    updated_at = datetime('now')
WHERE id = ? AND user_id = ? AND status='pending';

-- name: GetAllFilesOfUser :many
SELECT file_name, encrypted_size_bytes, status, key_management_mode, created_at, id
FROM files
WHERE user_id = ? AND status = 'uploaded'
ORDER BY created_at DESC;

-- name: GetDetailedFile :many
SELECT file_name, id, status, plaintext_size_bytes, encrypted_size_bytes, s3_key, key_management_mode, plaintext_hash
FROM files
WHERE user_id = ? AND status = 'uploaded' AND id LIKE CAST(? AS TEXT);

-- name: GetPassphraseSalt :one
SELECT passphrase_salt
FROM files
WHERE user_id = ? AND status = 'uploaded' AND key_management_mode = 'passphrase' AND id = ?;

-- name: GetFileHash :one
SELECT plaintext_hash
FROM files
WHERE user_id = ? AND status = 'uploaded' AND id = ?;

-- name: DeleteFileFromId :exec
UPDATE files
SET
    status = 'deleted',
    updated_at = datetime('now')
WHERE user_id = ? AND status = 'uploaded' AND id LIKE CAST(? AS TEXT);

-- name: GetS3KeyFromFileID :many
SELECT s3_key
FROM files
WHERE user_id = ? AND status='uploaded' AND id LIKE CAST(? AS TEXT);

-- name: CheckShortFileIDConflict :many
SELECT file_name, id
FROM files
WHERE user_id = ? AND status = 'uploaded' AND id LIKE CAST(? AS TEXT);

-- name: GetAnyS3KeyOfUser :one
SELECT s3_key
FROM files
WHERE user_id = ? AND status = 'uploaded'
LIMIT 1;

-- name: GetTotalBytesOfUploadedFiles :one
SELECT CAST(COALESCE(SUM(encrypted_size_bytes), 0) AS INTEGER) AS total_bytes
FROM files
WHERE status = 'uploaded';

-- name: GetUsersTotalBytesOfUploadedFiles :one
SELECT CAST(COALESCE(SUM(encrypted_size_bytes), 0) AS INTEGER) AS total_bytes
FROM files
WHERE user_id = ? AND status = 'uploaded';

-- name: GetStalePendingFiles :many
SELECT id, user_id, s3_key
FROM files
WHERE status = 'pending'
    AND created_at < :cutoff_time;

-- name: CleanDeletedAndFailedFiles :exec
DELETE
FROM files
WHERE (status = 'deleted' OR status = 'failed')
    AND updated_at < :cutoff_time;