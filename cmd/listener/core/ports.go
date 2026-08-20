package core

import (
	"context"
	"errors"
	"math/big"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/config"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

// TransactionRepository defines the interface for transaction persistence
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *domain.Transaction) error
	CreateTransactions(ctx context.Context, transactions []*domain.Transaction) error
	CreateTransactionsWithEnygmaData(
		ctx context.Context,
		transactions []*domain.Transaction,
		enygmaTransactions []*domain.EnygmaTransaction,
	) error
	CreateCreditMemo(ctx context.Context, tx *domain.Transaction) error
	GetTransactionByMessageID(ctx context.Context, messageID string, preload bool) (*domain.Transaction, error)
	GetTransactionsByMessageIDs(ctx context.Context, messageIDs []string, preload bool) ([]domain.Transaction, error)
	GetTransactionsBySharedIDs(ctx context.Context, sharedIDs []string, preload bool) ([]domain.Transaction, error)
	UpdateTransaction(ctx context.Context, tx *domain.Transaction) error
	UpdateTransactionsBulk(ctx context.Context, transactions []*domain.Transaction, columns ...string) error
	UpdateTeleportStatusBySharedID(ctx context.Context, sharedID string, status uint8) error
	UpdateDvpTeleportConfirmation(
		ctx context.Context,
		sharedID string,
		chainID *big.Int,
		to string,
		pnTxHash string,
		pnTxTimestamp time.Time,
		calldata string,
	) error
	PersistUpdateTransaction(
		ctx context.Context,
		ut *types.UpdateTransaction,
		fromChain, toChain string,
	) (uuid.UUID, error)
	PersistUpdateTransactions(
		ctx context.Context,
		uts *[]types.UpdateTransaction,
		fromChain, toChain string,
	) ([]uuid.UUID, error)

	// CreateTransactionsWithPromotion creates transactions and promotes batches atomically
	// This ensures either all operations succeed or none do
	CreateTransactionsWithPromotion(
		ctx context.Context,
		transactions []*domain.Transaction,
		batchIds []string,
	) error
}

// EnygmaTransactionRepository defines the interface for enygma transaction persistence
type EnygmaTransactionRepository interface {
	UpsertEnygmaTransactionsBulk(ctx context.Context, enygma_transactions []*domain.EnygmaTransaction) error
}

// RevertDataTransactionRepository defines the interface for revert data transaction persistence
type RevertDataTransactionRepository interface {
	CreateRevertTransactions(ctx context.Context, revert_transactions []*domain.RevertDataTransaction) error
}

// BlockRepository defines the interface for block tracking
type BlockRepository interface {
	GetLatestProcessedBlock(ctx context.Context) (*big.Int, error)
	UpdateLatestProcessedBlock(ctx context.Context, blockNumber *big.Int) error
}

// HeaderProofEventRepository defines the interface for header proof event persistence
type HeaderProofEventRepository interface {
	Create(ctx context.Context, event *domain.HeaderProofEvent) error
	CreateBatch(ctx context.Context, events []*domain.HeaderProofEvent) error
	GetByBlockNumber(ctx context.Context, chainID *big.Int, blockNumber *big.Int) (*domain.HeaderProofEvent, error)
}

// LogParser defines the interface for parsing smart contract logs
type LogParser interface {
	ParseLogs(ctx context.Context, fromBlock, toBlock *big.Int) ([]ContractLog, error)
}

// ConfigProvider defines the interface for getting configuration values
type ConfigProvider interface {
	GetStartingBlockNumber() *big.Int
	GetBatchSize() int64
	GetConfig(ctx context.Context) (*config.Config, error)
}

// Provider defines the interface for interacting with the blockchain
type Provider interface {
	GetLatestBlock(ctx context.Context) (*ethTypes.Block, error)
	GetBlockByNumber(ctx context.Context, number uint64) (*ethTypes.Block, error)
}

// PNodeDataAndSecrets represents Privacy Node data and secrets
type PNodeDataAndSecrets = map[string]map[string]*types.IPNodeDataAndSecrets

// Decryptor defines the interface for decryption operations
type Decryptor interface {
	GatherParticipantsData(ctx context.Context, conf *config.Config) (PNodeDataAndSecrets, error)
	DecryptPayloadBytes(
		payload []byte,
		blockNumber uint64,
		pnData PNodeDataAndSecrets,
		secretType types.SecretType,
	) ([]byte, error)
	// DecryptSwapPayload decrypts a SwapInitiated payload. The ciphertext is
	// encapsulated against the destination-chain participant's view public
	// key — not the operator's — so the caller must pass the gathered
	// pnData and blockNumber, and the implementation tries each candidate
	// participant DK until one succeeds. Returns the plaintext and the salt
	// reduced from the Ctxt-derived shared secret (InitiatorCtxtSalt), which
	// the caller persists alongside the payload-carried InitiatorSelfSalt.
	DecryptSwapPayload(
		ctxt []byte,
		encryptedData []byte,
		blockNumber uint64,
		pnData PNodeDataAndSecrets,
	) (plaintext []byte, initiatorCtxtSalt []byte, err error)
	// DecryptWithSalt decrypts a GCM payload with a caller-supplied salt.
	// Used for DVP swap payloads where the salt is persisted per swap rather
	// than fetched from pnData.
	DecryptWithSalt(payload []byte, salt []byte) ([]byte, error)
}

// SwapSaltsStore persists the per-swap salts recovered on SwapInitiated so
// subsequent events for the same swap (SwapCompleted, and potential future
// responder→initiator messages) can reuse them across listener restarts.
type SwapSaltsStore interface {
	Put(ctx context.Context, sharedID string, salts types.DvpSwapSalts) error
	Get(ctx context.Context, sharedID string) (types.DvpSwapSalts, error)
}

// ErrSwapSaltsNotFound is returned by SwapSaltsStore.Get when the sharedID
// has no stored salts.
var ErrSwapSaltsNotFound = errors.New("swap salts not found")

// TokenService defines the interface for token-related operations
type TokenService interface {
	GetTokenByResourceId(ctx context.Context, resourceId string) (*domain.Token, error)
}

// TokenRepository defines the interface for token persistence
type TokenRepository interface {
	GetTokenByResourceID(ctx context.Context, resourceID string) (*domain.Token, error)
	GetByIssuerAndName(ctx context.Context, issuerChainId string, name string) (*domain.Token, error)
	Upsert(ctx context.Context, token *domain.Token) error
}

// ParticipantRepository defines the interface for participant data persistence
type ParticipantRepository interface {
	GetByChainId(ctx context.Context, chainId string) (domain.Participant, error)
	Upsert(ctx context.Context, participant domain.Participant) error
}

// TokenFreezeRepository defines the interface for token freeze operations
type TokenFreezeRepository interface {
	UpdateTokenFreezeStatus(
		ctx context.Context,
		resourceId string,
		chainIds []*big.Int,
		action uint8,
		blockNumber *big.Int,
		txHash string,
		blockTimestamp time.Time,
	) error
}

type EventHandler interface {
	ContractName() string
	Handle(ctx context.Context, log ContractLog) error
	Name() string
}

// Message wraps a consumed contract log with its acknowledgment callback.
// It is a non-generic equivalent of msgqueue.Message[ContractLog], keeping
// the core package free of msgqueue imports.
type Message struct {
	Log ContractLog
	Ack func(context.Context) error
}

// ContractMQ defines the interface for publishing contract logs to a message queue.
type ContractMQ interface {
	Push(ctx context.Context, log ContractLog) error
}

// EventConsumer defines the interface for consuming contract logs from a message queue.
type EventConsumer interface {
	Next(ctx context.Context) (Message, error)
}
