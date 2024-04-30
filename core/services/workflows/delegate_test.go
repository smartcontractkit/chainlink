package workflows_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
)

func TestDelegate_JobSpecValidator(t *testing.T) {
	t.Parallel()

	var tt = []struct {
		name  string
		toml  string
		valid bool
	}{
		{
			"valid spec",
			`
type = "workflow"
schemaVersion = 1
workflowId = "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef"
workflowOwner = "00000000000000000000000000000000000000aa"
`,
			true,
		},
		{
			"parse error",
			`
invalid syntax{{{{
`,
			false,
		},
		{
			"invalid job type",
			`
type = "work flows"
schemaVersion = 1
`,
			false,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := workflows.ValidatedWorkflowSpec(tc.toml)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
