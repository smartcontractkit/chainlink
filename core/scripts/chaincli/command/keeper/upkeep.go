package keeper

import (
	"log"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

var (
	// rangeRegexp is the regexp for the block range string format, e.g. 10-100
	rangeRegexp = regexp.MustCompile(`^([0-9]*)?-?([0-9]*)$`)
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
		fromBlock, toBlock := getRange(cmd)
		upkeepId, err := cmd.Flags().GetInt64("upkeep-id")
		if err != nil {
			log.Fatal("failed to get upkeep-id flag: ", err)
		}

		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)

		hdlr.UpkeepHistory(cmd.Context(), upkeepId, fromBlock, toBlock)
	},
}

func getRange(cmd *cobra.Command) (uint64, uint64) {
	rangeRaw, err := cmd.Flags().GetString("range")
	if err != nil {
		log.Fatal("failed to get range flag: ", err)
	}

	matches := rangeRegexp.FindStringSubmatch(rangeRaw)
	if len(matches) != 3 {
		log.Fatal("unexpected matches: ", matches)
	}

	var from, to int

	if len(matches[1]) > 0 {
		from, err = strconv.Atoi(matches[1])
		if err != nil {
			log.Fatal("failed to parse from block: ", err)
		}
	}

	if len(matches[2]) > 0 {
		to, err = strconv.Atoi(matches[2])
		if err != nil {
			log.Fatal("failed to parse to block: ", err)
		}
	}

	return uint64(from), uint64(to)
}

func init() {
	upkeepHistoryCmd.Flags().Int64("upkeep-id", 0, "upkeep ID")
	upkeepHistoryCmd.Flags().String("range", "", "block range")
}
