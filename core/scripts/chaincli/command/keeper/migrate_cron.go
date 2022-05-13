package keeper

import (
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"

	forwarder "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/authorized_forwarder"
)

var inputFile string

var MigrateCronCmd = &cobra.Command{
	Use:   "migrate-cron",
	Short: "Migrate feed util jobs to use Cron upkeep",
	Long: `This command reads in a list of feed contracts from input file, creates a new cron keeper contract and registers the upkeep` +
		`Creates an output file migrate_cron_output.csv with format: targetAddress,targetFunction,cronSchedule,` +
		`upkeep_name,cronUpkeepAddress,gasLimit,admin,RegistrationHash,blockNum`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)

		if inputFile == "" {
			log.Fatal("Input file should be provided")
		}

		fetchIds, err := cmd.Flags().GetBool("fetch-ids")
		if err != nil {
			log.Fatal("failed to get fetch-ids flag: ", err)
		}

		if fetchIds {
			hdlr.FetchUpkeepIds(cmd.Context(), inputFile)
		} else {
			forwarderAbi, err := abi.JSON(strings.NewReader(forwarder.AuthorizedForwarderABI))
			if err != nil {
				log.Fatalln("Error generating forwarder ABI", err)
			}

			hdlr.MigrateCron(cmd.Context(), inputFile, forwarderAbi)
		}
	},
}

func init() {
	MigrateCronCmd.Flags().StringVar(&inputFile, "input-file", "", "path to csv file in format: targetAddress,targetFunction,cronSchedule,fundingAmountLink,name,encryptedEmail,admin,gasLimit")
	MigrateCronCmd.Flags().BoolP("fetch-ids", "f", false, "Specify to fetch upkeep IDs for registration requests given in input")
}
