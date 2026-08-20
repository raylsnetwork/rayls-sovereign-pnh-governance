package repositories

import (
	"context"

	"gorm.io/gorm"
)

type txContextKey struct{}

func withTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

func getTransaction(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return defaultDB
}
