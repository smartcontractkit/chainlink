package ccipsolana

import (
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

func NewGasEstimateProvider() EstimateProvider {
	return EstimateProvider{}
}

type EstimateProvider struct {
}

// CalculateMerkleTreeGas is not implemented
func (gp EstimateProvider) CalculateMerkleTreeGas(numRequests int) uint64 {
	return 0
}

// CalculateMessageMaxGas is not implemented.
func (gp EstimateProvider) CalculateMessageMaxGas(msg cciptypes.Message) uint64 {
	return 0
}
