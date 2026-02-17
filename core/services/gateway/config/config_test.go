package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ServiceConfig
		wantErr string
	}{
		{
			name: "valid config",
			config: ServiceConfig{
				ServiceName: "workflows",
				Handlers:    []Handler{{Name: "http_handler"}},
				DONs:        []string{"don1"},
			},
		},
		{
			name: "missing service name",
			config: ServiceConfig{
				Handlers: []Handler{{Name: "http_handler"}},
				DONs:     []string{"don1"},
			},
			wantErr: "service name is required",
		},
		{
			name: "no DONs",
			config: ServiceConfig{
				ServiceName: "workflows",
				Handlers:    []Handler{{Name: "http_handler"}},
				DONs:        []string{},
			},
			wantErr: "must have at least one DON",
		},
		{
			name: "handler without name",
			config: ServiceConfig{
				ServiceName: "workflows",
				Handlers:    []Handler{{Name: ""}},
				DONs:        []string{"don1"},
			},
			wantErr: "handler 0 in service \"workflows\" has no name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestShardedDONConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ShardedDONConfig
		wantErr string
	}{
		{
			name: "valid single shard",
			config: ShardedDONConfig{
				DonName: "don1",
				F:       1,
				Shards: []Shard{
					{Nodes: []NodeConfig{
						{Name: "node1", Address: "0x1111111111111111111111111111111111111111"},
						{Name: "node2", Address: "0x2222222222222222222222222222222222222222"},
						{Name: "node3", Address: "0x3333333333333333333333333333333333333333"},
						{Name: "node4", Address: "0x4444444444444444444444444444444444444444"},
					}},
				},
			},
		},
		{
			name: "valid multiple shards",
			config: ShardedDONConfig{
				DonName: "don1",
				F:       0,
				Shards: []Shard{
					{Nodes: []NodeConfig{{Name: "node1", Address: "0x1111111111111111111111111111111111111111"}}},
					{Nodes: []NodeConfig{{Name: "node2", Address: "0x2222222222222222222222222222222222222222"}}},
				},
			},
		},
		{
			name:    "missing DON name",
			config:  ShardedDONConfig{F: 1, Shards: []Shard{{Nodes: []NodeConfig{}}}},
			wantErr: "DON name is required",
		},
		{
			name: "negative F",
			config: ShardedDONConfig{
				DonName: "don1",
				F:       -1,
				Shards:  []Shard{{Nodes: []NodeConfig{}}},
			},
			wantErr: "invalid F value",
		},
		{
			name: "no shards",
			config: ShardedDONConfig{
				DonName: "don1",
				F:       1,
				Shards:  []Shard{},
			},
			wantErr: "must have at least one shard",
		},
		{
			name: "empty shard",
			config: ShardedDONConfig{
				DonName: "don1",
				F:       0,
				Shards:  []Shard{{Nodes: []NodeConfig{}}},
			},
			wantErr: "shard 0 has no nodes",
		},
		{
			name: "insufficient nodes for F",
			config: ShardedDONConfig{
				DonName: "don1",
				F:       1,
				Shards: []Shard{
					{Nodes: []NodeConfig{
						{Name: "node1", Address: "0x1111111111111111111111111111111111111111"},
						{Name: "node2", Address: "0x2222222222222222222222222222222222222222"},
					}},
				},
			},
			wantErr: "has 2 nodes, but requires at least 4",
		},
		{
			name: "node without address",
			config: ShardedDONConfig{
				DonName: "don1",
				F:       0,
				Shards: []Shard{
					{Nodes: []NodeConfig{{Name: "node1", Address: ""}}},
				},
			},
			wantErr: "node 0 has no address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGatewayConfig_Validate(t *testing.T) {
	validDON := ShardedDONConfig{
		DonName: "don1",
		F:       0,
		Shards: []Shard{
			{Nodes: []NodeConfig{{Name: "node1", Address: "0x1111111111111111111111111111111111111111"}}},
		},
	}

	validService := ServiceConfig{
		ServiceName: "workflows",
		Handlers:    []Handler{{Name: "http_handler"}},
		DONs:        []string{"don1"},
	}

	tests := []struct {
		name    string
		config  GatewayConfig
		wantErr string
	}{
		{
			name: "valid config",
			config: GatewayConfig{
				Services:    []ServiceConfig{validService},
				ShardedDONs: []ShardedDONConfig{validDON},
			},
		},
		{
			name: "multiple services sharing DON",
			config: GatewayConfig{
				Services: []ServiceConfig{
					{ServiceName: "workflows", Handlers: []Handler{{Name: "h1"}}, DONs: []string{"don1"}},
					{ServiceName: "vault", Handlers: []Handler{{Name: "h2"}}, DONs: []string{"don1"}},
				},
				ShardedDONs: []ShardedDONConfig{validDON},
			},
		},
		{
			name: "service with multiple DONs",
			config: GatewayConfig{
				Services: []ServiceConfig{
					{ServiceName: "workflows", Handlers: []Handler{{Name: "h1"}}, DONs: []string{"don1", "don2"}},
				},
				ShardedDONs: []ShardedDONConfig{
					validDON,
					{DonName: "don2", F: 0, Shards: []Shard{{Nodes: []NodeConfig{{Name: "n", Address: "0x2222222222222222222222222222222222222222"}}}}},
				},
			},
		},
		{
			name: "duplicate DON name",
			config: GatewayConfig{
				Services: []ServiceConfig{validService},
				ShardedDONs: []ShardedDONConfig{
					validDON,
					{DonName: "don1", F: 0, Shards: []Shard{{Nodes: []NodeConfig{{Name: "n", Address: "0x2222222222222222222222222222222222222222"}}}}},
				},
			},
			wantErr: "duplicate DON name: don1",
		},
		{
			name: "duplicate service name",
			config: GatewayConfig{
				Services: []ServiceConfig{
					validService,
					{ServiceName: "workflows", Handlers: []Handler{{Name: "h2"}}, DONs: []string{"don1"}},
				},
				ShardedDONs: []ShardedDONConfig{validDON},
			},
			wantErr: "duplicate service name: workflows",
		},
		{
			name: "service references unknown DON",
			config: GatewayConfig{
				Services: []ServiceConfig{
					{ServiceName: "workflows", Handlers: []Handler{{Name: "h1"}}, DONs: []string{"unknown_don"}},
				},
				ShardedDONs: []ShardedDONConfig{validDON},
			},
			wantErr: "references unknown DON: unknown_don",
		},
		{
			name:   "empty config is valid",
			config: GatewayConfig{},
		},
		{
			name: "legacy Dons and Services cannot be used together",
			config: GatewayConfig{
				Dons:     []DONConfig{{DonId: "legacy"}},
				Services: []ServiceConfig{validService},
			},
			wantErr: "legacy Dons config and Services/ShardedDONs cannot be used together",
		},
		{
			name: "legacy Dons and ShardedDONs cannot be used together",
			config: GatewayConfig{
				Dons:        []DONConfig{{DonId: "legacy"}},
				ShardedDONs: []ShardedDONConfig{validDON},
			},
			wantErr: "legacy Dons config and Services/ShardedDONs cannot be used together",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
