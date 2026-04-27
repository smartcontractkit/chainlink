package repo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootFrom_chainlinkModule(t *testing.T) {
	// This test runs from package dir; walk up to module root (tools/test), then repo root.
	here, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// internal/repo -> internal -> tools/test
	testMod := filepath.Clean(filepath.Join(here, "..", ".."))
	root, err := RootFrom(testMod)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(root) != "chainlink" && !strings.HasSuffix(root, "chainlink") {
		t.Fatalf("unexpected root %q", root)
	}
	// go.mod at root must exist
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatal(err)
	}
}

func TestRootFrom_skipsLeadingCommentsAndBlankLinesInGoMod(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	goMod := `// workspace root

` + "\n\n" + `module github.com/smartcontractkit/chainlink/v2

go 1.26
`
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte(goMod), 0600))
	nested := filepath.Join(root, "a", "b")
	require.NoError(t, os.MkdirAll(nested, 0700))

	got, err := RootFrom(nested)
	require.NoError(t, err)
	require.Equal(t, root, got)
}

func TestModulePathFromGoMod(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		data   string
		want   string
		wantOK bool
	}{
		{name: "leading comment", data: "// hi\n\nmodule github.com/smartcontractkit/chainlink/v2\n", want: rootModulePath, wantOK: true},
		{name: "no module", data: "go 1.26\n", wantOK: false},
		{name: "inline comment", data: "module github.com/smartcontractkit/chainlink/v2 // chainlink\n", want: rootModulePath, wantOK: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, ok := modulePathFromGoMod(tc.data)
			require.Equal(t, tc.wantOK, ok)
			if tc.wantOK {
				require.Equal(t, tc.want, got)
			}
		})
	}
}
