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

const missingEnvVarMsg = `USE_LEGACY_ETH_ENV_VARS is on but a required env var was missing: %s

PLEASE READ THIS ADDITIONAL INFO

Chainlink now supports multiple chains, so the way ETH nodes are configured has changed. From version 1.1.0 and up, eth node configuration is stored in the database.

The following environment variables are deprecated:

- ETH_URL
- ETH_HTTP_URL
- ETH_SECONDARY_URLS

If you wish to continue using these environment variables (as it used to work in 1.0.0 and below) you must ensure that the following are set:

- USE_LEGACY_ETH_ENV_VARS=true
- ETH_CHAIN_ID (mandatory)
- ETH_URL (mandatory)
- ETH_HTTP_URL (optional)
- ETH_SECONDARY_URLS (optional)

This will automatically overwrite the database records with the given ENV values every time Chainlink boots.

If, instead, you wish to use the API/CLI/GUI to configure your chains and eth nodes (recommended) you must set the following:

- USE_LEGACY_ETH_ENV_VARS=false (will be default in a future release)

This will cause Chainlink to ignore the values for ETH_URL, ETH_HTTP_URL and ETH_SECONDARY_URLS, and use only the database for its node configuration.

For more information on configuring your node, check the docs: https://docs.chain.link/docs/configuration-variables/
`

func ClobberDBFromEnv(db *sqlx.DB, config LegacyEthNodeConfig, lggr logger.Logger) error {
	ethChainID := utils.NewBig(config.DefaultChainID())
	if ethChainID == nil {
		return errors.Errorf(missingEnvVarMsg, "ETH_CHAIN_ID")
	}
	lggr.Infow(fmt.Sprintf("USE_LEGACY_ETH_ENV_VARS is on, upserting chain %s and replacing primary/send-only nodes. It is recommended "+
		"to set USE_LEGACY_ETH_ENV_VARS=false on subsequent runs and use the API to administer chains/nodes instead", ethChainID.String()),
		"evmChainID", ethChainID.String())

	if _, err := db.Exec("INSERT INTO evm_chains (id, created_at, updated_at) VALUES ($1, NOW(), NOW()) ON CONFLICT DO NOTHING;", ethChainID.String()); err != nil {
		return errors.Wrap(err, "failed to insert evm_chain")
	}

	if _, err := db.Exec("DELETE FROM nodes WHERE evm_chain_id = $1", ethChainID.String()); err != nil {
		return errors.Wrap(err, "failed to insert evm_chain")
	}

	stmt := `INSERT INTO nodes (name, evm_chain_id, ws_url, http_url, send_only, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,NOW(),NOW())`
	primaryWS := config.EthereumURL()
	if primaryWS == "" {
		return errors.Errorf(missingEnvVarMsg, "ETH_URL")
	}
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
