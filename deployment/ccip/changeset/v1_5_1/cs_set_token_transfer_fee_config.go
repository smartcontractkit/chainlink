package v1_5_1

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
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
	// A mapping from src chain selector => dst chain selector => token transfer fee input
	InputsByChain map[uint64]map[uint64]SetTokenTransferFeeArgs `json:"inputsByChain"`

	// The timelock config - all updates can be merged into one MCMS proposal with this setting
	MCMS *proposalutils.TimelockConfig `json:"mcms"`
}

type SetTokenTransferFeeArgs struct {
	// Tokens specified here will be given custom token transfer fee configs (isEnabled will be set to true on-chain)
	TokenTransferFeeConfigArgs map[common.Address]TokenTransferFeeArgs

	// Tokens specified here will have their custom token transfer fee configs reset (isEnabled will be set to false on-chain)
	TokensToUseDefaultFeeConfigs []common.Address
}

// NOTE: the _setTokenTransferFeeConfig method in the Solidity contract overwrites all config values
// for each token included in TokenTransferFeeConfigArgs (it doesn't upsert values) so we need to be
// extra careful here. In Go, it is *very* easy to unintentionally pass an input struct with missing
// fields to a func without realizing that zero values are really being used for the missing fields.
// To avoid these types of situations, we use pointers for each config field in the struct below. If
// a field is undefined or set to nil in the struct, then we will fallback to using any pre-existing
// values from the chain before sending the transaction. Otherwise the user's input values are used.
// If a token has no pre-existing config values on-chain (i.e. isEnabled == false), then every field
// must be explicitly provided by the caller.
type TokenTransferFeeArgs struct {
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

func setTokenTransferFeeConfigPrecondition(env cldf.Environment, cfg SetTokenTransferFeeConfig) error {
	if len(cfg.InputsByChain) == 0 {
		env.Logger.Warn("no inputs were provided - exiting precondition stage gracefully")
		return nil
	}

	state, err := stateview.LoadOnchainState(env, stateview.WithLoadLegacyContracts(true))
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for srcSelector, inputs := range cfg.InputsByChain {
		err := stateview.ValidateChain(env, state, srcSelector, cfg.MCMS)
		if err != nil {
			return fmt.Errorf("failed to validate src chain (src = %d): %w", srcSelector, err)
		}
		chainState, ok := state.EVMChainState(srcSelector)
		if !ok {
			return fmt.Errorf("selector does not exist in EVM chain state (src = %d)", srcSelector)
		}

		for dstSelector, input := range inputs {
			if err := stateview.ValidateChain(env, state, dstSelector, cfg.MCMS); err != nil {
				return fmt.Errorf("failed to validate dst chain (src = %d, dst = %d): %w", srcSelector, dstSelector, err)
			}
			if _, exists := chainState.EVM2EVMOnRamp[dstSelector]; !exists {
				return fmt.Errorf("no EVM2EVMOnRamp exists (src = %d, dst = %d)", srcSelector, dstSelector)
			}

			tokensToReset := map[common.Address]bool{}
			for _, tokenAddress := range input.TokensToUseDefaultFeeConfigs {
				if _, exists := tokensToReset[tokenAddress]; exists {
					return fmt.Errorf("duplicate address in TokensToUseDefaultFeeConfigs (src = %d, dst = %d, addr = %s)", srcSelector, dstSelector, tokenAddress.Hex())
				}
				if tokenAddress == utils.ZeroAddress {
					return fmt.Errorf("zero address not allowed in TokensToUseDefaultFeeConfigs (src = %d, dst = %d)", srcSelector, dstSelector)
				}
				tokensToReset[tokenAddress] = true
			}

			for tokenAddress := range input.TokenTransferFeeConfigArgs {
				if tokenAddress == utils.ZeroAddress {
					return fmt.Errorf("zero address not allowed in TokenTransferFeeConfigArgs (src = %d, dst = %d)", srcSelector, dstSelector)
				}
				if _, exists := tokensToReset[tokenAddress]; exists {
					return fmt.Errorf(
						"the same address cannot be referenced in both TokensToUseDefaultFeeConfigs and TokenTransferFeeConfigArgs (src = %d, dst = %d, addr = %s)",
						srcSelector,
						dstSelector,
						tokenAddress.Hex(),
					)
				}
			}
		}
	}

	return nil
}

func setTokenTransferFeeConfigLogic(env cldf.Environment, cfg SetTokenTransferFeeConfig) (cldf.ChangesetOutput, error) {
	if len(cfg.InputsByChain) == 0 {
		env.Logger.Warn("no inputs were provided - exiting apply stage gracefully")
		return cldf.ChangesetOutput{}, nil
	}

	state, err := stateview.LoadOnchainState(env, stateview.WithLoadLegacyContracts(true))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.
		NewDeployerGroup(env, state, cfg.MCMS).
		WithDeploymentContext("SetTokenTransferFeeConfig")

	env.Logger.Info("preparing deployer group transactions")
	for srcSelector, inputs := range cfg.InputsByChain {
		env.Logger.Infof("processing src %d", srcSelector)
		if len(inputs) == 0 {
			env.Logger.Infof("no inputs were detected for src %d - skipping", srcSelector)
			continue
		}

		srcChain, exists := env.BlockChains.EVMChains()[srcSelector]
		if !exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("could not find src EVM chain in environment (src = %d)", srcSelector)
		}

		chainState, exists := state.Chains[srcSelector]
		if !exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("could not find chain in state (src = %s)", srcChain.String())
		}

		opts, err := deployerGroup.GetDeployer(srcSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer (src = %s): %w", srcChain.String(), err)
		}

		for dstSelector, input := range inputs {
			dstChain, exists := env.BlockChains.EVMChains()[dstSelector]
			if !exists {
				return cldf.ChangesetOutput{}, fmt.Errorf("could not find dst EVM chain in environment (src = %s, dst = %d)", srcChain.String(), dstSelector)
			}

			onramp, exists := chainState.EVM2EVMOnRamp[dstSelector]
			if !exists {
				return cldf.ChangesetOutput{}, fmt.Errorf("no EVM2EVMOnRamp (src = %s, dst = %s)", srcChain.String(), dstChain.String())
			}

			tokenTransferFeeConfigArgs := []evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{}
			for tokenAddress, args := range input.TokenTransferFeeConfigArgs {
				// This gets the token transfer fee config for the given token - if it doesn't exist, then the zero struct will be returned and `IsEnabled` will be `false`
				env.Logger.Infof("fetching token transfer fee config (src = %s, dst = %s, token = %s)", srcChain.String(), dstChain.String(), tokenAddress.Hex())
				curConfig, err := onramp.GetTokenTransferFeeConfig(&bind.CallOpts{Context: env.GetContext()}, tokenAddress)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch token transfer fee config (src = %s, dst = %s, token = %s): %w", srcChain.String(), dstChain.String(), tokenAddress.Hex(), err)
				}

				// If no custom config already exists on-chain for the token, then we have no fallback values to use - in this case the caller must explicitly provide all fields
				env.Logger.Infof("fetched token transfer fee config (src = %s, dst = %s, token = %s, cfg = %+v)", srcChain.String(), dstChain.String(), tokenAddress.Hex(), curConfig)
				if !curConfig.IsEnabled && args.HasMissingFields() {
					return cldf.ChangesetOutput{}, fmt.Errorf("invalid args - when enabling a new token, all fields must be provided (src = %s, dst = %s, token = %s)", srcChain.String(), dstChain.String(), tokenAddress.Hex())
				}

				// At this point, we're either using fallback values from the chain or the caller has explicitly provided the inputs
				newConfig := evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{
					Token:                     tokenAddress,
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
				env.Logger.Infof("constructed token transfer fee config (src = %s, dst = %s, token = %s, new_cfg = %+v)", srcChain.String(), dstChain.String(), tokenAddress.Hex(), newConfig)
				if isDifferent {
					tokenTransferFeeConfigArgs = append(tokenTransferFeeConfigArgs, newConfig)
				} else {
					env.Logger.Infof("skipping update since input config is the same as on-chain config (src = %s, dst = %s, token = %s, cfg = %+v)", srcChain.String(), dstChain.String(), tokenAddress.Hex(), curConfig)
				}
			}

			resetsCount := len(input.TokensToUseDefaultFeeConfigs)
			updateCount := len(tokenTransferFeeConfigArgs)
			if updateCount == 0 && resetsCount == 0 {
				env.Logger.Infof("no changes detected (src = %s, dst = %s) - skipping", srcChain.String(), dstChain.String())
				continue
			}

			env.Logger.Infof("setting token transfer fee configs (src = %s, dst = %s, updates = %d, resets = %d)", srcChain.String(), dstChain.String(), updateCount, resetsCount)
			_, err = onramp.SetTokenTransferFeeConfig(opts,
				tokenTransferFeeConfigArgs,
				input.TokensToUseDefaultFeeConfigs,
			)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf(
					"failed to create SetTokenTransferFeeConfig transaction (src = %s, dst = %s): %w",
					srcChain.String(),
					dstChain.String(),
					err,
				)
			}
		}
	}

	env.Logger.Info("running deployer group")
	return deployerGroup.Enact()
}
