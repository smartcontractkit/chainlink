package test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type staticReportingPluginConfig struct {
	ReportContext          libocr.ReportContext
	Query                  libocr.Query
	Observation            libocr.Observation
	AttributedObservations []libocr.AttributedObservation
	Report                 libocr.Report
	ShouldReport           bool
	ShouldAccept           bool
	ShouldTransmit         bool
}

var _ libocr.ReportingPlugin = staticReportingPlugin{}

type staticReportingPlugin struct {
	staticReportingPluginConfig
}

func (s staticReportingPlugin) Query(ctx context.Context, timestamp libocr.ReportTimestamp) (libocr.Query, error) {
	if timestamp != s.staticReportingPluginConfig.ReportContext.ReportTimestamp {
		return nil, errExpected(s.staticReportingPluginConfig.ReportContext.ReportTimestamp, timestamp)
	}
	return s.staticReportingPluginConfig.Query, nil
}

func (s staticReportingPlugin) Observation(ctx context.Context, timestamp libocr.ReportTimestamp, q libocr.Query) (libocr.Observation, error) {
	if timestamp != s.staticReportingPluginConfig.ReportContext.ReportTimestamp {
		return nil, errExpected(s.staticReportingPluginConfig.ReportContext.ReportTimestamp, timestamp)
	}
	if !bytes.Equal(q, s.staticReportingPluginConfig.Query) {
		return nil, errExpected(s.staticReportingPluginConfig.Query, q)
	}
	return s.staticReportingPluginConfig.Observation, nil
}

func (s staticReportingPlugin) Report(ctx context.Context, timestamp libocr.ReportTimestamp, q libocr.Query, observations []libocr.AttributedObservation) (bool, libocr.Report, error) {
	if timestamp != s.staticReportingPluginConfig.ReportContext.ReportTimestamp {
		return false, nil, errExpected(s.staticReportingPluginConfig.ReportContext.ReportTimestamp, timestamp)
	}
	if !bytes.Equal(q, s.staticReportingPluginConfig.Query) {
		return false, nil, errExpected(s.staticReportingPluginConfig.Query, q)
	}
	if !assert.ObjectsAreEqual(s.staticReportingPluginConfig.AttributedObservations, observations) {
		return false, nil, errExpected(s.staticReportingPluginConfig.AttributedObservations, observations)
	}
	return s.staticReportingPluginConfig.ShouldReport, s.staticReportingPluginConfig.Report, nil
}

func (s staticReportingPlugin) ShouldAcceptFinalizedReport(ctx context.Context, timestamp libocr.ReportTimestamp, r libocr.Report) (bool, error) {
	if timestamp != s.staticReportingPluginConfig.ReportContext.ReportTimestamp {
		return false, errExpected(s.staticReportingPluginConfig.ReportContext.ReportTimestamp, timestamp)
	}
	if !bytes.Equal(r, s.staticReportingPluginConfig.Report) {
		return false, errExpected(s.staticReportingPluginConfig.Report, r)
	}
	return shouldAccept, nil
}

func (s staticReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, timestamp libocr.ReportTimestamp, r libocr.Report) (bool, error) {
	if timestamp != s.staticReportingPluginConfig.ReportContext.ReportTimestamp {
		return false, errExpected(s.staticReportingPluginConfig.ReportContext.ReportTimestamp, timestamp)
	}
	if !bytes.Equal(r, s.staticReportingPluginConfig.Report) {
		return false, errExpected(s.staticReportingPluginConfig.Report, r)
	}
	return shouldTransmit, nil
}

func (s staticReportingPlugin) Close() error { return nil }

func (s staticReportingPlugin) AssertEqual(ctx context.Context, t *testing.T, rp libocr.ReportingPlugin) {
	gotQuery, err := rp.Query(ctx, reportContext.ReportTimestamp)
	require.NoError(t, err)
	assert.Equal(t, query, []byte(gotQuery))
	gotObs, err := rp.Observation(ctx, reportContext.ReportTimestamp, query)
	require.NoError(t, err)
	assert.Equal(t, observation, gotObs)
	gotOk, gotReport, err := rp.Report(ctx, reportContext.ReportTimestamp, query, obs)
	require.NoError(t, err)
	assert.True(t, gotOk)
	assert.Equal(t, report, gotReport)
	gotShouldAccept, err := rp.ShouldAcceptFinalizedReport(ctx, reportContext.ReportTimestamp, report)
	require.NoError(t, err)
	assert.True(t, gotShouldAccept)
	gotShouldTransmit, err := rp.ShouldTransmitAcceptedReport(ctx, reportContext.ReportTimestamp, report)
	require.NoError(t, err)
	assert.True(t, gotShouldTransmit)
}

func errExpected(expected, got any) error {
	return fmt.Errorf("expected %v but got %v", expected, got)
}
