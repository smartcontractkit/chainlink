package cmd_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/cmd"
)

func TestSampleOutputsCmd_ExecutesSuccessfully(t *testing.T) {
	t.Parallel()

	root := cmd.NewRootCmd()
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"sample-outputs", "--no-color"})

	err := root.Execute()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "PR Diff Size")
	assert.Contains(t, out, "LARGE PR")
	assert.Contains(t, out, "LINT")
	assert.Contains(t, out, "TEST")
	assert.Contains(t, out, "FIXED")
}
