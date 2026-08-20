package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/api/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
)

const pgUniqueViolation = "23505"

var _ core.PrivateNetworkRepository = (*privateNetworkRepository)(nil)

// privateNetworkRepository implements core.PrivateNetworkRepository
type privateNetworkRepository struct {
	db *gorm.DB
}

// NewPrivateNetworkRepository creates a new private network repository
func NewPrivateNetworkRepository(db *gorm.DB) core.PrivateNetworkRepository {
	return &privateNetworkRepository{db: db}
}

// FindByUsername retrieves a private network operator by username
func (r *privateNetworkRepository) FindByUsername(
	ctx context.Context,
	username string,
) (*domain.PrivateNetwork, error) {
	var pn domain.PrivateNetwork

	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&pn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrRecordNotFound
		}
		return nil, err
	}

	return &pn, nil
}

// Create creates a new private network operator
func (r *privateNetworkRepository) Create(ctx context.Context, username, hashedPassword string) error {
	pn := &domain.PrivateNetwork{
		Username: username,
		Password: hashedPassword,
	}

	err := r.db.WithContext(ctx).Create(pn).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return core.ErrRecordConflict
		}
		return err
	}
	return nil
}
