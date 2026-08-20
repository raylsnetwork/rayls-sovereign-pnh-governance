package infrastructure

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/config"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/infrastructure/database"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
)

// Infrastructure holds all infrastructure-related components for the flagger
type Infrastructure struct {
	Config   *config.Config
	DBClient *gorm.DB
}

// SetupInfrastructure initializes all infrastructure components for the flagger service
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

	infra := &Infrastructure{
		Config:   conf,
		DBClient: dbClient,
	}

	return infra, nil
}

// Cleanup handles infrastructure cleanup
func (i *Infrastructure) Cleanup() error {
	// Close database connection
	if i.DBClient != nil {
		return database.Disconnect(i.DBClient)
	}
	return nil
}
