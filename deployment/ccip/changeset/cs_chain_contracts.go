package changeset

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

var (
	_ deployment.ChangeSet[UpdateOnRampDestsConfig]    = UpdateOnRampsDests
	_ deployment.ChangeSet[UpdateOffRampSourcesConfig] = UpdateOffRampSources
	_ deployment.ChangeSet[UpdateRouterRampsConfig]    = UpdateRouterRamps
	_ deployment.ChangeSet[UpdateFeeQuoterDestsConfig] = UpdateFeeQuoterDests
	_ deployment.ChangeSet[SetOCR3OffRampConfig]       = SetOCR3OffRamp
)

type UpdateOnRampDestsConfig struct {
	UpdatesByChain map[uint64]map[uint64]OnRampDestinationUpdate
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be acheived by calling this function multiple times)
	MCMS *MCMSConfig
}

type OnRampDestinationUpdate struct {
	IsEnabled        bool // If false, disables the destination by setting router to 0x0.
	TestRouter       bool // Flag for safety only allow specifying either router or testRouter.
	AllowListEnabled bool
}

func (cfg UpdateOnRampDestsConfig) Validate(e deployment.Environment) error {
	state, err := LoadOnchainState(e)
	if err != nil {
		return err
	}
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
		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.Chains[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.OnRamp); err != nil {
			return err
		}

		for destination := range updates {
			// Destination cannot be an unknown destination.
			if _, ok := supportedChains[destination]; !ok {
				return fmt.Errorf("destination chain %d is not a supported %s", destination, chainState.OnRamp.Address())
			}
			sc, err := chainState.OnRamp.GetStaticConfig(&bind.CallOpts{Context: e.GetContext()})
			if err != nil {
				return fmt.Errorf("failed to get onramp static config %s: %w", chainState.OnRamp.Address(), err)
			}
			if destination == sc.ChainSelector {
				return fmt.Errorf("cannot update onramp destination to the same chain")
			}
		}
	}
	return nil
}

// UpdateOnRampsDests updates the onramp destinations for each onramp
// in the chains specified. Multichain support is important - consider when we add a new chain
// and need to update the onramp destinations for all chains to support the new chain.
func UpdateOnRampsDests(e deployment.Environment, cfg UpdateOnRampDestsConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	s, err := LoadOnchainState(e)
	if err != nil {
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
		for destination, update := range updates {
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
				DestChainSelector: destination,
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

type UpdateFeeQuoterDestsConfig struct {
	UpdatesByChain map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be acheived by calling this function multiple times)
	MCMS *MCMSConfig
}

func (cfg UpdateFeeQuoterDestsConfig) Validate(e deployment.Environment) error {
	state, err := LoadOnchainState(e)
	if err != nil {
		return err
	}
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
		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.Chains[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.FeeQuoter); err != nil {
			return err
		}

		for destination := range updates {
			// Destination cannot be an unknown destination.
			if _, ok := supportedChains[destination]; !ok {
				return fmt.Errorf("destination chain %d is not a supported %s", destination, chainState.OnRamp.Address())
			}
			sc, err := chainState.OnRamp.GetStaticConfig(&bind.CallOpts{Context: e.GetContext()})
			if err != nil {
				return fmt.Errorf("failed to get onramp static config %s: %w", chainState.OnRamp.Address(), err)
			}
			if destination == sc.ChainSelector {
				return fmt.Errorf("cannot update onramp destination to the same chain")
			}
		}
	}
	return nil
}

func UpdateFeeQuoterDests(e deployment.Environment, cfg UpdateFeeQuoterDestsConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	s, err := LoadOnchainState(e)
	if err != nil {
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
		fq := s.Chains[chainSel].FeeQuoter
		var args []fee_quoter.FeeQuoterDestChainConfigArgs
		for destination, dc := range updates {
			args = append(args, fee_quoter.FeeQuoterDestChainConfigArgs{
				DestChainSelector: destination,
				DestChainConfig:   dc,
			})
		}
		tx, err := fq.ApplyDestChainConfigUpdates(txOpts, args)
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
						To:    fq.Address(),
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
		"Update fq destinations",
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{
		*p,
	}}, nil
}

type UpdateOffRampSourcesConfig struct {
	UpdatesByChain map[uint64]map[uint64]OffRampSourceUpdate
	MCMS           *MCMSConfig
}

type OffRampSourceUpdate struct {
	IsEnabled  bool // If false, disables the source by setting router to 0x0.
	TestRouter bool // Flag for safety only allow specifying either router or testRouter.
}

func (cfg UpdateOffRampSourcesConfig) Validate(e deployment.Environment) error {
	state, err := LoadOnchainState(e)
	if err != nil {
		return err
	}
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
		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.Chains[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.OffRamp); err != nil {
			return err
		}

		for source := range updates {
			// Source cannot be an unknown
			if _, ok := supportedChains[source]; !ok {
				return fmt.Errorf("source chain %d is not a supported chain %s", source, chainState.OffRamp.Address())
			}

			if source == chainSel {
				return fmt.Errorf("cannot update offramp source to the same chain %d", source)
			}
			sourceChain := state.Chains[source]
			// Source chain must have the onramp deployed.
			// Note this also validates the specified source selector.
			if sourceChain.OnRamp == nil {
				return fmt.Errorf("missing onramp for source %d", source)
			}
		}
	}
	return nil
}

// UpdateOffRampSources updates the offramp sources for each offramp.
func UpdateOffRampSources(e deployment.Environment, cfg UpdateOffRampSourcesConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	s, err := LoadOnchainState(e)
	if err != nil {
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
		for source, update := range updates {
			router := common.HexToAddress("0x0")
			if update.IsEnabled {
				if update.TestRouter {
					router = s.Chains[chainSel].TestRouter.Address()
				} else {
					router = s.Chains[chainSel].Router.Address()
				}
			}
			onRamp := s.Chains[source].OnRamp
			args = append(args, offramp.OffRampSourceChainConfigArgs{
				SourceChainSelector: source,
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

type UpdateRouterRampsConfig struct {
	// TestRouter means the updates will be applied to the test router
	// on all chains. Disallow mixing test router/non-test router per chain for simplicity.
	TestRouter     bool
	UpdatesByChain map[uint64]RouterUpdates
	MCMS           *MCMSConfig
}

type RouterUpdates struct {
	OffRampUpdates map[uint64]bool
	OnRampUpdates  map[uint64]bool
}

func (cfg UpdateRouterRampsConfig) Validate(e deployment.Environment) error {
	state, err := LoadOnchainState(e)
	if err != nil {
		return err
	}
	supportedChains := state.SupportedChains()
	for chainSel, update := range cfg.UpdatesByChain {
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
		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.Chains[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.Router); err != nil {
			return err
		}

		for source := range update.OffRampUpdates {
			// Source cannot be an unknown
			if _, ok := supportedChains[source]; !ok {
				return fmt.Errorf("source chain %d is not a supported chain %s", source, chainState.OffRamp.Address())
			}
			if source == chainSel {
				return fmt.Errorf("cannot update offramp source to the same chain %d", source)
			}
			sourceChain := state.Chains[source]
			// Source chain must have the onramp deployed.
			// Note this also validates the specified source selector.
			if sourceChain.OnRamp == nil {
				return fmt.Errorf("missing onramp for source %d", source)
			}
		}
		for destination := range update.OnRampUpdates {
			// Source cannot be an unknown
			if _, ok := supportedChains[destination]; !ok {
				return fmt.Errorf("dest chain %d is not a supported chain %s", destination, chainState.OffRamp.Address())
			}
			if destination == chainSel {
				return fmt.Errorf("cannot update onRamp dest to the same chain %d", destination)
			}
			destChain := state.Chains[destination]
			if destChain.OffRamp == nil {
				return fmt.Errorf("missing offramp for dest %d", destination)
			}
		}

	}
	return nil
}

// UpdateRouterRamps updates the on/offramps
// in either the router or test router for a series of chains. Use cases include:
// - Ramp upgrade. After deploying new ramps you can enable them on the test router and
// ensure it works e2e. Then enable the ramps on the real router.
// - New chain support. When adding a new chain, you can enable the new destination
// on all chains to support the new chain through the test router first. Once tested,
// Enable the new destination on the real router.
func UpdateRouterRamps(e deployment.Environment, cfg UpdateRouterRampsConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	s, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	var batches []timelock.BatchChainOperation
	timelocks := make(map[uint64]common.Address)
	proposers := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for chainSel, update := range cfg.UpdatesByChain {
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
		for source, enabled := range update.OffRampUpdates {
			if enabled {
				adds = append(adds, router.RouterOffRamp{
					SourceChainSelector: source,
					OffRamp:             offRamp.Address(),
				})
			} else {
				removes = append(removes, router.RouterOffRamp{
					SourceChainSelector: source,
					OffRamp:             offRamp.Address(),
				})
			}
		}
		// Ditto here, only one onramp expected until 1.7.
		onRamp := s.Chains[chainSel].OnRamp
		var onRampUpdates []router.RouterOnRamp
		for dest, enabled := range update.OnRampUpdates {
			if enabled {
				onRampUpdates = append(onRampUpdates, router.RouterOnRamp{
					DestChainSelector: dest,
					OnRamp:            onRamp.Address(),
				})
			} else {
				onRampUpdates = append(onRampUpdates, router.RouterOnRamp{
					DestChainSelector: dest,
					OnRamp:            common.HexToAddress("0x0"),
				})
			}
		}
		tx, err := routerC.ApplyRampUpdates(txOpts, onRampUpdates, removes, adds)
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

type SetOCR3OffRampConfig struct {
	HomeChainSel    uint64
	RemoteChainSels []uint64
	MCMS            *MCMSConfig
}

func (c SetOCR3OffRampConfig) Validate(e deployment.Environment) error {
	state, err := LoadOnchainState(e)
	if err != nil {
		return err
	}
	if _, ok := state.Chains[c.HomeChainSel]; !ok {
		return fmt.Errorf("home chain %d not found in onchain state", c.HomeChainSel)
	}
	for _, remote := range c.RemoteChainSels {
		chainState, ok := state.Chains[remote]
		if !ok {
			return fmt.Errorf("remote chain %d not found in onchain state", remote)
		}
		if err := commoncs.ValidateOwnership(e.GetContext(), c.MCMS != nil, e.Chains[remote].DeployerKey.From, chainState.Timelock.Address(), chainState.OffRamp); err != nil {
			return err
		}
	}
	return nil
}

// SetOCR3OffRamp will set the OCR3 offramp for the given chain.
// to the active configuration on CCIPHome. This
// is used to complete the candidate->active promotion cycle, it's
// run after the candidate is confirmed to be working correctly.
// Multichain is especially helpful for NOP rotations where we have
// to touch all the chain to change signers.
func SetOCR3OffRamp(e deployment.Environment, cfg SetOCR3OffRampConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	var batches []timelock.BatchChainOperation
	timelocks := make(map[uint64]common.Address)
	proposers := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for _, remote := range cfg.RemoteChainSels {
		donID, err := internal.DonIDForChain(
			state.Chains[cfg.HomeChainSel].CapabilityRegistry,
			state.Chains[cfg.HomeChainSel].CCIPHome,
			remote)
		args, err := internal.BuildSetOCR3ConfigArgs(donID, state.Chains[cfg.HomeChainSel].CCIPHome, remote)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		set, err := isOCR3ConfigSetOnOffRamp(e.Logger, e.Chains[remote], state.Chains[remote].OffRamp, args)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		if set {
			e.Logger.Infof("OCR3 config already set on offramp for chain %d", remote)
			continue
		}
		txOpts := e.Chains[remote].DeployerKey
		if cfg.MCMS != nil {
			txOpts = deployment.SimTransactOpts()
		}
		offRamp := state.Chains[remote].OffRamp
		tx, err := offRamp.SetOCR3Configs(txOpts, args)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		if cfg.MCMS == nil {
			if _, err := deployment.ConfirmIfNoError(e.Chains[remote], tx, err); err != nil {
				return deployment.ChangesetOutput{}, err
			}
		} else {
			batches = append(batches, timelock.BatchChainOperation{
				ChainIdentifier: mcms.ChainIdentifier(remote),
				Batch: []mcms.Operation{
					{
						To:    offRamp.Address(),
						Data:  tx.Data(),
						Value: big.NewInt(0),
					},
				},
			})
			timelocks[remote] = state.Chains[remote].Timelock.Address()
			proposers[remote] = state.Chains[remote].ProposerMcm
		}
	}
	if cfg.MCMS == nil {
		return deployment.ChangesetOutput{}, nil
	}
	p, err := proposalutils.BuildProposalFromBatches(
		timelocks,
		proposers,
		batches,
		"Update OCR3 config",
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infof("Proposing OCR3 config update for", cfg.RemoteChainSels)
	return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{
		*p,
	}}, nil
}

func isOCR3ConfigSetOnOffRamp(
	lggr logger.Logger,
	chain deployment.Chain,
	offRamp *offramp.OffRamp,
	offrampOCR3Configs []offramp.MultiOCR3BaseOCRConfigArgs,
) (bool, error) {
	mapOfframpOCR3Configs := make(map[cctypes.PluginType]offramp.MultiOCR3BaseOCRConfigArgs)
	for _, config := range offrampOCR3Configs {
		mapOfframpOCR3Configs[cctypes.PluginType(config.OcrPluginType)] = config
	}

	for _, pluginType := range []cctypes.PluginType{cctypes.PluginTypeCCIPCommit, cctypes.PluginTypeCCIPExec} {
		ocrConfig, err := offRamp.LatestConfigDetails(&bind.CallOpts{
			Context: context.Background(),
		}, uint8(pluginType))
		if err != nil {
			return false, fmt.Errorf("error fetching OCR3 config for plugin %s chain %s: %w", pluginType.String(), chain.String(), err)
		}
		lggr.Debugw("Fetched OCR3 Configs",
			"MultiOCR3BaseOCRConfig.F", ocrConfig.ConfigInfo.F,
			"MultiOCR3BaseOCRConfig.N", ocrConfig.ConfigInfo.N,
			"MultiOCR3BaseOCRConfig.IsSignatureVerificationEnabled", ocrConfig.ConfigInfo.IsSignatureVerificationEnabled,
			"Signers", ocrConfig.Signers,
			"Transmitters", ocrConfig.Transmitters,
			"configDigest", hex.EncodeToString(ocrConfig.ConfigInfo.ConfigDigest[:]),
			"chain", chain.String(),
		)
		// TODO: assertions to be done as part of full state
		// resprentation validation CCIP-3047
		if mapOfframpOCR3Configs[pluginType].ConfigDigest != ocrConfig.ConfigInfo.ConfigDigest {
			lggr.Infow("OCR3 config digest mismatch", "pluginType", pluginType.String())
			return false, nil
		}
		if mapOfframpOCR3Configs[pluginType].F != ocrConfig.ConfigInfo.F {
			lggr.Infow("OCR3 config F mismatch", "pluginType", pluginType.String())
			return false, nil
		}
		if mapOfframpOCR3Configs[pluginType].IsSignatureVerificationEnabled != ocrConfig.ConfigInfo.IsSignatureVerificationEnabled {
			lggr.Infow("OCR3 config signature verification mismatch", "pluginType", pluginType.String())
			return false, nil
		}
		if pluginType == cctypes.PluginTypeCCIPCommit {
			// only commit will set signers, exec doesn't need them.
			for i, signer := range mapOfframpOCR3Configs[pluginType].Signers {
				if !bytes.Equal(signer.Bytes(), ocrConfig.Signers[i].Bytes()) {
					lggr.Infow("OCR3 config signer mismatch", "pluginType", pluginType.String())
					return false, nil
				}
			}
		}
		for i, transmitter := range mapOfframpOCR3Configs[pluginType].Transmitters {
			if !bytes.Equal(transmitter.Bytes(), ocrConfig.Transmitters[i].Bytes()) {
				lggr.Infow("OCR3 config transmitter mismatch", "pluginType", pluginType.String())
				return false, nil
			}
		}
	}
	return true, nil
}
