
# DOCUMENTATION AND DEPLOYMENT WORK IN PROGRESS...
# Hashdrop

> Hashdrop is a zero-trust secure command-line storage tool built in Go. It allows users to encrypt any kind of data (e.g. text, audio, video, image), on the client side and uploads the encrypted blob on a secure storage. It uses unique separate keys to encrypt each file, allowing stronger security. The encrypted files can be shared by links, and authorized people can decrypt and download the file safely using the shared key. It also allows the user to check file's integrity to ensure the data is not tampered with.

## How It Works

Hashdrop follows a simple but secure flow:

1. **Encrypt locally** — files are encrypted on your machine before anything leaves it
2. **Upload safely** — only the encrypted blob reaches cloud storage (AWS S3)
3. **Share a link** — recipients get a short-lived signed download URL via CloudFront
4. **Decrypt on their end** — the recipient uses the shared key or passphrase to decrypt and verify the file

**Your plaintext data never touches the server.**

---

## Features

- **Token-based user authentication** — JWT access tokens + revocable refresh tokens
- **Client-side AES-GCM encryption** — each file gets its own unique Data Encryption Key (DEK)
- **Integrity verification** — plaintext hash is stored so you can confirm data wasn't tampered with
- **Local vault** — store your keys in an AES-GCM encrypted vault at `~/.hashdrop/vault.enc`
- **Passphrase mode (Optional)** — use a passphrase instead of a vault if you prefer using your own key

---

## Architecture Overview

Hashdrop has two main components:

**Server** — a Go HTTPS API handling auth, file metadata, presigned S3 upload URLs, and CloudFront signed download URLs.

**CLI** — a Go Cobra-based client (`hashdrop`) that handles encryption, uploads, downloads, decryption, and local key/token management.

---

## Installation
```bash
git clone https://github.com/your-username/hashdrop
cd hashdrop
go build ./cli/...
```

Move the binary somewhere in your `$PATH`:
```bash
mv hashdrop /usr/local/bin/hashdrop
```

---

## Usage

### Register and verify your account
```bash
hashdrop auth register
```

### Log in
```bash
hashdrop auth login
```

### Upload a file
```bash
hashdrop upload ./secret.pdf
```

### List your files
```bash
hashdrop files list
```

### Share and download
```bash
# The upload command returns a shareable link
# The recipient runs:
hashdrop decrypt --key <shared-key> <download-url> -o output.pdf
```

### Delete a file
```bash
hashdrop files delete <file-id>
```

---

## Local Storage

Hashdrop keeps all local state under `~/.hashdrop/`:

| File | Contents |
|---|---|
| `tokens.json` | Access + refresh tokens
| `vault.enc` | AES-GCM encrypted key vault

---

## Tech Stack

- **Go** — server and CLI
- **SQLite** — metadata storage (WAL mode)
- **AWS S3** — encrypted object storage
- **AWS CloudFront** — signed URL delivery
- **AWS SES** — OTP email delivery
- **sqlc** — type-safe SQL query generation

---