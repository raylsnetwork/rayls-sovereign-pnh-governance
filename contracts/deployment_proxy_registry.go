package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	dpv1 "github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/DeploymentProxyRegistryV1"
)

type DeploymentProxyRegistryClient struct {
	address  common.Address
	contract *dpv1.DeploymentProxyRegistryV1
	backend  bind.ContractBackend
}

func NewDeploymentProxyRegistryClient(
	address common.Address,
	backend bind.ContractBackend,
) *DeploymentProxyRegistryClient {
	return &DeploymentProxyRegistryClient{
		address:  address,
		contract: dpv1.NewDeploymentProxyRegistryV1(),
		backend:  backend,
	}
}

func (c *DeploymentProxyRegistryClient) GetAllContracts(opts *bind.CallOpts) (dpv1.GetAllContractsOutput, error) {
	ctx := context.Background()
	if opts != nil && opts.Context != nil {
		ctx = opts.Context
	}
	msg := ethereum.CallMsg{To: &c.address, Data: c.contract.PackGetAllContracts()}
	raw, err := c.backend.CallContract(ctx, msg, nil)
	if err != nil {
		return dpv1.GetAllContractsOutput{}, err
	}
	return c.contract.UnpackGetAllContracts(raw)
}
