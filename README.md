# Hashdrop

Hashdrop is a work-in-progress secure file sharing and integrity verification tool with a Go backend and CLI client.

It allows users to authenticate, upload files securely, encrypt them client-side with unique per-file keys, verify integrity using hashing, and generate secure downloadable access links. Files are uploaded to AWS S3 using presigned URLs, and encryption keys are handled entirely client-side with an optional encrypted local vault.

> **Status:** Under active development  
> Core functionality is being built. Documentation, tests, and polishing are in progress.

---

## Features (Planned / In Progress)

- Go backend + CLI client  
- User authentication (access + refresh tokens)  
- Client-side encryption with per-file DEKs  
- Hash-based tamper verification  
- Secure upload to S3 via presigned URLs  
- Secure shareable download links  
- Local encrypted key vault (with optional manual key mode)

---

## Tech

- Go  
- AWS S3  
- Client-side encryption + hashing  

---

## Note
This project is still evolving. Not production-ready yet.