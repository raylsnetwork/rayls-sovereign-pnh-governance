package contracts

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/raylsnetwork/rayls-sovereign-pnh-governance/contracts/ParticipantStorageV1"
)

type ParticipantStorageClient struct {
	address  common.Address
	contract *ParticipantStorageV1.ParticipantStorageV1
	backend  bind.ContractBackend
}

func NewParticipantStorageClient(address common.Address, backend bind.ContractBackend) *ParticipantStorageClient {
	return &ParticipantStorageClient{
		address:  address,
		contract: ParticipantStorageV1.NewParticipantStorageV1(),
		backend:  backend,
	}
}

func (c *ParticipantStorageClient) call(opts *bind.CallOpts, calldata []byte) ([]byte, error) {
	ctx := context.Background()
	if opts != nil && opts.Context != nil {
		ctx = opts.Context
	}
	msg := ethereum.CallMsg{To: &c.address, Data: calldata}
	return c.backend.CallContract(ctx, msg, nil)
}

func (c *ParticipantStorageClient) GetParticipantDataBatch(
	opts *bind.CallOpts,
) (ParticipantStorageV1.GetParticipantDataBatchOutput, error) {
	raw, err := c.call(opts, c.contract.PackGetParticipantDataBatch())
	if err != nil {
		return ParticipantStorageV1.GetParticipantDataBatchOutput{}, err
	}
	return c.contract.UnpackGetParticipantDataBatch(raw)
}

func (c *ParticipantStorageClient) GetKeyAgreements(
	opts *bind.CallOpts,
	chainId *big.Int,
) ([]ParticipantStorageV1.ParticipantStructsKeyAgreementData, error) {
	raw, err := c.call(opts, c.contract.PackGetKeyAgreements(chainId))
	if err != nil {
		return nil, err
	}
	return c.contract.UnpackGetKeyAgreements(raw)
}
