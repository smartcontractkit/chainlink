package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommandPathShowsGoCInvocation(t *testing.T) {
	t.Parallel()

	if got := rootCmd.CommandPath(); got != "go tool test" {
		t.Fatalf("root CommandPath (help / errors): got %q want %q", got, "go tool test")
	}
	if got := rootCmd.DisplayName(); got != "go tool test" {
		t.Fatalf("DisplayName: got %q want %q", got, "go tool test")
	}
	if got := rootCmd.Name(); got != "test" {
		t.Fatalf("internal Name (subcommand paths use CommandPath + Name): got %q want %q", got, "test")
	}
}

func TestSubcommandCommandPaths(t *testing.T) {
	t.Parallel()

	var gotestsum *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "gotestsum" {
			gotestsum = c
			break
		}
	}
	if gotestsum == nil {
		t.Fatal("gotestsum subcommand not found")
	}
	want := "go tool test gotestsum"
	if got := gotestsum.CommandPath(); got != want {
		t.Fatalf("gotestsum CommandPath: got %q want %q", got, want)
	}
}
