# Self-Hosting

This guide covers how to run the Hashdrop server locally using Docker — for development or testing purposes.

---

## Prerequisites

Before proceeding, you will need the following AWS resources set up:

- **S3 bucket** — for encrypted file storage
- **CloudFront distribution** — connected to your S3 bucket, with a key pair configured for signed URLs
- **SES** — with a verified sender email address for OTP delivery
- **IAM credentials** — with appropriate permissions for S3, CloudFront, and SES

You will also need a CloudFront private key (`.pem` file) corresponding to the key pair ID configured in your CloudFront distribution.

---

## Running with Docker

### 1. Install Docker

Ensure Docker is installed and running on your machine. See the [official Docker installation guide](https://docs.docker.com/get-docker/) if needed.

### 2. Pull the image
```bash
docker pull anxhukumar/hashdrop-server
```

### 3. Set up the directory structure

Create a directory to hold the server's data:
```
hashdrop/
└── storage/
```

The `storage/` directory should be left empty — the database file will be created automatically and all migrations will run on first boot.

### 4. Run the container

From the directory containing your `hashdrop/` folder, run:
```bash
docker run -d \
  -p 8080:8080 \
  -v ./hashdrop/storage:/app/storage \
  -e PORT="" \
  -e PLATFORM="" \
  -e DB="" \
  -e JWT_SECRET="" \
  -e CLOUDFRONT_KEY_PAIR_ID="" \
  -e CLOUDFRONT_PRIVATE_KEY="$(cat /path/to/hashdrop-private.pem)" \
  -e USERID_HASHING_SALT="" \
  -e OTP_HASHING_SECRET="" \
  -e REFRESH_TOKEN_HASHING_SECRET_VERSION_1="" \
  -e AWS_ACCESS_KEY_ID="" \
  -e AWS_SECRET_ACCESS_KEY="" \
  -e AWS_REGION="" \
  anxhukumar/hashdrop-server:latest
```

Replace `/path/to/hashdrop-private.pem` with the actual path to your CloudFront private key file. Set `PLATFORM` to `dev` for local development or `prod` for production. Fill in the remaining values with your AWS credentials and configuration.

Update the port mapping if you chose a different value for `PORT`.

The server will start and the database will be initialised automatically.

> **Note:** If you are also using the Hashdrop CLI for local testing, update the API base URL in the CLI config to point to your local server (e.g. `http://localhost:8080`).

---

## Related Documentation

- [Installation](./installation.md) — CLI installation and Go setup
- [Testing](./testing.md) — running the test suite
- [Architecture](./architecture.md) — system overview
- [API Reference](./api.md) — endpoint reference for local testing