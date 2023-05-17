package keeper

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

// upkeepEventsCmd represents the command to run the upkeep events counter command
// In order to use this command, deploy, register, and fund the UpkeepCounter contract and run this command after it
// emits events on chain.
var verifiableLoad = &cobra.Command{
	Use:   "verifiable-load",
	Short: "Print upkeep perform events(stdout and csv file)",
	Long:  `Print upkeep perform events and write to a csv file. args = hexaddr, fromBlock, toBlock`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)
		hdlr.GetVerifiableLoadStats(cmd.Context())
	},
}
