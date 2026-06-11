package manifest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeManifest(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestParseLine(t *testing.T) {
	t.Parallel()
	tests := []struct {
		line        string
		wantName    string
		wantVersion string
		wantOK      bool
	}{
		{"", "", "", false},
		{"# comment", "", "", false},
		{"golang        1.26.4", "golang", "1.26.4", true},
		{"github.com/jmank88/gomods  0.1.7", "github.com/jmank88/gomods", "0.1.7", true},
	}
	for _, tt := range tests {
		name, version, ok := parseLine(tt.line)
		assert.Equal(t, tt.wantOK, ok, "parseLine(%q) ok", tt.line)
		assert.Equal(t, tt.wantName, name, "parseLine(%q) name", tt.line)
		assert.Equal(t, tt.wantVersion, version, "parseLine(%q) version", tt.line)
	}
}

func TestLookup(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	tv := writeManifest(t, dir, ".tool-versions", `golang 1.26.4
mockery 2.53.0
protoc 29.3
`)
	gt := writeManifest(t, dir, "go-tools.txt", `# comment
github.com/jmank88/gomods 0.1.7
`)

	store, err := New(tv, gt)
	require.NoError(t, err)

	tests := []struct {
		key     string
		want    string
		wantErr bool
	}{
		{"mockery", "2.53.0", false},
		{"protoc", "29.3", false},
		{"github.com/jmank88/gomods", "0.1.7", false},
		{"missing", "", true},
		{"github.com/missing/tool", "", true},
	}
	for _, tt := range tests {
		got, err := store.Lookup(tt.key)
		if tt.wantErr {
			require.Error(t, err, "Lookup(%q)", tt.key)
			continue
		}
		require.NoError(t, err, "Lookup(%q)", tt.key)
		assert.Equal(t, tt.want, got, "Lookup(%q)", tt.key)
	}
}

func TestList(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	tv := writeManifest(t, dir, ".tool-versions", "mockery 2.53.0\n")
	gt := writeManifest(t, dir, "go-tools.txt", "github.com/jmank88/gomods 0.1.7\n")

	store, err := New(tv, gt)
	require.NoError(t, err)

	entries := store.List()
	require.Len(t, entries, 2)
	assert.Equal(t, "mockery", entries[0].Name)
	assert.Equal(t, "github.com/jmank88/gomods", entries[1].Name)
}

func TestListOrder(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	tv := writeManifest(t, dir, ".tool-versions", "golang 1.26.4\nmockery 2.53.0\nprotoc 29.3\n")
	gt := writeManifest(t, dir, "go-tools.txt", "github.com/a/tool 1.0.0\ngithub.com/b/tool 2.0.0\n")

	store, err := New(tv, gt)
	require.NoError(t, err)

	entries := store.List()
	require.Len(t, entries, 5)
	// .tool-versions entries come first, in file order
	assert.Equal(t, "golang", entries[0].Name)
	assert.Equal(t, "mockery", entries[1].Name)
	assert.Equal(t, "protoc", entries[2].Name)
	// go-tools.txt entries follow, in file order
	assert.Equal(t, "github.com/a/tool", entries[3].Name)
	assert.Equal(t, "github.com/b/tool", entries[4].Name)
}

func TestGoToolModulesStable(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	tv := writeManifest(t, dir, ".tool-versions", "golang 1.26.4\n")
	gt := writeManifest(t, dir, "go-tools.txt", "github.com/a/tool 1.0.0\ngithub.com/b/tool 2.0.0\ngithub.com/c/tool 3.0.0\n")

	store, err := New(tv, gt)
	require.NoError(t, err)

	got1 := store.GoToolModules()
	got2 := store.GoToolModules()
	assert.Equal(t, got1, got2, "GoToolModules must be deterministic")
	assert.Equal(t, []string{"github.com/a/tool", "github.com/b/tool", "github.com/c/tool"}, got1)
}

func TestNewFileNotFound(t *testing.T) {
	t.Parallel()
	_, err := New("/nonexistent/.tool-versions", "/nonexistent/go-tools.txt")
	require.Error(t, err)
}
