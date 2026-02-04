-- name: CheckAndUpdateDownloadAttemptsCount :one
INSERT INTO download_attempts_count(id, file_id, day, attempts)
VALUES (
    ?,
    ?,
    date('now'),
    1
)
ON CONFLICT(file_id, day)
DO UPDATE SET attempts = attempts + 1
RETURNING attempts;