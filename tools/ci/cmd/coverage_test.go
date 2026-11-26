package cmd

import (
	"bytes"
	"testing"
)

func TestCoverageCommandIsRegistered(t *testing.T) {
	t.Parallel()
	root := NewRootCmd()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"coverage", "--help"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, command := range []string{"pr", "publish-pr", "collect", "badge", "publish-badge"} {
		if !bytes.Contains(out.Bytes(), []byte(command)) {
			t.Errorf("coverage help does not list %q: %s", command, out.String())
		}
	}
}

func TestCoveragePRValidatesRequiredInput(t *testing.T) {
	t.Parallel()
	root := NewRootCmd()
	root.SetArgs([]string{"coverage", "pr", "--repo", t.TempDir()})
	if err := root.Execute(); err == nil {
		t.Fatal("expected missing base/head SHA to fail")
	}
}
