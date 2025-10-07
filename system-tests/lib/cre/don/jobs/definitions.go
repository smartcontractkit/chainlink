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
			"ocr3-bootstrap-"+name+fmt.Sprintf("-%d", chainID),
			ocr3CapabilityAddress,
			chainID),
	}
}

type GatewayHandler struct {
	Name   string
	Config string
}

func AnyGateway(bootstrapNodeID string, chainID uint64, extraAllowedPorts []int, extraAllowedIps, extrAallowedIPsCIDR []string, gatewayConfiguration *cre.DonGatewayConfiguration) *jobv1.ProposeJobRequest {
	var gatewayDons string

	for _, don := range gatewayConfiguration.Dons {
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

		var handlersConfig string
		for name, config := range don.Handlers {
			handlersConfig += fmt.Sprintf(`
	[[gatewayConfig.Dons.Handlers]]
	Name = "%s"
	%s
		`, name, config)
		}

		gatewayDons += fmt.Sprintf(`
	[[gatewayConfig.Dons]]
	DonId = "%s"
	F = 1
	%s
	%s
		`, don.ID, gatewayMembers, handlersConfig)
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
	AuthGatewayId = "%s"
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
	ReadTimeoutMillis = 80_000
	RequestTimeoutMillis = 80_000
	WriteTimeoutMillis = 80_000
	CORSEnabled = false
	CORSAllowedOrigins = []
	[gatewayConfig.HTTPClientConfig]
	MaxResponseBytes = 100_000_000
`,
		uuid,
		"cre-gateway",
		gatewayConfiguration.AuthGatewayID,
		gatewayDons,
		gatewayConfiguration.Outgoing.Path,
		gatewayConfiguration.Outgoing.Port,
		gatewayConfiguration.Incoming.Path,
		gatewayConfiguration.Incoming.InternalPort,
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

func DonTimeJob(nodeID string, ocr3CapabilityAddress, nodeEthAddress, ocr2KeyBundleID string, ocrPeeringData cre.OCRPeeringData, chainID uint64) *jobv1.ProposeJobRequest {
	uuid := uuid.NewString()
	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
	type = "offchainreporting2"
	schemaVersion = 1
	externalJobID = "%s"
	name = "dontime"
	forwardingAllowed = false
	maxTaskDuration = "0s"
	contractID = "%s"
	relay = "evm"
	pluginType = "dontime"
	ocrKeyBundleID = "%s"
	p2pv2Bootstrappers = [
		"%s@%s",
	]
	transmitterID = "%s"

	[relayConfig]
	chainID = "%d"
	providerType = "dontime"

	[pluginConfig]
	pluginName = "dontime"
	ocrVersion = 3
	telemetryType = "plugin"

	[onchainSigningStrategy]
	strategyName = 'multi-chain'
	[onchainSigningStrategy.config]
	evm = "%s"
`,
			uuid,
			ocr3CapabilityAddress, // re-use OCR3Capability contract
			ocr2KeyBundleID,
			ocrPeeringData.OCRBootstraperPeerID,
			fmt.Sprintf("%s:%d", ocrPeeringData.OCRBootstraperHost, ocrPeeringData.Port),
			nodeEthAddress, // transmitterID (although this shouldn't be used for this plugin?)
			chainID,
			ocr2KeyBundleID,
		),
	}
}

func WorkerOCR3(nodeID string, ocr3CapabilityAddress, nodeEthAddress, offchainBundleID string, ocr2KeyBundles map[string]string, ocrPeeringData cre.OCRPeeringData, chainID uint64) *jobv1.ProposeJobRequest {
	uuid := uuid.NewString()

	spec := fmt.Sprintf(`
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
	strategyName = "multi-chain"
	[onchainSigningStrategy.config]
`,
		uuid,
		cre.ConsensusCapability,
		ocr3CapabilityAddress,
		offchainBundleID,
		ocrPeeringData.OCRBootstraperPeerID,
		fmt.Sprintf("%s:%d", ocrPeeringData.OCRBootstraperHost, ocrPeeringData.Port),
		nodeEthAddress,
		chainID,
	)
	for family, key := range ocr2KeyBundles {
		spec += fmt.Sprintf(`
        %s = "%s"`, family, key)
		spec += "\n"
	}

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec:   spec,
	}
}

func WorkerVaultOCR3(nodeID string, vaultCapabilityAddress, dkgAddress, nodeEthAddress, ocr2KeyBundleID string, ocrPeeringData cre.OCRPeeringData, chainID uint64) *jobv1.ProposeJobRequest {
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
	dkgContractID = "%s"
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
			dkgAddress,
		),
	}
}
