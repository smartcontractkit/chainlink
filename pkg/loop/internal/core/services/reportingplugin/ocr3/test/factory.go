package ocr3_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

var Factory = ocr3staticPluginFactory{
	ReportingPluginConfig: ocr3reportingPluginConfig,
	reportingPlugin:       ReportingPlugin,
}

// OCR3
type ocr3staticPluginFactory struct {
	ocr3types.ReportingPluginConfig
	reportingPlugin ocr3staticReportingPlugin
}

var _ types.OCR3ReportingPluginFactory = (*ocr3staticPluginFactory)(nil)

func (o ocr3staticPluginFactory) Name() string { panic("implement me") }

func (o ocr3staticPluginFactory) Start(ctx context.Context) error { return nil }

func (o ocr3staticPluginFactory) Close() error { return nil }

func (o ocr3staticPluginFactory) Ready() error { panic("implement me") }

func (o ocr3staticPluginFactory) HealthReport() map[string]error { panic("implement me") }

func (o ocr3staticPluginFactory) NewReportingPlugin(config ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[[]byte], ocr3types.ReportingPluginInfo, error) {
	err := o.equalConfig(config)
	if err != nil {
		return nil, ocr3types.ReportingPluginInfo{}, fmt.Errorf("config mismatch: %w", err)
	}
	return o.reportingPlugin, ocr3rpi, nil
}

func (o ocr3staticPluginFactory) equalConfig(other ocr3types.ReportingPluginConfig) error {
	if other.ConfigDigest != o.ConfigDigest {
		return fmt.Errorf("expected ConfigDigest %x but got %x", o.ConfigDigest, other.ConfigDigest)
	}
	if other.OracleID != o.OracleID {
		return fmt.Errorf("expected OracleID %d but got %d", o.OracleID, other.OracleID)
	}
	if other.F != o.F {
		return fmt.Errorf("expected F %d but got %d", o.F, other.F)
	}
	if other.N != o.N {
		return fmt.Errorf("expected N %d but got %d", o.N, other.N)
	}
	if !bytes.Equal(other.OnchainConfig, o.OnchainConfig) {
		return fmt.Errorf("expected OnchainConfig %x but got %x", o.OnchainConfig, other.OnchainConfig)
	}
	if !bytes.Equal(other.OffchainConfig, o.OffchainConfig) {
		return fmt.Errorf("expected OffchainConfig %x but got %x", o.OffchainConfig, other.OffchainConfig)
	}
	if other.EstimatedRoundInterval != o.EstimatedRoundInterval {
		return fmt.Errorf("expected EstimatedRoundInterval %d but got %d", o.EstimatedRoundInterval, other.EstimatedRoundInterval)
	}
	if other.MaxDurationQuery != o.MaxDurationQuery {
		return fmt.Errorf("expected MaxDurationQuery %d but got %d", o.MaxDurationQuery, other.MaxDurationQuery)
	}
	if other.MaxDurationObservation != o.MaxDurationObservation {
		return fmt.Errorf("expected MaxDurationObservation %d but got %d", o.MaxDurationObservation, other.MaxDurationObservation)
	}
	if other.MaxDurationShouldAcceptAttestedReport != o.MaxDurationShouldAcceptAttestedReport {
		return fmt.Errorf("expected MaxDurationShouldAcceptAttestedReport %d but got %d", o.MaxDurationShouldAcceptAttestedReport, other.MaxDurationShouldAcceptAttestedReport)
	}
	if other.MaxDurationShouldTransmitAcceptedReport != o.MaxDurationShouldTransmitAcceptedReport {
		return fmt.Errorf("expected MaxDurationShouldTransmitAcceptedReport %d but got %d", o.MaxDurationShouldTransmitAcceptedReport, other.MaxDurationShouldTransmitAcceptedReport)
	}
	return nil
}

func OCR3ReportingPluginFactory(t *testing.T, factory types.OCR3ReportingPluginFactory) {
	expectedFactory := Factory
	t.Run("OCR3ReportingPluginFactory", func(t *testing.T) {
		ctx := tests.Context(t)
		rp, gotRPI, err := factory.NewReportingPlugin(ocr3reportingPluginConfig)
		require.NoError(t, err)
		assert.Equal(t, ocr3rpi, gotRPI)
		t.Cleanup(func() { assert.NoError(t, rp.Close()) })
		t.Run("OCR3ReportingPlugin", func(t *testing.T) {
			expectedFactory.reportingPlugin.AssertEqual(ctx, t, rp)
		})
	})
}

var (
	//OCR3
	ocr3reportingPluginConfig = ocr3types.ReportingPluginConfig{
		ConfigDigest:                            libocr.ConfigDigest([32]byte{1: 1, 3: 3, 5: 5}),
		OracleID:                                commontypes.OracleID(10),
		N:                                       12,
		F:                                       42,
		OnchainConfig:                           []byte{17: 11},
		OffchainConfig:                          []byte{32: 64},
		EstimatedRoundInterval:                  time.Second,
		MaxDurationQuery:                        time.Hour,
		MaxDurationObservation:                  time.Millisecond,
		MaxDurationShouldAcceptAttestedReport:   10 * time.Second,
		MaxDurationShouldTransmitAcceptedReport: time.Minute,
	}

	ocr3rpi = ocr3types.ReportingPluginInfo{
		Name: "test",
		Limits: ocr3types.ReportingPluginLimits{
			MaxQueryLength:       42,
			MaxObservationLength: 13,
			MaxOutcomeLength:     33,
			MaxReportLength:      17,
			MaxReportCount:       41,
		},
	}
)
