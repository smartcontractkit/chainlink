package config

import (
	"math/big"

	chainsel "github.com/smartcontractkit/chain-selectors"
	aptos_fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	evm_fee_quoter "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
)

type IChainDefinition interface {
	GetChainFamily() string
}

type EVMChainDefinition struct {
	v1_6.ChainDefinition
}

func (c EVMChainDefinition) GetChainFamily() string {
	return chainsel.FamilyEVM
}

func (c EVMChainDefinition) GetConvertedAptosFeeQuoterConfig() aptos_fee_quoter.DestChainConfig {
	efqc := c.FeeQuoterDestChainConfig
	// Handle the byte slice to fixed-size array conversion
	return aptos_fee_quoter.DestChainConfig{
		IsEnabled:                         efqc.IsEnabled,
		MaxNumberOfTokensPerMsg:           efqc.MaxNumberOfTokensPerMsg,
		MaxDataBytes:                      efqc.MaxDataBytes,
		MaxPerMsgGasLimit:                 efqc.MaxPerMsgGasLimit,
		DestGasOverhead:                   efqc.DestGasOverhead,
		DestGasPerPayloadByteBase:         efqc.DestGasPerPayloadByteBase,
		DestGasPerPayloadByteHigh:         efqc.DestGasPerPayloadByteHigh,
		DestGasPerPayloadByteThreshold:    efqc.DestGasPerPayloadByteThreshold,
		DestDataAvailabilityOverheadGas:   efqc.DestDataAvailabilityOverheadGas,
		DestGasPerDataAvailabilityByte:    efqc.DestGasPerDataAvailabilityByte,
		DestDataAvailabilityMultiplierBps: efqc.DestDataAvailabilityMultiplierBps,
		ChainFamilySelector:               efqc.ChainFamilySelector[:],
		EnforceOutOfOrder:                 efqc.EnforceOutOfOrder,
		DefaultTokenFeeUsdCents:           efqc.DefaultTokenFeeUSDCents,
		DefaultTokenDestGasOverhead:       efqc.DefaultTokenDestGasOverhead,
		DefaultTxGasLimit:                 efqc.DefaultTxGasLimit,
		GasMultiplierWeiPerEth:            efqc.GasMultiplierWeiPerEth,
		GasPriceStalenessThreshold:        efqc.GasPriceStalenessThreshold,
		NetworkFeeUsdCents:                efqc.NetworkFeeUSDCents,
	}
}

type AptosChainDefinition struct {
	// ConnectionConfig holds configuration for connection.
	v1_6.ConnectionConfig `json:"connectionConfig"`
	// Selector is the chain selector of this chain.
	Selector uint64 `json:"selector"`
	// GasPrice defines the USD price (18 decimals) per unit gas for this chain as a destination.
	GasPrice *big.Int `json:"gasPrice"`
	// FeeQuoterDestChainConfig is the configuration on a fee quoter for this chain as a destination.
	FeeQuoterDestChainConfig aptos_fee_quoter.DestChainConfig `json:"feeQuoterDestChainConfig"`
}

func (c AptosChainDefinition) GetChainFamily() string {
	return chainsel.FamilyAptos
}

func (c AptosChainDefinition) GetConvertedEVMFeeQuoterConfig() evm_fee_quoter.FeeQuoterDestChainConfig {
	afqc := c.FeeQuoterDestChainConfig
	// Handle the byte slice to fixed-size array conversion
	var chainFamilySelector [4]byte
	// Copy up to 4 bytes, zero-padding if source is shorter
	copy(chainFamilySelector[:], afqc.ChainFamilySelector)

	return evm_fee_quoter.FeeQuoterDestChainConfig{
		IsEnabled:                         afqc.IsEnabled,
		MaxNumberOfTokensPerMsg:           afqc.MaxNumberOfTokensPerMsg,
		MaxDataBytes:                      afqc.MaxDataBytes,
		MaxPerMsgGasLimit:                 afqc.MaxPerMsgGasLimit,
		DestGasOverhead:                   afqc.DestGasOverhead,
		DestGasPerPayloadByteBase:         uint8(afqc.DestGasPerPayloadByteBase),
		DestGasPerPayloadByteHigh:         uint8(afqc.DestGasPerPayloadByteHigh),
		DestGasPerPayloadByteThreshold:    afqc.DestGasPerPayloadByteThreshold,
		DestDataAvailabilityOverheadGas:   afqc.DestDataAvailabilityOverheadGas,
		DestGasPerDataAvailabilityByte:    afqc.DestGasPerDataAvailabilityByte,
		DestDataAvailabilityMultiplierBps: afqc.DestDataAvailabilityMultiplierBps,
		ChainFamilySelector:               chainFamilySelector,
		EnforceOutOfOrder:                 afqc.EnforceOutOfOrder,
		DefaultTokenFeeUSDCents:           afqc.DefaultTokenFeeUsdCents,
		DefaultTokenDestGasOverhead:       afqc.DefaultTokenDestGasOverhead,
		DefaultTxGasLimit:                 afqc.DefaultTxGasLimit,
		GasMultiplierWeiPerEth:            afqc.GasMultiplierWeiPerEth,
		GasPriceStalenessThreshold:        afqc.GasPriceStalenessThreshold,
		NetworkFeeUSDCents:                afqc.NetworkFeeUsdCents,
	}
}
