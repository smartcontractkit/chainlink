package median

import (
	"context"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

const typeName = "MedianReport"

type reportCodec struct {
	codec types.Codec
}

var _ median.ReportCodec = &reportCodec{}

func (r *reportCodec) BuildReport(observations []median.ParsedAttributedObservation) (ocrtypes.Report, error) {
	if len(observations) == 0 {
		return nil, fmt.Errorf("cannot build report from empty attributed observations")
	}

	return r.codec.Encode(context.Background(), aggregate(observations), typeName)
}

func (r *reportCodec) MedianFromReport(report ocrtypes.Report) (*big.Int, error) {
	agg := &aggregatedAttributedObservation{}
	if err := r.codec.Decode(context.Background(), report, agg, typeName); err != nil {
		return nil, err
	}
	observations := make([]*big.Int, len(agg.Observations))
	copy(observations, agg.Observations)
	medianObservation := len(agg.Observations) / 2
	return agg.Observations[medianObservation], nil
}

func (r *reportCodec) MaxReportLength(n int) (int, error) {
	return r.codec.GetMaxDecodingSize(context.Background(), n, typeName)
}
