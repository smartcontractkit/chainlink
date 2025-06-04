package config

import (
	"fmt"
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk"
	chainsel "github.com/smartcontractkit/chain-selectors"

	aptos_fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	evm_fee_quoter "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	aptosstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/aptos"
)

// ChainDefinition is an interface that defines a chain config for lane deployment
// It is used to convert between Aptos and EVM fee quoter configs.
type ChainDefinition interface {
	GetChainFamily() string
	GetSelector() uint64
}

type EVMChainDefinition struct {
	v1_6.ChainDefinition
	OnRampVersion []byte
}

func (c EVMChainDefinition) GetChainFamily() string {
	return chainsel.FamilyEVM
}

func (c EVMChainDefinition) GetSelector() uint64 {
	return c.Selector
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
	// AddTokenTransferFeeConfigs is the configuration for token transfer fees to be added to fee quoter
	AddTokenTransferFeeConfigs []aptos_fee_quoter.TokenTransferFeeConfigAdded
	// RemoveTokenTransferFeeConfigs holds token transfer fees to be removed from fee quoter
	RemoveTokenTransferFeeConfigs []aptos_fee_quoter.TokenTransferFeeConfigRemoved
}

func (c AptosChainDefinition) GetChainFamily() string {
	return chainsel.FamilyAptos
}

func (c AptosChainDefinition) GetSelector() uint64 {
	return c.Selector
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
		DestGasPerPayloadByteBase:         afqc.DestGasPerPayloadByteBase,
		DestGasPerPayloadByteHigh:         afqc.DestGasPerPayloadByteHigh,
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

func (c AptosChainDefinition) Validate(client aptos.AptosRpcClient, state aptosstate.CCIPChainState) error {
	// Check CCIP Package
	if state.CCIPAddress == (aptos.AccountAddress{}) {
		return fmt.Errorf("package CCIP is not deployed on Aptos chain %d", c.Selector)
	}
	// Check OnRamp module
	hasOnramp, err := utils.IsModuleDeployed(client, state.CCIPAddress, "onramp")
	if err != nil || !hasOnramp {
		return fmt.Errorf("onRamp module is not deployed on Aptos chain %d: %w", c.Selector, err)
	}
	// Check OffRamp module
	hasOfframp, err := utils.IsModuleDeployed(client, state.CCIPAddress, "offramp")
	if err != nil || !hasOfframp {
		return fmt.Errorf("offRamp module is not deployed on Aptos chain %d: %w", c.Selector, err)
	}
	// Check Router module
	hasRouter, err := utils.IsModuleDeployed(client, state.CCIPAddress, "router")
	if err != nil || !hasRouter {
		return fmt.Errorf("router module is not deployed on Aptos chain %d: %w", c.Selector, err)
	}

	return nil
}
