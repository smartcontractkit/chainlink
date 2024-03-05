package mercury_common_test

import (
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var (
	configDigest = libocr.ConfigDigest([32]byte{2: 10, 12: 16})
	obs          = []libocr.AttributedObservation{{Observation: []byte{21: 19}, Observer: commontypes.OracleID(99)}}

	previousReport = libocr.Report([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	report         = libocr.Report{42: 101}
	reportContext  = libocr.ReportContext{
		ReportTimestamp: libocr.ReportTimestamp{
			ConfigDigest: libocr.ConfigDigest([32]byte{1: 7, 31: 3}),
			Epoch:        79,
			Round:        17,
		},
		ExtraHash: [32]byte{1: 2, 3: 4, 5: 6},
	}
)

var (
	mercuryPluginConfig = ocr3types.MercuryPluginConfig{
		ConfigDigest:           configDigest,
		OracleID:               commontypes.OracleID(11),
		N:                      12,
		F:                      42,
		OnchainConfig:          []byte{17: 11},
		OffchainConfig:         []byte{32: 64},
		EstimatedRoundInterval: time.Second,
		MaxDurationObservation: time.Millisecond,
	}
	mercuryPluginInfo = ocr3types.MercuryPluginInfo{
		Name: "test",
		Limits: ocr3types.MercuryPluginLimits{
			MaxObservationLength: 13,
			MaxReportLength:      17,
		},
	}
)
