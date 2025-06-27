package migration

import (
	"github.com/ethereum/go-ethereum/common"

	onramp1_5 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
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

// NewFeeQuoterDestChainConfigParams defines fee_quoter.FeeQuoterDestChainConfig parameters that do not have a 1.5.0 equivalent.
// They need to be explicitly defined as part of the input.
type NewFeeQuoterDestChainConfigParams struct {
	DestGasPerPayloadByteBase      uint8
	DestGasPerPayloadByteHigh      uint8
	DestGasPerPayloadByteThreshold uint16
	DefaultTxGasLimit              uint32
	ChainFamilySelector            [4]byte
	GasPriceStalenessThreshold     uint32
	GasMultiplierWeiPerEth         uint64
	NetworkFeeUSDCents             uint32
}

// Translate the dynamic config fields from the 1.5.0 OnRamp to the FeeQuoterDestChainConfig on 1.6 FeeQuoter
// Start with default base values & then override with the values from the 1.5.0 OnRamp
func (m *EVM2EVMOnRampMigrateDestChainConfig) TranslateOnrampToFeequoterDynamicConfig(destChainSel uint64, destChainEVM2EVMDynamicCfg onramp1_5.EVM2EVMOnRampDynamicConfig, newParams NewFeeQuoterDestChainConfigParams) {
	m.MaxNumberOfTokensPerMsg = destChainEVM2EVMDynamicCfg.MaxNumberOfTokensPerMsg
	m.DestGasOverhead = destChainEVM2EVMDynamicCfg.DestGasOverhead
	m.DestDataAvailabilityOverheadGas = destChainEVM2EVMDynamicCfg.DestDataAvailabilityOverheadGas
	m.DestGasPerDataAvailabilityByte = destChainEVM2EVMDynamicCfg.DestGasPerDataAvailabilityByte
	m.DestDataAvailabilityMultiplierBps = destChainEVM2EVMDynamicCfg.DestDataAvailabilityMultiplierBps
	m.MaxDataBytes = destChainEVM2EVMDynamicCfg.MaxDataBytes
	m.MaxPerMsgGasLimit = destChainEVM2EVMDynamicCfg.MaxPerMsgGasLimit
	m.EnforceOutOfOrder = destChainEVM2EVMDynamicCfg.EnforceOutOfOrder
	m.DefaultTokenFeeUSDCents = destChainEVM2EVMDynamicCfg.DefaultTokenFeeUSDCents
	m.DefaultTokenDestGasOverhead = destChainEVM2EVMDynamicCfg.DefaultTokenDestGasOverhead

	m.IsEnabled = true
	m.DestGasPerPayloadByteBase = newParams.DestGasPerPayloadByteBase
	m.DestGasPerPayloadByteHigh = newParams.DestGasPerPayloadByteHigh
	m.DestGasPerPayloadByteThreshold = newParams.DestGasPerPayloadByteThreshold
	m.DefaultTxGasLimit = newParams.DefaultTxGasLimit
	m.ChainFamilySelector = newParams.ChainFamilySelector
	m.GasPriceStalenessThreshold = newParams.GasPriceStalenessThreshold
	m.GasMultiplierWeiPerEth = newParams.GasMultiplierWeiPerEth
	m.NetworkFeeUSDCents = newParams.NetworkFeeUSDCents
}

func (m *EVM2EVMOnRampMigratePremiumMultiplierCfg) TranslateOnrampToFeeQFeePremiumCfg(token common.Address, feetokenCfg onramp1_5.EVM2EVMOnRampFeeTokenConfig) {
	m.Token = token
	m.PremiumMultiplierWeiPerEth = feetokenCfg.PremiumMultiplierWeiPerEth
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
