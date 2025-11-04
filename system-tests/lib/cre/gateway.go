package cre

import (
	"strconv"

	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const (
	gatewayIncomingPort = 5002
	gatewayOutgoingPort = 5003
)

func NewGatewayConfig(p infra.Provider, id, gatewayNodeIdx int, isBootstrap bool, uuid, donName string) *GatewayConfiguration {
	return &GatewayConfiguration{
		NodeUUID: uuid,
		Outgoing: Outgoing{
			Path: "/node",
			Port: gatewayOutgoingPort,
			Host: p.InternalGatewayHost(id, isBootstrap, donName),
		},
		Incoming: Incoming{
			Protocol:     "http",
			Path:         "/",
			InternalPort: gatewayIncomingPort,
			ExternalPort: p.ExternalGatewayPort(gatewayIncomingPort),
		},
		AuthGatewayID: "gateway-node-" + strconv.Itoa(gatewayNodeIdx), // reflects what is done in deployment/cre/jobs/pkg/gateway_job.go
	}
}

type GatewayConfiguration struct {
	NodeUUID      string   `toml:"node_uuid" json:"node_uuid"`
	Outgoing      Outgoing `toml:"outgoing" json:"outgoing"`
	Incoming      Incoming `toml:"incoming" json:"incoming"`
	AuthGatewayID string   `toml:"auth_gateway_id" json:"auth_gateway_id"`
}

type Outgoing struct {
	Host string `toml:"host" json:"host"` // do not set, it will be set dynamically
	Path string `toml:"path" json:"path"`
	Port int    `toml:"port" json:"port"`
}

type Incoming struct {
	Protocol     string `toml:"protocol" json:"protocol"` // do not set, it will be set dynamically
	Host         string `toml:"host" json:"host"`         // do not set, it will be set dynamically
	Path         string `toml:"path" json:"path"`
	InternalPort int    `toml:"internal_port" json:"internal_port"`
	ExternalPort int    `toml:"external_port" json:"external_port"`
}
