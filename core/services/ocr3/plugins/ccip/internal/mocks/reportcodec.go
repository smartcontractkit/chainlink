package mocks

import (
	"context"
	"encoding/json"

	"github.com/smartcontractkit/ccipocr3/internal/model"
)

type CommitPluginJSONReportCodec struct{}

func NewCommitPluginJSONReportCodec() *CommitPluginJSONReportCodec {
	return &CommitPluginJSONReportCodec{}
}

func (c CommitPluginJSONReportCodec) Encode(ctx context.Context, report model.CommitPluginReport) ([]byte, error) {
	return json.Marshal(report)
}

func (c CommitPluginJSONReportCodec) Decode(ctx context.Context, bytes []byte) (model.CommitPluginReport, error) {
	report := model.CommitPluginReport{}
	err := json.Unmarshal(bytes, &report)
	return report, err
}
