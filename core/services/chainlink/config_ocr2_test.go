package chainlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOCR2Config(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	ocr2Cfg := cfg.OCR2()

	expectedContractTransmitterTransmitTimeout, err := time.ParseDuration("1m0s")
	require.NoError(t, err)
	expectedBlockchainTimeout, err := time.ParseDuration("3s")
	require.NoError(t, err)
	expectedDatabaseTimeout, err := time.ParseDuration("8s")
	require.NoError(t, err)
	expectedContractPollInterval, err := time.ParseDuration("1h0m0s")
	require.NoError(t, err)
	expectedContractSubscribeInterval, err := time.ParseDuration("1m0s")
	require.NoError(t, err)

	require.True(t, ocr2Cfg.Enabled())
	require.Equal(t, uint16(11), ocr2Cfg.ContractConfirmations())
	require.Equal(t, expectedContractTransmitterTransmitTimeout, ocr2Cfg.ContractTransmitterTransmitTimeout())
	require.Equal(t, expectedBlockchainTimeout, ocr2Cfg.BlockchainTimeout())
	require.Equal(t, expectedDatabaseTimeout, ocr2Cfg.DatabaseTimeout())
	require.Equal(t, expectedContractPollInterval, ocr2Cfg.ContractPollInterval())
	require.Equal(t, expectedContractSubscribeInterval, ocr2Cfg.ContractSubscribeInterval())
	require.False(t, ocr2Cfg.SimulateTransactions())
	require.False(t, ocr2Cfg.TraceLogging())
	require.Equal(t, uint32(1), ocr2Cfg.DefaultTransactionQueueDepth())
	require.False(t, ocr2Cfg.CaptureEATelemetry())
	require.True(t, ocr2Cfg.CaptureAutomationCustomTelemetry())
	require.True(t, ocr2Cfg.AllowNoBootstrappers())

	keyBundleID, err := ocr2Cfg.KeyBundleID()
	require.NoError(t, err)
	require.Equal(t, "7a5f66bbe6594259325bf2b4f5b1a9c900000000000000000000000000000000", keyBundleID)
}
