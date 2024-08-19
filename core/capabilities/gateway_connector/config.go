package gateway_connector

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
)

// TODO: this will be likely a part of node's TOML config
// additionally, we might need to figure out how to import config from Standard Capabilities' job specs
type WorkflowConnectorConfig struct {
	// TODO: more specific config goes here (e.g. allowlist, rate limits, anything capability specific)
	ChainIDForNodeKey      *big.Int
	GatewayConnectorConfig *connector.ConnectorConfig `json:"gatewayConnectorConfig"`
}
