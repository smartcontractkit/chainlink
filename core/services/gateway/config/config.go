package config

import (
	"encoding/json"

	gw_net "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

type GatewayConfig struct {
	UserServerConfig        gw_net.HTTPServerConfig
	NodeServerConfig        gw_net.WebSocketServerConfig
	ConnectionManagerConfig ConnectionManagerConfig
	Dons                    []DONConfig
}

type ConnectionManagerConfig struct {
	AuthTimestampToleranceSec uint32
	AuthChallengeLen          uint32
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
