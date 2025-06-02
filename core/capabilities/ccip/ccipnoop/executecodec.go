package ccipnoop

import (
	"context"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// NoopExecutePluginCodecV1 is a codec for encoding and decoding execute plugin reports.
// Compatible with:
// - "OffRamp 1.6.0-dev"
type NoopExecutePluginCodecV1 struct {
	extraDataCodec common.ExtraDataCodec
}

func NewNoopExecutePluginCodecV1(extraDataCodec common.ExtraDataCodec) *NoopExecutePluginCodecV1 {
	return &NoopExecutePluginCodecV1{
		extraDataCodec: extraDataCodec,
	}
}

func (e *NoopExecutePluginCodecV1) Encode(ctx context.Context, report cciptypes.ExecutePluginReport) ([]byte, error) {
	return []byte{}, nil
}

func (e *NoopExecutePluginCodecV1) Decode(ctx context.Context, encodedReport []byte) (cciptypes.ExecutePluginReport, error) {
	return cciptypes.ExecutePluginReport{}, nil
}

// Ensure ExecutePluginCodec implements the ExecutePluginCodec interface
var _ cciptypes.ExecutePluginCodec = (*NoopExecutePluginCodecV1)(nil)
