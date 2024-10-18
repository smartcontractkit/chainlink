package framework

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

type DonContext struct {
	EthBlockchain      *EthBlockchain
	p2pNetwork         *MockRageP2PNetwork
	capabilityRegistry *CapabilitiesRegistry
}

func CreateDonContext(ctx context.Context, t *testing.T) DonContext {
	ethBlockchain := NewEthBlockchain(t, 1000, 1*time.Second)
	rageP2PNetwork := NewMockRageP2PNetwork(t, 1000)
	capabilitiesRegistry := NewCapabilitiesRegistry(ctx, t, ethBlockchain)

	servicetest.Run(t, rageP2PNetwork)
	servicetest.Run(t, ethBlockchain)
	return DonContext{EthBlockchain: ethBlockchain, p2pNetwork: rageP2PNetwork, capabilityRegistry: capabilitiesRegistry}
}

type capabilityNode struct {
	*cltest.TestApplication
	registry  *capabilities.Registry
	key       ethkey.KeyV2
	KeyBundle ocr2key.KeyBundle
	peerID    peer
	start     func()
}

type DON struct {
	t                      *testing.T
	config                 DonConfiguration
	lggr                   logger.Logger
	nodes                  []*capabilityNode
	standardCapabilityJobs []*job.Job
	externalCapabilities   []capability
	capabilitiesRegistry   *CapabilitiesRegistry

	nodeConfigModifiers []func(c *chainlink.Config, node *capabilityNode)

	addOCR3NonStandardCapability bool

	triggerFactories []TriggerFactory
	targetFactories  []TargetFactory
}

func NewDON(ctx context.Context, t *testing.T, lggr logger.Logger, donConfig DonConfiguration,
	dependentDONs []commoncap.DON, donContext DonContext, supportsOCR bool) *DON {
	don := &DON{t: t, lggr: lggr.Named(donConfig.name), config: donConfig, capabilitiesRegistry: donContext.capabilityRegistry}

	var newOracleFactoryFn standardcapabilities.NewOracleFactoryFn
	var libOcr *MockLibOCR
	if supportsOCR {
		libOcr = NewMockLibOCR(t, lggr, donConfig.F, 1*time.Second)
		servicetest.Run(t, libOcr)
	}

	for i, member := range donConfig.Members {
		dispatcher := donContext.p2pNetwork.NewDispatcherForNode(member)
		capabilityRegistry := capabilities.NewRegistry(lggr)

		nodeInfo := commoncap.Node{
			PeerID:         &member,
			WorkflowDON:    donConfig.DON,
			CapabilityDONs: dependentDONs,
		}

		cn := &capabilityNode{
			registry:  capabilityRegistry,
			key:       donConfig.keys[i],
			KeyBundle: donConfig.KeyBundles[i],
			peerID:    donConfig.peerIDs[i],
		}
		don.nodes = append(don.nodes, cn)

		if supportsOCR {
			factory := newMockLibOcrOracleFactory(libOcr, donConfig.KeyBundles[i], len(donConfig.Members), int(donConfig.F))
			newOracleFactoryFn = factory.NewOracleFactory
		}

		cn.start = func() {
			node := startNewNode(ctx, t, lggr.Named(donConfig.name+"-"+strconv.Itoa(i)), nodeInfo, donContext.EthBlockchain,
				donContext.capabilityRegistry.getAddress(), dispatcher,
				peerWrapper{peer: p2pPeer{member}}, capabilityRegistry, newOracleFactoryFn,
				donConfig.keys[i], func(c *chainlink.Config) {
					for _, modifier := range don.nodeConfigModifiers {
						modifier(c, cn)
					}
				})

			require.NoError(t, node.Start(testutils.Context(t)))
			cn.TestApplication = node
		}
	}

	return don
}

// Initialise must be called after all capabilities have been added to the DONs and before Start is called
func (d *DON) Initialise() {
	if len(d.externalCapabilities) > 0 {
		id := d.capabilitiesRegistry.setupDON(d.config, d.externalCapabilities)

		//nolint:gosec // disable G115
		d.config.DON.ID = uint32(id)
	}
}

func (d *DON) GetID() uint32 {
	if d.config.DON.ID == 0 {
		panic("DON ID not set, call Initialise() first")
	}

	return d.config.ID
}

func (d *DON) GetConfigVersion() uint32 {
	return d.config.ConfigVersion
}

func (d *DON) GetF() uint8 {
	return d.config.F
}

func (d *DON) GetPeerIDs() []peer {
	return d.config.peerIDs
}

func (d *DON) Start(ctx context.Context, t *testing.T) {
	for _, triggerFactory := range d.triggerFactories {
		for _, node := range d.nodes {
			trigger := triggerFactory.CreateNewTrigger(t)
			err := node.registry.Add(ctx, trigger)
			require.NoError(t, err)
		}
	}

	for _, targetFactory := range d.targetFactories {
		for _, node := range d.nodes {
			target := targetFactory.CreateNewTarget(t)
			err := node.registry.Add(ctx, target)
			require.NoError(t, err)
		}
	}

	for _, node := range d.nodes {
		node.start()
	}

	if d.addOCR3NonStandardCapability {
		libocr := NewMockLibOCR(t, d.lggr, d.config.F, 1*time.Second)
		servicetest.Run(t, libocr)

		for _, node := range d.nodes {
			addOCR3Capability(ctx, t, d.lggr, node.registry, libocr, d.config.F, node.KeyBundle)
		}
	}

	for _, capabilityJob := range d.standardCapabilityJobs {
		err := d.AddJob(ctx, capabilityJob)
		require.NoError(t, err)
	}
}

const StandardCapabilityTemplateJobSpec = `
type = "standardcapabilities"
schemaVersion = 1
name = "%s"
command="%s"
config="%s"
`

func (d *DON) AddStandardCapability(name string, command string, config string) {
	spec := fmt.Sprintf(StandardCapabilityTemplateJobSpec, name, command, config)
	capabilitiesSpecJob, err := standardcapabilities.ValidatedStandardCapabilitiesSpec(spec)
	require.NoError(d.t, err)

	d.standardCapabilityJobs = append(d.standardCapabilityJobs, &capabilitiesSpecJob)
}

// TODO - add configuration for remote support - do this for each capability as an option
func (d *DON) AddTargetCapability(targetFactory TargetFactory) {
	d.targetFactories = append(d.targetFactories, targetFactory)
}

func (d *DON) AddExternalTriggerCapability(triggerFactory TriggerFactory) {
	d.triggerFactories = append(d.triggerFactories, triggerFactory)

	// Arguably this should be a parameter to AddExternalTriggerCapability, but for now we're just using the default
	// See TODO about local/remote exposure
	defaultTriggerCapabilityConfig := newCapabilityConfig()
	defaultTriggerCapabilityConfig.RemoteConfig = &pb.CapabilityConfig_RemoteTriggerConfig{
		RemoteTriggerConfig: &pb.RemoteTriggerConfig{
			RegistrationRefresh: durationpb.New(1000 * time.Millisecond),
			RegistrationExpiry:  durationpb.New(60000 * time.Millisecond),
			// F + 1
			MinResponsesToAggregate: uint32(d.config.F) + 1,
		},
	}

	triggerCapability := capability{
		donCapabilityConfig: defaultTriggerCapabilityConfig,
		registryConfig: kcr.CapabilitiesRegistryCapability{
			LabelledName:   triggerFactory.GetTriggerName(),
			Version:        triggerFactory.GetTriggerVersion(),
			CapabilityType: uint8(registrysyncer.ContractCapabilityTypeTrigger),
		},
	}

	d.externalCapabilities = append(d.externalCapabilities, triggerCapability)
}

func (d *DON) AddJob(ctx context.Context, j *job.Job) error {
	for _, node := range d.nodes {
		err := node.AddJobV2(ctx, j)
		if err != nil {
			return fmt.Errorf("failed to add job: %w", err)
		}
	}

	return nil
}

type TriggerFactory interface {
	CreateNewTrigger(t *testing.T) commoncap.TriggerCapability
	GetTriggerID() string
	GetTriggerName() string
	GetTriggerVersion() string
}

type TargetFactory interface {
	CreateNewTarget(t *testing.T) commoncap.TargetCapability
	GetTargetID() string
	GetTargetName() string
	GetTargetVersion() string
}

func startNewNode(ctx context.Context,
	t *testing.T, lggr logger.Logger, nodeInfo commoncap.Node,
	ethBlockchain *EthBlockchain, capRegistryAddr common.Address,
	dispatcher remotetypes.Dispatcher,
	peerWrapper p2ptypes.PeerWrapper,
	localCapabilities *capabilities.Registry,
	newOracleFactoryFn standardcapabilities.NewOracleFactoryFn,
	keyV2 ethkey.KeyV2,
	setupCfg func(c *chainlink.Config),
) *cltest.TestApplication {
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Capabilities.ExternalRegistry.ChainID = ptr(fmt.Sprintf("%d", testutils.SimulatedChainID))
		c.Capabilities.ExternalRegistry.Address = ptr(capRegistryAddr.String())
		c.Capabilities.Peering.V2.Enabled = ptr(true)
		c.Feature.FeedsManager = ptr(false)

		if setupCfg != nil {
			setupCfg(c)
		}
	})

	n, err := ethBlockchain.NonceAt(ctx, ethBlockchain.transactionOpts.From, nil)
	require.NoError(t, err)

	tx := cltest.NewLegacyTransaction(
		n, keyV2.Address,
		assets.Ether(1).ToInt(),
		21000,
		assets.GWei(1).ToInt(),
		nil)
	signedTx, err := ethBlockchain.transactionOpts.Signer(ethBlockchain.transactionOpts.From, tx)
	require.NoError(t, err)
	err = ethBlockchain.SendTransaction(ctx, signedTx)
	require.NoError(t, err)
	ethBlockchain.Commit()

	return cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, ethBlockchain.SimulatedBackend, nodeInfo,
		dispatcher, peerWrapper, newOracleFactoryFn, localCapabilities, keyV2, lggr)
}

// Functions below this point are for adding non-standard capabilities to a DON, deliberately verbose. Eventually these
// should be replaced with standard capabilities.

func (d *DON) AddOCR3NonStandardCapability() {
	d.addOCR3NonStandardCapability = true

	ocr := kcr.CapabilitiesRegistryCapability{
		LabelledName:   "offchain_reporting",
		Version:        "1.0.0",
		CapabilityType: uint8(registrysyncer.ContractCapabilityTypeConsensus),
	}

	d.externalCapabilities = append(d.externalCapabilities, capability{
		donCapabilityConfig: newCapabilityConfig(),
		registryConfig:      ocr,
	})
}

func (d *DON) AddEthereumWriteTargetNonStandardCapability(forwarderAddr common.Address) error {
	d.nodeConfigModifiers = append(d.nodeConfigModifiers, func(c *chainlink.Config, node *capabilityNode) {
		eip55Address := types.EIP55AddressFromAddress(forwarderAddr)
		c.EVM[0].Chain.Workflow.ForwarderAddress = &eip55Address
		c.EVM[0].Chain.Workflow.FromAddress = &node.key.EIP55Address
	})

	writeChain := kcr.CapabilitiesRegistryCapability{
		LabelledName:   "write_geth-testnet",
		Version:        "1.0.0",
		CapabilityType: uint8(registrysyncer.ContractCapabilityTypeTarget),
	}

	targetCapabilityConfig := newCapabilityConfig()

	configWithLimit, err := values.WrapMap(map[string]any{"gasLimit": 500000})
	if err != nil {
		return fmt.Errorf("failed to wrap map: %w", err)
	}

	targetCapabilityConfig.DefaultConfig = values.Proto(configWithLimit).GetMapValue()

	targetCapabilityConfig.RemoteConfig = &pb.CapabilityConfig_RemoteTargetConfig{
		RemoteTargetConfig: &pb.RemoteTargetConfig{
			RequestHashExcludedAttributes: []string{"signed_report.Signatures"},
		},
	}

	d.externalCapabilities = append(d.externalCapabilities, capability{
		donCapabilityConfig: targetCapabilityConfig,
		registryConfig:      writeChain,
	})

	return nil
}

func addOCR3Capability(ctx context.Context, t *testing.T, lggr logger.Logger, capabilityRegistry *capabilities.Registry,
	libocr *MockLibOCR, donF uint8, ocr2KeyBundle ocr2key.KeyBundle) {
	requestTimeout := 10 * time.Minute
	cfg := ocr3.Config{
		Logger:            lggr,
		EncoderFactory:    capabilities.NewEncoder,
		AggregatorFactory: capabilities.NewAggregator,
		RequestTimeout:    &requestTimeout,
	}

	ocr3Capability := ocr3.NewOCR3(cfg)
	servicetest.Run(t, ocr3Capability)

	pluginCfg := coretypes.ReportingPluginServiceConfig{}
	pluginFactory, err := ocr3Capability.NewReportingPluginFactory(ctx, pluginCfg, nil,
		nil, nil, nil, capabilityRegistry, nil, nil)
	require.NoError(t, err)

	repConfig := ocr3types.ReportingPluginConfig{
		F: int(donF),
	}
	plugin, _, err := pluginFactory.NewReportingPlugin(ctx, repConfig)
	require.NoError(t, err)

	transmitter := ocr3.NewContractTransmitter(lggr, capabilityRegistry, "")

	libocr.AddNode(plugin, transmitter, ocr2KeyBundle)
}

func Context(tb testing.TB) context.Context {
	return testutils.Context(tb)
}
