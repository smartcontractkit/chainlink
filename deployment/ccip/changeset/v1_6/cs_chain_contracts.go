package v1_6

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
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"

	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	ccipseqs "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
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
	EVMFamilySelector   = "2812d52c"
	SVMFamilySelector   = "1e10bdc4"
	AptosFamilySelector = "ac77ffec"
)

var (
	_ cldf.ChangeSet[UpdateOnRampDestsConfig]                  = UpdateOnRampsDestsChangeset
	_ cldf.ChangeSet[UpdateOnRampDynamicConfig]                = UpdateOnRampDynamicConfigChangeset
	_ cldf.ChangeSet[UpdateOnRampAllowListConfig]              = UpdateOnRampAllowListChangeset
	_ cldf.ChangeSet[WithdrawOnRampFeeTokensConfig]            = WithdrawOnRampFeeTokensChangeset
	_ cldf.ChangeSet[UpdateOffRampSourcesConfig]               = UpdateOffRampSourcesChangeset
	_ cldf.ChangeSet[UpdateRouterRampsConfig]                  = UpdateRouterRampsChangeset
	_ cldf.ChangeSet[UpdateFeeQuoterDestsConfig]               = UpdateFeeQuoterDestsChangeset
	_ cldf.ChangeSet[SetOCR3OffRampConfig]                     = SetOCR3OffRampChangeset
	_ cldf.ChangeSet[UpdateDynamicConfigOffRampConfig]         = UpdateDynamicConfigOffRampChangeset
	_ cldf.ChangeSet[UpdateFeeQuoterPricesConfig]              = UpdateFeeQuoterPricesChangeset
	_ cldf.ChangeSet[UpdateNonceManagerConfig]                 = UpdateNonceManagersChangeset
	_ cldf.ChangeSet[ApplyFeeTokensUpdatesConfig]              = ApplyFeeTokensUpdatesFeeQuoterChangeset
	_ cldf.ChangeSet[UpdateTokenPriceFeedsConfig]              = UpdateTokenPriceFeedsFeeQuoterChangeset
	_ cldf.ChangeSet[PremiumMultiplierWeiPerEthUpdatesConfig]  = ApplyPremiumMultiplierWeiPerEthUpdatesFeeQuoterChangeset
	_ cldf.ChangeSet[ApplyTokenTransferFeeConfigUpdatesConfig] = ApplyTokenTransferFeeConfigUpdatesFeeQuoterChangeset
)

type UpdateNonceManagerConfig struct {
	UpdatesByChain map[uint64]NonceManagerUpdate // source -> dest -> update
	MCMS           *proposalutils.TimelockConfig
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

func (cfg UpdateNonceManagerConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
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
		sourceChain, ok := e.BlockChains.EVMChains()[sourceSel]
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

func UpdateNonceManagersChangeset(e cldf.Environment, cfg UpdateNonceManagerConfig) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	s, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	batches := []mcmstypes.BatchOperation{}
	timelocks := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)

	for chainSel, updates := range cfg.UpdatesByChain {
		txOpts := e.BlockChains.EVMChains()[chainSel].DeployerKey
		if cfg.MCMS != nil {
			txOpts = cldf.SimTransactOpts()
		}
		nm := s.Chains[chainSel].NonceManager
		var authTx, prevRampsTx *types.Transaction
		if len(updates.AddedAuthCallers) > 0 || len(updates.RemovedAuthCallers) > 0 {
			authTx, err = nm.ApplyAuthorizedCallerUpdates(txOpts, nonce_manager.AuthorizedCallersAuthorizedCallerArgs{
				AddedCallers:   updates.AddedAuthCallers,
				RemovedCallers: updates.RemovedAuthCallers,
			})
			if cfg.MCMS == nil {
				if _, err := cldf.ConfirmIfNoErrorWithABI(e.BlockChains.EVMChains()[chainSel], authTx, nonce_manager.NonceManagerABI, err); err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("error updating authorized callers for chain %s: %w",
						e.BlockChains.EVMChains()[chainSel].String(), err)
				}
			} else {
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("error updating previous ramps for chain %s: %w", e.BlockChains.EVMChains()[chainSel].String(), err)
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
				if _, err := cldf.ConfirmIfNoErrorWithABI(e.BlockChains.EVMChains()[chainSel], prevRampsTx, nonce_manager.NonceManagerABI, err); err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("error updating previous ramps for chain %s: %w", e.BlockChains.EVMChains()[chainSel].String(), err)
				}
			} else {
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("error updating previous ramps for chain %s: %w", e.BlockChains.EVMChains()[chainSel].String(), err)
				}
			}
		}
		if cfg.MCMS != nil {
			mcmsTransactions := make([]mcmstypes.Transaction, 0)
			if authTx != nil {
				mcmsTx, err := proposalutils.TransactionForChain(chainSel, nm.Address().Hex(), authTx.Data(), big.NewInt(0),
					string(shared.NonceManager), []string{})
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction for chain %d: %w", chainSel, err)
				}

				mcmsTransactions = append(mcmsTransactions, mcmsTx)
			}
			if prevRampsTx != nil {
				mcmsTx, err := proposalutils.TransactionForChain(chainSel, nm.Address().Hex(), prevRampsTx.Data(), big.NewInt(0),
					string(shared.NonceManager), []string{})
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction for chain %d: %w", chainSel, err)
				}

				mcmsTransactions = append(mcmsTransactions, mcmsTx)
			}
			if len(mcmsTransactions) == 0 {
				return cldf.ChangesetOutput{}, errors.New("no operations to batch")
			}

			batches = append(batches, mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(chainSel),
				Transactions:  mcmsTransactions,
			})

			timelocks[chainSel] = s.Chains[chainSel].Timelock.Address().Hex()
			inspectors[chainSel], err = proposalutils.McmsInspectorForChain(e, chainSel)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get inspector for chain %d: %w", chainSel, err)
			}
		}
	}
	if cfg.MCMS == nil {
		return cldf.ChangesetOutput{}, nil
	}
	mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(e, s, cfg.MCMS)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error getting mcms contract by chain: %w", err)
	}
	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmsContractByChain,
		inspectors,
		batches,
		"Update nonce manager for previous ramps and authorized callers",
		*cfg.MCMS,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
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
	MCMS *proposalutils.TimelockConfig
	// SkipOwnershipCheck allows you to bypass the ownership check for the onRamp.
	// WARNING: This should only be used when running this changeset within another changeset that is managing contract ownership!
	// Never use this option when running this changeset in isolation.
	SkipOwnershipCheck bool
}

func (cfg UpdateOnRampDestsConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return err
	}
	supportedChains := state.SupportedChains()
	for chainSel, updates := range cfg.UpdatesByChain {
		if err := stateview.ValidateChain(e, state, chainSel, cfg.MCMS); err != nil {
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
		if !cfg.SkipOwnershipCheck {
			if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.BlockChains.EVMChains()[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.OnRamp); err != nil {
				return err
			}
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

func (cfg UpdateOnRampDestsConfig) ToSequenceInput(state stateview.CCIPOnChainState) ccipseqs.OnRampApplyDestChainConfigUpdatesSequenceInput {
	updatesByChain := make(map[uint64]opsutil.EVMCallInput[[]onramp.OnRampDestChainConfigArgs], len(cfg.UpdatesByChain))
	for chainSel, updates := range cfg.UpdatesByChain {
		var args []onramp.OnRampDestChainConfigArgs
		for destination, update := range updates {
			router := common.HexToAddress("0x0")
			// If not enabled, set router to 0x0.
			if update.IsEnabled {
				if update.TestRouter {
					router = state.Chains[chainSel].TestRouter.Address()
				} else {
					router = state.Chains[chainSel].Router.Address()
				}
			}
			args = append(args, onramp.OnRampDestChainConfigArgs{
				DestChainSelector: destination,
				Router:            router,
				AllowlistEnabled:  update.AllowListEnabled,
			})
		}
		updatesByChain[chainSel] = opsutil.EVMCallInput[[]onramp.OnRampDestChainConfigArgs]{
			Address:       state.Chains[chainSel].OnRamp.Address(),
			ChainSelector: chainSel,
			CallInput:     args,
			NoSend:        cfg.MCMS != nil,
		}
	}

	return ccipseqs.OnRampApplyDestChainConfigUpdatesSequenceInput{
		UpdatesByChain: updatesByChain,
	}
}

// UpdateOnRampsDestsChangeset updates the onramp destinations for each onramp
// in the chains specified. Multichain support is important - consider when we add a new chain
// and need to update the onramp destinations for all chains to support the new chain.
func UpdateOnRampsDestsChangeset(e cldf.Environment, cfg UpdateOnRampDestsConfig) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	s, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.OnRampApplyDestChainConfigUpdatesSequence,
		e.BlockChains.EVMChains(),
		cfg.ToSequenceInput(s),
	)
	return opsutil.AddEVMCallSequenceToCSOutput(e, s, cldf.ChangesetOutput{}, report, err, cfg.MCMS, "Call ApplyDestChainConfigUpdates on OnRamps")
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
	MCMS *proposalutils.TimelockConfig
}

func (cfg UpdateOnRampDynamicConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	for chainSel, config := range cfg.UpdatesByChain {
		if err := stateview.ValidateChain(e, state, chainSel, cfg.MCMS); err != nil {
			return err
		}
		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.BlockChains.EVMChains()[chainSel].DeployerKey.From, state.Chains[chainSel].Timelock.Address(), state.Chains[chainSel].OnRamp); err != nil {
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

func UpdateOnRampDynamicConfigChangeset(e cldf.Environment, cfg UpdateOnRampDynamicConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	batches := []mcmstypes.BatchOperation{}
	timelocks := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)

	for chainSel, update := range cfg.UpdatesByChain {
		txOps := e.BlockChains.EVMChains()[chainSel].DeployerKey
		if cfg.MCMS != nil {
			txOps = cldf.SimTransactOpts()
		}
		onRamp := state.Chains[chainSel].OnRamp
		dynamicConfig, err := onRamp.GetDynamicConfig(nil)
		if err != nil {
			return cldf.ChangesetOutput{}, err
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
			if _, err := cldf.ConfirmIfNoErrorWithABI(e.BlockChains.EVMChains()[chainSel], tx, onramp.OnRampABI, err); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error updating onramp dynamic config for chain %s: %w", e.BlockChains.EVMChains()[chainSel].String(), err)
			}
		} else {
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}

			batchOperation, err := proposalutils.BatchOperationForChain(chainSel, onRamp.Address().Hex(), tx.Data(),
				big.NewInt(0), string(shared.OnRamp), []string{})
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			batches = append(batches, batchOperation)

			timelocks[chainSel] = state.Chains[chainSel].Timelock.Address().Hex()
			inspectors[chainSel], err = proposalutils.McmsInspectorForChain(e, chainSel)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get inspector for chain %d: %w", chainSel, err)
			}
		}
	}
	if cfg.MCMS == nil {
		return cldf.ChangesetOutput{}, nil
	}
	mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(e, state, cfg.MCMS)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error getting mcms contract by chain: %w", err)
	}
	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e, timelocks, mcmsContractByChain, inspectors, batches,
		"update onramp dynamic config",
		*cfg.MCMS)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
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
	MCMS *proposalutils.TimelockConfig
}

func (cfg UpdateOnRampAllowListConfig) Validate(env cldf.Environment) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for srcSel, updates := range cfg.UpdatesByChain {
		if err := stateview.ValidateChain(env, state, srcSel, cfg.MCMS); err != nil {
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
			signer = env.BlockChains.EVMChains()[srcSel].DeployerKey.From
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
			if err := stateview.ValidateChain(env, state, srcSel, cfg.MCMS); err != nil {
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

func UpdateOnRampAllowListChangeset(e cldf.Environment, cfg UpdateOnRampAllowListConfig) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	onchain, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	batches := []mcmstypes.BatchOperation{}
	timelocks := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)

	for srcSel, updates := range cfg.UpdatesByChain {
		txOps := e.BlockChains.EVMChains()[srcSel].DeployerKey
		if cfg.MCMS != nil {
			txOps = cldf.SimTransactOpts()
		}
		onRamp := onchain.Chains[srcSel].OnRamp
		args := make([]onramp.OnRampAllowlistConfigArgs, len(updates))
		for destSel, update := range updates {
			allowedSendersResp, err := onRamp.GetAllowedSendersList(nil, destSel)
			if err != nil {
				return cldf.ChangesetOutput{}, err
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
			if _, err := cldf.ConfirmIfNoErrorWithABI(e.BlockChains.EVMChains()[srcSel], tx, onramp.OnRampABI, err); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error updating allowlist for chain %d: %w", srcSel, err)
			}
		} else {
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}

			batchOperation, err := proposalutils.BatchOperationForChain(srcSel, onRamp.Address().Hex(), tx.Data(),
				big.NewInt(0), string(shared.OnRamp), []string{})
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			batches = append(batches, batchOperation)

			timelocks[srcSel] = onchain.Chains[srcSel].Timelock.Address().Hex()
			inspectors[srcSel], err = proposalutils.McmsInspectorForChain(e, srcSel)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get inspector for chain %d: %w", srcSel, err)
			}
		}
	}
	if cfg.MCMS == nil {
		return cldf.ChangesetOutput{}, nil
	}
	mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(e, onchain, cfg.MCMS)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error getting mcms contract by chain: %w", err)
	}
	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmsContractByChain,
		inspectors,
		batches,
		"update onramp allowlist",
		*cfg.MCMS,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

type WithdrawOnRampFeeTokensConfig struct {
	FeeTokensByChain map[uint64][]common.Address
	MCMS             *proposalutils.TimelockConfig
}

func (cfg WithdrawOnRampFeeTokensConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	for chainSel, feeTokens := range cfg.FeeTokensByChain {
		if err := stateview.ValidateChain(e, state, chainSel, cfg.MCMS); err != nil {
			return err
		}
		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.BlockChains.EVMChains()[chainSel].DeployerKey.From, state.Chains[chainSel].Timelock.Address(), state.Chains[chainSel].OnRamp); err != nil {
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

func WithdrawOnRampFeeTokensChangeset(e cldf.Environment, cfg WithdrawOnRampFeeTokensConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	batches := []mcmstypes.BatchOperation{}
	timelocks := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)

	for chainSel, feeTokens := range cfg.FeeTokensByChain {
		txOps := e.BlockChains.EVMChains()[chainSel].DeployerKey
		onRamp := state.Chains[chainSel].OnRamp
		tx, err := onRamp.WithdrawFeeTokens(txOps, feeTokens)
		if cfg.MCMS == nil {
			if _, err := cldf.ConfirmIfNoErrorWithABI(e.BlockChains.EVMChains()[chainSel], tx, onramp.OnRampABI, err); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error withdrawing fee tokens for chain %s: %w", e.BlockChains.EVMChains()[chainSel].String(), err)
			}
		} else {
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}

			batchOperation, err := proposalutils.BatchOperationForChain(chainSel, onRamp.Address().Hex(), tx.Data(),
				big.NewInt(0), string(shared.OnRamp), []string{})
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			batches = append(batches, batchOperation)

			timelocks[chainSel] = state.Chains[chainSel].Timelock.Address().Hex()
			inspectors[chainSel], err = proposalutils.McmsInspectorForChain(e, chainSel)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get inspector for chain %d: %w", chainSel, err)
			}
		}
	}
	if cfg.MCMS == nil {
		return cldf.ChangesetOutput{}, nil
	}
	mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(e, state, cfg.MCMS)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error getting mcms contract by chain: %w", err)
	}
	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmsContractByChain,
		inspectors,
		batches,
		"withdraw onramp fee tokens",
		*cfg.MCMS,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

type UpdateFeeQuoterPricesConfig struct {
	PricesByChain map[uint64]FeeQuoterPriceUpdatePerSource // source -> PriceDetails
	MCMS          *proposalutils.TimelockConfig
}

type FeeQuoterPriceUpdatePerSource struct {
	TokenPrices map[common.Address]*big.Int // token address -> price
	GasPrices   map[uint64]*big.Int         // dest chain -> gas price
}

func (cfg UpdateFeeQuoterPricesConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return err
	}
	for chainSel, initialPrice := range cfg.PricesByChain {
		if err := cldf.IsValidChainSelector(chainSel); err != nil {
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
		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.BlockChains.EVMChains()[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.FeeQuoter); err != nil {
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
		expectedAuthCaller := e.BlockChains.EVMChains()[chainSel].DeployerKey.From
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
			contains, err := cldf.AddressBookContains(e.ExistingAddresses, chainSel, token.String())
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
			if err := cldf.IsValidChainSelector(dest); err != nil {
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

func (cfg UpdateFeeQuoterPricesConfig) ToSequenceInput(state stateview.CCIPOnChainState) ccipseqs.FeeQuoterUpdatePricesSequenceInput {
	updates := make(map[uint64]opsutil.EVMCallInput[fee_quoter.InternalPriceUpdates], len(cfg.PricesByChain))
	for chainSel, prices := range cfg.PricesByChain {
		tokenPriceUpdates := make([]fee_quoter.InternalTokenPriceUpdate, len(prices.TokenPrices))
		i := 0
		for tokenAddress, price := range prices.TokenPrices {
			tokenPriceUpdates[i] = fee_quoter.InternalTokenPriceUpdate{
				SourceToken: tokenAddress,
				UsdPerToken: price,
			}
			i++
		}
		gasPriceUpdates := make([]fee_quoter.InternalGasPriceUpdate, len(prices.GasPrices))
		i = 0
		for destChainSelector, price := range prices.GasPrices {
			gasPriceUpdates[i] = fee_quoter.InternalGasPriceUpdate{
				DestChainSelector: destChainSelector,
				UsdPerUnitGas:     price,
			}
			i++
		}
		updates[chainSel] = opsutil.EVMCallInput[fee_quoter.InternalPriceUpdates]{
			ChainSelector: chainSel,
			Address:       state.Chains[chainSel].FeeQuoter.Address(),
			CallInput: fee_quoter.InternalPriceUpdates{
				TokenPriceUpdates: tokenPriceUpdates,
				GasPriceUpdates:   gasPriceUpdates,
			},
			NoSend: cfg.MCMS != nil, // If MCMS exists, we do not want to send the transaction.
		}
	}

	return ccipseqs.FeeQuoterUpdatePricesSequenceInput{
		UpdatesByChain: updates,
	}
}

func UpdateFeeQuoterPricesChangeset(e cldf.Environment, cfg UpdateFeeQuoterPricesConfig) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	s, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.FeeQuoterUpdatePricesSequence,
		e.BlockChains.EVMChains(),
		cfg.ToSequenceInput(s),
	)
	return opsutil.AddEVMCallSequenceToCSOutput(e, s, cldf.ChangesetOutput{}, report, err, cfg.MCMS, "Call UpdatePrices on FeeQuoters")
}

type UpdateFeeQuoterDestsConfig struct {
	// UpdatesByChain is a mapping from source -> dest -> config update.
	UpdatesByChain map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig
	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMS *proposalutils.TimelockConfig
}

func (cfg UpdateFeeQuoterDestsConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
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
		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.BlockChains.EVMChains()[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.FeeQuoter); err != nil {
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

func (cfg UpdateFeeQuoterDestsConfig) ToSequenceInput(state stateview.CCIPOnChainState) ccipseqs.FeeQuoterApplyDestChainConfigUpdatesSequenceInput {
	updates := make(map[uint64]opsutil.EVMCallInput[[]fee_quoter.FeeQuoterDestChainConfigArgs], len(cfg.UpdatesByChain))
	for chainSel, destChainUpdates := range cfg.UpdatesByChain {
		args := make([]fee_quoter.FeeQuoterDestChainConfigArgs, len(destChainUpdates))
		i := 0
		for destChainSel, destChainUpdate := range destChainUpdates {
			args[i] = fee_quoter.FeeQuoterDestChainConfigArgs{
				DestChainSelector: destChainSel,
				DestChainConfig:   destChainUpdate,
			}
			i++
		}
		updates[chainSel] = opsutil.EVMCallInput[[]fee_quoter.FeeQuoterDestChainConfigArgs]{
			Address:       state.Chains[chainSel].FeeQuoter.Address(),
			ChainSelector: chainSel,
			CallInput:     args,
			NoSend:        cfg.MCMS != nil, // If MCMS exists, we do not want to send the transaction.
		}
	}

	return ccipseqs.FeeQuoterApplyDestChainConfigUpdatesSequenceInput{
		UpdatesByChain: updates,
	}
}

func UpdateFeeQuoterDestsChangeset(e cldf.Environment, cfg UpdateFeeQuoterDestsConfig) (cldf.ChangesetOutput, error) {
	output := cldf.ChangesetOutput{}

	if err := cfg.Validate(e); err != nil {
		return output, err
	}
	s, err := stateview.LoadOnchainState(e)
	if err != nil {
		return output, err
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.FeeQuoterApplyDestChainConfigUpdatesSequence,
		e.BlockChains.EVMChains(),
		cfg.ToSequenceInput(s),
	)
	return opsutil.AddEVMCallSequenceToCSOutput(e, s, output, report, err, cfg.MCMS, "Call ApplyDestChainConfigUpdates on FeeQuoters")
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
	MCMS           *proposalutils.TimelockConfig
	// SkipOwnershipCheck allows you to bypass the ownership check for the offRamp.
	// WARNING: This should only be used when running this changeset within another changeset that is managing contract ownership!
	// Never use this option when running this changeset in isolation.
	SkipOwnershipCheck bool
}

func (cfg UpdateOffRampSourcesConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
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
		if !cfg.SkipOwnershipCheck {
			if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.BlockChains.EVMChains()[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.OffRamp); err != nil {
				return err
			}
		}

		for source := range updates {
			// Source cannot be an unknown
			if _, ok := supportedChains[source]; !ok {
				return fmt.Errorf("source chain %d is not a supported chain %s", source, chainState.OffRamp.Address())
			}

			if source == chainSel {
				return fmt.Errorf("cannot update offramp source to the same chain %d", source)
			}

			if err := state.ValidateRamp(source, shared.OnRamp); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cfg UpdateOffRampSourcesConfig) ToSequenceInput(state stateview.CCIPOnChainState) ccipseqs.OffRampApplySourceChainConfigUpdatesSequenceInput {
	updatesByChain := make(map[uint64]opsutil.EVMCallInput[[]offramp.OffRampSourceChainConfigArgs], len(cfg.UpdatesByChain))
	for chainSel, updates := range cfg.UpdatesByChain {
		var args []offramp.OffRampSourceChainConfigArgs
		for source, update := range updates {
			var router common.Address
			if update.TestRouter {
				router = state.Chains[chainSel].TestRouter.Address()
			} else {
				router = state.Chains[chainSel].Router.Address()
			}
			// can ignore err as validation checks for nil addresses
			onRampBytes, _ := state.GetOnRampAddressBytes(source)

			args = append(args, offramp.OffRampSourceChainConfigArgs{
				SourceChainSelector: source,
				Router:              router,
				IsEnabled:           update.IsEnabled,
				// TODO: how would this work when the onRamp is nonEVM?
				OnRamp:                    common.LeftPadBytes(onRampBytes, 32),
				IsRMNVerificationDisabled: update.IsRMNVerificationDisabled,
			})
		}
		updatesByChain[chainSel] = opsutil.EVMCallInput[[]offramp.OffRampSourceChainConfigArgs]{
			ChainSelector: chainSel,
			Address:       state.Chains[chainSel].OffRamp.Address(),
			CallInput:     args,
			NoSend:        cfg.MCMS != nil,
		}
	}

	return ccipseqs.OffRampApplySourceChainConfigUpdatesSequenceInput{
		UpdatesByChain: updatesByChain,
	}
}

// UpdateOffRampSourcesChangeset updates the offramp sources for each offramp.
func UpdateOffRampSourcesChangeset(e cldf.Environment, cfg UpdateOffRampSourcesConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.OffRampApplySourceChainConfigUpdatesSequence,
		e.BlockChains.EVMChains(),
		cfg.ToSequenceInput(state),
	)

	return opsutil.AddEVMCallSequenceToCSOutput(e, state, cldf.ChangesetOutput{}, report, err, cfg.MCMS, "Call ApplySourceChainConfigUpdates on OffRamps")
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
	MCMS           *proposalutils.TimelockConfig
	// SkipOwnershipCheck allows you to bypass the ownership check for the router.
	// WARNING: This should only be used when running this changeset within another changeset that is managing contract ownership!
	// Never use this option when running this changeset in isolation.
	SkipOwnershipCheck bool
}

func (cfg UpdateRouterRampsConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	if !cfg.TestRouter {
		// If not using the test router, we need to enforce MCMS usage if the state calls for it.
		err := state.EnforceMCMSUsageIfProd(e.GetContext(), cfg.MCMS)
		if err != nil {
			return err
		}
	}
	supportedChains := state.SupportedChains()
	for chainSel, update := range cfg.UpdatesByChain {
		if err := stateview.ValidateChain(e, state, chainSel, cfg.MCMS); err != nil {
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
		if !cfg.SkipOwnershipCheck {
			if cfg.TestRouter {
				// If activating on the test router, we don't need the other CCIP contracts to have proper ownership.
				if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.BlockChains.EVMChains()[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.TestRouter); err != nil {
					return err
				}
			} else if cfg.MCMS == nil {
				// If we are not using MCMS, then we know we aren't in a production environment given the EnforceMCMSUsageIfProd check above.
				// In this case, we only need to validate that the router contract is owned by the deployer key.
				// We don't care about uniform ownership in non-production environments.
				if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.BlockChains.EVMChains()[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.Router); err != nil {
					return err
				}
			} else {
				// If we are activating ramps on the main router in a production environment, we should validate two things:
				//   1. All expected CCIP contracts exist on the chain.
				//   2. All contracts have the expected owner.
				// That way, if cfg.MCMS exists, we ensure that every contract is owned by MCMS.
				// Calling this function will ensure that both these checks are done.
				err := state.ValidateOwnershipOfChain(e, chainSel, cfg.MCMS)
				if err != nil {
					return fmt.Errorf("failed to validate ownership of contracts on %s: %w", e.BlockChains.EVMChains()[chainSel], err)
				}
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
			if err := state.ValidateRamp(source, shared.OnRamp); err != nil {
				return err
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
			if err := state.ValidateRamp(destination, shared.OffRamp); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cfg UpdateRouterRampsConfig) ToSequenceInput(state stateview.CCIPOnChainState) ccipseqs.RouterApplyRampUpdatesSequenceInput {
	input := make(map[uint64]opsutil.EVMCallInput[ccipops.RouterApplyRampUpdatesOpInput], len(cfg.UpdatesByChain))
	for chainSel, update := range cfg.UpdatesByChain {
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
		input[chainSel] = opsutil.EVMCallInput[ccipops.RouterApplyRampUpdatesOpInput]{
			ChainSelector: chainSel,
			Address:       routerC.Address(),
			CallInput: ccipops.RouterApplyRampUpdatesOpInput{
				OnRampUpdates:  onRampUpdates,
				OffRampRemoves: removes,
				OffRampAdds:    adds,
			},
			NoSend: cfg.MCMS != nil, // If MCMS exists, we do not want to send the transaction.
		}
	}

	return ccipseqs.RouterApplyRampUpdatesSequenceInput{
		UpdatesByChain: input,
	}
}

// UpdateRouterRampsChangeset updates the on/offramps
// in either the router or test router for a series of chains. Use cases include:
// - Ramp upgrade. After deploying new ramps you can enable them on the test router and
// ensure it works e2e. Then enable the ramps on the real router.
// - New chain support. When adding a new chain, you can enable the new destination
// on all chains to support the new chain through the test router first. Once tested,
// Enable the new destination on the real router.
func UpdateRouterRampsChangeset(e cldf.Environment, cfg UpdateRouterRampsConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.RouterApplyRampUpdatesSequence,
		e.BlockChains.EVMChains(),
		cfg.ToSequenceInput(state),
	)
	return opsutil.AddEVMCallSequenceToCSOutput(e, state, cldf.ChangesetOutput{}, report, err, cfg.MCMS, "Call ApplyRampUpdates on Routers")
}

type SetOCR3OffRampConfig struct {
	HomeChainSel       uint64
	RemoteChainSels    []uint64
	CCIPHomeConfigType globals.ConfigType
	MCMS               *proposalutils.TimelockConfig
}

func (c SetOCR3OffRampConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	if err := stateview.ValidateChain(e, state, c.HomeChainSel, c.MCMS); err != nil {
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

func (c SetOCR3OffRampConfig) validateRemoteChain(e *cldf.Environment, state *stateview.CCIPOnChainState, chainSelector uint64) error {
	family, err := chain_selectors.GetSelectorFamily(chainSelector)
	if err != nil {
		return err
	}
	switch family {
	case chain_selectors.FamilySolana:
		chainState, ok := state.SolChains[chainSelector]
		if !ok {
			return fmt.Errorf("remote chain %d not found in onchain state", chainSelector)
		}
		if chainState.OffRamp.IsZero() {
			return fmt.Errorf("missing OffRamp for chain %d", chainSelector)
		}
	case chain_selectors.FamilyEVM:
		chainState, ok := state.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("remote chain %d not found in onchain state", chainSelector)
		}
		if err := commoncs.ValidateOwnership(e.GetContext(), c.MCMS != nil, e.BlockChains.EVMChains()[chainSelector].DeployerKey.From, chainState.Timelock.Address(), chainState.OffRamp); err != nil {
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
func SetOCR3OffRampChangeset(e cldf.Environment, cfg SetOCR3OffRampConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	batches := []mcmstypes.BatchOperation{}
	timelocks := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)

	for _, remote := range cfg.RemoteChainSels {
		donID, err := internal.DonIDForChain(
			state.Chains[cfg.HomeChainSel].CapabilityRegistry,
			state.Chains[cfg.HomeChainSel].CCIPHome,
			remote)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		args, err := internal.BuildSetOCR3ConfigArgs(
			donID, state.Chains[cfg.HomeChainSel].CCIPHome, remote, cfg.CCIPHomeConfigType)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		set, err := isOCR3ConfigSetOnOffRamp(e.Logger, e.BlockChains.EVMChains()[remote], state.Chains[remote].OffRamp, args)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if set {
			e.Logger.Infof("OCR3 config already set on offramp for chain %d", remote)
			continue
		}
		txOpts := e.BlockChains.EVMChains()[remote].DeployerKey
		if cfg.MCMS != nil {
			txOpts = cldf.SimTransactOpts()
		}
		offRamp := state.Chains[remote].OffRamp
		tx, err := offRamp.SetOCR3Configs(txOpts, args)
		if cfg.MCMS == nil {
			if _, err := cldf.ConfirmIfNoErrorWithABI(e.BlockChains.EVMChains()[remote], tx, offramp.OffRampABI, err); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error setting OCR3 config for chain %d: %w", remote, err)
			}
			e.Logger.Infow("Set OCR3 config on offramp", "chain", remote,
				"offRamp", offRamp.Address().Hex(), "donID", donID,
				"config", args)
		} else {
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}

			batchOperation, err := proposalutils.BatchOperationForChain(remote, offRamp.Address().Hex(), tx.Data(),
				big.NewInt(0), string(shared.OffRamp), []string{})
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			batches = append(batches, batchOperation)

			timelocks[remote] = state.Chains[remote].Timelock.Address().Hex()
			inspectors[remote], err = proposalutils.McmsInspectorForChain(e, remote)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get inspector for chain %d: %w", remote, err)
			}
		}
	}
	if cfg.MCMS == nil {
		return cldf.ChangesetOutput{}, nil
	}
	mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(e, state, cfg.MCMS)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error getting mcms contract by chain: %w", err)
	}
	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmsContractByChain,
		inspectors,
		batches,
		"Update OCR3 config",
		*cfg.MCMS,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	e.Logger.Info("Proposing OCR3 config update for", cfg.RemoteChainSels)
	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

type UpdateDynamicConfigOffRampConfig struct {
	Updates map[uint64]ccipops.OffRampParams
	MCMS    *proposalutils.TimelockConfig
}

func (cfg UpdateDynamicConfigOffRampConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return err
	}
	for chainSel, params := range cfg.Updates {
		if err := stateview.ValidateChain(e, state, chainSel, cfg.MCMS); err != nil {
			return fmt.Errorf("chain %d: %w", chainSel, err)
		}
		if state.Chains[chainSel].OffRamp == nil {
			return fmt.Errorf("missing offramp for chain %d", chainSel)
		}
		if state.Chains[chainSel].FeeQuoter == nil {
			return fmt.Errorf("missing fee quoter for chain %d", chainSel)
		}
		if params.GasForCallExactCheck > 0 {
			e.Logger.Infow(
				"GasForCallExactCheck is set, please note it's a static config and will be ignored for this changeset",
				"chain", chainSel, "gas", params.GasForCallExactCheck)
		}
		if err := commoncs.ValidateOwnership(
			e.GetContext(),
			cfg.MCMS != nil,
			e.BlockChains.EVMChains()[chainSel].DeployerKey.From,
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

func UpdateDynamicConfigOffRampChangeset(e cldf.Environment, cfg UpdateDynamicConfigOffRampConfig) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	batches := []mcmstypes.BatchOperation{}
	timelocks := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)

	for chainSel, params := range cfg.Updates {
		chain := e.BlockChains.EVMChains()[chainSel]
		txOpts := e.BlockChains.EVMChains()[chainSel].DeployerKey
		if cfg.MCMS != nil {
			txOpts = cldf.SimTransactOpts()
		}
		offRamp := state.Chains[chainSel].OffRamp
		dCfg := offramp.OffRampDynamicConfig{
			FeeQuoter:                               state.Chains[chainSel].FeeQuoter.Address(),
			PermissionLessExecutionThresholdSeconds: params.PermissionLessExecutionThresholdSeconds,
			MessageInterceptor:                      params.MessageInterceptor,
		}
		tx, err := offRamp.SetDynamicConfig(txOpts, dCfg)
		if cfg.MCMS == nil {
			if _, err := cldf.ConfirmIfNoErrorWithABI(e.BlockChains.EVMChains()[chainSel], tx, offramp.OffRampABI, err); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error updating offramp dynamic config for chain %d: %w", chainSel, err)
			}
			e.Logger.Infow("Updated offramp dynamic config", "chain", chain.String(), "config", dCfg)
		} else {
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}

			batchOperation, err := proposalutils.BatchOperationForChain(chainSel, offRamp.Address().Hex(), tx.Data(),
				big.NewInt(0), string(shared.OffRamp), []string{})
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			batches = append(batches, batchOperation)

			timelocks[chainSel] = state.Chains[chainSel].Timelock.Address().Hex()
			inspectors[chainSel], err = proposalutils.McmsInspectorForChain(e, chainSel)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get inspector for chain %d: %w", chainSel, err)
			}
		}
	}
	if cfg.MCMS == nil {
		return cldf.ChangesetOutput{}, nil
	}
	mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(e, state, cfg.MCMS)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error getting mcms contract by chain: %w", err)
	}
	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmsContractByChain,
		inspectors,
		batches,
		"Update offramp dynamic config",
		*cfg.MCMS,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	e.Logger.Infow("Proposing offramp dynamic config update", "config", cfg.Updates)
	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

func isOCR3ConfigSetOnOffRamp(
	lggr logger.Logger,
	chain cldf_evm.Chain,
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
		} else if destFamily == chain_selectors.FamilyAptos {
			familySelector, _ = hex.DecodeString(AptosFamilySelector) // aptos
		}
	}
	return fee_quoter.FeeQuoterDestChainConfig{
		IsEnabled:                         configEnabled,
		MaxNumberOfTokensPerMsg:           10,
		MaxDataBytes:                      30_000,
		MaxPerMsgGasLimit:                 3_000_000,
		DestGasOverhead:                   ccipevm.DestGasOverhead,
		DefaultTokenFeeUSDCents:           25,
		DestGasPerPayloadByteBase:         ccipevm.CalldataGasPerByteBase,
		DestGasPerPayloadByteHigh:         ccipevm.CalldataGasPerByteHigh,
		DestGasPerPayloadByteThreshold:    ccipevm.CalldataGasPerByteThreshold,
		DestDataAvailabilityOverheadGas:   100,
		DestGasPerDataAvailabilityByte:    16,
		DestDataAvailabilityMultiplierBps: 1,
		DefaultTokenDestGasOverhead:       90_000,
		DefaultTxGasLimit:                 200_000,
		GasMultiplierWeiPerEth:            11e17, // Gas multiplier in wei per eth is scaled by 1e18, so 11e17 is 1.1 = 110%
		NetworkFeeUSDCents:                10,
		ChainFamilySelector:               [4]byte(familySelector),
	}
}

type ApplyFeeTokensUpdatesConfig struct {
	UpdatesByChain map[uint64]ApplyFeeTokensUpdatesConfigPerChain
	MCMSConfig     *proposalutils.TimelockConfig
}

type ApplyFeeTokensUpdatesConfigPerChain struct {
	TokensToRemove []shared.TokenSymbol
	TokensToAdd    []shared.TokenSymbol
}

func (cfg ApplyFeeTokensUpdatesConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return err
	}
	for chainSel, updates := range cfg.UpdatesByChain {
		if err := stateview.ValidateChain(e, state, chainSel, cfg.MCMSConfig); err != nil {
			return err
		}
		chainState := state.Chains[chainSel]
		if chainState.FeeQuoter == nil {
			return fmt.Errorf("missing fee quoter for chain %d", chainSel)
		}
		tokenAddresses, err := chainState.TokenAddressBySymbol()
		if err != nil {
			return fmt.Errorf("error getting token addresses for chain %d: %w", chainSel, err)
		}
		for _, token := range updates.TokensToRemove {
			if _, ok := tokenAddresses[token]; !ok {
				return fmt.Errorf("token %s not found in state for chain %d", token, chainSel)
			}
		}
		for _, token := range updates.TokensToAdd {
			if _, ok := tokenAddresses[token]; !ok {
				return fmt.Errorf("token %s not found for in state chain %d", token, chainSel)
			}
		}
		if err := commoncs.ValidateOwnership(
			e.GetContext(),
			cfg.MCMSConfig != nil,
			e.BlockChains.EVMChains()[chainSel].DeployerKey.From,
			state.Chains[chainSel].Timelock.Address(),
			state.Chains[chainSel].FeeQuoter,
		); err != nil {
			return err
		}
	}
	return nil
}

// ApplyFeeTokensUpdatesFeeQuoterChangeset applies the token updates to the fee quoter to add or remove fee tokens.
// If MCMSConfig is provided, it will create a proposal to apply the changes assuming the fee quoter is owned by the timelock.
// If MCMSConfig is nil, it will apply the changes directly using the deployer key for each chain.
func ApplyFeeTokensUpdatesFeeQuoterChangeset(e cldf.Environment, cfg ApplyFeeTokensUpdatesConfig) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	var batches []mcmstypes.BatchOperation
	timelocks := make(map[uint64]string)
	inspectorPerChain := map[uint64]mcmssdk.Inspector{}
	for chainSel, updates := range cfg.UpdatesByChain {
		txOpts := e.BlockChains.EVMChains()[chainSel].DeployerKey
		if cfg.MCMSConfig != nil {
			txOpts = cldf.SimTransactOpts()
		}
		fq := state.Chains[chainSel].FeeQuoter
		tokenAddresses, err := state.Chains[chainSel].TokenAddressBySymbol()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error getting token addresses for chain %d: %w", chainSel, err)
		}
		var tokensToRemove, tokensToAdd []common.Address
		for _, token := range updates.TokensToRemove {
			tokensToRemove = append(tokensToRemove, tokenAddresses[token])
		}
		for _, token := range updates.TokensToAdd {
			tokensToAdd = append(tokensToAdd, tokenAddresses[token])
		}
		tx, err := fq.ApplyFeeTokensUpdates(txOpts, tokensToRemove, tokensToAdd)
		if cfg.MCMSConfig == nil {
			if _, err := cldf.ConfirmIfNoErrorWithABI(e.BlockChains.EVMChains()[chainSel], tx, fee_quoter.FeeQuoterABI, err); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error applying token updates for chain %d: %w", chainSel, err)
			}
		} else {
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			op, err := proposalutils.BatchOperationForChain(
				chainSel, fq.Address().String(), tx.Data(), big.NewInt(0), shared.FeeQuoter.String(), nil)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error creating batch operation for chain %d: %w", chainSel, err)
			}
			batches = append(batches, op)
			timelocks[chainSel] = state.Chains[chainSel].Timelock.Address().String()
			inspector, err := proposalutils.McmsInspectorForChain(e, chainSel)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error creating inspector for chain %d: %w", chainSel, err)
			}
			inspectorPerChain[chainSel] = inspector
		}
	}
	if cfg.MCMSConfig == nil {
		return cldf.ChangesetOutput{}, nil
	}
	mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(e, state, cfg.MCMSConfig)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error getting mcms contract by chain: %w", err)
	}
	p, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmsContractByChain,
		inspectorPerChain,
		batches,
		"Apply fee tokens updates",
		*cfg.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error building proposal: %w", err)
	}
	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcmslib.TimelockProposal{*p},
	}, nil
}

type UpdateTokenPriceFeedsConfig struct {
	Updates           map[uint64][]UpdateTokenPriceFeedsConfigPerChain
	FeedChainSelector uint64
	MCMS              *proposalutils.TimelockConfig
}

type UpdateTokenPriceFeedsConfigPerChain struct {
	SourceToken shared.TokenSymbol
	IsEnabled   bool
}

func (cfg UpdateTokenPriceFeedsConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return err
	}
	feedChainState, ok := state.Chains[cfg.FeedChainSelector]
	if !ok {
		return fmt.Errorf("feed chain %d not found in state", cfg.FeedChainSelector)
	}
	for chainSel, updates := range cfg.Updates {
		if err := stateview.ValidateChain(e, state, chainSel, cfg.MCMS); err != nil {
			return err
		}
		chainState := state.Chains[chainSel]
		if chainState.FeeQuoter == nil {
			return fmt.Errorf("missing fee quoter for chain %d", chainSel)
		}
		if feedChainState.USDFeeds == nil {
			return fmt.Errorf("missing token price feed for chain %d", chainSel)
		}
		tokenAddresses, err := chainState.TokenAddressBySymbol()
		if err != nil {
			return fmt.Errorf("error getting token addresses for chain %d: %w", chainSel, err)
		}
		for _, update := range updates {
			if _, ok := tokenAddresses[update.SourceToken]; !ok {
				return fmt.Errorf("token %s not found in state for chain %d", update.SourceToken, chainSel)
			}
			if _, ok := feedChainState.USDFeeds[update.SourceToken]; !ok {
				return fmt.Errorf("price feed for token %s not found in state for chain %d", update.SourceToken, chainSel)
			}
		}
		if err := commoncs.ValidateOwnership(
			e.GetContext(),
			cfg.MCMS != nil,
			e.BlockChains.EVMChains()[chainSel].DeployerKey.From,
			state.Chains[chainSel].Timelock.Address(),
			state.Chains[chainSel].FeeQuoter,
		); err != nil {
			return err
		}
	}
	return nil
}

// UpdateTokenPriceFeedsFeeQuoterChangeset applies the token price feed updates to the fee quoter.
// Before applying the changeset, ensure that the environment state/addressbook is up to date with latest token and price feed addresses.
// If MCMS is provided, it will create a proposal to apply the changes assuming the fee quoter is owned by the timelock.
// If MCMS is nil, it will apply the changes directly using the deployer key for each chain.
func UpdateTokenPriceFeedsFeeQuoterChangeset(e cldf.Environment, cfg UpdateTokenPriceFeedsConfig) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	var batches []mcmstypes.BatchOperation
	timelocks := make(map[uint64]string)
	inspectorPerChain := map[uint64]mcmssdk.Inspector{}
	for chainSel, updates := range cfg.Updates {
		txOpts := e.BlockChains.EVMChains()[chainSel].DeployerKey
		if cfg.MCMS != nil {
			txOpts = cldf.SimTransactOpts()
		}
		fq := state.Chains[chainSel].FeeQuoter
		tokenAddresses, err := state.Chains[chainSel].TokenAddressBySymbol()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error getting token addresses for chain %d: %w", chainSel, err)
		}
		tokenDetails, err := state.Chains[chainSel].TokenDetailsBySymbol()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error getting token decimal for chain %d: %w", chainSel, err)
		}
		var priceFeedUpdates []fee_quoter.FeeQuoterTokenPriceFeedUpdate
		for _, update := range updates {
			_, ok := tokenDetails[update.SourceToken]
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("token details %s not found in state for chain %d", update.SourceToken, chainSel)
			}
			decimal, err := tokenDetails[update.SourceToken].Decimals(&bind.CallOpts{
				Context: e.GetContext(),
			})
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error getting token decimal for chain %d: %w", chainSel, err)
			}
			feed := state.Chains[cfg.FeedChainSelector].USDFeeds[update.SourceToken]
			priceFeedUpdates = append(priceFeedUpdates, fee_quoter.FeeQuoterTokenPriceFeedUpdate{
				SourceToken: tokenAddresses[update.SourceToken],
				FeedConfig: fee_quoter.FeeQuoterTokenPriceFeedConfig{
					DataFeedAddress: feed.Address(),
					TokenDecimals:   decimal,
					IsEnabled:       update.IsEnabled,
				},
			})
		}
		tx, err := fq.UpdateTokenPriceFeeds(txOpts, priceFeedUpdates)
		if cfg.MCMS == nil {
			if _, err := cldf.ConfirmIfNoErrorWithABI(e.BlockChains.EVMChains()[chainSel], tx, fee_quoter.FeeQuoterABI, err); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error applying token price feed update for chain %d: %w", chainSel, err)
			}
		} else {
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			op, err := proposalutils.BatchOperationForChain(
				chainSel, fq.Address().String(), tx.Data(), big.NewInt(0), shared.FeeQuoter.String(), nil)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error creating batch operation for chain %d: %w", chainSel, err)
			}
			batches = append(batches, op)
			timelocks[chainSel] = state.Chains[chainSel].Timelock.Address().String()
			inspector, err := proposalutils.McmsInspectorForChain(e, chainSel)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error getting inspector for chain %d: %w", chainSel, err)
			}
			inspectorPerChain[chainSel] = inspector
		}
	}
	if cfg.MCMS == nil {
		return cldf.ChangesetOutput{}, nil
	}
	mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(e, state, cfg.MCMS)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error getting mcms contract by chain: %w", err)
	}
	p, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmsContractByChain,
		inspectorPerChain,
		batches,
		"Update token price feeds",
		*cfg.MCMS,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error building proposal: %w", err)
	}
	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcmslib.TimelockProposal{*p},
	}, nil
}

type PremiumMultiplierWeiPerEthUpdatesConfig struct {
	Updates map[uint64][]PremiumMultiplierWeiPerEthUpdatesConfigPerChain
	MCMS    *proposalutils.TimelockConfig
}

func (cfg PremiumMultiplierWeiPerEthUpdatesConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return err
	}
	for chainSel, updates := range cfg.Updates {
		if err := stateview.ValidateChain(e, state, chainSel, cfg.MCMS); err != nil {
			return err
		}
		chainState := state.Chains[chainSel]
		if chainState.FeeQuoter == nil {
			return fmt.Errorf("missing fee quoter for chain %d", chainSel)
		}
		tokenAddresses, err := chainState.TokenAddressBySymbol()
		if err != nil {
			return fmt.Errorf("error getting token addresses for chain %d: %w", chainSel, err)
		}
		for _, update := range updates {
			if _, ok := tokenAddresses[update.Token]; !ok {
				return fmt.Errorf("token %s not found in state for chain %d", update.Token, chainSel)
			}
			if update.PremiumMultiplierWeiPerEth == 0 {
				return fmt.Errorf("missing premium multiplier for chain %d", chainSel)
			}
		}
		if err := commoncs.ValidateOwnership(
			e.GetContext(),
			cfg.MCMS != nil,
			e.BlockChains.EVMChains()[chainSel].DeployerKey.From,
			state.Chains[chainSel].Timelock.Address(),
			state.Chains[chainSel].FeeQuoter,
		); err != nil {
			return err
		}
	}
	return nil
}

type PremiumMultiplierWeiPerEthUpdatesConfigPerChain struct {
	Token                      shared.TokenSymbol
	PremiumMultiplierWeiPerEth uint64
}

// ApplyPremiumMultiplierWeiPerEthUpdatesFeeQuoterChangeset applies the premium multiplier updates for mentioned tokens to the fee quoter.
// If MCMS is provided, it will create a proposal to apply the changes assuming the fee quoter is owned by the timelock.
// If MCMS is nil, it will apply the changes directly using the deployer key for each chain.
func ApplyPremiumMultiplierWeiPerEthUpdatesFeeQuoterChangeset(e cldf.Environment, cfg PremiumMultiplierWeiPerEthUpdatesConfig) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	var batches []mcmstypes.BatchOperation
	timelocks := make(map[uint64]string)
	inspectorPerChain := map[uint64]mcmssdk.Inspector{}
	for chainSel, updates := range cfg.Updates {
		txOpts := e.BlockChains.EVMChains()[chainSel].DeployerKey
		if cfg.MCMS != nil {
			txOpts = cldf.SimTransactOpts()
		}
		fq := state.Chains[chainSel].FeeQuoter
		tokenAddresses, err := state.Chains[chainSel].TokenAddressBySymbol()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error getting token addresses for chain %d: %w", chainSel, err)
		}
		var premiumMultiplierUpdates []fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs
		for _, update := range updates {
			premiumMultiplierUpdates = append(premiumMultiplierUpdates, fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs{
				Token:                      tokenAddresses[update.Token],
				PremiumMultiplierWeiPerEth: update.PremiumMultiplierWeiPerEth,
			})
		}
		tx, err := fq.ApplyPremiumMultiplierWeiPerEthUpdates(txOpts, premiumMultiplierUpdates)
		if cfg.MCMS == nil {
			if _, err := cldf.ConfirmIfNoErrorWithABI(e.BlockChains.EVMChains()[chainSel], tx, fee_quoter.FeeQuoterABI, err); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error applying premium multiplier updates for chain %d: %w", chainSel, err)
			}
		} else {
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			op, err := proposalutils.BatchOperationForChain(
				chainSel, fq.Address().String(), tx.Data(), big.NewInt(0), shared.FeeQuoter.String(), nil)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error creating batch operation for chain %d: %w", chainSel, err)
			}
			batches = append(batches, op)
			timelocks[chainSel] = state.Chains[chainSel].Timelock.Address().String()
			inspector, err := proposalutils.McmsInspectorForChain(e, chainSel)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error getting inspector for chain %d: %w", chainSel, err)
			}
			inspectorPerChain[chainSel] = inspector
		}
	}
	if cfg.MCMS == nil {
		return cldf.ChangesetOutput{}, nil
	}
	mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(e, state, cfg.MCMS)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error getting mcms contract by chain: %w", err)
	}
	p, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmsContractByChain,
		inspectorPerChain,
		batches,
		"Apply premium multiplier updates",
		*cfg.MCMS,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error building proposal: %w", err)
	}
	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcmslib.TimelockProposal{*p},
	}, nil
}

type ApplyTokenTransferFeeConfigUpdatesConfig struct {
	UpdatesByChain map[uint64]ApplyTokenTransferFeeConfigUpdatesConfigPerChain
	MCMS           *proposalutils.TimelockConfig
}

func (cfg ApplyTokenTransferFeeConfigUpdatesConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return err
	}
	for chainSel, updates := range cfg.UpdatesByChain {
		if err := stateview.ValidateChain(e, state, chainSel, cfg.MCMS); err != nil {
			return err
		}
		chainState := state.Chains[chainSel]
		if chainState.FeeQuoter == nil {
			return fmt.Errorf("missing fee quoter for chain %d", chainSel)
		}
		tokenAddresses, err := chainState.TokenAddressBySymbol()
		if err != nil {
			return fmt.Errorf("error getting token addresses for chain %d: %w", chainSel, err)
		}
		for _, update := range updates.TokenTransferFeeConfigArgs {
			if update.DestChain == chainSel {
				return fmt.Errorf("dest chain %d cannot be the same as source chain %d", update.DestChain, chainSel)
			}
			for token, feeConfig := range update.TokenTransferFeeConfigPerToken {
				if _, ok := tokenAddresses[token]; !ok {
					return fmt.Errorf("token %s not found in state for chain %d", token, chainSel)
				}
				if feeConfig.MinFeeUSDCents >= feeConfig.MaxFeeUSDCents {
					return fmt.Errorf("min fee must be less than max fee for token %s in chain %d", token, chainSel)
				}
				if feeConfig.DestBytesOverhead < globals.CCIPLockOrBurnV1RetBytes {
					return fmt.Errorf("dest bytes overhead must be at least %d for token %s in chain %d", globals.CCIPLockOrBurnV1RetBytes, token, chainSel)
				}
			}
			if err := stateview.ValidateChain(e, state, update.DestChain, nil); err != nil {
				return fmt.Errorf("dest chain %d: %w", update.DestChain, err)
			}
		}
		for _, remove := range updates.TokenTransferFeeConfigRemoveArgs {
			if remove.DestChain == chainSel {
				return fmt.Errorf("dest chain %d cannot be the same as source chain %d", remove.DestChain, chainSel)
			}
			if _, ok := tokenAddresses[remove.Token]; !ok {
				return fmt.Errorf("token %s not found in state for chain %d", remove.Token, chainSel)
			}
			if err := stateview.ValidateChain(e, state, remove.DestChain, nil); err != nil {
				return fmt.Errorf("dest chain %d: %w", remove.DestChain, err)
			}
			_, err := chainState.FeeQuoter.GetTokenTransferFeeConfig(&bind.CallOpts{
				Context: e.GetContext(),
			}, remove.DestChain, tokenAddresses[remove.Token])
			if err != nil {
				return fmt.Errorf("is the token already updated with token transfer fee config ?"+
					"error getting token transfer fee config for token %s in chain %d: %w", remove.Token, chainSel, err)
			}
		}
		if err := commoncs.ValidateOwnership(
			e.GetContext(),
			cfg.MCMS != nil,
			e.BlockChains.EVMChains()[chainSel].DeployerKey.From,
			chainState.Timelock.Address(),
			chainState.FeeQuoter,
		); err != nil {
			return err
		}
	}
	return nil
}

type ApplyTokenTransferFeeConfigUpdatesConfigPerChain struct {
	TokenTransferFeeConfigArgs       []TokenTransferFeeConfigArg
	TokenTransferFeeConfigRemoveArgs []TokenTransferFeeConfigRemoveArg
}

type TokenTransferFeeConfigArg struct {
	DestChain                      uint64
	TokenTransferFeeConfigPerToken map[shared.TokenSymbol]fee_quoter.FeeQuoterTokenTransferFeeConfig
}

type TokenTransferFeeConfigRemoveArg struct {
	DestChain uint64
	Token     shared.TokenSymbol
}

// ApplyTokenTransferFeeConfigUpdatesFeeQuoterChangeset applies the token transfer fee config updates for provided tokens to the fee quoter.
// If TokenTransferFeeConfigRemoveArgs is provided, it will remove the token transfer fee config for the provided tokens and dest chains.
// If TokenTransferFeeConfigArgs is provided, it will update the token transfer fee config for the provided tokens and dest chains.
// Use this changeset whenever there is a need to update custom token transfer fee config for a chain, dest chain and token.
// If MCMS is provided, it will create a proposal to apply the changes assuming the fee quoter is owned by the timelock.
// If MCMS is nil, it will apply the changes directly using the deployer key for each chain.
func ApplyTokenTransferFeeConfigUpdatesFeeQuoterChangeset(e cldf.Environment, cfg ApplyTokenTransferFeeConfigUpdatesConfig) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	var batches []mcmstypes.BatchOperation
	timelocks := make(map[uint64]string)
	inspectorPerChain := map[uint64]mcmssdk.Inspector{}
	for chainSel, updates := range cfg.UpdatesByChain {
		txOpts := e.BlockChains.EVMChains()[chainSel].DeployerKey
		if cfg.MCMS != nil {
			txOpts = cldf.SimTransactOpts()
		}
		fq := state.Chains[chainSel].FeeQuoter
		tokenAddresses, err := state.Chains[chainSel].TokenAddressBySymbol()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("error getting token addresses for chain %d: %w", chainSel, err)
		}
		var tokenTransferFeeConfigs []fee_quoter.FeeQuoterTokenTransferFeeConfigArgs
		for _, update := range updates.TokenTransferFeeConfigArgs {
			var tokenTransferFeeConfigPerToken []fee_quoter.FeeQuoterTokenTransferFeeConfigSingleTokenArgs
			for token, feeConfig := range update.TokenTransferFeeConfigPerToken {
				tokenTransferFeeConfigPerToken = append(tokenTransferFeeConfigPerToken, fee_quoter.FeeQuoterTokenTransferFeeConfigSingleTokenArgs{
					Token:                  tokenAddresses[token],
					TokenTransferFeeConfig: feeConfig,
				})
			}
			tokenTransferFeeConfigs = append(tokenTransferFeeConfigs, fee_quoter.FeeQuoterTokenTransferFeeConfigArgs{
				DestChainSelector:       update.DestChain,
				TokenTransferFeeConfigs: tokenTransferFeeConfigPerToken,
			})
		}
		var tokenTransferFeeConfigsRemove []fee_quoter.FeeQuoterTokenTransferFeeConfigRemoveArgs
		for _, remove := range updates.TokenTransferFeeConfigRemoveArgs {
			tokenTransferFeeConfigsRemove = append(tokenTransferFeeConfigsRemove, fee_quoter.FeeQuoterTokenTransferFeeConfigRemoveArgs{
				DestChainSelector: remove.DestChain,
				Token:             tokenAddresses[remove.Token],
			})
		}
		tx, err := fq.ApplyTokenTransferFeeConfigUpdates(txOpts, tokenTransferFeeConfigs, tokenTransferFeeConfigsRemove)
		if cfg.MCMS == nil {
			if _, err := cldf.ConfirmIfNoErrorWithABI(e.BlockChains.EVMChains()[chainSel], tx, fee_quoter.FeeQuoterABI, err); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error applying token transfer fee config updates for chain %d: %w", chainSel, err)
			}
		} else {
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			op, err := proposalutils.BatchOperationForChain(
				chainSel, fq.Address().String(), tx.Data(), big.NewInt(0), shared.FeeQuoter.String(), nil)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error creating batch operation for chain %d: %w", chainSel, err)
			}
			batches = append(batches, op)
			timelocks[chainSel] = state.Chains[chainSel].Timelock.Address().String()
			inspector, err := proposalutils.McmsInspectorForChain(e, chainSel)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("error getting inspector for chain %d: %w", chainSel, err)
			}
			inspectorPerChain[chainSel] = inspector
		}
	}
	if cfg.MCMS == nil {
		return cldf.ChangesetOutput{}, nil
	}
	mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(e, state, cfg.MCMS)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error getting mcms contract by chain: %w", err)
	}
	p, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmsContractByChain,
		inspectorPerChain,
		batches,
		"Apply token transfer fee config updates",
		*cfg.MCMS,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error building proposal: %w", err)
	}
	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcmslib.TimelockProposal{*p},
	}, nil
}
