package runner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func makeTestRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	// Root module
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module github.com/smartcontractkit/chainlink/v2\n"), 0600))
	// core/ — plain subdirectory, no go.mod
	require.NoError(t, os.MkdirAll(filepath.Join(root, "core", "store"), 0700))
	// deployment/ — separate module
	require.NoError(t, os.MkdirAll(filepath.Join(root, "deployment", "ccip"), 0700))
	require.NoError(t, os.WriteFile(filepath.Join(root, "deployment", "go.mod"), []byte("module github.com/smartcontractkit/chainlink/deployment\n"), 0600))
	return root
}

func TestResolveModuleDir(t *testing.T) {
	t.Parallel()
	root := makeTestRepo(t)

	tests := []struct {
		name       string
		goTestArgs []string
		wantDir    string
		wantArgs   []string
		wantErr    bool
	}{
		{
			name:       "core package stays at repo root",
			goTestArgs: []string{"./core/..."},
			wantDir:    root,
			wantArgs:   []string{"./core/..."},
		},
		{
			name:       "deployment top-level rewrites pattern to dot-slash-dot-dot-dot",
			goTestArgs: []string{"./deployment/..."},
			wantDir:    filepath.Join(root, "deployment"),
			wantArgs:   []string{"./..."},
		},
		{
			name:       "deployment subdirectory rewrites pattern relative to deployment root",
			goTestArgs: []string{"./deployment/ccip/..."},
			wantDir:    filepath.Join(root, "deployment"),
			wantArgs:   []string{"./ccip/..."},
		},
		{
			name:       "flags before pattern are preserved unchanged",
			goTestArgs: []string{"-v", "-timeout=10m", "./deployment/..."},
			wantDir:    filepath.Join(root, "deployment"),
			wantArgs:   []string{"-v", "-timeout=10m", "./..."},
		},
		{
			name:       "dot-slash-dot-dot-dot at repo root stays at repo root",
			goTestArgs: []string{"./..."},
			wantDir:    root,
			wantArgs:   []string{"./..."},
		},
		{
			name:       "no package patterns returns repo root unchanged",
			goTestArgs: []string{"-v", "-count=1"},
			wantDir:    root,
			wantArgs:   []string{"-v", "-count=1"},
		},
		{
			name:       "mixed core and deployment patterns error",
			goTestArgs: []string{"./core/...", "./deployment/..."},
			wantErr:    true,
		},
		{
			name:       "specific deployment package without wildcard",
			goTestArgs: []string{"./deployment/ccip"},
			wantDir:    filepath.Join(root, "deployment"),
			wantArgs:   []string{"./ccip"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir, args, err := resolveModuleDir(root, tc.goTestArgs)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantDir, dir)
			require.Equal(t, tc.wantArgs, args)
		})
	}
}
