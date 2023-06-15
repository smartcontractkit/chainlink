package gateway

import (
	"encoding/json"

	gw_net "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

type GatewayConfig struct {
	UserServerConfig gw_net.HTTPServerConfig
	NodeServerConfig gw_net.WebSocketServerConfig
	Dons             []DONConfig
}

type DONConfig struct {
	DonId         string
	HandlerName   string
	HandlerConfig json.RawMessage
	Members       []NodeConfig
}

type NodeConfig struct {
	Name    string
	Address string
}
