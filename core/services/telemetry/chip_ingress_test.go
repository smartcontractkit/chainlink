package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	syncmocks "github.com/smartcontractkit/chainlink/v2/core/services/synchronization/mocks"
)

func TestNewChipIngressAgent(t *testing.T) {
	t.Run("Success - Ethereum Mainnet", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		contractID := "0x1234"

		adapter, err := NewChipIngressAgent(syncmocks.NewChipIngressService(t), "EVM", "1", contractID, synchronization.OCR2Median, lggr)
		require.NoError(t, err)
		require.NotNil(t, adapter)

		// Verify chain selector was derived correctly (Ethereum mainnet)
		assert.Equal(t, uint64(5009297550715157269), adapter.ChainSelector)
		assert.Equal(t, "EVM", adapter.Network)
		assert.Equal(t, "1", adapter.ChainID)
		assert.Equal(t, contractID, adapter.ContractID)
	})

	t.Run("Success - Polygon", func(t *testing.T) {
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgent(syncmocks.NewChipIngressService(t), "EVM", "137", "0xabc", synchronization.OCR3CCIPCommit, lggr)
		require.NoError(t, err)

		// Verify chain selector was derived correctly (Polygon)
		assert.Equal(t, uint64(4051577828743386545), adapter.ChainSelector)
		assert.Equal(t, "EVM", adapter.Network)
		assert.Equal(t, "137", adapter.ChainID)
	})

	t.Run("Success - Arbitrum One", func(t *testing.T) {
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgent(syncmocks.NewChipIngressService(t), "EVM", "42161", "0xdef", synchronization.OCR2CCIPExec, lggr)
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
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgent(syncmocks.NewChipIngressService(t), "EVM", "1234567890123456", "0x1234", synchronization.OCR2Median, lggr)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "failed to get chain details")
	})

	t.Run("Error - invalid network family", func(t *testing.T) {
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgent(syncmocks.NewChipIngressService(t), "INVALID_NETWORK", "1", "0x1234", synchronization.OCR2Median, lggr)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "failed to get chain details")
	})

	t.Run("Error - invalid telemetry type", func(t *testing.T) {
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgent(syncmocks.NewChipIngressService(t), "EVM", "1", "0x1234", synchronization.TelemetryType("unknown"), lggr)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "failed to map telemetry type to domain/entity")
	})
}

func TestNewChipIngressAgentMultitype(t *testing.T) {
	t.Run("Success - creates multitype agent", func(t *testing.T) {
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgentMultitype(syncmocks.NewChipIngressService(t), "EVM", "1", "0x1234", lggr)
		require.NoError(t, err)
		require.NotNil(t, adapter)

		// Verify chain selector was derived correctly
		assert.Equal(t, uint64(5009297550715157269), adapter.ChainSelector)
		assert.Equal(t, "EVM", adapter.Network)
		assert.Equal(t, "1", adapter.ChainID)
		assert.Empty(t, adapter.TelemType)
	})

	t.Run("Error - nil telemetry service", func(t *testing.T) {
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgentMultitype(nil, "EVM", "1", "0x1234", lggr)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "telemetry service cannot be nil")
	})

	t.Run("Error - invalid chainID", func(t *testing.T) {
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgentMultitype(syncmocks.NewChipIngressService(t), "EVM", "invalid", "0x1234", lggr)
		require.Error(t, err)
		assert.Nil(t, adapter)
		assert.Contains(t, err.Error(), "failed to get chain details")
	})
}

func TestChipIngressAgent_SendLog(t *testing.T) {
	t.Run("Success - sends to telemetry service", func(t *testing.T) {
		svc := syncmocks.NewChipIngressService(t)
		lggr := logger.TestLogger(t)
		contractID := "0x1234567890abcdef"
		telemType := synchronization.OCR2Median
		telemetryLog := []byte("test telemetry data")

		adapter, err := NewChipIngressAgent(svc, "EVM", "1", contractID, telemType, lggr)
		require.NoError(t, err)

		svc.On("Send", mock.Anything, synchronization.TelemPayload{
			Telemetry:     telemetryLog,
			TelemType:     telemType,
			ContractID:    contractID,
			ChainSelector: adapter.ChainSelector,
			Network:       adapter.Network,
		}).Once()

		adapter.SendLog(telemetryLog)
	})

	t.Run("Multiple SendLog calls", func(t *testing.T) {
		svc := syncmocks.NewChipIngressService(t)
		lggr := logger.TestLogger(t)
		contractID := "0x999"
		telemType := synchronization.OCR2CCIPCommit

		adapter, err := NewChipIngressAgent(svc, "EVM", "42161", contractID, telemType, lggr)
		require.NoError(t, err)

		logs := [][]byte{[]byte("log 1"), []byte("log 2"), []byte("log 3")}
		for _, l := range logs {
			svc.On("Send", mock.Anything, synchronization.TelemPayload{
				Telemetry:     l,
				TelemType:     telemType,
				ContractID:    contractID,
				ChainSelector: adapter.ChainSelector,
				Network:       adapter.Network,
			}).Once()
		}

		for _, l := range logs {
			adapter.SendLog(l)
		}
	})

	t.Run("SendLog on multitype agent does not send", func(t *testing.T) {
		t.Parallel()
		svc := syncmocks.NewChipIngressService(t)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgentMultitype(svc, "EVM", "1", "0x1234", lggr)
		require.NoError(t, err)

		adapter.SendLog([]byte("test"))
		svc.AssertNotCalled(t, "Send", mock.Anything, mock.Anything)
	})
}

func TestChipIngressAgent_SendTypedLog(t *testing.T) {
	t.Run("Success - sends typed log", func(t *testing.T) {
		svc := syncmocks.NewChipIngressService(t)
		lggr := logger.TestLogger(t)
		contractID := "0x1234"

		adapter, err := NewChipIngressAgentMultitype(svc, "EVM", "1", contractID, lggr)
		require.NoError(t, err)

		telemType := synchronization.OCR2Median
		telemetryLog := []byte("test data")

		svc.On("Send", mock.Anything, synchronization.TelemPayload{
			Telemetry:     telemetryLog,
			TelemType:     telemType,
			ContractID:    contractID,
			ChainSelector: adapter.ChainSelector,
			Network:       adapter.Network,
		}).Once()

		adapter.SendTypedLog(telemType, telemetryLog)
	})

	t.Run("SendTypedLog with different types", func(t *testing.T) {
		svc := syncmocks.NewChipIngressService(t)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgentMultitype(svc, "EVM", "137", "0x5678", lggr)
		require.NoError(t, err)

		for _, telemType := range []synchronization.TelemetryType{synchronization.OCR2Median, synchronization.OCR3Mercury} {
			svc.On("Send", mock.Anything, mock.MatchedBy(func(p synchronization.TelemPayload) bool {
				return p.TelemType == telemType
			})).Once()
			adapter.SendTypedLog(telemType, []byte("data"))
		}
	})

	t.Run("SendTypedLog works on single-type agent too", func(t *testing.T) {
		svc := syncmocks.NewChipIngressService(t)
		lggr := logger.TestLogger(t)

		adapter, err := NewChipIngressAgent(svc, "EVM", "1", "0x1234", synchronization.OCR2Median, lggr)
		require.NoError(t, err)

		svc.On("Send", mock.Anything, mock.MatchedBy(func(p synchronization.TelemPayload) bool {
			return p.TelemType == synchronization.OCR3Mercury
		})).Once()

		adapter.SendTypedLog(synchronization.OCR3Mercury, []byte("different type"))
	})
}
