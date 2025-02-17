package changeset

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commonState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/fee_quoter"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/onramp"
)

const (
	// https://github.com/smartcontractkit/chainlink/blob/1423e2581e8640d9e5cd06f745c6067bb2893af2/contracts/src/v0.8/ccip/libraries/Internal.sol#L275-L279
	/*
				```Solidity
					// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
					bytes4 public constant CHAIN_FAMILY_SELECTOR_EVM = 0x2812d52c;
					// bytes4(keccak256("CCIP ChainFamilySelector SVM"));
		  		bytes4 public constant CHAIN_FAMILY_SELECTOR_SVM = 0x1e10bdc4;
				```
	*/
	EVMFamilySelector = "2812d52c"
	SVMFamilySelector = "1e10bdc4"
)

var (
	_ deployment.ChangeSet[UpdateOnRampDestsConfig]          = UpdateOnRampsDestsChangeset
	_ deployment.ChangeSet[UpdateOnRampDynamicConfig]        = UpdateOnRampDynamicConfigChangeset
	_ deployment.ChangeSet[UpdateOnRampAllowListConfig]      = UpdateOnRampAllowListChangeset
	_ deployment.ChangeSet[WithdrawOnRampFeeTokensConfig]    = WithdrawOnRampFeeTokensChangeset
	_ deployment.ChangeSet[UpdateOffRampSourcesConfig]       = UpdateOffRampSourcesChangeset
	_ deployment.ChangeSet[UpdateRouterRampsConfig]          = UpdateRouterRampsChangeset
	_ deployment.ChangeSet[UpdateFeeQuoterDestsConfig]       = UpdateFeeQuoterDestsChangeset
	_ deployment.ChangeSet[SetOCR3OffRampConfig]             = SetOCR3OffRampChangeset
	_ deployment.ChangeSet[UpdateDynamicConfigOffRampConfig] = UpdateDynamicConfigOffRampChangeset
	_ deployment.ChangeSet[UpdateFeeQuoterPricesConfig]      = UpdateFeeQuoterPricesChangeset
	_ deployment.ChangeSet[UpdateNonceManagerConfig]         = UpdateNonceManagersChangeset
)

type UpdateNonceManagerConfig struct {
	UpdatesByChain map[uint64]NonceManagerUpdate // source -> dest -> update
	MCMS           *MCMSConfig
}

type NonceManagerUpdate struct {
	AddedAuthCallers   []common.Address
	RemovedAuthCallers []common.Address
	PreviousRampsArgs  []PreviousRampCfg
}

type PreviousRampCfg struct {
	RemoteChainSelector uint64
	OverrideExisting    bool
	// Set these only if the prevOnRamp or prevOffRamp addresses are not required to be in nonce manager.
	// If one of the onRamp or OffRamp is set with non-zero address and other is set with zero address,
	// it will not be possible to update the previous ramps later unless OverrideExisting is set to true.
	AllowEmptyOnRamp  bool // If true, the prevOnRamp address can be 0x0.
	AllowEmptyOffRamp bool // If true, the prevOffRamp address can be 0x0.
}

func (cfg UpdateNonceManagerConfig) Validate(e deployment.Environment) error {
	state, err := LoadOnchainState(e)
	if err != nil {
		return err
	}
	for sourceSel, update := range cfg.UpdatesByChain {
		sourceChainState, ok := state.Chains[sourceSel]
		if !ok {
			return fmt.Errorf("chain %d not found in onchain state", sourceSel)
		}
		if sourceChainState.NonceManager == nil {
			return fmt.Errorf("missing nonce manager for chain %d", sourceSel)
		}
		sourceChain, ok := e.Chains[sourceSel]
		if !ok {
			return fmt.Errorf("missing chain %d in environment", sourceSel)
		}
		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, sourceChain.DeployerKey.From, sourceChainState.Timelock.Address(), sourceChainState.OnRamp); err != nil {
			return fmt.Errorf("chain %s: %w", sourceChain.String(), err)
		}
		for _, prevRamp := range update.PreviousRampsArgs {
			if prevRamp.RemoteChainSelector == sourceSel {
				return errors.New("source and dest chain cannot be the same")
			}
			if _, ok := state.Chains[prevRamp.RemoteChainSelector]; !ok {
				return fmt.Errorf("dest chain %d not found in onchain state for chain %d", prevRamp.RemoteChainSelector, sourceSel)
			}
			// If one of the onRamp or OffRamp is set with non-zero address and other is set with zero address,
			// it will not be possible to update the previous ramps later.
			// Allow blank onRamp or offRamp only if AllowEmptyOnRamp or AllowEmptyOffRamp is set to true.
			// see https://github.com/smartcontractkit/chainlink/blob/develop/contracts/src/v0.8/ccip/NonceManager.sol#L139-L142
			if !prevRamp.AllowEmptyOnRamp {
				if prevOnRamp := state.Chains[sourceSel].EVM2EVMOnRamp; prevOnRamp == nil ||
					prevOnRamp[prevRamp.RemoteChainSelector] == nil ||
					prevOnRamp[prevRamp.RemoteChainSelector].Address() == (common.Address{}) {
					return fmt.Errorf("no previous onramp for source chain %d and dest chain %d, "+
						"If you want to set zero address for onRamp, set AllowEmptyOnRamp to true", sourceSel, prevRamp.RemoteChainSelector)
				}
			}
			if !prevRamp.AllowEmptyOffRamp {
				if prevOffRamp := state.Chains[sourceSel].EVM2EVMOffRamp; prevOffRamp == nil ||
					prevOffRamp[prevRamp.RemoteChainSelector] == nil ||
					prevOffRamp[prevRamp.RemoteChainSelector].Address() == (common.Address{}) {
					return fmt.Errorf("no previous offramp for source chain %d and dest chain %d"+
						"If you want to set zero address for offRamp, set AllowEmptyOffRamp to true", prevRamp.RemoteChainSelector, sourceSel)
				}
			}
		}
	}
	return nil
}

func UpdateNonceManagersChangeset(e deployment.Environment, cfg UpdateNonceManagerConfig) (deployment.ChangesetOutput, error) {
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
		if cfg.MCMS != nil {
			txOpts = deployment.SimTransactOpts()
		}
		nm := s.Chains[chainSel].NonceManager
		var authTx, prevRampsTx *types.Transaction
		if len(updates.AddedAuthCallers) > 0 || len(updates.RemovedAuthCallers) > 0 {
			authTx, err = nm.ApplyAuthorizedCallerUpdates(txOpts, nonce_manager.AuthorizedCallersAuthorizedCallerArgs{
				AddedCallers:   updates.AddedAuthCallers,
				RemovedCallers: updates.RemovedAuthCallers,
			})
			if cfg.MCMS == nil {
				if _, err := deployment.ConfirmIfNoErrorWithABI(e.Chains[chainSel], authTx, nonce_manager.NonceManagerABI, err); err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("error updating authorized callers for chain %s: %w",
						e.Chains[chainSel].String(), err)
				}
			} else {
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("error updating previous ramps for chain %s: %w", e.Chains[chainSel].String(), err)
				}
			}
		}
		if len(updates.PreviousRampsArgs) > 0 {
			previousRampsArgs := make([]nonce_manager.NonceManagerPreviousRampsArgs, 0)
			for _, prevRamp := range updates.PreviousRampsArgs {
				var onRamp, offRamp common.Address
				if !prevRamp.AllowEmptyOnRamp {
					onRamp = s.Chains[chainSel].EVM2EVMOnRamp[prevRamp.RemoteChainSelector].Address()
				}
				if !prevRamp.AllowEmptyOffRamp {
					offRamp = s.Chains[chainSel].EVM2EVMOffRamp[prevRamp.RemoteChainSelector].Address()
				}
				previousRampsArgs = append(previousRampsArgs, nonce_manager.NonceManagerPreviousRampsArgs{
					RemoteChainSelector:   prevRamp.RemoteChainSelector,
					OverrideExistingRamps: prevRamp.OverrideExisting,
					PrevRamps: nonce_manager.NonceManagerPreviousRamps{
						PrevOnRamp:  onRamp,
						PrevOffRamp: offRamp,
					},
				})
			}
			prevRampsTx, err = nm.ApplyPreviousRampsUpdates(txOpts, previousRampsArgs)
			if cfg.MCMS == nil {
				if _, err := deployment.ConfirmIfNoErrorWithABI(e.Chains[chainSel], prevRampsTx, nonce_manager.NonceManagerABI, err); err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("error updating previous ramps for chain %s: %w", e.Chains[chainSel].String(), err)
				}
			} else {
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("error updating previous ramps for chain %s: %w", e.Chains[chainSel].String(), err)
				}
			}
		}
		if cfg.MCMS != nil {
			ops := make([]mcms.Operation, 0)
			if authTx != nil {
				ops = append(ops, mcms.Operation{
					To:    nm.Address(),
					Data:  authTx.Data(),
					Value: big.NewInt(0),
				})
			}
			if prevRampsTx != nil {
				ops = append(ops, mcms.Operation{
					To:    nm.Address(),
					Data:  prevRampsTx.Data(),
					Value: big.NewInt(0),
				})
			}
			if len(ops) == 0 {
				return deployment.ChangesetOutput{}, errors.New("no operations to batch")
			}
			batches = append(batches, timelock.BatchChainOperation{
				ChainIdentifier: mcms.ChainIdentifier(chainSel),
				Batch:           ops,
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
		"Update nonce manager for previous ramps and authorized callers",
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{
		*p,
	}}, nil
}

type OnRampDestinationUpdate struct {
	IsEnabled        bool // If false, disables the destination by setting router to 0x0.
	TestRouter       bool // Flag for safety only allow specifying either router or testRouter.
	AllowListEnabled bool
}

type UpdateOnRampDestsConfig struct {
	// UpdatesByChain is a mapping of source -> dest -> update.
	UpdatesByChain map[uint64]map[uint64]OnRampDestinationUpdate

	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMS *MCMSConfig
}

func (cfg UpdateOnRampDestsConfig) Validate(e deployment.Environment) error {
	state, err := LoadOnchainState(e)
	if err != nil {
		return err
	}
	supportedChains := state.SupportedChains()
	for chainSel, updates := range cfg.UpdatesByChain {
		if err := ValidateChain(e, state, chainSel, cfg.MCMS != nil); err != nil {
			return err
		}
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
		sc, err := chainState.OnRamp.GetStaticConfig(&bind.CallOpts{Context: e.GetContext()})
		if err != nil {
			return fmt.Errorf("failed to get onramp static config %s: %w", chainState.OnRamp.Address(), err)
		}
		for destination := range updates {
			// Destination cannot be an unknown destination.
			if _, ok := supportedChains[destination]; !ok {
				return fmt.Errorf("destination chain %d is not a supported %s", destination, chainState.OnRamp.Address())
			}
			if destination == sc.ChainSelector {
				return errors.New("cannot update onramp destination to the same chain")
			}
		}
	}
	return nil
}

// UpdateOnRampsDestsChangeset updates the onramp destinations for each onramp
// in the chains specified. Multichain support is important - consider when we add a new chain
// and need to update the onramp destinations for all chains to support the new chain.
func UpdateOnRampsDestsChangeset(e deployment.Environment, cfg UpdateOnRampDestsConfig) (deployment.ChangesetOutput, error) {
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
		if cfg.MCMS == nil {
			if _, err := deployment.ConfirmIfNoErrorWithABI(e.Chains[chainSel], tx, onramp.OnRampABI, err); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("error updating onramp destinations for chain %s: %w", e.Chains[chainSel].String(), err)
			}
		} else {
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
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

type OnRampDynamicConfigUpdate struct {
	MessageInterceptor common.Address
	FeeAggregator      common.Address
	AllowlistAdmin     common.Address
}

type UpdateOnRampDynamicConfig struct {
	// UpdatesByChain is a mapping of source -> update.
	UpdatesByChain map[uint64]OnRampDynamicConfigUpdate
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMS *MCMSConfig
}

func (cfg UpdateOnRampDynamicConfig) Validate(e deployment.Environment, state CCIPOnChainState) error {
	for chainSel, config := range cfg.UpdatesByChain {
		if err := ValidateChain(e, state, chainSel, cfg.MCMS != nil); err != nil {
			return err
		}
		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.Chains[chainSel].DeployerKey.From, state.Chains[chainSel].Timelock.Address(), state.Chains[chainSel].OnRamp); err != nil {
			return err
		}
		if state.Chains[chainSel].FeeQuoter == nil {
			return fmt.Errorf("FeeQuoter is not on state of chain %d", chainSel)
		}
		if config.FeeAggregator == (common.Address{}) {
			return fmt.Errorf("FeeAggregator is not specified for chain %d", chainSel)
		}
	}
	return nil
}

func UpdateOnRampDynamicConfigChangeset(e deployment.Environment, cfg UpdateOnRampDynamicConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	var batches []timelock.BatchChainOperation
	timelocks := make(map[uint64]common.Address)
	proposers := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for chainSel, update := range cfg.UpdatesByChain {
		txOps := e.Chains[chainSel].DeployerKey
		if cfg.MCMS != nil {
			txOps = deployment.SimTransactOpts()
		}
		onRamp := state.Chains[chainSel].OnRamp
		dynamicConfig, err := onRamp.GetDynamicConfig(nil)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		// Do not update dynamic config if it is already in desired state
		if dynamicConfig.FeeQuoter == state.Chains[chainSel].FeeQuoter.Address() &&
			dynamicConfig.MessageInterceptor == update.MessageInterceptor &&
			dynamicConfig.FeeAggregator == update.FeeAggregator &&
			dynamicConfig.AllowlistAdmin == update.AllowlistAdmin {
			continue
		}
		tx, err := onRamp.SetDynamicConfig(txOps, onramp.OnRampDynamicConfig{
			FeeQuoter:              state.Chains[chainSel].FeeQuoter.Address(),
			ReentrancyGuardEntered: false,
			MessageInterceptor:     update.MessageInterceptor,
			FeeAggregator:          update.FeeAggregator,
			AllowlistAdmin:         update.AllowlistAdmin,
		})

		if cfg.MCMS == nil {
			if _, err := deployment.ConfirmIfNoErrorWithABI(e.Chains[chainSel], tx, onramp.OnRampABI, err); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("error updating onramp dynamic config for chain %s: %w", e.Chains[chainSel].String(), err)
			}
		} else {
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
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
			timelocks[chainSel] = state.Chains[chainSel].Timelock.Address()
			proposers[chainSel] = state.Chains[chainSel].ProposerMcm
		}
	}
	if cfg.MCMS == nil {
		return deployment.ChangesetOutput{}, nil
	}
	proposal, err := proposalutils.BuildProposalFromBatches(
		timelocks, proposers, batches,
		"update onramp dynamic config",
		cfg.MCMS.MinDelay)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*proposal},
	}, nil
}

type OnRampAllowListUpdate struct {
	AllowListEnabled          bool
	AddedAllowlistedSenders   []common.Address
	RemovedAllowlistedSenders []common.Address
}

type UpdateOnRampAllowListConfig struct {
	// UpdatesByChain is a mapping of source -> dest -> update.
	UpdatesByChain map[uint64]map[uint64]OnRampAllowListUpdate
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMS *MCMSConfig
}

func (cfg UpdateOnRampAllowListConfig) Validate(env deployment.Environment) error {
	state, err := LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for srcSel, updates := range cfg.UpdatesByChain {
		if err := ValidateChain(env, state, srcSel, cfg.MCMS != nil); err != nil {
			return err
		}
		onRamp := state.Chains[srcSel].OnRamp
		if onRamp == nil {
			return fmt.Errorf("missing onRamp on %d", srcSel)
		}
		config, err := onRamp.GetDynamicConfig(nil)
		if err != nil {
			return err
		}
		owner, err := onRamp.Owner(nil)
		if err != nil {
			return fmt.Errorf("failed to get owner: %w", err)
		}
		var signer common.Address
		if cfg.MCMS == nil {
			signer = env.Chains[srcSel].DeployerKey.From
			if signer != config.AllowlistAdmin && signer != owner {
				return fmt.Errorf("deployer key is not onramp's %s owner nor allowlist admin", onRamp.Address())
			}
		} else {
			signer = state.Chains[srcSel].Timelock.Address()
			if signer != config.AllowlistAdmin && signer != owner {
				return fmt.Errorf("timelock is not onramp's %s owner nor allowlist admin", onRamp.Address())
			}
		}
		for destSel, update := range updates {
			if err := ValidateChain(env, state, srcSel, false); err != nil {
				return err
			}
			if len(update.AddedAllowlistedSenders) > 0 && !update.AllowListEnabled {
				return fmt.Errorf("can't allowlist senders with disabled allowlist for src=%d, dest=%d", srcSel, destSel)
			}
			for _, sender := range update.AddedAllowlistedSenders {
				if sender == (common.Address{}) {
					return fmt.Errorf("can't allowlist 0-address sender for src=%d, dest=%d", srcSel, destSel)
				}
			}
		}
	}
	return nil
}

func UpdateOnRampAllowListChangeset(e deployment.Environment, cfg UpdateOnRampAllowListConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	onchain, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	var batches []timelock.BatchChainOperation
	timelocks := make(map[uint64]common.Address)
	proposers := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for srcSel, updates := range cfg.UpdatesByChain {
		txOps := e.Chains[srcSel].DeployerKey
		if cfg.MCMS != nil {
			txOps = deployment.SimTransactOpts()
		}
		onRamp := onchain.Chains[srcSel].OnRamp
		args := make([]onramp.OnRampAllowlistConfigArgs, len(updates))
		for destSel, update := range updates {
			allowedSendersResp, err := onRamp.GetAllowedSendersList(nil, destSel)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			if allowedSendersResp.IsEnabled == update.AllowListEnabled {
				desiredState := make(map[common.Address]bool)
				for _, address := range update.AddedAllowlistedSenders {
					desiredState[address] = true
				}
				for _, address := range update.RemovedAllowlistedSenders {
					desiredState[address] = false
				}
				needUpdate := false
				for _, allowedSender := range allowedSendersResp.ConfiguredAddresses {
					if !desiredState[allowedSender] {
						needUpdate = true
					}
				}
				if !needUpdate {
					continue
				}
			}
			args = append(args, onramp.OnRampAllowlistConfigArgs{
				DestChainSelector:         destSel,
				AllowlistEnabled:          update.AllowListEnabled,
				AddedAllowlistedSenders:   update.AddedAllowlistedSenders,
				RemovedAllowlistedSenders: update.RemovedAllowlistedSenders,
			})
		}
		if len(args) == 0 {
			continue
		}
		tx, err := onRamp.ApplyAllowlistUpdates(txOps, args)
		if cfg.MCMS == nil {
			if _, err := deployment.ConfirmIfNoErrorWithABI(e.Chains[srcSel], tx, onramp.OnRampABI, err); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("error updating allowlist for chain %d: %w", srcSel, err)
			}
		} else {
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			batches = append(batches, timelock.BatchChainOperation{
				ChainIdentifier: mcms.ChainIdentifier(srcSel),
				Batch: []mcms.Operation{
					{
						To:    onRamp.Address(),
						Data:  tx.Data(),
						Value: big.NewInt(0),
					},
				},
			})
			timelocks[srcSel] = onchain.Chains[srcSel].Timelock.Address()
			proposers[srcSel] = onchain.Chains[srcSel].ProposerMcm
		}
	}
	if cfg.MCMS == nil {
		return deployment.ChangesetOutput{}, nil
	}
	proposal, err := proposalutils.BuildProposalFromBatches(
		timelocks,
		proposers,
		batches,
		"update onramp allowlist",
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*proposal},
	}, nil
}

type WithdrawOnRampFeeTokensConfig struct {
	FeeTokensByChain map[uint64][]common.Address
	MCMS             *MCMSConfig
}

func (cfg WithdrawOnRampFeeTokensConfig) Validate(e deployment.Environment, state CCIPOnChainState) error {
	for chainSel, feeTokens := range cfg.FeeTokensByChain {
		if err := ValidateChain(e, state, chainSel, cfg.MCMS != nil); err != nil {
			return err
		}
		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.Chains[chainSel].DeployerKey.From, state.Chains[chainSel].Timelock.Address(), state.Chains[chainSel].OnRamp); err != nil {
			return err
		}
		feeQuoter := state.Chains[chainSel].FeeQuoter
		if feeQuoter == nil {
			return fmt.Errorf("no fee quoter for chain %d", chainSel)
		}
		onchainFeeTokens, err := feeQuoter.GetFeeTokens(nil)
		if len(onchainFeeTokens) == 0 {
			return fmt.Errorf("no fee tokens configured on fee quoter %s for chain %d", feeQuoter.Address().Hex(), chainSel)
		}
		if err != nil {
			return err
		}
		for _, feeToken := range feeTokens {
			found := false
			for _, onchainFeeToken := range onchainFeeTokens {
				if onchainFeeToken == feeToken {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("unknown fee token address=%s on chain=%d", feeToken.Hex(), chainSel)
			}
		}
	}
	return nil
}

func WithdrawOnRampFeeTokensChangeset(e deployment.Environment, cfg WithdrawOnRampFeeTokensConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	var batches []timelock.BatchChainOperation
	timelocks := make(map[uint64]common.Address)
	proposers := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for chainSel, feeTokens := range cfg.FeeTokensByChain {
		txOps := e.Chains[chainSel].DeployerKey
		onRamp := state.Chains[chainSel].OnRamp
		tx, err := onRamp.WithdrawFeeTokens(txOps, feeTokens)
		if cfg.MCMS == nil {
			if _, err := deployment.ConfirmIfNoErrorWithABI(e.Chains[chainSel], tx, onramp.OnRampABI, err); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("error withdrawing fee tokens for chain %s: %w", e.Chains[chainSel].String(), err)
			}
		} else {
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
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
			timelocks[chainSel] = state.Chains[chainSel].Timelock.Address()
			proposers[chainSel] = state.Chains[chainSel].ProposerMcm
		}
	}
	if cfg.MCMS == nil {
		return deployment.ChangesetOutput{}, nil
	}
	proposal, err := proposalutils.BuildProposalFromBatches(
		timelocks,
		proposers,
		batches,
		"withdraw onramp fee tokens",
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*proposal},
	}, nil
}

type UpdateFeeQuoterPricesConfig struct {
	PricesByChain map[uint64]FeeQuoterPriceUpdatePerSource // source -> PriceDetails
	MCMS          *MCMSConfig
}

type FeeQuoterPriceUpdatePerSource struct {
	TokenPrices map[common.Address]*big.Int // token address -> price
	GasPrices   map[uint64]*big.Int         // dest chain -> gas price
}

func (cfg UpdateFeeQuoterPricesConfig) Validate(e deployment.Environment) error {
	state, err := LoadOnchainState(e)
	if err != nil {
		return err
	}
	for chainSel, initialPrice := range cfg.PricesByChain {
		if err := deployment.IsValidChainSelector(chainSel); err != nil {
			return fmt.Errorf("invalid chain selector: %w", err)
		}
		chainState, ok := state.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain %d not found in onchain state", chainSel)
		}
		fq := chainState.FeeQuoter
		if fq == nil {
			return fmt.Errorf("missing fee quoter for chain %d", chainSel)
		}
		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.Chains[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.FeeQuoter); err != nil {
			return err
		}
		// check that whether price updaters are set
		authCallers, err := fq.GetAllAuthorizedCallers(&bind.CallOpts{Context: e.GetContext()})
		if err != nil {
			return fmt.Errorf("failed to get authorized callers for chain %d: %w", chainSel, err)
		}
		if len(authCallers) == 0 {
			return fmt.Errorf("no authorized callers for chain %d", chainSel)
		}
		expectedAuthCaller := e.Chains[chainSel].DeployerKey.From
		if cfg.MCMS != nil {
			expectedAuthCaller = chainState.Timelock.Address()
		}
		foundCaller := false
		for _, authCaller := range authCallers {
			if authCaller.Cmp(expectedAuthCaller) == 0 {
				foundCaller = true
			}
		}
		if !foundCaller {
			return fmt.Errorf("expected authorized caller %s not found for chain %d", expectedAuthCaller.String(), chainSel)
		}
		for token, price := range initialPrice.TokenPrices {
			if price == nil {
				return fmt.Errorf("token price for chain %d is nil", chainSel)
			}
			if token == (common.Address{}) {
				return fmt.Errorf("token address for chain %d is empty", chainSel)
			}
			contains, err := deployment.AddressBookContains(e.ExistingAddresses, chainSel, token.String())
			if err != nil {
				return fmt.Errorf("error checking address book for token %s: %w", token.String(), err)
			}
			if !contains {
				return fmt.Errorf("token %s not found in address book for chain %d", token.String(), chainSel)
			}
		}
		for dest, price := range initialPrice.GasPrices {
			if chainSel == dest {
				return errors.New("source and dest chain cannot be the same")
			}
			if err := deployment.IsValidChainSelector(dest); err != nil {
				return fmt.Errorf("invalid dest chain selector: %w", err)
			}
			if price == nil {
				return fmt.Errorf("gas price for chain %d is nil", chainSel)
			}
			if _, ok := state.SupportedChains()[dest]; !ok {
				return fmt.Errorf("dest chain %d not found in onchain state for chain %d", dest, chainSel)
			}
		}
	}

	return nil
}

func UpdateFeeQuoterPricesChangeset(e deployment.Environment, cfg UpdateFeeQuoterPricesConfig) (deployment.ChangesetOutput, error) {
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
	for chainSel, initialPrice := range cfg.PricesByChain {
		txOpts := e.Chains[chainSel].DeployerKey
		if cfg.MCMS != nil {
			txOpts = deployment.SimTransactOpts()
		}
		fq := s.Chains[chainSel].FeeQuoter
		var tokenPricesArgs []fee_quoter.InternalTokenPriceUpdate
		for token, price := range initialPrice.TokenPrices {
			tokenPricesArgs = append(tokenPricesArgs, fee_quoter.InternalTokenPriceUpdate{
				SourceToken: token,
				UsdPerToken: price,
			})
		}
		var gasPricesArgs []fee_quoter.InternalGasPriceUpdate
		for dest, price := range initialPrice.GasPrices {
			gasPricesArgs = append(gasPricesArgs, fee_quoter.InternalGasPriceUpdate{
				DestChainSelector: dest,
				UsdPerUnitGas:     price,
			})
		}
		tx, err := fq.UpdatePrices(txOpts, fee_quoter.InternalPriceUpdates{
			TokenPriceUpdates: tokenPricesArgs,
			GasPriceUpdates:   gasPricesArgs,
		})
		if cfg.MCMS == nil {
			if _, err := deployment.ConfirmIfNoErrorWithABI(e.Chains[chainSel], tx, fee_quoter.FeeQuoterABI, err); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("error confirming transaction for chain %s: %w", e.Chains[chainSel].String(), err)
			}
		} else {
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("error updating prices for chain %s: %w", e.Chains[chainSel].String(), err)
			}
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
		"Update fq prices",
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
	// UpdatesByChain is a mapping from source -> dest -> config update.
	UpdatesByChain map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
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
				return errors.New("source and destination chain cannot be the same")
			}
		}
	}
	return nil
}

func UpdateFeeQuoterDestsChangeset(e deployment.Environment, cfg UpdateFeeQuoterDestsConfig) (deployment.ChangesetOutput, error) {
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
			if _, err := deployment.ConfirmIfNoErrorWithABI(e.Chains[chainSel], tx, fee_quoter.FeeQuoterABI, err); err != nil {
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

type OffRampSourceUpdate struct {
	IsEnabled  bool // If false, disables the source by setting router to 0x0.
	TestRouter bool // Flag for safety only allow specifying either router or testRouter.
	// IsRMNVerificationDisabled is a flag to disable RMN verification for this source chain.
	IsRMNVerificationDisabled bool
}

type UpdateOffRampSourcesConfig struct {
	// UpdatesByChain is a mapping from dest chain -> source chain -> source chain
	// update on the dest chain offramp.
	UpdatesByChain map[uint64]map[uint64]OffRampSourceUpdate
	MCMS           *MCMSConfig
}

func (cfg UpdateOffRampSourcesConfig) Validate(e deployment.Environment, state CCIPOnChainState) error {
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

// UpdateOffRampSourcesChangeset updates the offramp sources for each offramp.
func UpdateOffRampSourcesChangeset(e deployment.Environment, cfg UpdateOffRampSourcesConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
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
		offRamp := state.Chains[chainSel].OffRamp
		var args []offramp.OffRampSourceChainConfigArgs
		for source, update := range updates {
			router := common.HexToAddress("0x0")
			if update.TestRouter {
				router = state.Chains[chainSel].TestRouter.Address()
			} else {
				router = state.Chains[chainSel].Router.Address()
			}
			onRamp := state.Chains[source].OnRamp
			args = append(args, offramp.OffRampSourceChainConfigArgs{
				SourceChainSelector: source,
				Router:              router,
				IsEnabled:           update.IsEnabled,
				// TODO: how would this work when the onRamp is nonEVM?
				OnRamp:                    common.LeftPadBytes(onRamp.Address().Bytes(), 32),
				IsRMNVerificationDisabled: update.IsRMNVerificationDisabled,
			})
		}
		tx, err := offRamp.ApplySourceChainConfigUpdates(txOpts, args)

		if cfg.MCMS == nil {
			if _, err := deployment.ConfirmIfNoErrorWithABI(e.Chains[chainSel], tx, offramp.OffRampABI, err); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("error applying source chain config updates for chain %d: %w", chainSel, err)
			}
		} else {
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
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
			timelocks[chainSel] = state.Chains[chainSel].Timelock.Address()
			proposers[chainSel] = state.Chains[chainSel].ProposerMcm
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

type RouterUpdates struct {
	OffRampUpdates map[uint64]bool
	OnRampUpdates  map[uint64]bool
}

type UpdateRouterRampsConfig struct {
	// TestRouter means the updates will be applied to the test router
	// on all chains. Disallow mixing test router/non-test router per chain for simplicity.
	TestRouter     bool
	UpdatesByChain map[uint64]RouterUpdates
	MCMS           *MCMSConfig
}

func (cfg UpdateRouterRampsConfig) Validate(e deployment.Environment, state CCIPOnChainState) error {
	supportedChains := state.SupportedChains()
	for chainSel, update := range cfg.UpdatesByChain {
		if err := ValidateChain(e, state, chainSel, cfg.MCMS != nil); err != nil {
			return err
		}
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
		if cfg.TestRouter {
			if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.Chains[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.TestRouter); err != nil {
				return err
			}
		} else {
			if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.Chains[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.Router); err != nil {
				return err
			}
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
			if err := state.ValidateOffRamp(destination); err != nil {
				return err
			}
		}
	}
	return nil
}

// UpdateRouterRampsChangeset updates the on/offramps
// in either the router or test router for a series of chains. Use cases include:
// - Ramp upgrade. After deploying new ramps you can enable them on the test router and
// ensure it works e2e. Then enable the ramps on the real router.
// - New chain support. When adding a new chain, you can enable the new destination
// on all chains to support the new chain through the test router first. Once tested,
// Enable the new destination on the real router.
func UpdateRouterRampsChangeset(e deployment.Environment, cfg UpdateRouterRampsConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
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
		routerC := state.Chains[chainSel].Router
		if cfg.TestRouter {
			routerC = state.Chains[chainSel].TestRouter
		}
		// Note if we add distinct offramps per source to the state,
		// we'll need to add support here for looking them up.
		// For now its simple, all sources use the same offramp.
		offRamp := state.Chains[chainSel].OffRamp
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
		onRamp := state.Chains[chainSel].OnRamp
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
		if cfg.MCMS == nil {
			if _, err := deployment.ConfirmIfNoErrorWithABI(e.Chains[chainSel], tx, router.RouterABI, err); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("error applying ramp updates for chain %d: %w", chainSel, err)
			}
		} else {
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
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
			timelocks[chainSel] = state.Chains[chainSel].Timelock.Address()
			proposers[chainSel] = state.Chains[chainSel].ProposerMcm
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
	HomeChainSel       uint64
	RemoteChainSels    []uint64
	CCIPHomeConfigType globals.ConfigType
	MCMS               *MCMSConfig
}

func (c SetOCR3OffRampConfig) Validate(e deployment.Environment, state CCIPOnChainState) error {
	if _, ok := state.Chains[c.HomeChainSel]; !ok {
		return fmt.Errorf("home chain %d not found in onchain state", c.HomeChainSel)
	}
	if err := ValidateChain(e, state, c.HomeChainSel, c.MCMS != nil); err != nil {
		return err
	}
	if c.CCIPHomeConfigType != globals.ConfigTypeActive &&
		c.CCIPHomeConfigType != globals.ConfigTypeCandidate {
		return fmt.Errorf("invalid CCIPHomeConfigType should be either %s or %s", globals.ConfigTypeActive, globals.ConfigTypeCandidate)
	}
	for _, remote := range c.RemoteChainSels {
		if err := c.validateRemoteChain(&e, &state, remote); err != nil {
			return err
		}
	}
	return nil
}

func (c SetOCR3OffRampConfig) validateRemoteChain(e *deployment.Environment, state *CCIPOnChainState, chainSelector uint64) error {
	family, err := chain_selectors.GetSelectorFamily(chainSelector)
	if err != nil {
		return err
	}
	switch family {
	case chain_selectors.FamilySolana:
		chain, ok := e.SolChains[chainSelector]
		if !ok {
			return fmt.Errorf("chain %d not found in environment", chainSelector)
		}
		chainState, ok := state.SolChains[chainSelector]
		if !ok {
			return fmt.Errorf("remote chain %d not found in onchain state", chainSelector)
		}
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return err
		}
		mcmState, err := commonState.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
		if err != nil {
			return fmt.Errorf("error loading MCMS state for chain %d: %w", chainSelector, err)
		}
		if err := commoncs.ValidateOwnershipSolana(e.GetContext(), c.MCMS != nil, e.SolChains[chainSelector].DeployerKey.PublicKey(), mcmState.TimelockProgram, mcmState.TimelockSeed, chainState.Router); err != nil {
			return err
		}
	case chain_selectors.FamilyEVM:
		chainState, ok := state.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("remote chain %d not found in onchain state", chainSelector)
		}
		if err := commoncs.ValidateOwnership(e.GetContext(), c.MCMS != nil, e.Chains[chainSelector].DeployerKey.From, chainState.Timelock.Address(), chainState.OffRamp); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported chain family %s", family)
	}
	return nil
}

// SetOCR3OffRampChangeset will set the OCR3 offramp for the given chain.
// to the active configuration on CCIPHome. This
// is used to complete the candidate->active promotion cycle, it's
// run after the candidate is confirmed to be working correctly.
// Multichain is especially helpful for NOP rotations where we have
// to touch all the chain to change signers.
func SetOCR3OffRampChangeset(e deployment.Environment, cfg SetOCR3OffRampConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
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
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		args, err := internal.BuildSetOCR3ConfigArgs(
			donID, state.Chains[cfg.HomeChainSel].CCIPHome, remote, cfg.CCIPHomeConfigType)
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
		if cfg.MCMS == nil {
			if _, err := deployment.ConfirmIfNoErrorWithABI(e.Chains[remote], tx, offramp.OffRampABI, err); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("error setting OCR3 config for chain %d: %w", remote, err)
			}
		} else {
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
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
	e.Logger.Info("Proposing OCR3 config update for", cfg.RemoteChainSels)
	return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{
		*p,
	}}, nil
}

type UpdateDynamicConfigOffRampConfig struct {
	Updates map[uint64]OffRampParams
	MCMS    *MCMSConfig
}

func (cfg UpdateDynamicConfigOffRampConfig) Validate(e deployment.Environment) error {
	state, err := LoadOnchainState(e)
	if err != nil {
		return err
	}
	for chainSel, params := range cfg.Updates {
		if deployment.IsValidChainSelector(chainSel) != nil {
			return fmt.Errorf("invalid chain selector %d", chainSel)
		}
		if _, ok := state.Chains[chainSel]; !ok {
			return fmt.Errorf("chain %d not found in onchain state", chainSel)
		}
		if state.Chains[chainSel].OffRamp == nil {
			return fmt.Errorf("missing offramp for chain %d", chainSel)
		}
		if state.Chains[chainSel].FeeQuoter == nil {
			return fmt.Errorf("missing fee quoter for chain %d", chainSel)
		}
		if state.Chains[chainSel].Timelock == nil {
			return fmt.Errorf("missing timelock for chain %d", chainSel)
		}
		if params.GasForCallExactCheck > 0 {
			e.Logger.Infow(
				"GasForCallExactCheck is set, please note it's a static config and will be ignored for this changeset",
				"chain", chainSel, "gas", params.GasForCallExactCheck)
		}
		if err := commoncs.ValidateOwnership(
			e.GetContext(),
			cfg.MCMS != nil,
			e.Chains[chainSel].DeployerKey.From,
			state.Chains[chainSel].Timelock.Address(),
			state.Chains[chainSel].OffRamp,
		); err != nil {
			return err
		}
		if err := params.Validate(true); err != nil {
			return fmt.Errorf("chain %d: %w", chainSel, err)
		}
	}
	return nil
}

func UpdateDynamicConfigOffRampChangeset(e deployment.Environment, cfg UpdateDynamicConfigOffRampConfig) (deployment.ChangesetOutput, error) {
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
	for chainSel, params := range cfg.Updates {
		chain := e.Chains[chainSel]
		txOpts := e.Chains[chainSel].DeployerKey
		if cfg.MCMS != nil {
			txOpts = deployment.SimTransactOpts()
		}
		offRamp := state.Chains[chainSel].OffRamp
		dCfg := offramp.OffRampDynamicConfig{
			FeeQuoter:                               state.Chains[chainSel].FeeQuoter.Address(),
			PermissionLessExecutionThresholdSeconds: params.PermissionLessExecutionThresholdSeconds,
			MessageInterceptor:                      params.MessageInterceptor,
		}
		tx, err := offRamp.SetDynamicConfig(txOpts, dCfg)
		if cfg.MCMS == nil {
			if _, err := deployment.ConfirmIfNoErrorWithABI(e.Chains[chainSel], tx, offramp.OffRampABI, err); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("error updating offramp dynamic config for chain %d: %w", chainSel, err)
			}
			e.Logger.Infow("Updated offramp dynamic config", "chain", chain.String(), "config", dCfg)
		} else {
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
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
			timelocks[chainSel] = state.Chains[chainSel].Timelock.Address()
			proposers[chainSel] = state.Chains[chainSel].ProposerMcm
		}
	}
	if cfg.MCMS == nil {
		return deployment.ChangesetOutput{}, nil
	}
	p, err := proposalutils.BuildProposalFromBatches(
		timelocks,
		proposers,
		batches,
		"Update offramp dynamic config",
		cfg.MCMS.MinDelay,
	)

	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infow("Proposing offramp dynamic config update", "config", cfg.Updates)
	return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{
		*p,
	}}, nil
}

func isOCR3ConfigSetOnOffRamp(
	lggr logger.Logger,
	chain deployment.Chain,
	offRamp offramp.OffRampInterface,
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

// DefaultFeeQuoterDestChainConfig returns the default FeeQuoterDestChainConfig
// with the config enabled/disabled based on the configEnabled flag.
func DefaultFeeQuoterDestChainConfig(configEnabled bool, destChainSelector ...uint64) fee_quoter.FeeQuoterDestChainConfig {
	familySelector, _ := hex.DecodeString(EVMFamilySelector) // evm
	if len(destChainSelector) > 0 {
		destFamily, _ := chain_selectors.GetSelectorFamily(destChainSelector[0])
		if destFamily == chain_selectors.FamilySolana {
			familySelector, _ = hex.DecodeString(SVMFamilySelector) // solana
		}
	}
	return fee_quoter.FeeQuoterDestChainConfig{
		IsEnabled:                         configEnabled,
		MaxNumberOfTokensPerMsg:           10,
		MaxDataBytes:                      256,
		MaxPerMsgGasLimit:                 3_000_000,
		DestGasOverhead:                   ccipevm.DestGasOverhead,
		DefaultTokenFeeUSDCents:           1,
		DestGasPerPayloadByteBase:         ccipevm.CalldataGasPerByteBase,
		DestGasPerPayloadByteHigh:         ccipevm.CalldataGasPerByteHigh,
		DestGasPerPayloadByteThreshold:    ccipevm.CalldataGasPerByteThreshold,
		DestDataAvailabilityOverheadGas:   100,
		DestGasPerDataAvailabilityByte:    100,
		DestDataAvailabilityMultiplierBps: 1,
		DefaultTokenDestGasOverhead:       125_000,
		DefaultTxGasLimit:                 200_000,
		GasMultiplierWeiPerEth:            11e17, // Gas multiplier in wei per eth is scaled by 1e18, so 11e17 is 1.1 = 110%
		NetworkFeeUSDCents:                1,
		ChainFamilySelector:               [4]byte(familySelector),
	}
}
