package median

import (
	"math/big"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type ReportCodec struct {
	ChainReader types.ChainReader
}

var _ median.ReportCodec = &ReportCodec{}

func (r *ReportCodec) BuildReport(observations []median.ParsedAttributedObservation) (ocrtypes.Report, error) {
	//TODO aggregate and send the report off to encode
	panic("implement me")
}

func (r *ReportCodec) MedianFromReport(report ocrtypes.Report) (*big.Int, error) {
	// TODO decode and return the median
	panic("implement me")
}

func (r *ReportCodec) MaxReportLength(n int) (int, error) {
	// TODO?
	panic("implement me")
}
