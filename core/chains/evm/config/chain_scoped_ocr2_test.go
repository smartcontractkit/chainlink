package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
)

func Test_ocr2Config(t *testing.T) {
	evmOcrCfg := cltest.NewTestChainScopedConfig(t) //fallback.toml values
	require.Equal(t, uint32(5400000), evmOcrCfg.EVM().OCR2().Automation().GasLimit())
}
