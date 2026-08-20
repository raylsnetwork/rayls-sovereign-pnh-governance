package infrastructure

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/config"
	baseInfra "github.com/raylsnetwork/rayls-privacy-pnh-governance-api/infrastructure"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/infrastructure/database"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
)

// Infrastructure holds all infrastructure-related components for the API
type Infrastructure struct {
	Config    *config.Config
	DBClient  *gorm.DB
	EthClient *ethclient.Client
}

// SetupInfrastructure initializes all infrastructure components for the API service
func SetupInfrastructure(configPath string) (*Infrastructure, error) {
	// Initialize configuration
	conf, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	// Initialize logger
	logger.InitializeLogger(conf)

	// Connect to database
	dbClient, err := database.Connect(conf.Database.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Connect to the Private Network Hub client
	var ethClient *ethclient.Client
	if conf.PrivateHub.URL != "" {
		ethClient, err = baseInfra.CreateClient(conf.PrivateHub.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to private network hub: %w", err)
		}

		// Load contract addresses dynamically from DeploymentProxyRegistry
		if err := baseInfra.LoadContractAddresses(context.Background(), ethClient, conf); err != nil {
			return nil, fmt.Errorf("failed to load contract addresses: %w", err)
		}
	}

	infra := &Infrastructure{
		Config:    conf,
		DBClient:  dbClient,
		EthClient: ethClient,
	}

	return infra, nil
}

// Cleanup handles infrastructure cleanup
func (i *Infrastructure) Cleanup() error {
	// Close Ethereum client connection
	if i.EthClient != nil {
		i.EthClient.Close()
	}

	// Close database connection
	if i.DBClient != nil {
		if err := database.Disconnect(i.DBClient); err != nil {
			logger.Error("Error disconnecting from the DB", "error", err)
			return err
		}
	}
	return nil
}
