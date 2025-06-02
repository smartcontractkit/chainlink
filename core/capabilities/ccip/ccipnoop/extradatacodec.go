package ccipnoop

import (
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// NoopExtraDataDecoder is a helper struct for decoding extra data
type NoopExtraDataDecoder struct{}

// DecodeExtraArgsToMap is a helper function for converting Borsh encoded extra args bytes into map[string]any
func (d NoopExtraDataDecoder) DecodeExtraArgsToMap(extraArgs cciptypes.Bytes) (map[string]any, error) {
	outputMap := make(map[string]any)
	return outputMap, nil
}

// DecodeDestExecDataToMap is a helper function for converting dest exec data bytes into map[string]any
func (d NoopExtraDataDecoder) DecodeDestExecDataToMap(destExecData cciptypes.Bytes) (map[string]any, error) {
	outputMap := make(map[string]any)
	return outputMap, nil
}

// Ensure NoopExtraDataDecoder implements the SourceChainExtraDataCodec interface
var _ ccipcommon.SourceChainExtraDataCodec = &NoopExtraDataDecoder{}
