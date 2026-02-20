# Authentication

## Registration

When you register, you provide your email and password. A 6-digit one-time code is sent to your email and must be verified before your account becomes active. Hashdrop stores only a hashed version of your password — the plaintext is never persisted.

---

## Login and tokens

When you log in, Hashdrop issues two tokens:

- **Access token** — a JWT that expires every **15 minutes**
- **Refresh token** — a random 32-byte hex string that expires after **30 days**

Both are stored locally at `~/.hashdrop/tokens.json` and the CLI reads them from there for every request. The refresh token is hashed before being stored in the database, so even if the database were leaked, the raw token would not be exposed.

---

## How token refresh works

The short-lived access token limits the damage if it is ever compromised — an attacker can only use it for up to 15 minutes. To avoid asking you to re-authenticate every 15 minutes, the refresh token is used to obtain a new access token whenever the current one has expired.

If the refresh token itself is compromised, it can be revoked by logging out, which immediately invalidates it on the server.

You will need to re-authenticate every **30 days** when the refresh token expires.