package ccipnoop

import (
	"context"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// ExecutePluginCodecV1 is a codec for encoding and decoding execute plugin reports.
// Compatible with:
// - "OffRamp 1.6.0-dev"
type ExecutePluginCodecV1 struct {
	extraDataCodec common.ExtraDataCodec
}

func NewExecutePluginCodecV1(extraDataCodec common.ExtraDataCodec) *ExecutePluginCodecV1 {
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
