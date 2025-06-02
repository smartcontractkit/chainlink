package ccipnoop

import (
	"context"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// NoopCommitPluginCodecV1 is a codec for encoding and decoding commit plugin reports.
// Compatible with:
// - "OffRamp 1.6.0-dev"
type NoopCommitPluginCodecV1 struct{}

func NewNoopCommitPluginCodecV1() *NoopCommitPluginCodecV1 {
	return &NoopCommitPluginCodecV1{}
}

func (c *NoopCommitPluginCodecV1) Encode(ctx context.Context, report cciptypes.CommitPluginReport) ([]byte, error) {
	return []byte{}, nil
}

func (c *NoopCommitPluginCodecV1) Decode(ctx context.Context, bytes []byte) (cciptypes.CommitPluginReport, error) {
	return cciptypes.CommitPluginReport{}, nil
}

// Ensure CommitPluginCodec implements the CommitPluginCodec interface
var _ cciptypes.CommitPluginCodec = (*NoopCommitPluginCodecV1)(nil)
