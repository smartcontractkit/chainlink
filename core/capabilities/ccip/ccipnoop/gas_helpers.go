package ccipnoop

import (
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

func NewGasEstimateProvider(codec common.ExtraDataCodec) ccipocr3.EstimateProvider {
	return estimateProvider{
		extraDataCodec: codec,
	}
}

type estimateProvider struct {
	extraDataCodec common.ExtraDataCodec
}

// CalculateMerkleTreeGas is not implemented
func (gp estimateProvider) CalculateMerkleTreeGas(numRequests int) uint64 {
	return 1
}

// CalculateMessageMaxGas is not implemented.
func (gp estimateProvider) CalculateMessageMaxGas(msg ccipocr3.Message) uint64 {
	return 1
}
