package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/flagger/adapters/repositories"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/flagger/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/flagger/infrastructure"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
)

const (
	DefaultBatchSize     = 100
	MaxBatchSize         = 800
	DefaultCheckInterval = 1 * time.Second
	MaxCheckInterval     = 6 * time.Second
)

func Run(envPath string) error {
	// Setup infrastructure layer - Loading config and initializing DB connection
	infra, err := infrastructure.SetupInfrastructure(envPath)
	if err != nil {
		return fmt.Errorf("failed to setup infrastructure: %w", err)
	}
	dbClient := infra.DBClient
	conf := infra.Config

	defer func() {
		if cleanupErr := infra.Cleanup(); cleanupErr != nil {
			logger.Error("Error cleaning up infrastructure", "error", cleanupErr)
		}
	}()

	// Create repositories
	headerProofRepo := repositories.NewHeaderProofRepository(dbClient)
	headerFlagEventRepo := repositories.NewHeaderFlagEventRepository(dbClient)
	balanceRepo := repositories.NewBalanceRepository(infra.DBClient)
	txRepo := repositories.NewTransactionRepository(infra.DBClient)

	// Create logger
	log := logger.NewLogger()

	// Create core handlers
	batchSize := conf.TransactionProcessor.BatchSize
	if batchSize <= 0 {
		batchSize = DefaultBatchSize
	}
	if batchSize > MaxBatchSize {
		batchSize = MaxBatchSize
	}

	transactionProcessor := core.NewTransactionProcessor(
		balanceRepo,
		txRepo,
		log,
		batchSize,
	)

	// Parse expiration period from config
	expirationPeriod, err := time.ParseDuration(conf.PrivateHub.HeaderProofExpirationPeriod)
	if err != nil {
		return fmt.Errorf("failed to parse expiration period: %w", err)
	}

	// Create core HeaderLivelinessProcessor
	headerLivelinessProcessor := core.NewHeaderLivelinessProcessor(
		headerProofRepo,
		headerFlagEventRepo,
		expirationPeriod,
		log,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create errgroup for managing concurrent goroutines
	eg, ctx := errgroup.WithContext(ctx)

	// Start header liveliness checking with ticker
	logger.Info("Starting the header liveliness processor")
	checkInterval := 1 * time.Minute
	eg.Go(func() error {
		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()
		return headerLivelinessProcessor.Run(ctx, ticker.C)
	})

	// Start transaction processor
	logger.Info("Starting the transaction processor")
	transactionCheckInterval := DefaultCheckInterval
	if conf.TransactionProcessor.CheckInterval != "" {
		parsed, err := time.ParseDuration(conf.TransactionProcessor.CheckInterval)
		if err != nil {
			return fmt.Errorf("failed to parse transaction check interval: %w", err)
		}
		transactionCheckInterval = parsed
	}
	if transactionCheckInterval > MaxCheckInterval {
		transactionCheckInterval = MaxCheckInterval
	}
	if transactionCheckInterval <= DefaultCheckInterval {
		transactionCheckInterval = DefaultCheckInterval
	}
	eg.Go(func() error {
		ticker := time.NewTicker(transactionCheckInterval)
		defer ticker.Stop()
		return transactionProcessor.Run(ctx, ticker.C)
	})

	// Start header proof purge processor (opt-in via PNH_HEADER_PROOF_PURGE_PERIOD)
	if conf.PrivateHub.HeaderProofPurgePeriod != "" {
		purgePeriod, err := time.ParseDuration(conf.PrivateHub.HeaderProofPurgePeriod)
		if err != nil {
			return fmt.Errorf("failed to parse purge period: %w", err)
		}
		purgeProcessor := core.NewHeaderProofPurgeProcessor(headerProofRepo, purgePeriod, log)
		logger.Info("Starting header proof purge processor", "retention", purgePeriod)
		eg.Go(func() error {
			ticker := time.NewTicker(24 * time.Hour)
			defer ticker.Stop()
			return purgeProcessor.Run(ctx, ticker.C)
		})
	}

	// Setup signal handling
	shutdownChan := make(chan os.Signal, 1)
	defer close(shutdownChan)

	// Use os.Interrupt and syscall.SIGTERM for Notify
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Start signal handler goroutine
	eg.Go(func() error {
		select {
		case sig := <-shutdownChan:
			//nolint:contextcheck // intentionally uses package-level logger; ctx is being cancelled here
			logger.Info(
				"Shutdown signal received",
				"signal",
				sig.String(),
			)
			cancel()
			return nil
		case <-ctx.Done():
			return nil
		}
	})

	// Wait for all goroutines to finish or first error
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("Error running services", "error", err)
		return fmt.Errorf("error running services: %w", err)
	}

	logger.Info("Shutdown complete")
	return nil
}
