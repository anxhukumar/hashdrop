-- +goose Up
CREATE INDEX idx_files_user_id_created_at
    ON files(user_id, created_at DESC);

-- +goose Down
DROP INDEX idx_files_user_id_created_at;