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

func TestGHAction_GetInputOrEnv(t *testing.T) {
	t.Setenv("INPUT_MY-INPUT", "")
	t.Setenv("MY_ENV_VAR", "")
	var stdout bytes.Buffer
	act := ghaction.NewAction(&stdout)

	assert.Empty(t, act.GetInputOrEnv("my-input", "MY_ENV_VAR"))

	t.Setenv("MY_ENV_VAR", "from-env")
	assert.Equal(t, "from-env", act.GetInputOrEnv("my-input", "MY_ENV_VAR"))

	t.Setenv("INPUT_MY-INPUT", "from-input")
	assert.Equal(t, "from-input", act.GetInputOrEnv("my-input", "MY_ENV_VAR"))
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
