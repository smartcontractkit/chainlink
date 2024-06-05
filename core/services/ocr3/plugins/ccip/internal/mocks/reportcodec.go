package mocks

import (
	"context"
	"encoding/json"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

type CommitPluginJSONReportCodec struct{}

func NewCommitPluginJSONReportCodec() *CommitPluginJSONReportCodec {
	return &CommitPluginJSONReportCodec{}
}

func (c CommitPluginJSONReportCodec) Encode(ctx context.Context, report cciptypes.CommitPluginReport) ([]byte, error) {
	return json.Marshal(report)
}

func (c CommitPluginJSONReportCodec) Decode(ctx context.Context, bytes []byte) (cciptypes.CommitPluginReport, error) {
	report := cciptypes.CommitPluginReport{}
	err := json.Unmarshal(bytes, &report)
	return report, err
}
