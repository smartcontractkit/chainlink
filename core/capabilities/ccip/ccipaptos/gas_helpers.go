package ccipaptos

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

func NewGasEstimateProvider() EstimateProvider {
	return EstimateProvider{}
}

type EstimateProvider struct{}

// CalculateMerkleTreeGas is not implemented
func (gp EstimateProvider) CalculateMerkleTreeGas(numRequests int) uint64 {
	return 0
}

// CalculateMessageMaxGas is not implemented.
func (gp EstimateProvider) CalculateMessageMaxGas(msg ccipocr3.Message) uint64 {
	return 0
}
