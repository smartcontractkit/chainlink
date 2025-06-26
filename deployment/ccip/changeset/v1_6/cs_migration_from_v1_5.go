package v1_6

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	migrate_seq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/migration"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

var (
	// InitChainUpgratesChangeset sets candidates for the commit and exec DONs for multiple destination chains.
	// It then identifies all existing 1.5.0 source chains for each chain in the batch.
	// For each 1.5.0 OnRamp connecting to a destination, configuration gets translated to the 1.6.0 FeeQuoter.
	// In addition, OnRamps are connected to destination chains via test routers.
	// We do NOT connect the destinations back to the source chains, as DONs are not guaranteed to exist for sources.
	// This changeset is NOT IDEMPOTENT - if AddDON is called more than once for the same chain it will revert.
	// This changeset also requires all FeeQuoters and NonceManagers to be owned by the MCMS timelock contract.
	InitChainUpgradesChangeset = cldf.CreateChangeSet(
		initChainUpgradesLogic,
		initChainUpgradesPrecondition,
	)
	// PromoteChainUpgradesChangeset promotes the commit and exec DON candidates for multiple destination chains.
	// It then connects the source chains to the destination chains via main routers.
	// Before running PromoteChainUpgradesChangeset for a batch, you must run InitChainUpgradesOnTestRoutersChangeset followed by SetOCR3OffRampChangeset.
	// SetOCR3OffRampChangeset should be run with ConfigType set to candidate, since the config won't be promoted until this changeset is run.
	// This changeset is NOT IDEMPOTENT - re-promoting will result in clearing the active digest, which is not desired.
	PromoteChainUpgradesChangeset = cldf.CreateChangeSet(
		promoteChainUpgradesLogic,
		promoteChainUpgradesPrecondition,
	)
)

// ChainUpgradeConfig defines the commit and exec OCR parameters for a given chain.
type ChainUpgradeConfig struct {
	// FeedChainSelector is the selector of the chain housing the feeds used by the commit plugin.
	FeedChainSelector uint64
	// CommitOCRParams defines the OCR parameters for the commit plugin.
	CommitOCRParams CCIPOCRParams
	// ExecOCRParams defines the OCR parameters for the exec plugin.
	ExecOCRParams CCIPOCRParams
	// NewFeeQuoterParamsPerSource defines the new FeeQuoter parameters for each source chain connected to this destination chain.
	NewFeeQuoterParamsPerSource map[uint64]migrate_seq.NewFeeQuoterDestChainConfigParams
}

// InitChainUpgradesConfig defines the configuration for the InitChainUpgradesChangeset.
type InitChainUpgradesConfig struct {
	// HomeChainSelector is the selector of the home chain.
	HomeChainSelector uint64
	// ChainsToUpgrade is a map of destination chain selectors to their upgrade configurations.
	ChainsToUpgrade map[uint64]ChainUpgradeConfig
	// MCMSConfig is the configuration for the MCMS.
	MCMSConfig *proposalutils.TimelockConfig
}

func initChainUpgradesPrecondition(e cldf.Environment, c InitChainUpgradesConfig) error {
	if c.MCMSConfig == nil {
		return errors.New("MCMSConfig must be defined")
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	err = ValidateHomeChainState(e, c.HomeChainSelector, state)
	if err != nil {
		return fmt.Errorf("failed to validate home chain state: %w", err)
	}
	// Home chain contracts are owned by MCMS.
	err = commoncs.ValidateOwnership(e.GetContext(), true, e.BlockChains.EVMChains()[c.HomeChainSelector].DeployerKey.From, state.Chains[c.HomeChainSelector].Timelock.Address(), state.Chains[c.HomeChainSelector].CCIPHome)
	if err != nil {
		return fmt.Errorf("failed to validate ownership of CCIPHome on %s: %w", e.BlockChains.EVMChains()[c.HomeChainSelector], err)
	}
	err = commoncs.ValidateOwnership(e.GetContext(), true, e.BlockChains.EVMChains()[c.HomeChainSelector].DeployerKey.From, state.Chains[c.HomeChainSelector].Timelock.Address(), state.Chains[c.HomeChainSelector].CapabilityRegistry)
	if err != nil {
		return fmt.Errorf("failed to validate ownership of CapabilityRegistry on %s: %w", e.BlockChains.EVMChains()[c.HomeChainSelector], err)
	}

	allChainSels := make(map[uint64]struct{})
	for chainSel, chainUpgradeCfg := range c.ChainsToUpgrade {
		// Chain selector is a valid EVM chain selector & all MCMS contracts exist.
		err := stateview.ValidateChain(e, state, chainSel, c.MCMSConfig)
		if err != nil {
			return fmt.Errorf("failed to validate chain %d: %w", chainSel, err)
		}

		allChainSels[chainSel] = struct{}{}
		sourceChainSels := getSourceChainsForSelector(state, chainSel)
		for _, sourceChainSel := range sourceChainSels {
			// Source chain selector is a valid EVM chain selector & all MCMS contracts exist.
			err := stateview.ValidateChain(e, state, sourceChainSel, c.MCMSConfig)
			if err != nil {
				return fmt.Errorf("failed to validate chain %d: %w", sourceChainSel, err)
			}
			allChainSels[sourceChainSel] = struct{}{}

			// Price registry exists on source if 1.5.0 OnRamps exist
			if len(state.Chains[sourceChainSel].EVM2EVMOnRamp) > 0 && state.Chains[sourceChainSel].PriceRegistry == nil {
				return fmt.Errorf("price registry does not exist on source chain %d, but 1.5.0 OnRamps exist", sourceChainSel)
			}
		}

		// Commit OCR params are valid.
		err = chainUpgradeCfg.CommitOCRParams.Validate(e, chainSel, chainUpgradeCfg.FeedChainSelector, state)
		if err != nil {
			return fmt.Errorf("failed to validate commit OCR params for chain %d: %w", chainSel, err)
		}

		// Exec OCR params are valid.
		err = chainUpgradeCfg.ExecOCRParams.Validate(e, chainSel, chainUpgradeCfg.FeedChainSelector, state)
		if err != nil {
			return fmt.Errorf("failed to validate exec OCR params for chain %d: %w", chainSel, err)
		}

		// ARMProxy contracts are owned by MCMS.
		err = commoncs.ValidateOwnership(e.GetContext(), true, e.BlockChains.EVMChains()[chainSel].DeployerKey.From, state.Chains[chainSel].Timelock.Address(), state.Chains[chainSel].RMNProxy)
		if err != nil {
			return fmt.Errorf("failed to validate ownership of RMNProxy on %s: %w", e.BlockChains.EVMChains()[chainSel], err)
		}
	}

	for chainSel := range allChainSels {
		// FeeQuoter and NonceManager are owned by MCMS on both source and destination chains.
		err := commoncs.ValidateOwnership(e.GetContext(), true, e.BlockChains.EVMChains()[chainSel].DeployerKey.From, state.Chains[chainSel].Timelock.Address(), state.Chains[chainSel].FeeQuoter)
		if err != nil {
			return fmt.Errorf("failed to validate ownership of FeeQuoter on %s: %w", e.BlockChains.EVMChains()[chainSel], err)
		}
		err = commoncs.ValidateOwnership(e.GetContext(), true, e.BlockChains.EVMChains()[chainSel].DeployerKey.From, state.Chains[chainSel].Timelock.Address(), state.Chains[chainSel].NonceManager)
		if err != nil {
			return fmt.Errorf("failed to validate ownership of NonceManager on %s: %w", e.BlockChains.EVMChains()[chainSel], err)
		}
	}

	return nil
}

func initChainUpgradesLogic(e cldf.Environment, c InitChainUpgradesConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	allProposals := make([]mcms.TimelockProposal, 0)
	allReports := make([]operations.Report[any, any], 0)

	// Collect all chain names for reporting purposes.
	allChainNames := make([]string, 0, len(c.ChainsToUpgrade))
	for destChainSel := range c.ChainsToUpgrade {
		allChainNames = append(allChainNames, e.BlockChains.EVMChains()[destChainSel].String())
	}

	// Some chains already have commit and exec DONs created on 1.6.0.
	// We need to filter these out before calling AddDonAndSetCandidateChangeset.
	existingDONs := make(map[uint64]uint32)
	for destChainSel := range c.ChainsToUpgrade {
		donID, err := internal.DonIDForChain(
			state.Chains[c.HomeChainSelector].CapabilityRegistry,
			state.Chains[c.HomeChainSelector].CCIPHome,
			destChainSel,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch DON ID for chain %d: %w", destChainSel, err)
		}
		if donID != 0 {
			existingDONs[destChainSel] = donID
		}
	}

	// Fetch the next DON ID from the home chain's capability registry.
	// This is so we can assign a DON ID to each chain in the batch.
	// TODO: Possibility of conflict with new chain integration workstream.
	nextDonID, err := state.Chains[c.HomeChainSelector].CapabilityRegistry.GetNextDONId(&bind.CallOpts{Context: e.GetContext()})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get next DON ID: %w", err)
	}

	for destChainSel, chainUpgradeCfg := range c.ChainsToUpgrade {
		chain := e.BlockChains.EVMChains()[destChainSel]
		donIDToUse := nextDonID
		setCandidateBase := SetCandidateConfigBase{
			HomeChainSelector: c.HomeChainSelector,
			FeedChainSelector: chainUpgradeCfg.FeedChainSelector,
			MCMS:              c.MCMSConfig,
		}
		commitCandidatePluginInfo := SetCandidatePluginInfo{
			PluginType: types.PluginTypeCCIPCommit,
			OCRConfigPerRemoteChainSelector: map[uint64]CCIPOCRParams{
				destChainSel: chainUpgradeCfg.CommitOCRParams,
			},
		}

		// (Add DON &) set candidate for commit plugin
		if existingID, ok := existingDONs[destChainSel]; ok {
			donIDToUse = existingID

			// DON already exists, so we just run SetCandidateChangeset for commit plugin.
			out, err := SetCandidateChangeset(e, SetCandidateChangesetConfig{
				SetCandidateConfigBase: setCandidateBase,
				PluginInfo:             []SetCandidatePluginInfo{commitCandidatePluginInfo},
				DonIDOverrides:         map[uint64]uint32{destChainSel: donIDToUse},
			})
			allReports = append(allReports, out.Reports...)
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run SetCandidateChangeset for commit plugin on %s: %w", chain, err)
			}
			allProposals = append(allProposals, out.MCMSTimelockProposals...)
		} else {
			out, err := AddDonAndSetCandidateChangeset(e, AddDonAndSetCandidateChangesetConfig{
				SetCandidateConfigBase: setCandidateBase,
				PluginInfo:             commitCandidatePluginInfo,
				DonIDOverride:          donIDToUse,
			})
			allReports = append(allReports, out.Reports...)
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run AddDonAndSetCandidateChangeset for commit plugin on %s: %w", chain, err)
			}
			allProposals = append(allProposals, out.MCMSTimelockProposals...)

			// Increment the DON ID since addDON was / will be called.
			nextDonID++
		}

		// Set candidate for exec plugin
		out, err := SetCandidateChangeset(e, SetCandidateChangesetConfig{
			SetCandidateConfigBase: setCandidateBase,
			PluginInfo: []SetCandidatePluginInfo{
				{
					PluginType: types.PluginTypeCCIPExec,
					OCRConfigPerRemoteChainSelector: map[uint64]CCIPOCRParams{
						destChainSel: chainUpgradeCfg.ExecOCRParams,
					},
				},
			},
			DonIDOverrides: map[uint64]uint32{destChainSel: donIDToUse},
		})
		allReports = append(allReports, out.Reports...)
		if err != nil {
			return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run SetCandidateChangeset for exec plugin on %s: %w", chain, err)
		}
		allProposals = append(allProposals, out.MCMSTimelockProposals...)
	}

	for destChainSel, chainUpgradeCfg := range c.ChainsToUpgrade {
		destChain := e.BlockChains.EVMChains()[destChainSel]

		// Ensure that RMNRemote is owned by the timelock contract
		out, err := ensureTimelockOwnership(e, destChainSel, []commoncs.Ownable{
			state.Chains[destChainSel].RMNRemote,
		}, *c.MCMSConfig)
		allReports = append(allReports, out.Reports...)
		if err != nil {
			return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to ensure timelock ownership of RMNRemote on %s: %w", destChain, err)
		}
		allProposals = append(allProposals, out.MCMSTimelockProposals...)

		// Point ARMProxy to RMNRemote
		out, err = SetRMNRemoteOnRMNProxyChangeset(e, SetRMNRemoteOnRMNProxyConfig{
			ChainSelectors: []uint64{destChainSel},
			MCMSConfig:     c.MCMSConfig,
		})
		allReports = append(allReports, out.Reports...)
		if err != nil {
			return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run SetRMNRemoteOnRMNProxyChangeset on %s: %w", destChain, err)
		}
		allProposals = append(allProposals, out.MCMSTimelockProposals...)

		// Transfer 1.5.0 OnRamp configs to FeeQuoter
		out, err = TranslateEVM2EVMOnRampsToFeeQuoterChangeset(e, TranslateEVM2EVMOnRampsToFeeQuoterConfig{
			NewFeeQuoterParamsPerSource: chainUpgradeCfg.NewFeeQuoterParamsPerSource,
			DestChainSelector:           destChainSel,
			MCMS:                        c.MCMSConfig,
		})
		allReports = append(allReports, out.Reports...)
		if err != nil {
			return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run TranslateEVM2EVMOnRampsToFeeQuoterChangeset for source chains of %s: %w", destChain, err)
		}
		allProposals = append(allProposals, out.MCMSTimelockProposals...)

		// Transfer token transfer fee configs to FeeQuoter
		out, err = TranslateEVM2EVMOnRampsToFeeQTokenTransferFeeConfigChangeset(e, TranslateEVM2EVMOnRampsToFeeQuoterConfig{
			NewFeeQuoterParamsPerSource: chainUpgradeCfg.NewFeeQuoterParamsPerSource,
			DestChainSelector:           destChainSel,
			MCMS:                        c.MCMSConfig,
		})
		allReports = append(allReports, out.Reports...)
		if err != nil {
			return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run TranslateEVM2EVMOnRampsToFeeQTokenTransferFeeConfigChangeset for source chains of %s: %w", destChain, err)
		}
		allProposals = append(allProposals, out.MCMSTimelockProposals...)

		// Loop through each source connected to the destination chain.
		sourceChainsToConnect := getSourceChainsForSelector(state, destChainSel)
		for _, sourceChainSel := range sourceChainsToConnect {
			sourceChain := e.BlockChains.EVMChains()[sourceChainSel]

			// Add 1.5.0 OnRamps and 1.5.0 OffRamps to NonceManager
			// This is done to ensure that nonces are correctly managed across versions.
			// OverrideExisting is set to true because previous a remote chain can show up twice for a given local chain (as a source or dest).
			// We don't want applyPreviousRampsUpdates to fail when it is called for the second time for a given remote chain selector.
			out, err = UpdateNonceManagersChangeset(e, UpdateNonceManagerConfig{
				MCMS: c.MCMSConfig,
				UpdatesByChain: map[uint64]NonceManagerUpdate{
					destChainSel: {
						PreviousRampsArgs: []PreviousRampCfg{
							{
								RemoteChainSelector: sourceChainSel,
								OverrideExisting:    true,
							},
						},
					},
					sourceChainSel: {
						PreviousRampsArgs: []PreviousRampCfg{
							{
								RemoteChainSelector: destChainSel,
								OverrideExisting:    true,
							},
						},
					},
				},
			})
			allReports = append(allReports, out.Reports...)
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run UpdateNonceManagersChangeset on %s and %s: %w", sourceChain, destChain, err)
			}
			allProposals = append(allProposals, out.MCMSTimelockProposals...)

			// Update OnRamp 1.6.0 on source chain (use test router).
			// OnRamp may be owned by timelock or deployer key here, so we need to check.
			mcmsConfig := c.MCMSConfig
			owner, err := state.Chains[sourceChainSel].OnRamp.Owner(&bind.CallOpts{Context: e.GetContext()})
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to get OnRamp owner on %s: %w", sourceChain, err)
			}
			if owner == e.BlockChains.EVMChains()[sourceChainSel].DeployerKey.From {
				mcmsConfig = nil // If OnRamp is owned by deployer key, we don't use MCMS.
			}
			out, err = UpdateOnRampsDestsChangeset(e, UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]OnRampDestinationUpdate{
					sourceChainSel: {
						destChainSel: {
							TestRouter:       true,
							AllowListEnabled: false,
							IsEnabled:        true,
						},
					},
				},
				MCMS: mcmsConfig,
			})
			allReports = append(allReports, out.Reports...)
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run UpdateOnRampsDestsChangeset on %s: %w", sourceChain, err)
			}
			allProposals = append(allProposals, out.MCMSTimelockProposals...)

			// Update OffRamp 1.6.0 on destination chain (use test router, no RMN verification).
			// OffRamp may be owned by timelock or deployer key here, so we need to check.
			mcmsConfig = c.MCMSConfig
			owner, err = state.Chains[destChainSel].OffRamp.Owner(&bind.CallOpts{Context: e.GetContext()})
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to get OffRamp owner on %s: %w", destChain, err)
			}
			if owner == e.BlockChains.EVMChains()[destChainSel].DeployerKey.From {
				mcmsConfig = nil // If OffRamp is owned by deployer key, we don't use MCMS.
			}
			out, err = UpdateOffRampSourcesChangeset(e, UpdateOffRampSourcesConfig{
				UpdatesByChain: map[uint64]map[uint64]OffRampSourceUpdate{
					destChainSel: {
						sourceChainSel: {
							TestRouter:                true,
							IsRMNVerificationDisabled: true, // TODO: Might not always be true, but for now we assume it is.
							IsEnabled:                 true,
						},
					},
				},
				MCMS: mcmsConfig,
			})
			allReports = append(allReports, out.Reports...)
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run UpdateOffRampSourcesChangeset on %s: %w", destChain, err)
			}
			allProposals = append(allProposals, out.MCMSTimelockProposals...)

			// Set OnRamp on source router, OffRamp on dest router (use test router).
			// Test routers are never owned by MCMS.
			out, err = UpdateRouterRampsChangeset(e, UpdateRouterRampsConfig{
				TestRouter: true,
				UpdatesByChain: map[uint64]RouterUpdates{
					destChainSel: {
						OffRampUpdates: map[uint64]bool{sourceChainSel: true},
					},
					sourceChainSel: {
						OnRampUpdates: map[uint64]bool{destChainSel: true},
					},
				},
			})
			allReports = append(allReports, out.Reports...)
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run UpdateRouterRampsChangeset and connect %s with %s: %w", sourceChain, destChain, err)
			}
			allProposals = append(allProposals, out.MCMSTimelockProposals...)
		}
	}

	proposal, err := proposalutils.AggregateProposals(
		e,
		state.EVMMCMSStateByChain(),
		nil,
		allProposals,
		"InitChainUpgradesOnTestRouters: "+strings.Join(allChainNames, ","),
		c.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to aggregate proposals: %w", err)
	}

	if proposal == nil {
		return cldf.ChangesetOutput{Reports: allReports}, nil
	}

	return cldf.ChangesetOutput{Reports: allReports, MCMSTimelockProposals: []mcms.TimelockProposal{*proposal}}, nil
}

// PromoteChainUpgradesConfig defines the configuration for the PromoteChainUpgradesChangeset.
type PromoteChainUpgradesConfig struct {
	// HomeChainSelector is the selector of the home chain.
	HomeChainSelector uint64
	// ChainsToPromote are the chain selectors to promote
	ChainsToPromote []uint64
	// MCMSConfig is the configuration for MCMS.
	MCMSConfig *proposalutils.TimelockConfig
}

func promoteChainUpgradesPrecondition(e cldf.Environment, c PromoteChainUpgradesConfig) error {
	if c.MCMSConfig == nil {
		return errors.New("MCMSConfig must be defined")
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	err = ValidateHomeChainState(e, c.HomeChainSelector, state)
	if err != nil {
		return fmt.Errorf("failed to validate home chain state: %w", err)
	}
	// Home chain contracts are owned by MCMS.
	err = commoncs.ValidateOwnership(e.GetContext(), true, e.BlockChains.EVMChains()[c.HomeChainSelector].DeployerKey.From, state.Chains[c.HomeChainSelector].Timelock.Address(), state.Chains[c.HomeChainSelector].CCIPHome)
	if err != nil {
		return fmt.Errorf("failed to validate ownership of CCIPHome on %s: %w", e.BlockChains.EVMChains()[c.HomeChainSelector], err)
	}
	err = commoncs.ValidateOwnership(e.GetContext(), true, e.BlockChains.EVMChains()[c.HomeChainSelector].DeployerKey.From, state.Chains[c.HomeChainSelector].Timelock.Address(), state.Chains[c.HomeChainSelector].CapabilityRegistry)
	if err != nil {
		return fmt.Errorf("failed to validate ownership of CapabilityRegistry on %s: %w", e.BlockChains.EVMChains()[c.HomeChainSelector], err)
	}

	allChainSels := make(map[uint64]struct{})
	for _, chainSel := range c.ChainsToPromote {
		// Chain selector is a valid EVM chain selector & all MCMS contracts exist.
		err := stateview.ValidateChain(e, state, chainSel, c.MCMSConfig)
		if err != nil {
			return fmt.Errorf("failed to validate chain %d: %w", chainSel, err)
		}

		allChainSels[chainSel] = struct{}{}
		sourceChainSels := getSourceChainsForSelector(state, chainSel)
		for _, sourceChainSel := range sourceChainSels {
			// Source chain selector is a valid EVM chain selector & all MCMS contracts exist.
			err := stateview.ValidateChain(e, state, sourceChainSel, c.MCMSConfig)
			if err != nil {
				return fmt.Errorf("failed to validate chain %d: %w", sourceChainSel, err)
			}
			allChainSels[sourceChainSel] = struct{}{}
		}
	}

	for chainSel := range allChainSels {
		// Routers are owned by MCMS on both source and destination chains.
		err := commoncs.ValidateOwnership(e.GetContext(), true, e.BlockChains.EVMChains()[chainSel].DeployerKey.From, state.Chains[chainSel].Timelock.Address(), state.Chains[chainSel].Router)
		if err != nil {
			return fmt.Errorf("failed to validate ownership of Router on %s: %w", e.BlockChains.EVMChains()[chainSel], err)
		}
	}

	return nil
}

func promoteChainUpgradesLogic(e cldf.Environment, c PromoteChainUpgradesConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	allProposals := make([]mcms.TimelockProposal, 0)
	allReports := make([]operations.Report[any, any], 0)

	// Collect all chain names for reporting purposes.
	allChainNames := make([]string, 0, len(c.ChainsToPromote))
	for _, destChainSel := range c.ChainsToPromote {
		allChainNames = append(allChainNames, e.BlockChains.EVMChains()[destChainSel].String())
	}

	// Promote candidates commit and exec plugins for all chains.
	out, err := PromoteCandidateChangeset(e, PromoteCandidateChangesetConfig{
		HomeChainSelector: c.HomeChainSelector,
		PluginInfo: []PromoteCandidatePluginInfo{
			{
				PluginType:              types.PluginTypeCCIPCommit,
				RemoteChainSelectors:    c.ChainsToPromote,
				AllowEmptyConfigPromote: false,
			},
			{
				PluginType:              types.PluginTypeCCIPExec,
				RemoteChainSelectors:    c.ChainsToPromote,
				AllowEmptyConfigPromote: false,
			},
		},
		MCMS: c.MCMSConfig,
	})
	allReports = append(allReports, out.Reports...)
	if err != nil {
		return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run PromoteCandidateChangeset: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	// Connect each destination to each of its sources via main routers
	ownershipAlreadyEnsured := make(map[uint64]map[common.Address]struct{})
	for _, destChainSel := range c.ChainsToPromote {
		// Assemble source chains for the destination chain, using 1.5.0 OnRamps.
		destChain := e.BlockChains.EVMChains()[destChainSel]

		// Transfer ownership of OffRamp on destination chain.
		out, err := ensureTimelockOwnership(e, destChainSel, []commoncs.Ownable{
			state.Chains[destChainSel].OffRamp,
		}, *c.MCMSConfig)
		allReports = append(allReports, out.Reports...)
		if err != nil {
			return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to ensure timelock ownership of OffRamp on %s: %w", destChain, err)
		}
		allProposals = append(allProposals, out.MCMSTimelockProposals...)

		// Loop through each source connected to the destination chain.
		sourceChainsToConnect := getSourceChainsForSelector(state, destChainSel)
		for _, sourceChainSel := range sourceChainsToConnect {
			sourceChain := e.BlockChains.EVMChains()[sourceChainSel]

			// Transfer ownership of OnRamp on source chain.
			// We need to track if we already ensured ownership of the OnRamp on the source chain,
			// as multiple destination chains can have overlapping source chains.
			if ownershipAlreadyEnsured[sourceChainSel] == nil {
				ownershipAlreadyEnsured[sourceChainSel] = make(map[common.Address]struct{})
			}
			if _, ok := ownershipAlreadyEnsured[sourceChainSel][state.Chains[sourceChainSel].OnRamp.Address()]; !ok {
				out, err = ensureTimelockOwnership(e, sourceChainSel, []commoncs.Ownable{
					state.Chains[sourceChainSel].OnRamp,
				}, *c.MCMSConfig)
				allReports = append(allReports, out.Reports...)
				if err != nil {
					return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to ensure timelock ownership of OnRamp on %s: %w", sourceChain, err)
				}
				allProposals = append(allProposals, out.MCMSTimelockProposals...)
				ownershipAlreadyEnsured[sourceChainSel][state.Chains[sourceChainSel].OnRamp.Address()] = struct{}{}
			}

			/*
				The ordering of the following changesets is important so we don't disrupt traffic:

				On source:
				1. Update the OnRamp destination config, pointing at the main router.
				2. Set the OnRamp on the main router (if we did this first, OnRamp wouldn't be ready and users would see reverts).

				On dest:
				1. Add the OffRamp to the main router.
				2. Update the OffRamp source config, pointing at the main router (if we did this first, there is a chance that incoming traffic would hit the new OffRamp before the router gets updated & would see "OffRamp not supported" errors).
			*/

			// Update OnRamp 1.6.0 on source chain (use main router).
			out, err = UpdateOnRampsDestsChangeset(e, UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]OnRampDestinationUpdate{
					sourceChainSel: {
						destChainSel: {
							TestRouter:       false,
							AllowListEnabled: false,
							IsEnabled:        true,
						},
					},
				},
				SkipOwnershipCheck: true, // We already ensured desired ownership above.
				MCMS:               c.MCMSConfig,
			})
			allReports = append(allReports, out.Reports...)
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run UpdateOnRampsDestsChangeset on %s: %w", sourceChain, err)
			}
			allProposals = append(allProposals, out.MCMSTimelockProposals...)

			// Update OnRamp to router on source, OffRamp to router on destination (use main router).
			out, err = UpdateRouterRampsChangeset(e, UpdateRouterRampsConfig{
				TestRouter:         false,
				MCMS:               c.MCMSConfig,
				SkipOwnershipCheck: true, // We already ensured desired ownership above.
				UpdatesByChain: map[uint64]RouterUpdates{
					destChainSel: {
						OffRampUpdates: map[uint64]bool{sourceChainSel: true},
					},
					sourceChainSel: {
						OnRampUpdates: map[uint64]bool{destChainSel: true},
					},
				},
			})
			allReports = append(allReports, out.Reports...)
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run UpdateRouterRampsChangeset and connect %s with %s: %w", sourceChain, destChain, err)
			}
			allProposals = append(allProposals, out.MCMSTimelockProposals...)

			// Update OffRamp 1.6.0 on destination chain (use main router, no RMN verification).
			out, err = UpdateOffRampSourcesChangeset(e, UpdateOffRampSourcesConfig{
				UpdatesByChain: map[uint64]map[uint64]OffRampSourceUpdate{
					destChainSel: {
						sourceChainSel: {
							TestRouter:                false,
							IsRMNVerificationDisabled: true, // TODO: Might not always be true, but for now we assume it is.
							IsEnabled:                 true,
						},
					},
				},
				SkipOwnershipCheck: true, // We already ensured desired ownership above.
				MCMS:               c.MCMSConfig,
			})
			allReports = append(allReports, out.Reports...)
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run UpdateOffRampSourcesChangeset on %s: %w", destChain, err)
			}
			allProposals = append(allProposals, out.MCMSTimelockProposals...)
		}
	}

	proposal, err := proposalutils.AggregateProposals(
		e,
		state.EVMMCMSStateByChain(),
		nil,
		allProposals,
		"PromoteChainUpgradesToMainRoutersChangeset: "+strings.Join(allChainNames, ","),
		c.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to aggregate proposals: %w", err)
	}

	if proposal == nil {
		return cldf.ChangesetOutput{Reports: allReports}, nil
	}

	return cldf.ChangesetOutput{Reports: allReports, MCMSTimelockProposals: []mcms.TimelockProposal{*proposal}}, nil
}

func getSourceChainsForSelector(state stateview.CCIPOnChainState, chainSel uint64) []uint64 {
	sourceChains := make([]uint64, 0)

	for sourceChainSel, sourceChainState := range state.Chains {
		if sourceChainSel == chainSel {
			continue // Skip the destination chain itself.
		}
		for destChainSel := range sourceChainState.EVM2EVMOnRamp {
			if destChainSel == chainSel {
				// Source chain has a 1.5.0 OnRamp to the destination chain.
				sourceChains = append(sourceChains, sourceChainSel)
			}
		}
	}

	return sourceChains
}

func ensureTimelockOwnership(e cldf.Environment, chainSel uint64, contracts []commoncs.Ownable, mcmsCfg proposalutils.TimelockConfig) (cldf.ChangesetOutput, error) {
	addressesToTransfer := make([]common.Address, 0, len(contracts))
	for _, contract := range contracts {
		if contract == nil {
			continue
		}
		owner, err := contract.Owner(&bind.CallOpts{Context: e.GetContext()})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get owner of contract %s on %d: %w", contract.Address().Hex(), chainSel, err)
		}
		if owner == e.BlockChains.EVMChains()[chainSel].DeployerKey.From {
			addressesToTransfer = append(addressesToTransfer, contract.Address())
		}
	}
	if len(addressesToTransfer) == 0 {
		return cldf.ChangesetOutput{}, nil // Nothing to transfer, no ownership change needed.
	}
	return commoncs.TransferToMCMSWithTimelockV2(e, commoncs.TransferToMCMSWithTimelockConfig{
		ContractsByChain: map[uint64][]common.Address{
			chainSel: addressesToTransfer,
		},
		MCMSConfig: mcmsCfg,
	})
}
