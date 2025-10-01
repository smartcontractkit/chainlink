package ccipnoop

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

func NewGasEstimateProvider(codec ccipocr3.ExtraDataCodecBundle) ccipocr3.EstimateProvider {
	return estimateProvider{
		extraDataCodec: codec,
	}
}

type estimateProvider struct {
	extraDataCodec ccipocr3.ExtraDataCodecBundle
}

// CalculateMerkleTreeGas is not implemented
func (gp estimateProvider) CalculateMerkleTreeGas(numRequests int) uint64 {
	return 1
}

// CalculateMessageMaxGas is not implemented.
func (gp estimateProvider) CalculateMessageMaxGas(msg ccipocr3.Message) uint64 {
	return 1
}
