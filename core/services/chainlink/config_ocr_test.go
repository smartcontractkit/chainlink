package chainlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOCRConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	ocrCfg := cfg.OCR()

	expectedObservationTimeout, err := time.ParseDuration("11s")
	require.NoError(t, err)
	expectedBlockchainTimeout, err := time.ParseDuration("3s")
	require.NoError(t, err)
	expectedContractPollInterval, err := time.ParseDuration("1h0m0s")
	require.NoError(t, err)
	expectedContractSubscribeInterval, err := time.ParseDuration("1m0s")
	require.NoError(t, err)

	require.Equal(t, true, ocrCfg.Enabled())
	require.Equal(t, expectedObservationTimeout, ocrCfg.ObservationTimeout())
	require.Equal(t, expectedBlockchainTimeout, ocrCfg.BlockchainTimeout())
	require.Equal(t, expectedContractPollInterval, ocrCfg.ContractPollInterval())
	require.Equal(t, expectedContractSubscribeInterval, ocrCfg.ContractSubscribeInterval())
	require.Equal(t, true, ocrCfg.SimulateTransactions())
	require.Equal(t, false, ocrCfg.TraceLogging())
	require.Equal(t, uint32(12), ocrCfg.DefaultTransactionQueueDepth())
	require.Equal(t, false, ocrCfg.CaptureEATelemetry())

	keyBundleID, err := ocrCfg.KeyBundleID()
	require.NoError(t, err)
	require.Equal(t, "acdd42797a8b921b2910497badc5000600000000000000000000000000000000", keyBundleID)

	transmitterAddress, err := ocrCfg.TransmitterAddress()
	require.NoError(t, err)
	require.Equal(t, "0xa0788FC17B1dEe36f057c42B6F373A34B014687e", transmitterAddress.String())
}
