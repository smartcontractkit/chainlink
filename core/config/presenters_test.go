package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfigPrinter(t *testing.T) {
	cfg := NewGeneralConfig()
	printer, err := NewConfigPrinter(cfg)
	require.NoError(t, err)
	require.NotNil(t, printer)
}
