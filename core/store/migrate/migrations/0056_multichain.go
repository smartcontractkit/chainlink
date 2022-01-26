package migrations

import (
	"database/sql"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/pressly/goose/v3"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func init() {
	goose.AddMigration(Up56, Down56)
}

const up56 = `
CREATE TABLE evm_chains (
	id numeric(78,0) PRIMARY KEY,
	cfg jsonb NOT NULL DEFAULT '{}',
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL
);

CREATE TABLE nodes (
	id serial PRIMARY KEY,
	name varchar(255) NOT NULL CHECK (name != ''),
	evm_chain_id numeric(78,0) NOT NULL REFERENCES evm_chains (id),
	ws_url text CHECK (ws_url != ''),
	http_url text CHECK (http_url != ''),
	send_only bool NOT NULL CONSTRAINT primary_or_sendonly CHECK (
		(send_only AND ws_url IS NULL AND http_url IS NOT NULL)
		OR
		(NOT send_only AND ws_url IS NOT NULL)
	),
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL
);

CREATE INDEX idx_nodes_evm_chain_id ON nodes (evm_chain_id);
CREATE UNIQUE INDEX idx_nodes_unique_name ON nodes (lower(name));
`

const down56 = `
DROP TABLE nodes;
DROP TABLE evm_chains;
`

func Up56(tx *sql.Tx) error {
	if _, err := tx.Exec(up56); err != nil {
		return err
	}
	dbURL := os.Getenv("DATABASE_URL")
	if !strings.Contains(dbURL, "_test") {
		chainIDStr := os.Getenv("ETH_CHAIN_ID")
		if chainIDStr == "" {
			logger.Warn("ETH_CHAIN_ID was not specified, auto-creating chain with id 1")
			chainIDStr = "1"
		}
		chainID, ok := new(big.Int).SetString(chainIDStr, 10)
		if !ok {
			panic(fmt.Sprintf("ETH_CHAIN_ID was invalid, expected a number, got: %s", chainIDStr))
		}
		_, err := tx.Exec("INSERT INTO evm_chains (id, created_at, updated_at) VALUES ($1, NOW(), NOW());", chainID.String())
		return err
	}
	return nil
}

func Down56(tx *sql.Tx) error {
	_, err := tx.Exec(down56)
	if err != nil {
		return err
	}
	return nil
}
