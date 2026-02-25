-- +goose Up
CREATE TABLE download_attempts_count (
    id TEXT PRIMARY KEY,
    file_id TEXT NOT NULL,
    day DATE NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 1,
    UNIQUE (file_id, day)
);

-- +goose Down
DROP TABLE download_attempts_count;