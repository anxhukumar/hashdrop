# API Reference

---

## Authentication

Endpoints marked as **authenticated** require the following header:
```
Authorization: Bearer <access_token>
```

Requests with a missing, malformed, or expired token will receive a `401 Unauthorized` response.

---

## Base URL

All API endpoints are relative to your Hashdrop server base URL, e.g. `https://api.hashdrop.dev`.

## Endpoints

---

### Health check
```
GET /api/healthz
```

Returns a plain text `OK` response indicating the API is running and ready to serve requests.

**Authentication:** None

**Request body:** None

**Response `200 OK`**
```
OK
```

Content-Type: `text/plain; charset=utf-8`

---

### Register a new user
```
POST /api/user/register
```

Creates a new user account. On success, a verification OTP is sent to the provided email address. The account must be verified before login is possible.

**Authentication:** None

**Request body**
```json
{
  "email": "user@example.com",
  "password": "yourpassword"
}
```

**Response `201 Created`**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com"
}
```

**Error responses**

| Status | Description |
|---|---|
| `400 Bad Request` | Malformed JSON or invalid password |
| `409 Conflict` | An account with this email already exists |
| `500 Internal Server Error` | Unexpected server error |

---

### Verify account
```
PATCH /api/user/verify
```

Verifies a newly registered account using the OTP sent to the user's email. Once verified, the account can be used to log in.

**Authentication:** None

**Request body**
```json
{
  "email": "user@example.com",
  "otp": "123456"
}
```

**Response `204 No Content`**

No response body.

**Error responses**

| Status | Description |
|---|---|
| `400 Bad Request` | Malformed JSON or OTP has expired |
| `401 Unauthorized` | Account not found or invalid OTP |
| `500 Internal Server Error` | Unexpected server error |

---

### Login
```
POST /api/user/login
```

Authenticates a verified user and returns an access token and a refresh token. The account must be verified before login is possible.

**Authentication:** None

**Request body**
```json
{
  "email": "user@example.com",
  "password": "yourpassword"
}
```

**Response `200 OK`**
```json
{
  "access_token": "eyJ...",
  "refresh_token": "a3f8..."
}
```

**Error responses**

| Status | Description |
|---|---|
| `400 Bad Request` | Malformed JSON |
| `401 Unauthorized` | Invalid email or password, or account not verified |
| `500 Internal Server Error` | Unexpected server error |

---

### Delete account
```
DELETE /api/user
```

Permanently deletes the authenticated user's account. This removes all uploaded files from storage, all associated database records, and invalidates all tokens. This action is irreversible.

**Authentication:** Required

**Request body**
```json
{
  "email": "user@example.com",
  "password": "yourpassword"
}
```

Email and password are required as explicit confirmation before deletion proceeds.

**Response `204 No Content`**

No response body.

**Error responses**

| Status | Description |
|---|---|
| `400 Bad Request` | Malformed JSON |
| `401 Unauthorized` | Invalid email or password |
| `500 Internal Server Error` | Unexpected server error |

---

### Refresh access token
```
POST /api/token/refresh
```

Exchanges a valid refresh token for a new access token. Use this when the current access token has expired.

**Authentication:** None

**Request body**
```json
{
  "refresh_token": "a3f8..."
}
```

**Response `200 OK`**
```json
{
  "access_token": "eyJ..."
}
```

**Error responses**

| Status | Description |
|---|---|
| `400 Bad Request` | Malformed JSON |
| `401 Unauthorized` | Refresh token is invalid, expired, or revoked |
| `500 Internal Server Error` | Unexpected server error |

---

### Revoke refresh token
```
POST /api/token/revoke
```

Revokes a refresh token, immediately invalidating it. Any subsequent attempts to use it to obtain a new access token will fail. This is called automatically when the user logs out.

**Authentication:** None

**Request body**
```json
{
  "refresh_token": "a3f8..."
}
```

**Response `204 No Content`**

No response body.

**Error responses**

| Status | Description |
|---|---|
| `400 Bad Request` | Malformed JSON |
| `500 Internal Server Error` | Unexpected server error |

---

### Request presigned upload URL
```
POST /api/files/presign
```

Validates storage quotas and returns a presigned S3 URL to upload an encrypted file directly to storage, along with a file ID to reference the upload in subsequent requests.

**Authentication:** Required

**Request body**
```json
{
  "file_name": "document.pdf",
  "mime_type": "application/pdf"
}
```

**Response `200 OK`**
```json
{
  "file_id": "uuid",
  "upload_resource": {
    "url": "https://s3.amazonaws.com/..."
  }
}
```

The `url` in `upload_resource` is a presigned S3 PUT URL. The encrypted file should be uploaded directly to this URL with the content type set to `application/octet-stream`. The URL has a limited expiry window — upload should begin immediately after receiving it.

**Error responses**

| Status | Description |
|---|---|
| `400 Bad Request` | Malformed JSON |
| `401 Unauthorized` | Missing or invalid access token |
| `403 Forbidden` | User storage quota exceeded |
| `503 Service Unavailable` | Global storage limit reached, uploads temporarily disabled |
| `500 Internal Server Error` | Unexpected server error |

---

### Complete file upload
```
POST /api/files/complete
```

Called after the encrypted file has been successfully uploaded to S3 via the presigned URL. The server verifies the uploaded file size, finalises the file record, and marks it as uploaded. If the file exceeds the allowed size limit, it is deleted from storage and the upload is marked as failed.

**Authentication:** Required

**Request body**
```json
{
  "file_id": "uuid",
  "plaintext_hash": "sha256...",
  "plaintext_size_bytes": 102400,
  "passphrase_salt": "hex-encoded-salt"
}
```

`passphrase_salt` should be provided only when the file was uploaded in no-vault mode. Leave it as an empty string or omit it for vault mode uploads.

**Response `200 OK`**
```json
{
  "s3_object_key": "usrh-abc123/uuid",
  "uploaded_file_size": 102400
}
```

`uploaded_file_size` reflects the server-verified encrypted file size in bytes as reported by S3, not the client-reported value.

**Error responses**

| Status | Description |
|---|---|
| `400 Bad Request` | Malformed JSON |
| `401 Unauthorized` | Missing or invalid access token |
| `404 Not Found` | File ID not found or does not belong to this user |
| `413 Content Too Large` | Uploaded file exceeds the allowed size limit |
| `500 Internal Server Error` | Unexpected server error |

---

### List all files
```
GET /api/files/all
```

Returns a list of all files belonging to the authenticated user.

**Authentication:** Required

**Request body:** None

**Response `200 OK`**
```json
[
  {
    "file_id": "uuid",
    "file_name": "document.pdf",
    "encrypted_size_bytes": 102400,
    "status": "uploaded",
    "key_management_mode": "vault",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

`key_management_mode` is either `vault` or `passphrase`. `status` reflects the current state of the file and can be `pending`, `uploaded`, `failed`, or `deleted`.

**Error responses**

| Status | Description |
|---|---|
| `401 Unauthorized` | Missing or invalid access token |
| `404 Not Found` | No files found for this user |
| `500 Internal Server Error` | Unexpected server error |

---

### Resolve file ID conflicts
```
GET /api/files/resolve?id={fileIDPrefix}
```

Returns all files whose IDs begin with the given prefix. Used when a shortened file ID matches more than one file, allowing the correct file to be identified by its full ID and name.

**Authentication:** Required

**Query parameters**

| Parameter | Required | Description |
|---|---|---|
| `id` | Yes | A file ID prefix to match against |

**Response `200 OK`**
```json
[
  {
    "file_name": "document.pdf",
    "file_id": "uuid"
  }
]
```

**Error responses**

| Status | Description |
|---|---|
| `400 Bad Request` | Missing `id` query parameter |
| `401 Unauthorized` | Missing or invalid access token |
| `404 Not Found` | No files found matching the given prefix |
| `500 Internal Server Error` | Unexpected server error |

---

### Get file details
```
GET /api/files?id={fileID}
```

Returns detailed metadata for files matching the given file ID. The `id` parameter is matched as a prefix, so a full or partial file ID can be provided. If multiple files match, all matches are returned.

**Authentication:** Required

**Query parameters**

| Parameter | Required | Description |
|---|---|---|
| `id` | Yes | Full or partial file ID to match against |

**Response `200 OK`**
```json
[
  {
    "file_id": "uuid",
    "file_name": "document.pdf",
    "status": "uploaded",
    "plaintext_size_bytes": 98304,
    "encrypted_size_bytes": 102400,
    "s3_key": "usrh-abc123/uuid",
    "key_management_mode": "vault",
    "plaintext_hash": "sha256..."
  }
]
```

**Error responses**

| Status | Description |
|---|---|
| `400 Bad Request` | Missing `id` query parameter |
| `401 Unauthorized` | Missing or invalid access token |
| `404 Not Found` | No files found matching the given ID |
| `500 Internal Server Error` | Unexpected server error |

---

### Get file hash
```
GET /api/files/hash?id={fileID}
```

Returns the plaintext hash of a file. Used to verify the integrity of a decrypted file against the hash computed at upload time.

**Authentication:** Required

**Query parameters**

| Parameter | Required | Description |
|---|---|---|
| `id` | Yes | The exact full file UUID |

**Response `200 OK`**
```json
{
  "hash": "sha256..."
}
```

**Error responses**

| Status | Description |
|---|---|
| `400 Bad Request` | Missing or invalid `id` query parameter |
| `401 Unauthorized` | Missing or invalid access token |
| `404 Not Found` | No file found matching the given ID |
| `500 Internal Server Error` | Unexpected server error |

---

### Get passphrase salt
```
GET /api/files/salt?id={fileID}
```

Returns the passphrase salt for a file that was uploaded in no-vault mode. This salt is needed to re-derive the DEK from the original passphrase during decryption. Only available for files encrypted with a passphrase — vault mode files do not have a salt stored server-side.

**Authentication:** Required

**Query parameters**

| Parameter | Required | Description |
|---|---|---|
| `id` | Yes | The exact full file UUID |

**Response `200 OK`**
```json
{
  "salt": "hex-encoded-salt"
}
```

**Error responses**

| Status | Description |
|---|---|
| `400 Bad Request` | Missing or invalid `id` query parameter |
| `401 Unauthorized` | Missing or invalid access token |
| `404 Not Found` | No passphrase salt found for this file |
| `500 Internal Server Error` | Unexpected server error |

---

### Delete file
```
DELETE /api/files?id={fileID}
```

Permanently deletes a file from both storage and the database. The `id` parameter is matched as a prefix. If the prefix matches more than one file, the request is rejected to prevent accidental deletion — use a longer prefix or the full file UUID to disambiguate.

**Authentication:** Required

**Query parameters**

| Parameter | Required | Description |
|---|---|---|
| `id` | Yes | Full or partial file ID to match against |

**Response `204 No Content`**

No response body.

**Error responses**

| Status | Description |
|---|---|
| `400 Bad Request` | Missing `id` query parameter |
| `401 Unauthorized` | Missing or invalid access token |
| `404 Not Found` | No file found matching the given ID |
| `409 Conflict` | The given ID prefix matches multiple files |
| `500 Internal Server Error` | Unexpected server error |

---

### Download file
```
GET /files/download/{userIDHash}/{fileID}
```

Validates the daily download limit for the file and redirects to a short-lived signed CloudFront URL serving the encrypted file. This is a public endpoint — no authentication is required. The file URL is typically obtained from the `files show` command or the file details endpoint.

**Authentication:** None

**Path parameters**

| Parameter | Description |
|---|---|
| `userIDHash` | Hashed user ID of the file owner |
| `fileID` | Full file UUID |

**Response `303 See Other`**

Redirects directly to a signed CloudFront URL. The response body contains the encrypted file content.

**Error responses**

| Status | Description |
|---|---|
| `400 Bad Request` | Invalid file ID format |
| `429 Too Many Requests` | Daily download limit exhausted for this file |
| `500 Internal Server Error` | Unexpected server error |

---

### Reset (development only)
```
DELETE /admin/reset
```

Deletes all users and associated data from the database. This endpoint is only available when the server is running in a `dev` environment and will reject requests in any other environment.

**Authentication:** None

**Request body:** None

**Response `204 No Content`**

No response body.

**Error responses**

| Status | Description |
|---|---|
| `403 Forbidden` | Server is not running in a `dev` environment |
| `500 Internal Server Error` | Unexpected server error |