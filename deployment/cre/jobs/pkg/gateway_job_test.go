package pkg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGateway_Validate(t *testing.T) {
	t.Parallel()

	g := GatewayJob{}
	require.ErrorContains(t, g.Validate(), "must provide job name")

	g.JobName = "AGatewayJob"
	require.ErrorContains(t, g.Validate(), "must provide at least one target DON")
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

[[gatewayConfig.Dons]]
DonId = 'workflow_1'

[[gatewayConfig.Dons.Handlers]]
Name = 'web-api-capabilities'

[gatewayConfig.Dons.Handlers.Config]
maxAllowedMessageAgeSec = 1000

[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10

[[gatewayConfig.Dons.Members]]
Address = '0xabc'
Name = 'Node 1'

[[gatewayConfig.Dons.Members]]
Address = '0xdef'
Name = 'Node 2'

[[gatewayConfig.Dons.Members]]
Address = '0xghi'
Name = 'Node 3'

[[gatewayConfig.Dons.Members]]
Address = '0xjkl'
Name = 'Node 4'

[[gatewayConfig.Dons]]
DonId = 'workflow_2'

[[gatewayConfig.Dons.Handlers]]
Name = 'web-api-capabilities'

[gatewayConfig.Dons.Handlers.Config]
maxAllowedMessageAgeSec = 1000

[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10

[[gatewayConfig.Dons.Members]]
Address = '0x2abc'
Name = 'Node 1'

[[gatewayConfig.Dons.Members]]
Address = '0x2def'
Name = 'Node 2'

[[gatewayConfig.Dons.Members]]
Address = '0x2ghi'
Name = 'Node 3'

[[gatewayConfig.Dons.Members]]
Address = '0x2jkl'
Name = 'Node 4'

[gatewayConfig.HTTPClientConfig]
MaxResponseBytes = 50000000

[gatewayConfig.NodeServerConfig]
HandshakeTimeoutMillis = 1000
MaxRequestBytes = 100000
Path = '/'
Port = 5003
ReadTimeoutMillis = 1000
RequestTimeoutMillis = 10000
WriteTimeoutMillis = 1000

[gatewayConfig.UserServerConfig]
ContentTypeHeader = 'application/jsonrpc'
MaxRequestBytes = 100000
Path = '/'
Port = 5002
ReadTimeoutMillis = 80000
RequestTimeoutMillis = 80000
WriteTimeoutMillis = 80000
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

[[gatewayConfig.Dons]]
DonId = 'workflow_1'

[[gatewayConfig.Dons.Handlers]]
Name = 'web-api-capabilities'

[gatewayConfig.Dons.Handlers.Config]
maxAllowedMessageAgeSec = 1000

[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10

[[gatewayConfig.Dons.Handlers]]
Name = 'vault'
ServiceName = 'vault'

[gatewayConfig.Dons.Handlers.Config]
requestTimeoutSec = 70

[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10

[[gatewayConfig.Dons.Members]]
Address = '0xabc'
Name = 'Node 1'

[[gatewayConfig.Dons.Members]]
Address = '0xdef'
Name = 'Node 2'

[[gatewayConfig.Dons.Members]]
Address = '0xghi'
Name = 'Node 3'

[[gatewayConfig.Dons.Members]]
Address = '0xjkl'
Name = 'Node 4'

[[gatewayConfig.Dons]]
DonId = 'workflow_2'

[[gatewayConfig.Dons.Handlers]]
Name = 'web-api-capabilities'

[gatewayConfig.Dons.Handlers.Config]
maxAllowedMessageAgeSec = 1000

[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10

[[gatewayConfig.Dons.Members]]
Address = '0x2abc'
Name = 'Node 1'

[[gatewayConfig.Dons.Members]]
Address = '0x2def'
Name = 'Node 2'

[[gatewayConfig.Dons.Members]]
Address = '0x2ghi'
Name = 'Node 3'

[[gatewayConfig.Dons.Members]]
Address = '0x2jkl'
Name = 'Node 4'

[gatewayConfig.HTTPClientConfig]
MaxResponseBytes = 50000000

[gatewayConfig.NodeServerConfig]
HandshakeTimeoutMillis = 1000
MaxRequestBytes = 100000
Path = '/'
Port = 5003
ReadTimeoutMillis = 1000
RequestTimeoutMillis = 10000
WriteTimeoutMillis = 1000

[gatewayConfig.UserServerConfig]
ContentTypeHeader = 'application/jsonrpc'
MaxRequestBytes = 100000
Path = '/'
Port = 5002
ReadTimeoutMillis = 80000
RequestTimeoutMillis = 80000
WriteTimeoutMillis = 80000
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

[[gatewayConfig.Dons]]
DonId = 'workflow_1'

[[gatewayConfig.Dons.Handlers]]
Name = 'http-capabilities'
ServiceName = 'workflows'

[gatewayConfig.Dons.Handlers.Config]
CleanUpPeriodMs = 86400000

[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10

[[gatewayConfig.Dons.Members]]
Address = '0xabc'
Name = 'Node 1'

[[gatewayConfig.Dons.Members]]
Address = '0xdef'
Name = 'Node 2'

[[gatewayConfig.Dons]]
DonId = 'workflow_2'

[[gatewayConfig.Dons.Handlers]]
Name = 'vault'
ServiceName = 'vault'

[gatewayConfig.Dons.Handlers.Config]
requestTimeoutSec = 70

[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10

[[gatewayConfig.Dons.Members]]
Address = '0xghi'
Name = 'Node 3'

[[gatewayConfig.Dons.Members]]
Address = '0xjkl'
Name = 'Node 4'

[gatewayConfig.HTTPClientConfig]
MaxResponseBytes = 50000000

[gatewayConfig.NodeServerConfig]
HandshakeTimeoutMillis = 1000
MaxRequestBytes = 100000
Path = '/'
Port = 5003
ReadTimeoutMillis = 1000
RequestTimeoutMillis = 10000
WriteTimeoutMillis = 1000

[gatewayConfig.UserServerConfig]
ContentTypeHeader = 'application/jsonrpc'
MaxRequestBytes = 100000
Path = '/'
Port = 5002
ReadTimeoutMillis = 80000
RequestTimeoutMillis = 80000
WriteTimeoutMillis = 80000
`
)

func TestGateway_Resolve(t *testing.T) {
	t.Parallel()

	g := GatewayJob{
		JobName: "Gateway1",
		TargetDONs: []TargetDON{
			{
				ID:       "workflow_1",
				Handlers: []string{GatewayHandlerTypeWebAPICapabilities},
				Members: []TargetDONMember{
					{
						Address: "0xabc",
						Name:    "Node 1",
					},
					{
						Address: "0xdef",
						Name:    "Node 2",
					},
					{
						Address: "0xghi",
						Name:    "Node 3",
					},
					{
						Address: "0xjkl",
						Name:    "Node 4",
					},
				},
			},
			{
				ID:       "workflow_2",
				Handlers: []string{GatewayHandlerTypeWebAPICapabilities},
				Members: []TargetDONMember{
					{
						Address: "0x2abc",
						Name:    "Node 1",
					},
					{
						Address: "0x2def",
						Name:    "Node 2",
					},
					{
						Address: "0x2ghi",
						Name:    "Node 3",
					},
					{
						Address: "0x2jkl",
						Name:    "Node 4",
					},
				},
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
		JobName: "Gateway1",
		TargetDONs: []TargetDON{
			{
				ID:       "workflow_1",
				Handlers: []string{GatewayHandlerTypeWebAPICapabilities, GatewayHandlerTypeVault},
				Members: []TargetDONMember{
					{
						Address: "0xabc",
						Name:    "Node 1",
					},
					{
						Address: "0xdef",
						Name:    "Node 2",
					},
					{
						Address: "0xghi",
						Name:    "Node 3",
					},
					{
						Address: "0xjkl",
						Name:    "Node 4",
					},
				},
			},
			{
				ID:       "workflow_2",
				Handlers: []string{GatewayHandlerTypeWebAPICapabilities},
				Members: []TargetDONMember{
					{
						Address: "0x2abc",
						Name:    "Node 1",
					},
					{
						Address: "0x2def",
						Name:    "Node 2",
					},
					{
						Address: "0x2ghi",
						Name:    "Node 3",
					},
					{
						Address: "0x2jkl",
						Name:    "Node 4",
					},
				},
			},
		},
	}

	spec, err := g.Resolve(1)
	fmt.Println(spec)
	require.NoError(t, err)
	assert.Equal(t, expectedWithVault, spec)
}

func TestGateway_Resolve_WithHTTPCapabilitiesHandler(t *testing.T) {
	t.Parallel()

	g := GatewayJob{
		JobName: "Gateway1",
		TargetDONs: []TargetDON{
			{
				ID:       "workflow_1",
				Handlers: []string{GatewayHandlerTypeHTTPCapabilities},
				Members: []TargetDONMember{
					{
						Address: "0xabc",
						Name:    "Node 1",
					},
					{
						Address: "0xdef",
						Name:    "Node 2",
					},
				},
			},
			{
				ID:       "workflow_2",
				Handlers: []string{GatewayHandlerTypeVault},
				Members: []TargetDONMember{
					{
						Address: "0xghi",
						Name:    "Node 3",
					},
					{
						Address: "0xjkl",
						Name:    "Node 4",
					},
				},
			},
		},
	}

	spec, err := g.Resolve(1)
	require.NoError(t, err)
	assert.Equal(t, expectedWithHTTPCapabilities, spec)
}
