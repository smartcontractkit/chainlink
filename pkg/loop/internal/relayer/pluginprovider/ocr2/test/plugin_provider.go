package pluginprovider

import (
	"context"
	"testing"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"

	chainreadertest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/chainreader/test"
	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var _ types.PluginProvider = staticPluginProvider{}

var AgnosticProvider = staticPluginProvider{
	offchainConfigDigester: OffchainConfigDigester,
	contractConfigTracker:  ContractConfigTracker,
	contractTransmitter:    ContractTransmitter,
	chainReader:            chainreadertest.ChainReader,
	codec:                  chainreadertest.Codec,
}

// staticPluginProvider is a static implementation of PluginProviderTester
type staticPluginProvider struct {
	offchainConfigDigester staticOffchainConfigDigester
	contractConfigTracker  staticContractConfigTracker
	contractTransmitter    testtypes.ContractTransmitterEvaluator
	chainReader            testtypes.ChainReaderTester
	codec                  testtypes.CodecEvaluator
}

var _ testtypes.PluginProviderTester = staticPluginProvider{}

func (s staticPluginProvider) Start(ctx context.Context) error { return nil }

func (s staticPluginProvider) Close() error { return nil }

func (s staticPluginProvider) Ready() error { panic("unimplemented") }

func (s staticPluginProvider) Name() string { panic("unimplemented") }

func (s staticPluginProvider) HealthReport() map[string]error { panic("unimplemented") }

func (s staticPluginProvider) OffchainConfigDigester() libocr.OffchainConfigDigester {
	return s.offchainConfigDigester
}

func (s staticPluginProvider) ContractConfigTracker() libocr.ContractConfigTracker {
	return s.contractConfigTracker
}

func (s staticPluginProvider) ContractTransmitter() libocr.ContractTransmitter {
	return s.contractTransmitter
}

func (s staticPluginProvider) ChainReader() types.ContractReader {
	return s.chainReader
}

func (s staticPluginProvider) Codec() types.Codec {
	return s.codec
}

func (s staticPluginProvider) AssertEqual(ctx context.Context, t *testing.T, provider types.PluginProvider) {
	t.Run("OffchainConfigDigester", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.offchainConfigDigester.Evaluate(ctx, provider.OffchainConfigDigester()))
	})

	t.Run("ContractConfigTracker", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.contractConfigTracker.Evaluate(ctx, provider.ContractConfigTracker()))
	})

	t.Run("ContractTransmitter", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.contractTransmitter.Evaluate(ctx, provider.ContractTransmitter()))
	})

	t.Run("ChainReader", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.chainReader.Evaluate(ctx, provider.ChainReader()))
	})
}

func (s staticPluginProvider) Evaluate(ctx context.Context, provider types.PluginProvider) error {
	err := s.offchainConfigDigester.Evaluate(ctx, provider.OffchainConfigDigester())
	if err != nil {
		return err
	}

	err = s.contractConfigTracker.Evaluate(ctx, provider.ContractConfigTracker())
	if err != nil {
		return err
	}

	err = s.contractTransmitter.Evaluate(ctx, provider.ContractTransmitter())
	if err != nil {
		return err
	}

	err = s.chainReader.Evaluate(ctx, provider.ChainReader())
	if err != nil {
		return err
	}

	return nil
}
