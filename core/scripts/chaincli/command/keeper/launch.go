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

		bootstrap, err := cmd.Flags().GetBool("bootstrap")
		if err != nil {
			log.Fatal("failed to get bootstrap flag: ", err)
		}

		printLogs, err := cmd.Flags().GetBool("export-logs")
		if err != nil {
			log.Fatal("failed to get export-logs flag: ", err)
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			log.Fatal("failed to get force flag: ", err)
		}

		hdlr.LaunchAndTest(cmd.Context(), withdraw, printLogs, force, bootstrap)
	},
}

func init() {
	launchAndTestCmd.Flags().BoolP("withdraw", "w", true, "Specify if funds should be withdrawn and upkeeps should be canceled")
	launchAndTestCmd.Flags().BoolP("bootstrap", "b", false, "Specify if launching bootstrap node is required. Default listen ports(5688, 8000) are used, if you need to use custom ports, please use bootstrap command")
	launchAndTestCmd.Flags().BoolP("export-logs", "l", false, "Specify if container logs should be exported to ./")
	launchAndTestCmd.Flags().BoolP("force", "f", false, "Specify if existing containers should be forcefully removed ./")
}
