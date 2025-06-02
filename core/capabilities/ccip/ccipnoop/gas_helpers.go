package ccipnoop

import (
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

func NewNoopGasEstimateProvider(codec ccipcommon.ExtraDataCodec) NoopEstimateProvider {
	return NoopEstimateProvider{
		extraDataCodec: codec,
	}
}

type NoopEstimateProvider struct {
	extraDataCodec ccipcommon.ExtraDataCodec
}

// CalculateMerkleTreeGas is not implemented
func (gp NoopEstimateProvider) CalculateMerkleTreeGas(numRequests int) uint64 {
	return 1
}

// CalculateMessageMaxGas is not implemented.
func (gp NoopEstimateProvider) CalculateMessageMaxGas(msg cciptypes.Message) uint64 {
	return 1
}
