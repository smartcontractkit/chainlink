package keeper

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Print the logs of your keeper nodes",
	Long:  `This command prints the logs of all keeper nodes.`,

	Run: func(cmd *cobra.Command, args []string) {
		containerPattern, err := cmd.Flags().GetString("container-pattern")
		if err != nil {
			panic(err)
		}
		grep, err := cmd.Flags().GetStringSlice("grep")
		if err != nil {
			panic(err)
		}
		grepv, err := cmd.Flags().GetStringSlice("grepv")
		if err != nil {
			panic(err)
		}
		cfg := config.New()
		keeper := handler.NewKeeper(cfg)

		keeper.PrintLogs(cmd.Context(), containerPattern, grep, grepv)
	},
}

func init() {
	logsCmd.Flags().String("container-pattern", `^/keeper-\d+$`, "Regex pattern of container names to listen to for logs")
	logsCmd.Flags().StringSlice("grep", []string{}, "comma separated list of terms logs must include")
	logsCmd.Flags().StringSlice("grepv", []string{}, "comma separated list of terms logs must not include")
}
