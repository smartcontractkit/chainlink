package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/chainlink/deployment"
)

type FeeQuoterView struct {
	DestinationChainConfig map[uint64]FeeQuoterDestChainConfig `json:"destinationChainConfig,omitempty"`
}

type FeeQuoterDestChainConfig struct {
	IsEnabled                         bool   `json:"isEnabled,omitempty"`
	LaneCodeVersion                   string `json:"laneCodeVersion,omitempty"`
	MaxNumberOfTokensPerMsg           uint16 `json:"maxNumberOfTokensPerMsg,omitempty"`
	MaxDataBytes                      uint32 `json:"maxDataBytes,omitempty"`
	MaxPerMsgGasLimit                 uint32 `json:"maxPerMsgGasLimit,omitempty"`
	DestGasOverhead                   uint32 `json:"destGasOverhead,omitempty"`
	DestGasPerPayloadByteBase         uint32 `json:"destGasPerPayloadByteBase,omitempty"`
	DestGasPerPayloadByteHigh         uint32 `json:"destGasPerPayloadByteHigh,omitempty"`
	DestGasPerPayloadByteThreshold    uint32 `json:"destGasPerPayloadByteThreshold,omitempty"`
	DestDataAvailabilityOverheadGas   uint32 `json:"destDataAvailabilityOverheadGas,omitempty"`
	DestGasPerDataAvailabilityByte    uint16 `json:"destGasPerDataAvailabilityByte,omitempty"`
	DestDataAvailabilityMultiplierBps uint16 `json:"destDataAvailabilityMultiplierBps,omitempty"`
	DefaultTokenFeeUSDCents           uint16 `json:"defaultTokenFeeUSDCents,omitempty"`
	DefaultTokenDestGasOverhead       uint32 `json:"defaultTokenDestGasOverhead,omitempty"`
	DefaultTxGasLimit                 uint32 `json:"defaultTxGasLimit,omitempty"`
	GasMultiplierWeiPerEth            uint64 `json:"gasMultiplierWeiPerEth,omitempty"`
	NetworkFeeUSDCents                uint32 `json:"networkFeeUSDCents,omitempty"`
	GasPriceStalenessThreshold        uint32 `json:"gasPriceStalenessThreshold,omitempty"`
	EnforceOutOfOrder                 bool   `json:"enforceOutOfOrder,omitempty"`
	ChainFamilySelector               string `json:"chainFamilySelector,omitempty"`
}

func GenerateFeeQuoterView(chain deployment.SolChain, program solana.PublicKey, remoteChains []uint64) (FeeQuoterView, error) {
	fq := FeeQuoterView{}
	fq.DestinationChainConfig = make(map[uint64]FeeQuoterDestChainConfig)
	for _, remote := range remoteChains {
		fqRemoteChainPDA, _, err := solState.FindFqDestChainPDA(remote, program)
		if err != nil {
			return fq, fmt.Errorf("failed to find dest chain state pda for remote chain %d: %w", remote, err)
		}
		var destChainStateAccount solFeeQuoter.DestChain
		if err = chain.GetAccountDataBorshInto(context.Background(), fqRemoteChainPDA, &destChainStateAccount); err != nil {
			return fq, fmt.Errorf("remote %d is not configured on solana chain feequoter %d", remote, chain.Selector)
		}
		fq.DestinationChainConfig[remote] = FeeQuoterDestChainConfig{
			IsEnabled:                         destChainStateAccount.Config.IsEnabled,
			LaneCodeVersion:                   destChainStateAccount.Config.LaneCodeVersion.String(),
			MaxNumberOfTokensPerMsg:           destChainStateAccount.Config.MaxNumberOfTokensPerMsg,
			MaxDataBytes:                      destChainStateAccount.Config.MaxDataBytes,
			MaxPerMsgGasLimit:                 destChainStateAccount.Config.MaxPerMsgGasLimit,
			DestGasOverhead:                   destChainStateAccount.Config.DestGasOverhead,
			DestGasPerPayloadByteBase:         destChainStateAccount.Config.DestGasPerPayloadByteBase,
			DestGasPerPayloadByteHigh:         destChainStateAccount.Config.DestGasPerPayloadByteHigh,
			DestGasPerPayloadByteThreshold:    destChainStateAccount.Config.DestGasPerPayloadByteThreshold,
			DestDataAvailabilityOverheadGas:   destChainStateAccount.Config.DestDataAvailabilityOverheadGas,
			DestGasPerDataAvailabilityByte:    destChainStateAccount.Config.DestGasPerDataAvailabilityByte,
			DestDataAvailabilityMultiplierBps: destChainStateAccount.Config.DestDataAvailabilityMultiplierBps,
			DefaultTokenFeeUSDCents:           destChainStateAccount.Config.DefaultTokenFeeUsdcents,
			DefaultTokenDestGasOverhead:       destChainStateAccount.Config.DefaultTokenDestGasOverhead,
			DefaultTxGasLimit:                 destChainStateAccount.Config.DefaultTxGasLimit,
			GasMultiplierWeiPerEth:            destChainStateAccount.Config.GasMultiplierWeiPerEth,
			NetworkFeeUSDCents:                destChainStateAccount.Config.NetworkFeeUsdcents,
			GasPriceStalenessThreshold:        destChainStateAccount.Config.GasPriceStalenessThreshold,
			EnforceOutOfOrder:                 destChainStateAccount.Config.EnforceOutOfOrder,
			ChainFamilySelector:               fmt.Sprintf("%x", destChainStateAccount.Config.ChainFamilySelector),
		}
	}

	return fq, nil
}
