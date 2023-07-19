package command

import (
	"fmt"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
	"github.com/smartcontractkit/sqlx"
	"github.com/spf13/cobra"
	"strings"
)

var (
	chainType     string
	chainID       int
	clDatabaseURL string
)

var tablesToCleanup = map[string][]string{
	"EVM": {"evm_log_poller_filters",
		"evm_log_poller_blocks",
		"log_broadcasts",
		"block_header_feeder_specs",
		"direct_request_specs",
		"evm_logs",
		"vrf_specs",
		"evm_heads",
		"evm_forwarders",
		"blockhash_store_specs",
		"evm_key_states",
		"log_broadcasts_pending",
		"eth_txes",
		"keeper_specs",
		"flux_monitor_specs",
		"ocr_oracle_specs"},
}

var ChainTablesCleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Deletes rows from chain tables",
	Long:  "Deletes rows from chain tables based on input",
	Run: func(cmd *cobra.Command, args []string) {
		if !chainTypeExists(&chainType) {
			fmt.Printf("Chain type %s doesn't exist\n", chainType)
			return
		}

		db, err := getDBConnection()
		if err != nil {
			fmt.Println("Error connecting to the database:", err)
			return
		}
		defer db.Close()

		switch chainType {
		case "EVM":
			// Delete rows from each table based on the chain_id.
			for _, tableName := range tablesToCleanup[chainType] {
				query := fmt.Sprintf("DELETE FROM %s WHERE evm_chain_id = %d", tableName, chainID)
				_, err := db.Exec(query)
				if err != nil {
					fmt.Printf("Error deleting rows from %s: %v\n", tableName, err)
				} else {
					fmt.Printf("Rows with chain_id %d deleted from %s.\n", chainID, tableName)
				}
			}
		}
	},
}

func chainTypeExists(chainType *string) bool {
	for name := range tablesToCleanup {
		if strings.EqualFold(name, *chainType) {
			*chainType = name
			return true
		}
	}
	return false
}

func getDBConnection() (*sqlx.DB, error) {
	cfg, err := chainlink.GeneralConfigOpts{}.New()
	if err != nil {
		return nil, err
	}

	dbConn, err := pg.NewConnection(clDatabaseURL, dialects.Postgres, cfg.Database())
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

func init() {
	ChainTablesCleanupCmd.Flags().StringVar(&chainType, "chainType", "", "Chain Type")
	ChainTablesCleanupCmd.Flags().IntVar(&chainID, "chainID", 0, "Chain ID")
	ChainTablesCleanupCmd.Flags().StringVar(&clDatabaseURL, "clDatabaseURL", "", "Database URL")

	ChainTablesCleanupCmd.MarkFlagRequired("chainType")
	ChainTablesCleanupCmd.MarkFlagRequired("chainID")
	ChainTablesCleanupCmd.MarkFlagRequired("clDatabaseURL")
	fmt.Println("chainType ", chainType)
	fmt.Println("chainID ", chainID)
	fmt.Println("clDatabaseURL ", clDatabaseURL)
}
