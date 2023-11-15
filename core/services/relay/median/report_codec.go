package median

import (
	"context"
	"math/big"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

const MedianTypeName = "MedianReport"

func NewReportCodec(codec types.Codec) (median.ReportCodec, error) {
	if codec == nil {
		return nil, errors.New("codec cannot be nil")
	}
	return &reportCodec{codec: codec}, nil
}

type reportCodec struct {
	codec types.Codec
}

var _ median.ReportCodec = reportCodec{}

func (r reportCodec) BuildReport(observations []median.ParsedAttributedObservation) (ocrtypes.Report, error) {
	agg := aggregate(observations)
	return r.codec.Encode(context.Background(), agg, MedianTypeName)
}

func (r reportCodec) MedianFromReport(report ocrtypes.Report) (*big.Int, error) {
	agg := &AggregatedAttributedObservation{}
	if err := r.codec.Decode(context.Background(), report, agg, MedianTypeName); err != nil {
		return nil, err
	}
	medianObservation := len(agg.Observations) / 2
	return agg.Observations[medianObservation], nil
}

func (r reportCodec) MaxReportLength(n int) (int, error) {
	return r.codec.GetMaxDecodingSize(context.Background(), n, MedianTypeName)
}
