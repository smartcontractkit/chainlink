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

// MockChipIngressService is a mock implementation of the ChipIngressService interface
type MockChipIngressService struct {
	mock.Mock
}

func (m *MockChipIngressService) Start(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockChipIngressService) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockChipIngressService) Send(ctx context.Context, payload synchronization.TelemPayload) {
	m.Called(ctx, payload)
}

func (m *MockChipIngressService) HealthReport() map[string]error {
	args := m.Called()
	return args.Get(0).(map[string]error)
}

func (m *MockChipIngressService) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockChipIngressService) Ready() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewChipIngressAgent(t *testing.T) {
	t.Run("Success - Ethereum Mainnet", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)
		chainID := "1"
		network := "EVM"
		contractID := "0x1234"
		telemType := synchronization.OCR2Median

		adapter, err := NewChipIngressAgent(mockTelemService, network, chainID, contractID, telemType, lggr)
		require.NoError(t, err)
		require.NotNil(t, adapter)

		// Verify chain selector was derived correctly (Ethereum mainnet)
		assert.Equal(t, uint64(5009297550715157269), adapter.ChainSelector)
		assert.Equal(t, "EVM", adapter.Network)
		assert.Equal(t, "1", adapter.ChainID)
		assert.Equal(t, contractID, adapter.ContractID)

		// Verify domain and entity were derived correctly
		assert.Equal(t, "data-feeds.telemetry.ocr2-median", adapter.Domain)
		assert.Equal(t, "offchainreporting2.TelemetryWrapper", adapter.Entity)
	})

	t.Run("Success - Polygon", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgent(mockTelemService, "EVM", "137", "0xabc", synchronization.OCR3CCIPCommit, lggr)
		require.NoError(t, err)

		// Verify chain selector was derived correctly (Polygon)
		assert.Equal(t, uint64(4051577828743386545), adapter.ChainSelector)
		assert.Equal(t, "EVM", adapter.Network)
		assert.Equal(t, "137", adapter.ChainID)

		// Verify domain and entity for CCIP commit
		assert.Equal(t, "ccip.telemetry.ocr3", adapter.Domain)
		assert.Equal(t, "offchainreporting3.TelemetryWrapper", adapter.Entity)
	})

	t.Run("Success - Arbitrum One", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgent(mockTelemService, "EVM", "42161", "0xdef", synchronization.OCR2CCIPExec, lggr)
		require.NoError(t, err)

		// Verify chain selector was derived correctly (Arbitrum One)
		assert.Equal(t, uint64(4949039107694359620), adapter.ChainSelector)
	})

	t.Run("Error - nil telemetry service", func(t *testing.T) {
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgent(nil, "EVM", "1", "0x1234", synchronization.OCR2Median, lggr)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "telemetry service cannot be nil")
	})

	t.Run("Error - invalid chainID", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgent(mockTelemService, "EVM", "1234567890123456", "0x1234", synchronization.OCR2Median, lggr)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "failed to get chain details")
	})

	t.Run("Error - invalid network family", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgent(mockTelemService, "INVALID_NETWORK", "1", "0x1234", synchronization.OCR2Median, lggr)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "failed to get chain details")
	})

	t.Run("Error - invalid telemetry type", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgent(mockTelemService, "EVM", "1", "0x1234", synchronization.TelemetryType("unknown"), lggr)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "failed to map telemetry type to domain/entity")
	})
}

func TestNewChipIngressAgentMultitype(t *testing.T) {
	t.Run("Success - creates multitype agent", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgentMultitype(mockTelemService, "EVM", "1", "0x1234", lggr)
		require.NoError(t, err)
		require.NotNil(t, adapter)

		// Verify chain selector was derived correctly
		assert.Equal(t, uint64(5009297550715157269), adapter.ChainSelector)
		assert.Equal(t, "EVM", adapter.Network)
		assert.Equal(t, "1", adapter.ChainID)

		// Verify TelemType, Domain, Entity are empty for multitype
		assert.Empty(t, adapter.TelemType)
		assert.Empty(t, adapter.Domain)
		assert.Empty(t, adapter.Entity)
	})

	t.Run("Error - nil telemetry service", func(t *testing.T) {
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgentMultitype(nil, "EVM", "1", "0x1234", lggr)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "telemetry service cannot be nil")
	})

	t.Run("Error - invalid chainID", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgentMultitype(mockTelemService, "EVM", "invalid", "0x1234", lggr)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "failed to get chain details")
	})
}

func TestChipIngressAgent_SendLog(t *testing.T) {
	t.Run("Success - sends to telemetry service", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)
		contractID := "0x1234567890abcdef"
		telemType := synchronization.OCR2Median
		telemetryLog := []byte("test telemetry data")

		adapter, err := NewChipIngressAgent(mockTelemService, "EVM", "1", contractID, telemType, lggr)
		require.NoError(t, err)

		expectedPayload := synchronization.TelemPayload{
			Telemetry:     telemetryLog,
			TelemType:     telemType,
			ContractID:    contractID,
			ChainSelector: adapter.ChainSelector,
			Domain:        adapter.Domain,
			Entity:        adapter.Entity,
			Network:       adapter.Network,
		}
		mockTelemService.On("Send", mock.Anything, expectedPayload).Once()

		// Call SendLog
		adapter.SendLog(telemetryLog)

		// Verify Send was called
		mockTelemService.AssertExpectations(t)
	})

	t.Run("SendLog continues even if send fails", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)
		contractID := "0xabcdef"
		telemType := synchronization.OCR3Automation
		telemetryLog := []byte("test data")

		adapter, err := NewChipIngressAgent(mockTelemService, "EVM", "137", contractID, telemType, lggr)
		require.NoError(t, err)

		expectedPayload := synchronization.TelemPayload{
			Telemetry:     telemetryLog,
			TelemType:     telemType,
			ContractID:    contractID,
			ChainSelector: adapter.ChainSelector,
			Domain:        adapter.Domain,
			Entity:        adapter.Entity,
			Network:       adapter.Network,
		}
		mockTelemService.On("Send", mock.Anything, expectedPayload).Once()

		// Call SendLog - should not panic
		assert.NotPanics(t, func() {
			adapter.SendLog(telemetryLog)
		})

		mockTelemService.AssertExpectations(t)
	})

	t.Run("Multiple SendLog calls", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)
		contractID := "0x999"
		telemType := synchronization.OCR2CCIPCommit

		adapter, err := NewChipIngressAgent(mockTelemService, "EVM", "42161", contractID, telemType, lggr)
		require.NoError(t, err)

		// Setup expectations for multiple calls
		log1 := []byte("log 1")
		log2 := []byte("log 2")
		log3 := []byte("log 3")

		payload1 := synchronization.TelemPayload{
			Telemetry:     log1,
			TelemType:     telemType,
			ContractID:    contractID,
			ChainSelector: adapter.ChainSelector,
			Domain:        adapter.Domain,
			Entity:        adapter.Entity,
			Network:       adapter.Network,
		}
		payload2 := synchronization.TelemPayload{
			Telemetry:     log2,
			TelemType:     telemType,
			ContractID:    contractID,
			ChainSelector: adapter.ChainSelector,
			Domain:        adapter.Domain,
			Entity:        adapter.Entity,
			Network:       adapter.Network,
		}
		payload3 := synchronization.TelemPayload{
			Telemetry:     log3,
			TelemType:     telemType,
			ContractID:    contractID,
			ChainSelector: adapter.ChainSelector,
			Domain:        adapter.Domain,
			Entity:        adapter.Entity,
			Network:       adapter.Network,
		}

		mockTelemService.On("Send", mock.Anything, payload1).Once()
		mockTelemService.On("Send", mock.Anything, payload2).Once()
		mockTelemService.On("Send", mock.Anything, payload3).Once()

		// Send multiple logs
		adapter.SendLog(log1)
		adapter.SendLog(log2)
		adapter.SendLog(log3)

		mockTelemService.AssertExpectations(t)
	})

	t.Run("SendLog on multitype agent logs warning", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgentMultitype(mockTelemService, "EVM", "1", "0x1234", lggr)
		require.NoError(t, err)

		// SendLog should not call Send on multitype agent
		assert.NotPanics(t, func() {
			adapter.SendLog([]byte("test"))
		})

		// Verify Send was NOT called
		mockTelemService.AssertNotCalled(t, "Send", mock.Anything, mock.Anything)
	})
}

func TestChipIngressAgent_SendTypedLog(t *testing.T) {
	t.Run("Success - sends typed log", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)
		contractID := "0x1234"

		adapter, err := NewChipIngressAgentMultitype(mockTelemService, "EVM", "1", contractID, lggr)
		require.NoError(t, err)

		telemType := synchronization.OCR2Median
		telemetryLog := []byte("test data")

		expectedPayload := synchronization.TelemPayload{
			Telemetry:     telemetryLog,
			TelemType:     telemType,
			ContractID:    contractID,
			ChainSelector: adapter.ChainSelector,
			Domain:        "data-feeds.telemetry.ocr2-median",
			Entity:        "offchainreporting2.TelemetryWrapper",
			Network:       adapter.Network,
		}
		mockTelemService.On("Send", mock.Anything, expectedPayload).Once()

		adapter.SendTypedLog(telemType, telemetryLog)

		mockTelemService.AssertExpectations(t)
	})

	t.Run("SendTypedLog with different types", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)
		contractID := "0x5678"

		adapter, err := NewChipIngressAgentMultitype(mockTelemService, "EVM", "137", contractID, lggr)
		require.NoError(t, err)

		// Send with OCR2Median
		mockTelemService.On("Send", mock.Anything, mock.MatchedBy(func(p synchronization.TelemPayload) bool {
			return p.TelemType == synchronization.OCR2Median && p.Domain == "data-feeds.telemetry.ocr2-median"
		})).Once()
		adapter.SendTypedLog(synchronization.OCR2Median, []byte("median data"))

		// Send with OCR3Mercury
		mockTelemService.On("Send", mock.Anything, mock.MatchedBy(func(p synchronization.TelemPayload) bool {
			return p.TelemType == synchronization.OCR3Mercury && p.Domain == "data-streams.telemetry.ocr3-mercury"
		})).Once()
		adapter.SendTypedLog(synchronization.OCR3Mercury, []byte("mercury data"))

		mockTelemService.AssertExpectations(t)
	})

	t.Run("SendTypedLog with invalid type logs error", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgentMultitype(mockTelemService, "EVM", "1", "0x1234", lggr)
		require.NoError(t, err)

		// Should not panic, just log error
		assert.NotPanics(t, func() {
			adapter.SendTypedLog(synchronization.TelemetryType("invalid"), []byte("test"))
		})

		// Verify Send was NOT called due to invalid type
		mockTelemService.AssertNotCalled(t, "Send", mock.Anything, mock.Anything)
	})

	t.Run("SendTypedLog works on single-type agent too", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)
		contractID := "0x1234"

		// Create single-type agent
		adapter, err := NewChipIngressAgent(mockTelemService, "EVM", "1", contractID, synchronization.OCR2Median, lggr)
		require.NoError(t, err)

		// SendTypedLog should still work
		mockTelemService.On("Send", mock.Anything, mock.MatchedBy(func(p synchronization.TelemPayload) bool {
			return p.TelemType == synchronization.OCR3Mercury
		})).Once()

		adapter.SendTypedLog(synchronization.OCR3Mercury, []byte("different type"))

		mockTelemService.AssertExpectations(t)
	})
}

func TestChipIngressAgent_ExportedFields(t *testing.T) {
	mockTelemService := new(MockChipIngressService)
	lggr := logger.TestLogger(t)
	contractID := "0x1234567890"
	telemType := synchronization.OCR3Mercury

	adapter, err := NewChipIngressAgent(mockTelemService, "EVM", "1", contractID, telemType, lggr)
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
		assert.Equal(t, "data-streams.telemetry.ocr3-mercury", adapter.Domain)
	})

	t.Run("Entity", func(t *testing.T) {
		assert.Equal(t, "offchainreporting3.TelemetryWrapper", adapter.Entity)
	})
}

func TestChipIngressAgent_InterfaceCompliance(t *testing.T) {
	t.Run("single-type agent implements MonitoringEndpoint", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgent(mockTelemService, "EVM", "1", "0x123", synchronization.OCR2Median, lggr)
		require.NoError(t, err)

		// Verify it can be assigned to the interface
		var _ interface{} = adapter

		// Call the interface method
		mockTelemService.On("Send", mock.Anything, mock.Anything)
		adapter.SendLog([]byte("test"))
		mockTelemService.AssertCalled(t, "Send", mock.Anything, mock.Anything)
	})

	t.Run("multitype agent implements MultitypeMonitoringEndpoint", func(t *testing.T) {
		mockTelemService := new(MockChipIngressService)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgentMultitype(mockTelemService, "EVM", "1", "0x123", lggr)
		require.NoError(t, err)

		// Verify it implements MultitypeMonitoringEndpoint
		var me MultitypeMonitoringEndpoint = adapter
		require.NotNil(t, me)

		// Call the interface method
		mockTelemService.On("Send", mock.Anything, mock.Anything)
		adapter.SendTypedLog(synchronization.OCR2Median, []byte("test"))
		mockTelemService.AssertCalled(t, "Send", mock.Anything, mock.Anything)
	})
}
