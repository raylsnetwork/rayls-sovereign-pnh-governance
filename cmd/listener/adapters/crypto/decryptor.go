package crypto

import (
	"context"
	"crypto/mlkem"
	"fmt"
	"math/big"
	"slices"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/listener/core"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/config"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/ParticipantStorageV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/cryptography"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/types"
)

// Ensure Decryptor implements core.Decryptor at compile time
var _ core.Decryptor = (*Decryptor)(nil)

// participantKeyInfo holds per-participant decapsulation keys recovered during phase 1.
type participantKeyInfo struct {
	chainId  string
	blockNum string
	dk       *mlkem.DecapsulationKey768
}

// Decryptor provides a concrete implementation of the core.Decryptor interface.
type Decryptor struct {
	client *ethclient.Client
	config *config.Config
	log    logger.Logger

	operatorDKOnce sync.Once
	operatorDK     *mlkem.DecapsulationKey768
	operatorDKErr  error
}

// NewDecryptor creates a new decryptor adapter with injected ethclient, config, and logger
func NewDecryptor(client *ethclient.Client, conf *config.Config, log logger.Logger) core.Decryptor {
	return &Decryptor{
		client: client,
		config: conf,
		log:    log,
	}
}

// getOperatorDK returns the operator's ML-KEM decapsulation key, parsing it
// on the first call and caching the result for subsequent calls.
func (d *Decryptor) getOperatorDK() (*mlkem.DecapsulationKey768, error) {
	d.operatorDKOnce.Do(func() {
		d.operatorDK, d.operatorDKErr = cryptography.ImportDecapsulationKey(
			d.config.PrivateHub.RaylsViewSecretKey,
		)
	})
	if d.operatorDKErr != nil {
		return nil, fmt.Errorf("failed to import operator ML-KEM key: %w", d.operatorDKErr)
	}
	return d.operatorDK, nil
}

// GatherParticipantsData fetches all cryptographic data required for decryption from the ParticipantStorage smart contract.
// It recovers shared secrets via ML-KEM decapsulation of on-chain key agreement ciphertexts.
func (d *Decryptor) GatherParticipantsData(ctx context.Context, conf *config.Config) (core.PNodeDataAndSecrets, error) {
	participantsContract, err := contracts.CreateParticipantStorage(
		ctx,
		conf.PrivateHub.ParticipantStorageContract,
		d.client,
	)
	if err != nil {
		return nil, err
	}

	operatorDK, err := d.getOperatorDK()
	if err != nil {
		return nil, err
	}

	data, err := participantsContract.GetParticipantDataBatch(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get participant data batch: %w", err)
	}

	// Fetch key agreements directed at operator (ciphertexts from participants → operator)
	operatorChainId := big.NewInt(conf.PrivateHub.OperatorChainId)
	operatorKeyAgreements, err := participantsContract.GetKeyAgreements(nil, operatorChainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get operator key agreements: %w", err)
	}

	// Build lookup: participantChainId → key agreement data (for operator)
	operatorKAByChain := buildKAByInitiator(operatorKeyAgreements)
	pnDataAndSecretsForChainAndBlock := make(map[string]map[string]*types.IPNodeDataAndSecrets)
	participantKeys := make([]participantKeyInfo, 0, len(data.AuditInfo))

	// Phase 1: Decrypt participant private keys and recover operator-participant shared secrets
	for _, auditData := range data.AuditInfo {
		info := d.processAuditEntry(
			auditData,
			conf.PrivateHub.OperatorChainId,
			operatorDK,
			operatorKAByChain,
			pnDataAndSecretsForChainAndBlock,
		)
		if info != nil {
			participantKeys = append(participantKeys, *info)
		}
	}

	// Phase 2: Compute pairwise participant-to-participant shared secrets.
	// Key agreements are established once per pair (whichever participant does it first).
	// Both directions must use the same secret, so we only decapsulate the direction
	// from the oldest block and skip the other.

	// Sort participants by block number ascending so the oldest key agreement is processed first.
	slices.SortFunc(participantKeys, func(a, b participantKeyInfo) int {
		blockA, _ := new(big.Int).SetString(a.blockNum, 10)
		blockB, _ := new(big.Int).SetString(b.blockNum, 10)
		return blockA.Cmp(blockB)
	})

	// Track which cross-participant pairs already have a secret (canonical key: "smallId-largeId")
	resolvedPairs := make(map[string]bool)
	for _, pkB := range participantKeys {
		chainIdB, _ := new(big.Int).SetString(pkB.chainId, 10)
		keyAgreements, err := participantsContract.GetKeyAgreements(nil, chainIdB)
		if err != nil {
			d.log.Warn("Failed to get key agreements for participant", "chainId", pkB.chainId, "error", err)
			continue
		}

		// Build lookup of ciphertexts directed at B: initiatorChainId → ciphertext
		kaByInitiator := buildKAByInitiator(keyAgreements)
		for _, pkA := range participantKeys {
			d.processPairSecret(pkA, pkB, kaByInitiator, resolvedPairs, pnDataAndSecretsForChainAndBlock)
		}
	}

	return pnDataAndSecretsForChainAndBlock, nil
}

// processAuditEntry decrypts a single participant's private key and recovers the operator-participant shared secret.
// Returns nil if the entry should be skipped or any cryptographic step fails.
func (d *Decryptor) processAuditEntry(
	auditData ParticipantStorageV1.ParticipantStructsAuditInfoData,
	operatorChainId int64,
	operatorDK *mlkem.DecapsulationKey768,
	operatorKAByChain map[string][]byte,
	pnData map[string]map[string]*types.IPNodeDataAndSecrets,
) *participantKeyInfo {
	if auditData.ChainId.Int64() == operatorChainId || auditData.ChainId.Int64() == 0 {
		d.log.Debug("Skipping OperatorChainId:", "chainId", auditData.ChainId.String())
		return nil
	}

	// Find operator's key agreement ciphertext for this participant
	ciphertext, ok := operatorKAByChain[auditData.ChainId.String()]
	if !ok {
		d.log.Error("No key agreement found for participant", "chainId", auditData.ChainId.String())
		return nil
	}

	// Decapsulate to recover shared secret with this participant
	sharedSecret, err := cryptography.Decapsulate(operatorDK, ciphertext)
	if err != nil {
		d.log.Error("Failed to decapsulate key agreement", "chainId", auditData.ChainId.String(), "error", err)
		return nil
	}

	participantDKBytes, err := cryptography.DecryptAuditPrivateKey(
		auditData.EncryptedRaylsViewPrivateKey,
		auditData.Mac,
		sharedSecret,
	)
	if err != nil {
		d.log.Error("Failed to decrypt participant private key", "chainId", auditData.ChainId.String(), "error", err)
		return nil
	}

	// Reconstruct participant's ML-KEM decapsulation key
	participantDK, err := mlkem.NewDecapsulationKey768(participantDKBytes)
	if err != nil {
		d.log.Error("Failed to reconstruct participant ML-KEM key", "chainId", auditData.ChainId.String(), "error", err)
		return nil
	}

	chainIdStr := auditData.ChainId.String()
	blockNumStr := auditData.BlockNumber.String()

	if pnData[chainIdStr] == nil {
		pnData[chainIdStr] = make(map[string]*types.IPNodeDataAndSecrets)
	}
	pnData[chainIdStr][blockNumStr] = &types.IPNodeDataAndSecrets{
		ChainId:         auditData.ChainId,
		BlockNumber:     auditData.BlockNumber,
		HubSharedSecret: sharedSecret,
		ParticipantDK:   participantDKBytes,
	}

	return &participantKeyInfo{chainId: chainIdStr, blockNum: blockNumStr, dk: participantDK}
}

// processPairSecret computes the pairwise shared secret between pkA and pkB and stores it in pnData.
// resolvedPairs prevents the same cross-participant pair from being processed twice.
func (d *Decryptor) processPairSecret(
	pkA, pkB participantKeyInfo,
	kaByInitiator map[string][]byte,
	resolvedPairs map[string]bool,
	pnDataAndSecretsForChainAndBlock map[string]map[string]*types.IPNodeDataAndSecrets,
) {
	ciphertext, ok := kaByInitiator[pkA.chainId]
	if !ok {
		return
	}

	// For cross-participant pairs, skip if already resolved from the older block
	if pkA.chainId != pkB.chainId {
		canonKey := pkA.chainId + "-" + pkB.chainId
		if pkB.chainId < pkA.chainId {
			canonKey = pkB.chainId + "-" + pkA.chainId
		}
		if resolvedPairs[canonKey] {
			return
		}
		resolvedPairs[canonKey] = true
	}

	pairwiseSecret, err := cryptography.Decapsulate(pkB.dk, ciphertext)
	if err != nil {
		d.log.Warn("Failed to decapsulate pairwise secret", "from", pkA.chainId, "to", pkB.chainId, "error", err)
		return
	}

	entry := pnDataAndSecretsForChainAndBlock[pkB.chainId][pkB.blockNum]
	secretData := &types.IPNodeDataAndSecrets{
		ChainId:      entry.ChainId,
		BlockNumber:  entry.BlockNumber,
		SharedSecret: pairwiseSecret,
	}

	// Store under A-B
	pairKey := pkA.chainId + "-" + pkB.chainId
	if pnDataAndSecretsForChainAndBlock[pairKey] == nil {
		pnDataAndSecretsForChainAndBlock[pairKey] = make(map[string]*types.IPNodeDataAndSecrets)
	}
	pnDataAndSecretsForChainAndBlock[pairKey][entry.BlockNumber.String()] = secretData

	// For cross-participant pairs, also store under B-A
	if pkA.chainId != pkB.chainId {
		reverseKey := pkB.chainId + "-" + pkA.chainId
		if pnDataAndSecretsForChainAndBlock[reverseKey] == nil {
			pnDataAndSecretsForChainAndBlock[reverseKey] = make(map[string]*types.IPNodeDataAndSecrets)
		}
		pnDataAndSecretsForChainAndBlock[reverseKey][entry.BlockNumber.String()] = secretData
	}
}

// buildKAByInitiator builds a chainId→ciphertext lookup from a list of key agreements.
func buildKAByInitiator(kas []ParticipantStorageV1.ParticipantStructsKeyAgreementData) map[string][]byte {
	m := make(map[string][]byte, len(kas))
	for _, ka := range kas {
		m[ka.ChainId.String()] = ka.Ciphertext
	}
	return m
}

// DecryptPayloadBytes decrypts a given payload using the appropriate set of shared secrets
func (d *Decryptor) DecryptPayloadBytes(
	payload []byte,
	blockNumber uint64,
	pnData core.PNodeDataAndSecrets,
	secretType types.SecretType,
) ([]byte, error) {
	var sharedSecrets [][]byte

	switch secretType {
	case types.ParticipantSecret:
		sharedSecrets = getParticipantsSharedSecrets(blockNumber, pnData)
	case types.AtomicSecret:
		sharedSecrets = getAtomicSharedSecrets(blockNumber, pnData)
	default:
		return nil, fmt.Errorf("unrecognized secretType (%d): cannot determine key derivation strategy", secretType)
	}

	for _, secret := range sharedSecrets {
		sharedKey, err := cryptography.DeriveSymmetricKey(secret)
		if err != nil {
			continue
		}
		_, decryptedData, err := cryptography.DecryptGCM(payload, sharedKey)
		if err == nil {
			return decryptedData, nil
		}
	}

	return nil, fmt.Errorf("no key could decrypt the payload")
}

// DecryptSwapPayload decrypts a SwapInitiated payload. The ciphertext is
// encapsulated against the destination-chain participant's view public key,
// but the event doesn't identify which participant that is. So we iterate
// every participant DK whose audit entry was recorded at or before
// blockNumber and try decapsulation + GCM decryption with each; the first
// success wins. Returns the plaintext and the salt reduced from the
// Ctxt-derived shared secret (InitiatorCtxtSalt) so the caller can persist it.
func (d *Decryptor) DecryptSwapPayload(
	ctxt []byte,
	encryptedData []byte,
	blockNumber uint64,
	pnData core.PNodeDataAndSecrets,
) ([]byte, []byte, error) {
	candidates := filterValidDataAndSecrets(pnData, blockNumber)

	var tried int
	for _, byBlock := range candidates {
		for _, entry := range byBlock {
			if len(entry.ParticipantDK) == 0 {
				continue
			}
			tried++
			dk, err := mlkem.NewDecapsulationKey768(entry.ParticipantDK)
			if err != nil {
				continue
			}
			secret, err := cryptography.Decapsulate(dk, ctxt)
			if err != nil {
				continue
			}
			initiatorCtxtSalt := cryptography.ReduceToSalt(secret)
			plaintext, err := d.DecryptWithSalt(encryptedData, initiatorCtxtSalt)
			if err != nil {
				continue
			}
			return plaintext, initiatorCtxtSalt, nil
		}
	}
	return nil, nil, fmt.Errorf("no participant DK could decrypt swap payload (tried %d)", tried)
}

// DecryptWithSalt decrypts a GCM payload using a caller-supplied salt,
// applying the same HKDF derivation used elsewhere.
func (d *Decryptor) DecryptWithSalt(payload []byte, salt []byte) ([]byte, error) {
	key, err := cryptography.DeriveSymmetricKey(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to derive symmetric key: %w", err)
	}
	_, plaintext, err := cryptography.DecryptGCM(payload, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt payload: %w", err)
	}
	return plaintext, nil
}

func getParticipantsSharedSecrets(blockNumber uint64, pnDataAndSecrets core.PNodeDataAndSecrets) [][]byte {
	pnDataAndSecretsForChainAndBlock := filterValidDataAndSecrets(pnDataAndSecrets, blockNumber)
	var sharedSecrets [][]byte

	for _, pnDataAndSecretsPerChain := range pnDataAndSecretsForChainAndBlock {
		for _, pnDataAndSecretsPerBlock := range pnDataAndSecretsPerChain {
			if pnDataAndSecretsPerBlock.SharedSecret != nil {
				sharedSecrets = append(sharedSecrets, pnDataAndSecretsPerBlock.SharedSecret)
			}
			// Also include HubSharedSecret as it may be used for participant decryption in some cases
			if pnDataAndSecretsPerBlock.HubSharedSecret != nil {
				sharedSecrets = append(sharedSecrets, pnDataAndSecretsPerBlock.HubSharedSecret)
			}
		}
	}

	return sharedSecrets
}

func getAtomicSharedSecrets(blockNumber uint64, pnDataAndSecrets core.PNodeDataAndSecrets) [][]byte {
	pnDataAndSecretsForChainAndBlock := filterValidDataAndSecrets(pnDataAndSecrets, blockNumber)

	var sharedSecrets [][]byte
	for _, pnDataAndSecretsMap := range pnDataAndSecretsForChainAndBlock {
		for _, pnDataAndSecret := range pnDataAndSecretsMap {
			if pnDataAndSecret.HubSharedSecret != nil {
				sharedSecrets = append(sharedSecrets, pnDataAndSecret.HubSharedSecret)
			}
		}
	}

	return sharedSecrets
}

func filterValidDataAndSecrets(
	pnDataAndSecretsForChainAndBlock core.PNodeDataAndSecrets,
	currentBlockNumber uint64,
) map[string]map[string]*types.IPNodeDataAndSecrets {
	validData := make(map[string]map[string]*types.IPNodeDataAndSecrets)

	for chainId, pnDataAndSecretsForBlock := range pnDataAndSecretsForChainAndBlock {
		var maxBlockNumber *big.Int
		var maxBlockData *types.IPNodeDataAndSecrets

		for _, data := range pnDataAndSecretsForBlock {
			if data.BlockNumber.Cmp(big.NewInt(0).SetUint64(currentBlockNumber)) <= 0 {
				if maxBlockNumber == nil || data.BlockNumber.Cmp(maxBlockNumber) > 0 {
					maxBlockNumber = data.BlockNumber
					maxBlockData = data
				}
			}
		}

		if maxBlockData != nil {
			if validData[chainId] == nil {
				validData[chainId] = make(map[string]*types.IPNodeDataAndSecrets)
			}
			validData[chainId][maxBlockNumber.String()] = maxBlockData
		}
	}

	return validData
}
