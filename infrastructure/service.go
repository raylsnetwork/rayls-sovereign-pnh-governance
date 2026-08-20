package infrastructure

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/config"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts"
	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
)

func CreateClient(url string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// LoadContractAddresses dynamically loads contract addresses from DeploymentProxyRegistry
// and populates them into the config
func LoadContractAddresses(ctx context.Context, client *ethclient.Client, conf *config.Config) error {
	if conf.PrivateHub.DeploymentProxyRegistry == "" {
		return fmt.Errorf("DeploymentProxyRegistry address not configured")
	}

	deploymentProxyInstance, err := contracts.CreateDeploymentProxyRegistry(
		ctx,
		conf.PrivateHub.DeploymentProxyRegistry,
		client,
	)
	if err != nil {
		return fmt.Errorf("failed to create deployment proxy instance: %w", err)
	}

	// Get all contracts from the registry
	namesAndAddresses, err := deploymentProxyInstance.GetAllContracts(nil)
	if err != nil {
		return fmt.Errorf("failed to get contracts from registry: %w", err)
	}

	// Build address map name -> address
	addresses := make(map[string]common.Address)
	for index, name := range namesAndAddresses.Names {
		if index < len(namesAndAddresses.Addresses) {
			addresses[name] = namesAndAddresses.Addresses[index]
		}
	}

	// Populate all known contract addresses
	if addr, exists := addresses["Teleport"]; exists {
		conf.PrivateHub.Teleport = addr.Hex()
	}
	if addr, exists := addresses["EnygmaPNHEvents"]; exists {
		conf.PrivateHub.EnygmaPNHEvents = addr.Hex()
	}
	if addr, exists := addresses["EnygmaTeleport"]; exists {
		conf.PrivateHub.EnygmaTeleport = addr.Hex()
	}
	if addr, exists := addresses["EnygmaTokenManager"]; exists {
		conf.PrivateHub.EnygmaTokenManager = addr.Hex()
	}
	if addr, exists := addresses["TokenCore"]; exists {
		conf.PrivateHub.TokenCore = addr.Hex()
	}
	if addr, exists := addresses["TokenFreezeManager"]; exists {
		conf.PrivateHub.TokenFreezeManager = addr.Hex()
	}
	if addr, exists := addresses["ParticipantCore"]; exists {
		conf.PrivateHub.ParticipantCore = addr.Hex()
	}
	if addr, exists := addresses["AuditManager"]; exists {
		conf.PrivateHub.AuditManager = addr.Hex()
	}
	if addr, exists := addresses["ParticipantStorage"]; exists {
		conf.PrivateHub.ParticipantStorageContract = addr.Hex()
	}
	if addr, exists := addresses["TokenRegistry"]; exists {
		conf.PrivateHub.TokenRegistry = addr.Hex()
	}
	if addr, exists := addresses["DvpTeleport"]; exists {
		conf.PrivateHub.DvpTeleport = addr.Hex()
	}
	if addr, exists := addresses["Proofs"]; exists {
		conf.PrivateHub.ProofsAddress = addr.Hex()
	}

	logger.Info("Loaded contract addresses from registry",
		"Proofs", conf.PrivateHub.ProofsAddress,
		"Teleport", conf.PrivateHub.Teleport,
		"TokenCore", conf.PrivateHub.TokenCore,
		"ParticipantCore", conf.PrivateHub.ParticipantCore,
	)

	return nil
}
