package migration

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/price_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type MigrateOnRampToFQDeps struct {
	Chain cldf_evm.Chain
}

type MigrateDestChainCfgInput struct {
	OnRamp            common.Address
	DestChainSelector uint64
}

type MigrateDestDynamicCfgOutput struct {
	FeeQuoterDestChainConfig   fee_quoter.FeeQuoterDestChainConfig
	EVM2EVMOnRampDynamicConfig evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig
}

type OnRampGetFeeTokenCfgInput struct {
	OnRamp          common.Address
	FeeTokenAddress common.Address
}

type GetAllConfiguredTokensInput struct {
	TokenAdminRegistry common.Address
}

type GetPoolBySourceTokenInput struct {
	OnRamp            common.Address
	FeeTokenAddress   common.Address
	DestChainSelector uint64
}

var (
	TranslateDestDynamicCfgOps = operations.NewOperation( //todo: rename this
		"TranslateDestDynamicCfgOps",
		semver.MustParse("1.0.0"),
		"Ports dynamic config from all 1.5.0 OnRamps into FeeQuoter DestChain Dynamic Config",
		func(b operations.Bundle, deps MigrateOnRampToFQDeps, input MigrateDestChainCfgInput) (MigrateDestDynamicCfgOutput, error) {
			onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(input.OnRamp, deps.Chain.Client)
			if err != nil {
				return MigrateDestDynamicCfgOutput{}, fmt.Errorf("Failed to create EVM2EVMOnRamp contract binding at address %s: %w", input.OnRamp.Hex(), err)
			}
			// migrateOnRamp := migration.EVM2EVMOnRampMigrate{EVM2EVMOnRamp: onRamp}
			destChainEVM2EVMDynamicCfg, err2 := onRamp.GetDynamicConfig(nil)
			if err2 != nil && destChainEVM2EVMDynamicCfg.PriceRegistry == (common.Address{}) {
				return MigrateDestDynamicCfgOutput{}, fmt.Errorf("Cannot GetDynamicConfig for destination Chain: %d, for 1.5.0 OnRamp %s: %w", input.DestChainSelector, onRamp.Address().Hex(), err2)
			}

			// fqDestConfig := migrateOnRamp.TranslateOnrampToFeequoterDynamicConfig(input.DestChainSelector, destChainEVM2EVMDynamicCfg)
			return MigrateDestDynamicCfgOutput{
				// FeeQuoterDestChainConfig:   nil,
				EVM2EVMOnRampDynamicConfig: destChainEVM2EVMDynamicCfg,
			}, nil
		})

	/* 	TranslateOnRampToFQDestDynamicCfgOps = operations.NewOperation(
	"TranslateOnRampToFQDestDynamicCfgOps",
	semver.MustParse("1.0.0"),
	"Ports dynamic config from all 1.5.0 OnRamps into FeeQuoter DestChain Dynamic Config",
	func(b operations.Bundle, deps MigrateOnRampToFQDeps, input MigrateDestChainCfgInput) (MigrateDestDynamicCfgOutput, error) {
		onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(input.OnRamp, deps.Chain.Client)
		if err != nil {
			return MigrateDestDynamicCfgOutput{}, fmt.Errorf("Failed to create EVM2EVMOnRamp contract binding at address %s: %w", input.OnRamp.Hex(), err)
		}
		migrateOnRamp := migration.EVM2EVMOnRampMigrate{EVM2EVMOnRamp: onRamp}
		destChainEVM2EVMDynamicCfg, err2 := onRamp.GetDynamicConfig(nil)
		if err2 != nil && destChainEVM2EVMDynamicCfg.PriceRegistry == (common.Address{}) {
			return MigrateDestDynamicCfgOutput{}, fmt.Errorf("Cannot GetDynamicConfig for destination Chain: %d, for 1.5.0 OnRamp %s: %w", input.DestChainSelector, onRamp.Address().Hex(), err2)
		}

		fqDestConfig := migrateOnRamp.TranslateOnrampToFeequoterDynamicConfig(input.DestChainSelector, destChainEVM2EVMDynamicCfg)
		return MigrateDestDynamicCfgOutput{
			FeeQuoterDestChainConfig:   fqDestConfig,
			EVM2EVMOnRampDynamicConfig: destChainEVM2EVMDynamicCfg,
		}, nil
	}) */

	GetAllFeeTokensOps = operations.NewOperation(
		"GetAllFeeTokensOps",
		semver.MustParse("1.0.0"),
		"Gets the FeeTokens for a price Registry",
		func(b operations.Bundle, deps MigrateOnRampToFQDeps, input common.Address) ([]common.Address, error) {
			priceRegistry, err := price_registry.NewPriceRegistry(input, deps.Chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create PriceRegistry contract binding on source chain %d: %w", deps.Chain.Selector, err)
			}

			allFeeTokens, err2 := priceRegistry.GetFeeTokens(nil)
			if err2 != nil {
				return nil, fmt.Errorf("failed to all tokens on PriceRegistry %d for  source chain %d: %w", input.Hex(), deps.Chain.Selector, err)
			}

			return allFeeTokens, nil
		})

	GetEVM2EVMOnRampGetFeeTokenConfigOps = operations.NewOperation(
		"GetEVM2EVMOnRampGetFeeTokenConfigOps",
		semver.MustParse("1.0.0"),
		"Gets the FeeTokenConfigs for a given fee token",
		func(b operations.Bundle, deps MigrateOnRampToFQDeps, input OnRampGetFeeTokenCfgInput) (evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfig, error) {
			onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(input.OnRamp, deps.Chain.Client)
			if err != nil {
				return evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfig{}, fmt.Errorf("Failed to create EVM2EVMOnRamp contract binding at address %s: %w", input.OnRamp.Hex(), err)
			}

			feeTokenCfg, err2 := onRamp.GetFeeTokenConfig(nil, input.FeeTokenAddress)
			if err2 != nil {
				return evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfig{}, fmt.Errorf("Cannot GetFeeTokenConfig for Feetoken address: %d, for 1.5.0 OnRamp %d: %w", input.FeeTokenAddress, onRamp.Address().Hex(), err2)
			}

			return feeTokenCfg, nil
		})

	GetAllConfiguredTokensOps = operations.NewOperation(
		"GetAllConfiguredTokensOps",
		semver.MustParse("1.0.0"),
		"Gets all configured tokens from the TokenAdminRegistry",
		func(b operations.Bundle, deps MigrateOnRampToFQDeps, input GetAllConfiguredTokensInput) ([]common.Address, error) {
			tokenAdminReg, err := token_admin_registry.NewTokenAdminRegistry(input.TokenAdminRegistry, deps.Chain.Client)
			if err != nil {
				return nil, fmt.Errorf("Failed to create TokenAdminRegistry contract binding", "chainSelector", deps.Chain.ChainSelector(), "TokenAdminRegistry Address", input.TokenAdminRegistry.Hex(), "error", err)
			}

			allTransferTokens := []common.Address{}
			var offset uint64 = 0
			const pageSize uint64 = 1000
			for {
				pageTokens, err := tokenAdminReg.GetAllConfiguredTokens(nil, offset, pageSize)
				if err != nil {
					return nil, fmt.Errorf("failed to get all configured tokens from TokenAdminRegistry at offset %d: %w", offset, err)
				}

				if len(pageTokens) == 0 { // No more tokens to fetch
					break
				}

				allTransferTokens = append(allTransferTokens, pageTokens...)
				offset += pageSize
			}
			return allTransferTokens, nil
		})

	GetEVM2EVMOnRampPoolBySourceTokenOps = operations.NewOperation(
		"GetPoolBySourceTokenOps",
		semver.MustParse("1.0.0"),
		"Gets all TokenPools for a given destination chain And source token",
		func(b operations.Bundle, deps MigrateOnRampToFQDeps, input GetPoolBySourceTokenInput) (common.Address, error) {
			onramp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(input.OnRamp, deps.Chain.Client)
			if err != nil {
				return common.Address{}, fmt.Errorf("Failed to create EVM2EVMOnRamp contract binding", "chainSelector", deps.Chain.ChainSelector(), "OnRamp Address", input.OnRamp.Hex(), "error", err)
			}
			tokenPoolAddress, err := onramp.GetPoolBySourceToken(nil, input.DestChainSelector, input.FeeTokenAddress)
			if err != nil {
				return common.Address{}, fmt.Errorf("Failed to get pool for token on 1.5.0 OnRamp", "destinationChainSelector", input.DestChainSelector, "Fee Token", input.FeeTokenAddress.Hex(), "error", err)

			}
			return tokenPoolAddress, nil
		})

	GetSupportedChainsForTokenPool = operations.NewOperation(
		"GetSupportedChainsForTokenPool",
		semver.MustParse("1.0.0"),
		"Gets all Supported Chains of a Token Pool",
		func(b operations.Bundle, deps MigrateOnRampToFQDeps, tokenPoolAddress common.Address) ([]uint64, error) {
			tokenPool, err := token_pool.NewTokenPool(tokenPoolAddress, deps.Chain.Client)
			if err != nil {
				return nil, fmt.Errorf("Failed to create tokenpool contract binding", "chainSelector", deps.Chain.ChainSelector(), "Token Pool", tokenPoolAddress.Hex(), "error", err)
			}

			supportedChains, err := tokenPool.GetSupportedChains(nil)
			if err != nil {
				return nil, fmt.Errorf("Failed to get supported chains from token pool", "chainSelector", deps.Chain.ChainSelector(), "Token Pool", tokenPoolAddress.Hex(), "error", err)
			}
			return supportedChains, nil
		})
)
