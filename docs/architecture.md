# Architecture

## Overview

Hashdrop is a client-server application. The backend is a Go API server hosted on a managed server instance, responsible for authentication, metadata management, and coordinating with AWS services. The CLI is the client - it handles all file encryption locally before any data leaves your device. The two communicate over HTTPS. For file transfers, the CLI uploads encrypted data directly to S3 and downloads via CloudFront.

![Hashdrop system architecture diagram](./assets/hashdrop-architecture.png)
---

## Components

### CLI (Client)

The CLI runs entirely on your machine. It is responsible for encrypting files before upload and decrypting them after download - your plaintext data never leaves your device. Local state is persisted in `~/.hashdrop/`:

- `tokens.json` - stores your access and refresh tokens (permissions: `0600`)
- `vault.enc` - an AES-GCM encrypted store of your file encryption keys (vault mode)
- `vault_meta.json` - vault metadata such as the salt used to derive the vault encryption key

The `~/.hashdrop` directory is created automatically on first use and is the only place Hashdrop writes persistent state to your machine. All files are written with `0600` permissions. The vault is encrypted at rest - even if `vault.enc` is stolen, it cannot be read without your vault passphrase.

---

### API Server

The backend is a Go HTTP server hosted on AWS EC2, running behind a reverse proxy that terminates HTTPS. It handles:

- User registration, authentication, and token lifecycle
- Storage quota enforcement and presigned URL generation
- File metadata management and lifecycle tracking (`pending → uploaded / failed / deleted`)
- Coordinating with AWS SES for OTP delivery, S3 for object management, and CloudFront for signed download URLs
- Protecting the service from abuse through request validation, per-user storage quotas, and per-file daily download limits

The server runs with an embedded database co-located on the same instance. This keeps the deployment simple - no separate database server is required. The database stores user records, file metadata, OTP state, token records, and download counters. If demand grows, the data layer can be migrated to a hosted database service without changes to the application logic.

---

### AWS Services

Hashdrop uses four AWS services:

**EC2** - hosts the Hashdrop API server. The server runs behind a reverse proxy on the EC2 instance, which handles HTTPS termination.

**S3** - stores the encrypted file blobs. Files are uploaded directly from the CLI to S3 via presigned PUT URLs. The server never proxies file data - it only coordinates the transfer.

**CloudFront** - serves encrypted files on download. The server generates a short-lived signed CloudFront URL and redirects the client to it. This keeps the download path fast and avoids routing large file transfers through the API server.

**SES** - sends OTP emails for account verification during registration.

---

## Request Flows at a Glance

**Registration and verification** - the CLI sends credentials to the server, which creates the account, generates an OTP, and delivers it via SES. The OTP hash is stored in the database and validated on verification.

**Login** - the server issues a short-lived JWT access token and a longer-lived refresh token. Both are stored locally by the CLI.

**Upload** - the CLI encrypts the file locally, requests a presigned S3 URL from the server, streams the encrypted blob directly to S3, then notifies the server to confirm and finalize the file record.

**Download** - the CLI requests the file from the server, which validates the daily download limit, generates a signed CloudFront URL, and redirects the CLI to it. The CLI then streams and decrypts the blob locally.

---

## Related Documentation

- [Security Model](./security_model.md) - encryption, key management, and the local vault
- [File Upload](./uploading.md) - detailed walkthrough of the upload flow
- [Downloading and Decryption](./decryption-and-downloading.md) - detailed walkthrough of the download and decryption flow
- [Authentication](./authentication.md) - token lifecycle and session management
- [CLI Usage](./usage.md) - command reference