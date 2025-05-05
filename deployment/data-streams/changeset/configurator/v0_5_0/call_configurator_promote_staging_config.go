package v0_5_0

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

var PromoteStagingConfigChangeset = cldf.CreateChangeSet(promoteStagingConfigLogic, promoteStagingConfigPrecondition)

type PromoteStagingConfigConfig struct {
	PromotionsByChain map[uint64][]PromoteStagingConfig
	MCMSConfig        *types.MCMSConfig
}

type PromoteStagingConfig struct {
	ConfiguratorAddress common.Address
	// 32-byte configId
	ConfigID [32]byte
	// Whether the current production config is considered green
	IsGreenProduction bool
}

type State struct {
	Configurator *configurator.Configurator
}

func (cfg PromoteStagingConfigConfig) Validate() error {
	if len(cfg.PromotionsByChain) == 0 {
		return errors.New("PromotionsByChain cannot be empty")
	}
	return nil
}

func promoteStagingConfigPrecondition(_ deployment.Environment, cc PromoteStagingConfigConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployConfiguratorConfig: %w", err)
	}
	return nil
}

func promoteStagingConfigLogic(e deployment.Environment, cfg PromoteStagingConfigConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	allBatches := []mcmstypes.BatchOperation{}

	var timelockAddressesPerChain map[uint64]string
	var proposerAddressPerChain map[uint64]string
	var inspectorPerChain map[uint64]mcmssdk.Inspector
	if cfg.MCMSConfig != nil {
		timelockAddressesPerChain = make(map[uint64]string)
		proposerAddressPerChain = make(map[uint64]string)
		inspectorPerChain = make(map[uint64]mcmssdk.Inspector)
	}

	for chainSelector, promotions := range cfg.PromotionsByChain {
		chain, chainExists := e.Chains[chainSelector]
		if !chainExists {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", chainSelector)
		}

		batch := mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chainSelector),
			Transactions:  []mcmstypes.Transaction{},
		}

		if cfg.MCMSConfig != nil {
			chainState := state.Chains[chainSelector]
			timelockAddressesPerChain[chainSelector] = chainState.Timelock.Address().String()
			proposerAddressPerChain[chainSelector] = chainState.ProposerMcm.Address().String()
			inspectorPerChain[chainSelector] = evm.NewInspector(chain.Client)
		}

		opts := getTransactOptsPromoteStagingConfig(e, chainSelector, cfg.MCMSConfig)
		for _, promotion := range promotions {
			confState, err := maybeLoadConfigurator(e, chainSelector, promotion.ConfiguratorAddress.String())
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}

			tx, err := promoteOrBuildTx(e, confState.Configurator, promotion, opts, chain, cfg.MCMSConfig)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}

			op := evm.NewTransaction(
				promotion.ConfiguratorAddress,
				tx.Data(),
				big.NewInt(0),
				string(types.Configurator),
				[]string{},
			)
			batch.Transactions = append(batch.Transactions, op)
		}

		allBatches = append(allBatches, batch)
	}

	if cfg.MCMSConfig != nil {
		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			e,
			timelockAddressesPerChain,
			proposerAddressPerChain,
			inspectorPerChain,
			allBatches,
			"PromoteStagingConfig proposal",
			proposalutils.TimelockConfig{
				MinDelay:     cfg.MCMSConfig.MinDelay,
				OverrideRoot: cfg.MCMSConfig.OverrideRoot,
			},
		)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}

		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal},
		}, nil
	}

	return deployment.ChangesetOutput{}, nil
}

func promoteOrBuildTx(
	e deployment.Environment,
	configuratorContract *configurator.Configurator,
	promotion PromoteStagingConfig,
	opts *bind.TransactOpts,
	chain deployment.Chain,
	mcmsConfig *types.MCMSConfig,
) (*ethTypes.Transaction, error) {
	tx, err := configuratorContract.PromoteStagingConfig(opts, promotion.ConfigID, promotion.IsGreenProduction)
	if err != nil {
		return nil, fmt.Errorf("error packing promoteStagingConfig tx data: %w", err)
	}

	if mcmsConfig == nil {
		if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
			e.Logger.Errorw("Failed to confirm promoteStagingConfig tx", "chain", chain.String(), "err", err)
			return nil, err
		}
	}
	return tx, nil
}

func maybeLoadConfigurator(e deployment.Environment, chainSel uint64, contractAddr string) (*State, error) {
	chain, ok := e.Chains[chainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found", chainSel)
	}

	addresses, err := e.ExistingAddresses.AddressesForChain(chainSel)
	if err != nil {
		return nil, err
	}

	tv, found := addresses[contractAddr]
	if !found {
		return nil, fmt.Errorf("unable to find Configurator contract on chain %s (selector %d)", chain.Name(), chain.Selector)
	}
	if tv.Type != types.Configurator || tv.Version != deployment.Version0_5_0 {
		return nil, fmt.Errorf("unexpected contract type %s for Configurator on chain %s (selector %d)", tv, chain.Name(), chain.Selector)
	}

	conf, err := configurator.NewConfigurator(common.HexToAddress(contractAddr), chain.Client)
	if err != nil {
		return nil, err
	}

	return &State{Configurator: conf}, nil
}

func getTransactOptsPromoteStagingConfig(e deployment.Environment, chainSel uint64, mcmsConfig *types.MCMSConfig) *bind.TransactOpts {
	if mcmsConfig == nil {
		return e.Chains[chainSel].DeployerKey
	}

	return deployment.SimTransactOpts()
}
