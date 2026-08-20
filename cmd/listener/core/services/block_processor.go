package services

import (
	"context"
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

// BlockProcessor is the core business logic component responsible for processing blockchain blocks
type BlockProcessor struct {
	blockRepo      core.BlockRepository
	logParser      core.LogParser
	logPublisher   *LogPublisher
	configProvider core.ConfigProvider
	provider       core.Provider
	log            logger.Logger
}

// NewBlockProcessor creates a new BlockProcessor instance
func NewBlockProcessor(
	blockRepo core.BlockRepository,
	logParser core.LogParser,
	logPublisher *LogPublisher,
	log logger.Logger,
	configProvider core.ConfigProvider,
	provider core.Provider,
) *BlockProcessor {
	return &BlockProcessor{
		blockRepo:      blockRepo,
		logParser:      logParser,
		logPublisher:   logPublisher,
		configProvider: configProvider,
		provider:       provider,
		log:            log,
	}
}

// Run processes blocks continuously
func (bp *BlockProcessor) Run(ctx context.Context, ticker <-chan time.Time) error {
	for {
		select {
		case <-ctx.Done():
			bp.log.Info("BlockProcessor stopped due to context cancellation")
			return nil
		case <-ticker:
			if err := bp.Start(ctx); err != nil {
				bp.log.Error("Error in block processing", "error", err)
			}
		}
	}
}

// Start processes blocks following the service logic
func (bp *BlockProcessor) Start(ctx context.Context) error {
	// Get latest processed block from the repository
	fromBlock := bp.getNextBlockToProcess(ctx)

	// if we just start, get the block from the config
	if fromBlock == nil {
		fromBlock = bp.configProvider.GetStartingBlockNumber()
	}

	// Process blocks
	logs, tipReached, lastProcessedBlock, err := bp.processBlocks(ctx, fromBlock)
	if err != nil {
		return withstack.Wrap(err)
	}

	// Publish logs to message queue
	if err := bp.logPublisher.Publish(ctx, logs); err != nil {
		return withstack.Wrap(err)
	}

	// Update the latest processed block
	if err := bp.blockRepo.UpdateLatestProcessedBlock(ctx, lastProcessedBlock); err != nil {
		return withstack.Wrap(err)
	}

	if tipReached {
		bp.log.Debug("Reached the tip of the blockchain")
	}

	return nil
}

// processBlocks processes a range of blocks
func (bp *BlockProcessor) processBlocks(
	ctx context.Context,
	fromBlockNumber *big.Int,
) ([]core.ContractLog, bool, *big.Int, error) {
	var logs []core.ContractLog

	// Get the latest block from the private network hub
	latestBlock, err := bp.provider.GetLatestBlock(ctx)
	if err != nil || latestBlock == nil {
		return logs, false, nil, withstack.Wrap(err)
	}

	// Calculate toBlock using batch size
	latestBlockNumber := latestBlock.Number()
	batchSize := bp.configProvider.GetBatchSize()

	toBlockNumber := new(big.Int).Add(fromBlockNumber, big.NewInt(batchSize))
	if toBlockNumber.Cmp(latestBlockNumber) > 0 {
		toBlockNumber = latestBlockNumber
	}

	bp.log.Debug("Processing blocks", "from", fromBlockNumber.String(), "to", toBlockNumber.String())

	// Parse logs from the block range
	logs, err = bp.logParser.ParseLogs(ctx, fromBlockNumber, toBlockNumber)
	if err != nil {
		return logs, false, nil, withstack.Wrap(err)
	}

	// Check if we reached the tip
	tipReached := toBlockNumber.Cmp(latestBlockNumber) == 0
	return logs, tipReached, toBlockNumber, nil
}

// getNextBlockToProcess determines the next block
func (bp *BlockProcessor) getNextBlockToProcess(ctx context.Context) *big.Int {
	// Try to get the latest processed block from the database
	latestBlock, err := bp.blockRepo.GetLatestProcessedBlock(ctx)
	if err != nil || latestBlock == nil {
		// If a block is not found in the DB (expected on first run), return nil so the caller can read it from the config
		return nil
	}

	// Return the next block after the latest processed one
	nextBlock := new(big.Int).Add(latestBlock, big.NewInt(1))
	return nextBlock
}
