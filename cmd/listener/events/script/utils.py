import re
from dataclasses import dataclass
from typing import Protocol


@dataclass(frozen=True)
class EventSpec:
    contract_name: str
    binding_name: str
    event_name: str
    event_signature: str

    @property
    def contract_camel(self) -> str:
        return to_camel_case(self.contract_name)

    @property
    def contract_snake(self) -> str:
        return to_snake_case(self.contract_name)

    @property
    def binding_package(self) -> str:
        return self.binding_name

    @property
    def binding_type(self) -> str:
        return self.binding_name

    @property
    def event_struct_type(self) -> str:
        return f"{self.binding_type}{self.event_name}"

    @property
    def event_signature_const(self) -> str:
        return f"{self.event_name}Sig"

    @property
    def parser_function_name(self) -> str:
        return f"{self.contract_camel}Parsers"

    @property
    def core_processor_function(self) -> str:
        return f"process{self.contract_name}Event"

    @property
    def event_handler_function(self) -> str:
        return f"process{self.event_name}"

    @property
    def core_file_name(self) -> str:
        return f"sc_{self.contract_snake}.go"


@dataclass(frozen=True)
class UpdateContext:
    contract_exists: bool


def prompt_non_empty(message: str) -> str:
    while True:
        value = input(message).strip()
        if value:
            return value
        print("Value required. Please try again.")


def to_words(value: str) -> list[str]:
    return re.findall(r"[A-Z]+(?=[A-Z][a-z]|[0-9]|$)|[A-Z]?[a-z]+|[0-9]+", value)


def to_camel_case(value: str) -> str:
    words = to_words(value)
    if not words:
        return value
    first, *rest = words
    tail = "".join(
        word if word.isupper() or word.isdigit() else word.capitalize() for word in rest
    )
    return first.lower() + tail


def to_snake_case(value: str) -> str:
    words = to_words(value)
    return "_".join(word.lower() for word in words)


class Spec(Protocol):
    contract_name: str
    contract_camel: str
    binding_name: str
    binding_package: str
    binding_type: str
    event_name: str
    event_signature: str
    event_signature_const: str
    event_struct_type: str
    parser_function_name: str
    core_processor_function: str
    event_handler_function: str


def parser_map_entry(spec: Spec) -> str:
    return (
        f"\n\t\tcommon.HexToAddress(cfg.CommitChain.{spec.contract_name}): {{\n"
        f"\t\t\tName:    events.Contract{spec.contract_name},\n"
        f"\t\t\tParsers: {spec.parser_function_name}(c),\n"
        "\t\t},"
    )


def parser_event_entry(spec: Spec) -> str:
    return (
        f"\n\t\t{{events.{spec.event_name}, Keccak256Hash(events.{spec.event_signature_const}), func(log types.Log) (any, error) {{\n"
        f"\t\t\tevent, err := c.{spec.contract_name}.Parse{spec.event_name}(log)\n"
        "\t\t\tif err != nil {\n"
        f"\t\t\t\treturn nil, fmt.Errorf(\"failed to parse {spec.event_name} event: %w\", err)\n"
        "\t\t\t}\n"
        "\t\t\treturn event, nil\n"
        "\t\t}},"
    )


def parser_function(spec: Spec) -> str:
    return (
        f"\n// {spec.contract_camel}Parsers returns the event parsers for the {spec.contract_name} contract\n"
        f"func {spec.parser_function_name}(c *Contracts) []EventParser {{\n"
        "\treturn []EventParser{\n"
        f"\t\t{{events.{spec.event_name}, Keccak256Hash(events.{spec.event_signature_const}), func(log types.Log) (any, error) {{\n"
        f"\t\t\tevent, err := c.{spec.contract_name}.Parse{spec.event_name}(log)\n"
        "\t\t\tif err != nil {\n"
        f"\t\t\t\treturn nil, fmt.Errorf(\"failed to parse {spec.event_name} event: %w\", err)\n"
        "\t\t\t}\n"
        "\t\t\treturn event, nil\n"
        "\t\t}},"
        "\t}\n"
        "}\n"
    )


def creator_function(spec: Spec) -> str:
    return (
        f"\nfunc Create{spec.contract_name}(address string, client *ethclient.Client) (*{spec.binding_package}.{spec.binding_type}, error) {{\n"
        "\tcontractAddress := common.HexToAddress(address)\n"
        "\tif err := infrastructure.CodeAt(client, contractAddress); err != nil {\n"
        "\t\treturn nil, cErr.WithStack(err)\n"
        "\t}\n\n"
        f"\tinstance, err := {spec.binding_package}.New{spec.binding_type}(contractAddress, client)\n"
        "\tif err != nil {\n"
        f"\t\treturn nil, fmt.Errorf(\"failed to create {spec.binding_type} instance for address %s: %w\", address, err)\n"
        "\t}\n"
        "\treturn instance, nil\n"
        "}"
    )


def block_processor_case(spec: Spec) -> str:
    return (
        f"\tcase events.Contract{spec.contract_name}:\n"
        f"\t\tif err := bp.{spec.core_processor_function}(ctx, log); err != nil {{\n"
        "\t\t\treturn cErr.WithStack(err)\n"
        "\t\t}\n"
    )


def core_switch_case(spec: Spec) -> str:
    return (
        f"\tcase events.{spec.event_name}:\n"
        f"\t\tif err := bp.{spec.event_handler_function}(ctx, log); err != nil {{\n"
        "\t\t\treturn cErr.WithStack(err)\n"
        "\t\t}\n"
    )


def core_event_handler(spec: Spec) -> str:
    return (
        f"\n// {spec.event_handler_function} processes {spec.event_name} events\n"
        f"func (bp *BlockProcessor) {spec.event_handler_function}(ctx context.Context, log ContractLog) error {{\n"
        f"\tevent, ok := log.EventData.(*{spec.binding_package}.{spec.event_struct_type})\n"
        "\tif !ok {\n"
        f"\t\treturn fmt.Errorf(\"unexpected event data type for {spec.event_name}: %T\", log.EventData)\n"
        "\t}\n\n"
        f"\tbp.logger.Info(\"{spec.event_name} event found\")\n\n"
        "\t// TODO: Persist the data here and make sure to also add meaningful entries to the log if applicable\n\n"
        f"\tbp.logger.Info(\"{spec.event_name} event processed\")\n"
        "\treturn nil\n"
        "}\n"
    )


def core_file(spec: Spec) -> str:
    return (
        "package core\n\n"
        "import (\n"
        "\t\"context\"\n"
        "\t\"fmt\"\n\n"
        "\t\"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/events\"\n"
        f"\t\"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/contracts/{spec.binding_package}\"\n\n"
        "\tcErr \"github.com/cockroachdb/errors\"\n"
        ")\n\n"
        f"// {spec.core_processor_function} processes {spec.contract_name} contract logs\n"
        f"func (bp *BlockProcessor) {spec.core_processor_function}(ctx context.Context, log ContractLog) error {{\n"
        "\tswitch log.EventName {\n"
        f"\tcase events.{spec.event_name}:\n"
        f"\t\tif err := bp.{spec.event_handler_function}(ctx, log); err != nil {{\n"
        "\t\t\treturn cErr.WithStack(err)\n"
        "\t\t}\n"
        f"\tdefault:\n\t\tbp.logger.Debug(\"No handler for {spec.contract_name} event\", \"event\", log.EventName)\n"
        "\t}\n"
        "\treturn nil\n"
        "}\n\n"
        f"// {spec.event_handler_function} processes {spec.event_name} events\n"
        f"func (bp *BlockProcessor) {spec.event_handler_function}(ctx context.Context, log ContractLog) error {{\n"
        f"\tevent, ok := log.EventData.(*{spec.binding_package}.{spec.event_struct_type})\n"
        "\tif !ok {\n"
        f"\t\treturn fmt.Errorf(\"unexpected event data type for {spec.event_name}: %T\", log.EventData)\n"
        "\t}\n\n"
        f"\tbp.logger.Info(\"{spec.event_name} event found\")\n\n"
        "\t// TODO: Persist the data here and make sure to also add meaningful entries to the log if applicable\n\n"
        f"\tbp.logger.Info(\"{spec.event_name} event processed\")\n"
        "\treturn nil\n"
        "}\n"
    )
