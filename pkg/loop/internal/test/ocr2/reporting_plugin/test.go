package reportingplugin

import (
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

const (
	epoch          = uint32(88)
	round          = uint8(74)
	shouldAccept   = true
	shouldReport   = true
	shouldTransmit = true
)

var (
	// ReportingPlugin is a static implementation of the ReportingPluginTester interface for testing
	ReportingPlugin = staticReportingPlugin{
		staticReportingPluginConfig: staticReportingPluginConfig{
			ReportContext:          reportContext,
			Query:                  query,
			Observation:            observation,
			AttributedObservations: obs,
			Report:                 report,
			ShouldAccept:           shouldAccept,
			ShouldReport:           shouldReport,
			ShouldTransmit:         shouldTransmit,
		},
	}

	configDigest = libocr.ConfigDigest([32]byte{2: 10, 12: 16})

	observation = libocr.Observation([]byte{21: 19})
	obs         = []libocr.AttributedObservation{{Observation: []byte{21: 19}, Observer: commontypes.OracleID(99)}}

	query = []byte{42: 42}

	report        = libocr.Report{42: 101}
	reportContext = libocr.ReportContext{
		ReportTimestamp: libocr.ReportTimestamp{
			ConfigDigest: configDigest,
			Epoch:        epoch,
			Round:        round,
		},
		ExtraHash: [32]byte{1: 2, 3: 4, 5: 6},
	}

	reportingPluginConfig = libocr.ReportingPluginConfig{
		ConfigDigest:                            libocr.ConfigDigest{}, // testpluginprovider.ConfigDigest,
		OracleID:                                commontypes.OracleID(10),
		N:                                       12,
		F:                                       42,
		OnchainConfig:                           []byte{17: 11},
		OffchainConfig:                          []byte{32: 64},
		EstimatedRoundInterval:                  time.Second,
		MaxDurationQuery:                        time.Hour,
		MaxDurationObservation:                  time.Millisecond,
		MaxDurationReport:                       time.Microsecond,
		MaxDurationShouldAcceptFinalizedReport:  10 * time.Second,
		MaxDurationShouldTransmitAcceptedReport: time.Minute,
	}

	rpi = libocr.ReportingPluginInfo{
		Name:          "test",
		UniqueReports: true,
		Limits: libocr.ReportingPluginLimits{
			MaxQueryLength:       42,
			MaxObservationLength: 13,
			MaxReportLength:      17,
		},
	}
)
