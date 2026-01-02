# Hashdrop

Hashdrop is a secure file sharing and integrity verification tool consisting of a Go backend and a Go CLI client.

Files are **encrypted locally** before upload using unique per-file Data Encryption Keys (DEKs), ensuring the server and storage provider never see plaintext. Uploaded data is stored in AWS S3 and delivered through AWS CloudFront. Users can verify file integrity via cryptographic hashing and manage encryption keys locally using an optional encrypted vault.

> **Status:** Work in progress  
> Core features are being actively built. APIs, CLI UX, testing, docs, and cleanup are still evolving.

---

## ✨ Features (Planned / In Progress)

- ⚙️ Go backend + Go CLI client
- 🔐 User authentication (access + refresh tokens)
- 🛡️ Client-side encryption  
  - unique per-file DEKs  
  - AES-GCM streaming encryption
- 🧾 Integrity verification via hashing
- ☁️ Secure uploads using AWS S3 presigned URLs
- 🌍 Secure file delivery via AWS CloudFront
- 🔑 Local encrypted vault for key storage  
  - optional **no-vault mode** (self-managed passphrase)

---

## 🧰 Tech Stack

- **Go**
- **AWS S3** (encrypted blob storage)
- **AWS CloudFront** (secure download delivery)
- **Client-side cryptography**
  - AES-GCM
  - Argon2 key derivation
  - SHA-256 hashing

---

## ⚠️ Note

This project is still evolving and is **not production-ready yet**.  
APIs may change, features may break, and security-critical components are still undergoing refinement.

---

More updates coming soon 🚀