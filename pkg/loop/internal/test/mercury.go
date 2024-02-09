package test

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	mercury_common_test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/mercury/common/test"
	mercury_v1_test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/mercury/v1/test"
	mercury_v2_test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/mercury/v2/test"
	mercury_v3_test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/mercury/v3/test"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	mercury_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	mercury_v1_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	mercury_v2_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	mercury_v3_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
)

func PluginMercury(t *testing.T, p types.PluginMercury) {
	PluginMercuryTest{&StaticMercuryProvider{}}.TestPluginMercury(t, p)
}

type PluginMercuryTest struct {
	types.MercuryProvider
}

func (m PluginMercuryTest) TestPluginMercury(t *testing.T, p types.PluginMercury) {
	t.Run("PluginMercuryV3", func(t *testing.T) {
		ctx := tests.Context(t)
		factory, err := p.NewMercuryV3Factory(ctx, m.MercuryProvider, mercury_v3_test.StaticDataSource{})
		require.NoError(t, err)
		require.NotNil(t, factory)

		MercuryPluginFactory(t, factory)
	})

	t.Run("PluginMercuryV2", func(t *testing.T) {
		ctx := tests.Context(t)
		factory, err := p.NewMercuryV2Factory(ctx, m.MercuryProvider, mercury_v2_test.StaticDataSource{})
		require.NoError(t, err)
		require.NotNil(t, factory)

		MercuryPluginFactory(t, factory)
	})

	t.Run("PluginMercuryV1", func(t *testing.T) {
		ctx := tests.Context(t)
		factory, err := p.NewMercuryV1Factory(ctx, m.MercuryProvider, mercury_v1_test.StaticDataSource{})
		require.NoError(t, err)
		require.NotNil(t, factory)

		MercuryPluginFactory(t, factory)
	})
}

type StaticPluginMercury struct{}

var _ types.PluginMercury = StaticPluginMercury{}

func (s StaticPluginMercury) commonValidation(ctx context.Context, provider types.MercuryProvider) error {
	ocd := provider.OffchainConfigDigester()
	gotDigestPrefix, err := ocd.ConfigDigestPrefix()
	if err != nil {
		return fmt.Errorf("failed to get ConfigDigestPrefix: %w", err)
	}
	if gotDigestPrefix != configDigestPrefix {
		return fmt.Errorf("expected ConfigDigestPrefix %x but got %x", configDigestPrefix, gotDigestPrefix)
	}
	gotDigest, err := ocd.ConfigDigest(contractConfig)
	if err != nil {
		return fmt.Errorf("failed to get ConfigDigest: %w", err)
	}
	if gotDigest != configDigest {
		return fmt.Errorf("expected ConfigDigest %x but got %x", configDigest, gotDigest)
	}
	cct := provider.ContractConfigTracker()
	gotBlockHeight, err := cct.LatestBlockHeight(ctx)
	if err != nil {
		return fmt.Errorf("failed to get LatestBlockHeight: %w", err)
	}
	if gotBlockHeight != blockHeight {
		return fmt.Errorf("expected LatestBlockHeight %d but got %d", blockHeight, gotBlockHeight)
	}
	gotChangedInBlock, gotConfigDigest, err := cct.LatestConfigDetails(ctx)
	if err != nil {
		return fmt.Errorf("failed to get LatestConfigDetails: %w", err)
	}
	if gotChangedInBlock != changedInBlock {
		return fmt.Errorf("expected changedInBlock %d but got %d", changedInBlock, gotChangedInBlock)
	}
	if gotConfigDigest != configDigest {
		return fmt.Errorf("expected ConfigDigest %s but got %s", configDigest, gotConfigDigest)
	}
	gotContractConfig, err := cct.LatestConfig(ctx, changedInBlock)
	if err != nil {
		return fmt.Errorf("failed to get LatestConfig: %w", err)
	}
	if !reflect.DeepEqual(gotContractConfig, contractConfig) {
		return fmt.Errorf("expected ContractConfig %v but got %v", contractConfig, gotContractConfig)
	}
	ct := provider.ContractTransmitter()
	gotAccount, err := ct.FromAccount()
	if err != nil {
		return fmt.Errorf("failed to get FromAccount: %w", err)
	}
	if gotAccount != account {
		return fmt.Errorf("expectd FromAccount %s but got %s", account, gotAccount)
	}
	gotConfigDigest, gotEpoch, err := ct.LatestConfigDigestAndEpoch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get LatestConfigDigestAndEpoch: %w", err)
	}
	if gotConfigDigest != configDigest {
		return fmt.Errorf("expected ConfigDigest %s but got %s", configDigest, gotConfigDigest)
	}
	if gotEpoch != epoch {
		return fmt.Errorf("expected Epoch %d but got %d", epoch, gotEpoch)
	}
	err = ct.Transmit(ctx, reportContext, report, sigs)
	if err != nil {
		return fmt.Errorf("failed to Transmit")
	}

	occ := provider.OnchainConfigCodec()
	gotEncoded, err := occ.Encode(mercury_common_test.StaticOnChainConfigCodecFixtures.Decoded)
	if err != nil {
		return fmt.Errorf("failed to Encode: %w", err)
	}
	if !bytes.Equal(gotEncoded, mercury_common_test.StaticOnChainConfigCodecFixtures.Encoded) {
		return fmt.Errorf("expected Encoded %s but got %s", encoded, gotEncoded)
	}
	gotDecoded, err := occ.Decode(mercury_common_test.StaticOnChainConfigCodecFixtures.Encoded)
	if err != nil {
		return fmt.Errorf("failed to Decode: %w", err)
	}
	if !reflect.DeepEqual(gotDecoded, mercury_common_test.StaticOnChainConfigCodecFixtures.Decoded) {
		return fmt.Errorf("expected OnchainConfig %s but got %s", onchainConfig, gotDecoded)
	}
	return nil
}

func (s StaticPluginMercury) NewMercuryV3Factory(ctx context.Context, provider types.MercuryProvider, dataSource mercury_v3_types.DataSource) (types.MercuryPluginFactory, error) {
	var err error
	defer func() {
		if err != nil {
			panic(fmt.Sprintf("provider %v, %T: %s", provider, provider, err))
		}
	}()
	err = s.commonValidation(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("failed commonValidation: %w", err)
	}
	rc := provider.ReportCodecV3()
	gotReport, err := rc.BuildReport(mercury_v3_test.Fixtures.ReportFields)
	if err != nil {
		return nil, fmt.Errorf("failed to BuildReport: %w", err)
	}
	if !bytes.Equal(gotReport, mercury_v3_test.Fixtures.Report) {
		return nil, fmt.Errorf("expected Report %x but got %x", report, gotReport)
	}
	gotMax, err := rc.MaxReportLength(n)
	if err != nil {
		return nil, fmt.Errorf("failed to get MaxReportLength: %w", err)
	}
	if gotMax != mercury_v3_test.Fixtures.MaxReportLength {
		return nil, fmt.Errorf("expected MaxReportLength %d but got %d", max, gotMax)
	}
	gotObservedTimestamp, err := rc.ObservationTimestampFromReport(gotReport)
	if err != nil {
		return nil, fmt.Errorf("failed to get ObservationTimestampFromReport: %w", err)
	}
	if gotObservedTimestamp != mercury_v3_test.Fixtures.ObservationTimestamp {
		return nil, fmt.Errorf("expected ObservationTimestampFromReport %d but got %d", mercury_v3_test.Fixtures.ObservationTimestamp, gotObservedTimestamp)
	}

	gotVal, err := dataSource.Observe(ctx, mercury_v3_test.Fixtures.ReportTimestamp, false)
	if err != nil {
		return nil, fmt.Errorf("failed to observe dataSource: %w", err)
	}
	if !assert.ObjectsAreEqual(mercury_v3_test.Fixtures.Observation, gotVal) {
		return nil, fmt.Errorf("expected Value %v but got %v", value, gotVal)
	}

	return staticMercuryPluginFactory{}, nil
}

func (s StaticPluginMercury) NewMercuryV2Factory(ctx context.Context, provider types.MercuryProvider, dataSource mercury_v2_types.DataSource) (types.MercuryPluginFactory, error) {
	var err error
	defer func() {
		if err != nil {
			panic(fmt.Sprintf("provider %v, %T: %s", provider, provider, err))
		}
	}()
	err = s.commonValidation(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("failed commonValidation: %w", err)
	}
	rc := provider.ReportCodecV2()
	gotReport, err := rc.BuildReport(mercury_v2_test.Fixtures.ReportFields)
	if err != nil {
		return nil, fmt.Errorf("failed to BuildReport: %w", err)
	}
	if !bytes.Equal(gotReport, mercury_v2_test.Fixtures.Report) {
		return nil, fmt.Errorf("expected Report %x but got %x", report, gotReport)
	}
	gotMax, err := rc.MaxReportLength(n)
	if err != nil {
		return nil, fmt.Errorf("failed to get MaxReportLength: %w", err)
	}
	if gotMax != mercury_v2_test.Fixtures.MaxReportLength {
		return nil, fmt.Errorf("expected MaxReportLength %d but got %d", max, gotMax)
	}
	gotObservedTimestamp, err := rc.ObservationTimestampFromReport(gotReport)
	if err != nil {
		return nil, fmt.Errorf("failed to get ObservationTimestampFromReport: %w", err)
	}
	if gotObservedTimestamp != mercury_v2_test.Fixtures.ObservationTimestamp {
		return nil, fmt.Errorf("expected ObservationTimestampFromReport %d but got %d", mercury_v2_test.Fixtures.ObservationTimestamp, gotObservedTimestamp)
	}

	gotVal, err := dataSource.Observe(ctx, mercury_v2_test.Fixtures.ReportTimestamp, false)
	if err != nil {
		return nil, fmt.Errorf("failed to observe dataSource: %w", err)
	}
	if !assert.ObjectsAreEqual(mercury_v2_test.Fixtures.Observation, gotVal) {
		return nil, fmt.Errorf("expected Value %v but got %v", value, gotVal)
	}

	return staticMercuryPluginFactory{}, nil
}

func (s StaticPluginMercury) NewMercuryV1Factory(ctx context.Context, provider types.MercuryProvider, dataSource mercury_v1_types.DataSource) (types.MercuryPluginFactory, error) {
	var err error
	defer func() {
		if err != nil {
			panic(fmt.Sprintf("provider %v, %T: %s", provider, provider, err))
		}
	}()
	err = s.commonValidation(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("failed commonValidation: %w", err)
	}
	rc := provider.ReportCodecV1()
	gotReport, err := rc.BuildReport(mercury_v1_test.Fixtures.ReportFields)
	if err != nil {
		return nil, fmt.Errorf("failed to BuildReport: %w", err)
	}
	if !bytes.Equal(gotReport, mercury_v1_test.Fixtures.Report) {
		return nil, fmt.Errorf("expected Report %x but got %x", report, gotReport)
	}
	gotMax, err := rc.MaxReportLength(n)
	if err != nil {
		return nil, fmt.Errorf("failed to get MaxReportLength: %w", err)
	}
	if gotMax != mercury_v1_test.Fixtures.MaxReportLength {
		return nil, fmt.Errorf("expected MaxReportLength %d but got %d", max, gotMax)
	}
	gotCurrentBlockNum, err := rc.CurrentBlockNumFromReport(gotReport)
	if err != nil {
		return nil, fmt.Errorf("failed to get ObservationTimestampFromReport: %w", err)
	}
	if gotCurrentBlockNum != mercury_v1_test.Fixtures.CurrentBlockNum {
		return nil, fmt.Errorf("expected ObservationTimestampFromReport %d but got %d", mercury_v1_test.Fixtures.CurrentBlockNum, gotCurrentBlockNum)
	}

	gotVal, err := dataSource.Observe(ctx, mercury_v1_test.Fixtures.ReportTimestamp, false)
	if err != nil {
		return nil, fmt.Errorf("failed to observe dataSource: %w", err)
	}
	if !assert.ObjectsAreEqual(mercury_v1_test.Fixtures.Observation, gotVal) {
		return nil, fmt.Errorf("expected Value %v but got %v", value, gotVal)
	}

	return staticMercuryPluginFactory{}, nil
}

type StaticMercuryProvider struct{}

var _ types.MercuryProvider = StaticMercuryProvider{}

func (s StaticMercuryProvider) Start(ctx context.Context) error { return nil }

func (s StaticMercuryProvider) Close() error { return nil }

func (s StaticMercuryProvider) Ready() error { panic("unimplemented") }

func (s StaticMercuryProvider) Name() string { panic("unimplemented") }

func (s StaticMercuryProvider) HealthReport() map[string]error { panic("unimplemented") }

func (s StaticMercuryProvider) OffchainConfigDigester() libocr.OffchainConfigDigester {
	return staticOffchainConfigDigester{}
}

func (s StaticMercuryProvider) ContractConfigTracker() libocr.ContractConfigTracker {
	return staticContractConfigTracker{}
}

func (s StaticMercuryProvider) ContractTransmitter() libocr.ContractTransmitter {
	return staticContractTransmitter{}
}

func (s StaticMercuryProvider) ReportCodecV1() mercury_v1_types.ReportCodec {
	return mercury_v1_test.StaticReportCodec{}
}

func (s StaticMercuryProvider) ReportCodecV2() mercury_v2_types.ReportCodec {
	return mercury_v2_test.StaticReportCodec{}
}

func (s StaticMercuryProvider) ReportCodecV3() mercury_v3_types.ReportCodec {
	return mercury_v3_test.StaticReportCodec{}
}

func (s StaticMercuryProvider) OnchainConfigCodec() mercury_types.OnchainConfigCodec {
	return mercury_common_test.StaticOnchainConfigCodec{}
}

func (s StaticMercuryProvider) MercuryChainReader() mercury_types.ChainReader {
	return mercury_common_test.StaticMercuryChainReader{}
}

func (s StaticMercuryProvider) ChainReader() types.ChainReader {
	//panic("mercury does not use the general ChainReader interface yet")
	return nil
}

func (s StaticMercuryProvider) MercuryServerFetcher() mercury_types.ServerFetcher {
	return mercury_common_test.StaticServerFetcher{}
}

func (s StaticMercuryProvider) Codec() types.Codec {
	return nil
}

type staticMercuryPluginFactory struct{}

func (s staticMercuryPluginFactory) Name() string { panic("implement me") }

func (s staticMercuryPluginFactory) Start(ctx context.Context) error { return nil }

func (s staticMercuryPluginFactory) Close() error { return nil }

func (s staticMercuryPluginFactory) Ready() error { panic("implement me") }

func (s staticMercuryPluginFactory) HealthReport() map[string]error { panic("implement me") }

func (s staticMercuryPluginFactory) NewMercuryPlugin(config ocr3types.MercuryPluginConfig) (ocr3types.MercuryPlugin, ocr3types.MercuryPluginInfo, error) {
	if config.ConfigDigest != mercuryPluginConfig.ConfigDigest {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected ConfigDigest %x but got %x", mercuryPluginConfig.ConfigDigest, config.ConfigDigest)
	}
	if config.OracleID != mercuryPluginConfig.OracleID {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected OracleID %d but got %d", mercuryPluginConfig.OracleID, config.OracleID)
	}
	if config.F != mercuryPluginConfig.F {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected F %d but got %d", mercuryPluginConfig.F, config.F)
	}
	if config.N != mercuryPluginConfig.N {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected N %d but got %d", mercuryPluginConfig.N, config.N)
	}
	if !bytes.Equal(config.OnchainConfig, mercuryPluginConfig.OnchainConfig) {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected OnchainConfig %x but got %x", mercuryPluginConfig.OnchainConfig, config.OnchainConfig)
	}
	if !bytes.Equal(config.OffchainConfig, mercuryPluginConfig.OffchainConfig) {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected OffchainConfig %x but got %x", mercuryPluginConfig.OffchainConfig, config.OffchainConfig)
	}
	if config.EstimatedRoundInterval != mercuryPluginConfig.EstimatedRoundInterval {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected EstimatedRoundInterval %d but got %d", mercuryPluginConfig.EstimatedRoundInterval, config.EstimatedRoundInterval)
	}

	if config.MaxDurationObservation != mercuryPluginConfig.MaxDurationObservation {
		return nil, ocr3types.MercuryPluginInfo{}, fmt.Errorf("expected MaxDurationObservation %d but got %d", mercuryPluginConfig.MaxDurationObservation, config.MaxDurationObservation)
	}

	return staticMercuryPlugin{}, mercuryPluginInfo, nil
}

func MercuryPluginFactory(t *testing.T, factory types.MercuryPluginFactory) {
	t.Run("ReportingPluginFactory", func(t *testing.T) {
		rp, gotRPI, err := factory.NewMercuryPlugin(mercuryPluginConfig)
		require.NoError(t, err)
		assert.Equal(t, mercuryPluginInfo, gotRPI)
		t.Cleanup(func() { assert.NoError(t, rp.Close()) })
		t.Run("ReportingPlugin", func(t *testing.T) {
			ctx := tests.Context(t)
			gotObs, err := rp.Observation(ctx, reportContext.ReportTimestamp, query)
			require.NoError(t, err)
			assert.Equal(t, observation, gotObs)
			gotOk, gotReport, err := rp.Report(reportContext.ReportTimestamp, query, obs)
			require.NoError(t, err)
			assert.True(t, gotOk)
			assert.Equal(t, report, gotReport)
		})
	})
}

type staticMercuryPlugin struct{}

var _ ocr3types.MercuryPlugin = staticMercuryPlugin{}

func (s staticMercuryPlugin) Observation(ctx context.Context, timestamp libocr.ReportTimestamp, previousReport libocr.Report) (libocr.Observation, error) {
	if timestamp != reportContext.ReportTimestamp {
		return nil, fmt.Errorf("expected report timestamp %v but got %v", reportContext.ReportTimestamp, timestamp)
	}
	if !bytes.Equal(previousReport, query) {
		return nil, fmt.Errorf("expected previous report %x but got %x", query, previousReport)
	}
	return observation, nil
}

func (s staticMercuryPlugin) Report(timestamp libocr.ReportTimestamp, previousReport libocr.Report, observations []libocr.AttributedObservation) (bool, libocr.Report, error) {
	if timestamp != reportContext.ReportTimestamp {
		return false, nil, fmt.Errorf("expected report timestamp %v but got %v", reportContext.ReportTimestamp, timestamp)
	}
	if !bytes.Equal(query, previousReport) {
		return false, nil, fmt.Errorf("expected previous report %x but got %x", query, previousReport)
	}
	if !assert.ObjectsAreEqual(obs, observations) {
		return false, nil, fmt.Errorf("expected %v but got %v", obs, observations)
	}
	return shouldReport, report, nil
}

func (s staticMercuryPlugin) Close() error { return nil }

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
