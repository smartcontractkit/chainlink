package cmd_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/cmd"
)

func BenchmarkCommandsHelp(b *testing.B) {
	commands := [][]string{
		{"--help"},
		{"lint", "--help"},
		{"test", "--help"},
		{"tidy", "--help"},
		{"generate", "--help"},
		{"end-of-file-fixer", "--help"},
		{"whitespace-fixer", "--help"},
	}

	for _, args := range commands {
		name := args[0]
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				root := cmd.NewRootCmd()
				root.SetOut(io.Discard)
				root.SetErr(io.Discard)
				root.SetArgs(args)

				err := root.Execute()
				require.NoError(b, err)
			}
		})
	}
}
