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
