package changeset

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
	Transfers        []LinkTransfer
	ChainSelector    uint64
	LinkTokenAddress common.Address
	TimelockAddress  common.Address
	MCMSAddress      common.Address
	ValidUntil       uint32        // unix time until the proposal will be valid
	MinDelay         time.Duration // delay for timelock worker to execute the transfers.
	OverrideRoot     bool
	StartingOpCount  uint64
}

var _ deployment.ChangeSet[*LinkTransferTimelockRequest] = LinkTransferTimelock

// packTransferFrom packs the transferFrom method call data
func packTransferFrom(parsedABI abi.ABI, to common.Address, value *big.Int) ([]byte, error) {
	data, err := parsedABI.Pack("transfer", to, value)
	if err != nil {
		return nil, fmt.Errorf("failed to pack transferFrom: %w", err)
	}
	return data, nil
}

// LinkTransferTimelock takes the given link transfers and creates an MCMS proposal for them.
func LinkTransferTimelock(_ deployment.Environment, req *LinkTransferTimelockRequest) (deployment.ChangesetOutput, error) {
	chainID := mcms.ChainIdentifier(req.ChainSelector)
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
	// Parse the ABI
	parsedABI, err := abi.JSON(strings.NewReader(link_token.LinkTokenMetaData.ABI))
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to parse ABI: %w", err)
	}
	for _, transfer := range req.Transfers {
		data, err := packTransferFrom(parsedABI, transfer.To, &transfer.Value)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error packing transferFrom data: %w", err)
		}
		op := mcms.Operation{
			To:           req.LinkTokenAddress,
			Data:         data,
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
