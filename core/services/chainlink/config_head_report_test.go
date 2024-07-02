package chainlink

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeadReportConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	hr := cfg.HeadReport()
	require.True(t, hr.TelemetryEnabled())
}
