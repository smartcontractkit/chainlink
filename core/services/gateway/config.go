package gateway

import "encoding/json"

type GatewayConfig struct {
	UserServerConfig UserServerConfig
	NodeServerConfig NodeServerConfig
	Dons             []DONConfig
}

type UserServerConfig struct {
	Port uint16
}

type NodeServerConfig struct {
	Port uint16
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
