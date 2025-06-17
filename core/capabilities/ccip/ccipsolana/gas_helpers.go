package ccipsolana

import (
	"fmt"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

const (
	// TODO: determine exact values
	PerTokenOverheadGas = 75_000  // releaseOrMint using callWithExactGas
	DestGasOverhead     = 350_000 // Commit and Exec costs
)

func NewGasEstimateProvider(codec ccipcommon.ExtraDataCodec) EstimateProvider {
	return EstimateProvider{
		extraDataCodec: codec,
	}
}

type EstimateProvider struct {
	extraDataCodec ccipcommon.ExtraDataCodec
}

// CalculateMerkleTreeGas is not implemented
func (gp EstimateProvider) CalculateMerkleTreeGas(numRequests int) uint64 {
	return 1
}

// CalculateMessageMaxGas is not implemented.
func (gp EstimateProvider) CalculateMessageMaxGas(msg cciptypes.Message) uint64 {
	maxGas, err := gp.CalculateMessageMaxGasWithError(msg)
	if err != nil {
		panic(err)
	}
	return maxGas
}

func (gp EstimateProvider) CalculateMessageMaxGasWithError(msg cciptypes.Message) (uint64, error) {
	decodedMap, err := gp.extraDataCodec.DecodeExtraArgs(msg.ExtraArgs, msg.Header.SourceChainSelector)
	if err != nil {
		return 0, fmt.Errorf("failed to decode extra args: %w", err)
	}

	extraData, err := parseExtraDataMap(decodedMap)
	if err != nil {
		return 0, fmt.Errorf("invalid extra args map: %w", err)
	}

	var totalTokenDestGasOverhead uint64
	for _, rampTokenAmount := range msg.TokenAmounts {
		decodedMap, err = gp.extraDataCodec.DecodeTokenAmountDestExecData(rampTokenAmount.DestExecData, msg.Header.SourceChainSelector)
		if err != nil {
			return 0, fmt.Errorf("failed to decode token dest gas overhead: %w", err)
		}

		tokenDestGasOverhead, err := extractDestGasAmountFromMap(decodedMap)
		if err != nil {
			return 0, fmt.Errorf("failed to extract dest gas amount from map: %w", err)
		}
		totalTokenDestGasOverhead += uint64(tokenDestGasOverhead)
	}

	// commit+exec overhead + token pool overhead + requested compute units
	// TODO: base cost in 5000 lamports, not compute units
	return DestGasOverhead +
		totalTokenDestGasOverhead +
		uint64(extraData.extraArgs.ComputeUnits), nil
}
