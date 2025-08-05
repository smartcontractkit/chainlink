package jobs

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

var (
	DefaultAllowedPorts = []int{80, 443}
)

type HandlerType string

const (
	WebAPIHandlerType HandlerType = "web-api-capabilities"
	HTTPHandlerType   HandlerType = "http-capabilities"
)

func BootstrapOCR3(nodeID string, name string, ocr3CapabilityAddress string, chainID uint64) *jobv1.ProposeJobRequest {
	uuid := uuid.NewString()

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
	type = "bootstrap"
	schemaVersion = 1
	externalJobID = "%s"
	name = "%s"
	contractID = "%s"
	contractConfigTrackerPollInterval = "1s"
	contractConfigConfirmations = 1
	relay = "evm"
	[relayConfig]
	chainID = %d
	providerType = "ocr3-capability"
`,
			uuid,
			"ocr3-bootstrap-"+name,
			ocr3CapabilityAddress,
			chainID),
	}
}

func AnyGateway(handlerType HandlerType, bootstrapNodeID string, chainID uint64, donID uint32, extraAllowedPorts []int, extraAllowedIps, extrAallowedIPsCIDR []string, gatewayConnectorData cre.GatewayConnectorOutput) *jobv1.ProposeJobRequest {
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

		if handlerType == HTTPHandlerType {
			gatewayDons += fmt.Sprintf(`
				[[gatewayConfig.Dons]]
				DonId = "workflows"
				F = 1
				HandlerName = "http-capabilities"
					[gatewayConfig.Dons.HandlerConfig]
					MaxAllowedMessageAgeSec = 1_000
					authPullIntervalMs = 1000
						[gatewayConfig.Dons.HandlerConfig.NodeRateLimiter]
						GlobalBurst = 10
						GlobalRPS = 50
						PerSenderBurst = 10
						PerSenderRPS = 10
						[gatewayConfig.Dons.HandlerConfig.UserRateLimiter]
						GlobalBurst = 10
						GlobalRPS = 50
						PerSenderBurst = 10
						PerSenderRPS = 10
					%s
			`, gatewayMembers)
		} else {
			gatewayDons += fmt.Sprintf(`
				[[gatewayConfig.Dons]]
				DonId = "%s"
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
	}

	uuid := uuid.NewString()

	gatewayJobSpec := fmt.Sprintf(`
	type = "gateway"
	schemaVersion = 1
	externalJobID = "%s"
	name = "%s"
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
	Path = "%s"
	Port = %d
	ReadTimeoutMillis = 1_000
	RequestTimeoutMillis = 10_000
	WriteTimeoutMillis = 1_000
	CORSEnabled = false
	CORSAllowedOrigins = []
	[gatewayConfig.HTTPClientConfig]
	MaxResponseBytes = 100_000_000
`,
		uuid,
		cre.GatewayJobName,
		gatewayDons,
		gatewayConnectorData.Outgoing.Path,
		gatewayConnectorData.Outgoing.Port,
		gatewayConnectorData.Incoming.Path,
		gatewayConnectorData.Incoming.InternalPort,
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

	if len(extrAallowedIPsCIDR) != 0 {
		allowedIPsCIDR := strings.Join(extrAallowedIPsCIDR, `", "`)

		gatewayJobSpec += fmt.Sprintf(`
	AllowedIPsCIDR = ["%s"]
`,
			allowedIPsCIDR,
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

func WorkerStandardCapability(nodeID, name, command, config, oracleFactoryConfig string) *jobv1.ProposeJobRequest {
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
	%s
`,
			uuid.NewString(),
			name,
			command,
			config,
			oracleFactoryConfig),
	}
}

func WorkerOCR3(nodeID string, ocr3CapabilityAddress, nodeEthAddress, ocr2KeyBundleID string, ocrPeeringData cre.OCRPeeringData, chainID uint64) *jobv1.ProposeJobRequest {
	uuid := uuid.NewString()

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
	type = "offchainreporting2"
	schemaVersion = 1
	externalJobID = "%s"
	name = "%s"
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
			cre.OCR3Capability,
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

func WorkerVaultOCR3(nodeID string, vaultCapabilityAddress, nodeEthAddress, ocr2KeyBundleID string, ocrPeeringData cre.OCRPeeringData, chainID uint64, masterPublicKey string, encryptedPrivateKeyShare string) *jobv1.ProposeJobRequest {
	uuid := uuid.NewString()

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
	type = "offchainreporting2"
	schemaVersion = 1
	externalJobID = "%s"
	name = "%s"
	contractID = "%s"
	ocrKeyBundleID = "%s"
	p2pv2Bootstrappers = [
		"%s@%s",
	]
	relay = "evm"
	pluginType = "%s"
	transmitterID = "%s"
	[relayConfig]
	chainID = "%d"
	[pluginConfig]
	requestExpiryDuration = "60s"
	[pluginConfig.dkg]
	masterPublicKey = "%s"
	encryptedPrivateKeyShare = "%s"
`,
			uuid,
			"Vault OCR3 Capability",
			vaultCapabilityAddress,
			ocr2KeyBundleID,
			ocrPeeringData.OCRBootstraperPeerID,
			fmt.Sprintf("%s:%d", ocrPeeringData.OCRBootstraperHost, ocrPeeringData.Port),
			types.VaultPlugin,
			nodeEthAddress,
			chainID,
			masterPublicKey,
			encryptedPrivateKeyShare,
		),
	}
}
