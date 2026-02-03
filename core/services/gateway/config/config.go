package config

import (
	"encoding/json"
	"errors"
	"fmt"

	gw_net "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

type GatewayConfig struct {
	UserServerConfig        gw_net.HTTPServerConfig
	NodeServerConfig        gw_net.WebSocketServerConfig
	ConnectionManagerConfig ConnectionManagerConfig
	// HTTPClientConfig is configuration for outbound HTTP calls to external endpoints
	HTTPClientConfig gw_net.HTTPClientConfig
	Dons             []DONConfig // Deprecated: use Services + ShardedDONs instead

	// Services defines logical groupings of handlers with attached DONs.
	// Each service can have multiple handlers and be associated with multiple DONs.
	Services    []ServiceConfig
	ShardedDONs []ShardedDONConfig
}

type ConnectionManagerConfig struct {
	AuthGatewayId             string
	AuthTimestampToleranceSec uint32
	AuthChallengeLen          uint32
	HeartbeatIntervalSec      uint32
}

type DONConfig struct {
	DonId         string
	HandlerName   string          // Deprecated: use Handlers instead
	HandlerConfig json.RawMessage // Deprecated: use Handlers instead
	Members       []NodeConfig
	F             int
	Handlers      []Handler
}

type Handler struct {
	Name        string
	ServiceName string
	Config      json.RawMessage
}

type NodeConfig struct {
	Name    string
	Address string
}

// ServiceConfig defines a logical service consisting of handlers attached to DONs.
// The service name is a prefix of JSONRPC method names (e.g., "workflows.execute" -> "workflows").
type ServiceConfig struct {
	ServiceName string
	Handlers    []Handler
	// DONs is a list of DON names (referencing ShardedDONConfig.DonName) attached to this service.
	DONs []string
}

type ShardedDONConfig struct {
	DonName string
	F       int
	Shards  []Shard
}
type Shard struct {
	Nodes []NodeConfig
}

func (c *GatewayConfig) Validate() error {
	if len(c.Dons) > 0 && (len(c.Services) > 0 || len(c.ShardedDONs) > 0) {
		return errors.New("legacy Dons config and Services/ShardedDONs cannot be used together")
	}

	donNames := make(map[string]bool)
	for _, don := range c.ShardedDONs {
		if err := don.Validate(); err != nil {
			return err
		}
		if donNames[don.DonName] {
			return fmt.Errorf("duplicate DON name: %s", don.DonName)
		}
		donNames[don.DonName] = true
	}

	serviceNames := make(map[string]bool)
	for _, svc := range c.Services {
		if err := svc.Validate(); err != nil {
			return err
		}
		if serviceNames[svc.ServiceName] {
			return fmt.Errorf("duplicate service name: %s", svc.ServiceName)
		}
		serviceNames[svc.ServiceName] = true

		for _, donRef := range svc.DONs {
			if !donNames[donRef] {
				return fmt.Errorf("service %q references unknown DON: %s", svc.ServiceName, donRef)
			}
		}
	}
	return nil
}

func (s *ServiceConfig) Validate() error {
	if s.ServiceName == "" {
		return errors.New("service name is required")
	}
	if len(s.DONs) == 0 {
		return fmt.Errorf("service %q must have at least one DON", s.ServiceName)
	}
	for i, h := range s.Handlers {
		if h.Name == "" {
			return fmt.Errorf("handler %d in service %q has no name", i, s.ServiceName)
		}
	}
	return nil
}

func (d *ShardedDONConfig) Validate() error {
	if d.DonName == "" {
		return errors.New("DON name is required")
	}
	if d.F < 0 {
		return fmt.Errorf("DON %q has invalid F value: %d", d.DonName, d.F)
	}
	if len(d.Shards) == 0 {
		return fmt.Errorf("DON %q must have at least one shard", d.DonName)
	}
	for i, shard := range d.Shards {
		if len(shard.Nodes) == 0 {
			return fmt.Errorf("DON %q shard %d has no nodes", d.DonName, i)
		}
		requiredNodes := 3*d.F + 1
		if len(shard.Nodes) < requiredNodes {
			return fmt.Errorf("DON %q shard %d has %d nodes, but requires at least %d (3F+1 where F=%d)",
				d.DonName, i, len(shard.Nodes), requiredNodes, d.F)
		}
		for j, node := range shard.Nodes {
			if node.Address == "" {
				return fmt.Errorf("DON %q shard %d node %d has no address", d.DonName, i, j)
			}
		}
	}
	return nil
}
