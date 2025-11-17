package main

import (
	"bytes"
	"testing"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/root"
)

func TestMyCommand(t *testing.T) {
	t.Skip("manual test")
	cmd := root.RootCmd
	b := &bytes.Buffer{}
	cmd.SetOut(b)
	// Set command arguments and flags as needed for your use case
	cmd.SetArgs([]string{"env", "start"})

	// Set a breakpoint here to debug, or in the main.go file on the appropriate subcommand
	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
