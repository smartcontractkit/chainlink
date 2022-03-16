package keeper

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

// upkeepEventsCmd represents the command to run the upkeep events counter command
var upkeepEventsCmd = &cobra.Command{
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

// upkeepHistoryCmd represents the command to run the upkeep history command
var upkeepHistoryCmd = &cobra.Command{
	Use:   "upkeep-history",
	Short: "Print checkUpkeep history",
	Long:  `Print checkUpkeep status and keeper responsibility for a given upkeep in a set block range`,
	Run: func(cmd *cobra.Command, args []string) {
		upkeepId, err := cmd.Flags().GetInt64("upkeep-id")
		if err != nil {
			log.Fatal("failed to get 'upkeep-id' flag: ", err)
		}

		fromBlock, err := cmd.Flags().GetUint64("from")
		if err != nil {
			log.Fatal("failed to get 'from' flag: ", err)
		}

		toBlock, err := cmd.Flags().GetUint64("to")
		if err != nil {
			log.Fatal("failed to get 'to' flag: ", err)
		}

		gasPrice, err := cmd.Flags().GetUint64("gas-price")
		if err != nil {
			log.Fatal("failed to get 'gas-price' flag: ", err)
		}

		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)

		hdlr.UpkeepHistory(cmd.Context(), upkeepId, fromBlock, toBlock, gasPrice)
	},
}

func init() {
	upkeepHistoryCmd.Flags().Int64("upkeep-id", 0, "upkeep ID")
	upkeepHistoryCmd.Flags().Uint64("from", 0, "from block")
	upkeepHistoryCmd.Flags().Uint64("to", 0, "to block")
	upkeepHistoryCmd.Flags().Uint64("gas-price", 0, "gas price to use")
}
