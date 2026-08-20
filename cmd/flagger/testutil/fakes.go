package testutil

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
)

// FakeBalanceRepository is an in-memory implementation of BalanceRepository for testing.
type FakeBalanceRepository struct {
	Balances map[string]*domain.Balance // key: "chainId:resourceId:ercId"
}

// NewFakeBalanceRepository creates a new FakeBalanceRepository with initialized maps
func NewFakeBalanceRepository() *FakeBalanceRepository {
	return &FakeBalanceRepository{
		Balances: make(map[string]*domain.Balance),
	}
}

// balanceKey generates a key for the balance map
func balanceKey(chainId, resourceId, ercId string) string {
	return fmt.Sprintf("%s:%s:%s", chainId, resourceId, ercId)
}

// bigIntToErcIdString converts a domain.BigInt to a string key
func bigIntToErcIdString(b domain.BigInt) string {
	if b.Int == nil {
		return ""
	}
	return b.String()
}

func (r *FakeBalanceRepository) CreateBalance(ctx context.Context, balance *domain.Balance) error {
	key := balanceKey(balance.ChainId, balance.ResourceId, balance.ErcId.String())
	r.Balances[key] = balance
	return nil
}

func (r *FakeBalanceRepository) UpdateBalance(
	ctx context.Context,
	chainId, resourceId string,
	ercId domain.BigInt,
	amount string,
) error {
	key := balanceKey(chainId, resourceId, bigIntToErcIdString(ercId))
	if b, ok := r.Balances[key]; ok {
		b.Amount = decimal.RequireFromString(amount)
	}
	return nil
}

func (r *FakeBalanceRepository) GetResourceBalance(
	ctx context.Context,
	chainId, resourceId, ercId string,
) (*domain.Balance, error) {
	key := balanceKey(chainId, resourceId, ercId)
	if b, ok := r.Balances[key]; ok {
		return b, nil
	}
	// Create and store balance if not found
	balance := &domain.Balance{
		ChainId:    chainId,
		ResourceId: resourceId,
		Amount:     decimal.Zero,
	}
	r.Balances[key] = balance
	return balance, nil
}

func (r *FakeBalanceRepository) UpdateSenderReceiverBalances(
	ctx context.Context,
	senderChainId, senderResourceId string,
	senderErcId domain.BigInt,
	senderNewAmount string,
	receiverChainId, receiverResourceId string,
	receiverErcId domain.BigInt,
	receiverNewAmount string,
) error {
	// Update sender balance
	senderKey := balanceKey(senderChainId, senderResourceId, bigIntToErcIdString(senderErcId))
	if b, ok := r.Balances[senderKey]; ok {
		b.Amount = decimal.RequireFromString(senderNewAmount)
	}

	// Update receiver balance
	receiverKey := balanceKey(receiverChainId, receiverResourceId, bigIntToErcIdString(receiverErcId))
	if b, ok := r.Balances[receiverKey]; ok {
		b.Amount = decimal.RequireFromString(receiverNewAmount)
	}

	return nil
}

// FakeTransactionRepository is an in-memory implementation of TransactionRepository for testing.
type FakeTransactionRepository struct {
	Transactions        []domain.Transaction
	PendingTransactions []domain.Transaction // PENDING txs not returned by GetTransactions but found by FindPendingCounterpart
	ProcessedIDs        map[string]bool
	FlaggedIDs          map[string]bool
}

// NewFakeTransactionRepository creates a new FakeTransactionRepository with initialized maps
func NewFakeTransactionRepository() *FakeTransactionRepository {
	return &FakeTransactionRepository{
		Transactions:        make([]domain.Transaction, 0),
		PendingTransactions: make([]domain.Transaction, 0),
		ProcessedIDs:        make(map[string]bool),
		FlaggedIDs:          make(map[string]bool),
	}
}

func (r *FakeTransactionRepository) GetTransactions(ctx context.Context, limit int) ([]domain.Transaction, error) {
	if limit > len(r.Transactions) {
		return r.Transactions, nil
	}
	return r.Transactions[:limit], nil
}

func (r *FakeTransactionRepository) MarkTransactionAsProcessed(ctx context.Context, transactionId string) error {
	r.ProcessedIDs[transactionId] = true
	return nil
}

func (r *FakeTransactionRepository) FlagTransaction(ctx context.Context, transactionId string) error {
	r.FlaggedIDs[transactionId] = true
	return nil
}

func (r *FakeTransactionRepository) FindPendingCounterpartByReferenceId(
	ctx context.Context,
	transactionId string,
) (*domain.Transaction, error) {
	allTxs := make([]domain.Transaction, 0, len(r.Transactions)+len(r.PendingTransactions))
	allTxs = append(allTxs, r.Transactions...)
	allTxs = append(allTxs, r.PendingTransactions...)

	// Find the current transaction and its reference_id
	var refId string
	var currentTx *domain.Transaction
	for i := range allTxs {
		tx := &allTxs[i]
		if tx.ID.String() == transactionId {
			currentTx = tx
			if tx.EnygmaTransaction != nil {
				refId = tx.EnygmaTransaction.ReferenceId
			}
			break
		}
	}
	if refId == "" || currentTx == nil {
		return nil, nil //nolint:nilnil
	}

	// Find a PENDING counterpart with the same reference_id in the reverse direction
	for i := range allTxs {
		tx := &allTxs[i]
		if tx.ID.String() != transactionId &&
			tx.EnygmaTransaction != nil &&
			tx.EnygmaTransaction.ReferenceId == refId &&
			tx.GetFromChainId().String() == currentTx.GetToChainId().String() &&
			tx.GetToChainId().String() == currentTx.GetFromChainId().String() &&
			tx.TeleportStatus != nil && *tx.TeleportStatus == 0 &&
			!r.ProcessedIDs[tx.ID.String()] {
			return tx, nil
		}
	}
	return nil, nil //nolint:nilnil
}

func (r *FakeTransactionRepository) Transact(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// FakeHeaderProofRepository is an in-memory implementation of HeaderProofRepository for testing.
type FakeHeaderProofRepository struct {
	HeaderProofs []domain.HeaderProofEvent
	Error        error
}

// NewFakeHeaderProofRepository creates a new FakeHeaderProofRepository
func NewFakeHeaderProofRepository() *FakeHeaderProofRepository {
	return &FakeHeaderProofRepository{
		HeaderProofs: make([]domain.HeaderProofEvent, 0),
	}
}

func (r *FakeHeaderProofRepository) GetLatestHeaderProofs(ctx context.Context) ([]domain.HeaderProofEvent, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	return r.HeaderProofs, nil
}

func (r *FakeHeaderProofRepository) DeleteOlderThan(_ context.Context, cutoff time.Time) (int64, error) {
	if r.Error != nil {
		return 0, r.Error
	}

	// Mirror SQL semantics: preserve the record(s) with MAX(block_number) per chain_id,
	// even if those records are older than the cutoff.
	maxBlock := map[string]*big.Int{}
	for _, h := range r.HeaderProofs {
		key := h.ChainID.String()
		if cur, ok := maxBlock[key]; !ok || h.BlockNumber.Cmp(cur) > 0 {
			maxBlock[key] = new(big.Int).Set(h.BlockNumber.Int)
		}
	}

	var kept []domain.HeaderProofEvent
	var deleted int64
	for _, h := range r.HeaderProofs {
		isLatestForChain := h.BlockNumber.Cmp(maxBlock[h.ChainID.String()]) == 0
		if h.CreatedAt.Before(cutoff) && !isLatestForChain {
			deleted++
		} else {
			kept = append(kept, h)
		}
	}
	r.HeaderProofs = kept
	return deleted, nil
}

// FakeHeaderFlagEventRepository is an in-memory implementation of HeaderFlagEventRepository for testing.
type FakeHeaderFlagEventRepository struct {
	// FlaggedParticipants maps "chainId:blockNumber" to the reason the flag was raised with,
	// mirroring the real repository's per-reason scoping so tests can assert that
	// UnflagParticipant clears only flags matching the given reason.
	FlaggedParticipants map[string]uint8
	Error               error
}

// NewFakeHeaderFlagEventRepository creates a new FakeHeaderFlagEventRepository
func NewFakeHeaderFlagEventRepository() *FakeHeaderFlagEventRepository {
	return &FakeHeaderFlagEventRepository{
		FlaggedParticipants: make(map[string]uint8),
	}
}

// flagKey generates a key for the flagged participants map
func flagKey(chainID, blockNumber *big.Int) string {
	return fmt.Sprintf("%s:%s", chainID.String(), blockNumber.String())
}

func (r *FakeHeaderFlagEventRepository) FlagParticipant(
	ctx context.Context,
	chainID *big.Int,
	blockNumber *big.Int,
	reason uint8,
	initiator uint8,
) (bool, error) {
	if r.Error != nil {
		return false, r.Error
	}
	key := flagKey(chainID, blockNumber)
	if _, ok := r.FlaggedParticipants[key]; ok {
		return false, nil // already flagged
	}
	r.FlaggedParticipants[key] = reason
	return true, nil // newly flagged
}

// IsFlagged checks if a participant is flagged (helper for test assertions)
func (r *FakeHeaderFlagEventRepository) IsFlagged(chainID, blockNumber *big.Int) bool {
	key := flagKey(chainID, blockNumber)
	_, ok := r.FlaggedParticipants[key]
	return ok
}

// UnflagParticipant clears flag(s) for the given chain whose recorded reason matches the
// supplied reason, mirroring the real repository's `WHERE chain_id = ? AND flag_reason = ?`
// scoping. Flags raised for any other reason are left untouched. Returns true if anything
// was cleared.
func (r *FakeHeaderFlagEventRepository) UnflagParticipant(
	ctx context.Context,
	chainID *big.Int,
	reason uint8,
) (bool, error) {
	if r.Error != nil {
		return false, r.Error
	}
	prefix := chainID.String() + ":"
	var unflagged bool
	for key, flagReason := range r.FlaggedParticipants {
		if strings.HasPrefix(key, prefix) && flagReason == reason {
			delete(r.FlaggedParticipants, key)
			unflagged = true
		}
	}
	return unflagged, nil
}
