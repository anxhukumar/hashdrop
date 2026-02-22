# Downloading and Decryption

## Overview

The decrypt command handles the full flow of downloading and decrypting a Hashdrop file locally. The file URL passed to the command is a Hashdrop backend endpoint that coordinates access control and delivery before your device ever receives any data. For details on how encryption and key generation work, see [Security Model](./security_model.md).

---

## 1. Specifying the decryption secret

After running:
```bash
hashdrop decrypt <file-url> [destination]
```

The CLI asks how you want to provide the decryption secret. You can either follow the prompts or use flags directly. The options are:

- **`--vault`** — the DEK is extracted automatically from your local vault (for files you uploaded in vault mode)
- **`--key`** — you provide the base64-encoded DEK directly (for recipients who were given the key by the file owner)
- **Passphrase** — you enter the original passphrase; the salt is fetched from the database to re-derive the correct DEK

Users who don't own the file will typically use `--key` or the passphrase option.

---

## 2. Download attempt validation

The CLI hits the Hashdrop backend endpoint `GET /files/download/{userIDHash}/{fileID}`. Before anything is served, the backend checks a per-file daily download counter stored in the database. If the file has already been downloaded more than the allowed number of times today, the request is rejected with a message indicating the daily limit has been exhausted. The counter resets after 24 hours.

Only if the count is within the limit does the backend proceed.

---

## 3. Signed CloudFront URL and delivery

Once the download is permitted, the backend generates a short-lived signed CloudFront URL pointing to the encrypted object in S3. The backend then redirects the request directly to that URL, so the encrypted content is streamed seamlessly — it appears to come from the Hashdrop backend even though it is served from the CDN.

---

## 4. Output path

By default, the decrypted file is saved to your system's Downloads directory. You can optionally provide a destination path or directory as the second argument to the command to control where the file is written.

---

## 5. Streaming decryption

The CLI opens an `io.Reader` from the CloudFront URL connection and decrypts the file in chunks using AES-GCM, mirroring exactly how it was encrypted during upload. The plaintext hash is also computed incrementally across chunks during this process.

If you pass the `--verify` flag, the computed hash is compared against the original plaintext hash stored in the database at upload time. A mismatch indicates the file's integrity cannot be confirmed.