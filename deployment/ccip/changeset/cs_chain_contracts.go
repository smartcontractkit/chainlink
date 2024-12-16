package changeset

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

var _ deployment.ChangeSet[UpdateOnRampsDestsConfig] = UpdateOnRampsDests
var _ deployment.ChangeSet[UpdateOffRampsSourcesConfig] = UpdateOffRampSources
var _ deployment.ChangeSet[UpdateRoutersOffRampsConfig] = UpdateRoutersOffRamps

type UpdateOnRampsDestsConfig struct {
	UpdatesByChain map[uint64][]OnRampDestinationUpdate
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be acheived by calling this function multiple times)
	MCMS *MCMSConfig
}

type OnRampDestinationUpdate struct {
	DestinationSelector uint64
	IsEnabled           bool // If false, disables the destination by setting router to 0x0.
	TestRouter          bool // Flag for safety only allow specifying either router or testRouter.
	AllowListEnabled    bool
}

func (cfg UpdateOnRampsDestsConfig) Validate(ctx context.Context, state CCIPOnChainState) error {
	supportedChains := state.SupportedChains()
	for chainSel, updates := range cfg.UpdatesByChain {
		chainState, ok := state.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain %d not found in onchain state", chainSel)
		}
		if chainState.TestRouter == nil {
			return fmt.Errorf("missing test router for chain %d", chainSel)
		}
		if chainState.Router == nil {
			return fmt.Errorf("missing router for chain %d", chainSel)
		}
		if chainState.OnRamp == nil {
			return fmt.Errorf("missing onramp onramp for chain %d", chainSel)
		}
		for _, update := range updates {
			if update.DestinationSelector == 0 {
				return fmt.Errorf("destination selector cannot be 0")
			}
			sc, err := chainState.OnRamp.GetStaticConfig(&bind.CallOpts{Context: ctx})
			if err != nil {
				return fmt.Errorf("failed to get onramp static config %s: %w", chainState.OnRamp.Address(), err)
			}
			if update.DestinationSelector == sc.ChainSelector {
				return fmt.Errorf("cannot update onramp destination to the same chain")
			}
			// Destination cannot be an unknown destination.
			if _, ok := supportedChains[update.DestinationSelector]; !ok {
				return fmt.Errorf("destination chain %d is not a supported %s", update.DestinationSelector, chainState.OnRamp.Address())
			}
		}
	}
	return nil
}

// UpdateOnRampsDests updates the onramp destinations for each onramp
// in the chains specified. Multichain support is important - consider when we add a new chain
// and need to update the onramp destinations for all chains to support the new chain.
func UpdateOnRampsDests(e deployment.Environment, cfg UpdateOnRampsDestsConfig) (deployment.ChangesetOutput, error) {
	s, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e.GetContext(), s); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	var batches []timelock.BatchChainOperation
	timelocks := make(map[uint64]common.Address)
	proposers := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for chainSel, updates := range cfg.UpdatesByChain {
		txOpts := e.Chains[chainSel].DeployerKey
		txOpts.Context = e.GetContext()
		if cfg.MCMS != nil {
			txOpts = deployment.SimTransactOpts()
		}
		onRamp := s.Chains[chainSel].OnRamp
		var args []onramp.OnRampDestChainConfigArgs
		for _, update := range updates {
			router := common.HexToAddress("0x0")
			// If not enabled, set router to 0x0.
			if update.IsEnabled {
				if update.TestRouter {
					router = s.Chains[chainSel].TestRouter.Address()
				} else {
					router = s.Chains[chainSel].Router.Address()
				}
			}
			args = append(args, onramp.OnRampDestChainConfigArgs{
				DestChainSelector: update.DestinationSelector,
				Router:            router,
				AllowlistEnabled:  update.AllowListEnabled,
			})
		}
		tx, err := onRamp.ApplyDestChainConfigUpdates(txOpts, args)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		if cfg.MCMS == nil {
			if _, err := deployment.ConfirmIfNoError(e.Chains[chainSel], tx, err); err != nil {
				return deployment.ChangesetOutput{}, err
			}
		} else {
			batches = append(batches, timelock.BatchChainOperation{
				ChainIdentifier: mcms.ChainIdentifier(chainSel),
				Batch: []mcms.Operation{
					{
						To:    onRamp.Address(),
						Data:  tx.Data(),
						Value: big.NewInt(0),
					},
				},
			})
			timelocks[chainSel] = s.Chains[chainSel].Timelock.Address()
			proposers[chainSel] = s.Chains[chainSel].ProposerMcm
		}
	}
	if cfg.MCMS == nil {
		return deployment.ChangesetOutput{}, nil
	}

	p, err := proposalutils.BuildProposalFromBatches(
		timelocks,
		proposers,
		batches,
		"Update onramp destinations",
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{
		*p,
	}}, nil
}

type UpdateOffRampsSourcesConfig struct {
	UpdatesByChain map[uint64][]OffRampSourceUpdate
	MCMS           *MCMSConfig
}

type OffRampSourceUpdate struct {
	SourceSelector uint64 // Note will look up the relevant onramp and set it
	IsEnabled      bool   // If false, disables the source by setting router to 0x0.
	TestRouter     bool   // Flag for safety only allow specifying either router or testRouter.
}

func (cfg UpdateOffRampsSourcesConfig) Validate(ctx context.Context, state CCIPOnChainState) error {
	supportedChains := state.SupportedChains()
	for chainSel, updates := range cfg.UpdatesByChain {
		chainState, ok := state.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain %d not found in onchain state", chainSel)
		}
		if chainState.TestRouter == nil {
			return fmt.Errorf("missing test router for chain %d", chainSel)
		}
		if chainState.Router == nil {
			return fmt.Errorf("missing router for chain %d", chainSel)
		}
		if chainState.OffRamp == nil {
			return fmt.Errorf("missing onramp onramp for chain %d", chainSel)
		}
		for _, update := range updates {
			if update.SourceSelector == chainSel {
				return fmt.Errorf("cannot update offramp source to the same chain %d", update.SourceSelector)
			}
			sourceChain := state.Chains[update.SourceSelector]
			// Source chain must have the onramp deployed.
			// Note this also validates the specified source selector.
			if sourceChain.OnRamp == nil {
				return fmt.Errorf("missing onramp for source %d", update.SourceSelector)
			}
			// Source cannot be an unknown
			if _, ok := supportedChains[update.SourceSelector]; !ok {
				return fmt.Errorf("source chain %d is not a supported chain %s", update.SourceSelector, chainState.OffRamp.Address())
			}
		}
	}
	return nil
}

// UpdateOffRampSources updates the offramp sources for each offramp.
func UpdateOffRampSources(e deployment.Environment, cfg UpdateOffRampsSourcesConfig) (deployment.ChangesetOutput, error) {
	s, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e.GetContext(), s); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	var batches []timelock.BatchChainOperation
	timelocks := make(map[uint64]common.Address)
	proposers := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for chainSel, updates := range cfg.UpdatesByChain {
		txOpts := e.Chains[chainSel].DeployerKey
		txOpts.Context = e.GetContext()
		if cfg.MCMS != nil {
			txOpts = deployment.SimTransactOpts()
		}
		offRamp := s.Chains[chainSel].OffRamp
		var args []offramp.OffRampSourceChainConfigArgs
		for _, update := range updates {
			router := common.HexToAddress("0x0")
			if update.IsEnabled {
				if update.TestRouter {
					router = s.Chains[chainSel].TestRouter.Address()
				} else {
					router = s.Chains[chainSel].Router.Address()
				}
			}
			onRamp := s.Chains[update.SourceSelector].OnRamp
			args = append(args, offramp.OffRampSourceChainConfigArgs{
				SourceChainSelector: update.SourceSelector,
				Router:              router,
				IsEnabled:           update.IsEnabled,
				OnRamp:              common.LeftPadBytes(onRamp.Address().Bytes(), 32),
			})
		}
		tx, err := offRamp.ApplySourceChainConfigUpdates(txOpts, args)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		if cfg.MCMS == nil {
			if _, err := deployment.ConfirmIfNoError(e.Chains[chainSel], tx, err); err != nil {
				return deployment.ChangesetOutput{}, err
			}
		} else {
			batches = append(batches, timelock.BatchChainOperation{
				ChainIdentifier: mcms.ChainIdentifier(chainSel),
				Batch: []mcms.Operation{
					{
						To:    offRamp.Address(),
						Data:  tx.Data(),
						Value: big.NewInt(0),
					},
				},
			})
			timelocks[chainSel] = s.Chains[chainSel].Timelock.Address()
			proposers[chainSel] = s.Chains[chainSel].ProposerMcm
		}
	}
	if cfg.MCMS == nil {
		return deployment.ChangesetOutput{}, nil
	}

	p, err := proposalutils.BuildProposalFromBatches(
		timelocks,
		proposers,
		batches,
		"Update offramp sources",
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{
		*p,
	}}, nil
}

type UpdateRoutersOffRampsConfig struct {
	UpdatesByChain map[uint64][]RouterOffRampUpdate
	TestRouter     bool // Disallow mixing test router/non-test router per chain for simplicity.
	MCMS           *MCMSConfig
}

type RouterOffRampUpdate struct {
	SourceSelector uint64 // Note will look up the relevant onramp and set it
	IsEnabled      bool   // If false, disables the source by setting router to 0x0.
}

func (cfg UpdateRoutersOffRampsConfig) Validate(ctx context.Context, state CCIPOnChainState) error {
	supportedChains := state.SupportedChains()
	for chainSel, updates := range cfg.UpdatesByChain {
		chainState, ok := state.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain %d not found in onchain state", chainSel)
		}
		if chainState.TestRouter == nil {
			return fmt.Errorf("missing test router for chain %d", chainSel)
		}
		if chainState.Router == nil {
			return fmt.Errorf("missing router for chain %d", chainSel)
		}
		if chainState.OffRamp == nil {
			return fmt.Errorf("missing onramp onramp for chain %d", chainSel)
		}
		for _, update := range updates {
			if update.SourceSelector == chainSel {
				return fmt.Errorf("cannot update offramp source to the same chain %d", update.SourceSelector)
			}
			sourceChain := state.Chains[update.SourceSelector]
			// Source chain must have the onramp deployed.
			// Note this also validates the specified source selector.
			if sourceChain.OnRamp == nil {
				return fmt.Errorf("missing onramp for source %d", update.SourceSelector)
			}
			// Source cannot be an unknown
			if _, ok := supportedChains[update.SourceSelector]; !ok {
				return fmt.Errorf("source chain %d is not a supported chain %s", update.SourceSelector, chainState.OffRamp.Address())
			}
		}
	}
	return nil
}

// UpdateRoutersOffRamps updates the enabled offramps
// in either the router or test router. The general upgrade path
// for an offramp would be first we enable it on the test router,
// ensure it works e2e. Then we enable that source on the real router.
func UpdateRoutersOffRamps(e deployment.Environment, cfg UpdateRoutersOffRampsConfig) (deployment.ChangesetOutput, error) {
	s, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e.GetContext(), s); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	var batches []timelock.BatchChainOperation
	timelocks := make(map[uint64]common.Address)
	proposers := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for chainSel, updates := range cfg.UpdatesByChain {
		txOpts := e.Chains[chainSel].DeployerKey
		txOpts.Context = e.GetContext()
		if cfg.MCMS != nil {
			txOpts = deployment.SimTransactOpts()
		}
		routerC := s.Chains[chainSel].Router
		if cfg.TestRouter {
			routerC = s.Chains[chainSel].TestRouter
		}
		// Note if we add distinct offramps per source to the state,
		// we'll need to add support here for looking them up.
		// For now its simple, all sources use the same offramp.
		offRamp := s.Chains[chainSel].OffRamp
		var removes, adds []router.RouterOffRamp
		for _, update := range updates {
			if update.IsEnabled {
				adds = append(adds, router.RouterOffRamp{
					SourceChainSelector: update.SourceSelector,
					OffRamp:             offRamp.Address(),
				})
			} else {
				removes = append(removes, router.RouterOffRamp{
					SourceChainSelector: update.SourceSelector,
					OffRamp:             offRamp.Address(),
				})
			}
		}
		tx, err := routerC.ApplyRampUpdates(txOpts, []router.RouterOnRamp{}, removes, adds)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		if cfg.MCMS == nil {
			if _, err := deployment.ConfirmIfNoError(e.Chains[chainSel], tx, err); err != nil {
				return deployment.ChangesetOutput{}, err
			}
		} else {
			batches = append(batches, timelock.BatchChainOperation{
				ChainIdentifier: mcms.ChainIdentifier(chainSel),
				Batch: []mcms.Operation{
					{
						To:    routerC.Address(),
						Data:  tx.Data(),
						Value: big.NewInt(0),
					},
				},
			})
			timelocks[chainSel] = s.Chains[chainSel].Timelock.Address()
			proposers[chainSel] = s.Chains[chainSel].ProposerMcm
		}
	}
	if cfg.MCMS == nil {
		return deployment.ChangesetOutput{}, nil
	}

	p, err := proposalutils.BuildProposalFromBatches(
		timelocks,
		proposers,
		batches,
		"Update router offramps",
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{
		*p,
	}}, nil
}
