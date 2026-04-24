package migration

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	migration_ops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_5"

	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

type OnRampToFeeQuoterDestChainConfigInput struct {
	OnRamps            map[uint64]common.Address
	NewFeeQuoterParams map[uint64]NewFeeQuoterDestChainConfigParams
	TokenAdminRegistry common.Address
}

type FeeQuoterUpdateTokenTransferConfig struct {
	UpdatesByChain map[uint64]opsutil.EVMCallInput[OnRampToFeeQuoterDestChainConfigInput]
}

type OnRampToFeeQuoterDestChainConfigOutput struct {
	FeeQuoterUpdates           map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig
	FeeTokens                  map[uint64][]common.Address
	FeeTokenPremiumMultipliers map[uint64][]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs
}

type OnRampToFeeQuoterTokenTransferFeeCfgOutput struct {
	FeeQuoterUpdates map[uint64][]fee_quoter.FeeQuoterTokenTransferFeeConfigArgs
}

var (
	SeqTranslateOnRampToFeeQDestConfig = operations.NewSequence(
		"translate-on-ramp-to-feequoter-dest-config",
		semver.MustParse("1.0.0"),
		"Translates existing 1.5.0 EVM2EVMOnRamp configs into appropriate 1.6.0 FeeQuoter Destination configs & returns all supported Fee tokens",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input FeeQuoterUpdateTokenTransferConfig) (OnRampToFeeQuoterDestChainConfigOutput, error) {
			feeQuoterUpdates := make(map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig)
			allFeeTokens := make(map[uint64][]common.Address)
			allFeetokenPremiumMultipliers := make(map[uint64][]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs)
			lggr := b.Logger
			for chainSel, update := range input.UpdatesByChain {
				srcChain, ok := chains[chainSel]
				if !ok {
					return OnRampToFeeQuoterDestChainConfigOutput{}, fmt.Errorf("chain with selector %d not defined", chainSel)
				}
				for destChainSel, onRamp1_5 := range update.CallInput.OnRamps {
					lggr.Infow("Processing 1.5.0 OnRamps", "sourceChainSelector", srcChain.ChainSelector(), "chainName", srcChain.Name, "destinationChainSelector", destChainSel, "onRampAddress", onRamp1_5)

					feeQuoterDestChainConfig := fee_quoter.FeeQuoterDestChainConfig{}
					feeQuoterTranslatedDestCfg := EVM2EVMOnRampMigrateDestChainConfig{FeeQuoterDestChainConfig: feeQuoterDestChainConfig}

					evm2evmOnRampDynamicCfgReport, err := operations.ExecuteOperation(
						b, migration_ops.EVM2EVMOnrampGetDynamicCfgOp,
						migration_ops.MigrateOnRampToFQDeps{
							Chain: srcChain,
						},
						onRamp1_5,
					)
					if err != nil {
						return OnRampToFeeQuoterDestChainConfigOutput{}, fmt.Errorf("failed to execute TranslateOnRampToFQDestDynamicCfgOps: %w", err)
					}

					feeQuoterTranslatedDestCfg.TranslateOnrampToFeequoterDynamicConfig(destChainSel, evm2evmOnRampDynamicCfgReport.Output, update.CallInput.NewFeeQuoterParams[destChainSel])

					allFeeTokensOp, err := operations.ExecuteOperation(
						b, migration_ops.PriceRegistryGetAllFeeTokensOps,
						migration_ops.MigrateOnRampToFQDeps{
							Chain: srcChain,
						},
						migration_ops.PriceRegistryGetAllFeeTokensIn{
							Address:       evm2evmOnRampDynamicCfgReport.Output.PriceRegistry,
							ChainSelector: chainSel,
						},
					)
					if err != nil {
						return OnRampToFeeQuoterDestChainConfigOutput{}, fmt.Errorf("failed to Execute GetAllFeeTokensOps: %w", err)
					}

					// add supported fee token config to FeeQuoter

					// This is per token in 1.5.0 onRamp, but in FeeQuoter its per destination chain,
					// But RDD values are just redundant & can be adjusted by the premium multiplier, so simplified in 1.6 FQ
					// So we can just use the any token's config (the last one in the loop here)
					onRampFeeTokenCfgReport := evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfig{}

					feeTokenPremiumMultipliers := make([]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs, len(allFeeTokensOp.Output))
					for idx, ft := range allFeeTokensOp.Output {
						feetokenCfgReport, err := operations.ExecuteOperation(
							b, migration_ops.EVM2EVMOnrampGetFeeTokenConfigOp,
							migration_ops.MigrateOnRampToFQDeps{
								Chain: srcChain,
							},
							migration_ops.OnRampGetTokenCfgIn{
								OnRamp:        onRamp1_5,
								Address:       ft,
								ChainSelector: chainSel,
							},
						)
						if err != nil {
							return OnRampToFeeQuoterDestChainConfigOutput{}, fmt.Errorf("failed to Execute GetOnRampGetFeeTokenConfigOps: %w", err)
						}

						// Translate the feeToken PremiumMultiplierCfg to 1.6 FeeQuoter config
						premiumMultiplierCfg := EVM2EVMOnRampMigratePremiumMultiplierCfg{}
						premiumMultiplierCfg.TranslateOnrampToFeeQFeePremiumCfg(ft, feetokenCfgReport.Output)
						feeTokenPremiumMultipliers[idx] = premiumMultiplierCfg.FeeQuoterPremiumMultiplierWeiPerEthArgs
						if onRampFeeTokenCfgReport == (evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfig{}) {
							onRampFeeTokenCfgReport = feetokenCfgReport.Output
						}
					}

					if _, ok := feeQuoterUpdates[chainSel]; !ok {
						feeQuoterUpdates[chainSel] = make(map[uint64]fee_quoter.FeeQuoterDestChainConfig)
					}
					feeQuoterUpdates[chainSel][destChainSel] = feeQuoterTranslatedDestCfg.FeeQuoterDestChainConfig
					allFeeTokens[chainSel] = append(allFeeTokens[chainSel], allFeeTokensOp.Output...)
					allFeetokenPremiumMultipliers[chainSel] = append(allFeetokenPremiumMultipliers[chainSel], feeTokenPremiumMultipliers...)
				}
			}

			return OnRampToFeeQuoterDestChainConfigOutput{
				FeeQuoterUpdates:           feeQuoterUpdates,
				FeeTokens:                  allFeeTokens,
				FeeTokenPremiumMultipliers: allFeetokenPremiumMultipliers,
			}, nil
		})

	SeqTranslateOnRampToFeeQTokenTransferFeeCfg = operations.NewSequence(
		"translate-on-ramp-to-feeQuoter-token-transfer-fee-configs",
		semver.MustParse("1.0.0"),
		"Translates existing 1.5.0 EVM2EVMOnRamp Token Transfer Fee Configs into appropriate 1.6.0 FeeQuoter Destination configs",
		func(b operations.Bundle, chains map[uint64]cldf_evm.Chain, input FeeQuoterUpdateTokenTransferConfig) (OnRampToFeeQuoterTokenTransferFeeCfgOutput, error) {
			lggr := b.Logger
			tokenTransferFeeConfigsPerSrcChain := make(map[uint64][]fee_quoter.FeeQuoterTokenTransferFeeConfigArgs)

			for chainSel, update := range input.UpdatesByChain {
				srcChain, ok := chains[chainSel]
				var tokenTransferFeeConfigsPerDestChain []fee_quoter.FeeQuoterTokenTransferFeeConfigArgs
				if !ok {
					return OnRampToFeeQuoterTokenTransferFeeCfgOutput{}, fmt.Errorf("chain with selector %d not defined", chainSel)
				}
				for destChainSel, onRamp1_5 := range update.CallInput.OnRamps {
					lggr.Infow("Processing 1.5.0 OnRamps", "sourceChainSelector", srcChain.ChainSelector(), "chainName", srcChain.Name, "destinationChainSelector", destChainSel, "onRampAddress", onRamp1_5)
					onRamp := evm_2_evm_onramp.EVM2EVMOnRamp{}
					migrateOnRamp := EVM2EVMOnRampMigrate{EVM2EVMOnRamp: &onRamp}
					allTransferTokensAndCfgs := make([]fee_quoter.FeeQuoterTokenTransferFeeConfigSingleTokenArgs, 0)

					// Port token transfer fee config args from all 1.5.0 OnRamps into FeeQuoter
					//	-> get alltokens from tokenAdminReg.getAllConfiguredTokens
					// 	-> for each token
					// 		-->  get tokenpool on onRamp.GetTokenTransferFeeConfig(token)
					//		--> if isEnabled
					// 		--> add this token to validTokens to process
					getAllConfiguredTokensOps, err := operations.ExecuteOperation(
						b, migration_ops.TokenAdminRegistryGetAllConfiguredTokensOp,
						migration_ops.MigrateOnRampToFQDeps{
							Chain: srcChain,
						},
						migration_ops.TokenAdminRegistryGetAllConfiguredTokensIn{
							Address:       update.CallInput.TokenAdminRegistry,
							ChainSelector: chainSel,
						},
					)
					if err != nil {
						return OnRampToFeeQuoterTokenTransferFeeCfgOutput{}, fmt.Errorf("failed to get all configured tokens from TokenAdminRegistry on source chain %d: %w", chainSel, err)
					}
					allTokens := getAllConfiguredTokensOps.Output
					for _, token := range allTokens {
						tokenTransferFeeCfgOp, err := operations.ExecuteOperation(
							b, migration_ops.EVM2EVMOnrampGetTokenTransferFeeConfigOp,
							migration_ops.MigrateOnRampToFQDeps{
								Chain: srcChain,
							},
							migration_ops.OnRampGetTokenCfgIn{
								OnRamp:        onRamp1_5,
								Address:       token,
								ChainSelector: chainSel,
							},
						)
						if err != nil {
							return OnRampToFeeQuoterTokenTransferFeeCfgOutput{}, fmt.Errorf("failed to get suported chains for the toksn Pool on source chain %d: %w", chainSel, err)
						}
						if !tokenTransferFeeCfgOp.Output.IsEnabled {
							continue // skip this token if the transfer fee config is not enabled
						}

						allTransferTokensAndCfgs = append(allTransferTokensAndCfgs,
							migrateOnRamp.TranslateOnrampToFeequoterTokenTransferFeeConfig(token, tokenTransferFeeCfgOp.Output),
						)
					}
					tokenTransferFeeConfigsPerDestChain = append(tokenTransferFeeConfigsPerDestChain, fee_quoter.FeeQuoterTokenTransferFeeConfigArgs{
						DestChainSelector:       destChainSel,
						TokenTransferFeeConfigs: allTransferTokensAndCfgs,
					})
				}
				tokenTransferFeeConfigsPerSrcChain[chainSel] = tokenTransferFeeConfigsPerDestChain
			}

			return OnRampToFeeQuoterTokenTransferFeeCfgOutput{
				FeeQuoterUpdates: tokenTransferFeeConfigsPerSrcChain,
			}, nil
		})
)
