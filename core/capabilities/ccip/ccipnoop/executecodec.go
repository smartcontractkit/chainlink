package ccipnoop

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// ExecutePluginCodecV1 is a codec for encoding and decoding execute plugin reports.
// Compatible with:
// - "OffRamp 1.6.0-dev"
type ExecutePluginCodecV1 struct {
	extraDataCodec ccipocr3.ExtraDataCodecBundle
}

func NewExecutePluginCodecV1(extraDataCodec ccipocr3.ExtraDataCodecBundle) *ExecutePluginCodecV1 {
	return &ExecutePluginCodecV1{
		extraDataCodec: extraDataCodec,
	}
}

func (e *ExecutePluginCodecV1) Encode(ctx context.Context, report ccipocr3.ExecutePluginReport) ([]byte, error) {
	return []byte{}, nil
}

func (e *ExecutePluginCodecV1) Decode(ctx context.Context, encodedReport []byte) (ccipocr3.ExecutePluginReport, error) {
	return ccipocr3.ExecutePluginReport{}, nil
}

// Ensure ExecutePluginCodec implements the ExecutePluginCodec interface
var _ ccipocr3.ExecutePluginCodec = (*ExecutePluginCodecV1)(nil)
