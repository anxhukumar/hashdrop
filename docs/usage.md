> Before proceeding, make sure you have the Hashdrop CLI installed. See the [Installation guide](./installation.md) to get set up.
# CLI Usage

This guide walks you through the core Hashdrop CLI workflows — from creating an account to uploading, sharing, and decrypting files.

---

## Authentication

**Register an account**
```bash
hashdrop auth register
```

You'll be prompted for your email, password, and a one-time code sent to your email to verify your account.

**Log in**
```bash
hashdrop auth login
```

Your tokens are stored locally so subsequent commands work without re-authenticating.

**Log out**
```bash
hashdrop auth logout
```

**Delete your account**
```bash
hashdrop auth delete-account
```

Permanently removes your account, all uploaded files, tokens, and shared links. This is irreversible.

---

## Uploading Files
```bash
hashdrop upload <file-path>
```

Files are encrypted client-side and uploaded via a secure presigned URL. By default, the encryption key is stored in your local vault.

**Common options**

| Flag | Description |
|---|---|
| `-n, --name` | Set a custom display name for the file |
| `-N, --no-vault` | Use a self-managed passphrase instead of the vault. If lost, the file cannot be decrypted. |

---

## Managing Files

**List all uploaded files**
```bash
hashdrop files list
```

Shows file names, sizes, upload status, key mode, and creation dates.

**Inspect a file**
```bash
hashdrop files show <file-id>
```

Displays detailed metadata including the download URL. Use `--reveal-key` to show the vault-stored encryption key.

**Delete a file**
```bash
hashdrop files delete <file-id>
```

Permanently removes the file from storage. It can no longer be downloaded or decrypted.

---

## Decrypting Files
```bash
hashdrop decrypt <file-url> [destination]
```

Downloads and decrypts a Hashdrop file locally. You'll be prompted for the decryption secret based on how the file was encrypted.

**Examples**
```bash
# Decrypt to default Downloads directory
hashdrop decrypt https://hashdrop.dev/...

# Decrypt to a specific path
hashdrop decrypt https://hashdrop.dev/... ~/Documents

# Decrypt with integrity verification
hashdrop decrypt --verify https://hashdrop.dev/...
```

**Decryption modes:** `--vault` (default for your own files), `--key` (raw key), or passphrase prompt for `--no-vault` uploads.

---

> Add `-v` / `--verbose` to any command for detailed output.