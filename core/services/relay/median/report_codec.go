package median

import (
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"math/big"
)

type ReportCodec struct {
	// TODO update all the chainlink products then reference my commit with the things I need
}

func (r ReportCodec) BuildReport(observations []median.ParsedAttributedObservation) (ocrtypes.Report, error) {
	//TODO aggregate and send the report off to encode
	panic("implement me")
}

func (r ReportCodec) MedianFromReport(report ocrtypes.Report) (*big.Int, error) {
	// TODO decode and return the median
	panic("implement me")
}

func (r ReportCodec) MaxReportLength(n int) (int, error) {
	// TODO?
	panic("implement me")
}

var _ median.ReportCodec = ReportCodec{}
