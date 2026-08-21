package ccipnoop

import (
	"context"

	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// CommitPluginCodecV1 is a codec for encoding and decoding commit plugin reports.
// Compatible with:
// - "OffRamp 1.6.0-dev"
type CommitPluginCodecV1 struct{}

func NewCommitPluginCodecV1() *CommitPluginCodecV1 {
	return &CommitPluginCodecV1{}
}

func (c *CommitPluginCodecV1) Encode(ctx context.Context, report ccipocr3common.CommitPluginReport) ([]byte, error) {
	return []byte{}, nil
}

func (c *CommitPluginCodecV1) Decode(ctx context.Context, bytes []byte) (ccipocr3common.CommitPluginReport, error) {
	return ccipocr3common.CommitPluginReport{}, nil
}

// Ensure CommitPluginCodec implements the CommitPluginCodec interface
var _ ccipocr3common.CommitPluginCodec = (*CommitPluginCodecV1)(nil)
