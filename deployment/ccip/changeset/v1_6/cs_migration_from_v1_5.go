package v1_6

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/mcms"
)

var (
	// InitChainUpgratesOnTestRoutersChangeset sets candidates for the commit and exec DONs for multiple destination chains.
	// It then identifies all existing 1.5.0 source chains for each chain in the batch.
	// For each 1.5.0 OnRamp connecting to a destination, configuration gets translated to the 1.6.0 FeeQuoter.
	// In addition, OnRamps are connected to destination chains via test routers.
	// We do NOT connect the destinations back to the source chains, as DONs are not guaranteed to exist for sources.
	// This changeset is NOT IDEMPOTENT - if AddDON is called more than once for the same chain it will revert.
	InitChainUpgradesOnTestRoutersChangeset = cldf.CreateChangeSet(
		initChainUpgradesOnTestRoutersLogic,
		initChainUpgradesOnTestRoutersPrecondition,
	)
	// PromoteChainUpgradesToMainRoutersChangeset promotes the commit and exec DON candidates for multiple destination chains.
	// It then connects the source chains to the destination chains via main routers.
	// Before running PromoteChainUpgradesToMainRoutersChangeset for a batch, you must run InitChainUpgradesOnTestRoutersChangeset followed by SetOCR3OffRampChangeset.
	// SetOCR3OffRampChangeset should be run with ConfigType set to candidate, since the config won't be promoted until this changeset is run.
	// This changeset is NOT IDEMPOTENT - re-promoting will result in clearing the active digest, which is not desired.
	PromoteChainUpgradesToMainRoutersChangeset = cldf.CreateChangeSet(
		promoteChainUpgradesToMainRoutersLogic,
		promoteChainUpgradesToMainRoutersPrecondition,
	)
)

type ChainUpgradeConfig struct {
	CommitOCRParams CCIPOCRParams
	ExecOCRParams   CCIPOCRParams
}

type InitChainUpgradesOnTestRoutersConfig struct {
	HomeChainSelector uint64
	FeedChainSelector uint64
	ChainsToUpgrade   map[uint64]ChainUpgradeConfig
	MCMSConfig        *proposalutils.TimelockConfig
}

func initChainUpgradesOnTestRoutersPrecondition(e cldf.Environment, c InitChainUpgradesOnTestRoutersConfig) error {
	// Chains must all be EVM

	return nil
}

func initChainUpgradesOnTestRoutersLogic(e cldf.Environment, c InitChainUpgradesOnTestRoutersConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	donID, err := state.Chains[c.HomeChainSelector].CapabilityRegistry.GetNextDONId(&bind.CallOpts{Context: e.GetContext()})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get next DON ID: %w", err)
	}
	allProposals := make([]mcms.TimelockProposal, 0)
	allReports := make([]operations.Report[any, any], 0)
	allChainNames := make([]string, 0, len(c.ChainsToUpgrade))

	for chainSel, chainUpgradeCfg := range c.ChainsToUpgrade {
		chain := e.BlockChains.EVMChains()[chainSel]
		allChainNames = append(allChainNames, chain.String())

		// Add DON & set candidate for commit plugin
		out, err := AddDonAndSetCandidateChangeset(e, AddDonAndSetCandidateChangesetConfig{
			SetCandidateConfigBase: SetCandidateConfigBase{
				HomeChainSelector: c.HomeChainSelector,
				FeedChainSelector: c.FeedChainSelector,
				MCMS:              c.MCMSConfig,
			},
			PluginInfo: SetCandidatePluginInfo{
				PluginType: types.PluginTypeCCIPCommit,
				OCRConfigPerRemoteChainSelector: map[uint64]CCIPOCRParams{
					chainSel: chainUpgradeCfg.CommitOCRParams,
				},
			},
			DonIDOverride: donID,
		})
		allReports = append(allReports, out.Reports...)
		if err != nil {
			return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run AddDonAndSetCandidateChangeset for commit plugin on %s: %w", chain, err)
		}
		allProposals = append(allProposals, out.MCMSTimelockProposals...)

		// Set candidate for exec plugin
		out, err = SetCandidateChangeset(e, SetCandidateChangesetConfig{
			SetCandidateConfigBase: SetCandidateConfigBase{
				HomeChainSelector: c.HomeChainSelector,
				FeedChainSelector: c.FeedChainSelector,
				MCMS:              c.MCMSConfig,
			},
			PluginInfo: []SetCandidatePluginInfo{
				{
					PluginType: types.PluginTypeCCIPExec,
					OCRConfigPerRemoteChainSelector: map[uint64]CCIPOCRParams{
						chainSel: chainUpgradeCfg.ExecOCRParams,
					},
					SkipChainConfigValidation: true,
				},
			},
			DonIDOverrides: map[uint64]uint32{chainSel: donID},
		})
		allReports = append(allReports, out.Reports...)
		if err != nil {
			return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run SetCandidateChangeset for exec plugin on %s: %w", chain, err)
		}
		allProposals = append(allProposals, out.MCMSTimelockProposals...)

		// Increment the DON ID for the next chain
		donID++
	}

	for range c.ChainsToUpgrade {
		// Process each chain concurrently, no sequential dependency here.

		// Import configuration from 1.5.0 (TODO)
		// Connect all sources to the destination via test routers (TODO)
	}

	proposal, err := proposalutils.AggregateProposals(
		e,
		state.EVMMCMSStateByChain(),
		nil,
		allProposals,
		fmt.Sprintf("InitChainUpgradesOnTestRouters: %s", strings.Join(allChainNames, ",")),
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

type PromoteChainUpgradesToMainRoutersConfig struct {
	HomeChainSelector uint64
	FeedChainSelector uint64
	ChainsToPromote   []uint64
	MCMSConfig        *proposalutils.TimelockConfig
}

func promoteChainUpgradesToMainRoutersPrecondition(e cldf.Environment, c PromoteChainUpgradesToMainRoutersConfig) error {
	// Chains must all be EVM

	return nil
}

func promoteChainUpgradesToMainRoutersLogic(e cldf.Environment, c PromoteChainUpgradesToMainRoutersConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	allProposals := make([]mcms.TimelockProposal, 0)
	allReports := make([]operations.Report[any, any], 0)
	allChainNames := make([]string, 0, len(c.ChainsToPromote))
	for _, chainSel := range c.ChainsToPromote {
		allChainNames = append(allChainNames, e.BlockChains.EVMChains()[chainSel].String())
	}

	// Promote candidates commit and exec plugins for all chains
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

	// Connect each destination to all its sources via main routers
	for _, chainSel := range c.ChainsToPromote {
		// TODO: This can be concurrent across each destination chain, no sequential dependency here.
		// Assemble source chains for the destination chain, using 1.5.0 OnRamps.
		chainState := state.Chains[chainSel]
		chain := e.BlockChains.EVMChains()[chainSel]
		sourceChainsToConnect := make([]uint64, 0)

		for sourceChainSel := range state.Chains {
			if sourceChainSel == chainSel {
				continue // Skip the destination chain itself.
			}
			for destChainSel := range chainState.EVM2EVMOnRamp {
				if destChainSel == chainSel {
					// Source chain has a 1.5.0 OnRamp to the destination chain.
					sourceChainsToConnect = append(sourceChainsToConnect, sourceChainSel)
				}
			}
		}

		for _, sourceChainSel := range sourceChainsToConnect {
			// TODO: This can be concurrent across each source chain, no sequential dependency here.
			sourceChain := e.BlockChains.EVMChains()[sourceChainSel]

			// Update OnRamp 1.6.0 on source chain (use main router).
			out, err := UpdateOnRampsDestsChangeset(e, UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]OnRampDestinationUpdate{
					sourceChainSel: {
						chainSel: {
							TestRouter:       false,
							AllowListEnabled: false,
							IsEnabled:        true,
						},
					},
				},
				MCMS: c.MCMSConfig,
			})
			allReports = append(allReports, out.Reports...)
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run UpdateOnRampsDestsChangeset on %s: %w", sourceChain, err)
			}
			allProposals = append(allProposals, out.MCMSTimelockProposals...)

			// Update OffRamp 1.6.0 on destination chain (use main router, no RMN verification).
			out, err = UpdateOffRampSourcesChangeset(e, UpdateOffRampSourcesConfig{
				UpdatesByChain: map[uint64]map[uint64]OffRampSourceUpdate{
					sourceChainSel: {
						chainSel: {
							TestRouter:                false,
							IsRMNVerificationDisabled: true,
							IsEnabled:                 true,
						},
					},
				},
				MCMS: c.MCMSConfig,
			})
			allReports = append(allReports, out.Reports...)
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run UpdateOffRampSourcesChangeset on %s: %w", chain, err)
			}
			allProposals = append(allProposals, out.MCMSTimelockProposals...)

			// Update Router on source and destination chains (use main router).
			out, err = UpdateRouterRampsChangeset(e, UpdateRouterRampsConfig{
				TestRouter: false,
				MCMS:       c.MCMSConfig,
				UpdatesByChain: map[uint64]RouterUpdates{
					chainSel: {
						OnRampUpdates:  map[uint64]bool{sourceChainSel: true},
						OffRampUpdates: map[uint64]bool{sourceChainSel: true},
					},
					sourceChainSel: {
						OnRampUpdates:  map[uint64]bool{chainSel: true},
						OffRampUpdates: map[uint64]bool{chainSel: true},
					},
				},
			})
			allReports = append(allReports, out.Reports...)
			if err != nil {
				return cldf.ChangesetOutput{Reports: allReports}, fmt.Errorf("failed to run UpdateRouterRampsChangeset and connect %s with %s: %w", sourceChain, chain, err)
			}
			allProposals = append(allProposals, out.MCMSTimelockProposals...)
		}
	}

	proposal, err := proposalutils.AggregateProposals(
		e,
		state.EVMMCMSStateByChain(),
		nil,
		allProposals,
		fmt.Sprintf("PromoteChainUpgradesToMainRoutersChangeset: %s", strings.Join(allChainNames, ",")),
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
