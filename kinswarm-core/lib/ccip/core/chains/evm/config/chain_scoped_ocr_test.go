package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
)

func Test_ocrConfig(t *testing.T) {
	cfg := testutils.NewTestChainScopedConfig(t, nil) //fallback.toml values

	require.Equal(t, uint16(4), cfg.EVM().OCR().ContractConfirmations())
	require.Equal(t, mustParseDuration(t, "10s"), cfg.EVM().OCR().ContractTransmitterTransmitTimeout())
	require.Equal(t, mustParseDuration(t, "10s"), cfg.EVM().OCR().DatabaseTimeout())
	require.Equal(t, mustParseDuration(t, "1s"), cfg.EVM().OCR().ObservationGracePeriod())
}

func mustParseDuration(t testing.TB, durationStr string) time.Duration {
	t.Helper()

	duration, err := time.ParseDuration(durationStr)
	require.NoError(t, err)
	return duration
}
