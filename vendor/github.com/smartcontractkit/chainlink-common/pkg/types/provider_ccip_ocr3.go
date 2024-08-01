package types

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

type CCIPOCR3CommitProvider interface {
	PluginProvider

	ReportCodec(ctx context.Context) (ccipocr3.CommitPluginCodec, error)
	MsgHasher(ctx context.Context) (ccipocr3.MessageHasher, error)
}

type CCIPOCR3ExecuteProvider interface {
	PluginProvider

	ReportCodec(ctx context.Context) (ccipocr3.ExecutePluginCodec, error)
	MsgHasher(ctx context.Context) (ccipocr3.MessageHasher, error)
}
