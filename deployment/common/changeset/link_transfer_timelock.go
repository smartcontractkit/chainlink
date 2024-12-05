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
	To    common.Address
	Value big.Int
}
type LinkTransferTimelockRequest struct {
	Transfers        map[uint64][]LinkTransfer
	LinkTokenAddress map[uint64]common.Address
	TimelockAddress  map[uint64]common.Address
	MCMSAddress      map[uint64]common.Address
	ValidUntil       uint32        // unix time until the proposal will be valid
	MinDelay         time.Duration // delay for timelock worker to execute the transfers.
	OverrideRoot     bool
	StartingOpCount  map[uint64]uint64
}

var _ deployment.ChangeSet[*LinkTransferTimelockRequest] = LinkTransferTimelock

// LinkTransferTimelock takes the given link transfers and creates an MCMS proposal for them.
func LinkTransferTimelock(e deployment.Environment, req *LinkTransferTimelockRequest) (deployment.ChangesetOutput, error) {
	chainMetadata := map[mcms.ChainIdentifier]mcms.ChainMetadata{}
	timelockAddresses := map[mcms.ChainIdentifier]common.Address{}
	allBatches := []timelock.BatchChainOperation{}
	for chainSelector := range req.Transfers {
		if _, ok := req.LinkTokenAddress[chainSelector]; !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("missing link token address for chain %d", chainSelector)
		}
		if _, ok := req.TimelockAddress[chainSelector]; !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("missing timelock address for chain %d", chainSelector)
		}
		if _, ok := req.MCMSAddress[chainSelector]; !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("missing MCMS address for chain %d", chainSelector)
		}
		chainID := mcms.ChainIdentifier(chainSelector)
		chain := e.Chains[chainSelector]
		linkContract, err := link_token.NewLinkToken(req.LinkTokenAddress[chainSelector], chain.Client)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get link contract: %w", err)
		}
		chainMetadata[chainID] = mcms.ChainMetadata{
			MCMAddress:      req.MCMSAddress[chainSelector],
			StartingOpCount: req.StartingOpCount[chainSelector],
		}
		timelockAddresses[chainID] = req.TimelockAddress[chainSelector]
		batch := timelock.BatchChainOperation{
			ChainIdentifier: chainID,
			Batch:           []mcms.Operation{},
		}

		for _, transfer := range req.Transfers[chainSelector] {
			tx, err := linkContract.Transfer(deployment.SimTransactOpts(), transfer.To, &transfer.Value)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("error packing transfer tx data: %w", err)
			}
			op := mcms.Operation{
				To:           req.LinkTokenAddress[chainSelector],
				Data:         tx.Data(),
				Value:        big.NewInt(0),
				ContractType: "LinkToken",
			}
			batch.Batch = append(batch.Batch, op)

		}
		allBatches = append(allBatches, batch)
	}

	proposal, err := timelock.NewMCMSWithTimelockProposal(
		"1",
		req.ValidUntil,
		[]mcms.Signature{},
		req.OverrideRoot,
		chainMetadata,
		timelockAddresses,
		"Value transfer proposal",
		allBatches,
		timelock.Schedule,
		req.MinDelay.String(),
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*proposal},
	}, nil

}
