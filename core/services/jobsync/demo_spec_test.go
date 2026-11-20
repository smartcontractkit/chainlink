package jobsync

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities"
)

// TestDemoCronSpecIsValid keeps the spec the NOP demo proposes in sync with
// what the node accepts. The demo lives in nop-management/cli/platform, which
// cannot import this module, so the spec is duplicated here on purpose: if the
// demo's spec changes, this test must change with it.
func TestDemoCronSpecIsValid(t *testing.T) {
	t.Parallel()
	const spec = `type = "standardcapabilities"
schemaVersion = 1
name = "dyncfg-cron-1"
externalJobID = "c61f4bde-0b4e-5198-b14c-7bc0b747d1f6"
forwardingAllowed = false
command = "cron"
config = """"""
`
	_, err := job.ValidateSpec(spec)
	require.NoError(t, err, "the demo's spec must parse")
	// ValidatedStandardCapabilitiesSpec is what the node runs on POST /v2/jobs.
	jb, err := standardcapabilities.ValidatedStandardCapabilitiesSpec(spec)
	require.NoError(t, err, "the demo's spec must validate")
	require.Equal(t, "dyncfg-cron-1", jb.Name.String)
	require.Equal(t, "cron", jb.StandardCapabilitiesSpec.Command)
}
