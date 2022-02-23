package keeper

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

// deployCmd represents the command to run the service
var upkeepCmd = &cobra.Command{
	Use:   "upkeep-events",
	Short: "Print upkeep perform events(stdout and csv file)",
	Long:  `Print upkeep perform events and write to a csv file. args = hexaddr, fromBlock, toBlock`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)
		fromBlock, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		toBlock, err := strconv.ParseUint(args[2], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		hdlr.UpkeepCounterEvents(cmd.Context(), args[0], fromBlock, toBlock)
	},
}
