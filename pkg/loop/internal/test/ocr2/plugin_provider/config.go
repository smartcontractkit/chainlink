package pluginprovider

import (
	"context"
	"testing"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/stretchr/testify/assert"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// ConfigProviderTester is a helper interface for testing ConfigProviders
type ConfigProviderTester interface {
	types.ConfigProvider
	// AssertEqual checks that the sub-components of the other ConfigProvider are equal to this one
	AssertEqual(ctx context.Context, t *testing.T, other types.ConfigProvider)
}

type staticConfigProviderConfig struct {
	offchainDigester      testtypes.OffchainConfigDigesterEvaluator
	contractConfigTracker testtypes.ContractConfigTrackerEvaluator
}

// staticConfigProvider is a static implementation of ConfigProviderTester
type staticConfigProvider struct {
	staticConfigProviderConfig
}

var _ ConfigProviderTester = staticConfigProvider{}

// TODO validate start/Close calls?
func (s staticConfigProvider) Start(ctx context.Context) error { return nil }

func (s staticConfigProvider) Close() error { return nil }

func (s staticConfigProvider) Ready() error { panic("unimplemented") }

func (s staticConfigProvider) Name() string { panic("unimplemented") }

func (s staticConfigProvider) HealthReport() map[string]error { panic("unimplemented") }

func (s staticConfigProvider) OffchainConfigDigester() libocr.OffchainConfigDigester {
	return s.offchainDigester
}

func (s staticConfigProvider) ContractConfigTracker() libocr.ContractConfigTracker {
	return s.contractConfigTracker
}

func (s staticConfigProvider) AssertEqual(ctx context.Context, t *testing.T, cp types.ConfigProvider) {
	t.Run("OffchainConfigDigester", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.offchainDigester.Evaluate(ctx, cp.OffchainConfigDigester()))
	})
	t.Run("ContractConfigTracker", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, s.contractConfigTracker.Evaluate(context.Background(), cp.ContractConfigTracker()))
	})
}
