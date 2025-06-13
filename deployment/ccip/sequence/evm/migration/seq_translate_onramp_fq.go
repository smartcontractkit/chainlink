package migration

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/migration"
)

type OnRampToFeeQuoterDestChainConfigDeps struct {
	Chain cldf_evm.Chain
}

type OnRampToFeeQuoterDestChainConfigInput struct {
	OnRamps             map[uint64]common.Address
	TokenAdminRegistry  common.Address
	SourceChainSelector uint64
}

type OnRampToFeeQuoterDestChainConfigOutput struct {
	FeeQuoterUpdates map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig
	Feetokens        []common.Address
}

var (
	SeqTranslateOnRampToFeeQDestConfig = operations.NewSequence(
		"translate-on-ramp-to-feequoter-dest-config",
		semver.MustParse("1.0.0"),
		"Translates existing 1.5.0 EVM2EVMOnRamp configs into appropriate 1.6.0 FeeQuoter Destination configs & returns all supported Fee tokens",
		func(b operations.Bundle, deps OnRampToFeeQuoterDestChainConfigDeps, input OnRampToFeeQuoterDestChainConfigInput) (OnRampToFeeQuoterDestChainConfigOutput, error) {
			feeQuoterUpdates := make(map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig)
			allFeeTokens := make([]common.Address, 0)
			lggr := b.Logger
			srcChain := deps.Chain

			for destChainSel, onRamp1_5 := range input.OnRamps {
				lggr.Infow("Processing 1.5.0 OnRamps", "sourceChainSelector", input.SourceChainSelector, "chainName", srcChain.Name, "destinationChainSelector", destChainSel, "onRampAddress", onRamp1_5)

				feeQuoterDestChainConfig := fee_quoter.FeeQuoterDestChainConfig{}
				feeQuoterTranslatedDestCfg := EVM2EVMOnRampMigrateDestChainConfig{FeeQuoterDestChainConfig: feeQuoterDestChainConfig}

				evm2evmOnRampDynamicCfgReport, err := operations.ExecuteOperation(
					b,
					ops.EVM2EVMOnrampGetDynamicCfgOp,
					ops.MigrateOnRampToFQDeps{
						Chain: srcChain,
					},
					ops.MigrateDestChainCfgInput{
						OnRamp:            onRamp1_5,
						DestChainSelector: destChainSel,
					},
				)
				if err != nil {
					return OnRampToFeeQuoterDestChainConfigOutput{}, fmt.Errorf("Failed to execute TranslateOnRampToFQDestDynamicCfgOps: %w", err)
				}

				feeQuoterTranslatedDestCfg.TranslateOnrampToFeequoterDynamicConfig(destChainSel, evm2evmOnRampDynamicCfgReport.Output)

				allFeeTokensOp, err := operations.ExecuteOperation(
					b, ops.PriceRegistryGetAllFeeTokensOps,
					ops.MigrateOnRampToFQDeps{
						Chain: srcChain,
					},
					evm2evmOnRampDynamicCfgReport.Output.PriceRegistry,
				)
				if err != nil {
					return OnRampToFeeQuoterDestChainConfigOutput{}, fmt.Errorf("failed to Execute GetAllFeeTokensOps: %w", err)
				}

				// add supported fee token config to FeeQuoter

				// This is per token in 1.5.0 onRamp, but in FeeQuoter its per destination chain,
				// But from RDD the values are constant (either 10/50 regardless of FeeToken)
				// So we can just use the first token's config
				// for _, ft := range allFeeTokensOp.Output {
				feetokenCfgReport, err := operations.ExecuteOperation(
					b, ops.EVM2EVMOnrampGetFeeTokenConfigOp,
					ops.MigrateOnRampToFQDeps{
						Chain: srcChain,
					},
					ops.OnRampGetFeeTokenCfgInput{
						OnRamp:          onRamp1_5,
						FeeTokenAddress: allFeeTokensOp.Output[0],
					},
				)
				if err != nil {
					return OnRampToFeeQuoterDestChainConfigOutput{}, fmt.Errorf("failed to Execute GetOnRampGetFeeTokenConfigOps: %w", err)
				}
				// }
				feeQuoterTranslatedDestCfg.TranslateOnrampToFeequoterFeeTokenCfg(feetokenCfgReport.Output)

				if _, ok := feeQuoterUpdates[input.SourceChainSelector]; !ok {
					feeQuoterUpdates[input.SourceChainSelector] = make(map[uint64]fee_quoter.FeeQuoterDestChainConfig)
				}
				feeQuoterUpdates[input.SourceChainSelector][destChainSel] = feeQuoterTranslatedDestCfg.FeeQuoterDestChainConfig
			}

			return OnRampToFeeQuoterDestChainConfigOutput{
				FeeQuoterUpdates: feeQuoterUpdates,
				Feetokens:        allFeeTokens,
			}, nil
		})

	SeqTranslateOnRampToFeeQTokenTransferFeeCfg = operations.NewSequence(
		"translate-on-ramp-to-feeQuoter-token-transfer-fee-configs",
		semver.MustParse("1.0.0"),
		"Translates existing 1.5.0 EVM2EVMOnRamp Token Transfer Fee Configs into appropriate 1.6.0 FeeQuoter Destination configs",
		func(b operations.Bundle, deps OnRampToFeeQuoterDestChainConfigDeps, input OnRampToFeeQuoterDestChainConfigInput) (map[uint64]fee_quoter.FeeQuoterTokenTransferFeeConfigArgs, error) {
			lggr := b.Logger
			tokenTransferFeeConfigsPerSrcChain := make(map[uint64]fee_quoter.FeeQuoterTokenTransferFeeConfigArgs)
			var tokenTransferFeeConfigsPerDestChain fee_quoter.FeeQuoterTokenTransferFeeConfigArgs

			for destChainSel, onRamp1_5 := range input.OnRamps {
				lggr.Infow("Processing 1.5.0 OnRamps", "sourceChainSelector", input.SourceChainSelector, "chainName", deps.Chain.Name, "destinationChainSelector", destChainSel, "onRampAddress", onRamp1_5)

				onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(onRamp1_5, deps.Chain.Client)
				if err != nil {
					return tokenTransferFeeConfigsPerSrcChain, fmt.Errorf("Failed to create EVM2EVMOnRamp contract binding at address %s: %w", onRamp1_5, err)
				}

				migrateOnRamp := EVM2EVMOnRampMigrate{EVM2EVMOnRamp: onRamp}
				allTransferTokensAndCfgs := make([]fee_quoter.FeeQuoterTokenTransferFeeConfigSingleTokenArgs, 0)

				// Port token transfer fee config args from all 1.5.0 OnRamps into FeeQuoter
				//	-> get alltokens from tokenAdminReg.getAllConfiguredTokens
				// 	-> for each token
				// 		-->  get tokenpool on onRamp.getPoolBySourceToken(selector, token)
				//		--> supportedChains := tokenPool.getSupportedChains()
				//		--> if targetSelector not in supportedChains then continue
				// 		--> add this token to validTokens to process
				getAllConfiguredTokensOps, err := operations.ExecuteOperation(
					b, ops.TokenAdminRegistryGetAllConfiguredTokensOp,
					ops.MigrateOnRampToFQDeps{
						Chain: deps.Chain,
					},
					input.TokenAdminRegistry,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to get all configured tokens from TokenAdminRegistry on source chain %d: %w", input.SourceChainSelector, err)
				}
				allTokens := getAllConfiguredTokensOps.Output
				for _, token := range allTokens {
					getPoolBySourceTokenOps, err := operations.ExecuteOperation(
						b, ops.EVM2EVMOnrampGetPoolBySourceTokenOp,
						ops.MigrateOnRampToFQDeps{
							Chain: deps.Chain,
						},
						ops.GetPoolBySourceTokenInput{
							OnRamp:            onRamp1_5,
							DestChainSelector: destChainSel,
							FeeTokenAddress:   token,
						},
					)
					if err != nil {
						return nil, fmt.Errorf("failed to get all configured tokens from TokenAdminRegistry on source chain %d: %w", input.SourceChainSelector, err)
					}
					if getPoolBySourceTokenOps.Output == (common.Address{}) {
						lggr.Warnw("Failed to get pool for token on 1.5.0 OnRamp", "sourceChainSelector", input.SourceChainSelector, "destinationChainSelector", destChainSel, "token", token.Hex(), "error", err)
						continue // continue or exit?
					}

					getSupportedChainsForTokenPool, err := operations.ExecuteOperation(
						b, ops.TokenPoolGetSupportedChainsOp,
						ops.MigrateOnRampToFQDeps{
							Chain: deps.Chain,
						},
						getPoolBySourceTokenOps.Output,
					)
					if err != nil {
						return nil, fmt.Errorf("failed to get suuported chains for the toksn Pool on source chain %d: %w", input.SourceChainSelector, err)
					}

					// Check if the destination chain selector is in the supported chains of token pool
					found := false
					for _, chainSel := range getSupportedChainsForTokenPool.Output {
						if chainSel == destChainSel {
							found = true
							break
						}
					}
					if !found {
						continue // skip this token if the destination chain is not supported
					}
					tokenTransferFeeCfg, err2 := onRamp.GetTokenTransferFeeConfig(nil, token)
					if err2 != nil {
						return nil, fmt.Errorf("failed to get suported chains for the toksn Pool on source chain %d: %w", input.SourceChainSelector, err)
					}
					allTransferTokensAndCfgs = append(allTransferTokensAndCfgs,
						migrateOnRamp.TranslateOnrampToFeequoterTokenTransferFeeConfig(destChainSel, token, tokenTransferFeeCfg),
					)
				}
				tokenTransferFeeConfigsPerDestChain = fee_quoter.FeeQuoterTokenTransferFeeConfigArgs{
					DestChainSelector:       destChainSel,
					TokenTransferFeeConfigs: allTransferTokensAndCfgs,
				}
			}

			if tokenTransferFeeConfigsPerDestChain.DestChainSelector != 0 {
				tokenTransferFeeConfigsPerSrcChain[input.SourceChainSelector] = tokenTransferFeeConfigsPerDestChain
			} else {
				return nil, fmt.Errorf("Empty token transfer fee configs for chain %d", input.SourceChainSelector)
			}
			return tokenTransferFeeConfigsPerSrcChain, nil
		})
)
