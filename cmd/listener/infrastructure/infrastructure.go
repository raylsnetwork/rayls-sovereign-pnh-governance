package infrastructure

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/config"
	baseInfra "github.com/raylsnetwork/rayls-privacy-pnh-governance-api/infrastructure"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/infrastructure/database"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/mtls"
)

// Infrastructure holds all infrastructure-related components
type Infrastructure struct {
	Config    *config.Config
	DBClient  *gorm.DB
	EthClient *ethclient.Client
	NATSConn  *nats.Conn
}

// SetupInfrastructure initializes all infrastructure components
func SetupInfrastructure(configPath string) (*Infrastructure, error) {
	// Initialize configuration
	conf, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	// Initialize logger
	logger.InitializeLogger(conf)

	// Connect to database
	var dbClient *gorm.DB
	if conf.Database.ConnectionString != "" {
		var err error
		dbClient, err = database.Connect(conf.Database.ConnectionString)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
	}

	// Connect to the Private Network Hub client
	var ethClient *ethclient.Client
	if conf.PrivateHub.URL != "" {
		var err error
		ethClient, err = baseInfra.CreateClient(conf.PrivateHub.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to ethereum client: %w", err)
		}
	}

	// Connect to NATS (mTLS)
	var natsConn *nats.Conn
	if conf.NATSUrl != "" {
		if conf.NATSTLS.CAFile == "" || conf.NATSTLS.CertFile == "" || conf.NATSTLS.KeyFile == "" {
			return nil, fmt.Errorf(
				"NATS_TLS_CA_FILE, NATS_TLS_CERT_FILE, and NATS_TLS_KEY_FILE are required when NATS_URL is set",
			)
		}
		natsTLS, err := mtls.LoadClientConfig(conf.NATSTLS.CAFile, conf.NATSTLS.CertFile, conf.NATSTLS.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load NATS client TLS config: %w", err)
		}
		natsConn, err = nats.Connect(conf.NATSUrl, nats.Secure(natsTLS))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to NATS: %w", err)
		}
	}

	infra := &Infrastructure{
		Config:    conf,
		DBClient:  dbClient,
		EthClient: ethClient,
		NATSConn:  natsConn,
	}

	// Load contract addresses dynamically from DeploymentProxyRegistry
	if err := baseInfra.LoadContractAddresses(context.Background(), ethClient, conf); err != nil {
		return nil, err
	}

	return infra, nil
}

// Cleanup handles infrastructure cleanup
func (i *Infrastructure) Cleanup() error {
	// Drain NATS connection (waits for pending messages)
	if i.NATSConn != nil {
		if err := i.NATSConn.Drain(); err != nil {
			logger.Error("Error draining NATS connection", "error", err)
		}
	}

	// Close Ethereum client connection
	if i.EthClient != nil {
		i.EthClient.Close()
	}

	// Close database connection
	if i.DBClient != nil {
		return database.Disconnect(i.DBClient)
	}
	return nil
}
