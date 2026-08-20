---
name: architecture
description: >
  Project architecture knowledge for the governance-api. Use when the user asks
  about project structure, how services connect, where to put new code, how the
  listener/flagger/api work, or when making architectural decisions.
---

# Governance API Architecture

## Overview

The governance-api enables Private Network operators to monitor Privacy Node
activity and verify transaction correctness on the Private Network Hub.

## Three Microservices

| Service | Entry | Port | Responsibility                                                                                                                      |
|---------|-------|------|-------------------------------------------------------------------------------------------------------------------------------------|
| **API** | `cmd/api/main.go` | 8080 | REST API (Gin) for audit queries, participants, tokens, balances, header proofs                                                     |
| **Listener** | `cmd/listener/main.go` | 8081 | Blockchain event listener: processes blocks from Private Network Hub, parses smart contract logs, decrypts payloads, persists to DB |
| **Flagger** | `cmd/flagger/main.go` | 8082 | Transaction validator: checks balances, flags invalid transactions (negative balances)                                              |

## Hexagonal Architecture (per service)

Each service follows the same layered structure:

```
cmd/<service>/
  main.go              # Entrypoint (Cobra CLI)
  app/
    app.go             # Bootstrap: wires adapters -> core -> infrastructure
    cmd.go             # Cobra command definition
  core/
    ports.go           # All interfaces (primary + secondary ports)
    services/          # Business logic implementations
    handlers/          # Event-specific handlers (listener only)
  adapters/
    handlers/          # HTTP handlers with Swagger annotations (API only)
    repositories/      # GORM-based PostgreSQL implementations
    config/            # Config adapter (listener only)
    crypto/            # Decryption adapter (listener only)
    indexer/           # Blockchain log parsing (listener only)
  infrastructure/
    infrastructure.go  # DB connection, Ethereum client setup
  middleware/          # HTTP middleware (API only)
  testutil/            # Fakes, stubs for testing
  mocks/               # Generated mocks (listener)
  msgqueue/            # NATS JetStream publisher/consumer (listener only)
```

### Dependency flow

```
main.go -> app.go (bootstrap)
  -> infrastructure.SetupInfrastructure()   # DB + Eth client
  -> adapters (repositories, indexer, etc.)  # Implement core interfaces
  -> core services                           # Business logic via interfaces
  -> adapters (HTTP handlers)                # Expose to outside world
```

All dependencies point **inward**: adapters depend on core interfaces, never the reverse.

## Shared Packages (root level)

| Package | Purpose |
|---------|---------|
| `domain/` | Domain models: `Transaction`, `Participant`, `Token`, `Balance`, `HeaderProofEvent`, `EnygmaTransaction`, `RevertDataTransaction`. Base `Model` with UUID PK + UTC timestamps. Custom types: `BigInt` |
| `dto/` | Data Transfer Objects: filters, responses, pagination structs |
| `types/` | Enums and bidirectional string mappings (`TxType`, `ProtocolType`, `AssetType`, etc.) |
| `config/` | Viper-based config loading with `ExperimentalBindStruct()` for env var binding |
| `contracts/` | Smart contract bindings generated via `abigen` (12 contracts) |
| `migrations/` | SQL migrations via `golang-migrate` |
| `logger/` | `slog`-based structured logging |
| `cryptography/` | DH key exchange utilities |
| `infrastructure/` | Shared DB setup |
| `adapters/` | Shared adapters (TokenRegistry) |
| `withstack/` | Error wrapping with stack traces |
| `flags/` | CLI flag definitions |

## Listener Pipeline

```
BlockProcessor.Run(ticker)
  -> getNextBlockToProcess()
  -> processBlocks(fromBlock)
     -> logParser.ParseLogs(from, to)
  -> logPublisher.Publish(logs)          # Publish to NATS JetStream
  -> blockRepo.UpdateLatestProcessedBlock()

EventDispatcher.Run()                    # Separate goroutine
  -> consumer.Next()                     # Consume from NATS JetStream
  -> dispatch by log.ContractName        # Hash map lookup
  -> handler.Handle(ctx, log)
  -> msg.Ack()
```

### Event Handlers (8 registered)

| Handler | Contract | Events Handled |
|---------|----------|----------------|
| `TeleportEventHandler` | Teleport | `EncryptedDataBatchStored`, `AtomicMessageAdditionalDataBatch`, `AtomicMessageStatusChangedBatch` |
| `EnygmaTeleportEventHandler` | EnygmaTeleport | `EnygmaTransfer`, `EnygmaTransferCompleted`, `EnygmaSupplyUpdated`, `EnygmaDvpBalanceUpdated` |
| `TokenCoreEventHandler` | TokenCore | `Erc20TokenRegistered`, `Erc721TokenRegistered`, `Erc1155TokenRegistered`, `DvpErc721TokenRegistered`, `DvpErc1155TokenRegistered`, `TokenStatusUpdated`, `TokenBalanceUpdated` |
| `ParticipantCoreEventHandler` | ParticipantCore | `ParticipantRegistered`, `ParticipantUpdated` |
| `AuditManagerEventHandler` | AuditManager | `NewAuditOrChainInfo` |
| `EnygmaTokenManagerEventHandler` | EnygmaTokenManager | `EnygmaTokenRegistered` |
| `ProofsEventHandler` | Proofs | `HeaderProofSubmitted` |
| `DvpTeleportEventHandler` | DvpTeleport | `TransferEncryptedData`, `SwapStateChanged` |

### Adding a new event handler

1. Define event name + signature constants in `cmd/listener/events/events.go`
2. Register the event signature in `cmd/listener/adapters/indexer/log_parser.go`
3. Create handler in `cmd/listener/core/handlers/` implementing `EventHandler` interface:
   - `ContractName() string`
   - `Handle(ctx context.Context, log ContractLog) error`
   - `Name() string`
4. Wire the handler in `cmd/listener/app/app.go` -> `createEventHandlers()`
5. Alternatively, use the Python script: `python3 cmd/listener/events/script/register_listener_event.py`

## Flagger Pipeline

Two concurrent processors run via `errgroup`:

```
TransactionProcessor.Run(ticker @ 1s)
  -> txRepo.GetTransactions(batchSize)    # Unprocessed txs (executed/completed/null status)
  -> processTransaction(tx)
     -> switch tx.TxType:
        CrossChain: debit source, credit dest, flag if source < 0
        Mint: credit source
        Burn: debit source, flag if < 0
     -> Each runs inside a DB transaction (txRepo.Transact)
  -> Mark processed / flag if negative balance

HeaderLivelinessProcessor.Run(ticker @ 1m)
  -> Check header proof freshness against expiration period
  -> Flag stale participants
```

## API Routes

All audit routes under `/audit`:

```
GET  /audit/transactions                          # List with filters + pagination
GET  /audit/transactions/:messageId               # By message ID
GET  /audit/transactions/dvp/:transactionId       # DVP by transaction ID
GET  /audit/transactions/dvp/swap/:sharedId       # DVP swap legs
GET  /audit/transactions/batch/:batchId           # Regular batch
GET  /audit/transactions/enygma/batch/:batchId    # Enygma batch
GET  /audit/participants                           # List with filters
GET  /audit/participants/:chainId                  # By chain ID
GET  /audit/tokens                                 # List with filters
GET  /audit/tokens/:resourceId                     # By resource ID
GET  /audit/header-proofs                          # List with filters + pagination
GET  /flagged                                      # Flagged transactions

# Auth-protected routes (JWT in cookie)
POST /ven/signup
POST /ven/login
GET  /resources/:chainid/*resourceid               # Balances in chain
GET  /resource_info_all_chains/:resourceid          # Balance across all chains
POST /resource_info_list_chains                     # Balance across specific chains
GET  /participant_info/:chainId                     # Participant info
GET  /token_status/:resourceid                      # Token registry status
```

## Key Patterns

- **Interface-driven**: All cross-boundary calls go through interfaces defined in `core/ports.go`
- **Compile-time checks**: `var _ Interface = (*Struct)(nil)` pattern in repositories
- **Receiver names**: Short lowercase (`r`, `p`, `t`, `s`, `h`, `tc`, `bp`)
- **Error handling**: `cockroachdb/errors` + `withstack.Wrap()` for stack traces
- **Enums**: `uint8` with bidirectional `map[string]Type` / `map[Type]string` in `types/`
- **Custom GORM types**: `BigInt`, `DecimalArray`, `StringArray` implementing `driver.Valuer` + `sql.Scanner`
- **Event handler pattern**: `ContractName` / `Handle` / `Name` with hash-map dispatching via EventDispatcher
- **Config**: Viper `ExperimentalBindStruct()` with `mapstructure` tags for env var binding

## Build and Test

```bash
make build                  # Build all 3 services
make lint                   # golangci-lint v2 (formatters + linters)
make swagger                # Regenerate OpenAPI docs
make listener-coverage      # Listener test coverage
make flagger-coverage       # Flagger test coverage
make api-coverage           # API test coverage
./test.sh                   # Integration tests (requires Docker)
```
