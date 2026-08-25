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
	assert.Equal(t, "foo=bar\n", string(content))
	assert.Empty(t, stdout.String())
}

func TestGHAction_SetOutput_Multiline(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "github_output")
	require.NoError(t, os.WriteFile(outputPath, []byte{}, 0o600))

	var stdout bytes.Buffer
	act := ghaction.New(&stdout, outputPath, "")

	err := act.SetOutput("matrix", "{\n  \"include\": [1, 2]\n}")
	require.NoError(t, err)

	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "matrix<<ghadelimiter_")
	assert.Contains(t, string(content), "{\n  \"include\": [1, 2]\n}")
}

func TestGHAction_SetOutput_FallbackStdout(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	act := ghaction.New(&stdout, "", "")

	err := act.SetOutput("key", "value")
	require.NoError(t, err)
	assert.Equal(t, "key=value\n", stdout.String())
}

func TestGHAction_SetEnv(t *testing.T) {
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
	assert.Equal(t, "MY_VAR=my_val\n", string(content))
}

func TestGHAction_Annotations(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	act := ghaction.New(&stdout, "", "")

	act.Errorf("failure %s", "test")
	act.Warningf("warning %d", 42)
	act.Group("my-group")
	act.EndGroup()

	expected := "::error::failure test\n::warning::warning 42\n::group::my-group\n::endgroup::\n"
	assert.Equal(t, expected, stdout.String())
}
