package v1_5_1

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/helpers/pointer"
)

var _ cldf.ChangeSetV2[SetTokenTransferFeeConfig] = SetTokenTransferFeeConfigChangeset

// SetTokenTransferFeeConfigChangeset is a changeset that allows you to set configurations such as DestGasOverhead
// for v1_5 lanes. The changeset is intended to replace the current approach where we use RDD CCIP and/or Gauntlet
// to perform this op: https://github.com/smartcontractkit/reference-data-directory-ccip/pull/1656/files. If you'd
// like to view the underlying solidity function that this changeset invokes, then the code for EVM2EVMOnRamp will
// be of interest: https://etherscan.io/address/0xb8a882f3B88bd52D1Ff56A873bfDB84b70431937#code.
var SetTokenTransferFeeConfigChangeset = cldf.CreateChangeSet(setTokenTransferFeeConfigLogic, setTokenTransferFeeConfigPrecondition)

type SetTokenTransferFeeConfig struct {
	// A map of chain selector => token transfer fee input which describes the updates to make on each chain
	InputsByChain map[uint64]SetTokenTransferFeeArgs `json:"inputsByChain"`

	// The timelock config - all updates can be merged into one MCMS proposal with this setting
	MCMS *proposalutils.TimelockConfig `json:"mcms"`
}

type SetTokenTransferFeeArgs struct {
	// Tokens specified here will be given custom token transfer fee configs (isEnabled will be set to true on-chain)
	TokenTransferFeeConfigArgs []TokenTransferFeeArgs

	// Tokens specified here will have their custom token transfer fee configs reset (isEnabled will be set to false on-chain)
	TokensToUseDefaultFeeConfigs []common.Address
}

// NOTE: the contract's _setTokenTransferFeeConfig method replaces the entire TokenTransferFeeConfig
// for each token included in TokenTransferFeeConfigArgs (it does not merge fields) so we need to be
// extra careful here. In Go, it is *very* easy to unintentionally pass an input struct with missing
// fields to a function without realizing that default values are being used for the missing fields.
// This isn't ideal for this changeset because it can lead to some configs being set to zero onchain
// (which more often than not is a mistake). To avoid these types of situations, we use pointers for
// each config field in the struct below to make things much more explicit. If a field isn't defined
// or set to nil then the Go code will fall back to the value onchain before sending. Otherwise, the
// user's provided value is used. When a token has no preexisting custom config (isEnabled == false)
// then all fields must be explicitly provided by the caller.
type TokenTransferFeeArgs struct {
	Token                     common.Address
	MinFeeUSDCents            *uint32
	MaxFeeUSDCents            *uint32
	DeciBps                   *uint16
	DestGasOverhead           *uint32
	DestBytesOverhead         *uint32
	AggregateRateLimitEnabled *bool
}

func (args TokenTransferFeeArgs) HasMissingFields() bool {
	return args.MinFeeUSDCents == nil ||
		args.MaxFeeUSDCents == nil ||
		args.DeciBps == nil ||
		args.DestGasOverhead == nil ||
		args.DestBytesOverhead == nil ||
		args.AggregateRateLimitEnabled == nil
}

func setTokenTransferFeeConfigPrecondition(e cldf.Environment, cfg SetTokenTransferFeeConfig) error {
	if len(cfg.InputsByChain) == 0 {
		e.Logger.Warn("no inputs were provided - exiting precondition stage gracefully")
		return nil
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for selector, input := range cfg.InputsByChain {
		chainState, ok := state.EVMChainState(selector)
		if !ok {
			return fmt.Errorf("selector %d does not exist in state", selector)
		}
		if onramp := chainState.EVM2EVMOnRamp[selector]; onramp == nil {
			return fmt.Errorf("missing EVM2EVMOnRamp on chain with selector %d", selector)
		}
		if err := stateview.ValidateChain(e, state, selector, cfg.MCMS); err != nil {
			return fmt.Errorf("failed to validate chain %d: %w", selector, err)
		}

		tokensToReset := map[common.Address]bool{}
		for _, address := range input.TokensToUseDefaultFeeConfigs {
			if _, exists := tokensToReset[address]; exists {
				return fmt.Errorf("found duplicate address (%s) in TokensToUseDefaultFeeConfigs for chain with selector %d", address.Hex(), selector)
			}
			if address == utils.ZeroAddress {
				return fmt.Errorf("for selector %d - zero address is not allowed in TokensToUseDefaultFeeConfigs", selector)
			}
			tokensToReset[address] = true
		}

		tokensToSetup := map[common.Address]bool{}
		for _, args := range input.TokenTransferFeeConfigArgs {
			if _, exists := tokensToSetup[args.Token]; exists {
				return fmt.Errorf("found duplicate token address (%s) in TokenTransferFeeConfigArgs for chain with selector %d", args.Token.Hex(), selector)
			}
			if args.Token == utils.ZeroAddress {
				return fmt.Errorf("for selector %d - zero address is not allowed in TokenTransferFeeConfigArgs", selector)
			}
			tokensToSetup[args.Token] = true
		}

		for addr := range tokensToSetup {
			if _, exists := tokensToReset[addr]; exists {
				return fmt.Errorf("for selector %d - TokensToUseDefaultFeeConfigs and TokenTransferFeeConfigArgs must not reference the same address (%s)", selector, addr.Hex())
			}
		}
	}

	return nil
}

func setTokenTransferFeeConfigLogic(e cldf.Environment, cfg SetTokenTransferFeeConfig) (cldf.ChangesetOutput, error) {
	if len(cfg.InputsByChain) == 0 {
		e.Logger.Warn("no inputs were provided - exiting apply stage gracefully")
		return cldf.ChangesetOutput{}, nil
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	e.Logger.Info("preparing deployer group transactions")
	deployerGroup := deployergroup.NewDeployerGroup(e, state, cfg.MCMS).WithDeploymentContext("SetTokenTransferFeeConfig")
	for selector, input := range cfg.InputsByChain {
		e.Logger.Infof("processing chain with selector %d", selector)

		blockChain, exists := e.BlockChains.EVMChains()[selector]
		if !exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("could not find EVM chain with selector %d in environment", selector)
		}

		chainState, exists := state.Chains[selector]
		if !exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("could not find chain %s in state", blockChain.String())
		}

		opts, err := deployerGroup.GetDeployer(selector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain %s: %w", blockChain.String(), err)
		}

		onramp, exists := chainState.EVM2EVMOnRamp[selector]
		if !exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("could not find EVM2EVMOnRamp for chain %s", blockChain.String())
		}

		tokenTransferFeeConfigArgs := []evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{}
		for _, args := range input.TokenTransferFeeConfigArgs {
			// This gets the token transfer fee config for the given token - if it doesn't exist, then the zero struct will be returned and `IsEnabled` will be `false`
			e.Logger.Infof("getting token transfer fee config for token %s on chain %s", args.Token.Hex(), blockChain.String())
			curConfig, err := onramp.GetTokenTransferFeeConfig(&bind.CallOpts{Context: e.GetContext()}, args.Token)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get TokenTransferFeeConfig for token %s on chain %s: %w", args.Token.Hex(), blockChain.String(), err)
			}

			// If no custom config already exists on-chain for the token, then we have no fallback values to use - in this case the caller must explicitly provide all fields
			if !curConfig.IsEnabled && args.HasMissingFields() {
				return cldf.ChangesetOutput{}, fmt.Errorf("token %s on %s: when enabling a new token, all fields must be provided", args.Token.Hex(), blockChain.String())
			}

			// At this point, we're either using fallback values from the chain or the caller has explicitly provided the inputs
			newConfig := evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{
				Token:                     args.Token,
				MinFeeUSDCents:            pointer.Coalesce(args.MinFeeUSDCents, curConfig.MinFeeUSDCents),
				MaxFeeUSDCents:            pointer.Coalesce(args.MaxFeeUSDCents, curConfig.MaxFeeUSDCents),
				DeciBps:                   pointer.Coalesce(args.DeciBps, curConfig.DeciBps),
				DestGasOverhead:           pointer.Coalesce(args.DestGasOverhead, curConfig.DestGasOverhead),
				DestBytesOverhead:         pointer.Coalesce(args.DestBytesOverhead, curConfig.DestBytesOverhead),
				AggregateRateLimitEnabled: pointer.Coalesce(args.AggregateRateLimitEnabled, curConfig.AggregateRateLimitEnabled),
			}

			// Check if the new config is different from the on-chain config
			isDifferent := !curConfig.IsEnabled
			if curConfig.IsEnabled {
				isDifferent = newConfig.MinFeeUSDCents != curConfig.MinFeeUSDCents ||
					newConfig.MaxFeeUSDCents != curConfig.MaxFeeUSDCents ||
					newConfig.DeciBps != curConfig.DeciBps ||
					newConfig.DestGasOverhead != curConfig.DestGasOverhead ||
					newConfig.DestBytesOverhead != curConfig.DestBytesOverhead ||
					newConfig.AggregateRateLimitEnabled != curConfig.AggregateRateLimitEnabled
			}

			// Only perform an update if the new config is different from the on-chain config
			if isDifferent {
				tokenTransferFeeConfigArgs = append(tokenTransferFeeConfigArgs, newConfig)
			} else {
				e.Logger.Infof("token %s on %s: input config is the same as on-chain config (%+v) - skipping", args.Token.Hex(), blockChain.String(), curConfig)
			}
		}

		resetsCount := len(input.TokensToUseDefaultFeeConfigs)
		updateCount := len(tokenTransferFeeConfigArgs)
		if updateCount == 0 && resetsCount == 0 {
			e.Logger.Infof("detected no changes for chain %s - skipping", blockChain.String())
			continue
		}

		e.Logger.Infof("setting token transfer fee config for chain %s (updates = %d, resets = %d)", blockChain.String(), updateCount, resetsCount)
		_, err = onramp.SetTokenTransferFeeConfig(opts,
			tokenTransferFeeConfigArgs,
			input.TokensToUseDefaultFeeConfigs,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf(
				"failed to create SetTokenTransferFeeConfig transaction on chain %s: %w",
				blockChain.String(),
				err,
			)
		}
	}

	e.Logger.Info("running deployer group")
	return deployerGroup.Enact()
}
