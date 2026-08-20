package repositories

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/flagger/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
)

const pgUniqueViolation = "23505"

// Ensure BalanceRepository implements core.BalanceRepository at compile time
var _ core.BalanceRepository = (*BalanceRepository)(nil)

// BalanceRepository is the adapter responsible for retrieving transactions from the database for flagging operations
type BalanceRepository struct {
	db *gorm.DB
}

// NewBalanceRepository creates a new BalanceRepository instance for managing transaction retrieval
func NewBalanceRepository(dbClient *gorm.DB) core.BalanceRepository {
	if dbClient == nil {
		panic("dbClient is nil")
	}
	return &BalanceRepository{db: dbClient}
}

// CreateBalance creates a single balance (idempotent based on message_id)
func (r *BalanceRepository) CreateBalance(ctx context.Context, balance *domain.Balance) error {
	if balance == nil {
		return fmt.Errorf("failed to create transaction: balance is nil")
	}

	err := getTransaction(ctx, r.db).WithContext(ctx).Create(balance).Error
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// GetResourceBalance retrieves the balance for a given resource on a specific chain.
// If the balance doesn't exist yet, it atomically inserts a zero-value row via
// FirstOrCreate and returns it, ensuring the row is always available for subsequent updates.
func (br *BalanceRepository) GetResourceBalance(
	ctx context.Context,
	chainId string,
	resourceId string,
	ercId string,
) (*domain.Balance, error) {
	var balance domain.Balance
	defaults, err := newEmptyBalance(resourceId, chainId, ercId)
	if err != nil {
		return nil, fmt.Errorf("failed to build default balance: %w", err)
	}

	query := getTransaction(ctx, br.db).WithContext(ctx).
		Where("chain_id = ? AND resource_id = ?", chainId, resourceId)

	if ercId == "" {
		query = query.Where("erc_id IS NULL")
	} else {
		query = query.Where("erc_id = ?", ercId)
	}

	result := query.Attrs(defaults).FirstOrCreate(&balance)
	if result.Error != nil {
		// On concurrent insert, the unique constraint rejects the duplicate.
		// Retry with a plain SELECT to fetch the row the other goroutine created.
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == pgUniqueViolation {
			if err := query.First(&balance).Error; err != nil {
				return nil, fmt.Errorf(
					"failed to re-fetch balance after conflict (chain_id=%s) (resource_id=%s): %w",
					chainId,
					resourceId,
					err,
				)
			}
			return &balance, nil
		}
		return nil, fmt.Errorf(
			"failed to get or create balance (chain_id=%s) (resource_id=%s): %w",
			chainId,
			resourceId,
			result.Error,
		)
	}

	return &balance, nil
}

// UpdateSenderReceiverBalances updates the balances of sender and receiver to the specified new amounts in a single transaction
func (br *BalanceRepository) UpdateSenderReceiverBalances(
	ctx context.Context,
	senderChainId string,
	senderResourceId string,
	senderErcId domain.BigInt,
	senderNewAmount string,
	receiverChainId string,
	receiverResourceId string,
	receiverErcId domain.BigInt,
	receiverNewAmount string,
) error {
	// Convert BigInt to string for SQL (nil stays nil)
	var senderErcIdVal, receiverErcIdVal any
	if senderErcId.Int != nil {
		senderErcIdVal = senderErcId.String()
	}
	if receiverErcId.Int != nil {
		receiverErcIdVal = receiverErcId.String()
	}

	query := `
		UPDATE balances 
		SET amount = CASE 
			WHEN chain_id = $1 AND resource_id = $2 AND erc_id IS NOT DISTINCT FROM $3::numeric THEN $4::numeric
			WHEN chain_id = $5 AND resource_id = $6 AND erc_id IS NOT DISTINCT FROM $7::numeric THEN $8::numeric
		END,
		updated_at = $9
		WHERE (chain_id = $1 AND resource_id = $2 AND erc_id IS NOT DISTINCT FROM $3::numeric)
		   OR (chain_id = $5 AND resource_id = $6 AND erc_id IS NOT DISTINCT FROM $7::numeric)
	`

	result := getTransaction(ctx, br.db).WithContext(ctx).Exec(query,
		senderChainId, senderResourceId, senderErcIdVal, senderNewAmount,
		receiverChainId, receiverResourceId, receiverErcIdVal, receiverNewAmount,
		time.Now(),
	)

	if result.Error != nil {
		return fmt.Errorf("failed to update sender/receiver balances: %w", result.Error)
	}

	if result.RowsAffected != 2 {
		return fmt.Errorf("expected 2 rows to be updated, got %d", result.RowsAffected)
	}

	return nil
}

// UpdateBalance updates the amount of a specific balance identified by chain_id, resource_id and erc_id
func (br *BalanceRepository) UpdateBalance(
	ctx context.Context,
	chainId string,
	resourceId string,
	ercId domain.BigInt,
	amount string,
) error {
	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	query := getTransaction(ctx, br.db).WithContext(ctx).Model(&domain.Balance{}).
		Where("chain_id = ? AND resource_id = ?", chainId, resourceId)

	if ercId.Int == nil {
		query = query.Where("erc_id IS NULL")
	} else {
		query = query.Where("erc_id = ?", ercId)
	}

	if err := query.Updates(map[string]any{
		"amount":     amountDecimal,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return err
	}

	return nil
}

// newEmptyBalance creates a zeroed balance for a given resource, chain and ercId string.
// Returns an error if ercId is non-empty but not a valid integer.
func newEmptyBalance(resourceId string, chainId string, ercId string) (*domain.Balance, error) {
	var ercBigInt *big.Int
	if ercId != "" {
		var ok bool
		ercBigInt, ok = new(big.Int).SetString(ercId, 10)
		if !ok {
			return nil, fmt.Errorf("invalid erc_id %q: not a valid integer", ercId)
		}
	}

	return &domain.Balance{
		ResourceId: resourceId,
		ChainId:    chainId,
		Amount:     decimal.Zero,
		ErcId:      domain.BigInt{Int: ercBigInt},
	}, nil
}
