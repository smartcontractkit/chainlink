package job

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	var tt = []struct {
		name      string
		spec      string
		assertion func(t *testing.T, err error)
	}{
		{
			name: "invalid job type",
			spec: `
type="blah"
schemaVersion=1
`,
			assertion: func(t *testing.T, err error) {
				require.True(t, errors.Is(errors.Cause(err), ErrInvalidJobType))
			},
		},
		{
			name: "invalid schema version",
			spec: `
type="vrf"
schemaVersion=2
`,
			assertion: func(t *testing.T, err error) {
				require.True(t, errors.Is(errors.Cause(err), ErrInvalidSchemaVersion))
			},
		},
		{
			name: "missing schema version",
			spec: `
type="vrf"
`,
			assertion: func(t *testing.T, err error) {
				require.True(t, errors.Is(errors.Cause(err), ErrInvalidSchemaVersion))
			},
		},
		{
			name: "missing pipeline spec key",
			spec: `
type="vrf"
schemaVersion=1
`,
			assertion: func(t *testing.T, err error) {
				require.True(t, errors.Is(errors.Cause(err), ErrNoPipelineSpec))
			},
		},
		{
			name: "missing pipeline spec value",
			spec: `
type="vrf"
schemaVersion=1
`,
			assertion: func(t *testing.T, err error) {
				require.True(t, errors.Is(errors.Cause(err), ErrNoPipelineSpec))
			},
		},
		{
			name: "invalid dot",
			spec: `
type="vrf"
schemaVersion=1
observationSource="""
sldkfjalskdjf
"""
`,
			assertion: func(t *testing.T, err error) {
				t.Log(err)
				require.Error(t, err)
			},
		},
		{
			name: "async check",
			spec: `
type="offchainreporting"
schemaVersion=1
observationSource="""
ds [type=bridge async=true]
"""
`,
			assertion: func(t *testing.T, err error) {
				t.Log(err)
				require.Error(t, err)
			},
		},
		{
			name: "happy path",
			spec: `
type="vrf"
schemaVersion=1
observationSource="""
ds [type=http]
"""
`,
			assertion: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := ValidateSpec(tc.spec)
			tc.assertion(t, err)
		})
	}
}
