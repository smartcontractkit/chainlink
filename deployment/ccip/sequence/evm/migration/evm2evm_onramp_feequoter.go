package migration

import (
	"github.com/ethereum/go-ethereum/common"

	onramp1_5 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
)

type EVM2EVMOnRampMigrate struct {
	*onramp1_5.EVM2EVMOnRamp
}

type EVM2EVMOnRampMigrateDestChainConfig struct {
	fee_quoter.FeeQuoterDestChainConfig
}

type EVM2EVMOnRampMigratePremiumMultiplierCfg struct {
	fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs
}

// Translate the dynamic config fields from the 1.5.0 OnRamp to the FeeQuoterDestChainConfig on 1.6 FeeQuoter
// Start with default base values & then override with the values from the 1.5.0 OnRamp
func (m *EVM2EVMOnRampMigrateDestChainConfig) TranslateOnrampToFeequoterDynamicConfig(destChainSel uint64, destChainEVM2EVMDynamicCfg onramp1_5.EVM2EVMOnRampDynamicConfig) {
	fqDestDefaults := ccipops.DefaultFeeQuoterDestChainConfig(true, destChainSel)

	m.MaxNumberOfTokensPerMsg = destChainEVM2EVMDynamicCfg.MaxNumberOfTokensPerMsg
	m.DestGasOverhead = destChainEVM2EVMDynamicCfg.DestGasOverhead
	m.DestGasPerPayloadByteBase = fqDestDefaults.DestGasPerPayloadByteBase
	m.DestGasPerPayloadByteHigh = fqDestDefaults.DestGasPerPayloadByteHigh
	m.DestGasPerPayloadByteThreshold = fqDestDefaults.DestGasPerPayloadByteThreshold
	m.DestDataAvailabilityOverheadGas = destChainEVM2EVMDynamicCfg.DestDataAvailabilityOverheadGas
	m.DestGasPerDataAvailabilityByte = destChainEVM2EVMDynamicCfg.DestGasPerDataAvailabilityByte
	m.DestDataAvailabilityMultiplierBps = destChainEVM2EVMDynamicCfg.DestDataAvailabilityMultiplierBps
	m.MaxDataBytes = destChainEVM2EVMDynamicCfg.MaxDataBytes
	m.MaxPerMsgGasLimit = destChainEVM2EVMDynamicCfg.MaxPerMsgGasLimit
	m.EnforceOutOfOrder = destChainEVM2EVMDynamicCfg.EnforceOutOfOrder
	m.DefaultTokenFeeUSDCents = destChainEVM2EVMDynamicCfg.DefaultTokenFeeUSDCents
	m.DefaultTokenDestGasOverhead = destChainEVM2EVMDynamicCfg.DefaultTokenDestGasOverhead
	m.DefaultTxGasLimit = fqDestDefaults.DefaultTxGasLimit
	m.ChainFamilySelector = fqDestDefaults.ChainFamilySelector
	m.IsEnabled = fqDestDefaults.IsEnabled
	m.GasPriceStalenessThreshold = fqDestDefaults.GasPriceStalenessThreshold
}

func (m *EVM2EVMOnRampMigrateDestChainConfig) TranslateOnrampToFeequoterFeeTokenCfg(feetokenCfg onramp1_5.EVM2EVMOnRampFeeTokenConfig) {
	m.GasMultiplierWeiPerEth = feetokenCfg.GasMultiplierWeiPerEth
	m.NetworkFeeUSDCents = feetokenCfg.NetworkFeeUSDCents
}

func (m *EVM2EVMOnRampMigratePremiumMultiplierCfg) TranslateOnrampToFeeQFeePremiumCfg(token common.Address, feetokenCfg onramp1_5.EVM2EVMOnRampFeeTokenConfig) {
	m.Token = token
	m.PremiumMultiplierWeiPerEth = feetokenCfg.GasMultiplierWeiPerEth
}

func (m EVM2EVMOnRampMigrate) TranslateOnrampToFeequoterTokenTransferFeeConfig(token common.Address, onRampTokenTransferFeeConfig onramp1_5.EVM2EVMOnRampTokenTransferFeeConfig) fee_quoter.FeeQuoterTokenTransferFeeConfigSingleTokenArgs {
	return fee_quoter.FeeQuoterTokenTransferFeeConfigSingleTokenArgs{
		Token: token,
		TokenTransferFeeConfig: fee_quoter.FeeQuoterTokenTransferFeeConfig{
			MinFeeUSDCents:    onRampTokenTransferFeeConfig.MinFeeUSDCents,
			MaxFeeUSDCents:    onRampTokenTransferFeeConfig.MaxFeeUSDCents,
			DeciBps:           onRampTokenTransferFeeConfig.DeciBps,
			DestGasOverhead:   onRampTokenTransferFeeConfig.DestGasOverhead,
			DestBytesOverhead: onRampTokenTransferFeeConfig.DestBytesOverhead,
			IsEnabled:         onRampTokenTransferFeeConfig.IsEnabled,
		},
	}
}
