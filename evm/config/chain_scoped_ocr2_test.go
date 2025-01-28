package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/evm/config/configtest"
)

func Test_ocr2Config(t *testing.T) {
	cfg := configtest.NewChainScopedConfig(t, nil) // fallback.toml values
	require.Equal(t, uint32(5400000), cfg.EVM().OCR2().Automation().GasLimit())
}
