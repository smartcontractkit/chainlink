package v1_5

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type MigrateOnRampToFQDeps struct {
	Chain cldf_evm.Chain
}

type OnRampGetTokenCfgIn struct {
	OnRamp        common.Address
	Address       common.Address
	ChainSelector uint64
}

type GetPoolBySourceTokenIn struct {
	OnRamp            common.Address
	FeeTokenAddress   common.Address
	ChainSelector     uint64
	DestChainSelector uint64
}

var (
	EVM2EVMOnrampGetDynamicCfgOp = operations.NewOperation(
		"EVM2EVMOnrampGetDynamicCfgOp",
		semver.MustParse("1.0.0"),
		"Get DynamicConfig from 1.5.0 OnRamps",
		func(b operations.Bundle, deps MigrateOnRampToFQDeps, input common.Address) (evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig, error) {
			onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(input, deps.Chain.Client)
			if err != nil {
				return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{}, fmt.Errorf("failed to create EVM2EVMOnRamp contract binding at address %s: %w", input.Hex(), err)
			}
			destChainEVM2EVMDynamicCfg, err2 := onRamp.GetDynamicConfig(nil)
			if err2 != nil && destChainEVM2EVMDynamicCfg.PriceRegistry == (common.Address{}) {
				return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{}, fmt.Errorf("cannot GetDynamicConfig for destination Chain from 1.5.0 OnRamp %s: %w", onRamp.Address().Hex(), err2)
			}
			return destChainEVM2EVMDynamicCfg, nil
		})

	EVM2EVMOnrampGetFeeTokenConfigOp = operations.NewOperation(
		"EVM2EVMOnrampGetFeeTokenConfigOp",
		semver.MustParse("1.0.0"),
		"Gets the FeeTokenConfigs for a given fee token",
		func(b operations.Bundle, deps MigrateOnRampToFQDeps, input OnRampGetTokenCfgIn) (evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfig, error) {
			onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(input.OnRamp, deps.Chain.Client)
			if err != nil {
				return evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfig{}, fmt.Errorf("failed to create EVM2EVMOnRamp contract binding at address %s: %w", input.OnRamp.Hex(), err)
			}
			feeTokenCfg, err2 := onRamp.GetFeeTokenConfig(nil, input.Address)
			if err2 != nil {
				return evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfig{}, fmt.Errorf("cannot GetFeeTokenConfig for Feetoken address: %d, for 1.5.0 OnRamp %s: %w", input.Address, onRamp.Address().Hex(), err2)
			}

			return feeTokenCfg, nil
		})

	EVM2EVMOnrampGetPoolBySourceTokenOp = operations.NewOperation(
		"EVM2EVMOnrampGetPoolBySourceTokenOp",
		semver.MustParse("1.0.0"),
		"Gets all TokenPools for a given destination chain And source token",
		func(b operations.Bundle, deps MigrateOnRampToFQDeps, input GetPoolBySourceTokenIn) (common.Address, error) {
			onramp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(input.OnRamp, deps.Chain.Client)
			if err != nil {
				return common.Address{}, fmt.Errorf("failed to create EVM2EVMOnRamp contract binding: chainSelector=%d, OnRamp Address=%s, error=%w", deps.Chain.ChainSelector(), input.OnRamp.Hex(), err)
			}
			tokenPoolAddress, err := onramp.GetPoolBySourceToken(nil, input.DestChainSelector, input.FeeTokenAddress)
			if err != nil {
				return common.Address{}, fmt.Errorf("failed to get pool for token on 1.5.0 OnRamp: destinationChainSelector=%d, Fee Token=%s, error=%w", input.DestChainSelector, input.FeeTokenAddress.Hex(), err)

			}
			return tokenPoolAddress, nil
		})

	EVM2EVMOnrampGetTokenTransferFeeConfigOp = operations.NewOperation(
		"EVM2EVMOnrampGetTokenTransferFeeConfigOp",
		semver.MustParse("1.0.0"),
		"Gets all TokenPools for a given destination chain And source token",
		func(b operations.Bundle, deps MigrateOnRampToFQDeps, input OnRampGetTokenCfgIn) (evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfig, error) {
			onramp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(input.OnRamp, deps.Chain.Client)
			if err != nil {
				return evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfig{}, fmt.Errorf("failed to create EVM2EVMOnRamp contract binding: chainSelector=%d, OnRamp Address=%s, error=%w", deps.Chain.ChainSelector(), input.OnRamp.Hex(), err)
			}
			tokenTransferFeeCfg, err := onramp.GetTokenTransferFeeConfig(nil, input.Address)
			if err != nil {
				return evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfig{}, fmt.Errorf("failed to Get Token Transfer Fee Config for token on 1.5.0 OnRamp Token=%s, error=%w", input.Address.Hex(), err)

			}
			return tokenTransferFeeCfg, nil
		})
)
