# Testing

Hashdrop uses Go's standard testing toolchain. Tests are organized alongside the code they cover.

---

## Running the tests

### Server
```bash
cd server
go test ./...
```

### CLI
```bash
cd cli
go test ./...
```

---

## Server — current coverage

### `internal/auth`

Covers the core authentication primitives:

- Password hashing and comparison — validates that empty passwords, passwords below the minimum length, and valid passwords are handled correctly, and that a hashed password correctly round-trips through comparison.
- JWT generation and validation — tests valid tokens, malformed token strings, and tokens verified against the wrong secret.
- Bearer token extraction — tests valid headers, missing headers, missing Bearer prefix, empty tokens, and tokens with extra whitespace.
- Refresh token generation — validates that the token is a 64-character hex string and that consecutive calls produce unique tokens.

### `internal/aws`

Covers presigned S3 PUT URL generation:

- Validates that a nil S3 client returns an error.
- Validates that a properly configured client returns a non-empty presigned URL.

### `internal/cloudfront_guard`

Covers CloudFront signing and download attempt validation:

- Private key loading — tests a valid PEM key, a missing file path, and an invalid PEM file.
- Signed URL generation — tests valid inputs producing a non-empty URL and an invalid key path returning an error.
- Download attempt validation — tests that a request within the daily limit is allowed and a request over the limit is rejected, using a real test database.

### `internal/otp`

Covers OTP generation and verification:

- OTP generation — validates that the generated code is exactly 6 digits and contains only numeric characters.
- OTP hashing and verification — tests correct OTP and secret, wrong OTP, wrong secret, and a tampered hash. Also validates that `HashOTP` is deterministic.

### `internal/storage_guard`

Covers S3 storage quota enforcement:

- Global quota — simulates a 4 GB upload and tests against limits both above and below that threshold.
- Per-user quota — simulates separate uploads for two users and validates that each user's quota is checked independently.

---

## Server — planned coverage

### `internal/handlers`

Handler integration tests are planned using Go's `httptest` package against a real test database. The following handlers are in scope:

`HandlerCreateUser`, `HandlerVerifyUser`, `HandlerLogin`, `HandlerRefreshToken`, `HandlerRevokeToken`, `HandlerDeleteUser`, `HandlerGeneratePresignLink`, `HandlerCompleteFileUpload`, `HandlerGetAllFiles`, `HandlerResolveFileMatches`, `HandlerGetDetailedFile`, `HandlerGetPassphraseSalt`, `HandlerGetFileHash`, `HandlerDeleteFile`, `HandlerGenerateDownloadLink`.

`HandlerReadiness` and `HandlerReset` are intentionally excluded — the former is trivial and the latter is dev-only.

### `internal/cleaners`

Tests are planned for all background cleanup routines — stale pending S3 files, stale file metadata, expired refresh tokens, stale download counters, unverified users, and expired OTPs. These will use the same test database infrastructure used by the storage guard and CloudFront guard tests.

---

## CLI — planned coverage

### `internal/encryption`

Covers the client-side encryption pipeline — DEK generation, AES-GCM file encryption and decryption, vault creation, vault encryption and decryption, and vault key derivation.

### `internal/auth`

Covers authentication helpers — token storage and retrieval, access token refresh logic, and login and registration flows against a mock server.

### `internal/upload`

Covers the upload pipeline helpers — plaintext hash generation, MIME type detection, and client-side file size validation.

### `internal/decrypt_command`

Covers decryption helpers — hash verification against the stored plaintext hash, vault key lookup, and passphrase-based key derivation.

### `internal/files`

Covers file ID resolution logic — short ID prefix matching and conflict detection.

---

## Test infrastructure

The server test suite uses a shared `SetupTestDB` helper in `internal/test_util` that creates a temporary SQLite database and runs all schema migrations in order before each test. This ensures tests run against a real schema without requiring any external dependencies or a running server instance.
