package ghaction_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/ghaction"
)

func TestGHAction_SetOutput_FallbackStdout(t *testing.T) {
	t.Setenv("GITHUB_OUTPUT", "")
	var stdout bytes.Buffer
	act := ghaction.NewAction(&stdout)

	err := act.SetOutput("key", "value")
	require.NoError(t, err)
	assert.Equal(t, "key=value\n", stdout.String())
}

func TestGHAction_SetOutputs_FallbackStdout(t *testing.T) {
	t.Setenv("GITHUB_OUTPUT", "")
	var stdout bytes.Buffer
	act := ghaction.NewAction(&stdout)

	outputs := map[string]string{
		"spot":     "co",
		"strategy": "capacity-optimized",
	}
	err := act.SetOutputs(outputs)
	require.NoError(t, err)
	assert.Equal(t, "spot=co\nstrategy=capacity-optimized\n", stdout.String())
}

func TestGHAction_SetEnv_FallbackStdout(t *testing.T) {
	t.Setenv("GITHUB_ENV", "")
	var stdout bytes.Buffer
	act := ghaction.NewAction(&stdout)

	err := act.SetEnv("MY_VAR", "my_val")
	require.NoError(t, err)
	assert.Equal(t, "MY_VAR=my_val\n", stdout.String())
}

func TestGHAction_AddStepSummary_FallbackStdout(t *testing.T) {
	t.Setenv("GITHUB_STEP_SUMMARY", "")
	var stdout bytes.Buffer
	act := ghaction.NewWithOptions(&stdout, "", "", "")

	err := act.AddStepSummary("### Summary Table")
	require.NoError(t, err)
	assert.Equal(t, "### Summary Table\n", stdout.String())
}

func TestGHAction_AddStepSummaryTemplate(t *testing.T) {
	t.Setenv("GITHUB_STEP_SUMMARY", "")
	var stdout bytes.Buffer
	act := ghaction.NewWithOptions(&stdout, "", "", "")

	data := struct {
		Name string
	}{Name: "Runner"}
	err := act.AddStepSummaryTemplate("### Hello {{.Name}}", data)
	require.NoError(t, err)
	assert.Equal(t, "### Hello Runner\n", stdout.String())
}

func TestGHAction_IsGitHubActions(t *testing.T) {
	var stdout bytes.Buffer
	act1 := ghaction.New(&stdout, "/path/to/output", "")
	assert.True(t, act1.IsGitHubActions())

	t.Setenv("GITHUB_ACTIONS", "")
	t.Setenv("GITHUB_OUTPUT", "")
	act2 := ghaction.NewWithOptions(&stdout, "", "", "")
	assert.False(t, act2.IsGitHubActions())
}

func TestGHAction_WithGroup(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	act := ghaction.NewAction(&stdout)

	ran := false
	act.WithGroup("grouped-task", func() {
		ran = true
	})

	assert.True(t, ran)
	assert.Contains(t, stdout.String(), "grouped-task")
}
