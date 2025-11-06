package v1_5_1

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	ccipcommoncs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

// UpdateAdminRoleChangesetV2 is a changeset that combines TransferAdminRoleChangesetV2 and ProposeAdminRoleChangesetV2 into
// one operation. It accepts the same inputs as TransferAdminRoleChangesetV2 + ProposeAdminRoleChangesetV2 then infers which
// one should be run by analyzing each input token's config from the TokenAdminRegistry. If the administrator is NOT defined
// then it'll try to propose one. Otherwise, it will assume that the existing admin should be transferred. Once the input is
// divided into one of two groups (either propose or transfer depending on whether an admin already exists on the token), it
// will reuse the pre-existing validation checks + logic for TransferAdminRoleChangesetV2 and ProposeAdminRoleChangesetV2 so
// that code duplication is kept at a minimum. This changeset also accepts an optional MCMS property that is reused for both
// underlying changesets. If it is defined, then OrchestrateChangesets will be used. Otherwise, things are run individually.
var UpdateAdminRoleChangesetV2 = cldf.CreateChangeSet(updateAdminRoleLogic, updateAdminRolePrecondition)

type UpdateAdminRoleConfig struct {
	// Only applicable when **proposing** a new admin - this allows the existing pending administrator to be
	// overridden if set to true. Use with caution as this will replace any existing pending admin proposals
	OverridePendingAdmin bool `json:"overridePendingAdmin"`

	// A map of chain selector => slice of TokenAdminInfo which describes the updates to make on each chain
	ChainUpdates map[uint64][]TokenAdminInfo `json:"ChainUpdates"`

	// The timelock config - all updates can be folded into one MCMS proposal with this setting
	MCMS *proposalutils.TimelockConfig `json:"mcms"`

	// Internal property for caching purposes
	configs *updateAdminRoleConfigs
}

type updateAdminRoleConfigs struct {
	orchestrateChangesetsConfig *ccipcommoncs.OrchestrateChangesetsConfig
	transferAdminRoleConfig     *TransferAdminRoleConfig
	proposeAdminRoleConfig      *ProposeAdminRoleConfig
}

func (cfg *UpdateAdminRoleConfig) populate(e cldf.Environment) (updateAdminRoleConfigs, error) {
	if cfg.configs != nil {
		e.Logger.Info("using cached configs")
		return *cfg.configs, nil
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return updateAdminRoleConfigs{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	transferAdminRoleConfig := TransferAdminRoleConfig{
		TransferAdminByChain: map[uint64][]TokenAdminInfo{},
		MCMS:                 cfg.MCMS,
	}

	proposeAdminRoleConfig := ProposeAdminRoleConfig{
		ProposeAdminByChain:  map[uint64][]TokenAdminInfo{},
		OverridePendingAdmin: cfg.OverridePendingAdmin,
		MCMS:                 cfg.MCMS,
	}

	for selector, updates := range cfg.ChainUpdates {
		chainState, ok := state.EVMChainState(selector)
		if !ok {
			return updateAdminRoleConfigs{}, fmt.Errorf("selector %d does not exist in state", selector)
		}

		tokenAddrSet := map[common.Address]bool{}
		transferInfo := []TokenAdminInfo{}
		proposeInfo := []TokenAdminInfo{}
		for _, info := range updates {
			// Ignore zero address (this will always have no token config)
			if info.TokenAddress == utils.ZeroAddress {
				e.Logger.Warnf("detected null token address for chain with selector '%d' - skipping", selector)
				continue
			}

			// Ignore duplicate token addresses
			exists := tokenAddrSet[info.TokenAddress]
			if exists {
				e.Logger.Warnf("detected duplicate token address (%s) for chain with selector '%d' - skipping", info.TokenAddress, selector)
				continue
			}

			e.Logger.Infof(
				"fetching token config for token '%s' from token admin registry at '%s' (chain selector = '%d')",
				info.TokenAddress.Hex(),
				chainState.TokenAdminRegistry.Address().Hex(),
				selector,
			)

			tokenConfig, err := chainState.TokenAdminRegistry.GetTokenConfig(&bind.CallOpts{Context: e.GetContext()}, info.TokenAddress)
			if err != nil {
				return updateAdminRoleConfigs{}, fmt.Errorf(
					"failed to get token config for token '%s' from chain with selector '%d'",
					info.TokenAddress.Hex(),
					selector,
				)
			}

			tokenAddrSet[info.TokenAddress] = true
			e.Logger.Infof(
				"found token config for token '%s' in token admin registry at '%s' (chain selector = '%d'): %+v",
				info.TokenAddress.Hex(),
				chainState.TokenAdminRegistry.Address().Hex(),
				selector,
				tokenConfig,
			)

			// If no admin exists for the token, then propose one otherwise transfer ownership
			if tokenConfig.Administrator == utils.ZeroAddress {
				proposeInfo = append(proposeInfo, info)
				continue
			}

			// Instead of throwing an error, ignore transfers to the same admin (this will make it
			// easier to re-run the changeset with the same inputs in case an error occurs midway)
			if tokenConfig.Administrator != info.AdminAddress {
				transferInfo = append(transferInfo, info)
				continue
			}
		}

		if len(transferInfo) > 0 {
			transferAdminRoleConfig.TransferAdminByChain[selector] = transferInfo
		}

		if len(proposeInfo) > 0 {
			proposeAdminRoleConfig.ProposeAdminByChain[selector] = proposeInfo
		}
	}

	changesets := []ccipcommoncs.WithConfig{}
	configs := updateAdminRoleConfigs{
		orchestrateChangesetsConfig: nil,
		transferAdminRoleConfig:     nil,
		proposeAdminRoleConfig:      nil,
	}

	if len(transferAdminRoleConfig.TransferAdminByChain) > 0 {
		configs.transferAdminRoleConfig = &transferAdminRoleConfig
		changesets = append(changesets, ccipcommoncs.CreateGenericChangeSetWithConfig(
			TransferAdminRoleChangesetV2,
			transferAdminRoleConfig,
		))
	}

	if len(proposeAdminRoleConfig.ProposeAdminByChain) > 0 {
		configs.proposeAdminRoleConfig = &proposeAdminRoleConfig
		changesets = append(changesets, ccipcommoncs.CreateGenericChangeSetWithConfig(
			ProposeAdminRoleChangesetV2,
			proposeAdminRoleConfig,
		))
	}

	if cfg.MCMS != nil {
		configs.orchestrateChangesetsConfig = &ccipcommoncs.OrchestrateChangesetsConfig{
			Description: "Propose or transfer admin roles",
			MCMS:        cfg.MCMS,
			ChangeSets:  changesets,
		}
	}

	cfg.configs = &configs
	return configs, nil
}

func updateAdminRolePrecondition(e cldf.Environment, cfg UpdateAdminRoleConfig) error {
	if len(cfg.ChainUpdates) == 0 {
		e.Logger.Warn("no chain updates were provided - exiting precondition stage gracefully")
		return nil
	}

	e.Logger.Info("populating internal configs for TransferAdminRoleChangesetV2 and ProposeAdminRoleChangesetV2")
	configs, err := cfg.populate(e)
	if err != nil {
		return fmt.Errorf("failed to populate internal configs: %w", err)
	}

	if configs.orchestrateChangesetsConfig == nil && configs.transferAdminRoleConfig == nil && configs.proposeAdminRoleConfig == nil {
		e.Logger.Warn("no operations to perform - exiting precondition stage gracefully")
		return nil
	}

	if configs.orchestrateChangesetsConfig != nil {
		e.Logger.Info("detected MCMS config - using OrchestrateChangesets to verify preconditions")
		return ccipcommoncs.OrchestrateChangesets.VerifyPreconditions(e, *configs.orchestrateChangesetsConfig)
	}

	e.Logger.Info("no MCMS config detected - verifying preconditions individually")
	if configs.transferAdminRoleConfig != nil {
		e.Logger.Info("verifying preconditions TransferAdminRoleChangesetV2...")
		err := TransferAdminRoleChangesetV2.VerifyPreconditions(e, *configs.transferAdminRoleConfig)
		if err != nil {
			return err
		}
		e.Logger.Info("successfully verified preconditions for TransferAdminRoleChangesetV2")
	}
	if configs.proposeAdminRoleConfig != nil {
		e.Logger.Info("verifying preconditions ProposeAdminRoleChangesetV2...")
		err := ProposeAdminRoleChangesetV2.VerifyPreconditions(e, *configs.proposeAdminRoleConfig)
		if err != nil {
			return err
		}
		e.Logger.Info("successfully verified preconditions for ProposeAdminRoleChangesetV2")
	}

	return nil
}

func updateAdminRoleLogic(e cldf.Environment, cfg UpdateAdminRoleConfig) (cldf.ChangesetOutput, error) {
	result := cldf.ChangesetOutput{}
	if len(cfg.ChainUpdates) == 0 {
		e.Logger.Warn("no chain updates were provided - exiting apply stage gracefully")
		return cldf.ChangesetOutput{}, nil
	}

	e.Logger.Info("populating internal configs for TransferAdminRoleChangesetV2 and ProposeAdminRoleChangesetV2")
	configs, err := cfg.populate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate internal configs: %w", err)
	}

	if configs.orchestrateChangesetsConfig == nil && configs.transferAdminRoleConfig == nil && configs.proposeAdminRoleConfig == nil {
		e.Logger.Warn("no operations to perform - exiting apply stage gracefully")
		return cldf.ChangesetOutput{}, nil
	}

	if configs.orchestrateChangesetsConfig != nil {
		e.Logger.Info("detected MCMS config - using OrchestrateChangesets to batch all operations into one MCMS proposal")
		return ccipcommoncs.OrchestrateChangesets.Apply(e, *configs.orchestrateChangesetsConfig)
	}

	e.Logger.Info("no MCMS config detected - applying changesets individually")
	if configs.transferAdminRoleConfig != nil {
		e.Logger.Info("applying TransferAdminRoleChangesetV2...")
		transferOutput, err := TransferAdminRoleChangesetV2.Apply(e, *configs.transferAdminRoleConfig)
		if err != nil {
			result.Reports = append(result.Reports, transferOutput.Reports...)
			return cldf.ChangesetOutput{Reports: result.Reports}, fmt.Errorf("failed to apply TransferAdminRoleChangesetV2: %w", err)
		}
		err = ccipcommoncs.MergeChangesetOutput(e, &result, transferOutput)
		if err != nil {
			result.Reports = append(result.Reports, transferOutput.Reports...)
			return cldf.ChangesetOutput{Reports: result.Reports}, fmt.Errorf("failed to merge output of TransferAdminRoleChangesetV2: %w", err)
		}
		e.Logger.Info("successfully applied TransferAdminRoleChangesetV2")
	}
	if configs.proposeAdminRoleConfig != nil {
		e.Logger.Info("applying ProposeAdminRoleChangesetV2...")
		proposeOutput, err := ProposeAdminRoleChangesetV2.Apply(e, *configs.proposeAdminRoleConfig)
		if err != nil {
			result.Reports = append(result.Reports, proposeOutput.Reports...)
			return cldf.ChangesetOutput{Reports: result.Reports}, fmt.Errorf("failed to apply ProposeAdminRoleChangesetV2: %w", err)
		}
		err = ccipcommoncs.MergeChangesetOutput(e, &result, proposeOutput)
		if err != nil {
			result.Reports = append(result.Reports, proposeOutput.Reports...)
			return cldf.ChangesetOutput{Reports: result.Reports}, fmt.Errorf("failed to merge output of ProposeAdminRoleChangesetV2: %w", err)
		}
		e.Logger.Info("successfully applied ProposeAdminRoleChangesetV2")
	}

	return result, nil
}
