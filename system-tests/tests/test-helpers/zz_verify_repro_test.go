package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"

	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
)

// Reproduces the CI failure: the merged settings doc (boot baseline carrying
// ShardingFailoverEnabled + the auto-failover override) must be rejected by
// the strict catalog check.
func TestVerifyCRESettingsReproCIUnknownFields(t *testing.T) {
	mergedTOML := `[global]
ShardingFailoverAutoExecutionEnabled = "true"
ShardingFailoverAutoWindow = "30s"
ShardingFailoverEnabled = "true"

[global.PerOrg]
  BaseTriggerRetransmitEnabled = "true"
`
	err := cre_jobs.VerifyCRESettings(mergedTOML)
	require.Error(t, err, "expected the CI unknown-fields rejection")
	t.Logf("reproduced CI error: %v", err)
}
