package ccipevm

import (
	"context"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

var _ cciptypes.CommitPluginCodec = (*CommitPluginCodec)(nil)

type CommitPluginCodec struct{}

func NewCommitPluginCodec() *CommitPluginCodec {
	return &CommitPluginCodec{}
}

func (c *CommitPluginCodec) Encode(ctx context.Context, report cciptypes.CommitPluginReport) ([]byte, error) {
	panic("implement me")
}

func (c *CommitPluginCodec) Decode(ctx context.Context, bytes []byte) (cciptypes.CommitPluginReport, error) {
	panic("implement me")
}
