package indexer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/config"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/AuditManagerV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/DvpTeleport"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/EnygmaTeleport"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/EnygmaTokenManagerV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/ParticipantCoreV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/Proofs"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/TeleportV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/TokenCoreV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/TokenFreezeManagerV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

// Ensure LogParser implements core.LogParser at compile time
var _ core.LogParser = (*LogParser)(nil)

// EventParser defines a function type for parsing events from logs
type EventParserFunction func(log ethTypes.Log) (any, error)

// EventParser represents a single event parser associated with a contract event by its signature hash
type EventParser struct {
	eventName          string
	eventSignatureHash common.Hash
	parser             EventParserFunction
}

// ContractParsers holds the contract name and its associated event parsers
type ContractParsers struct {
	Name    string
	Parsers []EventParser
}

// Contracts holds all the contract wrappers used by LogParser
type Contracts struct {
	TokenCore          *TokenCoreV1.TokenCoreV1
	TokenFreezeManager *TokenFreezeManagerV1.TokenFreezeManagerV1
	EnygmaTokenManager *EnygmaTokenManagerV1.EnygmaTokenManagerV1
	ParticipantCore    *ParticipantCoreV1.ParticipantCoreV1
	AuditManager       *AuditManagerV1.AuditManagerV1
	Teleport           *TeleportV1.TeleportV1
	Proofs             *Proofs.Proofs
	EnygmaTeleport     *EnygmaTeleport.EnygmaTeleport
	DvpTeleport        *DvpTeleport.DvpTeleport
}

// LogParser implements the core.LogParser interface for parsing logs from smart contracts
type LogParser struct {
	Client                   EthereumClient
	Config                   *config.Config
	AddressToContractParsers map[common.Address]ContractParsers
	Contracts                *Contracts
	provider                 core.Provider
	log                      logger.Logger
}

// NewLogParser creates a new LogParser instance with injected client, contracts, provider, config, and logger.
// It validates all event signatures against contract ABIs at startup to prevent silent event drops.
func NewLogParser(
	client EthereumClient,
	contracts *Contracts,
	provider core.Provider,
	config *config.Config,
	log logger.Logger,
) (*LogParser, error) {
	if err := ValidateEventSignatures(log); err != nil {
		return nil, fmt.Errorf("event signature validation failed: %w", err)
	}

	addressToContractParsers := NewParserRegistry(config, contracts)
	return &LogParser{
		Client:                   client,
		Config:                   config,
		Contracts:                contracts,
		AddressToContractParsers: addressToContractParsers,
		provider:                 provider,
		log:                      log,
	}, nil
}

// NewContracts instantiates all contract wrappers for LogParser.
// These are event-only bindings used for log parsing — no RPC calls.
func NewContracts() *Contracts {
	return &Contracts{
		TokenCore:          contracts.CreateTokenCore(),
		TokenFreezeManager: contracts.CreateTokenFreezeManager(),
		EnygmaTokenManager: contracts.CreateEnygmaTokenManager(),
		Teleport:           contracts.CreateTeleport(),
		ParticipantCore:    contracts.CreateParticipantCore(),
		AuditManager:       contracts.CreateAuditManager(),
		Proofs:             contracts.CreateProofs(),
		EnygmaTeleport:     contracts.CreateEnygmaTeleport(),
		DvpTeleport:        contracts.CreateDvpTeleport(),
	}
}

// ParseLogs parses logs from the Private Network Hub within a specified block range
func (l *LogParser) ParseLogs(ctx context.Context, fromBlock, toBlock *big.Int) ([]core.ContractLog, error) {
	logs, err := l.fetchContractLogs(ctx, fromBlock, toBlock)
	if err != nil {
		return nil, withstack.Wrap(err)
	}

	l.log.Debug("Filtered logs", "count", fmt.Sprint(len(logs)))

	blockTimestamps, err := l.fetchBlockTimestamps(ctx, logs)
	if err != nil {
		return nil, err
	}

	parsedLogs := make([]core.ContractLog, 0, len(logs))
	for _, log := range logs {
		parsedLog, err := l.parseSingleLog(log, blockTimestamps[log.BlockNumber])
		if err != nil {
			continue
		}
		parsedLogs = append(parsedLogs, *parsedLog)
	}
	return parsedLogs, nil
}

// fetchBlockTimestamps fetches timestamps for all unique block numbers present in logs.
func (l *LogParser) fetchBlockTimestamps(ctx context.Context, logs []ethTypes.Log) (map[uint64]time.Time, error) {
	uniqueBlocks := make(map[uint64]struct{})
	for _, log := range logs {
		uniqueBlocks[log.BlockNumber] = struct{}{}
	}

	timestamps := make(map[uint64]time.Time, len(uniqueBlocks))
	for blockNumber := range uniqueBlocks {
		block, err := l.provider.GetBlockByNumber(ctx, blockNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch block to retrieve timestamp (blockNumber=%d): %w", blockNumber, err)
		}
		blockTime := block.Time()
		if blockTime > math.MaxInt64 {
			return nil, fmt.Errorf("block timestamp overflows int64 (blockNumber=%d)", blockNumber)
		}
		timestamps[blockNumber] = time.Unix(int64(blockTime), 0).UTC()
	}
	return timestamps, nil
}

// FetchContractLogs fetches logs from monitored contracts within a specified block range
func (l *LogParser) fetchContractLogs(ctx context.Context, fromBlock, toBlock *big.Int) ([]ethTypes.Log, error) {
	contractsAddresses := MapKeys(l.AddressToContractParsers)
	query := ethereum.FilterQuery{
		Addresses: contractsAddresses,
		FromBlock: fromBlock,
		ToBlock:   toBlock,
	}
	logs, err := l.Client.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to filter logs from block %s to %s: %w",
			fromBlock.String(),
			toBlock.String(),
			err,
		)
	}
	return logs, nil
}

// parseSingleLog parses a single log entry and returns a ContractLog instance
func (l *LogParser) parseSingleLog(
	log ethTypes.Log,
	blockTimestamp time.Time,
) (*core.ContractLog, error) {
	// Some logs may not have topics if they are anonymous events
	if len(log.Topics) == 0 {
		malformedErr := fmt.Errorf(
			"%w (logAddress=%s, blockNumber=%d, txHash=%s)",
			core.ErrLogNoTopics,
			log.Address.String(),
			log.BlockNumber,
			log.TxHash.Hex(),
		)
		l.log.Error("Log has no topics",
			"logAddress", log.Address.String(),
			"blockNumber", log.BlockNumber,
			"txHash", log.TxHash.Hex())
		return nil, malformedErr
	}

	eventSignatureHash := log.Topics[0]
	contractName, eventName, parser, err := l.getEventParser(log.Address, eventSignatureHash)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrContractNotFound):
			l.log.Error("Contract not found at ParserRegistry",
				"eventSignature", eventSignatureHash.Hex(),
				"blockNumber", log.BlockNumber,
				"txHash", log.TxHash.Hex())
		case errors.Is(err, core.ErrUnwatchedEvent):
			l.log.Debug("Unwatched event found",
				"contract", contractName,
				"eventSignatureHash", eventSignatureHash.Hex())
		default:
			l.log.Error("Failed to get event parser",
				"eventSignature", eventSignatureHash.Hex(),
				"blockNumber", log.BlockNumber,
				"txHash", log.TxHash.Hex())
		}
		return nil, withstack.Wrap(err)
	}

	parsedEventData, err := parser(log)
	if err != nil {
		parseErr := fmt.Errorf(
			"%w (contract=%s event=%s): %w",
			core.ErrFailedToParseEvent,
			contractName,
			eventName,
			err,
		)
		l.log.Error("Failed to parse event",
			"contractName", contractName,
			"eventName", eventName)
		return nil, withstack.Wrap(parseErr)
	}

	rawEventData, err := json.Marshal(parsedEventData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data (contract=%s event=%s): %w", contractName, eventName, err)
	}

	return &core.ContractLog{
		ContractName:    contractName,
		ContractAddress: log.Address.String(),
		EventName:       eventName,
		RawEventData:    rawEventData,
		BlockNumber:     log.BlockNumber,
		TransactionHash: log.TxHash.String(),
		LogIndex:        uint64(log.Index),
		TxIndex:         log.TxIndex,
		BlockTimestamp:  blockTimestamp,
	}, nil
}

// getEventParser retrieves the event parser for a given contract address and event signature hash.
func (l *LogParser) getEventParser(
	address common.Address,
	eventSignatureHash common.Hash,
) (string, string, EventParserFunction, error) {
	contractParsers, ok := l.AddressToContractParsers[address]
	if !ok {
		return "", "", nil, fmt.Errorf(
			"%w (address=%s eventSignature=%s)",
			core.ErrContractNotFound,
			address.Hex(),
			eventSignatureHash.Hex(),
		)
	}
	for _, eventParser := range contractParsers.Parsers {
		if eventParser.eventSignatureHash == eventSignatureHash {
			return contractParsers.Name, eventParser.eventName, eventParser.parser, nil
		}
	}
	return contractParsers.Name, "", nil, fmt.Errorf(
		"%w (contract=%s eventSignatureHash=%s)",
		core.ErrUnwatchedEvent,
		contractParsers.Name,
		eventSignatureHash.Hex(),
	)
}
