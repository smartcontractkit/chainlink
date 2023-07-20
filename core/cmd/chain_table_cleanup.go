package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
	"github.com/smartcontractkit/sqlx"
	"github.com/urfave/cli"
	"net/url"
	"strings"
)

// CleanupChainTables deletes database table rows based on chain type and chain id input.
func (s *Shell) CleanupChainTables(c *cli.Context) error {
	cfg := s.Config.Database()
	parsed := cfg.URL()
	if parsed.String() == "" {
		return s.errorOut(errDBURLMissing)
	}

	fmt.Println("val ", c.Bool("dangerWillRobinson"))
	dbname := parsed.Path[1:]
	if !c.Bool("dangerWillRobinson") && !strings.HasSuffix(dbname, "_test") {
		return s.errorOut(fmt.Errorf("cannot reset database named `%s`. This command can only be run against databases with a name that ends in `_test`, to prevent accidental data loss. If you really want to delete chain specific data from this database, pass in the -dangerWillRobinson option", dbname))
	}

	db, err := getDBConnection(cfg.URL())
	if err != nil {

		return s.errorOut(errors.Wrap(err, "Error connecting to the database"))
	}

	defer db.Close()

	// Delete rows from each table based on the chain_id.
	if strings.EqualFold("EVM", c.String("chaintype")) {
		tables, err := getTablesContainingColumn(db, "evm_chain_id")
		if err != nil {
			return s.errorOut(errors.Wrap(err, "failed to get tables"))
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
	return nil
}

func getDBConnection(clDatabaseURL url.URL) (*sqlx.DB, error) {
	cfg, err := chainlink.GeneralConfigOpts{}.New()
	if err != nil {
		return nil, err
	}

	dbConn, err := pg.NewConnection(clDatabaseURL.String(), dialects.Postgres, cfg.Database())
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
