package smoke

import (
	"fmt"
	"testing"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/stretchr/testify/require"
)

func TestEnvVars(t *testing.T) {
	// cfg, err := tc.GetConfig([]string{"Smoke"}, tc.Automation)
	// require.NoError(t, err)

	cfg := &tc.TestConfig{}
	err := cfg.ReadConfigValuesFromEnvVars()
	require.NoError(t, err)
	fmt.Printf("cfg: %+v", cfg)
}
