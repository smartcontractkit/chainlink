package changeset

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

type LinkTransfer struct {
	From  common.Address
	To    common.Address
	Value big.Int
}
type LinkTransferTimelockRequest struct {
	Transfers        []LinkTransfer
	ChainSelector    uint64
	LinkTokenAddress common.Address
	TimelockAddress  common.Address
	MCMSAddress      common.Address
	ValidUntil       uint32        // unix time until the proposal will be valid
	MinDelay         time.Duration // delay for timelock worker to execute the transfers.
	OverrideRoot     bool
}

var _ deployment.ChangeSet[*LinkTransferTimelockRequest] = LinkTransferTimelock

// LinkTransferTimelock takes the given link transfers and creates an MCMS proposal for them.
func LinkTransferTimelock(env deployment.Environment, req *LinkTransferTimelockRequest) (deployment.ChangesetOutput, error) {
	chain := env.Chains[req.ChainSelector]
	link, err := link_token.NewLinkToken(req.LinkTokenAddress, chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	chainID := mcms.ChainIdentifier(req.ChainSelector)
	chainMetadata := map[mcms.ChainIdentifier]mcms.ChainMetadata{
		chainID: {MCMAddress: req.MCMSAddress, StartingOpCount: 0},
	}
	timelockAddresses := map[mcms.ChainIdentifier]common.Address{
		chainID: req.TimelockAddress,
	}
	batch := timelock.BatchChainOperation{
		ChainIdentifier: chainID,
		Batch:           []mcms.Operation{},
	}
	for _, transfer := range req.Transfers {
		transferTx, err := link.TransferFrom(chain.DeployerKey, transfer.From, transfer.To, &transfer.Value)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error constructing transfer tx: %w", err)
		}
		op := mcms.Operation{
			To:           req.LinkTokenAddress,
			Data:         transferTx.Data(),
			Value:        big.NewInt(0),
			ContractType: "LinkToken",
		}
		batch.Batch = append(batch.Batch, op)

	}
	proposal, err := timelock.NewMCMSWithTimelockProposal(
		"1",
		req.ValidUntil,
		[]mcms.Signature{},
		req.OverrideRoot,
		chainMetadata,
		timelockAddresses,
		"Value transfer proposal",
		[]timelock.BatchChainOperation{batch}, // 1 batch with all the transfers on the same batch
		timelock.Schedule,
		req.MinDelay.String(),
	)
	if err != nil {
		return deployment.ChangesetOutput{}, nil
	}
	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*proposal},
	}, nil

}
