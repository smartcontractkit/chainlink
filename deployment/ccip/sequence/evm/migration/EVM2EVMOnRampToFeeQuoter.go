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

type EVM2EVMOnRampMigratePremiumMultiplierCfg struct {
	fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs
}

// Translate the dynamic config fields from the 1.5.0 OnRamp to the FeeQuoterDestChainConfig on 1.6 FeeQuoter
// Start with default base values & then override with the values from the 1.5.0 OnRamp
func (m *EVM2EVMOnRampMigrateDestChainConfig) TranslateOnrampToFeequoterDynamicConfig(destChainSel uint64, destChainEVM2EVMDynamicCfg onramp1_5.EVM2EVMOnRampDynamicConfig) {
	fqDestDefaults := DefaultFeeQuoterDestChainConfig(true, destChainSel)

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
	/* return fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs{
		Token:                      token,
		PremiumMultiplierWeiPerEth: feetokenCfg.GasMultiplierWeiPerEth,
	} */
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

// TODO: the below needs to be moved to a common package, cannot right now due to circular dependencies
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
		switch destFamily {
		case chain_selectors.FamilySolana:
			familySelector, _ = hex.DecodeString(SVMFamilySelector) // solana
		case chain_selectors.FamilyAptos:
			familySelector, _ = hex.DecodeString(AptosFamilySelector) // aptos
		}
	}
	return fee_quoter.FeeQuoterDestChainConfig{
		IsEnabled:                         configEnabled,
		MaxNumberOfTokensPerMsg:           10,
		MaxDataBytes:                      30_000,
		MaxPerMsgGasLimit:                 3_000_000, // TODO: this needs to be updated based on RMN sig verification per chain?! 220/250K
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
