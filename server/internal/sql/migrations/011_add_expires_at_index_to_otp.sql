-- +goose Up
CREATE INDEX idx_otp_expires_at
    ON otp(expires_at);

-- +goose Down
DROP INDEX idx_otp_expires_at;