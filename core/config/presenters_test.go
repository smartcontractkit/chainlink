package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfigPrinter(t *testing.T) {
	cfg := NewGeneralConfig()
	printer := NewConfigPrinter(cfg)
	require.NotNil(t, printer)
}
