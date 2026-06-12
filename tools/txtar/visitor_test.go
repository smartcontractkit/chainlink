package txtar

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTreeFile(t *testing.T, path string) {
	t.Helper()

	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte("txtar"), 0o644))
}

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
		files    []string
		recurse  RecurseOpt
		wantDirs []string
	}{
		{
			name: "recurse finds nested txtar directories",
			files: []string{
				"a/test.txtar",
				"b/readme.md",
				"nested/deep/test.txtar",
			},
			recurse:  Recurse,
			wantDirs: []string{"a", "nested/deep"},
		},
		{
			name: "no recurse only visits root directory",
			files: []string{
				"root.txtar",
				"nested/test.txtar",
			},
			recurse:  NoRecurse,
			wantDirs: []string{"."},
		},
		{
			name: "no recurse skips nested txtar directories",
			files: []string{
				"nested/test.txtar",
			},
			recurse:  NoRecurse,
			wantDirs: nil,
		},
		{
			name: "recurse includes root when it contains txtar files",
			files: []string{
				"root.txtar",
				"child/test.txtar",
			},
			recurse:  Recurse,
			wantDirs: []string{".", "child"},
		},
		{
			name: "directories without txtar files are ignored",
			files: []string{
				"empty/readme.md",
				"scripts/test.txtar",
			},
			recurse:  Recurse,
			wantDirs: []string{"scripts"},
		},
		{
			name:     "empty root returns no directories",
			files:    nil,
			recurse:  Recurse,
			wantDirs: nil,
		},
		{
			name: "*txtar suffix matches non dotted extensions",
			files: []string{
				"weird/foo_txtar",
			},
			recurse:  Recurse,
			wantDirs: []string{"weird"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			root := t.TempDir()
			for _, file := range tt.files {
				writeTreeFile(t, filepath.Join(root, file))
			}

			got := collectTxtarDirs(t, root, tt.recurse)
			assert.Equal(t, tt.wantDirs, got)
		})
	}
}

func TestDirVisitor_Walk_callbackError(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeTreeFile(t, filepath.Join(root, "test.txtar"))

	wantErr := errors.New("callback failed")
	visitor := NewDirVisitor(root, Recurse, func(string) error {
		return wantErr
	})

	assert.ErrorIs(t, visitor.Walk(), wantErr)
}

func TestDirVisitor_Walk_missingRoot(t *testing.T) {
	t.Parallel()

	root := filepath.Join(t.TempDir(), "missing")
	visitor := NewDirVisitor(root, Recurse, func(string) error {
		t.Fatal("callback should not run")
		return nil
	})

	assert.Error(t, visitor.Walk())
}
