package database

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/config"
	appLogger "github.com/raylsnetwork/rayls-privacy-pnh-governance-api/logger"
)

func setupDockerPostgres(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	// Initialize logger for tests
	appLogger.InitializeLogger(&config.Config{Logging: "info"})

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=test_db",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}

	port := resource.GetPort("5432/tcp")
	connStr := fmt.Sprintf("host=localhost user=postgres password=secret dbname=test_db port=%s sslmode=disable", port)

	pool.MaxWait = 60 * time.Second
	var db *gorm.DB
	if retryErr := pool.Retry(func() error {
		db, err = Connect(connStr)
		if err != nil {
			log.Println("waiting for container...")
			return err
		}
		return nil
	}); retryErr != nil {
		t.Fatalf("Could not connect to docker postgres: %s", retryErr)
	}

	cleanup := func() {
		_ = pool.Purge(resource)
	}

	return db, cleanup
}

func TestParseDSN(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:  "GORM DSN - valid",
			input: "host=localhost user=postgres password=secret dbname=mydb port=5432",
			expected: map[string]string{
				"host":     "localhost",
				"user":     "postgres",
				"password": "secret",
				"dbname":   "mydb",
				"port":     "5432",
			},
		},
		{
			name:  "JDBC URI - valid",
			input: "jdbc:postgresql://localhost:5432/mydb",
			expected: map[string]string{
				"host":   "localhost",
				"port":   "5432",
				"dbname": "mydb",
			},
		},
		{
			name:     "JDBC URI - malformed",
			input:    "jdbc:postgresql://localhost5432mydb",
			expected: map[string]string{
				// Should be empty or partial
			},
		},
		{
			name:  "GORM DSN - malformed",
			input: "host:localhost user=postgres password",
			expected: map[string]string{
				"user": "postgres",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDSN(tt.input)
			for k, v := range tt.expected {
				assert.Equal(t, v, result[k], "expected key %s to be %s", k, v)
			}
		})
	}
}

func TestConnect(t *testing.T) {
	t.Run("[SUCCESS] connect to docker postgres and run automigrate", func(t *testing.T) {
		db, cleanup := setupDockerPostgres(t)
		defer cleanup()

		assert.NotNil(t, db)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		assert.NoError(t, sqlDB.Ping())

		expectedTables := []string{
			"tokens",
			"transactions",
			"enygma_transactions",
			"revert_data_transactions",
			"participants",
			"balances",
			"flagged_transactions",
			"last_processed_block",
			"header_proof_events",
			"header_flag_events",
		}

		var tableNames []string
		err = db.Raw(`
			SELECT tablename 
			FROM pg_tables 
			WHERE schemaname = 'public';
		`).Scan(&tableNames).Error
		assert.NoError(t, err)

		for _, expected := range expectedTables {
			assert.Contains(t, tableNames, expected, "table %s should exist", expected)
		}
	})

	t.Run("[ERROR] fail with invalid host", func(t *testing.T) {
		db, err := Connect("host=invalidhost user=postgres password=pass dbname=test port=5432 sslmode=disable")
		assert.Nil(t, db)
		assert.Error(t, err)
	})
}

func TestCreateDatabase_MissingParams(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]string
	}{
		{"missing host", map[string]string{
			"port": "5432", "user": "postgres", "password": "pass", "dbname": "test",
		}},
		{"missing dbname", map[string]string{
			"host": "localhost", "port": "5432", "user": "postgres", "password": "pass",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := createDatabase(tt.params)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "missing required DSN param")
		})
	}
}

func TestDisconnect(t *testing.T) {
	db, cleanup := setupDockerPostgres(t)
	defer cleanup()

	sqlDB, err := db.DB()
	assert.NoError(t, err)

	err = Disconnect(db)
	assert.NoError(t, err)

	err = sqlDB.Ping()
	assert.Error(t, err)
}
