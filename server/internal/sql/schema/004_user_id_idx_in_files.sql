-- +goose Up
CREATE INDEX idx_files_user_id_status_created_at
    ON files(user_id, status, created_at DESC);

-- +goose Down
DROP INDEX idx_files_user_id_status_created_at;