package mercury_common_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
)

// MercuryPlugin is a test helper for testing [ocr3types.MercuryPlugin] implementations.
// the name is adopted because ocr3 has a special mercury plugin type
var OCR3Plugin = staticMercuryPlugin{
	staticMercuryPluginConfig: staticMercuryPluginConfig{
		observationRequest: observationRequest{
			reportTimestamp: reportContext.ReportTimestamp,
			previousReport:  previousReport,
		},
		observationResponse: observationResponse{
			observation: libocr.Observation{1, 2, 3},
		},
		reportRequest: reportRequest{
			reportTimestamp: reportContext.ReportTimestamp,
			previousReport:  previousReport,
			observations:    obs,
		},
		reportResponse: reportResponse{
			shouldReport: true,
			report:       report,
		},
	},
}

type observationRequest struct {
	reportTimestamp libocr.ReportTimestamp
	previousReport  libocr.Report
}

type observationResponse struct {
	observation libocr.Observation
}

type reportRequest struct {
	reportTimestamp libocr.ReportTimestamp
	previousReport  libocr.Report
	observations    []libocr.AttributedObservation
}

type reportResponse struct {
	shouldReport bool
	report       libocr.Report
}
type staticMercuryPluginConfig struct {
	observationRequest
	observationResponse
	reportRequest
	reportResponse
}

type staticMercuryPlugin struct {
	staticMercuryPluginConfig
}

var _ ocr3types.MercuryPlugin = staticMercuryPlugin{}
var _ testtypes.AssertEqualer[ocr3types.MercuryPlugin] = staticMercuryPlugin{}

func (s staticMercuryPlugin) Observation(ctx context.Context, timestamp libocr.ReportTimestamp, previousReport libocr.Report) (libocr.Observation, error) {
	if timestamp != s.observationRequest.reportTimestamp {
		return nil, fmt.Errorf("expected report timestamp %v but got %v", s.observationRequest.reportTimestamp, timestamp)
	}
	if !bytes.Equal(previousReport, s.observationRequest.previousReport) {
		return nil, fmt.Errorf("expected previous report %x but got %x", s.observationRequest.previousReport, previousReport)
	}
	return s.observationResponse.observation, nil
}

func (s staticMercuryPlugin) Report(ctx context.Context, timestamp libocr.ReportTimestamp, previousReport libocr.Report, observations []libocr.AttributedObservation) (bool, libocr.Report, error) {
	if timestamp != s.reportRequest.reportTimestamp {
		return false, nil, fmt.Errorf("expected report timestamp %v but got %v", s.reportRequest.reportTimestamp, timestamp)
	}
	if !bytes.Equal(s.reportRequest.previousReport, previousReport) {
		return false, nil, fmt.Errorf("expected previous report %x but got %x", s.reportRequest.previousReport, previousReport)
	}
	if !assert.ObjectsAreEqual(s.reportRequest.observations, observations) {
		return false, nil, fmt.Errorf("expected %v but got %v", s.reportRequest.observations, observations)
	}
	return s.reportResponse.shouldReport, s.reportResponse.report, nil
}

func (s staticMercuryPlugin) Close() error { return nil }

func (s staticMercuryPlugin) AssertEqual(ctx context.Context, t *testing.T, other ocr3types.MercuryPlugin) {
	gotObs, err := other.Observation(ctx, s.observationRequest.reportTimestamp, s.observationRequest.previousReport)
	require.NoError(t, err)
	assert.Equal(t, s.observationResponse.observation, gotObs)
	gotOk, gotReport, err := other.Report(ctx, s.reportRequest.reportTimestamp, s.reportRequest.previousReport, s.reportRequest.observations)
	require.NoError(t, err)
	assert.Equal(t, s.reportResponse.shouldReport, gotOk)
	assert.Equal(t, s.reportResponse.report, gotReport)
}
