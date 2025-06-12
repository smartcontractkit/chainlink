package migration

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	onramp1_5 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
)

type EVM2EVMOnRampMigrate struct {
	*onramp1_5.EVM2EVMOnRamp
}

type EVM2EVMOnRampMigrateDestChainConfig struct {
	fee_quoter.FeeQuoterDestChainConfig
}

// Translate the dynamic config fields from the 1.5.0 OnRamp to the FeeQuoterDestChainConfig on 1.6 FeeQuoter
// Start with default base values & then override with the values from the 1.5.0 OnRamp
func (m EVM2EVMOnRampMigrate) TranslateOnrampToFeequoterDynamicConfig(destChainSel uint64, destChainEVM2EVMDynamicCfg onramp1_5.EVM2EVMOnRampDynamicConfig) fee_quoter.FeeQuoterDestChainConfig {
	fqDestConfig := DefaultFeeQuoterDestChainConfig(true, destChainSel)

	fqDestConfig.MaxNumberOfTokensPerMsg = destChainEVM2EVMDynamicCfg.MaxNumberOfTokensPerMsg
	fqDestConfig.MaxDataBytes = destChainEVM2EVMDynamicCfg.MaxDataBytes
	fqDestConfig.MaxPerMsgGasLimit = destChainEVM2EVMDynamicCfg.MaxPerMsgGasLimit
	fqDestConfig.DestGasOverhead = destChainEVM2EVMDynamicCfg.DestGasOverhead
	// fqDestConfig.DestGasPerPayloadByteBase = destChainEVM2EVMDynamicCfg.DestGasPerPayloadByte
	// fqDestConfig.DestGasPerPayloadByteHigh = destChainEVM2EVMDynamicCfg.DestGasPerPayloadByte
	// fqDestConfig.DestGasPerPayloadByteThreshold = destChainEVM2EVMDynamicCfg.DestGasPerPayloadByte
	fqDestConfig.DestDataAvailabilityOverheadGas = destChainEVM2EVMDynamicCfg.DestDataAvailabilityOverheadGas
	fqDestConfig.DestGasPerDataAvailabilityByte = destChainEVM2EVMDynamicCfg.DestGasPerDataAvailabilityByte
	fqDestConfig.DestDataAvailabilityMultiplierBps = destChainEVM2EVMDynamicCfg.DestDataAvailabilityMultiplierBps
	fqDestConfig.EnforceOutOfOrder = destChainEVM2EVMDynamicCfg.EnforceOutOfOrder
	fqDestConfig.DefaultTokenFeeUSDCents = destChainEVM2EVMDynamicCfg.DefaultTokenFeeUSDCents
	fqDestConfig.DefaultTokenDestGasOverhead = destChainEVM2EVMDynamicCfg.DefaultTokenDestGasOverhead
	// fqDestConfig.DefaultTxGasLimit = destChainEVM2EVMDynamicCfg.DefaultTxGasLimit // is this from static config?
	// fqDestConfig.GasPriceStalenessThreshold = destChainEVM2EVMDynamicCfg.GasPriceStalenessThreshold // where do we get this from?
	// fqDestConfig.GasMultiplierWeiPerEth = destChainEVM2EVMDynamicCfg.GasMultiplierWeiPerEth // where do we get this from? : FeeTokenConfig in onramp -- Probably not needed & can use the default instantiation
	// fqDestConfig.NetworkFeeUSDCents = destChainEVM2EVMDynamicCfg.NetworkFeeUSDCents // where do we get this from? : FeeTokenConfig in onramp - -- Probably not needed & can use the default instantiation

	return fqDestConfig
}

func (m *EVM2EVMOnRampMigrateDestChainConfig) TranslateOnrampToFeequoterDynamicConfig(destChainSel uint64, destChainEVM2EVMDynamicCfg onramp1_5.EVM2EVMOnRampDynamicConfig) {
	fqDestDefaults := DefaultFeeQuoterDestChainConfig(true, destChainSel)

	m.MaxNumberOfTokensPerMsg = destChainEVM2EVMDynamicCfg.MaxNumberOfTokensPerMsg
	m.MaxDataBytes = destChainEVM2EVMDynamicCfg.MaxDataBytes
	m.MaxPerMsgGasLimit = destChainEVM2EVMDynamicCfg.MaxPerMsgGasLimit
	m.DestGasOverhead = destChainEVM2EVMDynamicCfg.DestGasOverhead
	m.DestGasPerPayloadByteBase = fqDestDefaults.DestGasPerPayloadByteBase
	m.DestGasPerPayloadByteHigh = fqDestDefaults.DestGasPerPayloadByteHigh
	m.DestGasPerPayloadByteThreshold = fqDestDefaults.DestGasPerPayloadByteThreshold
	m.DestDataAvailabilityOverheadGas = destChainEVM2EVMDynamicCfg.DestDataAvailabilityOverheadGas
	m.DestGasPerDataAvailabilityByte = destChainEVM2EVMDynamicCfg.DestGasPerDataAvailabilityByte
	m.DestDataAvailabilityMultiplierBps = destChainEVM2EVMDynamicCfg.DestDataAvailabilityMultiplierBps
	m.EnforceOutOfOrder = destChainEVM2EVMDynamicCfg.EnforceOutOfOrder
	m.DefaultTokenFeeUSDCents = destChainEVM2EVMDynamicCfg.DefaultTokenFeeUSDCents
	m.DefaultTokenDestGasOverhead = destChainEVM2EVMDynamicCfg.DefaultTokenDestGasOverhead
	m.DefaultTxGasLimit = fqDestDefaults.DefaultTxGasLimit
	m.ChainFamilySelector = fqDestDefaults.ChainFamilySelector
	m.IsEnabled = fqDestDefaults.IsEnabled
	// m.GasPriceStalenessThreshold = destChainEVM2EVMDynamicCfg.GasPriceStalenessThreshold // where do we get this from?
	// m.GasMultiplierWeiPerEth = destChainEVM2EVMDynamicCfg.GasMultiplierWeiPerEth // where do we get this from? : FeeTokenConfig in onramp -- Probably not needed & can use the default instantiation
	// m.NetworkFeeUSDCents = destChainEVM2EVMDynamicCfg.NetworkFeeUSDCents // where do we get this from? : FeeTokenConfig in onramp - -- Probably not needed & can use the default instantiation
}

func (m *EVM2EVMOnRampMigrateDestChainConfig) TranslateOnrampToFeequoterFeeTokenCfg(feetokenCfg onramp1_5.EVM2EVMOnRampFeeTokenConfig) {

	// m.DefaultTxGasLimit = destChainEVM2EVMDynamicCfg.DefaultTxGasLimit // is this from static config?
	// m.GasPriceStalenessThreshold = destChainEVM2EVMDynamicCfg.GasPriceStalenessThreshold // where do we get this from?
	m.GasMultiplierWeiPerEth = feetokenCfg.GasMultiplierWeiPerEth // where do we get this from? : FeeTokenConfig in onramp -- Probably not needed & can use the default instantiation
	m.NetworkFeeUSDCents = feetokenCfg.NetworkFeeUSDCents         // where do we get this from? : FeeTokenConfig in onramp - -- Probably not needed & can use the default instantiation
}

func (m EVM2EVMOnRampMigrate) TranslateOnrampToFeequoterTokenTransferFeeConfig(destChainSel uint64, token common.Address, onRampTokenTransferFeeConfig onramp1_5.EVM2EVMOnRampTokenTransferFeeConfig) fee_quoter.FeeQuoterTokenTransferFeeConfigSingleTokenArgs {
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

// the below needs to be moved to a common package
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
		GasPriceStalenessThreshold:        90000,
	}
}
