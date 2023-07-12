package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
)

func Test_ocrConfig(t *testing.T) {
	evmOcrCfg := cltest.NewTestChainScopedConfig(t) //fallback.toml values
	require.Equal(t, uint16(4), evmOcrCfg.EVM().OCR().ContractConfirmations())
	require.Equal(t, cltest.MustParseDuration(t, "10s"), evmOcrCfg.EVM().OCR().ContractTransmitterTransmitTimeout())
	require.Equal(t, cltest.MustParseDuration(t, "10s"), evmOcrCfg.EVM().OCR().DatabaseTimeout())
	require.Equal(t, cltest.MustParseDuration(t, "1s"), evmOcrCfg.EVM().OCR().ObservationGracePeriod())
}
