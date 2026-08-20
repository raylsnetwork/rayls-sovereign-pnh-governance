# Listener Event Wiring Script

This directory contains the `register_listener_event.py` script, which automates wiring a new contract event into the listener. It creates and updates the Go files required for the event to be processed in the pipeline. The script works both when the contract has never been mapped before and when you are just adding another event to an already mapped contract, as long as you provide the exact contract name and binding package name that the codebase already uses.

## Prerequisites

- Generate the contract binding using `rayls-contracts` task and copy the resulting package into `contracts/`, e.g. `contracts/MyContractV1`.
- Ensure the package builds correctly and is tracked in source control; the script imports the binding using the package name you provide at the prompt.

## Canonical event signature

The script prompts for the event's *canonical signature*. You can find it directly in the generated binding, usually exposed as a constant or metadata entry (for example `MyContractV1EventCompleted`). Another reliable hint is the inline comment that precedes each generated event handler. Search for the pattern `// Solidity: event MyEvent(` within the binding file and copy the signature right after it, stripping parameter names as described below.

The canonical signature must include **only the types**, without parameter names, and preserve the original order. For example, if the binding exposes:

```go
"EventCompleted(bytes32 indexed requestId, uint256 amount, (uint256,uint8) payload, uint256 timestamp)"
```

Enter in the prompt:

```go
EventCompleted(bytes32,uint256,(uint256,uint8),uint256)
```

## How to run

1. Open a terminal at the repository root (`cmd/listener/events/script/register_listener_event.py`).
2. Run the script: `python3 register_listener_event.py`.
3. Answer the prompts as described in the table below.

## Prompts overview

| Prompt | Example | Description |
| --- | --- | --- |
| `Contract name (MyContract):` | `TokenCore` | Contract name in PascalCase. It's used to generate Go constants and function names. |
| `Binding package name (e.g., MyContractV1):` | `TokenCoreV1` | Go package name for the binding located in `contracts/`. Typically matches the directory produced by `abigen`. |
| `Event name (e.g., EventCompleted):` | `TokenMinted` | Event name exactly as declared in the contract/binding, in PascalCase. |
| `Full event signature (e.g., EventCompleted(bytes32,uint256,(uint256,uint8),uint256):` | `TokenMinted(address,uint256)` | Canonical signature described above, without parameter names. |

## Output

After confirming the inputs, the script updates the relevant Go files (parsers, block processor, binding creation, event constants, and optionally creates a new core handler). So, you just have to review the generated changes and proceed with the data persistence workflow.