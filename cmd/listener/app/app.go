package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"golang.org/x/sync/errgroup"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/adapters/tokenregistry"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/adapters/config"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/adapters/crypto"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/adapters/indexer"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/adapters/repositories"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core/handlers"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core/services"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/events"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/infrastructure"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
)

func Run(envPath string) error {
	// Setup infrastructure layer
	infra, err := infrastructure.SetupInfrastructure(envPath)
	if err != nil {
		return fmt.Errorf("failed to setup infrastructure: %w", err)
	}
	defer func() {
		if cleanupErr := infra.Cleanup(); cleanupErr != nil {
			logger.Error("Error cleaning up infrastructure", "error", cleanupErr)
		}
	}()

	// Create repositories
	blockRepo := repositories.NewLastProcessedBlockRepository(infra.DBClient)
	txRepo := repositories.NewTransactionRepository(infra.DBClient)
	revertRepo := repositories.NewRevertDataTransactionRepository(infra.DBClient)
	headerProofRepo := repositories.NewHeaderProofEventRepository(infra.DBClient)
	tokenFreezeRepo := repositories.NewTokenFreezeRepository(infra.DBClient)

	// Create other necessary adapters
	log := logger.NewLogger()
	provider := indexer.NewProvider(infra.EthClient, infra.Config)
	configProvider := config.NewConfigProvider(infra.Config)
	decryptor := crypto.NewDecryptor(infra.EthClient, infra.Config, log)
	contracts := indexer.NewContracts()
	logParser, err := indexer.NewLogParser(infra.EthClient, contracts, provider, infra.Config, log)
	if err != nil {
		return fmt.Errorf("failed to create log parser: %w", err)
	}

	initCtx := context.Background()

	// Create token-related adapters
	tokenService, err := tokenregistry.NewTokenRegistryAdapter(initCtx, infra.EthClient, infra.Config)
	if err != nil {
		return fmt.Errorf("failed to create token service: %w", err)
	}
	tokenRepo := repositories.NewTokenRepository(infra.DBClient)
	participantRepo := repositories.NewParticipantRepository(infra.DBClient)

	config, err := configProvider.GetConfig(initCtx)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}
	pnData, err := decryptor.GatherParticipantsData(initCtx, config)
	if err != nil {
		return fmt.Errorf("failed to gather participants data: %w", err)
	}

	swapSaltsStore := repositories.NewDvpSwapSaltsRepository(infra.DBClient)

	// Create all event handlers
	eventHandlers := createEventHandlers(
		txRepo,
		revertRepo,
		tokenRepo,
		tokenService,
		participantRepo,
		tokenFreezeRepo,
		configProvider,
		provider,
		decryptor,
		log,
		&pnData,
		swapSaltsStore,
	)

	// Initialize NATS JetStream and message queue
	if infra.NATSConn == nil {
		return fmt.Errorf("NATS connection is required for the listener service (set NATS_URL)")
	}

	js, err := jetstream.New(infra.NATSConn)
	if err != nil {
		return fmt.Errorf("failed to create jetstream: %w", err)
	}

	const mqInitTimeout = 30 * time.Second
	mqCtx, mqCancel := context.WithTimeout(context.Background(), mqInitTimeout)
	defer mqCancel()

	manager, err := msgqueue.NewManager(mqCtx, js, infra.Config.PrivateHub.ChainId)
	if err != nil {
		return fmt.Errorf("failed to create message queue manager: %w", err)
	}

	contractPublisher := msgqueue.NewPublisher[core.ContractLog](manager, "contract_logs")
	proofsPublisher := msgqueue.NewPublisher[core.ContractLog](manager, "proofs_logs")

	contractConsumer, err := msgqueue.NewConsumer[core.ContractLog](mqCtx, manager, "event_dispatcher", "contract_logs")
	if err != nil {
		return fmt.Errorf("failed to create contract consumer: %w", err)
	}
	proofsConsumer, err := msgqueue.NewConsumer[core.ContractLog](mqCtx, manager, "proofs_dispatcher", "proofs_logs")
	if err != nil {
		return fmt.Errorf("failed to create proofs consumer: %w", err)
	}

	routingPubAdapter := msgqueue.NewRoutingPublisherAdapter(
		contractPublisher,
		proofsPublisher,
		func(l core.ContractLog) bool { return l.ContractName == events.ContractProofs },
	)

	logPublisher := services.NewLogPublisher(routingPubAdapter, log)

	governanceHandlers := eventHandlers

	eventDispatcher := services.NewEventDispatcher(
		msgqueue.NewConsumerAdapter(contractConsumer),
		governanceHandlers,
		log,
	)
	proofsBatchProcessor := services.NewProofsBatchProcessor(
		msgqueue.NewConsumerAdapter(proofsConsumer),
		headerProofRepo,
		log,
		services.DefaultProofsBatchMaxSize,
		services.DefaultProofsFlushInterval,
	)

	// Create BlockProcessor with LogPublisher
	blockProcessor := services.NewBlockProcessor(
		blockRepo,
		logParser,
		logPublisher,
		log,
		configProvider,
		provider,
	)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Info("Shutdown signal received, initiating graceful shutdown")
		cancel()
	}()

	// Run BlockProcessor and EventDispatcher as concurrent goroutines
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		return blockProcessor.Run(gCtx, ticker.C)
	})

	g.Go(func() error {
		return eventDispatcher.Run(gCtx)
	})

	g.Go(func() error {
		return proofsBatchProcessor.Run(gCtx)
	})

	if err := g.Wait(); err != nil {
		logger.Error("Service stopped with error", "error", err)
		return err
	}

	logger.Info("Shutdown complete")
	return nil
}

// createEventHandlers creates governance event handlers (excludes ProofsHandler).
func createEventHandlers(
	txRepo core.TransactionRepository,
	revertRepo core.RevertDataTransactionRepository,
	tokenRepo core.TokenRepository,
	tokenService core.TokenService,
	participantRepo core.ParticipantRepository,
	tokenFreezeRepo core.TokenFreezeRepository,
	configProvider core.ConfigProvider,
	provider core.Provider,
	decryptor core.Decryptor,
	log logger.Logger,
	pnData *core.PNodeDataAndSecrets,
	swapSalts core.SwapSaltsStore,
) []core.EventHandler {
	teleportHandler := handlers.NewTeleportEventHandler(txRepo, revertRepo, decryptor, log, pnData)

	enygmaTeleportHandler := handlers.NewEnygmaTeleportEventHandler(
		txRepo,
		decryptor,
		provider,
		log,
		pnData,
	)

	tokenCoreHandler := handlers.NewTokenCoreEventHandler(txRepo, tokenRepo, tokenService, log)

	participantCoreHandler := handlers.NewParticipantCoreEventHandler(participantRepo, log)

	auditManagerHandler := handlers.NewAuditManagerEventHandler(configProvider, decryptor, log, pnData)

	enygmaTokenManagerHandler := handlers.NewEnygmaTokenManagerEventHandler(tokenRepo, tokenService, log)

	dvpTeleportHandler := handlers.NewDvpTeleportEventHandler(txRepo, decryptor, provider, log, pnData, swapSalts)

	tokenFreezeManagerHandler := handlers.NewTokenFreezeManagerEventHandler(tokenFreezeRepo, provider, log)

	return []core.EventHandler{
		teleportHandler,
		enygmaTeleportHandler,
		tokenCoreHandler,
		participantCoreHandler,
		auditManagerHandler,
		enygmaTokenManagerHandler,
		dvpTeleportHandler,
		tokenFreezeManagerHandler,
	}
}
