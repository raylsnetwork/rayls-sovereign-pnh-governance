---
name: investigate
description: >
  Investigate bugs and data issues across all three governance-api services
  (listener, flagger, api) and the NATS message queue. Use when the user reports
  a data issue, wrong API response, missing field, incorrect value, "record not found"
  error, stub token appearing, negative balance, missing token, wrong status, or any
  bug that could originate in any layer of the pipeline. Also use when the user pastes
  an error log, asks "why is this field empty", says "something looks wrong", or wants
  to trace how data flows through the system.
argument-hint: '[endpoint or entity] [symptom]'
---

# Investigate Bug

You are investigating a bug in the governance-api. Your goal is to find the **root cause**,
not just identify which layer has the problem. Trace data from its origin (blockchain event)
through every transformation until the point where the bug becomes visible.

## System Overview

Three services share a PostgreSQL database. A bug visible in the API may originate anywhere:

```
Blockchain → Listener → NATS JetStream → EventDispatcher → Handler → DB → Flagger → API
```

| Service | Entry | Role |
|---------|-------|------|
| **Listener** | `cmd/listener/` | Processes blocks, parses logs, publishes to NATS, dispatches events to handlers that persist data |
| **Flagger** | `cmd/flagger/` | Validates transactions, updates balances, flags anomalies |
| **API** | `cmd/api/` | Queries DB, transforms data, serves REST responses |

## Root Cause Investigation Method

Work backwards from the symptom to the source. At each layer, read the actual code — don't
guess based on file names alone.

### Step 1: Understand the symptom

Before reading any code, clarify:
- What entity is affected? (transaction, token, participant, balance, header proof)
- What's wrong? (missing, wrong value, stale, unexpected status)
- Is there an error log? (extract handler name, event name, and error message)

### Step 2: Trace backwards from the API

Start at the layer where the bug is visible and work toward the data source.

**API layer** — is the data being queried and transformed correctly?
- Handler: `cmd/api/adapters/handlers/` — params extraction, response format
- Service: `cmd/api/core/` — business logic, validation, error mapping
- Repository: `cmd/api/adapters/repositories/` — SQL query, JOINs, aggregation
- Mapper/DTO: `cmd/api/core/transaction_mapper.go`, `dto/` — field mapping

**Database** — is the stored data correct?
- Read the repository query to understand what columns and JOINs are used
- Check if the issue is: missing data (never written), wrong data (bad value), stale data (not updated)

**Listener layer** — was the data persisted correctly?
- Handler: `cmd/listener/core/handlers/` — event parsing, field extraction, persistence calls
- Repository: `cmd/listener/adapters/repositories/` — upsert logic, ON CONFLICT behavior
- On-chain adapter: `adapters/tokenregistry/` — blockchain RPC calls for token metadata

**NATS layer** — was the event delivered and processed?
- Publisher: `cmd/listener/msgqueue/publisher.go` — publishes contract logs after block processing
- Consumer: `cmd/listener/msgqueue/consumer.go` — delivers messages to EventDispatcher
- Manager: `cmd/listener/msgqueue/manager.go` — stream config, subject namespacing by chainId
- Dispatcher: `cmd/listener/core/services/event_dispatcher.go` — routes by `ContractName`, acks on success
- If a handler returns an error, the message is NOT acked and will be redelivered (up to 10 times)

**Block processing** — was the event captured from the blockchain?
- `cmd/listener/core/services/block_processor.go` — fetches blocks in batches
- `cmd/listener/adapters/indexer/log_parser.go` — parses raw logs into ContractLog structs
- `cmd/listener/adapters/indexer/utils.go` — event signature registry
- `cmd/listener/core/services/log_publisher.go` — publishes parsed logs to NATS

**Flagger layer** — was the data modified after persistence?
- `cmd/flagger/core/transaction_processor.go` — balance calculations by TxType (CrossChain/Mint/Burn)
- `cmd/flagger/core/header_liveliness_processor.go` — header proof expiration checks

### Step 3: Identify the root cause

Once you find where the data diverges from what's expected, determine **why**:
- Is it a logic error in the code?
- Is it a timing/ordering issue (event processed before dependency exists)?
- Is it an infrastructure issue (RPC node returning stale data)?
- Is it a data model issue (wrong type, missing field, enum mismatch)?

Always explain the full causal chain: "Event X happens → handler Y does Z → but Z fails because..."

## Data Flow by Entity

### Transactions

```
Blockchain event
  → BlockProcessor.processBlocks()
  → LogPublisher.Publish()          → NATS "events.<chainId>.contract_logs"
  → EventDispatcher.Run()           → routes to handler by ContractName
  → teleport_handler.go             → processEncryptedDataBatchStored()
     dvp_teleport_handler.go           processTransferEncryptedDataEvent()
     enygma_teleport_handler.go        processEnygmaTransferEvent()
  → transaction_repository.go       → Create/Update in DB
  → Flagger picks up               → transaction_processor.go
  → API serves                     → transaction_handler.go → transaction_mapper.go
```

**Listener handlers by protocol:**
| Protocol | Handler | Key method |
|----------|---------|------------|
| Regular (Teleport) | `teleport_handler.go` | `processEncryptedDataBatchStored()` |
| DVP Swap | `dvp_teleport_handler.go` | `processTransferEncryptedDataEvent()` |
| Enygma | `enygma_teleport_handler.go` | `processEnygmaTransferEvent()` |
| Status updates | `teleport_handler.go` | `processAtomicMessageStatusChangedBatch()` |
| Destination data | `teleport_handler.go` | `processAtomicMessageAdditionalDataBatch()` |

**API endpoints:**
| Endpoint | Repository method |
|----------|------------------|
| `GET /audit/transactions` | `FindWithFilters()` |
| `GET /audit/transactions/:messageId` | `FindByMessageId()` |
| `GET /audit/transactions/dvp/:transactionId` | `FindByTransactionId()` |
| `GET /audit/transactions/dvp/swap/:sharedId` | `FindBySharedId()` |
| `GET /audit/transactions/batch/:batchId` | `FindByBatchIdPaginated()` |
| `GET /audit/transactions/enygma/batch/:batchId` | `FindByEnygmaBatchId()` |
| `GET /flagged` | `FindFlagged()` |

### Tokens

```
Blockchain event (Erc20TokenRegistered, TokenStatusUpdated, TokenBalanceUpdated, etc.)
  → token_core_handler.go
  → persistTokenRegistry()          → calls on-chain TokenRegistry via RPC
     ├── Registry returns valid data → Upsert with real metadata
     ├── Registry returns no data    → Upsert STUB token (Name="Unknown", Symbol="Unknown")
     └── Issuer mismatch            → returns ErrIssuerMismatch, event skipped
  → token_repository.go::Upsert()   → ON CONFLICT (resource_id) DO UPDATE
  → API: token_handler.go → token_service.go → token_repository.go
```

**Stub token pattern (common issue):**
When `persistTokenRegistry()` can't get token data from the blockchain registry (RPC node
inconsistency, lagging replica), it persists a stub with `Name="Unknown"`, `Symbol="Unknown"`,
and the issuer from the event. The Upsert will overwrite existing data if it runs again.
If you see "Unknown" tokens in the API, this is the cause — check RPC node health and whether
a subsequent event updated the stub.

**Enygma tokens** follow a separate path:
- `enygma_token_manager_handler.go` → `processEnygmaTokenRegistered()` / `processTokenRegistered()`
- Also calls `persistTokenRegistry()` with the same stub token fallback
- Returns `ErrIssuerMismatch` if token belongs to another network (logged, skipped)

### Balances

```
Listener                    Flagger                         API
token_core_handler.go  →    transaction_processor.go   →    balance_handler.go
(TokenBalanceUpdated)       UpdateBalance()                 balance_service.go
(initial mint txs)          UpdateSenderReceiverBalances()  balance_repository.go
```

**Common issues:** Negative balances, wrong amounts — check `transaction_processor.go` switch
on `TxType` (CrossChain debit/credit, Mint credit, Burn debit). Check if the flagger processes
transactions out of order.

### Participants

```
Listener                          Flagger                          API
participant_core_handler.go  →    header_liveliness_processor.go → participant_handler.go
  ↓                               FlagParticipant()                participant_service.go
participant_repository.go                                          participant_repository.go
```

**Common issues:** Wrong status/role — check `types/` enum mappings (`ParticipantStatus`,
`ParticipantRole`). Wrong flag state — check flagger's expiration logic and whether header
proofs are being submitted.

### Header Proofs

```
Listener                    Flagger                            API
proofs_handler.go      →    header_liveliness_processor.go →   header_proof_handler.go
  ↓                         GetLatestHeaderProofs()
header_proof_repository.go  header_flag_event_repository.go
```

### Enygma Transactions

```
Listener                              API
enygma_teleport_handler.go       →    transaction_handler.go
  processEnygmaTransferEvent()        FindByEnygmaBatchId()
  processEnygmaTransferCompleted()      (JOIN enygma_transactions)
  ↓
transaction_repository.go::CreateTransactionsWithEnygmaData()
enygma_transaction_repository.go::UpsertEnygmaTransactionsBulk()
```

**Common issues:** Missing timestamps (listener doesn't set them initially), missing `RnHash`
(completion event not received yet).

### Token Freeze

```
Blockchain event (TokenFreezeStatusChanged)
  → token_freeze_manager_handler.go → processTokenFreezeStatusChanged()
  → token_freeze_repository.go     → UpdateTokenFreezeStatus() (audit + state in single tx)
  → DB: token_freeze_states, token_freeze_audits
  → API: token_handler.go          → RestrictedChainIds in token response
         token_repository.go       → JOIN with token_freeze_states (is_frozen = true)
```

**Common issues:** Token appears frozen but shouldn't be — check `token_freeze_states` table
for stale entries. `RestrictedChainIds` empty when expected — check if the JOIN in
`token_repository.go` filters correctly on `is_frozen = true`.

## Common Bug Patterns

| Symptom | Likely cause | Where to look |
|---------|-------------|---------------|
| "record not found" in handler logs | Token not yet persisted when dependent event arrives | `persistTokenRegistry()` in `token_core_handler.go` |
| Token shows "Unknown" name/symbol | Stub token persisted due to RPC inconsistency | `persistTokenRegistry()` stub fallback path |
| Missing transaction data | Event not parsed or handler not registered | `log_parser.go` event registry, `app.go` handler wiring |
| Event processed but data not in DB | Handler error → message not acked → redelivery | Check listener logs for handler errors |
| Stale data after update | Upsert ON CONFLICT clobbered with old values | Check column list in ON CONFLICT DO UPDATE |
| Wrong enum value in API response | Mismatch between DB uint8 and `types/` string map | `types/` bidirectional maps |
| Negative balance flagged | Flagger processed transactions out of order | `transaction_processor.go` TxType switch |
| Event not reaching handler | Wrong ContractName or event not in NATS subject | `event_dispatcher.go` handler map, `log_parser.go` |

## Key Files Reference

**Domain models:** `domain/{transaction,balance,token,participant,header_proof_event,enygma_transaction,token_freeze}.go`

**DTOs:** `dto/{transaction_detail,transaction_list,batch_transaction,dvp_swap_transaction,token_registry_status}.go`

**Type mappings:** `types/` — bidirectional enum maps for `TxType`, `ProtocolType`, `AssetType`,
`TeleportStatus`, `ErcStandard`, `ParticipantStatus`, `ParticipantRole`, `FreezeAction`

**Transaction mapper:** `cmd/api/core/transaction_mapper.go` — `ToTransactionDetailDto()`,
`ToBatchTransactionDtos()`, `ToDvpSwapTransactionDtos()`

**Error types:** `cmd/api/core/errors.go` — `NotFoundError`, `ValidationError`, `InternalError`

**Sentinel errors:** `cmd/listener/core/ports.go` — `ErrIssuerMismatch`; `cmd/api/core/ports.go` — `ErrRecordNotFound`

**NATS config:** `cmd/listener/msgqueue/manager.go` — stream "EVENTS", subjects "events.<chainId>.*",
WorkQueuePolicy, 1h dedup window, max 10 redeliveries

## Investigation Checklist

When investigating, read the actual code at each layer — don't just check file names.

- [ ] **Symptom clear**: entity, field, expected vs actual value identified
- [ ] **API layer**: handler params, service logic, repository query, DTO/mapper
- [ ] **Database**: data exists with expected values (or identify what's missing/wrong)
- [ ] **Listener layer**: handler sets all required fields, repository persists correctly
- [ ] **NATS layer**: event published, consumed, dispatched to correct handler, ack behavior
- [ ] **Blockchain**: event emitted, log parsed, event signature registered
- [ ] **Flagger layer**: balance calculation, flagging logic (if relevant)
- [ ] **Type mappings**: enum values match between `types/` maps and DB values
- [ ] **Root cause**: full causal chain identified, not just the layer where it breaks
