-- name: CheckAndUpdateDownloadAttemptsCount :one
INSERT INTO download_attempts_count(id, file_id, created_at, attempts)
VALUES (
    ?,
    ?,
    datetime('now'),
    1
)
ON CONFLICT(file_id, created_at)
DO UPDATE SET attempts = attempts + 1
RETURNING attempts;

-- name: CleanDownloadCount :exec
DELETE
FROM download_attempts_count
WHERE created_at < :cutoff_time;