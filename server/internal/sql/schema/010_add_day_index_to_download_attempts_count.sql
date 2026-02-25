-- +goose Up
CREATE INDEX idx_download_attempts_count_day
    ON download_attempts_count(day);

-- +goose Down
DROP INDEX idx_download_attempts_count_day;