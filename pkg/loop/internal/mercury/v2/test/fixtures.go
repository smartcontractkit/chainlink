package v2_test

import (
	"math/big"

	ocr2plus_types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	mercury_v2_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
)

type Parameters struct {
	// ReportCodec
	Report               ocr2plus_types.Report
	ReportFields         mercury_v2_types.ReportFields
	MaxReportLength      int
	ObservationTimestamp uint32

	// DataSource
	ReportTimestamp ocr2plus_types.ReportTimestamp
	Observation     mercury_v2_types.Observation
}

var Fixtures = Parameters{
	// ReportCodec
	Report: ocr2plus_types.Report([]byte("mercury v2 report")),
	ReportFields: mercury_v2_types.ReportFields{
		ValidFromTimestamp: 0,
		Timestamp:          1,
		NativeFee:          big.NewInt(2),
		LinkFee:            big.NewInt(3),
		ExpiresAt:          4,
		BenchmarkPrice:     big.NewInt(5),
	},
	MaxReportLength:      20,
	ObservationTimestamp: 23,

	// DataSource
	ReportTimestamp: ocr2plus_types.ReportTimestamp{
		ConfigDigest: [32]byte([]byte("mercury v2 configuration digest!")),
		Epoch:        0,
		Round:        1,
	},
	Observation: mercury_v2_types.Observation{
		BenchmarkPrice:        mercury.ObsResult[*big.Int]{Val: big.NewInt(50)},
		MaxFinalizedTimestamp: mercury.ObsResult[int64]{Val: 79},
		LinkPrice:             mercury.ObsResult[*big.Int]{Val: big.NewInt(30)},
		NativePrice:           mercury.ObsResult[*big.Int]{Val: big.NewInt(20)},
	},
}
