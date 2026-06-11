package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLIGetAndTarget(t *testing.T) {
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
	t.Setenv("CHAINLINK_ROOT", dir)
	t.Setenv("TOOL_VERSIONS_FILE", filepath.Join(dir, ".tool-versions"))
	t.Setenv("GO_TOOLS_FILE", filepath.Join(dir, "tools", "go-tools.txt"))

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
			var buf bytes.Buffer
			cmd := newRootCmd()
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs(tt.args)
			require.NoError(t, cmd.Execute(), "args %v", tt.args)
			assert.Equal(t, tt.want, string(bytes.TrimSpace(buf.Bytes())), "args %v", tt.args)
		})
	}
}

func writeManifest(t *testing.T, dir, name, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(dir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
}
