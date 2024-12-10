package changeset

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
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
	ValidUntil   uint32        // unix time until the proposal will be valid
	MinDelay     time.Duration // delay for timelock worker to execute the transfers.
	OverrideRoot bool
}

var _ deployment.ChangeSet[*LinkTransferTimelockConfig] = LinkTransferTimelock

// Validate checks that the LinkTransferTimelockConfig is valid.
func (cfg LinkTransferTimelockConfig) Validate() error {

	// Check that Transfers map has at least one key
	if len(cfg.Transfers) == 0 {
		return fmt.Errorf("transfers map must have at least one key")
	}

	// Check transfers config values.
	for key, transfers := range cfg.Transfers {
		if len(transfers) == 0 {
			return fmt.Errorf("transfers for key %d must have at least one LinkTransfer", key)
		}
		for _, transfer := range transfers {
			if transfer.To == (common.Address{}) {
				return fmt.Errorf("'to' address for transfers  must be set")
			}
			if transfer.Value == nil {
				return fmt.Errorf("value for transfers must be set")
			}
			if transfer.Value == big.NewInt(0) {
				return fmt.Errorf("value for transfers must be non-zero")
			}
		}
	}
	// Check that Transfers and StartingOpCount have the same keys
	for key := range cfg.Transfers {
		_, err := chain_selectors.GetSelectorFamily(key)
		if err != nil {
			return fmt.Errorf("invalid chain selector: %w", err)
		}
	}
	if !cfg.UseMCMS {
		return nil
	}
	// Mcms specific configs
	if cfg.McmsConfig == nil {
		return fmt.Errorf("mcmsConfig must be set when UseMCMS is true")
	}
	// Upper bound for min delay (7 days)
	if cfg.McmsConfig.MinDelay > 24*7*time.Hour {
		return fmt.Errorf("minDelay must be less than 7 days")
	}
	return nil
}

// LinkTransferTimelock takes the given link transfers and executes them or creates an MCMS proposal for them.
func LinkTransferTimelock(e deployment.Environment, req *LinkTransferTimelockConfig) (deployment.ChangesetOutput, error) {
	err := req.Validate()
	ctx := e.GetContext()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid LinkTransferTimelockConfig: %w", err)
	}
	chainSelectors := []uint64{}
	for chainSelector := range req.Transfers {
		chainSelectors = append(chainSelectors, chainSelector)
	}
	mcmsPerChain := map[uint64]*owner_helpers.ManyChainMultiSig{}

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
		timelockAddress := mcmsState.Timelock.Address()

		mcmsPerChain[uint64(chainID)] = mcmsState.ProposerMcm

		timelockAddresses[chainID] = timelockAddress
		batch := timelock.BatchChainOperation{
			ChainIdentifier: chainID,
			Batch:           []mcms.Operation{},
		}
		opts := deployment.SimTransactOpts()
		if !req.UseMCMS {
			opts = chain.DeployerKey
		}
		totalAmount := big.NewInt(0)
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
			totalAmount.Add(totalAmount, transfer.Value)
		}
		// check that from address has enough funds for the transfers
		balance, err := linkState.LinkToken.BalanceOf(&bind.CallOpts{Context: ctx}, timelockAddress)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		if balance.Cmp(totalAmount) < 0 {
			return deployment.ChangesetOutput{}, fmt.Errorf("timelock address does not have enough funds for transfers for chainID %d", chainSelector)
		}
		allBatches = append(allBatches, batch)
	}
	chainMetadata, err := proposalutils.BuildProposalMetadata(chainSelectors, mcmsPerChain)
	if err != nil {
		return deployment.ChangesetOutput{}, err
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
