package testutil

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/dto"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

// ErrRecordNotFound mirrors core.ErrRecordNotFound to avoid import cycle
var ErrRecordNotFound = errors.New("record not found")

// FakeParticipantRepository is an in-memory implementation for testing
type FakeParticipantRepository struct {
	Participants []domain.Participant
	Error        error
	NotFoundErr  error // override to inject core.ErrRecordNotFound so errors.Is matches in service layer
}

func NewFakeParticipantRepository() *FakeParticipantRepository {
	return &FakeParticipantRepository{
		Participants: []domain.Participant{},
		NotFoundErr:  ErrRecordNotFound,
	}
}

func (r *FakeParticipantRepository) FindByChainId(ctx context.Context, chainId string) (*domain.Participant, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	for i := range r.Participants {
		p := &r.Participants[i]
		if p.ChainId != nil && strconv.FormatUint(uint64(*p.ChainId), 10) == chainId {
			return p, nil
		}
	}
	return nil, r.NotFoundErr
}

func (r *FakeParticipantRepository) FindByFilters(
	ctx context.Context,
	filters dto.ParticipantListFilters,
) ([]domain.Participant, error) {
	if r.Error != nil {
		return nil, r.Error
	}

	result := make([]domain.Participant, 0, len(r.Participants))
	for _, p := range r.Participants {
		// Filter by name (case-insensitive contains)
		if filters.Name != "" && !strings.Contains(strings.ToLower(p.Name), strings.ToLower(filters.Name)) {
			continue
		}
		// Filter by chainId (exact match)
		if filters.ChainId != nil && (p.ChainId == nil || *p.ChainId != *filters.ChainId) {
			continue
		}
		// Filter by status (case-insensitive exact match)
		if filters.Status != "" && !strings.EqualFold(p.StatusStr, filters.Status) {
			continue
		}
		// Filter by role (case-insensitive exact match)
		if filters.Role != "" && !strings.EqualFold(p.RoleStr, filters.Role) {
			continue
		}
		result = append(result, p)
	}
	return result, nil
}

// HasParticipant checks if a participant with the given name exists (helper for assertions)
func (r *FakeParticipantRepository) HasParticipant(name string) bool {
	for _, p := range r.Participants {
		if p.Name == name {
			return true
		}
	}
	return false
}

// FakeBalanceRepository is an in-memory implementation for testing
type FakeBalanceRepository struct {
	Balances    []domain.Balance
	Error       error
	NotFoundErr error // override to inject core.ErrRecordNotFound so errors.Is matches in service layer
}

func NewFakeBalanceRepository() *FakeBalanceRepository {
	return &FakeBalanceRepository{
		Balances:    []domain.Balance{},
		NotFoundErr: ErrRecordNotFound,
	}
}

func (r *FakeBalanceRepository) FindAllInChain(ctx context.Context, chainId string) ([]domain.Balance, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	var result []domain.Balance
	for _, b := range r.Balances {
		if b.ChainId == chainId {
			result = append(result, b)
		}
	}
	return result, nil
}

func (r *FakeBalanceRepository) FindByChainAndResource(
	ctx context.Context,
	chainId, resourceId string,
) (*domain.Balance, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	for i := range r.Balances {
		b := &r.Balances[i]
		if b.ChainId == chainId && b.ResourceId == resourceId {
			return b, nil
		}
	}
	return nil, r.NotFoundErr
}

func (r *FakeBalanceRepository) FindAcrossAllChains(ctx context.Context, resourceId string) ([]domain.Balance, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	var result []domain.Balance
	for _, b := range r.Balances {
		if b.ResourceId == resourceId {
			result = append(result, b)
		}
	}
	return result, nil
}

func (r *FakeBalanceRepository) FindAcrossSpecificChains(
	ctx context.Context,
	resourceId string,
	chains []string,
) ([]domain.Balance, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	chainSet := make(map[string]bool)
	for _, c := range chains {
		chainSet[c] = true
	}
	var result []domain.Balance
	for _, b := range r.Balances {
		if b.ResourceId == resourceId && chainSet[b.ChainId] {
			result = append(result, b)
		}
	}
	return result, nil
}

// HasBalance checks if a balance exists for the given chainId and resourceId (helper for assertions)
func (r *FakeBalanceRepository) HasBalance(chainId, resourceId string) bool {
	for _, b := range r.Balances {
		if b.ChainId == chainId && b.ResourceId == resourceId {
			return true
		}
	}
	return false
}

// GetAmount returns the amount for a specific balance (helper for assertions)
func (r *FakeBalanceRepository) GetAmount(chainId, resourceId string) decimal.Decimal {
	for _, b := range r.Balances {
		if b.ChainId == chainId && b.ResourceId == resourceId {
			return b.Amount
		}
	}
	return decimal.Zero
}

// FakeTokenRepository is an in-memory implementation for testing tokens
type FakeTokenRepository struct {
	Tokens []domain.TokenWithBalancesAndFreezeState
	Error  error
}

func NewFakeTokenRepository() *FakeTokenRepository {
	return &FakeTokenRepository{
		Tokens: []domain.TokenWithBalancesAndFreezeState{},
	}
}

func (r *FakeTokenRepository) FindByResourceIdWithBalances(
	ctx context.Context,
	resourceId string,
) (*domain.TokenWithBalancesAndFreezeState, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	for i := range r.Tokens {
		t := &r.Tokens[i]
		if t.ResourceId == resourceId {
			return t, nil
		}
	}
	// mirror GORM adapter behaviour: return empty struct when not found
	return &domain.TokenWithBalancesAndFreezeState{}, nil
}

func (r *FakeTokenRepository) FindByFilters(
	ctx context.Context,
	filters dto.TokenListFilters,
) ([]domain.TokenWithBalancesAndFreezeState, int64, error) {
	if r.Error != nil {
		return nil, 0, r.Error
	}

	result := make([]domain.TokenWithBalancesAndFreezeState, 0, len(r.Tokens))
	for _, t := range r.Tokens {
		// Filter by name (case-insensitive contains)
		if filters.Name != "" && !strings.Contains(strings.ToLower(t.Name), strings.ToLower(filters.Name)) {
			continue
		}
		// Filter by symbol (case-insensitive contains)
		if filters.Symbol != "" && !strings.Contains(strings.ToLower(t.Symbol), strings.ToLower(filters.Symbol)) {
			continue
		}
		// Filter by status (exact match via domain mapping)
		if filters.Status != "" {
			statusValue, ok := domain.StringToTokenStatus[strings.ToLower(filters.Status)]
			if ok && t.Status != uint8(statusValue) { //nolint:gosec // statusValue is 0-2, safe for uint8
				continue
			}
		}
		// Filter by ercStandard (exact match via types mapping)
		if filters.ErcStandard != "" {
			ercValue, ok := types.StringToAssetType[strings.ToLower(filters.ErcStandard)]
			if ok && t.ErcStandard != ercValue {
				continue
			}
		}
		// Filter by issuerId (exact match)
		if filters.IssuerId != "" && t.IssuerId != filters.IssuerId {
			continue
		}
		// Filter by decimals (exact match)
		if filters.Decimals != nil && t.Decimals != *filters.Decimals {
			continue
		}
		result = append(result, t)
	}
	return result, int64(len(result)), nil
}

// HasToken checks if a token with the given name exists (helper for assertions)
func (r *FakeTokenRepository) HasToken(name string) bool {
	for _, t := range r.Tokens {
		if t.Name == name {
			return true
		}
	}
	return false
}

// FakeTokenMetadataService returns preconfigured metadata for ERC tokens
type FakeTokenMetadataService struct {
	Data  map[string]*dto.TokenMetadataInfoDto // key: ercId
	Error error
}

func NewFakeTokenMetadataService() *FakeTokenMetadataService {
	return &FakeTokenMetadataService{Data: make(map[string]*dto.TokenMetadataInfoDto)}
}

func (s *FakeTokenMetadataService) GetMetadata(
	ctx context.Context,
	baseURL, ercId string,
) (*dto.TokenMetadataInfoDto, error) {
	if s.Error != nil {
		return nil, s.Error
	}
	if m, ok := s.Data[ercId]; ok {
		return m, nil
	}
	return nil, errors.New("metadata not found")
}

// FakeTransactionRepository is an in-memory implementation for testing transactions
type FakeTransactionRepository struct {
	Transactions      []domain.Transaction
	FlaggedTxs        []domain.FlaggedTransaction
	BatchTransactions []domain.Transaction
	EnygmaBatchTxs    []domain.Transaction
	Error             error
	NotFoundErr       error // allows injecting the correct sentinel
}

func NewFakeTransactionRepository() *FakeTransactionRepository {
	return &FakeTransactionRepository{
		Transactions:      []domain.Transaction{},
		FlaggedTxs:        []domain.FlaggedTransaction{},
		BatchTransactions: []domain.Transaction{},
		EnygmaBatchTxs:    []domain.Transaction{},
		NotFoundErr:       ErrRecordNotFound,
	}
}

func (r *FakeTransactionRepository) FindByMessageId(
	ctx context.Context,
	messageId string,
) (*domain.Transaction, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	for i := range r.Transactions {
		tx := &r.Transactions[i]
		if tx.MessageId == messageId {
			return tx, nil
		}
	}
	return nil, r.NotFoundErr
}

func (r *FakeTransactionRepository) FindByTransactionId(
	ctx context.Context,
	transactionId string,
) (*domain.Transaction, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	for i := range r.Transactions {
		tx := &r.Transactions[i]
		if tx.ID.String() == transactionId {
			return tx, nil
		}
	}
	return nil, r.NotFoundErr
}

func (r *FakeTransactionRepository) FindBySharedId(
	ctx context.Context,
	sharedId string,
) ([]domain.Transaction, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	var result []domain.Transaction
	for _, tx := range r.Transactions {
		if tx.SharedId == sharedId {
			result = append(result, tx)
		}
	}
	return result, nil
}

func (r *FakeTransactionRepository) FindByBatchId(
	ctx context.Context,
	batchId string,
) ([]domain.Transaction, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	return r.BatchTransactions, nil
}

func (r *FakeTransactionRepository) FindByBatchIdPaginated(
	ctx context.Context,
	batchId string,
	page, limit int,
) ([]domain.Transaction, int64, error) {
	if r.Error != nil {
		return nil, 0, r.Error
	}
	return r.BatchTransactions, int64(len(r.BatchTransactions)), nil
}

func (r *FakeTransactionRepository) FindByEnygmaBatchId(
	ctx context.Context,
	batchId string,
) ([]domain.Transaction, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	return r.EnygmaBatchTxs, nil
}

func (r *FakeTransactionRepository) FindByEnygmaBatchIdPaginated(
	ctx context.Context,
	batchId string,
	page, limit int,
) ([]domain.Transaction, int64, error) {
	if r.Error != nil {
		return nil, 0, r.Error
	}
	return r.EnygmaBatchTxs, int64(len(r.EnygmaBatchTxs)), nil
}

func (r *FakeTransactionRepository) FindWithFilters(
	ctx context.Context,
	filters dto.MergedTransactionsFilters,
) ([]domain.Transaction, int64, error) {
	if r.Error != nil {
		return nil, 0, r.Error
	}

	result := filterTransactions(r.Transactions, filters)
	return result, int64(len(result)), nil
}

// filterTransactions applies MergedTransactionsFilters to a slice of transactions
func filterTransactions(txs []domain.Transaction, filters dto.MergedTransactionsFilters) []domain.Transaction {
	result := make([]domain.Transaction, 0, len(txs))
	for _, tx := range txs {
		if !matchesTransactionFilters(tx, filters) {
			continue
		}
		result = append(result, tx)
	}
	return result
}

// matchesTransactionFilters checks if a transaction matches all provided filters
func matchesTransactionFilters(tx domain.Transaction, filters dto.MergedTransactionsFilters) bool {
	// Filter by fromChainId
	if filters.FromChainId != "" && tx.FromChainId.String() != filters.FromChainId {
		return false
	}
	// Filter by toChainId
	if filters.ToChainId != "" && tx.ToChainId.String() != filters.ToChainId {
		return false
	}
	// Filter by from address (case-insensitive)
	if filters.From != "" && !strings.EqualFold(tx.From, filters.From) {
		return false
	}
	// Filter by to address (case-insensitive)
	if filters.To != "" && !strings.EqualFold(tx.To, filters.To) {
		return false
	}
	// Filter by resourceId
	if filters.ResourceId != "" && tx.ResourceId != filters.ResourceId {
		return false
	}
	// Filter by messageId
	if filters.MessageId != "" && tx.MessageId != filters.MessageId {
		return false
	}
	// Filter by messageType
	if filters.MessageType != "" {
		msgType, ok := types.StringToAssetType[strings.ToLower(filters.MessageType)]
		if ok && tx.MsgType != msgType {
			return false
		}
	}
	return true
}

func (r *FakeTransactionRepository) FindFlagged(ctx context.Context) ([]domain.FlaggedTransaction, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	return r.FlaggedTxs, nil
}

// HasTransaction checks if a transaction with the given messageId exists (helper for assertions)
func (r *FakeTransactionRepository) HasTransaction(messageId string) bool {
	for _, tx := range r.Transactions {
		if tx.MessageId == messageId {
			return true
		}
	}
	return false
}

// FakePrivateNetworkRepository is an in-memory implementation for testing auth
type FakePrivateNetworkRepository struct {
	Networks    map[string]*domain.PrivateNetwork
	Error       error
	NotFoundErr error // allows injecting the correct sentinel
	ConflictErr error // inject core.ErrRecordConflict to simulate duplicates
}

func NewFakePrivateNetworkRepository() *FakePrivateNetworkRepository {
	return &FakePrivateNetworkRepository{
		Networks:    make(map[string]*domain.PrivateNetwork),
		NotFoundErr: ErrRecordNotFound,
	}
}

func (r *FakePrivateNetworkRepository) FindByUsername(
	ctx context.Context,
	username string,
) (*domain.PrivateNetwork, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	if pn, ok := r.Networks[username]; ok {
		return pn, nil
	}
	return nil, r.NotFoundErr
}

func (r *FakePrivateNetworkRepository) Create(ctx context.Context, username, hashedPassword string) error {
	if r.Error != nil {
		return r.Error
	}
	if _, exists := r.Networks[username]; exists && r.ConflictErr != nil {
		return r.ConflictErr
	}
	r.Networks[username] = &domain.PrivateNetwork{
		Username: username,
		Password: hashedPassword,
	}
	return nil
}

// FakeHeaderProofRepository is an in-memory implementation for testing header proofs
type FakeHeaderProofRepository struct {
	HeaderProofs []domain.HeaderProofEvent
	Error        error
}

func NewFakeHeaderProofRepository() *FakeHeaderProofRepository {
	return &FakeHeaderProofRepository{
		HeaderProofs: []domain.HeaderProofEvent{},
	}
}

func (r *FakeHeaderProofRepository) FindByBlockRange(
	ctx context.Context,
	chainId string,
	startBlock, endBlock int64,
	page, pageSize int,
) ([]domain.HeaderProofEvent, int64, error) {
	if r.Error != nil {
		return nil, 0, r.Error
	}

	result := make([]domain.HeaderProofEvent, 0, len(r.HeaderProofs))
	for _, hp := range r.HeaderProofs {
		// Filter by chainId
		if hp.ChainID.String() != chainId {
			continue
		}
		// Filter by block range
		blockNum := hp.BlockNumber.Unwrap().Int64()
		if blockNum < startBlock || blockNum > endBlock {
			continue
		}
		result = append(result, hp)
	}
	return result, int64(len(result)), nil
}

// HasProof checks if a header proof with the given block number exists (helper for assertions)
func (r *FakeHeaderProofRepository) HasProof(blockNumber int64) bool {
	for _, hp := range r.HeaderProofs {
		if hp.BlockNumber.Unwrap().Int64() == blockNumber {
			return true
		}
	}
	return false
}
