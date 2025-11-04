package crossfamily

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/helpers/pointer"

	ccip_cs_common "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccip_cs_sol_v0_1_1 "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	ccip_cs_evm_v1_5_1 "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	ccip_cs_evm_v1_6_0 "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSetV2[SetTokenTransferFeeConfigInput] = SetTokenTransferFeeConfig

// SetTokenTransferFeeConfig is a cross-family changeset that updates per-token transfer-fee
// parameters for both EVM and Solana source chains in one operation. It accepts a two-level
// map of source -> destination chain selectors with per-token configs, validates each entry,
// and routes them to the appropriate version-specific changeset.
//
// Features:
//   - Version-aware: automatically detects or respects VersionHints for each family (EVM v1.6.0 / v1.5.1 and Solana v0.1.1 supported).
//   - Unified schema: validates and converts string token addresses (hex → EVM, base58 → Solana).
//   - Strict validation: rejects invalid selectors and same-chain (src==dst) updates.
//   - No-op friendly: skips empty updates and exits cleanly when no changes are needed.
//   - MCMS integration: batches all family/version changes into a single orchestrated proposal.
//
// This changeset performs only validation, version routing, and type conversion. It delegates
// all actual on-chain logic to the existing family/version changesets.
var SetTokenTransferFeeConfig = cldf.CreateChangeSet(setTokenTransferFeeLogic, setTokenTransferFeePrecondition)

var SetTokenTransferFeeLatestSupportedVersions = OptionalVersions{
	Solana: pointer.To(ccip_cs_sol_v0_1_1.VersionSolanaV0_1_1),
	Evm:    pointer.To(deployment.Version1_6_0.String()),
}

type OptionalVersions struct {
	Solana *string
	Evm    *string
}

type OptionalTokenTransferFeeConfig struct {
	DestBytesOverhead *uint32
	DestGasOverhead   *uint32
	MinFeeUsdCents    *uint32
	MaxFeeUsdCents    *uint32
	DeciBps           *uint16
	IsEnabled         *bool
}

type TokenTransferFeeConfigArgs struct {
	TokenAddressToFeeConfig map[string]OptionalTokenTransferFeeConfig
}

type SetTokenTransferFeeConfigInput struct {
	InputsByChain map[uint64]map[uint64]TokenTransferFeeConfigArgs `json:"inputsByChain"`
	VersionHints  *OptionalVersions                                `json:"versionHints"`
	MCMS          *proposalutils.TimelockConfig                    `json:"mcms"`
}

func (cfg SetTokenTransferFeeConfigInput) buildOrchestrateChangesetsConfig(env cldf.Environment) (ccip_cs_common.OrchestrateChangesetsConfig, error) {
	// Make sure MCMS config is specified
	env.Logger.Info("building orchestrate changesets config")
	if cfg.MCMS == nil {
		return ccip_cs_common.OrchestrateChangesetsConfig{}, errors.New("MCMS config is required")
	}

	// Load state
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return ccip_cs_common.OrchestrateChangesetsConfig{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	// Default to latest supported version for each chain family if no version hints were specified
	allVersionHints := pointer.Coalesce(cfg.VersionHints, SetTokenTransferFeeLatestSupportedVersions)
	solanaVersion := pointer.Coalesce(allVersionHints.Solana, *SetTokenTransferFeeLatestSupportedVersions.Solana)
	evmVersion := pointer.Coalesce(allVersionHints.Evm, *SetTokenTransferFeeLatestSupportedVersions.Evm)

	// Define the configs for each version
	solanaConfigV0_1_1 := map[uint64]map[uint64]ccip_cs_sol_v0_1_1.TokenTransferFeeForRemoteChainConfigArgsV2{}
	evmConfigV1_6_0 := map[uint64]map[uint64]ccip_cs_evm_v1_6_0.ApplyTokenTransferFeeConfigUpdatesConfigV2Input{}
	evmConfigV1_5_1 := map[uint64]map[uint64]ccip_cs_evm_v1_5_1.SetTokenTransferFeeArgs{}

	// Populate configs
	for srcSelector, inputs := range cfg.InputsByChain {
		err := stateview.ValidateChain(env, state, srcSelector, cfg.MCMS)
		if err != nil {
			return ccip_cs_common.OrchestrateChangesetsConfig{}, fmt.Errorf("failed to validate src chain (src = %d): %w", srcSelector, err)
		}
		srcFamily, err := chainsel.GetSelectorFamily(srcSelector)
		if err != nil {
			return ccip_cs_common.OrchestrateChangesetsConfig{}, fmt.Errorf("failed to get src chain family (src = %d): %w", srcSelector, err)
		}

		for dstSelector, input := range inputs {
			if err := stateview.ValidateChain(env, state, dstSelector, cfg.MCMS); err != nil {
				return ccip_cs_common.OrchestrateChangesetsConfig{}, fmt.Errorf(
					"failed to validate dst chain (src = %d, dst = %d): %w", srcSelector, dstSelector, err)
			}
			if srcSelector == dstSelector {
				return ccip_cs_common.OrchestrateChangesetsConfig{}, fmt.Errorf(
					"destination chain cannot be the same as source chain (src = %d, dst = %d)", srcSelector, dstSelector)
			}

			// Use the chain family and version to determine which config should be updated
			switch srcFamily {
			case chainsel.FamilyEVM:
				switch evmVersion {
				case deployment.Version1_6_0.String():
					if err = updateEvmTokenTransferFeeConfigV1_6_0(srcSelector, dstSelector, input, evmConfigV1_6_0); err != nil {
						return ccip_cs_common.OrchestrateChangesetsConfig{}, fmt.Errorf("failed to update EVM %s config: %w", evmVersion, err)
					}
				case deployment.Version1_5_1.String():
					if err = updateEvmTokenTransferFeeConfigV1_5_1(srcSelector, dstSelector, input, evmConfigV1_5_1); err != nil {
						return ccip_cs_common.OrchestrateChangesetsConfig{}, fmt.Errorf("failed to update EVM %s config: %w", evmVersion, err)
					}
				default:
					return ccip_cs_common.OrchestrateChangesetsConfig{}, fmt.Errorf("unsupported EVM version: %s", evmVersion)
				}
			case chainsel.FamilySolana:
				switch solanaVersion {
				case ccip_cs_sol_v0_1_1.VersionSolanaV0_1_1:
					if err = updateSolanaTokenTransferFeeConfigV0_1_1(srcSelector, dstSelector, input, solanaConfigV0_1_1); err != nil {
						return ccip_cs_common.OrchestrateChangesetsConfig{}, fmt.Errorf("failed to update Solana %s config: %w", solanaVersion, err)
					}
				default:
					return ccip_cs_common.OrchestrateChangesetsConfig{}, fmt.Errorf("unsupported Solana version: %s", solanaVersion)
				}
			default:
				return ccip_cs_common.OrchestrateChangesetsConfig{}, fmt.Errorf("unsupported source family for selector %d: %s", srcSelector, srcFamily)
			}
		}
	}

	// Combine the results
	changesets := []ccip_cs_common.WithConfig{}
	if len(evmConfigV1_6_0) > 0 {
		changesets = append(changesets, ccip_cs_common.CreateGenericChangeSetWithConfig(
			ccip_cs_evm_v1_6_0.ApplyTokenTransferFeeConfigUpdatesFeeQuoterChangesetV2,
			ccip_cs_evm_v1_6_0.ApplyTokenTransferFeeConfigUpdatesConfigV2{
				InputsByChain: evmConfigV1_6_0,
				MCMS:          cfg.MCMS,
			},
		))
	}
	if len(evmConfigV1_5_1) > 0 {
		changesets = append(changesets, ccip_cs_common.CreateGenericChangeSetWithConfig(
			ccip_cs_evm_v1_5_1.SetTokenTransferFeeConfigChangeset,
			ccip_cs_evm_v1_5_1.SetTokenTransferFeeConfig{
				InputsByChain: evmConfigV1_5_1,
				MCMS:          cfg.MCMS,
			},
		))
	}
	if len(solanaConfigV0_1_1) > 0 {
		changesets = append(changesets, ccip_cs_common.CreateGenericChangeSetWithConfig(
			ccip_cs_sol_v0_1_1.AddTokenTransferFeeForRemoteChainV2,
			ccip_cs_sol_v0_1_1.TokenTransferFeeForRemoteChainConfigV2{
				InputsByChain: solanaConfigV0_1_1,
				MCMS:          cfg.MCMS,
			},
		))
	}

	// Return the combined orchestrate config
	return ccip_cs_common.OrchestrateChangesetsConfig{
		Description: "Set Token Transfer Fee Config - Cross Family",
		MCMS:        cfg.MCMS,
		ChangeSets:  changesets,
	}, nil
}

func setTokenTransferFeePrecondition(env cldf.Environment, cfg SetTokenTransferFeeConfigInput) error {
	if len(cfg.InputsByChain) == 0 {
		env.Logger.Info("no inputs were provided - exiting precondition stage gracefully")
		return nil
	}

	input, err := cfg.buildOrchestrateChangesetsConfig(env)
	if err != nil {
		return fmt.Errorf("failed to build OrchestrateChangesetsConfig: %w", err)
	}

	if len(input.ChangeSets) == 0 {
		env.Logger.Info("no changesets to orchestrate - exiting precondition stage gracefully")
		return nil
	}

	return ccip_cs_common.OrchestrateChangesets.VerifyPreconditions(env, input)
}

func setTokenTransferFeeLogic(env cldf.Environment, cfg SetTokenTransferFeeConfigInput) (cldf.ChangesetOutput, error) {
	if len(cfg.InputsByChain) == 0 {
		env.Logger.Info("no inputs were provided - exiting apply stage gracefully")
		return cldf.ChangesetOutput{}, nil
	}

	input, err := cfg.buildOrchestrateChangesetsConfig(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build OrchestrateChangesetsConfig: %w", err)
	}

	if len(input.ChangeSets) == 0 {
		env.Logger.Info("no changesets to orchestrate - exiting apply stage gracefully")
		return cldf.ChangesetOutput{}, nil
	}

	return ccip_cs_common.OrchestrateChangesets.Apply(env, input)
}

// --- EVM helpers ----

func updateEvmTokenTransferFeeConfigV1_6_0(
	srcSelector uint64,
	dstSelector uint64,
	input TokenTransferFeeConfigArgs,
	config map[uint64]map[uint64]ccip_cs_evm_v1_6_0.ApplyTokenTransferFeeConfigUpdatesConfigV2Input,
) error {
	configs := map[common.Address]ccip_cs_evm_v1_6_0.OptionalFeeQuoterTokenTransferFeeConfig{}
	for rawTokenAddress, tokenFeeConfig := range input.TokenAddressToFeeConfig {
		if !common.IsHexAddress(rawTokenAddress) {
			return fmt.Errorf("invalid hex EVM address detected '%s'", rawTokenAddress)
		}
		configs[common.HexToAddress(rawTokenAddress)] = ccip_cs_evm_v1_6_0.OptionalFeeQuoterTokenTransferFeeConfig{
			DestBytesOverhead: tokenFeeConfig.DestBytesOverhead,
			DestGasOverhead:   tokenFeeConfig.DestGasOverhead,
			MinFeeUSDCents:    tokenFeeConfig.MinFeeUsdCents,
			MaxFeeUSDCents:    tokenFeeConfig.MaxFeeUsdCents,
			IsEnabled:         tokenFeeConfig.IsEnabled,
			DeciBps:           tokenFeeConfig.DeciBps,
		}
	}

	if len(configs) > 0 {
		if _, ok := config[srcSelector]; !ok {
			config[srcSelector] = map[uint64]ccip_cs_evm_v1_6_0.ApplyTokenTransferFeeConfigUpdatesConfigV2Input{}
		}
		config[srcSelector][dstSelector] = ccip_cs_evm_v1_6_0.ApplyTokenTransferFeeConfigUpdatesConfigV2Input{
			TokenTransferFeeConfigRemoveArgs: []common.Address{},
			TokenTransferFeeConfigArgs:       configs,
		}
	}
	return nil
}

func updateEvmTokenTransferFeeConfigV1_5_1(
	srcSelector uint64,
	dstSelector uint64,
	input TokenTransferFeeConfigArgs,
	config map[uint64]map[uint64]ccip_cs_evm_v1_5_1.SetTokenTransferFeeArgs,
) error {
	configs := map[common.Address]ccip_cs_evm_v1_5_1.TokenTransferFeeArgs{}
	for rawTokenAddress, tokenFeeConfig := range input.TokenAddressToFeeConfig {
		if !common.IsHexAddress(rawTokenAddress) {
			return fmt.Errorf("invalid hex EVM address detected '%s'", rawTokenAddress)
		}
		configs[common.HexToAddress(rawTokenAddress)] = ccip_cs_evm_v1_5_1.TokenTransferFeeArgs{
			DestBytesOverhead: tokenFeeConfig.DestBytesOverhead,
			DestGasOverhead:   tokenFeeConfig.DestGasOverhead,
			MinFeeUSDCents:    tokenFeeConfig.MinFeeUsdCents,
			MaxFeeUSDCents:    tokenFeeConfig.MaxFeeUsdCents,
			DeciBps:           tokenFeeConfig.DeciBps,

			// The sensible default should always be used for this field
			AggregateRateLimitEnabled: nil,
		}
	}

	if len(configs) > 0 {
		if _, ok := config[srcSelector]; !ok {
			config[srcSelector] = map[uint64]ccip_cs_evm_v1_5_1.SetTokenTransferFeeArgs{}
		}
		config[srcSelector][dstSelector] = ccip_cs_evm_v1_5_1.SetTokenTransferFeeArgs{
			TokensToUseDefaultFeeConfigs: []common.Address{},
			TokenTransferFeeConfigArgs:   configs,
		}
	}
	return nil
}

// --- Solana helpers ---

func updateSolanaTokenTransferFeeConfigV0_1_1(
	srcSelector uint64,
	dstSelector uint64,
	input TokenTransferFeeConfigArgs,
	config map[uint64]map[uint64]ccip_cs_sol_v0_1_1.TokenTransferFeeForRemoteChainConfigArgsV2,
) error {
	configs := map[solana.PublicKey]ccip_cs_sol_v0_1_1.OptionalFeeQuoterTokenTransferFeeConfig{}
	for rawTokenAddress, tokenFeeConfig := range input.TokenAddressToFeeConfig {
		tokenAddress, err := solana.PublicKeyFromBase58(rawTokenAddress)
		if err != nil {
			return fmt.Errorf("invalid base58 solana address detected '%s': %w", rawTokenAddress, err)
		}
		configs[tokenAddress] = ccip_cs_sol_v0_1_1.OptionalFeeQuoterTokenTransferFeeConfig{
			DestBytesOverhead: tokenFeeConfig.DestBytesOverhead,
			DestGasOverhead:   tokenFeeConfig.DestGasOverhead,
			MinFeeUsdcents:    tokenFeeConfig.MinFeeUsdCents,
			MaxFeeUsdcents:    tokenFeeConfig.MaxFeeUsdCents,
			IsEnabled:         tokenFeeConfig.IsEnabled,
			DeciBps:           tokenFeeConfig.DeciBps,
		}
	}

	if len(configs) > 0 {
		if _, ok := config[srcSelector]; !ok {
			config[srcSelector] = map[uint64]ccip_cs_sol_v0_1_1.TokenTransferFeeForRemoteChainConfigArgsV2{}
		}
		config[srcSelector][dstSelector] = ccip_cs_sol_v0_1_1.TokenTransferFeeForRemoteChainConfigArgsV2{
			TokenAddressToFeeConfig: configs,
		}
	}
	return nil
}
