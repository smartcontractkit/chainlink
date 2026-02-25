package generic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

func TestNewOracleFactory(t *testing.T) {
	params := OracleFactoryParams{
		JobID:   1,
		JobName: "test-job",
		Logger:  logger.TestLogger(t),
	}

	factory, err := NewOracleFactory(params)
	require.NoError(t, err)
	require.NotNil(t, factory)
}

func TestNewOracleFactory_WithOCRConfigService(t *testing.T) {
	mockService := &mockOCRConfigService{}

	params := OracleFactoryParams{
		JobID:            1,
		JobName:          "test-job",
		Logger:           logger.TestLogger(t),
		OCRConfigService: mockService,
		CapabilityID:     "offchain_reporting@1.0.0",
	}

	factory, err := NewOracleFactory(params)
	require.NoError(t, err)
	require.NotNil(t, factory)

	of := factory.(*oracleFactory)
	assert.Equal(t, mockService, of.ocrConfigService)
	assert.Equal(t, "offchain_reporting@1.0.0", of.capabilityID)
}

// mockOCRConfigService is a minimal mock for testing OracleFactory setup.
type mockOCRConfigService struct {
	services.Service
	registrysyncer.Listener
}

func (m *mockOCRConfigService) GetConfigTracker(
	capabilityID string,
	ocrConfigKey string,
	legacyTracker ocrtypes.ContractConfigTracker,
) (ocrtypes.ContractConfigTracker, error) {
	return &mockConfigTracker{}, nil
}

func (m *mockOCRConfigService) GetConfigDigester(
	capabilityID string,
	ocrConfigKey string,
	legacyDigester ocrtypes.OffchainConfigDigester,
) (ocrtypes.OffchainConfigDigester, error) {
	return &mockConfigDigester{}, nil
}

type mockConfigTracker struct{}

func (m *mockConfigTracker) Notify() <-chan struct{} { return nil }
func (m *mockConfigTracker) LatestConfigDetails(ctx context.Context) (uint64, ocrtypes.ConfigDigest, error) {
	return 0, ocrtypes.ConfigDigest{}, nil
}
func (m *mockConfigTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	return ocrtypes.ContractConfig{}, nil
}
func (m *mockConfigTracker) LatestBlockHeight(ctx context.Context) (uint64, error) {
	return 0, nil
}

type mockConfigDigester struct{}

func (m *mockConfigDigester) ConfigDigest(ctx context.Context, cc ocrtypes.ContractConfig) (ocrtypes.ConfigDigest, error) {
	return ocrtypes.ConfigDigest{}, nil
}
func (m *mockConfigDigester) ConfigDigestPrefix(ctx context.Context) (ocrtypes.ConfigDigestPrefix, error) {
	return ocrtypes.ConfigDigestPrefixKeystoneOCR3Capability, nil
}
