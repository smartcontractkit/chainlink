package presenters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/stretchr/testify/require"
)

func TestNewConfigPrinter(t *testing.T) {
	config := config.NewGeneralConfig()
	printer, err := presenters.NewConfigPrinter(config)
	require.NoError(t, err)
	require.Contains(t, printer.String(), "CHAINLINK_DEV")
}
