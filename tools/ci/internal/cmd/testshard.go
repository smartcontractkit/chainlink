package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/testshard"
)

func newTestshardCmd(stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	testshardCmd := &cobra.Command{
		Use:   "testshard",
		Short: "Shard and verify test package lists",
		Long:  "Commands to partition package lists across parallel CI runners using deterministic hashing.",
	}

	testshardCmd.AddCommand(newListCmd(stdin, stdout))
	testshardCmd.AddCommand(newVerifyCmd(stdin, stdout))

	return testshardCmd
}

func newListCmd(stdin io.Reader, stdout io.Writer) *cobra.Command {
	var (
		shardCount int
		shardIndex int
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Filter stdin package list to those assigned to shard index",
		Example: `  printf "pkg1\npkg2\n" | go run ./tools/ci testshard list --shard-count=4 --shard-index=0`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return testshard.List(stdin, stdout, shardCount, shardIndex)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&shardCount, "shard-count", 0, "Total number of shards (>= 1)")
	flags.IntVar(&shardIndex, "shard-index", 0, "0-based index of this shard [0, shard-count)")
	_ = cmd.MarkFlagRequired("shard-count")
	_ = cmd.MarkFlagRequired("shard-index")

	return cmd
}

func newVerifyCmd(stdin io.Reader, stdout io.Writer) *cobra.Command {
	var shardCount int

	cmd := &cobra.Command{
		Use:     "verify",
		Short:   "Verify complete, non-overlapping shard coverage of stdin package list",
		Example: `  printf "pkg1\npkg2\n" | go run ./tools/ci testshard verify --shard-count=4`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return testshard.Verify(stdin, stdout, shardCount)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&shardCount, "shard-count", 0, "Total number of shards (>= 1)")
	_ = cmd.MarkFlagRequired("shard-count")

	return cmd
}
