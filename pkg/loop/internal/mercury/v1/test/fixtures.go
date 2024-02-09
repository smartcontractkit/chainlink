package v1_test

import (
	"math/big"

	ocr2plus_types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	mercury_v1_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
)

type Parameters struct {
	// ReportCodec
	Report          ocr2plus_types.Report
	ReportFields    mercury_v1_types.ReportFields
	MaxReportLength int
	CurrentBlockNum int64

	// DataSource
	ReportTimestamp ocr2plus_types.ReportTimestamp
	Observation     mercury_v1_types.Observation
}

var Fixtures = Parameters{
	// ReportCodec
	Report: ocr2plus_types.Report([]byte("mercury v1 report")),
	ReportFields: mercury_v1_types.ReportFields{
		Timestamp:             0,
		BenchmarkPrice:        big.NewInt(5),
		Ask:                   big.NewInt(6),
		Bid:                   big.NewInt(7),
		CurrentBlockNum:       8,
		CurrentBlockHash:      []byte("mercury v1 current block hash"),
		ValidFromBlockNum:     9,
		CurrentBlockTimestamp: 10,
	},
	MaxReportLength: 20,
	CurrentBlockNum: 23,

	// DataSource
	ReportTimestamp: ocr2plus_types.ReportTimestamp{
		ConfigDigest: [32]byte([]byte("mercury v1 configuration digest!")),
		Epoch:        0,
		Round:        1,
	},
	Observation: mercury_v1_types.Observation{
		BenchmarkPrice:          mercury.ObsResult[*big.Int]{Val: big.NewInt(50)},
		Ask:                     mercury.ObsResult[*big.Int]{Val: big.NewInt(60)},
		Bid:                     mercury.ObsResult[*big.Int]{Val: big.NewInt(70)},
		CurrentBlockNum:         mercury.ObsResult[int64]{Val: 80},
		CurrentBlockHash:        mercury.ObsResult[[]byte]{Val: []byte("mercury v1 test block hash")},
		CurrentBlockTimestamp:   mercury.ObsResult[uint64]{Val: 90},
		LatestBlocks:            []mercury_v1_types.Block{{Num: 100, Hash: "fakehash", Ts: 101}, {Num: 102, Hash: "fakehash2", Ts: 103}},
		MaxFinalizedBlockNumber: mercury.ObsResult[int64]{Val: 79},
	},
}
