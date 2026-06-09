package dbdetect

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// makeModuleRepo creates a temp repo with:
//
//	root/go.mod          (module github.com/smartcontractkit/chainlink/v2)
//	root/core/           (no go.mod — stays in main module)
//	root/deployment/go.mod (module github.com/smartcontractkit/chainlink/deployment)
//	root/deployment/ccip/ (sub-package of deployment)
func makeModuleRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(root, "go.mod"),
		[]byte("module github.com/smartcontractkit/chainlink/v2\n"),
		0600,
	))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "core", "services"), 0700))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "deployment", "ccip"), 0700))
	require.NoError(t, os.WriteFile(
		filepath.Join(root, "deployment", "go.mod"),
		[]byte("module github.com/smartcontractkit/chainlink/deployment\n"),
		0600,
	))
	return root
}

func TestNearestGoModDir(t *testing.T) {
	t.Parallel()
	root := makeModuleRepo(t)

	tests := []struct {
		name    string
		dir     string
		wantDir string
	}{
		{
			name:    "repo root returns root",
			dir:     root,
			wantDir: root,
		},
		{
			name:    "core dir has no intermediate go.mod, returns root",
			dir:     filepath.Join(root, "core", "services"),
			wantDir: root,
		},
		{
			name:    "deployment dir returns deployment",
			dir:     filepath.Join(root, "deployment"),
			wantDir: filepath.Join(root, "deployment"),
		},
		{
			name:    "deployment subdir returns deployment",
			dir:     filepath.Join(root, "deployment", "ccip"),
			wantDir: filepath.Join(root, "deployment"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := nearestGoModDir(tc.dir, root)
			assert.Equal(t, tc.wantDir, got)
		})
	}
}

func TestPatternToBaseDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		pattern string
		want    string
	}{
		{"./deployment/...", "./deployment"},
		{"./deployment/ccip/...", "./deployment/ccip"},
		{"./deployment/ccip", "./deployment/ccip"},
		{"./core/...", "./core"},
		{"./...", "."},
		{".", "."},
	}

	for _, tc := range tests {
		t.Run(tc.pattern, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, patternToBaseDir(tc.pattern))
		})
	}
}

func TestRewritePatterns(t *testing.T) {
	t.Parallel()
	root := makeModuleRepo(t)
	depDir := filepath.Join(root, "deployment")

	tests := []struct {
		name     string
		patterns []string
		want     []string
	}{
		{
			name:     "top-level deployment wildcard",
			patterns: []string{"./deployment/..."},
			want:     []string{"./..."},
		},
		{
			name:     "deployment subpackage wildcard",
			patterns: []string{"./deployment/ccip/..."},
			want:     []string{"./ccip/..."},
		},
		{
			name:     "deployment package without wildcard",
			patterns: []string{"./deployment/ccip"},
			want:     []string{"./ccip"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := rewritePatterns(root, depDir, tc.patterns)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestResolveModulePatterns(t *testing.T) {
	t.Parallel()
	root := makeModuleRepo(t)
	depDir := filepath.Join(root, "deployment")

	tests := []struct {
		name     string
		patterns []string
		wantDir  string
		wantPats []string
		wantErr  bool
	}{
		{
			name:     "core patterns stay at root",
			patterns: []string{"./core/services/..."},
			wantDir:  root,
			wantPats: []string{"./core/services/..."},
		},
		{
			name:     "deployment top-level rewrites",
			patterns: []string{"./deployment/..."},
			wantDir:  depDir,
			wantPats: []string{"./..."},
		},
		{
			name:     "deployment subpackage rewrites",
			patterns: []string{"./deployment/ccip/..."},
			wantDir:  depDir,
			wantPats: []string{"./ccip/..."},
		},
		{
			name:     "no relative patterns uses root",
			patterns: []string{"github.com/foo/bar"},
			wantDir:  root,
			wantPats: []string{"github.com/foo/bar"},
		},
		{
			name:     "cross-module patterns error",
			patterns: []string{"./core/...", "./deployment/..."},
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir, pats, err := resolveModulePatterns(root, tc.patterns)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantDir, dir)
			assert.Equal(t, tc.wantPats, pats)
		})
	}
}
