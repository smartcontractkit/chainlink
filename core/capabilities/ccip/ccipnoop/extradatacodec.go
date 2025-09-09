package ccipnoop

import (
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// extraDataDecoder is a helper struct for decoding extra data
type extraDataDecoder struct{}

// DecodeExtraArgsToMap is a helper function for converting Borsh encoded extra args bytes into map[string]any
func (d extraDataDecoder) DecodeExtraArgsToMap(extraArgs cciptypes.Bytes) (map[string]any, error) {
	outputMap := make(map[string]any)
	return outputMap, nil
}

// DecodeDestExecDataToMap is a helper function for converting dest exec data bytes into map[string]any
func (d extraDataDecoder) DecodeDestExecDataToMap(destExecData cciptypes.Bytes) (map[string]any, error) {
	outputMap := make(map[string]any)
	return outputMap, nil
}

// Ensure extraDataDecoder implements the SourceChainExtraDataCodec interface
var _ ccipcommon.SourceChainExtraDataCodec = &extraDataDecoder{}
