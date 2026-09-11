package cmd_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/cmd"
)

func TestPRSizeCmdHelp(t *testing.T) {
	t.Parallel()

	root := cmd.NewRootCmd()
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"pr-size", "--help"})

	err := root.Execute()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, buf.String(), "Calculates the git diff of the current branch against the default branch")
	assert.Contains(t, out, "--strategy")
	assert.Contains(t, out, "--small-limit")
	assert.Contains(t, out, "--medium-limit")
	assert.Contains(t, out, "--fail-on-large")
	assert.Contains(t, out, "--base")
	assert.Contains(t, out, "--ignore-lockfiles")
	assert.Contains(t, out, "--ignore-generated")
}

func TestPRSizeCmdAliases(t *testing.T) {
	t.Parallel()

	aliases := []string{"big-pr-guard", "pr-guard", "diff-guard", "diff-size"}
	for _, alias := range aliases {
		t.Run(alias, func(t *testing.T) {
			t.Parallel()
			root := cmd.NewRootCmd()
			buf := new(bytes.Buffer)
			root.SetOut(buf)
			root.SetErr(buf)
			root.SetArgs([]string{alias, "--help"})

			err := root.Execute()
			require.NoError(t, err)
			assert.Contains(t, buf.String(), "Calculates the git diff of the current branch against the default branch")
		})
	}
}
