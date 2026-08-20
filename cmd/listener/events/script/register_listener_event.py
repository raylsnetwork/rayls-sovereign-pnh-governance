import re
import sys
from pathlib import Path
from typing import Callable
from utils import (
    EventSpec,
    UpdateContext,
    prompt_non_empty,
    block_processor_case,
    core_event_handler,
    core_file,
    core_switch_case,
    creator_function,
    parser_event_entry,
    parser_function,
    parser_map_entry,
)


SCRIPT_DIR = Path(__file__).resolve().parent
REPO_ROOT = SCRIPT_DIR.parents[3]
CMD_DIR = REPO_ROOT / "cmd"
LISTENER_DIR = CMD_DIR / "listener"
ADAPTERS_DIR = LISTENER_DIR / "adapters"
EVENTS_DIR = LISTENER_DIR / "events"
CORE_DIR = LISTENER_DIR / "core"
CONTRACTS_DIR = REPO_ROOT / "contracts"
CONTRACT_IMPORT_PREFIX = "github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/"
CONTRACT_NAMES_MARKER = "// Contract names"
EVENT_NAMES_MARKER = "// Event names and signatures"
CONST_BLOCK_START = "const ("
CONST_BLOCK_TERMINATOR = ")"
CONST_BLOCK_CLOSE = "\n)\n"
PATHS = {
    "log_parser": ADAPTERS_DIR / "indexer/log_parser.go",
    "events": EVENTS_DIR / "events.go",
    "utils": ADAPTERS_DIR / "indexer/utils.go",
    "creator": CONTRACTS_DIR / "creator.go",
    "block_processor": CORE_DIR / "services/block_processor.go",
}


def ensure_import(content: str, import_path: str) -> str:
    line = f'\t"{import_path}"'
    if line in content:
        return content

    match = re.search(r"import \(\n([\s\S]*?)\n\)", content)
    if not match:
        raise ValueError("Could not locate import block")

    start, end = match.start(1), match.end(1)
    body = content[start:end]
    lines = body.splitlines()

    insert_after = -1
    for idx, existing in enumerate(lines):
        if CONTRACT_IMPORT_PREFIX in existing:
            insert_after = idx

    if insert_after >= 0:
        lines.insert(insert_after + 1, line)
    else:
        lines.append(line)

    new_body = "\n".join(lines)
    return content[:start] + new_body + content[end:]


def contract_constant_exists(content: str, spec: EventSpec) -> bool:
    pattern = rf"Contract{spec.contract_name}"
    return re.search(pattern, content) is not None


def update_log_parser(spec: EventSpec, path: Path, ctx: UpdateContext) -> bool:
    if ctx.contract_exists:
        return False
    original = path.read_text(encoding="utf-8")
    updated = ensure_import(
        original,
        f"{CONTRACT_IMPORT_PREFIX}{spec.binding_package}",
    )

    contracts_struct_pattern = r"(type Contracts struct {\n)([\s\S]*?)(\n})"
    match = re.search(contracts_struct_pattern, updated)
    if not match:
        raise ValueError("Contracts struct not found in log_parser.go")
    struct_body = match.group(2)
    struct_lines = struct_body.splitlines()
    new_field = f"\t{spec.contract_name}\t   *{spec.binding_package}.{spec.binding_type}"
    if new_field not in struct_lines:
        struct_lines.append(new_field)
        struct_body = "\n".join(struct_lines)
        updated = (
            updated[: match.start(2)]
            + struct_body
            + updated[match.end(2) :]
        )

    creation_block = (
        f"\n\t{spec.contract_camel}, err := contracts.Create{spec.contract_name}(config.CommitChain.{spec.contract_name}, client)\n"
        "\tif err != nil {\n"
        "\t\treturn nil, cErr.WithStack(err)\n"
        "\t}\n"
    )
    if f"contracts.Create{spec.contract_name}" not in updated:
        anchor = "\n\treturn &Contracts{"
        if anchor not in updated:
            raise ValueError("newContracts return anchor not found")
        updated = updated.replace(anchor, f"{creation_block}{anchor}", 1)

    return_pattern = r"(\treturn &Contracts{\n)([\s\S]*?)(\n\t}, nil)"
    match = re.search(return_pattern, updated)
    if not match:
        raise ValueError("Contracts return block not found")
    return_body = match.group(2)
    return_line = f"\n\t\t{spec.contract_name}: \t{spec.contract_camel},"
    if return_line not in return_body:
        return_body = return_body + return_line
        updated = (
            updated[: match.start(2)]
            + return_body
            + updated[match.end(2) :]
        )

    if updated == original:
        return False

    path.write_text(updated, encoding="utf-8")
    return True


def insert_contract_constant(content: str, spec: EventSpec) -> str:
    if contract_constant_exists(content, spec):
        return content

    start = content.find(CONTRACT_NAMES_MARKER)
    if start == -1:
        raise ValueError("Contract constants marker not found")
    const_start = content.find(CONST_BLOCK_START, start)
    if const_start == -1:
        raise ValueError("Contract constants block not found")
    const_end = content.find(CONST_BLOCK_CLOSE, const_start)
    if const_end == -1:
        raise ValueError("Contract constants block end not found")

    line = f"\n\tContract{spec.contract_name} = \t\"{spec.contract_name}\""
    return content[:const_end] + line + content[const_end:]


def insert_event_constants(content: str, spec: EventSpec) -> str:
    event_pattern = rf"\t{spec.event_name}\\s*=\\s*\"{spec.event_name}\""
    if re.search(event_pattern, content):
        return content

    lines = content.splitlines()
    try:
        header_idx = next(i for i, line in enumerate(lines) if EVENT_NAMES_MARKER in line)
    except StopIteration as exc:
        raise ValueError("Event constants marker not found") from exc

    try:
        const_start_idx = next(
            i for i in range(header_idx, len(lines)) if lines[i].startswith(CONST_BLOCK_START)
        )
    except StopIteration as exc:
        raise ValueError("Event constants block not found") from exc

    try:
        const_end_idx = next(
            i for i in range(const_start_idx, len(lines)) if lines[i] == CONST_BLOCK_TERMINATOR
        )
    except StopIteration as exc:
        raise ValueError("Event constants block end not found") from exc

    comment_line = f"\t// {spec.contract_name} events"
    event_line = f"\t{spec.event_name}\t\t\t  = \"{spec.event_name}\""
    signature_line = f"\t{spec.event_signature_const}\t\t  = \"{spec.event_signature}\""

    block_slice = slice(const_start_idx + 1, const_end_idx)
    block_lines = lines[block_slice]

    if comment_line in block_lines:
        comment_idx = block_lines.index(comment_line) + const_start_idx + 1
        insert_idx = comment_idx + 1
        while (
            insert_idx < const_end_idx
            and lines[insert_idx].startswith("\t")
            and not lines[insert_idx].startswith("\t//")
            and lines[insert_idx].strip() != ""
        ):
            insert_idx += 1
        lines.insert(insert_idx, signature_line)
        lines.insert(insert_idx, event_line)
    else:
        insert_idx = const_end_idx
        if lines[insert_idx - 1] != "":
            lines.insert(insert_idx, "")
            insert_idx += 1
        lines.insert(insert_idx, comment_line)
        insert_idx += 1
        lines.insert(insert_idx, event_line)
        insert_idx += 1
        lines.insert(insert_idx, signature_line)

    return "\n".join(lines)


def update_events(spec: EventSpec, path: Path, ctx: UpdateContext) -> bool:
    original = path.read_text(encoding="utf-8")
    updated = original

    if not ctx.contract_exists:
        updated = insert_contract_constant(updated, spec)

    updated = insert_event_constants(updated, spec)

    if updated == original:
        return False

    path.write_text(updated, encoding="utf-8")
    return True


def insert_parser_map_entry(content: str, spec: EventSpec) -> str:
    entry_tag = f"common.HexToAddress(cfg.CommitChain.{spec.contract_name})"
    if entry_tag in content:
        return content

    pattern = r"(return map\[common.Address]ContractParsers{\n)([\s\S]*?)(\n\t}\n})"
    match = re.search(pattern, content)
    if not match:
        raise ValueError("Parser registry map not found")

    body = match.group(2)
    entry = parser_map_entry(spec)
    if entry not in body:
        body = body + entry
        content = content[: match.start(2)] + body + content[match.end(2) :]

    return content


def insert_event_parser(content: str, spec: EventSpec) -> str:
    signature = f"func {spec.parser_function_name}(c *Contracts) []EventParser {{"
    event_entry = parser_event_entry(spec)

    signature_idx = content.find(signature)
    if signature_idx != -1:
        function_end = content.find("\n}\n", signature_idx)
        if function_end == -1:
            raise ValueError("Core parser function closing brace not found")

        function_body = content[signature_idx:function_end]
        if f"Parse{spec.event_name}" in function_body:
            return content

        return_idx = content.find("return []EventParser{", signature_idx)
        if return_idx == -1 or return_idx > function_end:
            raise ValueError("Parser return block not found")

        closing_idx = content.find("\n\t}\n", return_idx)
        if closing_idx == -1 or closing_idx > function_end:
            raise ValueError("Parser return closing not found")

        insertion_point = closing_idx
        return content[:insertion_point] + event_entry + content[insertion_point:]

    return content + parser_function(spec)


def update_utils(spec: EventSpec, path: Path, ctx: UpdateContext) -> bool:
    original = path.read_text(encoding="utf-8")
    updated = original

    if not ctx.contract_exists:
        updated = insert_parser_map_entry(updated, spec)

    updated = insert_event_parser(updated, spec)

    if updated == original:
        return False

    path.write_text(updated, encoding="utf-8")
    return True


def update_creator(spec: EventSpec, path: Path, ctx: UpdateContext) -> bool:
    if ctx.contract_exists:
        return False

    original = path.read_text(encoding="utf-8")
    updated = ensure_import(
        original,
        f"{CONTRACT_IMPORT_PREFIX}{spec.binding_package}",
    )

    if f"Create{spec.contract_name}" not in updated:
        anchor = "\nfunc NewDeploymentProxyRegistryV1"
        if anchor not in updated:
            raise ValueError("DeploymentProxyRegistry anchor not found")
        updated = updated.replace(anchor, creator_function(spec) + "\n" + anchor, 1)

    if updated == original:
        return False

    path.write_text(updated, encoding="utf-8")
    return True


def update_block_processor(spec: EventSpec, path: Path, ctx: UpdateContext) -> bool:
    if ctx.contract_exists:
        return False

    original = path.read_text(encoding="utf-8")
    case_line = f"\tcase events.Contract{spec.contract_name}:"
    if case_line in original:
        return False

    header = "func (bp *BlockProcessor) processContractLog"
    header_idx = original.find(header)
    if header_idx == -1:
        raise ValueError("processContractLog function not found")

    switch_idx = original.find("switch log.ContractName {", header_idx)
    if switch_idx == -1:
        raise ValueError("processContractLog switch block not found")

    default_idx = original.find("\tdefault:", switch_idx)
    if default_idx == -1:
        raise ValueError("processContractLog default clause not found")

    prefix = "" if original[:default_idx].endswith("\n") else "\n"
    insertion = prefix + block_processor_case(spec)
    updated = original[:default_idx] + insertion + original[default_idx:]

    if original == updated:
        return False

    path.write_text(updated, encoding="utf-8")
    return True


def ensure_core_case(content: str, spec: EventSpec) -> tuple[str, bool]:
    case_line = f"\tcase events.{spec.event_name}:"
    if case_line in content:
        return content, False

    header = f"func (bp *BlockProcessor) {spec.core_processor_function}"
    header_idx = content.find(header)
    if header_idx == -1:
        raise ValueError(f"{spec.core_processor_function} not found in core file")

    switch_idx = content.find("switch log.EventName {", header_idx)
    if switch_idx == -1:
        raise ValueError("switch block not found in core file")

    default_idx = content.find("\tdefault:", switch_idx)
    if default_idx == -1:
        raise ValueError("default clause not found in core file switch")

    updated = content[:default_idx] + core_switch_case(spec) + content[default_idx:]
    return updated, True


def ensure_core_handler(content: str, spec: EventSpec) -> tuple[str, bool]:
    signature = f"func (bp *BlockProcessor) {spec.event_handler_function}(ctx context.Context, log ContractLog) error {{"
    if signature in content:
        return content, False

    return content + core_event_handler(spec), True


def update_core(spec: EventSpec, directory: Path, ctx: UpdateContext) -> bool:
    target = directory / spec.core_file_name
    if not target.exists():
        directory.mkdir(parents=True, exist_ok=True)
        target.write_text(core_file(spec), encoding="utf-8")
        return True

    original = target.read_text(encoding="utf-8")
    updated, case_changed = ensure_core_case(original, spec)
    updated, handler_changed = ensure_core_handler(updated, spec)

    if not case_changed and not handler_changed:
        return False

    target.write_text(updated, encoding="utf-8")
    return True


def apply_update(path: Path, updater: Callable[[EventSpec, Path, UpdateContext], bool], spec: EventSpec, ctx: UpdateContext) -> bool:
    return updater(spec, path, ctx)


def main() -> int:
    print("\nBefore running this script for the first time, review the README docs. \nIt covers the required contract binding, canonical signature format and naming rules.\n")
    contract_name = prompt_non_empty("Contract name (MyContract): ")
    binding_name = prompt_non_empty("Binding package name (e.g., MyContractV1): ")
    event_name = prompt_non_empty("Event name (e.g., EventCompleted): ")
    event_signature = prompt_non_empty("Full event signature (e.g., EventCompleted(bytes32,uint256,(uint256,uint8),uint256): ")

    spec = EventSpec(
        contract_name=contract_name,
        binding_name=binding_name,
        event_name=event_name,
        event_signature=event_signature,
    )

    missing = [str(p) for p in PATHS.values() if not p.exists()]
    if missing:
        print("Error: required files not found:")
        for item in missing:
            print(f" - {item}")
        return 1

    events_content = PATHS["events"].read_text(encoding="utf-8")
    contract_exists = contract_constant_exists(events_content, spec)
    if contract_exists:
        print(f"\nContract {spec.contract_name} is already mapped. Skipping contract setup steps.\n")
    else:
        print(f"\nContract {spec.contract_name} is not mapped yet. Full contract wiring will be added.\n")

    ctx = UpdateContext(contract_exists=contract_exists)

    updates = {
        "log_parser": update_log_parser,
        "events": update_events,
        "utils": update_utils,
        "creator": update_creator,
        "block_processor": update_block_processor,
    }

    changes = []
    for key, updater in updates.items():
        changed = apply_update(PATHS[key], updater, spec, ctx)
        status = "updated" if changed else "unchanged"
        changes.append((PATHS[key], status))

    core_dir = CORE_DIR
    core_relative_path = core_dir.relative_to(REPO_ROOT)
    core_path = core_dir / spec.core_file_name
    core_existed_before = core_path.exists()
    core_changed = update_core(spec, core_dir, ctx)

    for path, status in changes:
        print(f"[{status}] {path.relative_to(REPO_ROOT)}")

    if core_changed:
        note = "created" if not core_existed_before else "updated"
        print(f"[{note}] {core_relative_path / spec.core_file_name}")
    else:
        print(f"[unchanged] {core_relative_path / spec.core_file_name}")

    print("\nUpdates applied successfully.")

    return 0


if __name__ == "__main__":
    sys.exit(main())

