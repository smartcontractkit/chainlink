package main

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/manifest"
	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/paths"
	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/resolve"
)

func writeManifest(t *testing.T, dir, name, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(dir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
}

func makeTestLoader(t *testing.T) loaderFn {
	t.Helper()
	dir := t.TempDir()
	writeManifest(t, dir, ".tool-versions", `mockery 2.53.0
protoc 29.3
golangci-lint 2.12.2
golang 1.26.4
`)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "tools"), 0o755))
	writeManifest(t, filepath.Join(dir, "tools"), "go-tools.txt", `github.com/jmank88/gomods 0.1.7
github.com/smartcontractkit/gencodec 42dc7da8c2874db550e91c656f98d05fca3c2f98
`)
	cfg := paths.Config{
		Root:             dir,
		ToolVersionsFile: filepath.Join(dir, ".tool-versions"),
		GoToolsFile:      filepath.Join(dir, "tools", "go-tools.txt"),
	}
	store, err := manifest.New(cfg.ToolVersionsFile, cfg.GoToolsFile)
	require.NoError(t, err)
	r := resolve.New(store)
	return func() (*resolve.Resolver, paths.Config, error) {
		return r, cfg, nil
	}
}

func runCLI(t *testing.T, load loaderFn, args ...string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	cmd := newRootCmdWithLoader(load)
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return string(bytes.TrimSpace(buf.Bytes())), err
}

func TestCLIGetAndTarget(t *testing.T) {
	t.Parallel()
	load := makeTestLoader(t)

	tests := []struct {
		args []string
		want string
	}{
		{[]string{"get", "mockery"}, "2.53.0"},
		{[]string{"ref", "golangci-lint"}, "v2.12.2"},
		{[]string{"target", "mockery"}, "github.com/vektra/mockery/v2@v2.53.0"},
		{[]string{"target", "github.com/jmank88/gomods"}, "github.com/jmank88/gomods@v0.1.7"},
		{[]string{"target", "github.com/smartcontractkit/gencodec"}, "github.com/smartcontractkit/gencodec@42dc7da8c2874db550e91c656f98d05fca3c2f98"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			got, err := runCLI(t, load, tt.args...)
			require.NoError(t, err, "args %v", tt.args)
			assert.Equal(t, tt.want, got, "args %v", tt.args)
		})
	}
}

func TestCLIUnknownKey(t *testing.T) {
	t.Parallel()
	load := makeTestLoader(t)

	for _, args := range [][]string{
		{"get", "no-such-tool"},
		{"ref", "no-such-tool"},
		{"target", "no-such-tool"},
	} {
		_, err := runCLI(t, load, args...)
		require.Error(t, err, "args %v should fail for unknown key", args)
	}
}

func TestCLIList(t *testing.T) {
	t.Parallel()
	load := makeTestLoader(t)

	got, err := runCLI(t, load, "list")
	require.NoError(t, err)

	lines := strings.Split(got, "\n")
	require.GreaterOrEqual(t, len(lines), 6)
	// .tool-versions entries come first
	assert.Contains(t, lines[0], "mockery")
	// go-tools.txt entries follow
	assert.Contains(t, got, "github.com/jmank88/gomods")
	assert.Contains(t, got, "github.com/smartcontractkit/gencodec")
}

func TestCLIModules(t *testing.T) {
	t.Parallel()
	load := makeTestLoader(t)

	got, err := runCLI(t, load, "modules")
	require.NoError(t, err)

	lines := strings.Split(got, "\n")
	// Must include the modulemap entry and go-tools entries
	assert.Contains(t, lines, "github.com/vektra/mockery/v2")
	assert.Contains(t, lines, "github.com/jmank88/gomods")
	assert.Contains(t, lines, "github.com/smartcontractkit/gencodec")
	// Must be sorted
	assert.True(t, sort.StringsAreSorted(lines), "modules output must be sorted; got %v", lines)
}

func TestCLIMakeVars(t *testing.T) {
	t.Parallel()
	load := makeTestLoader(t)

	got, err := runCLI(t, load, "make-vars")
	require.NoError(t, err)

	assert.Contains(t, got, "GOLANGCI_LINT_VERSION=v2.12.2")
	// protoc uses raw get (no v-prefix)
	assert.Contains(t, got, "PROTOC_VERSION=29.3")
	assert.NotContains(t, got, "PROTOC_VERSION=v")
}
