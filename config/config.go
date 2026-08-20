package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Database struct {
	ConnectionString string `mapstructure:"DATABASE_CONNECTIONSTRING"`
	Type             string `mapstructure:"DATABASE_TYPE"`
}

type PrivateHub struct {
	URL                         string `mapstructure:"PNH_RPC_URL"`
	DeploymentProxyRegistry     string `mapstructure:"PNH_DEPLOYMENT_PROXY_REGISTRY"`
	PrivateKey                  string `mapstructure:"PNH_PRIVATE_KEY"`
	ChainId                     string `mapstructure:"PNH_CHAIN_ID"`
	RaylsViewSecretKey          string `mapstructure:"PNH_RAYLS_VIEW_SECRET_KEY"`
	StartingBlock               string `mapstructure:"PNH_STARTING_BLOCK"`
	BatchSize                   int64  `mapstructure:"PNH_BATCH_SIZE"`
	OperatorChainId             int64  `mapstructure:"PNH_OPERATOR_CHAIN_ID"`
	HeaderProofExpirationPeriod string `mapstructure:"PNH_HEADER_PROOF_EXPIRATION_PERIOD"`
	HeaderProofPurgePeriod      string `mapstructure:"PNH_HEADER_PROOF_PURGE_PERIOD"`
	// Active contracts
	EnygmaTokenManager         string `mapstructure:"PNH_ENYGMA_TOKEN_MANAGER"`
	TokenCore                  string `mapstructure:"PNH_TOKEN_CORE"`
	TokenFreezeManager         string `mapstructure:"PNH_TOKEN_FREEZE_MANAGER"`
	ParticipantCore            string `mapstructure:"PNH_PARTICIPANT_CORE"`
	AuditManager               string `mapstructure:"PNH_AUDIT_MANAGER"`
	Teleport                   string `mapstructure:"PNH_TELEPORT"`
	EnygmaPNHEvents            string `mapstructure:"PNH_ENYGMA_EVENTS"`
	EnygmaTeleport             string `mapstructure:"PNH_ENYGMA_TELEPORT"`
	ParticipantStorageContract string `mapstructure:"PNH_PARTICIPANT_STORAGE_CONTRACT"`
	TokenRegistry              string `mapstructure:"PNH_TOKEN_REGISTRY"`
	DvpTeleport                string `mapstructure:"PNH_DVP_TELEPORT"`
	// Keep Proofs for future use (currently commented out in BlockProcessor)
	ProofsAddress string `mapstructure:"PNH_PROOFS_ADDRESS"`
}

type TransactionProcessor struct {
	BatchSize     int    `mapstructure:"TRANSACTIONPROCESSOR_BATCHSIZE"`
	CheckInterval string `mapstructure:"TRANSACTIONPROCESSOR_CHECKINTERVAL"`
}

// NATSTLS holds the file paths for the mTLS keypair and trusted CA
// used when connecting to NATS. All three are required when NATSUrl
// is set — NATS in this stack always requires client certs.
type NATSTLS struct {
	CAFile   string `mapstructure:"NATS_TLS_CA_FILE"`
	CertFile string `mapstructure:"NATS_TLS_CERT_FILE"`
	KeyFile  string `mapstructure:"NATS_TLS_KEY_FILE"`
}

type Config struct {
	Database             Database             `mapstructure:",squash"`
	PrivateHub           PrivateHub           `mapstructure:",squash"`
	TransactionProcessor TransactionProcessor `mapstructure:",squash"`
	JWTSecret            string               `mapstructure:"JWTSECRET"`
	Logging              string               `mapstructure:"LOGGING"`
	CorsUrls             string               `mapstructure:"CORSURLS"`
	NATSUrl              string               `mapstructure:"NATS_URL"`
	NATSTLS              NATSTLS              `mapstructure:",squash"`
}

func Load(envPath string) (*Config, error) {
	// Create a new Viper instance with ExperimentalBindStruct enabled
	// This automatically binds all mapstructure tags to environment variables
	v := viper.NewWithOptions(viper.ExperimentalBindStruct())
	v.AutomaticEnv()

	if envPath != "" {
		// If a specific .env file path is provided, read from it.
		// Force "env" type so Viper correctly parses KEY=VALUE files
		// regardless of the file extension (e.g. .env.local, .env.docker).
		v.SetConfigType("env")
		v.SetConfigFile(envPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read .env file: %w", err)
		}
	}
	// If no path provided, the config will work purely with OS environment variables
	// ExperimentalBindStruct() handles the binding automatically

	// Unmarshal configuration
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}
