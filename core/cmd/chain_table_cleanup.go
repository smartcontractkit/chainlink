package cmd

import (
	"fmt"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
	"github.com/smartcontractkit/sqlx"
	"github.com/urfave/cli"
	"strings"
)

func initChainTablesCleanupCmd(s *Shell) cli.Command {
	return cli.Command{
		Name:   "chaintablecleanup",
		Usage:  "Deletes rows from chain tables based on input",
		Action: s.IndexTxAttempts,
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:     "chainid",
				Usage:    "chainID based on which table cleanup will be done",
				Required: true,
			},
			cli.StringFlag{
				Name:     "chaintype",
				Usage:    "Chain type based on which table cleanup will be done, eg. EVM",
				Required: true,
			},
			cli.StringFlag{
				Name:     "dburl",
				Usage:    "Database URL from which tables will be deleted",
				Required: true,
			},
		},
	}

}

// CleanupChainTables deletes database table rows based on sure chain type and chain id input.
func (s *Shell) CleanupChainTables(c *cli.Context) {
	db, err := getDBConnection(c.String("dburl"))
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	if strings.EqualFold("EVM", c.String("chaintype")) {
		// Delete rows from each table based on the chain_id.
		tables, err := getTablesContainingColumn(db, "evm_chain_id")
		if err != nil {
			fmt.Println("failed to get tables:", err)
			return
		}
		for _, tableName := range tables {
			query := fmt.Sprintf("DELETE FROM %s WHERE evm_chain_id = $1", tableName)
			_, err := db.Exec(query, c.Int("chainid"))
			if err != nil {
				fmt.Printf("Error deleting rows from %s: %v\n", tableName, err)
			} else {
				fmt.Printf("Rows with chain_id %d deleted from %s.\n", c.Int("chainid"), tableName)
			}
		}
	}

}

func getDBConnection(clDatabaseURL string) (*sqlx.DB, error) {
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

func getTablesContainingColumn(db *sqlx.DB, columnName string) ([]string, error) {
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err = rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	var tablesWithColumn []string
	for _, table := range tables {
		has, err := hasColumn(db, table, columnName)
		if err != nil {
			return nil, err
		}
		if has {
			tablesWithColumn = append(tablesWithColumn, table)
		}
	}
	return tablesWithColumn, nil
}

func hasColumn(db *sqlx.DB, table, columnName string) (bool, error) {
	query := fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name='%s' AND column_name='%s'", table, columnName)
	rows, err := db.Query(query)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}
