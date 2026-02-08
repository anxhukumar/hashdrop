-- +goose Up
CREATE TABLE otp (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    otp_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES users(id)
);

-- +goose DOWN
DROP TABLE otp