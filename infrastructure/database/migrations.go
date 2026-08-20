package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

func RunMigrations(connStr string) error {
	connStr = normalizeSSLMode(connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return withstack.Wrap(err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.Error("failed to close db", "err", closeErr)
		}
	}()

	driver, err := migratePostgres.WithInstance(db, &migratePostgres.Config{})
	if err != nil {
		return withstack.Wrap(err)
	}

	migrationsPath, err := findProjectMigrationsPath()
	if err != nil {
		return withstack.Wrap(err)
	}

	logger.Info("Running migrations", "path", migrationsPath)

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver,
	)
	if err != nil {
		return withstack.Wrap(err)
	}
	defer func() {
		sourceError, destinationError := m.Close()
		if sourceError != nil {
			logger.Error("failed to close db (source)", "err", sourceError)
		}
		if destinationError != nil {
			logger.Error("failed to close db (destination)", "err", destinationError)
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("All migrations already applied")
			return nil
		}
		return withstack.Wrap(err)
	}

	logger.Info("All pending migrations applied successfully")

	return nil
}

func normalizeSSLMode(connStr string) string {
	connStr = strings.ReplaceAll(connStr, "sslmode=enable", "")
	if !strings.Contains(connStr, "sslmode=") {
		connStr += " sslmode=disable"
	}
	return connStr
}

func findProjectMigrationsPath() (string, error) {
	startDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := startDir
	for {
		migrationsPath := filepath.Join(dir, "migrations")
		if _, err := os.Stat(migrationsPath); err == nil {
			return migrationsPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find migrations folder from: %s", startDir)
		}
		dir = parent
	}
}
