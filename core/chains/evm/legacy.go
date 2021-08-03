package evm

import (
	"fmt"
	"math/big"
	"net/url"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type LegacyEthNodeConfig interface {
	DefaultChainID() *big.Int
	EthereumURL() string
	EthereumHTTPURL() *url.URL
	EthereumSecondaryURLs() []url.URL
}

func ClobberNodesFromEnv(db *gorm.DB, config LegacyEthNodeConfig) error {
	ethChainID := utils.NewBig(config.DefaultChainID())

	stmt := `
INSERT INTO nodes (name, evm_chain_id, ws_url, http_url, send_only, created_at, updated_at) VALUES (?,?,?,?,?,NOW(),NOW())
ON CONFLICT (lower(name)) DO UPDATE SET
evm_chain_id = ?,
ws_url = ?,
http_url = ?,
updated_at = NOW()
`
	primaryWS := config.EthereumURL()
	var primaryHTTP null.String
	if config.EthereumHTTPURL() != nil {
		primaryHTTP = null.StringFrom(config.EthereumHTTPURL().String())
	}
	if err := db.Exec(stmt, "primary-0", ethChainID, primaryWS, primaryHTTP, false, ethChainID, primaryWS, primaryHTTP).Error; err != nil {
		return errors.Wrap(err, "failed to upsert primary-0")
	}

	for i, url := range config.EthereumSecondaryURLs() {
		name := fmt.Sprintf("secondary-%d", i)
		if err := db.Exec(stmt, name, ethChainID, nil, url.String(), false, ethChainID, nil, url).Error; err != nil {
			return errors.Wrapf(err, "failed to upsert %s", name)
		}
	}
	return nil

}
