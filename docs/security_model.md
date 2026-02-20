# Security Model

## Core philosophy

Hashdrop is built on a simple principle: **you trust no one with your files, not even the storage service.** Files are encrypted on your device before upload — the server and storage layer only ever see encrypted blobs. Your raw file data never touches the server.

---

## Encryption key generation

Every file gets its own unique Data Encryption Key (DEK), so sharing or compromising one file's key has no effect on any other file.

How the DEK is generated depends on the mode chosen at upload time:

**Vault mode (default)**
A random 32-byte DEK is generated automatically. You don't need to think about key management — Hashdrop handles it for you via the local vault (described below).

**No-vault mode**
You provide a passphrase and take responsibility for managing it. Hashdrop does not trust the passphrase alone for key strength, so a random 16-byte salt is generated and used with Argon2 to derive the final DEK. The salt is stored in the database so the key can be re-derived later for decryption.

> For a detailed walkthrough of the upload flow, see [Uploading Files](./uploading.md).

---

## File encryption

Before encrypting, Hashdrop computes a hash of the plaintext bytes for later integrity verification. The file is then encrypted in chunks using **AES-GCM**. Each chunk is written in the format:
```
[nonce][ciphertext length][ciphertext]
```

This structure makes the chunk boundaries unambiguous during decryption. AES-GCM also provides built-in integrity guarantees at the cipher level, and the plaintext hash provides an additional application-level check.

Only the encrypted blob is uploaded to S3.

---

## The local vault (vault mode)

After a successful upload in vault mode, the DEK is stored in a local encrypted vault at `~/.hashdrop/vault.enc`. The vault is a hashmap of file IDs to their base64-encoded DEKs.

The vault itself is protected by a passphrase you provide. That passphrase is processed using a random 16-byte salt and Argon2 to derive the vault encryption key, and the vault is encrypted with AES-GCM in the format:
```
[nonce][ciphertext]
```

The salt and other vault metadata are stored separately at `~/.hashdrop/vault_meta.json`. Even if `vault.enc` is stolen, it cannot be read without your vault passphrase.

Whenever you upload a new file or decrypt a file whose key lives in the vault, you are prompted for your vault passphrase to unlock it.

---

## Sharing and decryption

**If you own the file (vault mode):** The DEK is extracted from your vault and can be shared as a base64 string. The recipient uses it directly with `--key` to decrypt.

**If you own the file (no-vault mode):** You share the original passphrase.

---

## Summary

| What is stored on the server | What never leaves your device |
|---|---|
| Encrypted file blob (S3) | Plaintext file data |
| DEK salt (no-vault mode) | Raw DEK (vault mode) |
| File metadata and integrity hash | Vault passphrase |