package mercury_common_test

import (
	"context"
	"testing"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"

	mercury_v1_test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v1/test"
	mercury_v2_test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v2/test"
	mercury_v3_test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v3/test"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	mercury_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	mercury_v1_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	mercury_v2_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	mercury_v3_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"

	testpluginprovider "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/ocr2/plugin_provider"
	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
)

var MercuryProvider = staticMercuryProvider{
	staticMercuryProviderConfig: staticMercuryProviderConfig{
		offchainDigester:    testpluginprovider.OffchainConfigDigester,
		contractTracker:     testpluginprovider.ContractConfigTracker,
		contractTransmitter: testpluginprovider.ContractTransmitter,
		reportCodecV1:       mercury_v1_test.ReportCodec,
		reportCodecV2:       mercury_v2_test.ReportCodec,
		reportCodecV3:       mercury_v3_test.ReportCodec,
		onchainConfigCodec:  OnchainConfigCodec,
		mercuryChainReader:  ChainReader,
		serviceFetcher:      ServerFetcher,
	},
}

type MercuryProviderTester interface {
	types.MercuryProvider
	AssertEqual(ctx context.Context, t *testing.T, other types.MercuryProvider)
}

type staticMercuryProviderConfig struct {
	// we use the static implementation type not the interface type
	// because we always expect the static implementation to be used
	// and it facilitates testing.
	offchainDigester    testtypes.OffchainConfigDigesterEvaluator
	contractTracker     testtypes.ContractConfigTrackerEvaluator
	contractTransmitter testtypes.ContractTransmitterEvaluator
	reportCodecV1       mercury_v1_test.ReportCodecEvaluator
	reportCodecV2       mercury_v2_test.ReportCodecEvaluator
	reportCodecV3       mercury_v3_test.ReportCodecEvaluator
	onchainConfigCodec  OnchainConfigCodecEvaluator
	mercuryChainReader  MercuryChainReaderEvaluator
	serviceFetcher      ServerFetcherEvaluator
}

var _ types.MercuryProvider = staticMercuryProvider{}

type staticMercuryProvider struct {
	staticMercuryProviderConfig
}

func (s staticMercuryProvider) Start(ctx context.Context) error { return nil }

func (s staticMercuryProvider) Close() error { return nil }

func (s staticMercuryProvider) Ready() error { panic("unimplemented") }

func (s staticMercuryProvider) Name() string { panic("unimplemented") }

func (s staticMercuryProvider) HealthReport() map[string]error { panic("unimplemented") }

func (s staticMercuryProvider) OffchainConfigDigester() libocr.OffchainConfigDigester {
	return s.offchainDigester
}

func (s staticMercuryProvider) ContractConfigTracker() libocr.ContractConfigTracker {
	return s.contractTracker
}

func (s staticMercuryProvider) ContractTransmitter() libocr.ContractTransmitter {
	return s.contractTransmitter
}

func (s staticMercuryProvider) ReportCodecV1() mercury_v1_types.ReportCodec {
	return s.reportCodecV1
}

func (s staticMercuryProvider) ReportCodecV2() mercury_v2_types.ReportCodec {
	return s.reportCodecV2
}

func (s staticMercuryProvider) ReportCodecV3() mercury_v3_types.ReportCodec {
	return s.reportCodecV3
}

func (s staticMercuryProvider) OnchainConfigCodec() mercury_types.OnchainConfigCodec {
	return s.onchainConfigCodec
}

func (s staticMercuryProvider) MercuryChainReader() mercury_types.ChainReader {
	return s.mercuryChainReader
}

func (s staticMercuryProvider) ChainReader() types.ChainReader {
	//panic("mercury does not use the general ChainReader interface yet")
	return nil
}

func (s staticMercuryProvider) MercuryServerFetcher() mercury_types.ServerFetcher {
	return s.serviceFetcher
}

func (s staticMercuryProvider) Codec() types.Codec {
	return nil
}

func (s staticMercuryProvider) AssertEqual(ctx context.Context, t *testing.T, other types.MercuryProvider) {
	t.Run("OffchainConfigDigester", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.offchainDigester.Evaluate(ctx, other.OffchainConfigDigester()))
	})
	t.Run("ContractConfigTracker", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.contractTracker.Evaluate(ctx, other.ContractConfigTracker()))
	})
	t.Run("ContractTransmitter", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.contractTransmitter.Evaluate(ctx, other.ContractTransmitter()))
	})
	t.Run("ReportCodecV1", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.reportCodecV1.Evaluate(ctx, other.ReportCodecV1()))
	})
	t.Run("ReportCodecV2", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.reportCodecV2.Evaluate(ctx, other.ReportCodecV2()))
	})
	t.Run("ReportCodecV3", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.reportCodecV3.Evaluate(ctx, other.ReportCodecV3()))
	})
	t.Run("OnchainConfigCodec", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.onchainConfigCodec.Evaluate(ctx, other.OnchainConfigCodec()))
	})
	t.Run("MercuryChainReader", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.mercuryChainReader.Evaluate(ctx, other.MercuryChainReader()))
	})
	t.Run("MercuryServerFetcher", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.serviceFetcher.Evaluate(ctx, other.MercuryServerFetcher()))
	})
}
