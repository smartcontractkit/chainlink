package synchronization_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

func TestNewChipIngressBatchClient(t *testing.T) {
	t.Run("Success - creates client", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		var ks keystore.CSA

		client := synchronization.NewChipIngressBatchClient(nil, nil, ks, lggr, nil)

		require.NotNil(t, client)
		assert.Equal(t, "ChipIngressBatchClient", client.Name())
	})

	t.Run("Client implements ChipIngressService interface", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		var ks keystore.CSA

		client := synchronization.NewChipIngressBatchClient(nil, nil, ks, lggr, nil)

		// Verify it implements the interface
		var _ synchronization.ChipIngressService = client
	})
}

func TestChipIngressBatchClient_Send(t *testing.T) {
	t.Run("Success - sends telemetry", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		var ks keystore.CSA

		client := synchronization.NewChipIngressBatchClient(nil, nil, ks, lggr, nil)
		require.NotNil(t, client)

		servicetest.Run(t, client)

		// Send telemetry
		ctx := testutils.Context(t)
		telemetry := []byte("test telemetry data")
		contractID := "0x1234567890abcdef"
		telemType := synchronization.OCR2Median
		chainSelector := uint64(5009297550715157269) // Ethereum mainnet
		domain := "data-feeds"
		entity := "ocr.v2.median.telemetry"

		client.Send(ctx, telemetry, contractID, telemType, chainSelector, domain, entity)
	})

	t.Run("Multiple sends for same contract", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		var ks keystore.CSA

		client := synchronization.NewChipIngressBatchClient(nil, nil, ks, lggr, nil)
		require.NotNil(t, client)

		servicetest.Run(t, client)

		ctx := testutils.Context(t)

		// Send multiple telemetry messages for the same contract
		for i := 0; i < 5; i++ {
			client.Send(ctx, []byte("test data"), "0xabc", synchronization.OCR3CCIPCommit, 4051577828743386545, "ccip", "ocr.v3.ccip.commit.telemetry")
		}
	})

	t.Run("Multiple contracts with different telemetry types", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		var ks keystore.CSA

		client := synchronization.NewChipIngressBatchClient(nil, nil, ks, lggr, nil)
		require.NotNil(t, client)

		servicetest.Run(t, client)

		ctx := testutils.Context(t)

		testCases := []struct {
			name          string
			contractID    string
			telemType     synchronization.TelemetryType
			chainSelector uint64
			domain        string
			entity        string
		}{
			{
				name:          "OCR2Median - Ethereum",
				contractID:    "0x1",
				telemType:     synchronization.OCR2Median,
				chainSelector: 5009297550715157269,
				domain:        "data-feeds",
				entity:        "ocr.v2.median.telemetry",
			},
			{
				name:          "OCR3CCIPCommit - Polygon",
				contractID:    "0x2",
				telemType:     synchronization.OCR3CCIPCommit,
				chainSelector: 4051577828743386545,
				domain:        "ccip",
				entity:        "ocr.v3.ccip.commit.telemetry",
			},
			{
				name:          "OCR3Mercury - Ethereum",
				contractID:    "0x3",
				telemType:     synchronization.OCR3Mercury,
				chainSelector: 5009297550715157269,
				domain:        "data-streams",
				entity:        "ocr.v3.mercury.telemetry",
			},
			{
				name:          "OCR3Automation - Arbitrum",
				contractID:    "0x4",
				telemType:     synchronization.OCR3Automation,
				chainSelector: 4949039107694359620,
				domain:        "automation",
				entity:        "ocr.v3.automation.telemetry",
			},
			{
				name:          "OCR2CCIPExec - Ethereum",
				contractID:    "0x5",
				telemType:     synchronization.OCR2CCIPExec,
				chainSelector: 5009297550715157269,
				domain:        "ccip",
				entity:        "ocr.v2.ccip.exec.telemetry",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				client.Send(ctx, []byte("test telemetry"), tc.contractID, tc.telemType, tc.chainSelector, tc.domain, tc.entity)
			})
		}
	})
}

func TestChipIngressBatchClient_Lifecycle(t *testing.T) {
	t.Run("Start and Close", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		var ks keystore.CSA

		client := synchronization.NewChipIngressBatchClient(nil, nil, ks, lggr, nil)
		require.NotNil(t, client)

		ctx := testutils.Context(t)

		// Start
		err := client.Start(ctx)
		require.NoError(t, err)

		// Verify Ready
		err = client.Ready()
		assert.NoError(t, err)

		// Check HealthReport
		health := client.HealthReport()
		assert.NotNil(t, health)

		// Close
		err = client.Close()
		assert.NoError(t, err)
	})

	t.Run("Multiple Start calls", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		var ks keystore.CSA

		client := synchronization.NewChipIngressBatchClient(nil, nil, ks, lggr, nil)
		require.NotNil(t, client)

		ctx := testutils.Context(t)

		// First start should succeed
		err := client.Start(ctx)
		require.NoError(t, err)

		// Second start should fail (service already started)
		err = client.Start(ctx)
		assert.Error(t, err)

		err = client.Close()
		assert.NoError(t, err)
	})

	t.Run("Close without Start returns error", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		var ks keystore.CSA

		client := synchronization.NewChipIngressBatchClient(nil, nil, ks, lggr, nil)
		require.NotNil(t, client)

		// Close without starting should return an error
		err := client.Close()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "has not been started")
	})
}

func TestChipIngressBatchClient_ServiceMethods(t *testing.T) {
	t.Run("Name returns correct value", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		var ks keystore.CSA

		client := synchronization.NewChipIngressBatchClient(nil, nil, ks, lggr, nil)

		assert.Equal(t, "ChipIngressBatchClient", client.Name())
	})

	t.Run("HealthReport returns map", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		var ks keystore.CSA

		client := synchronization.NewChipIngressBatchClient(nil, nil, ks, lggr, nil)

		servicetest.Run(t, client)

		health := client.HealthReport()
		assert.NotNil(t, health)
		assert.IsType(t, map[string]error{}, health)
	})

	t.Run("Ready returns no error when started", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		var ks keystore.CSA

		client := synchronization.NewChipIngressBatchClient(nil, nil, ks, lggr, nil)

		servicetest.Run(t, client)

		err := client.Ready()
		assert.NoError(t, err)
	})
}
