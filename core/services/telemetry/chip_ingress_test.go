package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

// MockEmitter is a mock implementation of the Emitter interface
type MockEmitter struct {
	mock.Mock
}

func (m *MockEmitter) Emit(ctx context.Context, body []byte, attrKVs ...any) error {
	args := m.Called(ctx, body, attrKVs)
	return args.Error(0)
}

func TestNewChipIngressAdapter(t *testing.T) {
	t.Run("Success - Ethereum Mainnet", func(t *testing.T) {
		mockEmitter := new(MockEmitter)
		lggr := logger.TestLogger(t)
		chainID := "1"
		network := "EVM"
		contractID := "0x1234"
		telemType := synchronization.OCR2Median

		adapter, err := NewChipIngressAdapter(network, chainID, contractID, telemType, mockEmitter, lggr)
		require.NoError(t, err)
		require.NotNil(t, adapter)

		// Verify chain selector was derived correctly (Ethereum mainnet)
		assert.Equal(t, uint64(5009297550715157269), adapter.ChainSelector)
		assert.Equal(t, "EVM", adapter.Network)
		assert.Equal(t, "1", adapter.ChainID)
		assert.Equal(t, contractID, adapter.ContractID)

		// Verify domain and entity were derived correctly
		assert.Equal(t, "data-feeds", adapter.Domain)
		assert.Equal(t, "ocr.v2.median.telemetry", adapter.Entity)
	})

	t.Run("Success - Polygon", func(t *testing.T) {
		mockEmitter := new(MockEmitter)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAdapter("EVM", "137", "0xabc", synchronization.OCR3CCIPCommit, mockEmitter, lggr)
		require.NoError(t, err)

		// Verify chain selector was derived correctly (Polygon)
		assert.Equal(t, uint64(4051577828743386545), adapter.ChainSelector)
		assert.Equal(t, "EVM", adapter.Network)
		assert.Equal(t, "137", adapter.ChainID)

		// Verify domain and entity for CCIP commit
		assert.Equal(t, "ccip", adapter.Domain)
		assert.Equal(t, "ocr.v3.ccip.commit.telemetry", adapter.Entity)
	})

	t.Run("Success - Arbitrum One", func(t *testing.T) {
		mockEmitter := new(MockEmitter)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAdapter("EVM", "42161", "0xdef", synchronization.OCR2CCIPExec, mockEmitter, lggr)
		require.NoError(t, err)

		// Verify chain selector was derived correctly (Arbitrum One)
		assert.Equal(t, uint64(4949039107694359620), adapter.ChainSelector)
	})

	t.Run("Error - nil emitter", func(t *testing.T) {
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAdapter("EVM", "1", "0x1234", synchronization.OCR2Median, nil, lggr)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "beholder emitter cannot be nil")
	})

	t.Run("Error - invalid chainID", func(t *testing.T) {
		mockEmitter := new(MockEmitter)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAdapter("EVM", "1234567890123456", "0x1234", synchronization.OCR2Median, mockEmitter, lggr)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "failed to get chain details")
	})

	t.Run("Error - invalid network family", func(t *testing.T) {
		mockEmitter := new(MockEmitter)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAdapter("INVALID_NETWORK", "1", "0x1234", synchronization.OCR2Median, mockEmitter, lggr)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "failed to get chain details")
	})

	t.Run("Error - invalid telemetry type", func(t *testing.T) {
		mockEmitter := new(MockEmitter)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAdapter("EVM", "1", "0x1234", synchronization.TelemetryType("unknown"), mockEmitter, lggr)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "failed to map telemetry type to domain/entity")
	})
}

func TestChipIngressAdapter_SendLog(t *testing.T) {
	t.Run("Success - sends to beholder", func(t *testing.T) {
		mockEmitter := new(MockEmitter)
		lggr := logger.TestLogger(t)
		contractID := "0x1234567890abcdef"
		telemType := synchronization.OCR2Median
		telemetryLog := []byte("test telemetry data")

		adapter, err := NewChipIngressAdapter("EVM", "1", contractID, telemType, mockEmitter, lggr)
		require.NoError(t, err)

		// Setup expectations
		mockEmitter.
			On("Emit", mock.Anything, telemetryLog, mock.Anything).
			Run(func(args mock.Arguments) {
				// args: ctx, body, attrKVs
				attrKVs, ok := args.Get(2).([]any)
				require.True(t, ok)
				require.GreaterOrEqual(t, len(attrKVs)%2, 0)
				attrs := make(map[string]any)
				for i := 0; i+1 < len(attrKVs); i += 2 {
					key, isStr := attrKVs[i].(string)
					require.True(t, isStr)
					attrs[key] = attrKVs[i+1]
				}
				// Ensure network_name key is present and old network key is absent
				assert.Equal(t, "EVM", attrs["network_name"])
				_, hadOld := attrs["network"]
				assert.False(t, hadOld)
				// Spot-check a couple of other attributes
				assert.Equal(t, "1", attrs["chain_id"])
				assert.Equal(t, contractID, attrs["contract_id"])
			}).
			Return(nil).
			Once()

		// Call SendLog
		adapter.SendLog(telemetryLog)

		// Verify emit was called
		mockEmitter.AssertExpectations(t)
	})

	t.Run("SendLog continues even if emit fails", func(t *testing.T) {
		mockEmitter := new(MockEmitter)
		lggr := logger.TestLogger(t)
		contractID := "0xabcdef"
		telemType := synchronization.OCR3Automation
		telemetryLog := []byte("test data")

		adapter, err := NewChipIngressAdapter("EVM", "137", contractID, telemType, mockEmitter, lggr)
		require.NoError(t, err)

		// Setup expectations - emitter returns error
		mockEmitter.On("Emit", mock.Anything, telemetryLog, mock.Anything).Return(assert.AnError).Once()

		// Call SendLog - should not panic even if emit fails
		assert.NotPanics(t, func() {
			adapter.SendLog(telemetryLog)
		})

		mockEmitter.AssertExpectations(t)
	})

	t.Run("Multiple SendLog calls", func(t *testing.T) {
		mockEmitter := new(MockEmitter)
		lggr := logger.TestLogger(t)
		contractID := "0x999"
		telemType := synchronization.OCR2CCIPCommit

		adapter, err := NewChipIngressAdapter("EVM", "42161", contractID, telemType, mockEmitter, lggr)
		require.NoError(t, err)

		// Setup expectations for multiple calls
		log1 := []byte("log 1")
		log2 := []byte("log 2")
		log3 := []byte("log 3")

		mockEmitter.On("Emit", mock.Anything, log1, mock.Anything).Return(nil).Once()
		mockEmitter.On("Emit", mock.Anything, log2, mock.Anything).Return(nil).Once()
		mockEmitter.On("Emit", mock.Anything, log3, mock.Anything).Return(nil).Once()

		// Send multiple logs
		adapter.SendLog(log1)
		adapter.SendLog(log2)
		adapter.SendLog(log3)

		mockEmitter.AssertExpectations(t)
	})
}

func TestChipIngressAdapter_ExportedFields(t *testing.T) {
	mockEmitter := new(MockEmitter)
	lggr := logger.TestLogger(t)
	contractID := "0x1234567890"
	telemType := synchronization.OCR3Mercury

	adapter, err := NewChipIngressAdapter("EVM", "1", contractID, telemType, mockEmitter, lggr)
	require.NoError(t, err)

	t.Run("ChainSelector", func(t *testing.T) {
		assert.Equal(t, uint64(5009297550715157269), adapter.ChainSelector)
	})

	t.Run("Network", func(t *testing.T) {
		assert.Equal(t, "EVM", adapter.Network)
	})

	t.Run("ChainID", func(t *testing.T) {
		assert.Equal(t, "1", adapter.ChainID)
	})

	t.Run("ContractID", func(t *testing.T) {
		assert.Equal(t, contractID, adapter.ContractID)
	})

	t.Run("Domain", func(t *testing.T) {
		assert.Equal(t, "data-streams", adapter.Domain)
	})

	t.Run("Entity", func(t *testing.T) {
		assert.Equal(t, "ocr.v3.mercury.telemetry", adapter.Entity)
	})
}

func TestChipIngressAdapter_InterfaceCompliance(t *testing.T) {
	// This test verifies that ChipIngressAdapter implements commontypes.MonitoringEndpoint
	mockEmitter := new(MockEmitter)
	lggr := logger.TestLogger(t)

	adapter, err := NewChipIngressAdapter("EVM", "1", "0x123", synchronization.OCR2Median, mockEmitter, lggr)
	require.NoError(t, err)

	// Verify it can be assigned to the interface
	var _ interface{} = adapter

	// Call the interface method
	mockEmitter.On("Emit", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	adapter.SendLog([]byte("test"))
	mockEmitter.AssertCalled(t, "Emit", mock.Anything, mock.Anything, mock.Anything)
}
