package keeper

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

var launchAndTestCmd = &cobra.Command{
	Use:   "launch-and-test",
	Short: "Launches keepers and starts performing",
	Long:  `This command launches chainlink nodes, keeper setup and starts performing upkeeps.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		if err := cfg.Validate(); err != nil {
			log.Fatal(err)
		}

		hdlr := handler.NewKeeper(cfg)

		withdraw, err := cmd.Flags().GetBool("withdraw")
		if err != nil {
			log.Fatal("failed to get withdraw flag: ", err)
		}

		printLogs, err := cmd.Flags().GetBool("export-logs")
		if err != nil {
			log.Fatal("failed to get export-logs flag: ", err)
		}

		hdlr.LaunchAndTest(cmd.Context(), withdraw, printLogs)
	},
}

func init() {
	launchAndTestCmd.Flags().BoolP("withdraw", "w", false, "Specify if funds should be withdrawn and upkeeps should be canceled")
	launchAndTestCmd.Flags().BoolP("export-logs", "l", false, "Specify if container logs should be exported to ./")
}
