package v0_5_0

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"

	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"

	"github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

var SetStagingConfigChangeset = cldf.CreateChangeSet(setStagingConfigLogic, setStagingConfigPrecondition)

type Config interface {
	GetConfiguratorAddress() common.Address
}

type SetStagingConfigConfig struct {
	ConfigurationsByChain map[uint64][]SetStagingConfig
	MCMSConfig            *changeset.MCMSConfig
}

type SetStagingConfig struct {
	ConfiguratorAddress   common.Address
	ConfigID              [32]byte
	Signers               [][]byte
	OffchainTransmitters  [][32]byte
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

func (sc SetStagingConfig) GetConfiguratorAddress() common.Address {
	return sc.ConfiguratorAddress
}

func setStagingConfigPrecondition(_ deployment.Environment, ss SetStagingConfigConfig) error {
	if err := ss.Validate(); err != nil {
		return fmt.Errorf("invalid SetStagingConfigConfig: %w", err)
	}

	return nil
}

func (cfg SetStagingConfigConfig) Validate() error {
	if len(cfg.ConfigurationsByChain) == 0 {
		return errors.New("ConfigurationsByChain cannot be empty")
	}
	return nil
}

func callSetConfigCommon[T Config](
	e deployment.Environment,
	configurationsByChain map[uint64][]T,
	mcmsConfig *changeset.MCMSConfig,
	proposalName string,
	txBuilder func(e deployment.Environment, configuratorContract *configurator.Configurator, cfg T, opts *bind.TransactOpts, chain deployment.Chain, mcmsConfig *changeset.MCMSConfig) (*ethTypes.Transaction, error),
) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	timelockAddressesPerChain := make(map[uint64]string)
	proposerAddressPerChain := make(map[uint64]string)
	inspectorPerChain := make(map[uint64]mcmssdk.Inspector)

	allBatches := []mcmstypes.BatchOperation{}

	for chainSelector, configs := range configurationsByChain {
		chain, ok := e.Chains[chainSelector]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", chainSelector)
		}
		batch := mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chainSelector),
			Transactions:  []mcmstypes.Transaction{},
		}

		if mcmsConfig != nil {
			chainState := state.Chains[chainSelector]

			timelockAddressesPerChain[chainSelector] = chainState.Timelock.Address().String()
			proposerAddressPerChain[chainSelector] = chainState.ProposerMcm.Address().String()
			inspectorPerChain[chainSelector] = evm.NewInspector(chain.Client)
		}

		opts := getTransactOpts(e, chainSelector, mcmsConfig)
		for _, cfg := range configs {
			state, err := maybeLoadConfigurator(e, chainSelector, cfg.GetConfiguratorAddress().String())
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}

			tx, err := txBuilder(e, state.Configurator, cfg, opts, chain, mcmsConfig)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}

			op := evm.NewTransaction(cfg.GetConfiguratorAddress(), tx.Data(), big.NewInt(0), string(types.Configurator), []string{})

			batch.Transactions = append(batch.Transactions, op)
		}

		allBatches = append(allBatches, batch)
	}

	if mcmsConfig != nil {
		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			e,
			timelockAddressesPerChain,
			proposerAddressPerChain,
			inspectorPerChain,
			allBatches,
			proposalName,
			proposalutils.TimelockConfig{
				MinDelay:     mcmsConfig.MinDelay,
				OverrideRoot: mcmsConfig.OverrideRoot,
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

func setStagingConfigLogic(e deployment.Environment, cfg SetStagingConfigConfig) (deployment.ChangesetOutput, error) {
	return callSetConfigCommon(e, cfg.ConfigurationsByChain, cfg.MCMSConfig, "SetStagingConfig proposal", setStagingConfigTx)
}

func setStagingConfigTx(
	e deployment.Environment,
	configuratorContract *configurator.Configurator,
	stagingCfg SetStagingConfig,
	opts *bind.TransactOpts,
	chain deployment.Chain,
	mcmsConfig *changeset.MCMSConfig,
) (*ethTypes.Transaction, error) {
	tx, err := configuratorContract.SetStagingConfig(opts, stagingCfg.ConfigID, stagingCfg.Signers, stagingCfg.OffchainTransmitters, stagingCfg.F, stagingCfg.OnchainConfig, stagingCfg.OffchainConfigVersion, stagingCfg.OffchainConfig)
	if err != nil {
		return nil, fmt.Errorf("error packing setStagingConfig tx data: %w", err)
	}

	if mcmsConfig == nil {
		if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
			e.Logger.Errorw("Failed to confirm setStagingConfig tx", "chain", chain.String(), "err", err)
			return nil, err
		}
	}
	return tx, nil
}
