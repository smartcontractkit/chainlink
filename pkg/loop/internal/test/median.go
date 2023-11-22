package test

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func PluginMedian(t *testing.T, p types.PluginMedian) {
	PluginMedianTest{&StaticMedianProvider{}}.TestPluginMedian(t, p)
}

type PluginMedianTest struct {
	types.MedianProvider
}

func (m PluginMedianTest) TestPluginMedian(t *testing.T, p types.PluginMedian) {
	t.Run("PluginMedian", func(t *testing.T) {
		ctx := tests.Context(t)
		factory, err := p.NewMedianFactory(ctx, m.MedianProvider, &staticDataSource{value}, &staticDataSource{juelsPerFeeCoin}, &StaticErrorLog{})
		require.NoError(t, err)

		ReportingPluginFactory(t, factory)
	})
}

func ReportingPluginFactory(t *testing.T, factory types.ReportingPluginFactory) {
	t.Run("ReportingPluginFactory", func(t *testing.T) {
		rp, gotRPI, err := factory.NewReportingPlugin(reportingPluginConfig)
		require.NoError(t, err)
		assert.Equal(t, rpi, gotRPI)
		t.Cleanup(func() { assert.NoError(t, rp.Close()) })
		t.Run("ReportingPlugin", func(t *testing.T) {
			ctx := tests.Context(t)
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
		})
	})
}

type StaticPluginMedian struct{}

func (s StaticPluginMedian) NewMedianFactory(ctx context.Context, provider types.MedianProvider, dataSource, juelsPerFeeCoinDataSource median.DataSource, errorLog types.ErrorLog) (types.ReportingPluginFactory, error) {
	ocd := provider.OffchainConfigDigester()
	gotDigestPrefix, err := ocd.ConfigDigestPrefix()
	if err != nil {
		return nil, fmt.Errorf("failed to get ConfigDigestPrefix: %w", err)
	}
	if gotDigestPrefix != configDigestPrefix {
		return nil, fmt.Errorf("expected ConfigDigestPrefix %x but got %x", configDigestPrefix, gotDigestPrefix)
	}
	gotDigest, err := ocd.ConfigDigest(contractConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get ConfigDigest: %w", err)
	}
	if gotDigest != configDigest {
		return nil, fmt.Errorf("expected ConfigDigest %x but got %x", configDigest, gotDigest)
	}
	cct := provider.ContractConfigTracker()
	gotBlockHeight, err := cct.LatestBlockHeight(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get LatestBlockHeight: %w", err)
	}
	if gotBlockHeight != blockHeight {
		return nil, fmt.Errorf("expected LatestBlockHeight %d but got %d", blockHeight, gotBlockHeight)
	}
	gotChangedInBlock, gotConfigDigest, err := cct.LatestConfigDetails(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get LatestConfigDetails: %w", err)
	}
	if gotChangedInBlock != changedInBlock {
		return nil, fmt.Errorf("expected changedInBlock %d but got %d", changedInBlock, gotChangedInBlock)
	}
	if gotConfigDigest != configDigest {
		return nil, fmt.Errorf("expected ConfigDigest %s but got %s", configDigest, gotConfigDigest)
	}
	gotContractConfig, err := cct.LatestConfig(ctx, changedInBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to get LatestConfig: %w", err)
	}
	if !reflect.DeepEqual(gotContractConfig, contractConfig) {
		return nil, fmt.Errorf("expected ContractConfig %v but got %v", contractConfig, gotContractConfig)
	}
	ct := provider.ContractTransmitter()
	gotAccount, err := ct.FromAccount()
	if err != nil {
		return nil, fmt.Errorf("failed to get FromAccount: %w", err)
	}
	if gotAccount != account {
		return nil, fmt.Errorf("expectd FromAccount %s but got %s", account, gotAccount)
	}
	gotConfigDigest, gotEpoch, err := ct.LatestConfigDigestAndEpoch(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get LatestConfigDigestAndEpoch: %w", err)
	}
	if gotConfigDigest != configDigest {
		return nil, fmt.Errorf("expected ConfigDigest %s but got %s", configDigest, gotConfigDigest)
	}
	if gotEpoch != epoch {
		return nil, fmt.Errorf("expected Epoch %d but got %d", epoch, gotEpoch)
	}
	err = ct.Transmit(ctx, reportContext, report, sigs)
	if err != nil {
		return nil, fmt.Errorf("failed to Transmit")
	}
	rc := provider.ReportCodec()
	gotReport, err := rc.BuildReport(pobs)
	if err != nil {
		return nil, fmt.Errorf("failed to BuildReport: %w", err)
	}
	if !bytes.Equal(gotReport, report) {
		return nil, fmt.Errorf("expected Report %x but got %x", report, gotReport)
	}
	gotMedianValue, err := rc.MedianFromReport(report)
	if err != nil {
		return nil, fmt.Errorf("failed to get MedianFromReport: %w", err)
	}
	if medianValue.Cmp(gotMedianValue) != 0 {
		return nil, fmt.Errorf("expected MedianValue %s but got %s", medianValue, gotMedianValue)
	}
	gotMax, err := rc.MaxReportLength(n)
	if err != nil {
		return nil, fmt.Errorf("failed to get MaxReportLength: %w", err)
	}
	if gotMax != max {
		return nil, fmt.Errorf("expected MaxReportLength %d but got %d", max, gotMax)
	}
	mc := provider.MedianContract()
	gotConfigDigest, gotEpoch, gotRound, err := mc.LatestRoundRequested(ctx, lookbackDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to get LatestRoundRequested: %w", err)
	}
	if gotConfigDigest != configDigest {
		return nil, fmt.Errorf("expected ConfigDigest %s but got %s", configDigest, gotConfigDigest)
	}
	if gotEpoch != epoch {
		return nil, fmt.Errorf("expected Epoch %d but got %d", epoch, gotEpoch)
	}
	if gotRound != round {
		return nil, fmt.Errorf("expected Round %d but got %d", round, gotRound)
	}
	gotConfigDigest, gotEpoch, gotRound, gotLatestAnswer, gotLatestTimestamp, err := mc.LatestTransmissionDetails(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get LatestTransmissionDetails: %w", err)
	}
	if gotConfigDigest != configDigest {
		return nil, fmt.Errorf("expected ConfigDigest %s but got %s", configDigest, gotConfigDigest)
	}
	if gotEpoch != epoch {
		return nil, fmt.Errorf("expected Epoch %d but got %d", epoch, gotEpoch)
	}
	if gotRound != round {
		return nil, fmt.Errorf("expected Round %d but got %d", round, gotRound)
	}
	if latestAnswer.Cmp(gotLatestAnswer) != 0 {
		return nil, fmt.Errorf("expected LatestAnswer %s but got %s", latestAnswer, gotLatestAnswer)
	}
	if !gotLatestTimestamp.Equal(latestTimestamp) {
		return nil, fmt.Errorf("expected LatestTimestamp %s but got %s", latestTimestamp, gotLatestTimestamp)
	}
	occ := provider.OnchainConfigCodec()
	gotEncoded, err := occ.Encode(onchainConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to Encode: %w", err)
	}
	if !bytes.Equal(gotEncoded, encoded) {
		return nil, fmt.Errorf("expected Encoded %s but got %s", encoded, gotEncoded)
	}
	gotDecoded, err := occ.Decode(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to Decode: %w", err)
	}
	if !reflect.DeepEqual(gotDecoded, onchainConfig) {
		return nil, fmt.Errorf("expected OnchainConfig %s but got %s", onchainConfig, gotDecoded)
	}
	gotVal, err := dataSource.Observe(ctx, reportContext.ReportTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to observe dataSource: %w", err)
	}
	if !assert.ObjectsAreEqual(value, gotVal) {
		return nil, fmt.Errorf("expected Value %s but got %s", value, gotVal)
	}
	gotJuels, err := juelsPerFeeCoinDataSource.Observe(ctx, reportContext.ReportTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to observe juelsPerFeeCoin: %w", err)
	}
	if !assert.ObjectsAreEqual(juelsPerFeeCoin, gotJuels) {
		return nil, fmt.Errorf("expected JuelsPerFeeCoin %s but got %s", juelsPerFeeCoin, gotJuels)
	}
	if err := errorLog.SaveError(ctx, errMsg); err != nil {
		return nil, fmt.Errorf("failed to save error: %w", err)
	}
	return staticPluginFactory{}, nil
}

type staticPluginFactory struct{}

func (s staticPluginFactory) Name() string { panic("implement me") }

func (s staticPluginFactory) Start(ctx context.Context) error { return nil }

func (s staticPluginFactory) Close() error { return nil }

func (s staticPluginFactory) Ready() error { panic("implement me") }

func (s staticPluginFactory) HealthReport() map[string]error { panic("implement me") }

func (s staticPluginFactory) NewReportingPlugin(config libocr.ReportingPluginConfig) (libocr.ReportingPlugin, libocr.ReportingPluginInfo, error) {
	if config.ConfigDigest != reportingPluginConfig.ConfigDigest {
		return nil, libocr.ReportingPluginInfo{}, fmt.Errorf("expected ConfigDigest %x but got %x", reportingPluginConfig.ConfigDigest, config.ConfigDigest)
	}
	if config.OracleID != reportingPluginConfig.OracleID {
		return nil, libocr.ReportingPluginInfo{}, fmt.Errorf("expected OracleID %d but got %d", reportingPluginConfig.OracleID, config.OracleID)
	}
	if config.F != reportingPluginConfig.F {
		return nil, libocr.ReportingPluginInfo{}, fmt.Errorf("expected F %d but got %d", reportingPluginConfig.F, config.F)
	}
	if config.N != reportingPluginConfig.N {
		return nil, libocr.ReportingPluginInfo{}, fmt.Errorf("expected N %d but got %d", reportingPluginConfig.N, config.N)
	}
	if !bytes.Equal(config.OnchainConfig, reportingPluginConfig.OnchainConfig) {
		return nil, libocr.ReportingPluginInfo{}, fmt.Errorf("expected OnchainConfig %x but got %x", reportingPluginConfig.OnchainConfig, config.OnchainConfig)
	}
	if !bytes.Equal(config.OffchainConfig, reportingPluginConfig.OffchainConfig) {
		return nil, libocr.ReportingPluginInfo{}, fmt.Errorf("expected OffchainConfig %x but got %x", reportingPluginConfig.OffchainConfig, config.OffchainConfig)
	}
	if config.EstimatedRoundInterval != reportingPluginConfig.EstimatedRoundInterval {
		return nil, libocr.ReportingPluginInfo{}, fmt.Errorf("expected EstimatedRoundInterval %d but got %d", reportingPluginConfig.EstimatedRoundInterval, config.EstimatedRoundInterval)
	}
	if config.MaxDurationQuery != reportingPluginConfig.MaxDurationQuery {
		return nil, libocr.ReportingPluginInfo{}, fmt.Errorf("expected MaxDurationQuery %d but got %d", reportingPluginConfig.MaxDurationQuery, config.MaxDurationQuery)
	}
	if config.MaxDurationReport != reportingPluginConfig.MaxDurationReport {
		return nil, libocr.ReportingPluginInfo{}, fmt.Errorf("expected MaxDurationReport %d but got %d", reportingPluginConfig.MaxDurationReport, config.MaxDurationReport)
	}
	if config.MaxDurationObservation != reportingPluginConfig.MaxDurationObservation {
		return nil, libocr.ReportingPluginInfo{}, fmt.Errorf("expected MaxDurationObservation %d but got %d", reportingPluginConfig.MaxDurationObservation, config.MaxDurationObservation)
	}
	if config.MaxDurationShouldAcceptFinalizedReport != reportingPluginConfig.MaxDurationShouldAcceptFinalizedReport {
		return nil, libocr.ReportingPluginInfo{}, fmt.Errorf("expected MaxDurationShouldAcceptFinalizedReport %d but got %d", reportingPluginConfig.MaxDurationShouldAcceptFinalizedReport, config.MaxDurationShouldAcceptFinalizedReport)
	}
	if config.MaxDurationShouldTransmitAcceptedReport != reportingPluginConfig.MaxDurationShouldTransmitAcceptedReport {
		return nil, libocr.ReportingPluginInfo{}, fmt.Errorf("expected MaxDurationShouldTransmitAcceptedReport %d but got %d", reportingPluginConfig.MaxDurationShouldTransmitAcceptedReport, config.MaxDurationShouldTransmitAcceptedReport)
	}
	return staticReportingPlugin{}, rpi, nil
}

type StaticMedianProvider struct{}

func (s StaticMedianProvider) Start(ctx context.Context) error { return nil }

func (s StaticMedianProvider) Close() error { return nil }

func (s StaticMedianProvider) Ready() error { panic("unimplemented") }

func (s StaticMedianProvider) Name() string { panic("unimplemented") }

func (s StaticMedianProvider) HealthReport() map[string]error { panic("unimplemented") }

func (s StaticMedianProvider) OffchainConfigDigester() libocr.OffchainConfigDigester {
	return staticOffchainConfigDigester{}
}

func (s StaticMedianProvider) ContractConfigTracker() libocr.ContractConfigTracker {
	return staticContractConfigTracker{}
}

func (s StaticMedianProvider) ContractTransmitter() libocr.ContractTransmitter {
	return staticContractTransmitter{}
}

func (s StaticMedianProvider) ReportCodec() median.ReportCodec { return staticReportCodec{} }

func (s StaticMedianProvider) MedianContract() median.MedianContract { return staticMedianContract{} }

func (s StaticMedianProvider) OnchainConfigCodec() median.OnchainConfigCodec {
	return staticOnchainConfigCodec{}
}

type staticReportCodec struct{}

func (s staticReportCodec) BuildReport(os []median.ParsedAttributedObservation) (libocr.Report, error) {
	if !assert.ObjectsAreEqual(pobs, os) {
		return nil, fmt.Errorf("expected observations %v but got %v", pobs, os)
	}
	return report, nil
}

func (s staticReportCodec) MedianFromReport(r libocr.Report) (*big.Int, error) {
	if !bytes.Equal(report, r) {
		return nil, fmt.Errorf("expected report %x but got %x", report, r)
	}
	return medianValue, nil
}

func (s staticReportCodec) MaxReportLength(n2 int) (int, error) {
	if n != n2 {
		return -1, fmt.Errorf("expected n %d but got %d", n, n2)
	}
	return max, nil
}

type staticMedianContract struct{}

func (s staticMedianContract) LatestTransmissionDetails(ctx context.Context) (libocr.ConfigDigest, uint32, uint8, *big.Int, time.Time, error) {
	return configDigest, epoch, round, latestAnswer, latestTimestamp, nil
}

func (s staticMedianContract) LatestRoundRequested(ctx context.Context, lookback time.Duration) (libocr.ConfigDigest, uint32, uint8, error) {
	if lookbackDuration != lookback {
		return libocr.ConfigDigest{}, 0, 0, fmt.Errorf("expected lookback %s but got %s", lookbackDuration, lookback)
	}
	return configDigest, epoch, round, nil
}

type staticOnchainConfigCodec struct{}

func (s staticOnchainConfigCodec) Encode(c median.OnchainConfig) ([]byte, error) {
	if !assert.ObjectsAreEqual(onchainConfig.Max, c.Max) {
		return nil, fmt.Errorf("expected max %s but got %s", onchainConfig.Max, c.Max)
	}
	if !assert.ObjectsAreEqual(onchainConfig.Min, c.Min) {
		return nil, fmt.Errorf("expected min %s but got %s", onchainConfig.Min, c.Min)
	}
	return encoded, nil
}

func (s staticOnchainConfigCodec) Decode(b []byte) (median.OnchainConfig, error) {
	if !bytes.Equal(encoded, b) {
		return median.OnchainConfig{}, fmt.Errorf("expected encoded %x but got %x", encoded, b)
	}
	return onchainConfig, nil
}
