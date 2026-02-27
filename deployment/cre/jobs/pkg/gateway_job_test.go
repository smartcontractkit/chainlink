package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGateway_Validate(t *testing.T) {
	t.Parallel()

	g := GatewayJob{}
	require.ErrorContains(t, g.Validate(), "must provide job name")

	g.JobName = "AGatewayJob"
	require.ErrorContains(t, g.Validate(), "must provide at least one DON")
}

const (
	expected = `type = 'gateway'
schemaVersion = 1
name = 'Gateway1'
externalJobID = '4657f08a-e8cd-526f-9c13-66bbef7e4e03'
forwardingAllowed = false

[gatewayConfig]
[gatewayConfig.ConnectionManagerConfig]
AuthChallengeLen = 10
AuthGatewayId = 'gateway-node-1'
AuthTimestampToleranceSec = 5
HeartbeatIntervalSec = 20

[[gatewayConfig.ShardedDONs]]
DonName = 'workflow_1'
F = 1

[[gatewayConfig.ShardedDONs.Shards]]
[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0xabc'
Name = 'Node 1'

[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0xdef'
Name = 'Node 2'

[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0xghi'
Name = 'Node 3'

[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0xjkl'
Name = 'Node 4'

[[gatewayConfig.ShardedDONs]]
DonName = 'workflow_2'
F = 0

[[gatewayConfig.ShardedDONs.Shards]]
[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0x2abc'
Name = 'Node 1'

[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0x2def'
Name = 'Node 2'

[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0x2ghi'
Name = 'Node 3'

[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0x2jkl'
Name = 'Node 4'

[[gatewayConfig.Services]]
ServiceName = 'workflows'
DONs = ['workflow_1', 'workflow_2']

[[gatewayConfig.Services.Handlers]]
Name = 'web-api-capabilities'

[gatewayConfig.Services.Handlers.Config]
maxAllowedMessageAgeSec = 1000

[gatewayConfig.Services.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10

[gatewayConfig.HTTPClientConfig]
MaxResponseBytes = 50000000
AllowedPorts = [443]
AllowedSchemes = ['https']
AllowedIPsCIDR = []

[gatewayConfig.NodeServerConfig]
HandshakeTimeoutMillis = 1000
MaxRequestBytes = 100000
Path = '/'
Port = 5003
ReadTimeoutMillis = 1000
RequestTimeoutMillis = 15000
WriteTimeoutMillis = 1000

[gatewayConfig.UserServerConfig]
ContentTypeHeader = 'application/jsonrpc'
MaxRequestBytes = 100000
Path = '/'
Port = 5002
ReadTimeoutMillis = 15000
RequestTimeoutMillis = 15000
WriteTimeoutMillis = 16000
`

	expectedWithVault = `type = 'gateway'
schemaVersion = 1
name = 'Gateway1'
externalJobID = '4657f08a-e8cd-526f-9c13-66bbef7e4e03'
forwardingAllowed = false

[gatewayConfig]
[gatewayConfig.ConnectionManagerConfig]
AuthChallengeLen = 10
AuthGatewayId = 'gateway-node-1'
AuthTimestampToleranceSec = 5
HeartbeatIntervalSec = 20

[[gatewayConfig.ShardedDONs]]
DonName = 'workflow_1'
F = 1

[[gatewayConfig.ShardedDONs.Shards]]
[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0xabc'
Name = 'Node 1'

[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0xdef'
Name = 'Node 2'

[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0xghi'
Name = 'Node 3'

[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0xjkl'
Name = 'Node 4'

[[gatewayConfig.ShardedDONs]]
DonName = 'workflow_2'
F = 0

[[gatewayConfig.ShardedDONs.Shards]]
[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0x2abc'
Name = 'Node 1'

[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0x2def'
Name = 'Node 2'

[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0x2ghi'
Name = 'Node 3'

[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0x2jkl'
Name = 'Node 4'

[[gatewayConfig.Services]]
ServiceName = 'workflows'
DONs = ['workflow_1', 'workflow_2']

[[gatewayConfig.Services.Handlers]]
Name = 'web-api-capabilities'

[gatewayConfig.Services.Handlers.Config]
maxAllowedMessageAgeSec = 1000

[gatewayConfig.Services.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10

[[gatewayConfig.Services]]
ServiceName = 'vault'
DONs = ['workflow_1']

[[gatewayConfig.Services.Handlers]]
Name = 'vault'
ServiceName = 'vault'

[gatewayConfig.Services.Handlers.Config]
requestTimeoutSec = 14

[gatewayConfig.Services.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10

[gatewayConfig.HTTPClientConfig]
MaxResponseBytes = 50000000
AllowedPorts = [443]
AllowedSchemes = ['https']
AllowedIPsCIDR = []

[gatewayConfig.NodeServerConfig]
HandshakeTimeoutMillis = 1000
MaxRequestBytes = 100000
Path = '/'
Port = 5003
ReadTimeoutMillis = 1000
RequestTimeoutMillis = 15000
WriteTimeoutMillis = 1000

[gatewayConfig.UserServerConfig]
ContentTypeHeader = 'application/jsonrpc'
MaxRequestBytes = 100000
Path = '/'
Port = 5002
ReadTimeoutMillis = 15000
RequestTimeoutMillis = 15000
WriteTimeoutMillis = 16000
`

	expectedWithHTTPCapabilities = `type = 'gateway'
schemaVersion = 1
name = 'Gateway1'
externalJobID = '4657f08a-e8cd-526f-9c13-66bbef7e4e03'
forwardingAllowed = false

[gatewayConfig]
[gatewayConfig.ConnectionManagerConfig]
AuthChallengeLen = 10
AuthGatewayId = 'gateway-node-1'
AuthTimestampToleranceSec = 5
HeartbeatIntervalSec = 20

[[gatewayConfig.ShardedDONs]]
DonName = 'workflow_1'
F = 3

[[gatewayConfig.ShardedDONs.Shards]]
[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0xabc'
Name = 'Node 1'

[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0xdef'
Name = 'Node 2'

[[gatewayConfig.ShardedDONs]]
DonName = 'workflow_2'
F = 0

[[gatewayConfig.ShardedDONs.Shards]]
[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0xghi'
Name = 'Node 3'

[[gatewayConfig.ShardedDONs.Shards.Nodes]]
Address = '0xjkl'
Name = 'Node 4'

[[gatewayConfig.Services]]
ServiceName = 'workflows'
DONs = ['workflow_1']

[[gatewayConfig.Services.Handlers]]
Name = 'http-capabilities'
ServiceName = 'workflows'

[gatewayConfig.Services.Handlers.Config]
CleanUpPeriodMs = 600000

[gatewayConfig.Services.Handlers.Config.NodeRateLimiter]
globalBurst = 100
globalRPS = 500
perSenderBurst = 100
perSenderRPS = 100

[[gatewayConfig.Services]]
ServiceName = 'vault'
DONs = ['workflow_2']

[[gatewayConfig.Services.Handlers]]
Name = 'vault'
ServiceName = 'vault'

[gatewayConfig.Services.Handlers.Config]
requestTimeoutSec = 14

[gatewayConfig.Services.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10

[gatewayConfig.HTTPClientConfig]
MaxResponseBytes = 50000000
AllowedPorts = [443]
AllowedSchemes = ['https']
AllowedIPsCIDR = []

[gatewayConfig.NodeServerConfig]
HandshakeTimeoutMillis = 1000
MaxRequestBytes = 100000
Path = '/'
Port = 5003
ReadTimeoutMillis = 1000
RequestTimeoutMillis = 15000
WriteTimeoutMillis = 1000

[gatewayConfig.UserServerConfig]
ContentTypeHeader = 'application/jsonrpc'
MaxRequestBytes = 100000
Path = '/'
Port = 5002
ReadTimeoutMillis = 15000
RequestTimeoutMillis = 15000
WriteTimeoutMillis = 16000
`
)

func TestGateway_Resolve(t *testing.T) {
	t.Parallel()

	g := GatewayJob{
		JobName:           "Gateway1",
		RequestTimeoutSec: 15,
		DONs: []TargetDON{
			{
				ID: "workflow_1",
				F:  1,
				Members: []TargetDONMember{
					{Address: "0xabc", Name: "Node 1"},
					{Address: "0xdef", Name: "Node 2"},
					{Address: "0xghi", Name: "Node 3"},
					{Address: "0xjkl", Name: "Node 4"},
				},
			},
			{
				ID: "workflow_2",
				Members: []TargetDONMember{
					{Address: "0x2abc", Name: "Node 1"},
					{Address: "0x2def", Name: "Node 2"},
					{Address: "0x2ghi", Name: "Node 3"},
					{Address: "0x2jkl", Name: "Node 4"},
				},
			},
		},
		Services: []GatewayServiceConfig{
			{
				ServiceName: ServiceNameWorkflows,
				Handlers:    []string{GatewayHandlerTypeWebAPICapabilities},
				DONs:        []string{"workflow_1", "workflow_2"},
			},
		},
	}

	spec, err := g.Resolve(1)
	require.NoError(t, err)
	assert.Equal(t, expected, spec)
}

func TestGateway_Resolve_WithVaultHandler(t *testing.T) {
	t.Parallel()

	g := GatewayJob{
		JobName:           "Gateway1",
		RequestTimeoutSec: 15,
		DONs: []TargetDON{
			{
				ID: "workflow_1",
				F:  1,
				Members: []TargetDONMember{
					{Address: "0xabc", Name: "Node 1"},
					{Address: "0xdef", Name: "Node 2"},
					{Address: "0xghi", Name: "Node 3"},
					{Address: "0xjkl", Name: "Node 4"},
				},
			},
			{
				ID: "workflow_2",
				Members: []TargetDONMember{
					{Address: "0x2abc", Name: "Node 1"},
					{Address: "0x2def", Name: "Node 2"},
					{Address: "0x2ghi", Name: "Node 3"},
					{Address: "0x2jkl", Name: "Node 4"},
				},
			},
		},
		Services: []GatewayServiceConfig{
			{
				ServiceName: ServiceNameWorkflows,
				Handlers:    []string{GatewayHandlerTypeWebAPICapabilities},
				DONs:        []string{"workflow_1", "workflow_2"},
			},
			{
				ServiceName: ServiceNameVault,
				Handlers:    []string{GatewayHandlerTypeVault},
				DONs:        []string{"workflow_1"},
			},
		},
	}

	spec, err := g.Resolve(1)
	require.NoError(t, err)
	assert.Equal(t, expectedWithVault, spec)
}

func TestGateway_Resolve_WithHTTPCapabilitiesHandler(t *testing.T) {
	t.Parallel()

	g := GatewayJob{
		JobName:           "Gateway1",
		RequestTimeoutSec: 15,
		DONs: []TargetDON{
			{
				ID: "workflow_1",
				F:  3,
				Members: []TargetDONMember{
					{Address: "0xabc", Name: "Node 1"},
					{Address: "0xdef", Name: "Node 2"},
				},
			},
			{
				ID: "workflow_2",
				Members: []TargetDONMember{
					{Address: "0xghi", Name: "Node 3"},
					{Address: "0xjkl", Name: "Node 4"},
				},
			},
		},
		Services: []GatewayServiceConfig{
			{
				ServiceName: ServiceNameWorkflows,
				Handlers:    []string{GatewayHandlerTypeHTTPCapabilities},
				DONs:        []string{"workflow_1"},
			},
			{
				ServiceName: ServiceNameVault,
				Handlers:    []string{GatewayHandlerTypeVault},
				DONs:        []string{"workflow_2"},
			},
		},
	}

	spec, err := g.Resolve(1)
	require.NoError(t, err)
	assert.Equal(t, expectedWithHTTPCapabilities, spec)
}
