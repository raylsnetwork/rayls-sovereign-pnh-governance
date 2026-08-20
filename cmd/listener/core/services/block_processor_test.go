package services

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/mocks"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/testutil"
)

func TestRun(t *testing.T) {
	t.Run("context cancellation stops the loop", func(t *testing.T) {
		// Setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		ticker := make(chan time.Time)

		bp := &BlockProcessor{
			log: &testutil.StubLogger{},
		}

		// Execute in goroutine
		done := make(chan error, 1)
		go func() {
			done <- bp.Run(ctx, ticker)
		}()

		// Cancel immediately
		cancel()

		// Assert - should return nil within timeout
		select {
		case err := <-done:
			assert.NoError(t, err, "Run should return nil on context cancellation")
		case <-time.After(1 * time.Second):
			t.Fatal("Run did not return within timeout after context cancellation")
		}
	})

	t.Run("ticker triggers Start call", func(t *testing.T) {
		// Setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ticker := make(chan time.Time)

		// Track calls to verify Start was executed
		callCount := 0
		mockBlockRepo := mocks.NewMockBlockRepository(ctrl)
		stubConfigProvider := testutil.NewStubConfigProvider()

		mockBlockRepo.EXPECT().
			GetLatestProcessedBlock(gomock.Any()).
			DoAndReturn(func(ctx context.Context) (*big.Int, error) {
				callCount++
				// Return a block without error so getNextBlockToProcess succeeds
				return big.NewInt(100), nil
			}).
			AnyTimes()

		// Make processBlocks fail immediately after we know Start was called
		mockProvider := mocks.NewMockProvider(ctrl)
		mockProvider.EXPECT().
			GetLatestBlock(gomock.Any()).
			Return(nil, errors.New("stop here - we proved Start was called")).
			AnyTimes()

		bp := &BlockProcessor{
			blockRepo:      mockBlockRepo,
			configProvider: stubConfigProvider,
			provider:       mockProvider,
			log:            &testutil.StubLogger{},
		}

		// Execute in goroutine
		done := make(chan error, 1)
		go func() {
			done <- bp.Run(ctx, ticker)
		}()

		// Send a tick
		ticker <- time.Now()

		// Give it time to process
		time.Sleep(50 * time.Millisecond)

		// Cancel context to stop
		cancel()

		// Assert - should return nil
		select {
		case err := <-done:
			assert.NoError(t, err, "Run should return nil")
		case <-time.After(1 * time.Second):
			t.Fatal("Run did not return within timeout")
		}

		// Verify Start was called
		assert.Equal(t, 1, callCount, "Start should be called once after one tick")
	})
}

func TestNewBlockProcessor(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlockRepo := mocks.NewMockBlockRepository(ctrl)
	mockLogParser := mocks.NewMockLogParser(ctrl)
	logPub := NewLogPublisher(&stubContractMQ{}, &testutil.StubLogger{})
	mockLogger := &testutil.StubLogger{}
	stubConfigProvider := testutil.NewStubConfigProvider()
	mockProvider := mocks.NewMockProvider(ctrl)

	// Execute
	bp := NewBlockProcessor(
		mockBlockRepo,
		mockLogParser,
		logPub,
		mockLogger,
		stubConfigProvider,
		mockProvider,
	)

	// Assert
	assert.NotNil(t, bp, "BlockProcessor should not be nil")
	assert.Equal(t, mockBlockRepo, bp.blockRepo, "blockRepo should be set")
	assert.Equal(t, mockLogParser, bp.logParser, "logParser should be set")
	assert.Equal(t, logPub, bp.logPublisher, "logPublisher should be set")
	assert.Equal(t, mockLogger, bp.log, "logger should be set")
	assert.Equal(t, stubConfigProvider, bp.configProvider, "configProvider should be set")
	assert.Equal(t, mockProvider, bp.provider, "provider should be set")
}

func TestGetNextBlockToProcess(t *testing.T) {
	tests := []struct {
		name          string
		mockReturn    *big.Int
		mockError     error
		expectedBlock *big.Int
	}{
		{
			name:          "first run - no block in DB",
			mockReturn:    nil,
			mockError:     nil,
			expectedBlock: nil,
		},
		{
			name:          "DB returns error",
			mockReturn:    nil,
			mockError:     errors.New("db connection failed"),
			expectedBlock: nil,
		},
		{
			name:          "resume from block 100",
			mockReturn:    big.NewInt(100),
			mockError:     nil,
			expectedBlock: big.NewInt(101),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockBlockRepo := mocks.NewMockBlockRepository(ctrl)

			// Mock expectations
			mockBlockRepo.EXPECT().
				GetLatestProcessedBlock(ctx).
				Return(tt.mockReturn, tt.mockError)

			// Create BlockProcessor with minimal dependencies
			bp := &BlockProcessor{
				blockRepo: mockBlockRepo,
				log:       &testutil.StubLogger{},
			}

			// Execute
			result := bp.getNextBlockToProcess(ctx)

			// Assert
			if tt.expectedBlock == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedBlock.Int64(), result.Int64())
			}
		})
	}
}

func TestProcessBlocks(t *testing.T) {
	tests := []struct {
		name               string
		fromBlock          *big.Int
		batchSize          int64
		latestBlockNumber  *big.Int
		latestBlockErr     error
		parsedLogs         []core.ContractLog
		parseLogsErr       error
		expectedToBlock    *big.Int
		expectedTipReached bool
		expectedError      bool
	}{
		{
			name:               "normal batch - not at tip",
			fromBlock:          big.NewInt(100),
			batchSize:          20,
			latestBlockNumber:  big.NewInt(200),
			latestBlockErr:     nil,
			parsedLogs:         []core.ContractLog{{ContractName: "C1"}, {ContractName: "C2"}, {ContractName: "C3"}},
			parseLogsErr:       nil,
			expectedToBlock:    big.NewInt(120),
			expectedTipReached: false,
			expectedError:      false,
		},
		{
			name:               "reaches tip exactly",
			fromBlock:          big.NewInt(100),
			batchSize:          20,
			latestBlockNumber:  big.NewInt(120),
			latestBlockErr:     nil,
			parsedLogs:         []core.ContractLog{{ContractName: "C1"}, {ContractName: "C2"}},
			parseLogsErr:       nil,
			expectedToBlock:    big.NewInt(120),
			expectedTipReached: true,
			expectedError:      false,
		},
		{
			name:               "already at tip",
			fromBlock:          big.NewInt(100),
			batchSize:          20,
			latestBlockNumber:  big.NewInt(100),
			latestBlockErr:     nil,
			parsedLogs:         []core.ContractLog{{ContractName: "C1"}},
			parseLogsErr:       nil,
			expectedToBlock:    big.NewInt(100),
			expectedTipReached: true,
			expectedError:      false,
		},
		{
			name:               "GetLatestBlock returns error",
			fromBlock:          big.NewInt(100),
			batchSize:          20,
			latestBlockNumber:  nil,
			latestBlockErr:     errors.New("connection failed"),
			parsedLogs:         nil,
			parseLogsErr:       nil,
			expectedToBlock:    nil,
			expectedTipReached: false,
			expectedError:      true,
		},
		{
			name:               "ParseLogs returns error",
			fromBlock:          big.NewInt(100),
			batchSize:          20,
			latestBlockNumber:  big.NewInt(200),
			latestBlockErr:     nil,
			parsedLogs:         nil,
			parseLogsErr:       errors.New("log parsing failed"),
			expectedToBlock:    nil,
			expectedTipReached: false,
			expectedError:      true,
		},
		{
			name:               "ParseLogs returns empty logs",
			fromBlock:          big.NewInt(100),
			batchSize:          20,
			latestBlockNumber:  big.NewInt(200),
			latestBlockErr:     nil,
			parsedLogs:         []core.ContractLog{},
			parseLogsErr:       nil,
			expectedToBlock:    big.NewInt(120),
			expectedTipReached: false,
			expectedError:      false,
		},
		{
			name:               "large batch size caps at latest block",
			fromBlock:          big.NewInt(100),
			batchSize:          1000,
			latestBlockNumber:  big.NewInt(200),
			latestBlockErr:     nil,
			parsedLogs:         []core.ContractLog{{ContractName: "C1"}, {ContractName: "C2"}},
			parseLogsErr:       nil,
			expectedToBlock:    big.NewInt(200),
			expectedTipReached: true,
			expectedError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockProvider := mocks.NewMockProvider(ctrl)
			stubConfigProvider := testutil.NewStubConfigProvider()
			stubConfigProvider.BatchSize = tt.batchSize
			mockLogParser := mocks.NewMockLogParser(ctrl)

			// Setup GetLatestBlock expectation
			if tt.latestBlockErr != nil {
				mockProvider.EXPECT().
					GetLatestBlock(ctx).
					Return(nil, tt.latestBlockErr)
			} else if tt.latestBlockNumber != nil {
				mockProvider.EXPECT().
					GetLatestBlock(ctx).
					Return(ethTypes.NewBlockWithHeader(&ethTypes.Header{
						Number: tt.latestBlockNumber,
					}), nil)

				mockLogParser.EXPECT().
					ParseLogs(ctx, tt.fromBlock, gomock.Any()).
					Return(tt.parsedLogs, tt.parseLogsErr)
			}

			bp := &BlockProcessor{
				provider:       mockProvider,
				configProvider: stubConfigProvider,
				logParser:      mockLogParser,
				log:            &testutil.StubLogger{},
			}

			// Execute
			logs, tipReached, toBlock, err := bp.processBlocks(ctx, tt.fromBlock)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, toBlock)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, toBlock)
				assert.Equal(t, tt.expectedToBlock.Int64(), toBlock.Int64(), "toBlock should match expected")
				assert.Equal(t, tt.expectedTipReached, tipReached, "tipReached should match expected")
				assert.Equal(t, len(tt.parsedLogs), len(logs), "logs length should match")
			}
		})
	}
}

func TestStart(t *testing.T) {
	tests := []struct {
		name                    string
		latestProcessedBlock    *big.Int
		latestProcessedBlockErr error
		startingBlock           *big.Int
		expectedFromBlock       *big.Int
		latestChainBlock        *big.Int
		latestChainBlockErr     error
		batchSize               int64
		parsedLogs              []core.ContractLog
		parseLogsErr            error
		publishErr              error
		updateBlockErr          error
		expectedError           bool
	}{
		{
			name:                 "happy path - resume from existing block",
			latestProcessedBlock: big.NewInt(100),
			startingBlock:        nil,
			expectedFromBlock:    big.NewInt(101),
			latestChainBlock:     big.NewInt(200),
			batchSize:            20,
			parsedLogs:           []core.ContractLog{{ContractName: "C1"}},
			expectedError:        false,
		},
		{
			name:                 "happy path - first run (no block in DB)",
			latestProcessedBlock: nil,
			startingBlock:        big.NewInt(50),
			expectedFromBlock:    big.NewInt(50),
			latestChainBlock:     big.NewInt(200),
			batchSize:            20,
			parsedLogs:           []core.ContractLog{{ContractName: "C1"}},
			expectedError:        false,
		},
		{
			name:                 "processBlocks returns error - GetLatestBlock fails",
			latestProcessedBlock: big.NewInt(100),
			expectedFromBlock:    big.NewInt(101),
			latestChainBlockErr:  errors.New("chain connection failed"),
			expectedError:        true,
		},
		{
			name:                 "processBlocks returns error - ParseLogs fails",
			latestProcessedBlock: big.NewInt(100),
			expectedFromBlock:    big.NewInt(101),
			latestChainBlock:     big.NewInt(200),
			batchSize:            20,
			parseLogsErr:         errors.New("log parsing failed"),
			expectedError:        true,
		},
		{
			name:                 "Publish returns error",
			latestProcessedBlock: big.NewInt(100),
			expectedFromBlock:    big.NewInt(101),
			latestChainBlock:     big.NewInt(200),
			batchSize:            20,
			parsedLogs:           []core.ContractLog{{ContractName: "C1"}},
			publishErr:           errors.New("publish failed"),
			expectedError:        true,
		},
		{
			name:                 "UpdateLatestProcessedBlock returns error",
			latestProcessedBlock: big.NewInt(100),
			expectedFromBlock:    big.NewInt(101),
			latestChainBlock:     big.NewInt(200),
			batchSize:            20,
			parsedLogs:           []core.ContractLog{{ContractName: "C1"}},
			updateBlockErr:       errors.New("db update failed"),
			expectedError:        true,
		},
		{
			name:                 "tip reached - completes successfully",
			latestProcessedBlock: big.NewInt(100),
			expectedFromBlock:    big.NewInt(101),
			latestChainBlock:     big.NewInt(120),
			batchSize:            20,
			parsedLogs:           []core.ContractLog{{ContractName: "C1"}},
			expectedError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockBlockRepo := mocks.NewMockBlockRepository(ctrl)
			stubConfigProvider := testutil.NewStubConfigProvider()
			stubConfigProvider.StartingBlockNumber = tt.startingBlock
			stubConfigProvider.BatchSize = tt.batchSize
			mockProvider := mocks.NewMockProvider(ctrl)
			mockLogParser := mocks.NewMockLogParser(ctrl)
			mq := &stubContractMQ{failAt: 0, failErr: tt.publishErr}
			logPub := NewLogPublisher(mq, &testutil.StubLogger{})

			bp := &BlockProcessor{
				blockRepo:      mockBlockRepo,
				configProvider: stubConfigProvider,
				provider:       mockProvider,
				logParser:      mockLogParser,
				logPublisher:   logPub,
				log:            &testutil.StubLogger{},
			}

			if tt.latestProcessedBlockErr != nil {
				// Error in getNextBlockToProcess - early exit
				mockBlockRepo.EXPECT().
					GetLatestProcessedBlock(ctx).
					Return(nil, tt.latestProcessedBlockErr)
			} else {
				// GetLatestProcessedBlock succeeds
				mockBlockRepo.EXPECT().
					GetLatestProcessedBlock(ctx).
					Return(tt.latestProcessedBlock, nil)

				// Setup processBlocks dependencies
				if tt.latestChainBlockErr != nil {
					// GetLatestBlock fails
					mockProvider.EXPECT().
						GetLatestBlock(ctx).
						Return(nil, tt.latestChainBlockErr)
				} else if tt.latestChainBlock != nil {
					// GetLatestBlock succeeds
					latestBlock := ethTypes.NewBlockWithHeader(&ethTypes.Header{
						Number: tt.latestChainBlock,
					})
					mockProvider.EXPECT().
						GetLatestBlock(ctx).
						Return(latestBlock, nil)

					if tt.parseLogsErr != nil {
						// ParseLogs fails - use hardcoded expectedFromBlock to avoid replicating business logic
						mockLogParser.EXPECT().
							ParseLogs(ctx, tt.expectedFromBlock, gomock.Any()).
							Return(nil, tt.parseLogsErr)
					} else {
						// ParseLogs succeeds - use hardcoded expectedFromBlock to avoid replicating business logic
						mockLogParser.EXPECT().
							ParseLogs(ctx, tt.expectedFromBlock, gomock.Any()).
							Return(tt.parsedLogs, nil)

						// UpdateLatestProcessedBlock only called if Publish succeeds
						if tt.publishErr == nil {
							if tt.updateBlockErr != nil {
								mockBlockRepo.EXPECT().
									UpdateLatestProcessedBlock(ctx, gomock.Any()).
									Return(tt.updateBlockErr)
							} else {
								mockBlockRepo.EXPECT().
									UpdateLatestProcessedBlock(ctx, gomock.Any()).
									Return(nil)
							}
						}
					}
				}
			}

			// Execute
			err := bp.Start(ctx)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
