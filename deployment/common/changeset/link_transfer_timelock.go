package changeset

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

type LinkTransfer struct {
	To    common.Address
	Value *big.Int
}
type LinkTransferTimelockConfig struct {
	Transfers  map[uint64][]LinkTransfer
	UseMCMS    bool
	McmsConfig *MCMSConfig
}
type MCMSConfig struct {
	ValidUntil      uint32        // unix time until the proposal will be valid
	MinDelay        time.Duration // delay for timelock worker to execute the transfers.
	OverrideRoot    bool
	StartingOpCount map[uint64]uint64
}

var _ deployment.ChangeSet[*LinkTransferTimelockConfig] = LinkTransferTimelock

// Validate checks that the LinkTransferTimelockConfig is valid.
func (cfg LinkTransferTimelockConfig) Validate() error {
	// Check that Transfers map has at least one key
	if len(cfg.Transfers) == 0 {
		return fmt.Errorf("transfers map must have at least one key")
	}

	// Check that each key in Transfers has at least one LinkTransfer
	for key, transfers := range cfg.Transfers {
		if len(transfers) == 0 {
			return fmt.Errorf("transfers for key %d must have at least one LinkTransfer", key)
		}
	}
	if !cfg.UseMCMS {
		return nil
	}
	// Mcms specific configs
	if cfg.McmsConfig == nil {
		return fmt.Errorf("mcmsConfig must be set when UseMCMS is true")
	}
	// Check that StartingOpCount map has at least one key
	if len(cfg.McmsConfig.StartingOpCount) == 0 {
		return fmt.Errorf("startingOpCount map must have at least one key")
	}
	// Check that Transfers and StartingOpCount have the same keys
	for key := range cfg.Transfers {
		if _, exists := cfg.McmsConfig.StartingOpCount[key]; !exists {
			return fmt.Errorf("startingOpCount map is missing key %d from transfers map", key)
		}
	}
	for key := range cfg.McmsConfig.StartingOpCount {
		if _, exists := cfg.Transfers[key]; !exists {
			return fmt.Errorf("transfers map is missing key %d from startingOpCount map", key)
		}
	}

	return nil
}

// LinkTransferTimelock takes the given link transfers and executes them or creates an MCMS proposal for them.
func LinkTransferTimelock(e deployment.Environment, req *LinkTransferTimelockConfig) (deployment.ChangesetOutput, error) {
	err := req.Validate()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid LinkTransferTimelockConfig: %w", err)
	}
	chainMetadata := map[mcms.ChainIdentifier]mcms.ChainMetadata{}
	timelockAddresses := map[mcms.ChainIdentifier]common.Address{}
	allBatches := []timelock.BatchChainOperation{}
	for chainSelector := range req.Transfers {
		chainID := mcms.ChainIdentifier(chainSelector)
		chain := e.Chains[chainSelector]
		addrs, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		linkState, err := MaybeLoadLinkTokenState(chain, addrs)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		linkAddress := linkState.LinkToken.Address()
		mcmsState, err := MaybeLoadMCMSWithTimelockState(chain, addrs)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		mcmAddress := mcmsState.ProposerMcm.Address()
		timelockAddress := mcmsState.Timelock.Address()

		chainMetadata[chainID] = mcms.ChainMetadata{
			MCMAddress:      mcmAddress,
			StartingOpCount: req.McmsConfig.StartingOpCount[chainSelector],
		}
		timelockAddresses[chainID] = timelockAddress
		batch := timelock.BatchChainOperation{
			ChainIdentifier: chainID,
			Batch:           []mcms.Operation{},
		}
		opts := deployment.SimTransactOpts()
		if !req.UseMCMS {
			opts = chain.DeployerKey
		}
		for _, transfer := range req.Transfers[chainSelector] {
			tx, err := linkState.LinkToken.Transfer(opts, transfer.To, transfer.Value)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("error packing transfer tx data: %w", err)
			}
			op := mcms.Operation{
				To:           linkAddress,
				Data:         tx.Data(),
				Value:        big.NewInt(0),
				ContractType: string(types.LinkToken),
			}
			batch.Batch = append(batch.Batch, op)

		}
		allBatches = append(allBatches, batch)
	}
	if req.UseMCMS {
		proposal, err := timelock.NewMCMSWithTimelockProposal(
			"1",
			req.McmsConfig.ValidUntil,
			[]mcms.Signature{},
			req.McmsConfig.OverrideRoot,
			chainMetadata,
			timelockAddresses,
			"Value transfer proposal",
			allBatches,
			timelock.Schedule,
			req.McmsConfig.MinDelay.String(),
		)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}

		return deployment.ChangesetOutput{
			Proposals: []timelock.MCMSWithTimelockProposal{*proposal},
		}, nil
	}

	return deployment.ChangesetOutput{}, nil
}
