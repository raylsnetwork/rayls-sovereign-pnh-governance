package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/TokenRegistryV1"
)

type TokenRegistryClient struct {
	address  common.Address
	contract *TokenRegistryV1.TokenRegistryV1
	backend  bind.ContractBackend
}

func NewTokenRegistryClient(address common.Address, backend bind.ContractBackend) *TokenRegistryClient {
	return &TokenRegistryClient{
		address:  address,
		contract: TokenRegistryV1.NewTokenRegistryV1(),
		backend:  backend,
	}
}

func (c *TokenRegistryClient) GetTokenByResourceId(
	opts *bind.CallOpts,
	resourceId [32]byte,
) (TokenRegistryV1.TokenStructsToken, error) {
	ctx := context.Background()
	if opts != nil && opts.Context != nil {
		ctx = opts.Context
	}
	msg := ethereum.CallMsg{To: &c.address, Data: c.contract.PackGetTokenByResourceId(resourceId)}
	raw, err := c.backend.CallContract(ctx, msg, nil)
	if err != nil {
		return TokenRegistryV1.TokenStructsToken{}, err
	}
	return c.contract.UnpackGetTokenByResourceId(raw)
}
