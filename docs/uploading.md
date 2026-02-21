# File Upload

## Overview

When you run `hashdrop upload <file-path>`, a multi-step process coordinates between your device, the Hashdrop backend, and S3 to ensure your file is validated, encrypted, and securely stored — without your raw file data ever reaching the server. For details on how encryption and key generation work, see [Security Model](./security_model.md).

---

## 1. Client-side validation

Before anything is sent over the network, the CLI checks the file size against a configured per-file size limit. This is an early rejection for well-behaved clients — since this check runs on the client it can be bypassed, but the server performs its own authoritative size verification later in the flow.

---

## 2. Preparing file metadata

Concurrently, Hashdrop computes a hash of the plaintext file bytes and detects the MIME type of the file. These happen in parallel to keep the pre-upload phase fast.

The user is also given the option to provide a custom display name for the file. If none is provided, the original filename is kept.

---

## 3. Obtaining the encryption key

The CLI obtains the Data Encryption Key (DEK) for the file based on the mode selected — vault or no-vault. See [Security Model](./security_model.md) for how keys are generated in each mode.

---

## 4. Requesting a presigned upload URL

The CLI makes its first request to the backend, sending the filename and MIME type. The backend then:

- Checks whether the global S3 storage limit has been reached
- Checks whether the requesting user has exceeded their personal storage quota by summing the size of all their uploaded files
- Hashes the user ID to use as a path component in the S3 key, avoiding direct exposure of the user's identity in storage paths
- Generates a unique file ID
- Constructs the S3 key in the format `usrh-{userIdHash}/{fileId}`
- Generates a presigned S3 PUT URL with an expiry window
- Creates a `pending` file entry in the database

The presigned URL and file ID are returned to the CLI.

---

## 5. Streaming encrypted upload

The CLI encrypts the file in chunks inside a goroutine and simultaneously streams the encrypted bytes to S3 via the presigned URL using an `io.Pipe()`. The content type is set to `application/octet-stream`.

The pipe provides natural backpressure — if the S3 upload isn't keeping up with encryption, the writer blocks automatically. This means the file is never fully loaded into memory; everything flows on the fly with a bounded memory footprint.

If an error occurs during the upload, the file entry remains in `pending` status in the database.

---

## 6. Server-side confirmation and verification

Once the stream completes, the CLI makes a second request to the backend to confirm the upload. It sends:

- The plaintext file hash
- The client-reported plaintext file size
- The passphrase salt, if no-vault mode was used

The backend then requests the actual object size directly from S3 via a `HeadObject` call — this value cannot be manipulated by the client. This size is validated against the per-file limit. If it exceeds the limit, the object is deleted from S3 immediately and the file is marked as `failed` in the database. If validation passes, the file status is updated to `uploaded` and all remaining metadata is persisted.

---

## 7. Vault update

Finally, the CLI updates the local vault with the file ID and its DEK, and returns a success message to the user.

---

## Pending file cleanup

Files that remain in `pending` status beyond a certain duration are assumed to be either stale or the result of a failed or aborted upload and are automatically removed from S3. For details on how this cleanup and other guardrails are implemented, see [Resource Limits and Abuse Prevention](./resource-limits-and-abuse-prevention.md).