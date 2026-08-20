package repositories

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/infrastructure/database"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

var (
	txPool     *dockertest.Pool
	txResource *dockertest.Resource
	txDB       *gorm.DB
	txRepo     *TransactionRepository
	txOnce     sync.Once
)

func initTransactionTestEnv(t *testing.T) {
	t.Helper()
	txOnce.Do(func() {
		var err error
		txPool, err = dockertest.NewPool("")
		if err != nil {
			log.Fatalf("pool init: %v", err)
		}
		txResource, err = txPool.RunWithOptions(&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "16",
			Env: []string{
				"POSTGRES_PASSWORD=mysecretpassword",
				"POSTGRES_USER=postgres",
				"POSTGRES_DB=tx_repo_test",
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
		_ = txResource.Expire(240)
		txPool.MaxWait = 120 * time.Second

		dsn := "host=localhost user=postgres password=mysecretpassword dbname=tx_repo_test port=" + txResource.GetPort(
			"5432/tcp",
		) + " sslmode=disable"
		if err = txPool.Retry(func() error {
			txDB, err = database.Connect(dsn)
			if err != nil {
				return withstack.Wrap(err)
			}
			return nil
		}); err != nil {
			log.Fatalf("connect db: %v", err)
		}

		txRepo = NewTransactionRepository(txDB).(*TransactionRepository)
	})
}

func resetTxTables(t *testing.T) {
	t.Helper()
	txDB.Exec("DELETE FROM revert_data_transactions")
	txDB.Exec("DELETE FROM enygma_transactions")
	txDB.Exec("DELETE FROM transactions")
	txDB.Exec("DELETE FROM tokens")
}

func mustInsertToken(t *testing.T, id string) {
	t.Helper()
	token := domain.Token{
		ResourceId:  id,
		Name:        "Test Token",
		Symbol:      "TST",
		MetadataUrl: "https://example.invalid/token.json",
		ErcStandard: 0, // 0=ERC20, 1=ERC721, 2=ERC1155
		Decimals:    18,
		IssuerId:    "issuer-1",
		Status:      1,
	}
	if err := txDB.Create(&token).Error; err != nil {
		t.Fatalf("insert token: %v", err)
	}
}

func buildBaseTx(resourceId, msg string) *domain.Transaction {
	now := time.Now().UTC().Truncate(time.Second)
	return &domain.Transaction{
		ResourceId:           resourceId,
		MessageId:            msg,
		From:                 "0xfrom",
		To:                   "0xto",
		Protocol:             types.Atomic,
		SharedId:             "shared_1",
		BatchId:              "batch_1",
		TeleportStatus:       func() *uint8 { v := uint8(1); return &v }(),
		LogIndex:             42,
		Payload:              `{"ok":true}`,
		TxHashDestination:    "0xDEST",
		DestinationTimestamp: now,
		Amount:               decimal.NewFromInt(10),
		BlockNumber:          decimal.NewFromInt(100),
	}
}

func mustCreateTx(t *testing.T, tx *domain.Transaction) *domain.Transaction {
	if err := txRepo.CreateTransaction(context.Background(), tx); err != nil {
		t.Fatalf("create tx: %v", err)
	}
	return tx
}

func TestTransactionCRUD(t *testing.T) {
	initTransactionTestEnv(t)
	resetTxTables(t)
	resourceId := "res_123"
	mustInsertToken(t, resourceId)
	tx := buildBaseTx(resourceId, "msg_123")
	created := mustCreateTx(t, tx)
	assert.NotEqual(t, uuid.Nil, created.ID)

	gotByMsg, err := txRepo.GetTransactionByMessageID(context.Background(), "msg_123", false)
	assert.NoError(t, err)
	assert.Equal(t, resourceId, gotByMsg.ResourceId)
	assert.Equal(t, "msg_123", gotByMsg.MessageId)
	assert.Equal(t, created.ID, gotByMsg.ID)

	gotByMsg.IsFlagged = true
	gotByMsg.TxHashDestination = "0xUPDATED"
	err = txRepo.UpdateTransaction(context.Background(), gotByMsg)
	assert.NoError(t, err)

	refetched, err := txRepo.GetTransactionByMessageID(context.Background(), "msg_123", false)
	assert.NoError(t, err)
	assert.True(t, refetched.IsFlagged)
	assert.Equal(t, "0xUPDATED", refetched.TxHashDestination)
}

func TestEnygmaTransactionCRUD(t *testing.T) {
	initTransactionTestEnv(t)
	resetTxTables(t)
	resource := "res_e1"
	mustInsertToken(t, resource)
	baseTx := mustCreateTx(t, buildBaseTx(resource, "msg_e1"))
	data := &domain.EnygmaTransaction{
		TransactionId: baseTx.ID,
		ReferenceId:   "ref_1",
		UpdatedAt:     time.Now().UTC().Truncate(time.Second),
	}
	// Create
	assert.NoError(t, txDB.Create(&data).Error)
	assert.NotEqual(t, uuid.Nil, data.TransactionId)

	// Read
	var found domain.EnygmaTransaction
	assert.NoError(t, txDB.First(&found, "transaction_id = ?", baseTx.ID).Error)
	assert.Equal(t, "ref_1", found.ReferenceId)

	// Update
	assert.NoError(
		t,
		txDB.Model(&domain.EnygmaTransaction{}).
			Where("transaction_id = ?", baseTx.ID).
			Update("reference_id", "ref_2").
			Error,
	)
	var updated domain.EnygmaTransaction
	assert.NoError(t, txDB.First(&updated, "transaction_id = ?", baseTx.ID).Error)
	assert.Equal(t, "ref_2", updated.ReferenceId)

	// Delete
	assert.NoError(t, txDB.Delete(&domain.EnygmaTransaction{}, "transaction_id = ?", baseTx.ID).Error)
	var after domain.EnygmaTransaction
	err := txDB.First(&after, "transaction_id = ?", baseTx.ID).Error
	assert.Error(t, err)
}

func TestRevertDataTransactionCRUD(t *testing.T) {
	initTransactionTestEnv(t)
	resetTxTables(t)
	resource := "res_r1"
	mustInsertToken(t, resource)
	baseTx := mustCreateTx(t, buildBaseTx(resource, "msg_r1"))
	rev := &domain.RevertDataTransaction{
		TransactionId:                 baseTx.ID,
		TxHashDestinationRevert:       "0xDEST_REV",
		TxHashDestinationRevertStatus: 2,
		TxHashSourceRevert:            "0xSRC_REV",
		TxHashSourceRevertStatus:      3,
	}
	assert.NoError(t, txDB.Create(&rev).Error)
	assert.NotEqual(t, uuid.Nil, rev.TransactionId)

	var got domain.RevertDataTransaction
	assert.NoError(t, txDB.First(&got, "transaction_id = ?", baseTx.ID).Error)
	assert.Equal(t, "0xDEST_REV", got.TxHashDestinationRevert)

	assert.NoError(
		t,
		txDB.Model(&domain.RevertDataTransaction{}).
			Where("transaction_id = ?", baseTx.ID).
			Update("tx_hash_destination_revert_status", 4).
			Error,
	)
	var refetched domain.RevertDataTransaction
	assert.NoError(t, txDB.First(&refetched, "transaction_id = ?", baseTx.ID).Error)
	assert.Equal(t, uint8(4), refetched.TxHashDestinationRevertStatus)

	assert.NoError(t, txDB.Delete(&domain.RevertDataTransaction{}, "transaction_id = ?", baseTx.ID).Error)
	var after domain.RevertDataTransaction
	err := txDB.First(&after, "transaction_id = ?", baseTx.ID).Error
	assert.Error(t, err)
}
