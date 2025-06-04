package config

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aptos_fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
)

func TestToEVMUpdateLanesConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    UpdateAptosLanesConfig
		expected v1_6.UpdateBidirectionalLanesChangesetConfigs
	}{
		{
			name: "EVM <> Aptos Biderectional Lane",
			input: UpdateAptosLanesConfig{
				EVMMCMSConfig: &proposalutils.TimelockConfig{},
				Lanes: []LaneConfig{
					{
						Source:     getEVMDef(),
						Dest:       getAptosDef(t),
						IsDisabled: false,
					},
					{
						Source:     getAptosDef(t),
						Dest:       getEVMDef(),
						IsDisabled: false,
					},
				},
				TestRouter: false,
			},
			expected: v1_6.UpdateBidirectionalLanesChangesetConfigs{
				UpdateFeeQuoterDestsConfig: v1_6.UpdateFeeQuoterDestsConfig{
					MCMS: &proposalutils.TimelockConfig{},
					UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
						1881: {
							4457093679053095497: fee_quoter.FeeQuoterDestChainConfig{
								IsEnabled:                         true,
								MaxNumberOfTokensPerMsg:           11,
								MaxDataBytes:                      40_000,
								MaxPerMsgGasLimit:                 4_000_000,
								DestGasOverhead:                   ccipevm.DestGasOverhead,
								DefaultTokenFeeUSDCents:           30,
								DestGasPerPayloadByteBase:         ccipevm.CalldataGasPerByteBase,
								DestGasPerPayloadByteHigh:         ccipevm.CalldataGasPerByteHigh,
								DestGasPerPayloadByteThreshold:    ccipevm.CalldataGasPerByteThreshold,
								DestDataAvailabilityOverheadGas:   700,
								DestGasPerDataAvailabilityByte:    17,
								DestDataAvailabilityMultiplierBps: 2,
								DefaultTokenDestGasOverhead:       100_000,
								DefaultTxGasLimit:                 100_000,
								GasMultiplierWeiPerEth:            12e17,
								NetworkFeeUSDCents:                20,
								EnforceOutOfOrder:                 false,
								GasPriceStalenessThreshold:        2,
								ChainFamilySelector:               [4]byte{0xac, 0x77, 0xff, 0xec}, // From AptosFamilySelector "ac77ffec"
							},
						},
					},
				},
				UpdateFeeQuoterPricesConfig: v1_6.UpdateFeeQuoterPricesConfig{
					MCMS: &proposalutils.TimelockConfig{},
					PricesByChain: map[uint64]v1_6.FeeQuoterPriceUpdatePerSource{
						1881: {
							GasPrices: map[uint64]*big.Int{
								4457093679053095497: big.NewInt(100),
							},
						},
					},
				},
				UpdateOnRampDestsConfig: v1_6.UpdateOnRampDestsConfig{
					MCMS: &proposalutils.TimelockConfig{},
					UpdatesByChain: map[uint64]map[uint64]v1_6.OnRampDestinationUpdate{
						1881: {
							4457093679053095497: v1_6.OnRampDestinationUpdate{
								IsEnabled:        true,
								TestRouter:       false,
								AllowListEnabled: true,
							},
						},
					},
				},
				UpdateOffRampSourcesConfig: v1_6.UpdateOffRampSourcesConfig{
					MCMS: &proposalutils.TimelockConfig{},
					UpdatesByChain: map[uint64]map[uint64]v1_6.OffRampSourceUpdate{
						1881: {
							4457093679053095497: v1_6.OffRampSourceUpdate{
								IsEnabled:                 true,
								TestRouter:                false,
								IsRMNVerificationDisabled: false,
							},
						},
					},
				},
				UpdateRouterRampsConfig: v1_6.UpdateRouterRampsConfig{
					TestRouter: false,
					MCMS:       &proposalutils.TimelockConfig{},
					UpdatesByChain: map[uint64]v1_6.RouterUpdates{
						1881: {
							OnRampUpdates: map[uint64]bool{
								4457093679053095497: true,
							},
							OffRampUpdates: map[uint64]bool{
								4457093679053095497: true,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToEVMUpdateLanesConfig(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func getEVMDef() EVMChainDefinition {
	return EVMChainDefinition{
		ChainDefinition: v1_6.ChainDefinition{
			ConnectionConfig: v1_6.ConnectionConfig{
				RMNVerificationDisabled: true,
				AllowListEnabled:        false,
			},
			Selector:                 1881,
			GasPrice:                 big.NewInt(1e17),
			FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
		},
		OnRampVersion: []byte{1, 6, 0},
	}
}

func getAptosDef(t *testing.T) AptosChainDefinition {
	return AptosChainDefinition{
		ConnectionConfig: v1_6.ConnectionConfig{AllowListEnabled: true},
		Selector:         4457093679053095497,
		GasPrice:         big.NewInt(100),
		FeeQuoterDestChainConfig: aptos_fee_quoter.DestChainConfig{
			IsEnabled:                         true,
			MaxNumberOfTokensPerMsg:           11,
			MaxDataBytes:                      40_000,
			MaxPerMsgGasLimit:                 4_000_000,
			DestGasOverhead:                   ccipevm.DestGasOverhead,
			DefaultTokenFeeUsdCents:           30,
			DestGasPerPayloadByteBase:         ccipevm.CalldataGasPerByteBase,
			DestGasPerPayloadByteHigh:         ccipevm.CalldataGasPerByteHigh,
			DestGasPerPayloadByteThreshold:    ccipevm.CalldataGasPerByteThreshold,
			DestDataAvailabilityOverheadGas:   700,
			DestGasPerDataAvailabilityByte:    17,
			DestDataAvailabilityMultiplierBps: 2,
			DefaultTokenDestGasOverhead:       100_000,
			DefaultTxGasLimit:                 100_000,
			GasMultiplierWeiPerEth:            12e17,
			NetworkFeeUsdCents:                20,
			ChainFamilySelector:               hexMustDecode(t, v1_6.AptosFamilySelector),
			EnforceOutOfOrder:                 false,
			GasPriceStalenessThreshold:        2,
		},
	}
}

func hexMustDecode(t *testing.T, s string) []byte {
	b, err := hex.DecodeString(s)
	require.NoError(t, err)
	return b
}
