package repositories

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/config"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/infrastructure/database"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

var (
	testPool     *dockertest.Pool
	testResource *dockertest.Resource
	testDB       *gorm.DB
	testRepo     *LastProcessedBlockRepository
)

// TestMain sets up a single Postgres container reused across all tests in this file.
func TestMain(m *testing.M) {
	// Initialize logger for tests
	logger.InitializeLogger(&config.Config{Logging: "development"})

	var err error
	testPool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	testResource, err = testPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16",
		Env: []string{
			"POSTGRES_PASSWORD=mysecretpassword",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=block_repo_test",
			"POSTGRES_DB_PORT=5432",
			"listen_addresses = '*'",
		},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	_ = testResource.Expire(240)
	testPool.MaxWait = 120 * time.Second

	port := testResource.GetPort("5432/tcp")
	dsn := fmt.Sprintf(
		"host=localhost user=postgres password=mysecretpassword dbname=block_repo_test port=%s sslmode=disable",
		port,
	)

	if err = testPool.Retry(func() error {
		testDB, err = database.Connect(dsn)
		if err != nil {
			log.Println("Postgres not ready yet... retrying")
			return withstack.Wrap(err)
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to postgres: %s", err)
	}

	testRepo = NewLastProcessedBlockRepository(testDB).(*LastProcessedBlockRepository)

	code := m.Run()

	// Cleanup container
	if err := testPool.Purge(testResource); err != nil {
		log.Printf("failed to purge docker resource: %v", err)
	}
	os.Exit(code)
}

// helper to truncate table between subtests
func resetTable(t *testing.T) {
	t.Helper()
	testDB.Exec("DELETE FROM last_processed_block")
}

// helper to insert a block row
func insertBlock(t *testing.T, number string) {
	t.Helper()
	bi := new(big.Int)
	bi.SetString(number, 10)
	if err := testDB.Create(&domain.LastProcessedBlock{ID: 1, Number: domain.NewBigInt(bi)}).Error; err != nil {
		t.Fatalf("failed inserting block: %v", err)
	}
}

func fetchBlock(t *testing.T) *domain.LastProcessedBlock {
	t.Helper()
	var b domain.LastProcessedBlock
	if err := testDB.First(&b).Error; err != nil {
		t.Fatalf("fetch block failed: %v", err)
	}
	return &b
}

func countBlocks(t *testing.T) int64 {
	t.Helper()
	var c int64
	if err := testDB.Model(&domain.LastProcessedBlock{}).Count(&c).Error; err != nil {
		t.Fatalf("count blocks failed: %v", err)
	}
	return c
}

func TestGetLatestProcessedBlock(t *testing.T) {
	resetTable(t)
	t.Run("success returns latest processed block", func(t *testing.T) {
		resetTable(t)
		insertBlock(t, "12345")
		block, err := testRepo.GetLatestProcessedBlock(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, block)
		assert.Equal(t, "12345", block.String())
	})

	t.Run("error when no record", func(t *testing.T) {
		resetTable(t)
		block, err := testRepo.GetLatestProcessedBlock(context.Background())
		assert.Error(t, err)
		assert.Nil(t, block)
		assert.Contains(t, err.Error(), "no last processed block found in database")
		assert.Contains(t, err.Error(), "record not found")
	})
}

func TestUpdateLatestProcessedBlock(t *testing.T) {
	cases := []struct {
		name          string
		setup         func(t *testing.T)
		newValue      *big.Int
		expect        string
		expectCount   int64
		secondUpdates []*big.Int // optional chain of additional updates to validate id stability
	}{
		{
			name:        "create when empty",
			setup:       func(t *testing.T) { resetTable(t) },
			newValue:    big.NewInt(1000),
			expect:      "1000",
			expectCount: 1,
		},
		{
			name:        "update existing",
			setup:       func(t *testing.T) { resetTable(t); insertBlock(t, "1000") },
			newValue:    big.NewInt(2000),
			expect:      "2000",
			expectCount: 1,
		},
		{
			name:          "multiple sequential updates keep single row",
			setup:         func(t *testing.T) { resetTable(t) },
			newValue:      big.NewInt(1),
			secondUpdates: []*big.Int{big.NewInt(2), big.NewInt(5000), big.NewInt(5000)},
			expect:        "5000",
			expectCount:   1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(t)
			err := testRepo.UpdateLatestProcessedBlock(context.Background(), tc.newValue)
			assert.NoError(t, err)
			for _, upd := range tc.secondUpdates {
				err = testRepo.UpdateLatestProcessedBlock(context.Background(), upd)
				assert.NoError(t, err)
			}
			blk := fetchBlock(t)
			assert.Equal(t, tc.expect, blk.Number.Unwrap().String())
			assert.Equal(t, uint(1), blk.ID)
			assert.Equal(t, tc.expectCount, countBlocks(t))
		})
	}
}
