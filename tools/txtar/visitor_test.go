package txtar

import (
	"errors"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func collectTxtarDirs(t *testing.T, root string, recurse RecurseOpt) []string {
	t.Helper()

	var dirs []string
	visitor := NewDirVisitor(root, recurse, func(path string) error {
		rel, err := filepath.Rel(root, path)
		require.NoError(t, err)
		dirs = append(dirs, rel)
		return nil
	})
	require.NoError(t, visitor.Walk())

	slices.Sort(dirs)
	return dirs
}

func TestDirVisitor_Walk(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testDir  string
		recurse  RecurseOpt
		wantDirs []string
	}{
		{
			name:     "recurse finds nested txtar directories",
			testDir:  "recurse_nested",
			recurse:  Recurse,
			wantDirs: []string{"a", "nested/deep"},
		},
		{
			name:     "no recurse only visits root directory",
			testDir:  "no_recurse_root",
			recurse:  NoRecurse,
			wantDirs: []string{"."},
		},
		{
			name:     "no recurse skips nested txtar directories",
			testDir:  "no_recurse_nested",
			recurse:  NoRecurse,
			wantDirs: nil,
		},
		{
			name:     "recurse includes root when it contains txtar files",
			testDir:  "recurse_includes_root",
			recurse:  Recurse,
			wantDirs: []string{".", "child"},
		},
		{
			name:     "directories without txtar files are ignored",
			testDir:  "ignore_empty",
			recurse:  Recurse,
			wantDirs: []string{"scripts"},
		},
		{
			name:     "empty root returns no directories",
			testDir:  "empty_root",
			recurse:  Recurse,
			wantDirs: nil,
		},
		{
			name:     "*txtar suffix matches non dotted extensions",
			testDir:  "suffix_matches",
			recurse:  Recurse,
			wantDirs: []string{"weird"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			root := filepath.Join("testdata", tt.testDir)
			got := collectTxtarDirs(t, root, tt.recurse)
			assert.Equal(t, tt.wantDirs, got)
		})
	}
}

func TestDirVisitor_Walk_callbackError(t *testing.T) {
	t.Parallel()

	root := filepath.Join("testdata", "callback_error")

	wantErr := errors.New("callback failed")
	visitor := NewDirVisitor(root, Recurse, func(string) error {
		return wantErr
	})

	assert.ErrorIs(t, visitor.Walk(), wantErr)
}

func TestDirVisitor_Walk_missingRoot(t *testing.T) {
	t.Parallel()

	root := filepath.Join("testdata", "missing_root_dir_that_does_not_exist")
	visitor := NewDirVisitor(root, Recurse, func(string) error {
		t.Fatal("callback should not run")
		return nil
	})

	assert.Error(t, visitor.Walk())
}
