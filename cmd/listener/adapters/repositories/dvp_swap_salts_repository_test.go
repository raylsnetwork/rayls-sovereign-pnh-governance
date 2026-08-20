package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/infrastructure/database"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

var (
	swapSaltsPool     *dockertest.Pool
	swapSaltsResource *dockertest.Resource
	swapSaltsDB       *gorm.DB
	swapSaltsRepo     *DvpSwapSaltsRepository
	swapSaltsOnce     sync.Once
)

func initSwapSaltsTestEnv(t *testing.T) {
	t.Helper()
	swapSaltsOnce.Do(func() {
		var err error
		swapSaltsPool, err = dockertest.NewPool("")
		if err != nil {
			log.Fatalf("pool init: %v", err)
		}
		swapSaltsResource, err = swapSaltsPool.RunWithOptions(&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "16",
			Env: []string{
				"POSTGRES_PASSWORD=mysecretpassword",
				"POSTGRES_USER=postgres",
				"POSTGRES_DB=swap_salts_repo_test",
				"POSTGRES_DB_PORT=5432",
				"listen_addresses = '*'",
			},
		}, func(hc *docker.HostConfig) {
			hc.AutoRemove = true
			hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})
		if err != nil {
			log.Fatalf("start resource: %v", err)
		}
		_ = swapSaltsResource.Expire(240)
		swapSaltsPool.MaxWait = 120 * time.Second

		dsn := fmt.Sprintf(
			"host=localhost user=postgres password=mysecretpassword dbname=swap_salts_repo_test port=%s sslmode=disable",
			swapSaltsResource.GetPort("5432/tcp"),
		)

		if err = swapSaltsPool.Retry(func() error {
			swapSaltsDB, err = database.Connect(dsn)
			if err != nil {
				return withstack.Wrap(err)
			}
			return nil
		}); err != nil {
			log.Fatalf("connect db: %v", err)
		}

		swapSaltsRepo = NewDvpSwapSaltsRepository(swapSaltsDB)
	})
}

func resetSwapSalts(t *testing.T) {
	t.Helper()
	if err := swapSaltsDB.Exec("DELETE FROM dvp_swap_salts").Error; err != nil {
		t.Fatalf("failed to reset dvp_swap_salts: %v", err)
	}
}

func TestDvpSwapSaltsRepository_Put_InsertsNewRow(t *testing.T) {
	// Putting salts for a previously-unseen sharedID inserts a single row.
	initSwapSaltsTestEnv(t)
	resetSwapSalts(t)

	salts := types.DvpSwapSalts{
		InitiatorSelfSalt: []byte("self-1"),
		InitiatorCtxtSalt: []byte("ctxt-1"),
	}
	err := swapSaltsRepo.Put(context.Background(), "abc", salts)

	require.NoError(t, err)
	var rows []domain.DvpSwapSalts
	require.NoError(t, swapSaltsDB.Find(&rows).Error)
	require.Len(t, rows, 1)
	assert.Equal(t, "abc", rows[0].SharedID)
	assert.Equal(t, []byte("self-1"), rows[0].InitiatorSelfSalt)
	assert.Equal(t, []byte("ctxt-1"), rows[0].InitiatorCtxtSalt)
}

func TestDvpSwapSaltsRepository_Put_SameSaltsAreNoop(t *testing.T) {
	// Putting identical salts twice keeps exactly one row and leaves it unchanged.
	initSwapSaltsTestEnv(t)
	resetSwapSalts(t)

	salts := types.DvpSwapSalts{
		InitiatorSelfSalt: []byte("self-1"),
		InitiatorCtxtSalt: []byte("ctxt-1"),
	}
	require.NoError(t, swapSaltsRepo.Put(context.Background(), "abc", salts))
	require.NoError(t, swapSaltsRepo.Put(context.Background(), "abc", salts))

	var count int64
	require.NoError(t, swapSaltsDB.Model(&domain.DvpSwapSalts{}).Count(&count).Error)
	assert.Equal(t, int64(1), count)
}

func TestDvpSwapSaltsRepository_Put_ReplacesDifferentSalts(t *testing.T) {
	// Putting different salts for the same sharedID overwrites the existing values.
	initSwapSaltsTestEnv(t)
	resetSwapSalts(t)

	require.NoError(t, swapSaltsRepo.Put(context.Background(), "abc", types.DvpSwapSalts{
		InitiatorSelfSalt: []byte("self-1"),
		InitiatorCtxtSalt: []byte("ctxt-1"),
	}))
	require.NoError(t, swapSaltsRepo.Put(context.Background(), "abc", types.DvpSwapSalts{
		InitiatorSelfSalt: []byte("self-2"),
		InitiatorCtxtSalt: []byte("ctxt-2"),
	}))

	got, err := swapSaltsRepo.Get(context.Background(), "abc")
	require.NoError(t, err)
	assert.Equal(t, []byte("self-2"), got.InitiatorSelfSalt)
	assert.Equal(t, []byte("ctxt-2"), got.InitiatorCtxtSalt)
}

func TestDvpSwapSaltsRepository_Get_ReturnsStoredSalts(t *testing.T) {
	// Get returns the salts previously stored under the given sharedID.
	initSwapSaltsTestEnv(t)
	resetSwapSalts(t)

	require.NoError(t, swapSaltsRepo.Put(context.Background(), "xyz", types.DvpSwapSalts{
		InitiatorSelfSalt: []byte("deadbeef"),
		InitiatorCtxtSalt: []byte("cafebabe"),
	}))

	got, err := swapSaltsRepo.Get(context.Background(), "xyz")

	require.NoError(t, err)
	assert.Equal(t, []byte("deadbeef"), got.InitiatorSelfSalt)
	assert.Equal(t, []byte("cafebabe"), got.InitiatorCtxtSalt)
}

func TestDvpSwapSaltsRepository_Get_ReturnsErrSwapSaltsNotFoundForUnknownSharedID(t *testing.T) {
	// Get maps GORM's ErrRecordNotFound to core.ErrSwapSaltsNotFound so the
	// handler can detect foreign swaps and skip them.
	initSwapSaltsTestEnv(t)
	resetSwapSalts(t)

	_, err := swapSaltsRepo.Get(context.Background(), "missing")

	assert.True(t, errors.Is(err, core.ErrSwapSaltsNotFound))
}
