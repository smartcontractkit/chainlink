package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCLIGetAndTarget(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, ".tool-versions", `mockery 2.53.0
protoc 29.3
golangci-lint 2.12.2
golang 1.26.4
`)
	if err := os.MkdirAll(filepath.Join(dir, "tools"), 0o755); err != nil {
		t.Fatal(err)
	}
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
		var buf bytes.Buffer
		cmd := newRootCmd()
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)
		cmd.SetArgs(tt.args)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("args %v: %v", tt.args, err)
		}
		got := bytes.TrimSpace(buf.Bytes())
		if string(got) != tt.want {
			t.Errorf("args %v = %q, want %q", tt.args, got, tt.want)
		}
	}
}

func writeManifest(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
