package repositories

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/flagger/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
)

// Ensure HeaderProofRepository implements core.HeaderProofRepository at compile time
var _ core.HeaderProofRepository = (*HeaderProofRepository)(nil)

// HeaderProofRepository is the adapter responsible for retrieving header proof events
type HeaderProofRepository struct {
	db *gorm.DB
}

// NewHeaderProofRepository creates a new HeaderProofRepository instance
func NewHeaderProofRepository(dbClient *gorm.DB) core.HeaderProofRepository {
	if dbClient == nil {
		panic("dbClient is nil")
	}
	return &HeaderProofRepository{db: dbClient}
}

// DeleteOlderThan deletes header proof events older than cutoff, preserving the most recent proof per chain.
func (r *HeaderProofRepository) DeleteOlderThan(ctx context.Context, cutoff time.Time) (int64, error) {
	result := r.db.WithContext(ctx).Exec(`
		DELETE FROM header_proof_events
		WHERE created_at < ?
		  AND id NOT IN (
		    SELECT id FROM header_proof_events
		    WHERE (chain_id, block_number) IN (
		      SELECT chain_id, MAX(block_number) FROM header_proof_events GROUP BY chain_id
		    )
		  )
	`, cutoff)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to purge header proof events: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// GetLatestHeaderProofs returns the latest header proof per chain using a subquery
func (r *HeaderProofRepository) GetLatestHeaderProofs(ctx context.Context) ([]domain.HeaderProofEvent, error) {
	var headers []domain.HeaderProofEvent

	// Use raw SQL with subquery to get latest header per chain
	err := r.db.WithContext(ctx).Raw(`
		SELECT h.*
		FROM header_proof_events h
		INNER JOIN (
			SELECT chain_id, MAX(block_number) as max_block
			FROM header_proof_events
			GROUP BY chain_id
		) latest ON h.chain_id = latest.chain_id AND h.block_number = latest.max_block
	`).Scan(&headers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get latest header proofs: %w", err)
	}

	return headers, nil
}
