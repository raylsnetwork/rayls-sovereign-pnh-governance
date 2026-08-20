package contracts

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/AuditManagerV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/DvpTeleport"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/EnygmaTeleport"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/EnygmaTokenManagerV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/ParticipantCoreV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/Proofs"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/TeleportV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/TokenCoreV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/TokenFreezeManagerV1"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/withstack"
)

// Event-only contracts — no address/client needed, only used for log parsing.

func CreateTeleport() *TeleportV1.TeleportV1 {
	return TeleportV1.NewTeleportV1()
}

func CreateTokenFreezeManager() *TokenFreezeManagerV1.TokenFreezeManagerV1 {
	return TokenFreezeManagerV1.NewTokenFreezeManagerV1()
}

func CreateProofs() *Proofs.Proofs {
	return Proofs.NewProofs()
}

func CreateEnygmaTeleport() *EnygmaTeleport.EnygmaTeleport {
	return EnygmaTeleport.NewEnygmaTeleport()
}

func CreateTokenCore() *TokenCoreV1.TokenCoreV1 {
	return TokenCoreV1.NewTokenCoreV1()
}

func CreateEnygmaTokenManager() *EnygmaTokenManagerV1.EnygmaTokenManagerV1 {
	return EnygmaTokenManagerV1.NewEnygmaTokenManagerV1()
}

func CreateParticipantCore() *ParticipantCoreV1.ParticipantCoreV1 {
	return ParticipantCoreV1.NewParticipantCoreV1()
}

func CreateAuditManager() *AuditManagerV1.AuditManagerV1 {
	return AuditManagerV1.NewAuditManagerV1()
}

func CreateDvpTeleport() *DvpTeleport.DvpTeleport {
	return DvpTeleport.NewDvpTeleport()
}

// RPC-calling contracts — need address and client for on-chain reads.

func CreateParticipantStorage(
	ctx context.Context,
	address string,
	client bind.ContractBackend,
) (*ParticipantStorageClient, error) {
	contractAddress := common.HexToAddress(address)
	code, err := client.CodeAt(ctx, contractAddress, nil)
	if err != nil {
		return nil, withstack.Wrap(err)
	}
	if len(code) == 0 {
		return nil, fmt.Errorf("no contract code at address %s", contractAddress.Hex())
	}
	return NewParticipantStorageClient(contractAddress, client), nil
}

func CreateTokenRegistry(
	ctx context.Context,
	address string,
	client bind.ContractBackend,
) (*TokenRegistryClient, error) {
	contractAddress := common.HexToAddress(address)
	code, err := client.CodeAt(ctx, contractAddress, nil)
	if err != nil {
		return nil, withstack.Wrap(err)
	}
	if len(code) == 0 {
		return nil, fmt.Errorf("no contract code at address %s", contractAddress.Hex())
	}
	return NewTokenRegistryClient(contractAddress, client), nil
}

func CreateDeploymentProxyRegistry(
	ctx context.Context,
	address string,
	client bind.ContractBackend,
) (*DeploymentProxyRegistryClient, error) {
	contractAddress := common.HexToAddress(address)
	code, err := client.CodeAt(ctx, contractAddress, nil)
	if err != nil {
		return nil, withstack.Wrap(err)
	}
	if len(code) == 0 {
		return nil, fmt.Errorf("no contract code at address %s", contractAddress.Hex())
	}
	return NewDeploymentProxyRegistryClient(contractAddress, client), nil
}

// Keep ethclient import used — some callers pass *ethclient.Client which satisfies bind.ContractBackend.
var _ bind.ContractBackend = (*ethclient.Client)(nil)
