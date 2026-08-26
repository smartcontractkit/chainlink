package ghaction_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/ghaction"
)

func TestGHAction_SetOutput_File(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "github_output")
	require.NoError(t, os.WriteFile(outputPath, []byte{}, 0o600))

	var stdout bytes.Buffer
	act := ghaction.New(&stdout, outputPath, "")

	err := act.SetOutput("foo", "bar")
	require.NoError(t, err)

	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "foo")
	assert.Contains(t, string(content), "bar")
	assert.Empty(t, stdout.String())
}

func TestGHAction_SetOutput_FallbackStdout(t *testing.T) {
	if os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("skipping in CI: GITHUB_OUTPUT is set in CI environment")
	}
	t.Parallel()
	var stdout bytes.Buffer
	act := ghaction.New(&stdout, "", "")

	err := act.SetOutput("key", "value")
	require.NoError(t, err)
	assert.Equal(t, "key=value\n", stdout.String())
}

func TestGHAction_SetEnv_File(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, "github_env")
	require.NoError(t, os.WriteFile(envPath, []byte{}, 0o600))

	var stdout bytes.Buffer
	act := ghaction.New(&stdout, "", envPath)

	err := act.SetEnv("MY_VAR", "my_val")
	require.NoError(t, err)

	content, err := os.ReadFile(envPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "MY_VAR")
	assert.Contains(t, string(content), "my_val")
	assert.Empty(t, stdout.String())
}

func TestGHAction_SetEnv_FallbackStdout(t *testing.T) {
	if os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("skipping in CI: GITHUB_ENV is set in CI environment")
	}
	t.Parallel()
	var stdout bytes.Buffer
	act := ghaction.New(&stdout, "", "")

	err := act.SetEnv("MY_VAR", "my_val")
	require.NoError(t, err)
	assert.Equal(t, "MY_VAR=my_val\n", stdout.String())
}

func TestGHAction_AddStepSummary_File(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	summaryPath := filepath.Join(tmpDir, "github_step_summary")
	require.NoError(t, os.WriteFile(summaryPath, []byte{}, 0o600))

	var stdout bytes.Buffer
	act := ghaction.NewWithOptions(&stdout, "", "", summaryPath)

	err := act.AddStepSummary("### Summary Table\n- item 1")
	require.NoError(t, err)

	content, err := os.ReadFile(summaryPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "### Summary Table\n- item 1")
	assert.Empty(t, stdout.String())
}

func TestGHAction_AddStepSummary_FallbackStdout(t *testing.T) {
	if os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("skipping in CI: GITHUB_STEP_SUMMARY is set in CI environment")
	}
	t.Parallel()
	var stdout bytes.Buffer
	act := ghaction.NewWithOptions(&stdout, "", "", "")

	err := act.AddStepSummary("### Summary Table")
	require.NoError(t, err)
	assert.Equal(t, "### Summary Table\n", stdout.String())
}

func TestGHAction_WithGroup(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	act := ghaction.New(&stdout, "", "")

	ran := false
	act.WithGroup("grouped-task", func() {
		ran = true
	})

	assert.True(t, ran)
	assert.Contains(t, stdout.String(), "grouped-task")
}
