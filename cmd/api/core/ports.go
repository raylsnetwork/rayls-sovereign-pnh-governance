package core

import (
	"context"
	"errors"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/dto"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
)

// ============================================================================
// PRIMARY PORTS (Application Use Cases)
// These are what the Core provides - the business logic interface
// ============================================================================

// TransactionService defines the business operations for transactions
type TransactionService interface {
	// GetTransactionByMessageId retrieves a single transaction by messageId
	GetTransactionByMessageId(ctx context.Context, messageId string) (*dto.TransactionDetailDto, error)

	// GetTransactionByTransactionId retrieves a single transaction by transactionId
	GetTransactionByTransactionId(ctx context.Context, transactionId string) (*dto.TransactionDetailDto, error)

	// GetTransactionsList retrieves a paginated list of transactions
	GetTransactionsList(
		ctx context.Context,
		filters dto.MergedTransactionsFilters,
	) (*types.Paginated[dto.TransactionListDto], error)

	// GetEnygmaBatchTransactions retrieves paginated transactions in an enygma batch
	GetEnygmaBatchTransactions(
		ctx context.Context,
		batchId string,
		page, limit int,
	) (*types.Paginated[dto.BatchTransactionDto], error)

	// GetRegularBatchTransactions retrieves paginated transactions in a regular batch
	GetRegularBatchTransactions(
		ctx context.Context,
		batchId string,
		page, limit int,
	) (*types.Paginated[dto.BatchTransactionDto], error)

	// GetTransactionsBySharedId retrieves transactions by shared_id (zkdvp swaps)
	GetTransactionsBySharedId(
		ctx context.Context,
		sharedId string,
	) ([]dto.DvpSwapTransactionDto, error)

	// GetFlaggedTransactions retrieves all flagged transactions
	GetFlaggedTransactions(ctx context.Context) ([]domain.FlaggedTransaction, error)
}

// ParticipantService defines the business operations for participants
type ParticipantService interface {
	// GetParticipantByChainId retrieves a single participant by chain ID
	GetParticipantByChainId(ctx context.Context, chainId string) (*domain.Participant, error)

	// GetParticipantsList retrieves a list of participants with filters
	GetParticipantsList(ctx context.Context, filters dto.ParticipantListFilters) ([]domain.Participant, error)
}

// TokenService defines the business operations for tokens
type TokenService interface {
	// GetTokenByResourceId retrieves a single token by resource ID with balances and restrictions
	GetTokenByResourceId(ctx context.Context, resourceId string) (*domain.TokenWithBalancesAndFreezeState, error)

	// GetTokensList retrieves a paginated list of tokens with filters
	GetTokensList(
		ctx context.Context,
		filters dto.TokenListFilters,
	) (*types.Paginated[domain.TokenWithBalancesAndFreezeState], error)

	// GetTokenRegistryStatus retrieves token registry status (subset of fields) by resource ID
	GetTokenRegistryStatus(ctx context.Context, resourceId string) (*dto.TokenRegistryStatusDto, error)
}

// HeaderProofService defines the business operations for header proofs
type HeaderProofService interface {
	// GetHeaderProofsList retrieves header proofs for a block range with pagination
	GetHeaderProofsList(ctx context.Context, filters dto.HeaderProofFilters) (*dto.HeaderProofListResponse, error)
}

// AuthService defines the business operations for private network authentication
type AuthService interface {
	// SignUp creates a new private network account
	SignUp(ctx context.Context, username, password string) error

	// Login authenticates a private network and returns a JWT token
	Login(ctx context.Context, username, password string) (string, error)

	// ValidateToken validates a JWT token and returns the associated private network
	ValidateToken(ctx context.Context, tokenString string) (*domain.PrivateNetwork, error)
}

// BalanceService defines the business operations for token balance queries
type BalanceService interface {
	// GetBalancesInChain retrieves balance(s) for a chain
	// If resourceId is "/" or empty, returns all balances in the chain
	// Otherwise returns specific resource balance
	GetBalancesInChain(ctx context.Context, chainId, resourceId string) (interface{}, error)

	// GetBalanceAcrossAllChains retrieves balance for a resource across all chains
	GetBalanceAcrossAllChains(ctx context.Context, resourceId string) ([]domain.Balance, error)

	// GetBalanceAcrossSpecificChains retrieves balance for a resource across specific chains
	// If chains is empty, delegates to GetBalanceAcrossAllChains
	GetBalanceAcrossSpecificChains(ctx context.Context, resourceId string, chains []string) ([]domain.Balance, error)
}

// ============================================================================
// SECONDARY PORTS (Infrastructure Needs)
// These are what the Core needs - dependencies injected from outside
// ============================================================================

// ErrRecordNotFound should be returned by repositories when a queried record doesn't exist
var ErrRecordNotFound = errors.New("record not found")

// TransactionRepository handles database queries for transactions
type TransactionRepository interface {
	// FindByMessageId finds a single transaction by message_id field
	FindByMessageId(ctx context.Context, messageId string) (*domain.Transaction, error)

	// FindByTransactionId finds a single transaction by transaction_id field
	FindByTransactionId(ctx context.Context, transactionId string) (*domain.Transaction, error)

	// FindByBatchId finds transactions by batch_id field
	FindByBatchId(ctx context.Context, batchId string) ([]domain.Transaction, error)

	// FindBySharedId finds transactions by shared_id field
	FindBySharedId(ctx context.Context, sharedId string) ([]domain.Transaction, error)

	// FindByBatchIdPaginated finds transactions by batch_id with pagination
	FindByBatchIdPaginated(
		ctx context.Context,
		batchId string,
		page, limit int,
	) ([]domain.Transaction, int64, error)

	// FindByEnygmaBatchId finds transactions by enygma_transactions.batch_id
	FindByEnygmaBatchId(ctx context.Context, batchId string) ([]domain.Transaction, error)

	// FindByEnygmaBatchIdPaginated finds enygma transactions by batch_id with pagination
	FindByEnygmaBatchIdPaginated(
		ctx context.Context,
		batchId string,
		page, limit int,
	) ([]domain.Transaction, int64, error)

	// FindWithFilters finds transactions matching the provided filters with pagination
	FindWithFilters(ctx context.Context, filters dto.MergedTransactionsFilters) ([]domain.Transaction, int64, error)

	// FindFlagged retrieves all flagged transactions
	FindFlagged(ctx context.Context) ([]domain.FlaggedTransaction, error)
}

// TokenRepository handles database queries for tokens
type TokenRepository interface {
	// FindByResourceIdWithBalances retrieves a token with its balances and freeze restrictions
	FindByResourceIdWithBalances(
		ctx context.Context,
		resourceId string,
	) (*domain.TokenWithBalancesAndFreezeState, error)

	// FindByFilters finds tokens matching the provided filters with pagination
	FindByFilters(
		ctx context.Context,
		filters dto.TokenListFilters,
	) ([]domain.TokenWithBalancesAndFreezeState, int64, error)
}

// ParticipantRepository handles database queries for participants
type ParticipantRepository interface {
	// FindByChainId retrieves a single participant by chain ID
	FindByChainId(ctx context.Context, chainId string) (*domain.Participant, error)

	// FindByFilters finds participants matching the provided filters
	FindByFilters(ctx context.Context, filters dto.ParticipantListFilters) ([]domain.Participant, error)
}

// HeaderProofRepository handles database queries for header proofs
type HeaderProofRepository interface {
	// FindByBlockRange retrieves header proofs within a block range with pagination
	FindByBlockRange(
		ctx context.Context,
		chainId string,
		startBlock, endBlock int64,
		page, pageSize int,
	) ([]domain.HeaderProofEvent, int64, error)
}

// PrivateNetworkRepository handles database queries for private networks
type PrivateNetworkRepository interface {
	// FindByUsername retrieves a private network by username
	FindByUsername(ctx context.Context, username string) (*domain.PrivateNetwork, error)

	// Create creates a new private network
	Create(ctx context.Context, username, hashedPassword string) error
}

// BalanceRepository handles database queries for token balances
type BalanceRepository interface {
	// FindAllInChain retrieves all balances for a specific chain
	FindAllInChain(ctx context.Context, chainId string) ([]domain.Balance, error)

	// FindByChainAndResource retrieves balance for specific resource in chain
	FindByChainAndResource(ctx context.Context, chainId, resourceId string) (*domain.Balance, error)

	// FindAcrossAllChains retrieves balances for a resource across all chains
	FindAcrossAllChains(ctx context.Context, resourceId string) ([]domain.Balance, error)

	// FindAcrossSpecificChains retrieves balances for a resource across specific chains
	FindAcrossSpecificChains(ctx context.Context, resourceId string, chains []string) ([]domain.Balance, error)
}

// TokenMetadataService handles external HTTP calls to fetch NFT metadata
// This abstracts the external dependency
type TokenMetadataService interface {
	// GetMetadata fetches token metadata from external URL
	// Returns metadata info for ERC721/ERC1155 tokens
	GetMetadata(ctx context.Context, baseURL, ercId string) (*dto.TokenMetadataInfoDto, error)
}
