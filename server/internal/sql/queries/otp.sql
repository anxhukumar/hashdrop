-- name: CreateOtp :exec
INSERT INTO otp (id, user_id, otp_hash, created_at, expires_at)
VALUES (
    ?,
    ?,
    ?,
    datetime('now'),
    datetime('now', '+10 minutes')
);

-- name: GetOtpHash :one
SELECT otp_hash
FROM otp
WHERE user_id = ?
    AND expires_at > datetime('now')
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteOtpByUserID :exec
DELETE FROM otp
WHERE user_id = ?;

-- name: DeleteOtpByOtpID :exec
DELETE FROM otp
WHERE id = ?;

-- name: CleanExpiredOtp :exec
DELETE
FROM otp
WHERE expires_at < :cutoff_time;