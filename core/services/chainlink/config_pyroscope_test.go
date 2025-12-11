package chainlink

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestPyroscopeConfigTest(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	opts := GeneralConfigOpts{
		ConfigStrings:  []string{fullTOML},
		SecretsStrings: []string{secretsFullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	pcfg := cfg.Pyroscope()

	require.Equal(t, "pyroscope-token", pcfg.AuthToken())
	require.Equal(t, "http://localhost:4040", pcfg.ServerAddress())
	require.Equal(t, "tests", pcfg.Environment())
}
