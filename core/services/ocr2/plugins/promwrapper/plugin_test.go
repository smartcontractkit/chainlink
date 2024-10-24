package promwrapper

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/promwrapper/mocks"
)

var (
	// Intra-phase latencies.
	qDuration = time.Millisecond * 100 // duration of Query()
	oDuration = time.Millisecond * 200 // duration of Observation()
	rDuration = time.Millisecond * 300 // duration of Report()
	aDuration = time.Millisecond * 400 // duration of ShouldAcceptFinalizedReport()
	tDuration = time.Millisecond * 500 // duration of ShouldTransmitAcceptedReport()
	cDuration = time.Millisecond * 600 // duration of Close()

	// Inter-phase latencies.
	qToOLatency = time.Millisecond * 100 // latency between Query() and Observation()
	oToRLatency = time.Millisecond * 200 // latency between Observation() and Report()
	rToALatency = time.Millisecond * 300 // latency between Report() and ShouldAcceptFinalizedReport()
	aToTLatency = time.Millisecond * 400 // latency between ShouldAcceptFinalizedReport() and ShouldTransmitAcceptedReport()

	ceiling = time.Millisecond * 700
)

// fakeReportingPlugin has varied intra-phase latencies.
type fakeReportingPlugin struct{}

func (fakeReportingPlugin) Query(context.Context, types.ReportTimestamp) (types.Query, error) {
	time.Sleep(qDuration)
	return nil, nil
}
func (fakeReportingPlugin) Observation(context.Context, types.ReportTimestamp, types.Query) (types.Observation, error) {
	time.Sleep(oDuration)
	return nil, nil
}
func (fakeReportingPlugin) Report(context.Context, types.ReportTimestamp, types.Query, []types.AttributedObservation) (bool, types.Report, error) {
	time.Sleep(rDuration)
	return false, nil, nil
}
func (fakeReportingPlugin) ShouldAcceptFinalizedReport(context.Context, types.ReportTimestamp, types.Report) (bool, error) {
	time.Sleep(aDuration)
	return false, nil
}
func (fakeReportingPlugin) ShouldTransmitAcceptedReport(context.Context, types.ReportTimestamp, types.Report) (bool, error) {
	time.Sleep(tDuration)
	return false, nil
}
func (fakeReportingPlugin) Close() error {
	time.Sleep(cDuration)
	return nil
}

var _ types.ReportingPlugin = &fakeReportingPlugin{}

func TestPlugin_MustInstantiate(t *testing.T) {
	// Ensure instantiation without panic for no override backend.
	var reportingPlugin = &fakeReportingPlugin{}
	promPlugin := New(reportingPlugin, "test", "EVM", big.NewInt(1), types.ReportingPluginConfig{}, nil)
	require.NotNil(t, promPlugin)

	// Ensure instantiation without panic for override provided.
	backend := mocks.NewPrometheusBackend(t)
	promPlugin = New(reportingPlugin, "test-2", "EVM", big.NewInt(1), types.ReportingPluginConfig{}, backend)
	require.NotNil(t, promPlugin)
}

func TestPlugin_GetLatencies(t *testing.T) {
	// Use arbitrary report timestamp and label values.
	configDigest := common.BytesToHash(crypto.Keccak256([]byte("foobar")))
	reportTimestamp := types.ReportTimestamp{
		ConfigDigest: types.ConfigDigest(configDigest),
		Epoch:        1,
		Round:        1,
	}
	var assertCorrectLabelValues = func(labelValues []string) {
		require.Equal(
			t,
			[]string{
				"EVM",
				"1",
				"test-plugin",
				"0",
				common.Bytes2Hex(configDigest[:]),
			}, labelValues)
	}

	// Instantiate prometheus backend mock.
	backend := mocks.NewPrometheusBackend(t)

	// Assert intra-phase latencies.
	backend.On("SetQueryDuration", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		labelValues := args[0].([]string)
		duration := time.Duration(args[1].(float64))
		assertCorrectLabelValues(labelValues)
		require.Greater(t, duration, qDuration)
		require.Less(t, duration, oDuration)
	}).Return()
	backend.On("SetObservationDuration", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		labelValues := args[0].([]string)
		duration := time.Duration(args[1].(float64))
		assertCorrectLabelValues(labelValues)
		require.Greater(t, duration, oDuration)
		require.Less(t, duration, rDuration)
	}).Return()
	backend.On("SetReportDuration", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		labelValues := args[0].([]string)
		duration := time.Duration(args[1].(float64))
		assertCorrectLabelValues(labelValues)
		require.Greater(t, duration, rDuration)
		require.Less(t, duration, aDuration)
	}).Return()
	backend.On("SetShouldAcceptFinalizedReportDuration", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		labelValues := args[0].([]string)
		duration := time.Duration(args[1].(float64))
		assertCorrectLabelValues(labelValues)
		require.Greater(t, duration, aDuration)
		require.Less(t, duration, tDuration)
	}).Return()
	backend.On("SetShouldTransmitAcceptedReportDuration", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		labelValues := args[0].([]string)
		duration := time.Duration(args[1].(float64))
		assertCorrectLabelValues(labelValues)
		require.Greater(t, duration, tDuration)
		require.Less(t, duration, cDuration)
	}).Return()

	// Assert inter-phase latencies.
	backend.On("SetQueryToObservationLatency", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		labelValues := args[0].([]string)
		latency := time.Duration(args[1].(float64))
		assertCorrectLabelValues(labelValues)
		require.Greater(t, latency, qToOLatency)
		require.Less(t, latency, oToRLatency)
	}).Return()
	backend.On("SetObservationToReportLatency", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		labelValues := args[0].([]string)
		latency := time.Duration(args[1].(float64))
		assertCorrectLabelValues(labelValues)
		require.Greater(t, latency, oToRLatency)
		require.Less(t, latency, rToALatency)
	}).Return()
	backend.On("SetReportToAcceptFinalizedReportLatency", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		labelValues := args[0].([]string)
		latency := time.Duration(args[1].(float64))
		assertCorrectLabelValues(labelValues)
		require.Greater(t, latency, rToALatency)
		require.Less(t, latency, aToTLatency)
	}).Return()
	backend.On("SetAcceptFinalizedReportToTransmitAcceptedReportLatency", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		labelValues := args[0].([]string)
		latency := time.Duration(args[1].(float64))
		assertCorrectLabelValues(labelValues)
		require.Greater(t, latency, aToTLatency)
		require.Less(t, latency, cDuration)
	}).Return()

	// Assert close correctly reported.
	backend.On("SetCloseDuration", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		labelValues := args[0].([]string)
		latency := time.Duration(args[1].(float64))
		require.Equal(
			t,
			[]string{
				"EVM",
				"1",
				"test-plugin",
				"0",
				common.Bytes2Hex(configDigest[:]),
			}, labelValues)
		require.Greater(t, latency, cDuration)
		require.Less(t, latency, ceiling)
	}).Return()

	// Create promPlugin with mocked prometheus backend.
	var reportingPlugin = &fakeReportingPlugin{}
	var promPlugin *promPlugin = New(
		reportingPlugin,
		"test-plugin",
		"EVM",
		big.NewInt(1),
		types.ReportingPluginConfig{ConfigDigest: reportTimestamp.ConfigDigest},
		backend,
	).(*promPlugin)
	require.NotNil(t, promPlugin)

	ctx := testutils.Context(t)

	// Run OCR methods.
	_, err := promPlugin.Query(ctx, reportTimestamp)
	require.NoError(t, err)
	_, ok := promPlugin.queryEndTimes.Get(timestampToKey(reportTimestamp))
	require.True(t, ok)
	time.Sleep(qToOLatency)

	_, err = promPlugin.Observation(ctx, reportTimestamp, nil)
	require.NoError(t, err)
	_, ok = promPlugin.observationEndTimes.Get(timestampToKey(reportTimestamp))
	require.True(t, ok)
	time.Sleep(oToRLatency)

	_, _, err = promPlugin.Report(ctx, reportTimestamp, nil, nil)
	require.NoError(t, err)
	_, ok = promPlugin.reportEndTimes.Get(timestampToKey(reportTimestamp))
	require.True(t, ok)
	time.Sleep(rToALatency)

	_, err = promPlugin.ShouldAcceptFinalizedReport(ctx, reportTimestamp, nil)
	require.NoError(t, err)
	_, ok = promPlugin.acceptFinalizedReportEndTimes.Get(timestampToKey(reportTimestamp))
	require.True(t, ok)
	time.Sleep(aToTLatency)

	_, err = promPlugin.ShouldTransmitAcceptedReport(ctx, reportTimestamp, nil)
	require.NoError(t, err)

	// Close.
	err = promPlugin.Close()
	require.NoError(t, err)
}
