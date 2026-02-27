# Self-Hosting

This guide covers two ways to run the Hashdrop server locally — for development or testing purposes. The recommended approach is Docker.

---

## Prerequisites

Before proceeding with either method, you will need the following AWS resources set up:

- **S3 bucket** — for encrypted file storage
- **CloudFront distribution** — connected to your S3 bucket, with a key pair configured for signed URLs
- **SES** — with a verified sender email address for OTP delivery
- **IAM credentials** — with appropriate permissions for S3, CloudFront, and SES

You will also need a CloudFront private key (`.pem` file) corresponding to the key pair ID configured in your CloudFront distribution. This is the `hashdrop-private.pem` file referenced throughout this guide.

---

## Method 1 — Docker (Recommended)

### 1. Install Docker

Ensure Docker is installed and running on your machine. See the [official Docker installation guide](https://docs.docker.com/get-docker/) if needed.

### 2. Pull the image
```bash
docker pull anxhukumar/hashdrop-server
```

### 3. Set up the directory structure

Create a directory to hold the server's configuration and data, then create two subdirectories inside it:
```
hashdrop/
├── secrets/
│   ├── .env
│   └── hashdrop-private.pem
└── storage/
```

The `storage/` directory should be left empty — the database file will be created automatically and all migrations will run on first boot.

### 4. Configure the environment

Create a `.env` file inside `secrets/` with the following structure:
```dotenv
PORT=""

PLATFORM=""

DB="file:storage/storage.db?_foreign_keys=on"

JWT_SECRET=""

CLOUDFRONT_KEY_PAIR_ID=""
CLOUDFRONT_PRIVATE_KEY_PATH="secrets/hashdrop-private.pem"

USERID_HASHING_SALT=""
OTP_HASHING_SECRET=""
REFRESH_TOKEN_HASHING_SECRET_VERSION_1=""

AWS_ACCESS_KEY_ID=""
AWS_SECRET_ACCESS_KEY=""
AWS_REGION=""
```

Set `PLATFORM` to `dev` for local development or `prod` for production. Fill in the remaining values with your AWS credentials and configuration.

Place your CloudFront private key at `secrets/hashdrop-private.pem`.

### 5. Run the container

From the directory containing your `hashdrop/` folder, run:
```bash
docker run -d \
  -p 8080:8080 \
  -v ./hashdrop/secrets:/app/secrets:ro \
  -v ./hashdrop/storage:/app/storage \
  anxhukumar/hashdrop-server:latest
```

Update the port mapping if you chose a different value for `PORT`.

The server will start and the database will be initialised automatically.

> **Note:** If you are also using the Hashdrop CLI for local testing, update the API base URL in the CLI config to point to your local server (e.g. `http://localhost:8080`).

---

## Method 2 — Manual

### 1. Clone the repository
```bash
git clone https://github.com/anxhukumar/hashdrop
cd hashdrop/server
```

The `cli/` directory is included in the clone but is not needed for running the server.

### 2. Install required tools

Ensure the following are installed:

**Go** — see the [Installation guide](./installation.md) for setup instructions.

**sqlc** — used to generate the type-safe database query code:
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### 3. Generate database code

From inside the `server/` directory:
```bash
sqlc generate
```

This generates the `internal/database/` package from the SQL query definitions.

### 4. Set up secrets and storage

Create two directories inside `server/`:
```
server/
├── secrets/
│   ├── .env
│   └── hashdrop-private.pem
└── storage/
```

The `storage/` directory should be left empty — the database file will be created automatically and all migrations will run on first boot.

### 5. Configure the environment

Create a `.env` file inside `secrets/` with the following structure:
```dotenv
PORT=""

PLATFORM=""

DB="file:storage/storage.db?_foreign_keys=on"

JWT_SECRET=""

CLOUDFRONT_KEY_PAIR_ID=""
CLOUDFRONT_PRIVATE_KEY_PATH="secrets/hashdrop-private.pem"

USERID_HASHING_SALT=""
OTP_HASHING_SECRET=""
REFRESH_TOKEN_HASHING_SECRET_VERSION_1=""

AWS_ACCESS_KEY_ID=""
AWS_SECRET_ACCESS_KEY=""
AWS_REGION=""
```

Set `PLATFORM` to `dev` for local development or `prod` for production. Fill in the remaining values with your AWS credentials and configuration.

Place your CloudFront private key at `secrets/hashdrop-private.pem`.

### 6. Run the server

To run directly:
```bash
go run .
```

Or build and run the executable:
```bash
go build -o build/hashdrop-server .
./build/hashdrop-server
```

> **Note:** If you are also using the Hashdrop CLI for local testing, update the API base URL in the CLI config to point to your local server (e.g. `http://localhost:8080`).

---

## Related Documentation

- [Installation](./installation.md) — CLI installation and Go setup
- [Testing](./testing.md) — running the test suite
- [Architecture](./architecture.md) — system overview
- [API Reference](./api.md) — endpoint reference for local testing