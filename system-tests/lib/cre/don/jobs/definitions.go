package jobs

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

var (
	DefaultAllowedPorts = []int{80, 443}
)

func BootstrapOCR3(nodeID string, ocr3CapabilityAddress common.Address, chainID uint64) *jobv1.ProposeJobRequest {
	uuid := uuid.NewString()

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
	type = "bootstrap"
	schemaVersion = 1
	externalJobID = "%s"
	name = "Botostrap-%s"
	contractID = "%s"
	contractConfigTrackerPollInterval = "1s"
	contractConfigConfirmations = 1
	relay = "evm"
	[relayConfig]
	chainID = %d
	providerType = "ocr3-capability"
`,
			uuid,
			uuid[0:8],
			ocr3CapabilityAddress.Hex(),
			chainID),
	}
}

func AnyGateway(bootstrapNodeID string, chainID uint64, donID uint32, extraAllowedPorts []int, extraAllowedIps []string, gatewayConnectorData types.GatewayConnectorOutput) *jobv1.ProposeJobRequest {
	var gatewayDons string

	for _, don := range gatewayConnectorData.Dons {
		var gatewayMembers string

		for i := 0; i < len(don.MembersEthAddresses); i++ {
			gatewayMembers += fmt.Sprintf(`
	[[gatewayConfig.Dons.Members]]
	Address = "%s"
	Name = "Node %d"`,
				don.MembersEthAddresses[i],
				i+1,
			)
		}

		gatewayDons += fmt.Sprintf(`
		[[gatewayConfig.Dons]]
		DonId = "%d"
		F = 1
		HandlerName = "web-api-capabilities"
			[gatewayConfig.Dons.HandlerConfig]
			MaxAllowedMessageAgeSec = 1_000
				[gatewayConfig.Dons.HandlerConfig.NodeRateLimiter]
				GlobalBurst = 10
				GlobalRPS = 50
				PerSenderBurst = 10
				PerSenderRPS = 10
			%s
		`, don.ID, gatewayMembers)
	}

	uuid := uuid.NewString()

	gatewayJobSpec := fmt.Sprintf(`
	type = "gateway"
	schemaVersion = 1
	externalJobID = "%s"
	name = "Gateway-%s"
	forwardingAllowed = false
	[gatewayConfig.ConnectionManagerConfig]
	AuthChallengeLen = 10
	AuthGatewayId = "por_gateway"
	AuthTimestampToleranceSec = 5
	HeartbeatIntervalSec = 20
	%s
	[gatewayConfig.NodeServerConfig]
	HandshakeTimeoutMillis = 1_000
	MaxRequestBytes = 100_000
	# this is the path other nodes will use to connect to the gateway
	Path = "%s"
	# this is the port other nodes will use to connect to the gateway
	Port = %d
	ReadTimeoutMillis = 1_000
	RequestTimeoutMillis = 10_000
	WriteTimeoutMillis = 1_000
	[gatewayConfig.UserServerConfig]
	ContentTypeHeader = "application/jsonrpc"
	MaxRequestBytes = 100_000
	Path = "/"
	Port = 5_002
	ReadTimeoutMillis = 1_000
	RequestTimeoutMillis = 10_000
	WriteTimeoutMillis = 1_000
	[gatewayConfig.HTTPClientConfig]
	MaxResponseBytes = 100_000_000
`,
		uuid,
		uuid[0:8],
		gatewayDons,
		gatewayConnectorData.Path,
		gatewayConnectorData.Port,
	)

	if len(extraAllowedPorts) != 0 {
		var allowedPorts string
		allPorts := make([]int, 0, len(DefaultAllowedPorts)+len(extraAllowedPorts))
		allPorts = append(allPorts, append(extraAllowedPorts, DefaultAllowedPorts...)...)
		for _, port := range allPorts {
			allowedPorts += fmt.Sprintf("%d, ", port)
		}

		// when we pass custom allowed IPs, defaults are not used and we need to
		// pass HTTP and HTTPS explicitly
		gatewayJobSpec += fmt.Sprintf(`
	AllowedPorts = [%s]
`,
			allowedPorts,
		)
	}

	if len(extraAllowedIps) != 0 {
		allowedIPs := strings.Join(extraAllowedIps, `", "`)

		gatewayJobSpec += fmt.Sprintf(`
	AllowedIps = ["%s"]
`,
			allowedIPs,
		)
	}

	return &jobv1.ProposeJobRequest{
		NodeId: bootstrapNodeID,
		Spec:   gatewayJobSpec,
	}
}

const (
	EmptyStdCapConfig = "\"\""
)

func ExternalCapabilityPath(binaryName string) string {
	return "/home/capabilities/" + binaryName
}

func WorkerStandardCapability(nodeID, name, command, config string) *jobv1.ProposeJobRequest {
	uuid := uuid.NewString()

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
	type = "standardcapabilities"
	schemaVersion = 1
	externalJobID = "%s"
	name = "%s"
	forwardingAllowed = false
	command = "%s"
	config = %s
`,
			uuid,
			name+"-"+uuid[0:8],
			command,
			config),
	}
}

func WorkerOCR3(nodeID string, ocr3CapabilityAddress common.Address, nodeEthAddress, ocr2KeyBundleID string, ocrPeeringData types.OCRPeeringData, chainID uint64) *jobv1.ProposeJobRequest {
	uuid := uuid.NewString()

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
	type = "offchainreporting2"
	schemaVersion = 1
	externalJobID = "%s"
	name = "ocr3-consensus-%s"
	contractID = "%s"
	ocrKeyBundleID = "%s"
	p2pv2Bootstrappers = [
		"%s@%s",
	]
	relay = "evm"
	pluginType = "plugin"
	transmitterID = "%s"
	[relayConfig]
	chainID = "%d"
	[pluginConfig]
	command = "/usr/local/bin/chainlink-ocr3-capability"
	ocrVersion = 3
	pluginName = "ocr-capability"
	providerType = "ocr3-capability"
	telemetryType = "plugin"
	[onchainSigningStrategy]
	strategyName = 'multi-chain'
	[onchainSigningStrategy.config]
	evm = "%s"
`,
			uuid,
			uuid[0:8],
			ocr3CapabilityAddress,
			ocr2KeyBundleID,
			ocrPeeringData.OCRBootstraperPeerID,
			fmt.Sprintf("%s:%d", ocrPeeringData.OCRBootstraperHost, ocrPeeringData.Port),
			nodeEthAddress,
			chainID,
			ocr2KeyBundleID,
		),
	}
}

func MockCapabilities(nodeID string) *jobv1.ProposeJobRequest {
	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
			type = "standardcapabilities"
			schemaVersion = 1
			externalJobID = "%s"
			name = "mock-capabilitie"
			forwardingAllowed = false
			command = "/home/capabilities/amd64_mock"
			config = """
				port=7777
				[[DefaultMocks]]
				id="streams-trigger@1.0.0"
				description="stream trigger mock"
				type="trigger"
				[[DefaultMocks]]
				id="write_ethereum@1.0.0"
				description="write trigger mock"
				type="target"
"""
`,
			uuid.NewString(),
		),
	}
}

func TextWorkflow(nodeID string) *jobv1.ProposeJobRequest {
	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: `
type = "workflow"
 schemaVersion = 1
 name = "test-workflow-with-mock"
 externalJobID = "65a75326-65f3-4a40-b8f5-407d2376ce7d"
 forwardingAllowed = false
 workflowID = "2ab944834c5b3e58aa71e4a4a29d6382b613f833302550277cbc3f5cc2be5226"
 workflow = """
name: abcdefgasd
owner: '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512'
triggers:
  - id: streams-trigger@1.0.0
    config:
      maxFrequencyMs: 5000
      feedIds:
        - '0x000351de403f638036014add21a5abd5f464bf21d11aa356dfc6dbe4e2384e4e' # BTC/USD
        - '0x0003f2f4cae1891f647db8d73c87a7a03888bd176afdb7206853da9abfc92874' # ETH/USD
        - '0x00034db6355441c80b613f666757c63777dae7743885a9c594ca25d9f9b896ca' # LINK/USD

consensus:
  - id: offchain_reporting@1.0.0
    ref: ccip_feeds
    inputs:
      observations:
        - $(trigger.outputs)
    config:
      report_id: '0001'
      key_id: 'evm'
      aggregation_method: data_feeds
      aggregation_config:
        allowedPartialStaleness: '0.5'
        feeds:
          '0x000351de403f638036014add21a5abd5f464bf21d11aa356dfc6dbe4e2384e4e':  # BTC/USD
            deviation: '0.01'
            heartbeat: 600
            remappedID: '0x666666666666'
          '0x0003f2f4cae1891f647db8d73c87a7a03888bd176afdb7206853da9abfc92874': # ETH/USD
            deviation: '0.01'
            heartbeat: 600
            remappedID: '0x777777777777'
          '0x00034db6355441c80b613f666757c63777dae7743885a9c594ca25d9f9b896ca': # LINK/USD
            deviation: '0.01'
            heartbeat: 600
      encoder: EVM
      encoder_config:
        abi: (bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports

targets:
  - id: write_ethereum@1.0.0
    inputs:
      signed_report: $(ccip_feeds.outputs)
    config:
      address: '0x24309990d635A6C5FF711503BfCb942dd25F96A0'
      deltaStage: 10s
      schedule: oneAtATime
"""
`,
	}
}
