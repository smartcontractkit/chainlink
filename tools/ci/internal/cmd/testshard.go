package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/testshard"
)

func newTestshardCmd(stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	var (
		shardCount int
		shardIndex int
	)

	cmd := &cobra.Command{
		Use:   "testshard",
		Short: "Assign or verify Go packages into deterministic test shards",
		Long:  "Reads package paths from stdin and either filters for a specific shard index or verifies shard coverage.",
	}

	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "Filter stdin package paths for a specific shard index",
		Example: "  printf '%s\\n' pkgA pkgB | go run ./tools/ci testshard list --shard-count=4 --shard-index=0",
		RunE: func(cmd *cobra.Command, args []string) error {
			return testshard.List(stdin, stdout, shardCount, shardIndex)
		},
	}
	listCmd.Flags().IntVar(&shardCount, "shard-count", 0, "Total number of shards (must be >= 1)")
	listCmd.Flags().IntVar(&shardIndex, "shard-index", -1, "Target shard index to list (0-based, must be in [0, shard-count))")
	_ = listCmd.MarkFlagRequired("shard-count")
	_ = listCmd.MarkFlagRequired("shard-index")

	verifyCmd := &cobra.Command{
		Use:     "verify",
		Short:   "Verify that all stdin package paths are assigned across shards without duplicates",
		Example: "  printf '%s\\n' pkgA pkgB | go run ./tools/ci testshard verify --shard-count=4",
		RunE: func(cmd *cobra.Command, args []string) error {
			return testshard.Verify(stdin, stdout, shardCount)
		},
	}
	verifyCmd.Flags().IntVar(&shardCount, "shard-count", 0, "Total number of shards (must be >= 1)")
	_ = verifyCmd.MarkFlagRequired("shard-count")

	cmd.AddCommand(listCmd)
	cmd.AddCommand(verifyCmd)

	return cmd
}
