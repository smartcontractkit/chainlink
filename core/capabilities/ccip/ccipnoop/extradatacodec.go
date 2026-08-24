package ccipnoop

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// extraDataDecoder is a helper struct for decoding extra data
type extraDataDecoder struct{}

// DecodeExtraArgsToMap is a helper function for converting Borsh encoded extra args bytes into map[string]any
func (d extraDataDecoder) DecodeExtraArgsToMap(extraArgs ccipocr3.Bytes) (map[string]any, error) {
	outputMap := make(map[string]any)
	return outputMap, nil
}

// DecodeDestExecDataToMap is a helper function for converting dest exec data bytes into map[string]any
func (d extraDataDecoder) DecodeDestExecDataToMap(destExecData ccipocr3.Bytes) (map[string]any, error) {
	outputMap := make(map[string]any)
	return outputMap, nil
}

// Ensure extraDataDecoder implements the SourceChainExtraDataCodec interface
var _ ccipcommon.SourceChainExtraDataCodec = &extraDataDecoder{}
