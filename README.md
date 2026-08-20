<div align="center">

# Rayls Governance API

**Governance & audit services for the Private Network Hub — indexing PNH activity, answering audit queries, and flagging anomalous cross-chain transactions.**

[![License: Apache 2.0][license-badge]][license-url]
[![Go][go-badge]][go-url]

[![Discord][discord-badge]][discord-url]
[![X][x-badge]][x-url]
[![LinkedIn][linkedin-badge]][linkedin-url]
[![YouTube][youtube-badge]][youtube-url]

[Architecture](#architecture) | [Configuration](#configuration) | [Running](#running-locally) | [API](#api-endpoints)

</div>

## What is this?

A multi-service system that gives Private Network operators the ability to verify the activity done by different Privacy Nodes and validate the correctness of cross-chain transactions. It runs as three services sharing a PostgreSQL database: an **API** (audit/governance queries), a **Listener** (indexes PNH events), and a **Flagger** (anomaly detection).

## Architecture

The system is composed of three independent services that share a PostgreSQL database:

```
┌─────────────────────────────────────────────────────┐
│                  Private Network Hub                │
│          (Ethereum-compatible blockchain)           │
└───────────────────────┬─────────────────────────────┘
                        │ RPC polling
                        ▼
            ┌───────────────────────┐
            │       Listener        │  port 8081
            │  Block event indexer  │
            └───────────┬───────────┘
                        │ NATS JetStream
                        ▼
            ┌───────────────────────┐
            │       Listener        │
            │   Event Dispatcher    |
            └───────────┬───────────┘
                        │ Routes & persists
                        ▼
            ┌───────────────────────┐     ┌───────────────────────┐
            │      PostgreSQL       │◄────│       Flagger         │  port 8082
            │      Database         │     │  Anomaly detection    │
            └───────────┬───────────┘     └───────────────────────┘
                        │
                        ▼
            ┌───────────────────────┐
            │         API           │  port 8080
            │   REST audit queries  │
            └───────────────────────┘
```

### Services

| Service | Port | Description |
|---------|------|-------------|
| **API** | 8080 | REST API for governance queries, audit operations, and authentication |
| **Listener** | 8081 | Indexes blockchain events from the Private Network Hub via RPC polling and NATS JetStream |
| **Flagger** | 8082 | Validates transactions and flags anomalies (balance checks, stale header proofs) |

## Prerequisites

- Private Network Hub (PNH) running and accessible via RPC
- PostgreSQL database
- NATS server with JetStream enabled (required by Listener)
- Contracts deployed on the Private Network Hub
- Rayls View key pair for the Private Network Operator (used to decrypt privacy node data)

## Configuration

All services are configured via environment variables. Copy the example template as a starting point and populate the values (env files are git-ignored — never commit real values).

```bash
cp config/.env.example .env.local
```

### Required Variables

```bash
# Database
DATABASE_TYPE=postgresql
DATABASE_CONNECTIONSTRING=postgresql://user:pass@host:5432/governance?sslmode=disable

# Private Network Hub
PNH_RPC_URL=http://pnh:3445
PNH_CHAIN_ID=1337
PNH_PRIVATE_KEY=0x...
PNH_DEPLOYMENT_PROXY_REGISTRY=0x...   # auto-loads contract addresses on startup
PNH_RAYLS_VIEW_SECRET_KEY=...         # secret key for decrypting PN data
PNH_STARTING_BLOCK=0                  # block to start indexing from
PNH_BATCH_SIZE=20                     # blocks per polling batch
PNH_OPERATOR_CHAIN_ID=999
PNH_HEADER_PROOF_EXPIRATION_PERIOD=5m

# Message Queue (Listener only)
NATS_URL=nats://nats:4222

# API
JWTSECRET=...
CORSURLS=http://localhost;https://your-frontend.example.com
LOGGING=development   # or: production

# Flagger (optional, defaults shown)
TRANSACTIONPROCESSOR_BATCHSIZE=100
TRANSACTIONPROCESSOR_CHECKINTERVAL=1s
```

> **Contract addresses** (`PNH_TELEPORT`, `PNH_TOKEN_CORE`, etc.) are loaded automatically from the `DeploymentProxyRegistry` contract at startup. They can also be set explicitly via environment variables if needed.

> **Rayls View key**: Generate a ML-KEM key pair for the Private Network operator. The secret key (`PNH_RAYLS_VIEW_SECRET_KEY`) is used solely to decrypt encrypted participant data from the Private Network Hub, it is not used for signing or asset transfers.

## Running Locally

### Docker Compose

The repository includes a `docker-compose.yml` that starts all three services together with PostgreSQL and NATS.

**1. Create your environment file**

```bash
cp .env.docker.example .env.docker
```

Edit `.env.docker` and fill in the required values (RPC URL, private key, contract addresses, etc.). The file is git-ignored — never commit it.

**2. Start the stack**

```bash
docker compose up
```

This starts:

| Container | Port(s) | Description |
|-----------|---------|-------------|
| `postgres` | 5433 | PostgreSQL 16 |
| `nats` | 4222, 8222 | NATS JetStream (client + monitoring) |
| `nui` | 31311 | NATS UI (browser dashboard) |
| `listener` | 2345 (dlv) | Listener service with live reload |
| `flagger` | 2346 (dlv) | Flagger service with live reload |
| `api` | 8080, 2347 (dlv) | API service with live reload |

The `listener`, `flagger`, and `api` containers mount the local source tree and run with `dlv` — attach VS Code debugger to ports 2345, 2346, 2347 respectively.

**3. Stop the stack**

```bash
docker compose down           # keep volumes
docker compose down -v        # also remove database and NATS data
```

---

### VS Code Debugger

The repository already includes `.vscode/launch.json` with configurations for all three services. Open the **Run and Debug** tab (`Ctrl+Shift+D`), select the service you want to start, and press **F5**. Make sure `.env.local` exists and is populated before launching.

The services can run simultaneously, each binds to its own port. Expected output per service:

| Service | Ready message |
|---------|---------------|
| API | `[GIN-debug] Listening and serving HTTP on :8080` |
| Listener | `Starting block processor` |
| Flagger | `Starting transaction processor` |

### Terminal

```bash
go run ./cmd/api/main.go      run --config .env.local
go run ./cmd/listener/main.go run --config .env.local
go run ./cmd/flagger/main.go  run --config .env.local
```

## Build & Run

### Docker (standalone images)

```bash
# Build
docker build -f Dockerfile.api      -t governance-api      .
docker build -f Dockerfile.listener -t governance-listener .
docker build -f Dockerfile.flagger  -t governance-flagger  .

# Run (supply an env file — see .env.docker.example)
docker run --env-file .env.docker -p 8080:8080 governance-api
docker run --env-file .env.docker -p 8081:8081 governance-listener
docker run --env-file .env.docker -p 8082:8082 governance-flagger
```

### Makefile

```bash
make build      # compile all three binaries into build/
make api        # compile API (also regenerates Swagger docs)
make listener   # compile listener
make flagger    # compile flagger
make swagger    # regenerate OpenAPI docs (requires swag)
make lint       # run golangci-lint
```

### Binary

```bash
./build/api      run --config .env.local
./build/listener run --config .env.local
./build/flagger  run --config .env.local
```

## API Endpoints

The full Swagger UI is available at [`http://localhost:8080/swagger/index.html`](http://localhost:8080/swagger/index.html) when the API is running.

### Audit (public)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/audit/transactions` | List transactions with filters |
| `GET` | `/audit/transactions/message/:id` | Transaction by message ID |
| `GET` | `/audit/transactions/batch/:id` | Transactions in a batch |
| `GET` | `/audit/transactions/dvp/:id` | DVP swap transactions |
| `GET` | `/audit/participants/:chainId` | Participants on a chain |
| `GET` | `/audit/tokens` | Token registry info |
| `GET` | `/audit/header-proofs` | Header proof events |
| `GET` | `/flagged` | Flagged transactions |

### Governance (requires auth)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/resources/:chainId/*` | Token balances for a chain |
| `GET` | `/resource_info_all_chains/:resourceId` | Token status across all chains |
| `GET` | `/participant_info/:chainId` | Participant info (auth) |
| `GET` | `/token_status/:resourceId` | Token registry status |

### Authentication

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/private-network/signup` | Register a VEN operator |
| `POST` | `/private-network/login` | Authenticate and receive JWT cookie |

Authentication uses JWT tokens stored in HTTP-only cookies (`SameSite=Lax`, `Secure=true`).

## Database Migrations

Migrations run automatically on service startup using [`golang-migrate`](https://github.com/golang-migrate/migrate).

Migration files live in `/migrations` and follow the naming convention:

```
{timestamp}_{description}.up.sql
{timestamp}_{description}.down.sql
```

## Testing

```bash
# Unit + integration tests per service
make api-coverage
make listener-coverage
make flagger-coverage

# Or directly with go test
go test ./cmd/api/...
go test ./cmd/listener/...
go test ./cmd/flagger/...
```

Integration tests spin up a real PostgreSQL container via [`dockertest`](https://github.com/ory/dockertest) - Docker must be available on the machine.

## Contributing

We are not accepting external contributions at this time — see [CONTRIBUTING.md](./CONTRIBUTING.md). Please also read our [Code of Conduct](./CODE_OF_CONDUCT.md).

## Security

To report a security vulnerability, see [SECURITY.md](./SECURITY.md) — please do not open a public issue.

## License

Licensed under the Apache License, Version 2.0 — see [LICENSE](./LICENSE).

This project links third-party libraries that remain under their own licenses; notably [go-ethereum](https://github.com/ethereum/go-ethereum) under the LGPL-3.0 (library packages only) and the HashiCorp `errwrap` / `go-multierror` libraries under the MPL-2.0. See [NOTICE](./NOTICE).

Copyright 2026 Rayls Core Ltd.

[license-badge]: https://img.shields.io/badge/License-Apache_2.0-blue.svg
[license-url]: ./LICENSE
[go-badge]: https://img.shields.io/badge/Go-1.26.4-00ADD8?logo=go&logoColor=white
[go-url]: ./go.mod
[discord-badge]: https://img.shields.io/badge/Discord-join%20chat-5865F2?logo=discord&logoColor=white
[discord-url]: https://discord.gg/6THZ96357r
[x-badge]: https://img.shields.io/badge/X-%40RaylsLabs-000000?logo=x&logoColor=white
[x-url]: https://x.com/RaylsLabs
[linkedin-badge]: https://img.shields.io/badge/LinkedIn-Rayls-0A66C2?logo=linkedin&logoColor=white
[linkedin-url]: https://www.linkedin.com/company/rayls/
[youtube-badge]: https://img.shields.io/badge/YouTube-Rayls-FF0000?logo=youtube&logoColor=white
[youtube-url]: https://www.youtube.com/@Rayls_blockchain
