-- +goose Up
CREATE TABLE files (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    file_name TEXT NOT NULL,
    mime_type TEXT,
    hash TEXT,
    size_bytes INTEGER,
    hash_alg TEXT,
    enc_alg TEXT,
    kdf TEXT,
    kdf_params TEXT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    s3_key TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE files;