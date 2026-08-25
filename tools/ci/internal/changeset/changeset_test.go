package changeset_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/changeset"
)

func TestCheckTags_Valid(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test-changeset.md")
	content := `---
"chainlink": patch
---

Added some feature #added #db_update
`
	require.NoError(t, os.WriteFile(file, []byte(content), 0o600))

	res, err := changeset.CheckTags(file)
	require.NoError(t, err)
	assert.True(t, res.HasTags)
	assert.Equal(t, []string{"#added", "#db_update"}, res.FoundTags)
}

func TestCheckTags_NoTags(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test-changeset.md")
	content := `---
"chainlink": minor
---

Some text without tags
`
	require.NoError(t, os.WriteFile(file, []byte(content), 0o600))

	res, err := changeset.CheckTags(file)
	require.NoError(t, err)
	assert.False(t, res.HasTags)
	assert.Empty(t, res.FoundTags)
}

func TestCheckTags_InvalidSemver(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test-changeset.md")
	content := `---
"chainlink": invalid
---

Some text #added
`
	require.NoError(t, os.WriteFile(file, []byte(content), 0o600))

	_, err := changeset.CheckTags(file)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid changeset semvar value for 'chainlink'")
}

func TestCheckTags_FileNotFound(t *testing.T) {
	t.Parallel()
	_, err := changeset.CheckTags("/path/does/not/exist.md")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
