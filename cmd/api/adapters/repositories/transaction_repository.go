package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/utils"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/dto"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

var _ core.TransactionRepository = (*transactionRepository)(nil)

// transactionRepository implements core.TransactionRepository
type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new GORM-based transaction repository
func NewTransactionRepository(db *gorm.DB) core.TransactionRepository {
	return &transactionRepository{db: db}
}

// FindByMessageId finds a single transaction by message_id field
func (r *transactionRepository) FindByMessageId(ctx context.Context, messageId string) (*domain.Transaction, error) {
	var transaction domain.Transaction

	err := r.db.WithContext(ctx).
		Preload("Token").
		Preload("EnygmaTransaction").
		Where("message_id = ? AND message_id != ''", messageId).
		First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrRecordNotFound
		}
		return nil, err
	}

	return &transaction, nil
}

// FindByTransactionId finds a single transaction by transaction_id field
func (r *transactionRepository) FindByTransactionId(
	ctx context.Context,
	transactionId string,
) (*domain.Transaction, error) {
	var transaction domain.Transaction

	err := r.db.WithContext(ctx).
		Preload("Token").
		Where("id = ?", transactionId).
		First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrRecordNotFound
		}
		return nil, err
	}

	return &transaction, nil
}

// FindByBatchId finds transactions by batch_id field (regular batches)
func (r *transactionRepository) FindByBatchId(ctx context.Context, batchId string) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	err := r.db.WithContext(ctx).
		Preload("Token").
		Find(&transactions, "batch_id = ? AND batch_id != ''", batchId).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// FindBySharedId finds transactions by shared_id field (dvp swaps)
func (r *transactionRepository) FindBySharedId(ctx context.Context, sharedId string) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	err := r.db.WithContext(ctx).
		Preload("Token").
		Find(&transactions, "shared_id = ?", sharedId).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// FindByBatchIdPaginated finds transactions by aggregation_key for regular_batch with pagination
func (r *transactionRepository) FindByBatchIdPaginated(
	ctx context.Context,
	batchId string,
	page, limit int,
) ([]domain.Transaction, int64, error) {
	var total int64

	// Count total
	err := r.db.WithContext(ctx).
		Model(&domain.Transaction{}).
		Where("aggregation_key = ? AND aggregation_type = 'regular_batch'", batchId).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []domain.Transaction{}, 0, nil
	}

	// Fetch page
	offset := (page - 1) * limit
	var transactions []domain.Transaction
	err = r.db.WithContext(ctx).
		Preload("Token").
		Where("aggregation_key = ? AND aggregation_type = 'regular_batch'", batchId).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// FindByEnygmaBatchId finds transactions by enygma_transactions.batch_id
func (r *transactionRepository) FindByEnygmaBatchId(ctx context.Context, batchId string) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	err := r.db.WithContext(ctx).
		Preload("Token").
		Preload("EnygmaTransaction").
		Joins("JOIN enygma_transactions ON transactions.id = enygma_transactions.transaction_id").
		Where("transactions.batch_id = ?", batchId).
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// FindByEnygmaBatchIdPaginated finds enygma transactions by batch_id with pagination
func (r *transactionRepository) FindByEnygmaBatchIdPaginated(
	ctx context.Context,
	batchId string,
	page, limit int,
) ([]domain.Transaction, int64, error) {
	var total int64

	// Count total
	err := r.db.WithContext(ctx).
		Model(&domain.Transaction{}).
		Joins("JOIN enygma_transactions ON transactions.id = enygma_transactions.transaction_id").
		Where("transactions.batch_id = ?", batchId).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []domain.Transaction{}, 0, nil
	}

	// Fetch page
	offset := (page - 1) * limit
	var transactions []domain.Transaction
	err = r.db.WithContext(ctx).
		Preload("Token").
		Preload("EnygmaTransaction").
		Joins("JOIN enygma_transactions ON transactions.id = enygma_transactions.transaction_id").
		Where("transactions.batch_id = ?", batchId).
		Order("transactions.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// FindWithFilters finds transactions matching the provided filters with DB-level pagination.
// Returns: transactions (one per aggregation_key), total count, error
func (r *transactionRepository) FindWithFilters(
	ctx context.Context,
	filters dto.MergedTransactionsFilters,
) ([]domain.Transaction, int64, error) {
	// Build WHERE clause and args
	whereClause, args := r.buildWhereClause(filters)

	// Get total count of distinct aggregation_keys
	var total int64
	countSQL := `SELECT COUNT(DISTINCT aggregation_key) FROM transactions WHERE ` + whereClause
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []domain.Transaction{}, 0, nil
	}

	// Calculate offset
	offset := (filters.Page - 1) * filters.Limit

	// Query to get one representative transaction per aggregation_key with pagination
	dataSQL := `
		SELECT * FROM (
			SELECT DISTINCT ON (aggregation_key) *
			FROM transactions
			WHERE ` + whereClause + `
			ORDER BY aggregation_key, created_at DESC
		) AS aggregated
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	args = append(args, filters.Limit, offset)

	var transactions []domain.Transaction
	if err := r.db.WithContext(ctx).Raw(dataSQL, args...).Scan(&transactions).Error; err != nil {
		return nil, 0, err
	}

	if err := r.preloadRelatedData(ctx, transactions); err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// FindFlagged retrieves all flagged transactions
func (r *transactionRepository) FindFlagged(ctx context.Context) ([]domain.FlaggedTransaction, error) {
	flagged := make([]domain.FlaggedTransaction, 0)

	err := r.db.WithContext(ctx).
		Raw("SELECT * FROM flagged_transactions").
		Scan(&flagged).Error
	if err != nil {
		return nil, err
	}

	return flagged, nil
}

// preloadRelatedData loads tokens, enygma transactions, and revert data for the given transactions
func (r *transactionRepository) preloadRelatedData(ctx context.Context, transactions []domain.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	ids := make([]uuid.UUID, len(transactions))
	resourceIdSet := make(map[string]bool)
	for i, tx := range transactions {
		ids[i] = tx.ID
		resourceIdSet[tx.ResourceId] = true
	}

	resourceIds := make([]string, 0, len(resourceIdSet))
	for id := range resourceIdSet {
		resourceIds = append(resourceIds, id)
	}

	// Fetch tokens
	var tokens []domain.Token
	if err := r.db.WithContext(ctx).Where("resource_id IN ?", resourceIds).Find(&tokens).Error; err != nil {
		return err
	}
	tokenMap := make(map[string]domain.Token)
	for _, t := range tokens {
		tokenMap[t.ResourceId] = t
	}

	// Fetch enygma transactions
	var enygmaTxs []domain.EnygmaTransaction
	if err := r.db.WithContext(ctx).Where("transaction_id IN ?", ids).Find(&enygmaTxs).Error; err != nil {
		return err
	}
	enygmaMap := make(map[uuid.UUID]*domain.EnygmaTransaction)
	for i := range enygmaTxs {
		enygmaMap[enygmaTxs[i].TransactionId] = &enygmaTxs[i]
	}

	// Fetch revert data transactions
	var revertTxs []domain.RevertDataTransaction
	if err := r.db.WithContext(ctx).Where("transaction_id IN ?", ids).Find(&revertTxs).Error; err != nil {
		return err
	}
	revertMap := make(map[uuid.UUID]*domain.RevertDataTransaction)
	for i := range revertTxs {
		revertMap[revertTxs[i].TransactionId] = &revertTxs[i]
	}

	// Attach related data to transactions
	for i := range transactions {
		if token, ok := tokenMap[transactions[i].ResourceId]; ok {
			transactions[i].Token = token
		}
		if enygma, ok := enygmaMap[transactions[i].ID]; ok {
			transactions[i].EnygmaTransaction = enygma
		}
		if revert, ok := revertMap[transactions[i].ID]; ok {
			transactions[i].RevertDataTransaction = revert
		}
	}

	return nil
}

// buildWhereClause constructs the WHERE clause and arguments for transaction queries
func (r *transactionRepository) buildWhereClause(filters dto.MergedTransactionsFilters) (string, []interface{}) {
	conditions := []string{
		// Exclude Burn (0) and Mint (1) transactions that are not related with Dvp flow
		"(tx_type = 2 OR protocol IS NOT NULL)",
	}
	var args []any

	if filters.FromChainId != "" {
		conditions = append(conditions, "from_chain_id = ?")
		args = append(args, filters.FromChainId)
	}

	if filters.ToChainId != "" {
		conditions = append(conditions, "to_chain_id = ?")
		args = append(args, filters.ToChainId)
	}

	if filters.From != "" {
		conditions = append(conditions, `LOWER("from") = LOWER(?)`)
		args = append(args, filters.From)
	}

	if filters.To != "" {
		conditions = append(conditions, `LOWER("to") = LOWER(?)`)
		args = append(args, filters.To)
	}

	if filters.MessageId != "" {
		conditions = append(conditions, "message_id = ?")
		args = append(args, filters.MessageId)
	}

	if filters.ResourceId != "" {
		conditions = append(conditions, "resource_id = ?")
		args = append(args, filters.ResourceId)
	}

	if filters.MessageType != "" {
		if msgTypeInt, ok := types.StringToAssetType[filters.MessageType]; ok {
			conditions = append(conditions, "msg_type = ?")
			args = append(args, uint8(msgTypeInt))
		}
	}

	if filters.InitiatedAfter != "" {
		afterTime, _ := utils.ParseTime(filters.InitiatedAfter)
		conditions = append(conditions, "created_at > ?")
		args = append(args, afterTime)
	}

	if filters.InitiatedBefore != "" {
		beforeTime, _ := utils.ParseTime(filters.InitiatedBefore)
		conditions = append(conditions, "created_at < ?")
		args = append(args, beforeTime)
	}

	return strings.Join(conditions, " AND "), args
}
