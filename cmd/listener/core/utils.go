package core

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

// SafeUnixTime converts a uint64 timestamp to time.Time, capping at max int64
func SafeUnixTime(timestamp uint64) time.Time {
	if timestamp > math.MaxInt64 {
		return time.Unix(math.MaxInt64, 0)
	}
	return time.Unix(int64(timestamp), 0)
}

// DecryptPayload decrypts the given payload using the provided Decryptor and configuration, unmarshals the decrypted bytes into the generic type T
func DecryptPayload[T any](
	d Decryptor,
	payload []byte,
	blockNumber uint64,
	pnData PNodeDataAndSecrets,
	secretType types.SecretType,
) (T, error) {
	var emptyResult T

	decryptedBytes, err := d.DecryptPayloadBytes(payload, blockNumber, pnData, secretType)
	if err != nil {
		return emptyResult, withstack.Wrap(err)
	}

	var result T
	if err := json.Unmarshal(decryptedBytes, &result); err != nil {
		return emptyResult, fmt.Errorf("failed to unmarshal decrypted data: %w", err)
	}
	return result, nil
}

// IDExtractor extracts an ID string from a generic data item
type IDExtractor[T any] func(T) string

// TransactionIDType (msg_id or shared_id)
type TransactionIDType int

const (
	MessageID TransactionIDType = iota
	SharedID
)

// GetKey returns the appropriate key from a transaction based on the TransactionIDType
func (t TransactionIDType) GetKey(tx *domain.Transaction) string {
	switch t {
	case MessageID:
		return tx.MessageId
	case SharedID:
		return tx.SharedId
	default:
		return ""
	}
}

// BuildTransactionMap creates a transaction map for fast lookup
func BuildTransactionMap[T any](
	ctx context.Context,
	persister TransactionRepository,
	decryptedData []T,
	extractor IDExtractor[T],
	idType TransactionIDType,
) (map[string]*domain.Transaction, error) {
	if len(decryptedData) == 0 {
		return map[string]*domain.Transaction{}, nil
	}

	// Extract IDs using the provided extractor function
	ids := make([]string, 0, len(decryptedData))
	for _, data := range decryptedData {
		ids = append(ids, extractor(data))
	}

	// Fetch transactions based on identifier type
	var transactions []domain.Transaction
	var err error

	// Retrieve transactions in the database based on identifier type
	switch idType {
	case MessageID:
		transactions, err = persister.GetTransactionsByMessageIDs(ctx, ids, false)
	case SharedID:
		transactions, err = persister.GetTransactionsBySharedIDs(ctx, ids, false)
	default:
		return nil, fmt.Errorf("unsupported identifier type: %d", idType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	// Build the transaction map for quick lookup and further processing
	transactionMap := make(map[string]*domain.Transaction)
	for i := range transactions {
		transaction := &transactions[i]
		key := idType.GetKey(transaction)
		transactionMap[key] = transaction
	}

	return transactionMap, nil
}

// StringToDomainBigInt converts a string to our domain.BigInt
func StringToDomainBigInt(s string) domain.BigInt {
	id, _ := new(big.Int).SetString(s, 10)
	return domain.NewBigInt(id)
}

// dvpProtocolFromUpdateType maps a DVP balance update type to the corresponding protocol:
// Mint (withdrawal from DVP) → DvpWithdraw; Burn (deposit into DVP) → DvpDeposit.
func DvpProtocolFromUpdateType(updateType types.TxType) types.ProtocolType {
	if updateType == types.Mint {
		return types.DvpWithdraw
	}
	return types.DvpDeposit
}
