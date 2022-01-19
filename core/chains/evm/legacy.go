package evm

import (
	"fmt"
	"math/big"
	"net/url"

	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/sqlx"
)

type LegacyEthNodeConfig interface {
	DefaultChainID() *big.Int
	EthereumURL() string
	EthereumHTTPURL() *url.URL
	EthereumSecondaryURLs() []url.URL
}

const missingEthChainIDMsg = `missing ETH_CHAIN_ID; this env var is required if ETH_URL is set

PLEASE READ THIS ADDITIONAL INFO

Chainlink now supports multiple chains, so the way ETH nodes are configured has changed. From version 1.1.0 and up, eth node configuration is stored in the database.

The following environment variables are DEPRECATED:

- ETH_URL
- ETH_HTTP_URL
- ETH_SECONDARY_URLS

Setting ETH_URL will cause Chainlink to automatically overwrite the database records with the given ENV values every time Chainlink boots. This behavior is used mainly to ease the process of upgrading from older versions, and on subsequent runs (once your old settings have been written to the database) it is recommended to unset these ENV vars and use the API commands exclusively to administer chains and nodes.

If you wish to continue using these environment variables (as it used to work in 1.0.x and below) you must ensure that the following are set:

- ETH_CHAIN_ID (mandatory) <--- CURRENTLY MISSING
- ETH_URL (mandatory)
- ETH_HTTP_URL (optional)
- ETH_SECONDARY_URLS (optional)

If, instead, you wish to use the API/CLI/GUI to configure your chains and eth nodes (recommended) you must REMOVE the following environment variables:

- ETH_URL
- ETH_HTTP_URL
- ETH_SECONDARY_URLS

This will cause Chainlink to use the database for its node configuration.

NOTE: ETH_CHAIN_ID remains optional if ETH_URL is left unset. If provided, it specifies the default chain to use in a multichain environment (if you leave ETH_CHAIN_ID unset, the default chain is simply the "first").

For more information on configuring your node, check the docs: https://docs.chain.link/docs/configuration-variables/
`

func ClobberDBFromEnv(db *sqlx.DB, config LegacyEthNodeConfig, lggr logger.Logger) error {
	primaryWS := config.EthereumURL()
	if primaryWS == "" {
		return nil
	}

	ethChainID := utils.NewBig(config.DefaultChainID())
	if ethChainID == nil {
		return errors.New(missingEthChainIDMsg)
	}
	lggr.Warnw(fmt.Sprintf("ETH_URL was set, automatically inserting/updating chain %s and its primary/send-only nodes. It is recommended "+
		"to unset ETH_URL on subsequent runs and use the API to administer chains/nodes instead", ethChainID.String()),
		"evmChainID", ethChainID.String())

	if _, err := db.Exec("INSERT INTO evm_chains (id, created_at, updated_at) VALUES ($1, NOW(), NOW()) ON CONFLICT DO NOTHING;", ethChainID.String()); err != nil {
		return errors.Wrap(err, "failed to insert evm_chain")
	}

	if _, err := db.Exec("DELETE FROM nodes WHERE evm_chain_id = $1", ethChainID.String()); err != nil {
		return errors.Wrap(err, "failed to insert evm_chain")
	}

	stmt := `INSERT INTO nodes (name, evm_chain_id, ws_url, http_url, send_only, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,NOW(),NOW())`
	var primaryHTTP null.String
	if config.EthereumHTTPURL() != nil {
		primaryHTTP = null.StringFrom(config.EthereumHTTPURL().String())
	}
	if _, err := db.Exec(stmt, fmt.Sprintf("primary-0-%s", ethChainID), ethChainID, primaryWS, primaryHTTP, false); err != nil {
		return errors.Wrap(err, "failed to upsert primary-0")
	}

	for i, url := range config.EthereumSecondaryURLs() {
		name := fmt.Sprintf("sendonly-%d-%s", i, ethChainID)
		if _, err := db.Exec(stmt, name, ethChainID, nil, url.String(), true); err != nil {
			return errors.Wrapf(err, "failed to upsert %s", name)
		}
	}
	return nil

}
