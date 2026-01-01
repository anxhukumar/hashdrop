-- +goose Up
CREATE TABLE files (
    id TEXT PRIMARY KEY,

    user_id TEXT NOT NULL,
    file_name TEXT NOT NULL,
    mime_type TEXT,
    
    -- Integrity / verification
    plaintext_hash TEXT,
    plaintext_hash_version INTEGER NOT NULL DEFAULT 1,
    
    -- Data size stored in s3
    plaintext_size_bytes INTEGER,
    encrypted_size_bytes INTEGER,
    
    -- Ecryption metadata
    encryption_version INTEGER NOT NULL DEFAULT 1,
    encryption_chunk_kb INTEGER NOT NULL DEFAULT 64,

    -- vault metadata
    key_management_mode TEXT NOT NULL
        CHECK (key_management_mode IN ('vault','passphrase')),
    passphrase_salt TEXT,

    -- Key derivation 
    kdf_version INTEGER NOT NULL DEFAULT 1,

    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'uploaded', 'deleted', 'failed')),

    s3_key TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE files;