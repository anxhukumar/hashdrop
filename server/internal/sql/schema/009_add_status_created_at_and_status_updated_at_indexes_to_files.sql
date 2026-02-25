-- +goose Up
CREATE INDEX idx_files_status_created_at
    ON files(status, created_at);
CREATE INDEX idx_files_status_updated_at
    ON files(status, updated_at);

-- +goose Down
DROP INDEX idx_files_status_created_at;
DROP INDEX idx_files_status_updated_at;