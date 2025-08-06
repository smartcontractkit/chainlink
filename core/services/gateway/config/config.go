package config

import (
	"encoding/json"

	gw_net "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

type GatewayConfig struct {
	UserServerConfig        gw_net.HTTPServerConfig
	NodeServerConfig        gw_net.WebSocketServerConfig
	ConnectionManagerConfig ConnectionManagerConfig
	// HTTPClientConfig is configuration for outbound HTTP calls to external endpoints
	HTTPClientConfig gw_net.HTTPClientConfig
	Dons             []DONConfig
}

type ConnectionManagerConfig struct {
	AuthGatewayId             string
	AuthTimestampToleranceSec uint32
	AuthChallengeLen          uint32
	HeartbeatIntervalSec      uint32
}

type DONConfig struct {
	DonId   string
	Members []NodeConfig
	F       int

	// Deprecated: use HandlerConfigs instead
	HandlerName string
	// Deprecated: use HandlerConfigs instead
	HandlerConfig json.RawMessage

	Handlers []Handler
}

type Handler struct {
	Name   string
	Config json.RawMessage
}

type NodeConfig struct {
	Name    string
	Address string
}
