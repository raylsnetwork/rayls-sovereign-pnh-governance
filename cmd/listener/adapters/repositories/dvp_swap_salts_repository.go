package repositories

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

var _ core.SwapSaltsStore = (*DvpSwapSaltsRepository)(nil)

type DvpSwapSaltsRepository struct {
	db *gorm.DB
}

func NewDvpSwapSaltsRepository(dbClient *gorm.DB) *DvpSwapSaltsRepository {
	if dbClient == nil {
		panic("dbClient is nil")
	}
	return &DvpSwapSaltsRepository{db: dbClient}
}

// Put persists the salts for sharedID. If a row already exists with identical
// salts it is a no-op; otherwise the stored salts are replaced.
//
// The read-then-UPSERT is not atomic: two concurrent Put calls with matching
// salts can both pass the equality short-circuit and both issue the UPSERT.
// The end state is correct (UPSERT is idempotent for identical values) and
// we rely on the NATS consumer to serialize DvpTeleport events per sharedID,
// so this races only in pathological redelivery scenarios.
func (r *DvpSwapSaltsRepository) Put(ctx context.Context, sharedID string, salts types.DvpSwapSalts) error {
	var existing domain.DvpSwapSalts
	err := r.db.WithContext(ctx).Where("shared_id = ?", sharedID).First(&existing).Error
	if err == nil &&
		bytes.Equal(existing.InitiatorSelfSalt, salts.InitiatorSelfSalt) &&
		bytes.Equal(existing.InitiatorCtxtSalt, salts.InitiatorCtxtSalt) {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to read swap salts for shared_id %s: %w", sharedID, err)
	}

	row := domain.DvpSwapSalts{
		SharedID:          sharedID,
		InitiatorSelfSalt: salts.InitiatorSelfSalt,
		InitiatorCtxtSalt: salts.InitiatorCtxtSalt,
	}
	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "shared_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"initiator_self_salt", "initiator_ctxt_salt"}),
		}).
		Create(&row).Error; err != nil {
		return fmt.Errorf("failed to persist swap salts for shared_id %s: %w", sharedID, err)
	}
	return nil
}

// Get returns the salts stored for sharedID or core.ErrSwapSaltsNotFound.
func (r *DvpSwapSaltsRepository) Get(ctx context.Context, sharedID string) (types.DvpSwapSalts, error) {
	var row domain.DvpSwapSalts
	if err := r.db.WithContext(ctx).Where("shared_id = ?", sharedID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.DvpSwapSalts{}, core.ErrSwapSaltsNotFound
		}
		return types.DvpSwapSalts{}, fmt.Errorf("failed to read swap salts for shared_id %s: %w", sharedID, err)
	}
	return types.DvpSwapSalts{
		InitiatorSelfSalt: row.InitiatorSelfSalt,
		InitiatorCtxtSalt: row.InitiatorCtxtSalt,
	}, nil
}
