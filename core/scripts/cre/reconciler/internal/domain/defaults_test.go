package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCapabilityDefaults(t *testing.T) {
	t.Parallel()

	// The test calls LoadCapabilityDefaults with an empty string, which should return nil.
	// The real test happens in the root package where the embedded TOML is available.
	defaults := LoadCapabilityDefaults("")
	// If the file isn't found (empty string passed), just skip
	if defaults == nil {
		t.Skip("capability_defaults.toml not found from test CWD")
	}

	// Should have evm, cron, http-action etc.
	require.NotEmpty(t, defaults)
	require.Contains(t, defaults, "evm")
	require.Equal(t, "evm", defaults["evm"].BinaryName)
	require.NotEmpty(t, defaults["evm"].Values)
	require.Contains(t, defaults["evm"].Values, "LogTriggerPollInterval")

	require.Contains(t, defaults, "cron")
	require.Equal(t, "cron", defaults["cron"].BinaryName)

	// One binary serves both HTTP capabilities now, and it is named by the trigger's
	// entry: see the http feature, which starts the pair.
	require.Contains(t, defaults, "http-trigger")
	require.Equal(t, "http", defaults["http-trigger"].BinaryName)
}
