package mercury_common_test

import (
	"context"
	"testing"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"

	mercuryv1test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v1/test"
	mercuryv2test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v2/test"
	mercuryv3test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v3/test"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	mercurytypes "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	mercuryv1types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	mercuryv2types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	mercuryv3types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"

	ocr2test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ocr2/test"
	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
)

var MercuryProvider = staticMercuryProvider{
	staticMercuryProviderConfig: staticMercuryProviderConfig{
		offchainDigester:    ocr2test.OffchainConfigDigester,
		contractTracker:     ocr2test.ContractConfigTracker,
		contractTransmitter: ocr2test.ContractTransmitter,
		reportCodecV1:       mercuryv1test.ReportCodec,
		reportCodecV2:       mercuryv2test.ReportCodec,
		reportCodecV3:       mercuryv3test.ReportCodec,
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
	reportCodecV1       mercuryv1test.ReportCodecEvaluator
	reportCodecV2       mercuryv2test.ReportCodecEvaluator
	reportCodecV3       mercuryv3test.ReportCodecEvaluator
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

func (s staticMercuryProvider) ReportCodecV1() mercuryv1types.ReportCodec {
	return s.reportCodecV1
}

func (s staticMercuryProvider) ReportCodecV2() mercuryv2types.ReportCodec {
	return s.reportCodecV2
}

func (s staticMercuryProvider) ReportCodecV3() mercuryv3types.ReportCodec {
	return s.reportCodecV3
}

func (s staticMercuryProvider) OnchainConfigCodec() mercurytypes.OnchainConfigCodec {
	return s.onchainConfigCodec
}

func (s staticMercuryProvider) MercuryChainReader() mercurytypes.ChainReader {
	return s.mercuryChainReader
}

func (s staticMercuryProvider) ChainReader() types.ContractReader {
	//panic("mercury does not use the general ChainReader interface yet")
	return nil
}

func (s staticMercuryProvider) MercuryServerFetcher() mercurytypes.ServerFetcher {
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
