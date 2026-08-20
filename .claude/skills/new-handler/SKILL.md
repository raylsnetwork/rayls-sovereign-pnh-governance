---
name: new-handler
description: >
  Scaffold a new listener event handler for processing smart contract events.
  Use when the user wants to add a new event, wire a new contract, or handle
  a new blockchain event in the listener service.
argument-hint: "[ContractName] [EventName] [EventSignature]"
---

# New Listener Event Handler

You are adding a new event handler to the listener service. A Python script handles
the mechanical wiring. Your job is to gather inputs, run the script, then complete
what it can't do: persistence logic, tests, and verification.

The user will provide:
- **Contract name** (PascalCase, e.g., `EnygmaTeleport`) — use `$ARGUMENTS[0]` if provided
- **Event name** (PascalCase, e.g., `EnygmaTransferTimeout`) — use `$ARGUMENTS[1]` if provided
- **Event signature** (canonical, types only, e.g., `EnygmaTransferTimeout(bytes32,uint256)`) — use `$ARGUMENTS[2]` if provided

If any of these are missing, ask the user before proceeding.

## Prerequisites

Before starting, verify:
1. The contract binding package exists under `contracts/` (e.g., `contracts/EnygmaTeleportV1/`)
2. The binding has a `Parse<EventName>` method (e.g., `ParseEnygmaTransferTimeout`)
3. You know the **binding package name** (e.g., `EnygmaTeleportV1`) — this is the directory name under `contracts/`

## Step 1: Run the wiring script

The script at `cmd/listener/events/script/register_listener_event.py` automates all mechanical wiring.
It is interactive and requires four inputs. Pipe them via stdin:

```bash
printf '<ContractName>\n<BindingPackage>\n<EventName>\n<EventSignature>\n' | python3 cmd/listener/events/script/register_listener_event.py
```

Example:
```bash
printf 'EnygmaTeleport\nEnygmaTeleportV1\nEnygmaTransferTimeout\nEnygmaTransferTimeout(bytes32,uint256)\n' | python3 cmd/listener/events/script/register_listener_event.py
```

The script updates these files automatically:
- `cmd/listener/events/events.go` — event constants
- `cmd/listener/adapters/indexer/utils.go` — event parser + registry entry
- `cmd/listener/adapters/indexer/log_parser.go` — Contracts struct + instantiation (new contracts only)
- `contracts/creator.go` — creator function (new contracts only)
- `cmd/listener/core/handlers/<contract>_handler.go` — handler file with event switch + stub method

## Step 2: Review generated code

Read the files the script modified and verify:
- Event constants are correct
- Parser uses the right binding method (`Parse<EventName>`)
- For new contracts: Contracts struct field, creator function, and registry entry look correct

## Step 3: Implement persistence logic

The script generates a stub handler with a `// TODO` comment. Fill in the actual logic:
1. Type-assert the event data from the binding
2. Extract and transform event data into domain objects
3. Persist via repository
4. Add meaningful log entries

Use `cmd/listener/core/handlers/proofs_handler.go` as a minimal reference for the persistence pattern.

## Step 4: Write tests

Create or extend `cmd/listener/core/handlers/<contract>_handler_test.go`.

Follow the existing test pattern:
- Use mocks from `cmd/listener/mocks/`
- Test the happy path
- Test type assertion failure
- Test repository error propagation

## Step 5: Verify

```bash
make build
make lint
```

## Checklist

- [ ] Script ran successfully (all files show `[updated]` or `[created]`)
- [ ] Generated code reviewed
- [ ] Handler persistence logic implemented (TODO replaced)
- [ ] Tests written
- [ ] Build passes (`make build`)
- [ ] Lint passes (`make lint`)
