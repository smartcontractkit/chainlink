package framework

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/compute"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

type DonContext struct {
	EthBlockchain         *EthBlockchain
	p2pNetwork            *FakeRageP2PNetwork
	capabilityRegistry    *CapabilitiesRegistry
	workflowRegistry      *WorkflowRegistry
	syncerFetcherFunc     artifacts.FetcherFunc
	computeFetcherFactory compute.FetcherFactory
}

func CreateDonContext(ctx context.Context, t *testing.T) DonContext {
	ethBlockchain := NewEthBlockchain(t, 1000, 1*time.Second)
	rageP2PNetwork := NewFakeRageP2PNetwork(ctx, t, 1000)
	capabilitiesRegistry := NewCapabilitiesRegistry(ctx, t, ethBlockchain)

	servicetest.Run(t, rageP2PNetwork)
	servicetest.Run(t, ethBlockchain)
	return DonContext{EthBlockchain: ethBlockchain, p2pNetwork: rageP2PNetwork, capabilityRegistry: capabilitiesRegistry}
}

func CreateDonContextWithWorkflowRegistry(ctx context.Context, t *testing.T, syncerFetcherFunc artifacts.FetcherFunc,
	computeFetcherFactory compute.FetcherFactory) DonContext {
	donContext := CreateDonContext(ctx, t)
	workflowRegistry := NewWorkflowRegistry(ctx, t, donContext.EthBlockchain)
	donContext.workflowRegistry = workflowRegistry
	donContext.syncerFetcherFunc = syncerFetcherFunc
	donContext.computeFetcherFactory = computeFetcherFactory
	return donContext
}

func (c DonContext) WaitForCapabilitiesToBeExposed(t *testing.T, dons ...*DON) {
	allExpectedCapabilities := make(map[CapabilityRegistration]bool)
	for _, don := range dons {
		caps, err := don.GetExternalCapabilities()
		require.NoError(t, err)
		for k, v := range caps {
			allExpectedCapabilities[k] = v
		}
	}

	require.Eventually(t, func() bool {
		registrations := c.p2pNetwork.GetCapabilityRegistrations()

		for k := range allExpectedCapabilities {
			if _, ok := registrations[k]; !ok {
				return false
			}
		}

		return true
	}, 1*time.Minute, 1*time.Second, "timeout waiting for capabilities to be exposed")
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
	services.StateMachine
	t                      *testing.T
	id                     *uint32
	config                 DonConfiguration
	lggr                   logger.Logger
	nodes                  []*capabilityNode
	standardCapabilityJobs []*job.Job
	publishedCapabilities  []capability

	initialised bool

	capabilitiesRegistry *CapabilitiesRegistry
	workflowRegistry     *WorkflowRegistry

	nodeConfigModifiers []func(c *chainlink.Config, node *capabilityNode)

	fakeLibOcr                   *FakeLibOCR
	addOCR3NonStandardCapability bool

	triggerFactories []TriggerFactory
	targetFactories  []TargetFactory
}

func NewDON(ctx context.Context, t *testing.T, lggr logger.Logger, donConfig DonConfiguration,
	dependentDONs []commoncap.DON, donContext DonContext, supportsOCR bool, protocolRoundInterval time.Duration) *DON {
	don := &DON{t: t, lggr: lggr.Named(donConfig.name), config: donConfig, capabilitiesRegistry: donContext.capabilityRegistry,
		workflowRegistry: donContext.workflowRegistry}

	var newOracleFactoryFn standardcapabilities.NewOracleFactoryFn
	if supportsOCR {
		// This is required to support the non standard OCR3 capability - will be removed when required OCR3 behaviour is implemented as standard capabilities
		don.fakeLibOcr = NewFakeLibOCR(t, lggr, donConfig.F, protocolRoundInterval)
		servicetest.Run(t, don.fakeLibOcr)
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
			factory := newFakeOracleFactoryFactory(t, lggr, donConfig.KeyBundles[i], len(donConfig.Members), donConfig.F,
				protocolRoundInterval)
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
				}, donContext.syncerFetcherFunc, donContext.computeFetcherFactory)

			require.NoError(t, node.Start(testutils.Context(t)))
			cn.TestApplication = node
		}
	}

	return don
}

// Initialise must be called after all capabilities have been added to the DONs and before Start is called
func (d *DON) Initialise() {
	id := d.capabilitiesRegistry.setupDON(d.config, d.publishedCapabilities)

	//nolint:gosec // disable G115
	d.config.DON.ID = uint32(id)
	d.id = &d.config.DON.ID

	if d.config.AcceptsWorkflows && d.workflowRegistry != nil {
		d.workflowRegistry.UpdateAllowedDons([]uint32{d.config.DON.ID})
		d.nodeConfigModifiers = append(d.nodeConfigModifiers, func(c *chainlink.Config, node *capabilityNode) {
			workflowRegistryAddressStr := d.workflowRegistry.addr.String()
			c.Capabilities.WorkflowRegistry.Address = &workflowRegistryAddressStr
			c.Capabilities.WorkflowRegistry.ChainID = ptr(fmt.Sprintf("%d", testutils.SimulatedChainID))
		})
	}
	d.initialised = true
}

func (d *DON) GetID() uint32 {
	if d.config.DON.ID == 0 {
		panic("DON ID not set, call Initialise() first")
	}

	return d.config.ID
}

func (d *DON) GetExternalCapabilities() (map[CapabilityRegistration]bool, error) {
	result := map[CapabilityRegistration]bool{}
	for _, publishedCapability := range d.publishedCapabilities {
		if publishedCapability.internalOnly {
			continue
		}

		for _, node := range d.nodes {
			peerIDBytes, err := peerIDToBytes(node.peerID.PeerID)
			if err != nil {
				return nil, fmt.Errorf("failed to convert peer ID to bytes: %w", err)
			}
			result[CapabilityRegistration{
				nodePeerID:      hex.EncodeToString(peerIDBytes[:]),
				capabilityID:    publishedCapability.registryConfig.LabelledName + "@" + publishedCapability.registryConfig.Version,
				capabilityDonID: d.GetID(),
			}] = true
		}
	}

	return result, nil
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

func (d *DON) Start(ctx context.Context) error {
	for _, triggerFactory := range d.triggerFactories {
		for _, node := range d.nodes {
			trigger := triggerFactory.CreateNewTrigger(d.t)
			if err := node.registry.Add(ctx, trigger); err != nil {
				return fmt.Errorf("failed to add trigger: %w", err)
			}
		}
	}

	for _, targetFactory := range d.targetFactories {
		for _, node := range d.nodes {
			target := targetFactory.CreateNewTarget(d.t)
			if err := node.registry.Add(ctx, target); err != nil {
				return fmt.Errorf("failed to add target: %w", err)
			}
		}
	}

	for _, node := range d.nodes {
		node.start()
	}

	if d.addOCR3NonStandardCapability {
		if d.fakeLibOcr == nil {
			return errors.New("don does not support OCR")
		}

		for _, node := range d.nodes {
			addOCR3Capability(ctx, d.t, d.lggr, node.registry, d.fakeLibOcr, d.config.F, node.KeyBundle)
		}
	}

	for _, capabilityJob := range d.standardCapabilityJobs {
		if err := d.AddJob(ctx, capabilityJob); err != nil {
			return fmt.Errorf("failed to add standard capability job: %w", err)
		}
	}

	return nil
}

func (d *DON) Close() error {
	for _, node := range d.nodes {
		if err := node.Stop(); err != nil {
			return fmt.Errorf("failed to stop node: %w", err)
		}
	}

	return nil
}

const StandardCapabilityTemplateJobSpec = `
type = "standardcapabilities"
schemaVersion = 1
name = "%s"
command="%s"
config=%s
`

func (d *DON) AddStandardCapability(name string, command string, config string) {
	spec := fmt.Sprintf(StandardCapabilityTemplateJobSpec, name, command, config)
	capabilitiesSpecJob, err := standardcapabilities.ValidatedStandardCapabilitiesSpec(spec)
	require.NoError(d.t, err)

	d.standardCapabilityJobs = append(d.standardCapabilityJobs, &capabilitiesSpecJob)
}

func (d *DON) AddPublishedStandardCapability(name string, command string, config string,
	defaultCapabilityRequestConfig *pb.CapabilityConfig,
	registryConfig kcr.CapabilitiesRegistryCapability) {
	spec := fmt.Sprintf(StandardCapabilityTemplateJobSpec, name, command, config)
	capabilitiesSpecJob, err := standardcapabilities.ValidatedStandardCapabilitiesSpec(spec)
	require.NoError(d.t, err)

	d.standardCapabilityJobs = append(d.standardCapabilityJobs, &capabilitiesSpecJob)

	d.publishedCapabilities = append(d.publishedCapabilities, capability{
		donCapabilityConfig: defaultCapabilityRequestConfig,
		registryConfig:      registryConfig,
	})
}

// TODO - add configuration for remote support - do this for each capability as an option
func (d *DON) AddTargetCapability(targetFactory TargetFactory) {
	d.targetFactories = append(d.targetFactories, targetFactory)
}

func (d *DON) AddTriggerCapability(triggerFactory TriggerFactory) {
	d.triggerFactories = append(d.triggerFactories, triggerFactory)
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

	d.publishedCapabilities = append(d.publishedCapabilities, triggerCapability)
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

func (d *DON) AddWorkflow(workflow Workflow) error {
	if !d.config.AcceptsWorkflows {
		return errors.New("cannot add workflow to non-workflow DON")
	}

	if !d.initialised {
		return errors.New("cannot add workflow to non-initialised DON")
	}

	d.workflowRegistry.RegisterWorkflow(workflow, *d.id)

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
	fetcherFunc artifacts.FetcherFunc,
	fetcherFactoryFunc compute.FetcherFactory,
) *cltest.TestApplication {
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Capabilities.ExternalRegistry.ChainID = ptr(fmt.Sprintf("%d", testutils.SimulatedChainID))
		c.Capabilities.ExternalRegistry.Address = ptr(capRegistryAddr.String())
		c.Capabilities.Peering.V2.Enabled = ptr(true)
		c.Feature.FeedsManager = ptr(false)
		c.Feature.LogPoller = ptr(true)

		if setupCfg != nil {
			setupCfg(c)
		}
	})

	n, err := ethBlockchain.Client().NonceAt(ctx, ethBlockchain.transactionOpts.From, nil)
	require.NoError(t, err)

	tx := cltest.NewLegacyTransaction(
		n, keyV2.Address,
		assets.Ether(1).ToInt(),
		21000,
		assets.GWei(1).ToInt(),
		nil)
	signedTx, err := ethBlockchain.transactionOpts.Signer(ethBlockchain.transactionOpts.From, tx)
	require.NoError(t, err)
	err = ethBlockchain.Client().SendTransaction(ctx, signedTx)
	require.NoError(t, err)
	ethBlockchain.Commit()

	return cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, ethBlockchain.Backend,
		nodeInfo, dispatcher, peerWrapper, newOracleFactoryFn, localCapabilities, keyV2, lggr, fetcherFunc,
		fetcherFactoryFunc)
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

	d.publishedCapabilities = append(d.publishedCapabilities, capability{
		donCapabilityConfig: newCapabilityConfig(),
		registryConfig:      ocr,
		internalOnly:        true,
	})
}

func (d *DON) AddPublishedEthereumWriteTargetNonStandardCapability(forwarderAddr common.Address) (string, error) {
	published := true

	capabilityID, s, err := d.addEthereumWriteTarget(forwarderAddr, published)
	if err != nil {
		return s, err
	}

	return capabilityID, nil
}

func (d *DON) AddEthereumWriteTargetNonStandardCapability(forwarderAddr common.Address) (string, error) {
	published := false

	capabilityID, s, err := d.addEthereumWriteTarget(forwarderAddr, published)
	if err != nil {
		return s, err
	}

	return capabilityID, nil
}

func (d *DON) addEthereumWriteTarget(forwarderAddr common.Address, published bool) (string, string, error) {
	d.nodeConfigModifiers = append(d.nodeConfigModifiers, func(c *chainlink.Config, node *capabilityNode) {
		eip55Address := types.EIP55AddressFromAddress(forwarderAddr)
		c.EVM[0].Chain.Workflow.ForwarderAddress = &eip55Address
		c.EVM[0].Chain.Workflow.FromAddress = &node.key.EIP55Address
	})

	labelledName := "write_geth-testnet"
	version := "1.0.0"

	writeChain := kcr.CapabilitiesRegistryCapability{
		LabelledName:   labelledName,
		Version:        version,
		CapabilityType: uint8(registrysyncer.ContractCapabilityTypeTarget),
	}

	capabilityID := fmt.Sprintf("%s@%s", labelledName, version)

	targetCapabilityConfig := newCapabilityConfig()

	configWithLimit, err := values.WrapMap(map[string]any{"gasLimit": 500000})
	if err != nil {
		return "", "", fmt.Errorf("failed to wrap map: %w", err)
	}

	targetCapabilityConfig.DefaultConfig = values.Proto(configWithLimit).GetMapValue()

	targetCapabilityConfig.RemoteConfig = &pb.CapabilityConfig_RemoteTargetConfig{
		RemoteTargetConfig: &pb.RemoteTargetConfig{
			RequestHashExcludedAttributes: []string{"signed_report.Signatures"},
		},
	}

	if published {
		d.publishedCapabilities = append(d.publishedCapabilities, capability{
			donCapabilityConfig: targetCapabilityConfig,
			registryConfig:      writeChain,
		})
	}
	return capabilityID, "", nil
}

func addOCR3Capability(ctx context.Context, t *testing.T, lggr logger.Logger, capabilityRegistry *capabilities.Registry,
	libocr *FakeLibOCR, donF uint8, ocr2KeyBundle ocr2key.KeyBundle) {
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
