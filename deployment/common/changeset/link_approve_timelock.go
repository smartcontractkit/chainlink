package changeset

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

type LinkAllowances struct {
	Spender   common.Address
	Allowance big.Int
}
type LinkApproveTimelockRequest struct {
	Allowances       []LinkAllowances
	ChainSelector    uint64
	LinkTokenAddress common.Address
	TimelockAddress  common.Address
	MCMSAddress      common.Address
	ValidUntil       uint32        // unix time until the proposal will be valid
	MinDelay         time.Duration // delay for timelock worker to execute the transfers.
	OverrideRoot     bool
	StartingOpCount  uint64
}

var _ deployment.ChangeSet[*LinkApproveTimelockRequest] = LinkApproveTimelock

// packApprove packs the transferFrom method call data
func packApprove(parsedABI abi.ABI, spender common.Address, amount *big.Int) ([]byte, error) {
	data, err := parsedABI.Pack("approve", spender, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to pack transferFrom: %w", err)
	}
	return data, nil
}

// LinkApproveTimelock takes the given approvals for token transfers and creates an mcms proposal for them
func LinkApproveTimelock(_ deployment.Environment, req *LinkApproveTimelockRequest) (deployment.ChangesetOutput, error) {
	chainID := mcms.ChainIdentifier(req.ChainSelector)
	linkContract, err := link_token.NewLinkToken(req.LinkTokenAddress, nil)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get link contract: %w", err)
	}

	chainMetadata := map[mcms.ChainIdentifier]mcms.ChainMetadata{
		chainID: {MCMAddress: req.MCMSAddress, StartingOpCount: req.StartingOpCount},
	}
	timelockAddresses := map[mcms.ChainIdentifier]common.Address{
		chainID: req.TimelockAddress,
	}
	batch := timelock.BatchChainOperation{
		ChainIdentifier: chainID,
		Batch:           []mcms.Operation{},
	}

	for _, transfer := range req.Allowances {
		tx, err := linkContract.Approve(deployment.SimTransactOpts(), transfer.Spender, &transfer.Allowance)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error packing approve tx data: %w", err)
		}
		op := mcms.Operation{
			To:           req.LinkTokenAddress,
			Data:         tx.Data(),
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
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*proposal},
	}, nil

}
