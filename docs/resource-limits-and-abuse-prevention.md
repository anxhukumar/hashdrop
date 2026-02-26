# Resource Limits and Abuse Prevention

Hashdrop applies multiple layers of protection server-side to keep the service stable and prevent abuse. This document covers each layer and how they work together.

---

## Rate Limiting

Every endpoint in the API is rate limited. The rate limiter is applied as middleware before a request reaches any handler logic.

### Two types of limiters

**Global limiter** — a single shared limiter for an endpoint group, applied across all incoming requests regardless of who is making them. If the global limit is exhausted, all further requests to that group are rejected until it recovers. This protects the server from being overwhelmed at a service level.

**Key-based limiter** — a per-identity limiter that tracks each caller separately. Depending on the endpoint group, the key is either the caller's IP address or their authenticated user ID. This prevents any single identity from monopolising capacity even when the global limit has not been reached.

### How they are applied

Endpoints are grouped by their nature and cost, and each group gets its own pair of limiters sized appropriately for that group's expected traffic and resource cost.

Auth endpoints (registration, login, account deletion) are protected by both a global limiter and a per-IP limiter. Since these endpoints are public and unauthenticated, IP is the only available identity signal. Token endpoints (refresh, revoke) follow the same pattern but with more generous limits since they are lightweight and frequent. Upload endpoints (presign, complete) use a global limiter and a per-user limiter — uploads are the most resource-intensive operations involving S3 and database writes, so these limits are the most conservative. List and file metadata endpoints also use a global and per-user limiter pair, sized more permissively since they are read-heavy and less expensive. The health check endpoint has a global limiter only, to prevent the monitoring endpoint itself from being used as a DDoS vector.

For authenticated endpoint groups, the key-based check is applied first. If it passes, the request proceeds to the global check. A request must pass both to reach the handler.

### Key-based limiter memory management

The per-key limiters maintain an in-memory map of active identities. To prevent unbounded memory growth, a background goroutine periodically sweeps the map and removes entries that have been inactive beyond a set TTL. This runs for the lifetime of the server and exits cleanly on shutdown.

---

## Storage Quota Enforcement

Before a presigned S3 upload URL is issued, the server runs two storage quota checks. If either check fails, the presigned URL is never generated and the upload is rejected before anything reaches S3.

### Global quota

The server sums the total bytes of all uploaded files across all users and compares it against a configured global limit. This is a last-resort ceiling — a safeguard for the worst case where collective usage approaches infrastructure limits. If the global limit is exceeded, uploads are temporarily disabled for all users until storage is freed.

### Per-user quota

The server sums the total bytes of all uploaded files belonging to the requesting user and compares it against a per-user limit. If the user has exceeded their personal storage allowance, the request is rejected.

In addition to the byte limit, the server also enforces a cap on the total number of files a user can have. This prevents a specific abuse pattern where a user uploads a large number of tiny files, each within the byte quota, but collectively creating an excessive number of database records and S3 objects.

### How the checks are ordered

The global check runs first. If it passes, the per-user check runs second. Both must pass before the server proceeds to generate the presigned URL and create the pending file record in the database.

---

## CloudFront Download Protection

The download endpoint is public — no authentication is required to request a file download. To prevent abuse, the server enforces a per-file daily download limit before a signed CloudFront URL is ever generated.

### Daily download limit

When a download is requested, the server looks up the download counter for that file in the database and increments it. If the count for the current day exceeds the configured limit, the request is rejected and no URL is generated. The counter resets daily — stale counters are cleared by the cleanup process covered in the next section.

### Signed CloudFront URL

Once the download limit check passes, the server generates a short-lived signed CloudFront URL using a private RSA key stored on the server. The URL is valid for a limited window of time — long enough to complete the download, short enough that a leaked or shared URL becomes useless quickly.

The server then redirects the client directly to the signed URL. The encrypted file is streamed from CloudFront — the server never proxies the file data itself.

The private key is loaded from disk at signing time and is never exposed outside the server.

---

## Automated Cleanup

The server runs a set of background cleanup routines that start at boot and run for the lifetime of the process, exiting cleanly when the server shuts down. Each routine runs on its own schedule and targets a specific category of stale data. They run concurrently and independently — a failure in one does not affect the others.

### S3 — stale pending files

When a file upload is initiated, a `pending` record is created in the database and a presigned S3 URL is issued. If the upload never completes — due to a client crash, network failure, or an aborted upload — the object may still exist in S3 with no corresponding completed record. The cleaner periodically queries for pending file records older than a configured age, checks whether the object exists in S3 via a `HeadObject` call, deletes it if found, and marks the database record as `failed`. If the object is already gone, the record is marked `failed` directly.

### Database — stale file metadata

File records marked as `failed` or `deleted` are retained in the database for a short period before being permanently removed. This gives enough time for any in-flight processes to resolve before the record disappears entirely.

### Database — stale refresh tokens

Refresh tokens that have been revoked or have passed their expiry are periodically purged from the database. This keeps the tokens table lean and avoids accumulating rows that serve no purpose.

### Database — stale download counters

The per-file daily download counters used by the CloudFront download protection are cleared on a daily cycle. This is what resets the download limit each day.

### Database — unverified users

User accounts that were registered but never verified within the allowed verification window are removed. This prevents the accumulation of ghost accounts from incomplete registrations.

### Database — expired OTPs

OTP records that have passed their expiry timestamp are cleaned up periodically, keeping the OTP table clear of records that can no longer be used.

---

## Related Documentation

- [Architecture](./architecture.md) — system overview and component relationships
- [File Upload](./uploading.md) — how storage quotas interact with the upload flow
- [Downloading and Decryption](./decryption-and-downloading.md) — how the download limit fits into the download flow
- [API Reference](./api.md) — error responses returned when limits are exceeded