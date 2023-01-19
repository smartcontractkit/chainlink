package keeper

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
)

// upkeepEventsCmd represents the command to run the upkeep events counter command
// In order to use this command, deploy, register, and fund the UpkeepCounter contract and run this command after it
// emits events on chain.
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
		upkeepIdStr, err := cmd.Flags().GetString("upkeep-id")
		if err != nil {
			log.Fatal("failed to get 'upkeep-id' flag: ", err)
		}
		upkeepId, ok := keeper.ParseUpkeepId(upkeepIdStr)
		if !ok {
			log.Fatal("failed to parse upkeep-id")
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

var ocr2UpkeepReportHistoryCmd = &cobra.Command{
	Use:   "ocr2-reports",
	Short: "Print ocr2 automation reports",
	Long:  "Print ocr2 automation reports within specified range for registry address",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewBaseHandler(cfg)

		var hashes []string
		var err error
		var path string

		path, err = cmd.Flags().GetString("csv")
		if err == nil && len(path) != 0 {
			rec, err := readCsvFile(path)
			if err != nil {
				log.Fatal(err)
			}

			if len(rec) < 1 {
				log.Fatal("not enough records")
			}

			hashes = make([]string, len(rec))
			for i := 0; i < len(rec); i++ {
				hashes[i] = rec[i][0]
			}
		} else {
			hashes, err = cmd.Flags().GetStringSlice("tx-hashes")
			if err != nil {
				log.Fatalf("failed to get transaction hashes from input: %s", err)
			}
		}

		if err := handler.OCR2AutomationReports(hdlr, hashes); err != nil {
			log.Fatalf("failed to collect transaction data: %s", err)
		}
	},
}

var ocr2UpdateConfigCmd = &cobra.Command{
	Use:   "ocr2-get-config",
	Short: "Get OCR2 config parameters",
	Long:  "Get latest OCR2 config parameters from registry contract address",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.New()
		hdlr := handler.NewBaseHandler(cfg)

		if err := handler.OCR2GetConfig(hdlr, cfg.RegistryAddress); err != nil {
			log.Fatalf("failed to get config data: %s", err)
		}
	},
}

func init() {
	upkeepHistoryCmd.Flags().String("upkeep-id", "", "upkeep ID")
	upkeepHistoryCmd.Flags().Uint64("from", 0, "from block")
	upkeepHistoryCmd.Flags().Uint64("to", 0, "to block")
	upkeepHistoryCmd.Flags().Uint64("gas-price", 0, "gas price to use")

	ocr2UpkeepReportHistoryCmd.Flags().StringSlice("tx-hashes", []string{}, "list of transaction hashes to get information for")
	ocr2UpkeepReportHistoryCmd.Flags().String("csv", "", "path to csv file containing transaction hashes; first element per line should be transaction hash; file should not have headers")

	ocr2UpdateConfigCmd.Flags().String("tx", "", "transaction of last config update")
}

func readCsvFile(filePath string) ([][]string, error) {
	var records [][]string
	var err error

	f, err := os.Open(filePath)
	if err != nil {
		return records, fmt.Errorf("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = 0
	csvReader.LazyQuotes = false
	records, err = csvReader.ReadAll()
	if err != nil {
		return records, fmt.Errorf("Unable to parse file as CSV for "+filePath, err)
	}

	return records, nil
}
