package changeset

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

type (
	SetChannelDefinitionsConfig struct {
		// DefinitionsByChain is a map of chain selectors -> ChannelConfigStore addresses -> ChannelDefinitions to deploy.
		DefinitionsByChain map[uint64]map[string]ChannelDefinition // Use string for address because of future non-EVM chains.
		MCMSConfig         *MCMSConfig
	}

	ChannelDefinition struct {
		ChannelConfigStore common.Address

		DonID uint32
		S3URL string
		Hash  [32]byte
	}

	MCMSConfig struct {
		MinDelay     time.Duration
		OverrideRoot bool
	}

	ChannelConfigStoreState struct {
		ChannelConfigStore *channel_config_store.ChannelConfigStore
	}
)

func (cfg SetChannelDefinitionsConfig) Validate() error {
	if len(cfg.DefinitionsByChain) == 0 {
		return errors.New("DefinitionsByChain cannot be empty")
	}
	return nil
}

func CallSetChannelDefinitions(e deployment.Environment, cfg SetChannelDefinitionsConfig) (deployment.ChangesetOutput, error) {
	err := cfg.Validate()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid SetChannelDefinitionsConfig: %w", err)
	}

	chainSelectors := []uint64{}
	for chainSelector := range cfg.DefinitionsByChain {
		chainSelectors = append(chainSelectors, chainSelector)
	}

	// Initialize state for each chain
	var mcmsStatePerChain map[uint64]*changeset.MCMSWithTimelockState
	if cfg.MCMSConfig != nil {
		mcmsStatePerChain, err = changeset.MaybeLoadMCMSWithTimelockState(e, chainSelectors)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
	}

	mcmsPerChain := map[uint64]*owner_helpers.ManyChainMultiSig{}
	timelockAddresses := map[uint64]common.Address{}
	allBatches := []timelock.BatchChainOperation{}
	for chainSelector := range cfg.DefinitionsByChain {
		chainID := mcms.ChainIdentifier(chainSelector)

		batch := timelock.BatchChainOperation{
			ChainIdentifier: chainID,
			Batch:           []mcms.Operation{},
		}

		chain, ok := e.Chains[chainSelector]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain %d not found", chainSelector)
		}

		if cfg.MCMSConfig != nil {
			mcmsState := mcmsStatePerChain[chainSelector]
			timelockAddress := mcmsState.Timelock.Address()
			mcmsPerChain[uint64(chainID)] = mcmsState.ProposerMcm
			timelockAddresses[chainSelector] = timelockAddress
		}

		opts := getTransactOpts(e, chainSelector, cfg.MCMSConfig)

		for contractAddr, definition := range cfg.DefinitionsByChain[chainSelector] {
			ccsState, err := maybeLoadChannelConfigStoreState(e, chainSelector, contractAddr)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}

			tx, err := transferOrBuildTx(e, ccsState.ChannelConfigStore, definition, opts, chain, cfg.MCMSConfig)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			op := mcms.Operation{
				To:           common.HexToAddress(contractAddr),
				Data:         tx.Data(),
				Value:        big.NewInt(0),
				ContractType: string(types.ChannelConfigStore),
			}
			batch.Batch = append(batch.Batch, op)
		}

		allBatches = append(allBatches, batch)
	}

	if cfg.MCMSConfig != nil {
		proposal, err := proposalutils.BuildProposalFromBatches(
			timelockAddresses,
			mcmsPerChain,
			allBatches,
			"ChannelConfigStore setChannelDefinitions proposal",
			cfg.MCMSConfig.MinDelay,
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

func getTransactOpts(e deployment.Environment, chainSel uint64, mcmsConfig *MCMSConfig) *bind.TransactOpts {
	if mcmsConfig == nil {
		return e.Chains[chainSel].DeployerKey
	}

	return deployment.SimTransactOpts()
}

func transferOrBuildTx(
	e deployment.Environment,
	ccs *channel_config_store.ChannelConfigStore,
	definition ChannelDefinition,
	opts *bind.TransactOpts,
	chain deployment.Chain,
	mcmsConfig *MCMSConfig,
) (*ethTypes.Transaction, error) {
	if ccs == nil {
		return nil, errors.New("provided ChannelConfigStore is nil")
	}
	tx, err := ccs.ChannelConfigStoreTransactor.SetChannelDefinitions(opts, definition.DonID, definition.S3URL, definition.Hash)
	if err != nil {
		return nil, fmt.Errorf("error packing transfer tx data: %w", err)
	}
	// only wait for tx if we are not using MCMS
	if mcmsConfig == nil {
		if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
			e.Logger.Errorw("Failed to confirm transfer tx", "chain", chain.String(), "err", err)
			return nil, err
		}
	}
	return tx, nil
}

func maybeLoadChannelConfigStoreState(e deployment.Environment, chainSel uint64, contractAddr string) (*ChannelConfigStoreState, error) {
	chain, ok := e.Chains[chainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found", chainSel)
	}
	addresses, err := e.ExistingAddresses.AddressesForChain(chainSel)
	if err != nil {
		return nil, err
	}
	tv, ok := addresses[contractAddr]
	if !ok {
		return nil, fmt.Errorf("unable to find channlConfigStore contract on chain %s (chain selector %d)", chain.Name(), chain.Selector)
	}

	if tv.Type != types.ChannelConfigStore || tv.Version != deployment.Version1_0_0 {
		return nil, fmt.Errorf("unexpected contract type %s for channlConfigStore on chain %s (chain selector %d)", tv, chain.Name(), chain.Selector)
	}
	ccs, err := channel_config_store.NewChannelConfigStore(common.HexToAddress(contractAddr), chain.Client)
	if err != nil {
		return nil, err
	}

	return &ChannelConfigStoreState{
		ChannelConfigStore: ccs,
	}, nil
}
