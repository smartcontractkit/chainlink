package types

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/system-tests/lib/nix"
	"github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

type NodeType = string

const (
	BootstrapNode NodeType = "bootstrap"
	WorkerNode    NodeType = "worker"
	GatewayNode   NodeType = "gateway"
)

type JobDescription struct {
	Flag     CapabilityFlag
	NodeType string
}

type ConfigDescription struct {
	Flag     CapabilityFlag
	NodeType string
}

type DonJobs = map[JobDescription][]*jobv1.ProposeJobRequest
type DonsToJobSpecs = map[uint32]DonJobs

type NodeIndexToConfigOverride = map[int]string
type NodeIndexToSecretsOverride = map[int]string

type KeystoneContractsInput struct {
	ChainSelector uint64                  `toml:"-"`
	CldEnv        *deployment.Environment `toml:"-"`
	Out           *KeystoneContractOutput `toml:"out"`
}

func (k *KeystoneContractsInput) Validate() error {
	if k.ChainSelector == 0 {
		return errors.New("chain selector not set")
	}
	if k.CldEnv == nil {
		return errors.New("chainlink deployment env not set")
	}
	return nil
}

type KeystoneContractOutput struct {
	UseCache                    bool           `toml:"use_cache"`
	CapabilitiesRegistryAddress common.Address `toml:"capabilities_registry_address"`
	ForwarderAddress            common.Address `toml:"forwarder_address"`
	OCR3CapabilityAddress       common.Address `toml:"ocr3_capability_address"`
	WorkflowRegistryAddress     common.Address `toml:"workflow_registry_address"`
}

type WorkflowRegistryInput struct {
	ChainSelector  uint64                  `toml:"-"`
	CldEnv         *deployment.Environment `toml:"-"`
	AllowedDonIDs  []uint32                `toml:"-"`
	WorkflowOwners []common.Address        `toml:"-"`
	Out            *WorkflowRegistryOutput `toml:"out"`
}

func (w *WorkflowRegistryInput) Validate() error {
	if w.ChainSelector == 0 {
		return errors.New("chain selector not set")
	}
	if w.CldEnv == nil {
		return errors.New("chainlink deployment env not set")
	}
	if len(w.AllowedDonIDs) == 0 {
		return errors.New("allowed don ids not set")
	}
	if len(w.WorkflowOwners) == 0 {
		return errors.New("workflow owners not set")
	}

	return nil
}

type WorkflowRegistryOutput struct {
	UseCache       bool             `toml:"use_cache"`
	ChainSelector  uint64           `toml:"chain_selector"`
	AllowedDonIDs  []uint32         `toml:"allowed_don_ids"`
	WorkflowOwners []common.Address `toml:"workflow_owners"`
}

type DeployFeedConsumerInput struct {
	ChainSelector uint64                    `toml:"-"`
	CldEnv        *deployment.Environment   `toml:"-"`
	Out           *DeployFeedConsumerOutput `toml:"out"`
}

func (i *DeployFeedConsumerInput) Validate() error {
	if i.ChainSelector == 0 {
		return errors.New("chain selector not set")
	}
	if i.CldEnv == nil {
		return errors.New("chainlink deployment env not set")
	}
	return nil
}

type DeployFeedConsumerOutput struct {
	UseCache            bool           `toml:"use_cache"`
	FeedConsumerAddress common.Address `toml:"feed_consumer_address"`
}

type ConfigureFeedConsumerInput struct {
	SethClient            *seth.Client                 `toml:"-"`
	FeedConsumerAddress   common.Address               `toml:"-"`
	AllowedSenders        []common.Address             `toml:"-"`
	AllowedWorkflowOwners []common.Address             `toml:"-"`
	AllowedWorkflowNames  []string                     `toml:"-"`
	Out                   *ConfigureFeedConsumerOutput `toml:"out"`
}

func (c *ConfigureFeedConsumerInput) Validate() error {
	if c.SethClient == nil {
		return errors.New("seth client not set")
	}
	if c.FeedConsumerAddress == (common.Address{}) {
		return errors.New("feed consumer address not set")
	}
	if len(c.AllowedSenders) == 0 {
		return errors.New("allowed senders not set")
	}
	if len(c.AllowedWorkflowOwners) == 0 {
		return errors.New("allowed workflow owners not set")
	}
	if len(c.AllowedWorkflowNames) == 0 {
		return errors.New("allowed workflow names not set")
	}

	return nil
}

type ConfigureFeedConsumerOutput struct {
	UseCache              bool             `toml:"use_cache"`
	FeedConsumerAddress   common.Address   `toml:"feed_consumer_address"`
	AllowedSenders        []common.Address `toml:"allowed_senders"`
	AllowedWorkflowOwners []common.Address `toml:"allowed_workflow_owners"`
	AllowedWorkflowNames  []string         `toml:"allowed_workflow_names"`
}

type WrappedNodeOutput struct {
	*ns.Output
	NodeSetName  string
	Capabilities []string
}

type CreateJobsInput struct {
	CldEnv        *deployment.Environment
	DonTopology   *DonTopology
	DonToJobSpecs DonsToJobSpecs
}

func (c *CreateJobsInput) Validate() error {
	if c.CldEnv == nil {
		return errors.New("chainlink deployment env not set")
	}
	if c.DonTopology == nil {
		return errors.New("don topology not set")
	}
	if len(c.DonTopology.DonsWithMetadata) == 0 {
		return errors.New("topology dons not set")
	}
	if len(c.DonToJobSpecs) == 0 {
		return errors.New("don to job specs not set")
	}

	return nil
}

type DebugInput struct {
	DebugDons        []*DebugDon
	BlockchainOutput *blockchain.Output
	InfraInput       *types.InfraInput
}

type DebugDon struct {
	Flags          []string
	ContainerNames []string
	NodesMetadata  []*NodeMetadata
}

func (d *DebugInput) Validate() error {
	if d.DebugDons == nil {
		return errors.New("don topology not set")
	}
	if len(d.DebugDons) == 0 {
		return errors.New("debug don not set")
	}
	for _, don := range d.DebugDons {
		if len(don.ContainerNames) == 0 {
			return errors.New("container names not set")
		}
		if len(don.NodesMetadata) == 0 {
			return errors.New("nodes metadata not set")
		}
		if len(don.Flags) == 0 {
			return errors.New("flags not set")
		}
	}
	if d.BlockchainOutput == nil {
		return errors.New("blockchain output not set")
	}
	if d.InfraInput == nil {
		return errors.New("infra input not set")
	}

	return nil
}

type ConfigureKeystoneInput struct {
	ChainSelector uint64
	Topology      *Topology
	CldEnv        *deployment.Environment
}

func (c *ConfigureKeystoneInput) Validate() error {
	if c.ChainSelector == 0 {
		return errors.New("chain selector not set")
	}
	if c.Topology == nil {
		return errors.New("don topology not set")
	}
	if len(c.Topology.DonsMetadata) == 0 {
		return errors.New("meta dons not set")
	}
	if c.CldEnv == nil {
		return errors.New("chainlink deployment env not set")
	}

	return nil
}

type GatewayConnectorDons struct {
	MembersEthAddresses []string
	ID                  uint32
}

type GatewayConnectorOutput struct {
	Dons []GatewayConnectorDons // do not set, it will be set dynamically
	Host string                 // do not set, it will be set dynamically
	Path string
	Port int
}

type GeneratePoRJobSpecsInput struct {
	DonsWithMetadata       []*DonWithMetadata
	BlockchainOutput       *blockchain.Output
	OCR3CapabilityAddress  common.Address
	ExtraAllowedPorts      []int
	ExtraAllowedIPs        []string
	CronCapBinName         string
	GatewayConnectorOutput GatewayConnectorOutput
}

func (g *GeneratePoRJobSpecsInput) Validate() error {
	if len(g.DonsWithMetadata) == 0 {
		return errors.New("metadata dons not set")
	}
	if g.BlockchainOutput == nil {
		return errors.New("blockchain output not set")
	}
	if g.OCR3CapabilityAddress == (common.Address{}) {
		return errors.New("ocr3 capability address not set")
	}
	if g.CronCapBinName == "" {
		return errors.New("cron cap bin name not set")
	}
	if g.GatewayConnectorOutput.Path == "" {
		return errors.New("gateway connector path is not set")
	}
	if g.GatewayConnectorOutput.Port == 0 {
		return errors.New("gateway connector port is not set")
	}

	return nil
}

type GeneratePoRConfigsInput struct {
	DonMetadata                 *DonMetadata
	BlockchainOutput            *blockchain.Output
	DonID                       uint32
	Flags                       []string
	PeeringData                 CapabilitiesPeeringData
	CapabilitiesRegistryAddress common.Address
	WorkflowRegistryAddress     common.Address
	ForwarderAddress            common.Address
	GatewayConnectorOutput      *GatewayConnectorOutput
}

func (g *GeneratePoRConfigsInput) Validate() error {
	if len(g.DonMetadata.NodesMetadata) == 0 {
		return errors.New("don nodes not set")
	}
	if g.BlockchainOutput == nil {
		return errors.New("blockchain output not set")
	}
	if g.DonID == 0 {
		return errors.New("don id not set")
	}
	if len(g.Flags) == 0 {
		return errors.New("flags not set")
	}
	if g.PeeringData == (CapabilitiesPeeringData{}) {
		return errors.New("peering data not set")
	}
	if g.CapabilitiesRegistryAddress == (common.Address{}) {
		return errors.New("capabilities registry address not set")
	}
	if g.WorkflowRegistryAddress == (common.Address{}) {
		return errors.New("workflow registry address not set")
	}
	if g.ForwarderAddress == (common.Address{}) {
		return errors.New("forwarder address not set")
	}
	if g.GatewayConnectorOutput == nil {
		return errors.New("gateway connector output not set")
	}

	return nil
}

type ToplogyInput struct {
	NodeSetInput    []*CapabilitiesAwareNodeSet
	DonToEthAddress map[uint32][]common.Address
}

type DonWithMetadata struct {
	DON *devenv.DON
	*DonMetadata
}

type DonMetadata struct {
	NodesMetadata []*NodeMetadata
	Flags         []string
	ID            uint32
	Name          string
}

type Label struct {
	Key   string
	Value string
}

func LabelFromProto(p *ptypes.Label) (*Label, error) {
	if p.Value == nil {
		return nil, errors.New("value not set")
	}
	return &Label{
		Key:   p.Key,
		Value: *p.Value,
	}, nil
}

type NodeMetadata struct {
	Labels []*Label
}

type Topology struct {
	WorkflowDONID          uint32
	DonsMetadata           []*DonMetadata
	GatewayConnectorOutput *GatewayConnectorOutput
}

type DonTopology struct {
	WorkflowDonID    uint32
	DonsWithMetadata []*DonWithMetadata
}

type CapabilitiesAwareNodeSet struct {
	*ns.Input
	Capabilities       []string
	DONTypes           []string
	BootstrapNodeIndex int // -1 -> no bootstrap, only used if the DON doesn't hae the GatewayDON flag
	GatewayNodeIndex   int // -1 -> no gateway, only used if the DON has the GatewayDON flag
}

type CapabilitiesPeeringData struct {
	GlobalBootstraperPeerID string
	GlobalBootstraperHost   string
	Port                    int
}

type OCRPeeringData struct {
	OCRBootstraperPeerID string
	OCRBootstraperHost   string
	Port                 int
}

type GenerateKeysInput struct {
	GenerateEVMKeysForChainIDs []int
	GenerateP2PKeys            bool
	Topology                   *Topology
	Password                   string
}

func (g *GenerateKeysInput) Validate() error {
	if g.Topology == nil {
		return errors.New("topology not set")
	}
	if len(g.Topology.DonsMetadata) == 0 {
		return errors.New("metadata not set")
	}
	if g.Topology.WorkflowDONID == 0 {
		return errors.New("workflow don id not set")
	}
	return nil
}

type DonsToEVMKeys = map[uint32]*types.EVMKeys
type DonsToP2PKeys = map[uint32]*types.P2PKeys

type GenerateKeysOutput struct {
	EVMKeys DonsToEVMKeys
	P2PKeys DonsToP2PKeys
}

type GenerateSecretsInput struct {
	DonMetadata *DonMetadata
	EVMKeys     *types.EVMKeys
	P2PKeys     *types.P2PKeys
}

func (g *GenerateSecretsInput) Validate() error {
	if g.DonMetadata == nil {
		return errors.New("don metadata not set")
	}
	if g.EVMKeys != nil {
		if len(g.EVMKeys.ChainIDs) == 0 {
			return errors.New("chain ids not set")
		}
		if len(g.EVMKeys.EncryptedJSONs) == 0 {
			return errors.New("encrypted jsons not set")
		}
	}
	if g.P2PKeys != nil {
		if len(g.P2PKeys.EncryptedJSONs) == 0 {
			return errors.New("encrypted jsons not set")
		}
	}

	return nil
}

type FullCLDEnvironmentInput struct {
	JdOutput          *jd.Output
	BlockchainOutput  *blockchain.Output
	SethClient        *seth.Client
	NodeSetOutput     []*WrappedNodeOutput
	ExistingAddresses deployment.AddressBook
	Topology          *Topology
}

func (f *FullCLDEnvironmentInput) Validate() error {
	if f.JdOutput == nil {
		return errors.New("jd output not set")
	}
	if f.BlockchainOutput == nil {
		return errors.New("blockchain output not set")
	}
	if f.SethClient == nil {
		return errors.New("seth client not set")
	}
	if len(f.NodeSetOutput) == 0 {
		return errors.New("node set output not set")
	}
	if f.Topology == nil {
		return errors.New("topology not set")
	}
	if len(f.Topology.DonsMetadata) == 0 {
		return errors.New("metadata not set")
	}
	if f.Topology.WorkflowDONID == 0 {
		return errors.New("workflow don id not set")
	}
	return nil
}

type FullCLDEnvironmentOutput struct {
	Environment *deployment.Environment
	DonTopology *DonTopology
}

type DeployCribDonsInput struct {
	Topology       *Topology
	NodeSetInputs  []*CapabilitiesAwareNodeSet
	NixShell       *nix.Shell
	CribConfigsDir string
}

func (d *DeployCribDonsInput) Validate() error {
	if d.Topology == nil {
		return errors.New("topology not set")
	}
	if len(d.Topology.DonsMetadata) == 0 {
		return errors.New("metadata not set")
	}
	if d.NixShell == nil {
		return errors.New("nix shell not set")
	}
	if len(d.NodeSetInputs) == 0 {
		return errors.New("node set inputs not set")
	}
	if d.CribConfigsDir == "" {
		return errors.New("crib configs dir not set")
	}
	return nil
}

type DeployCribJdInput struct {
	JDInput        *jd.Input
	NixShell       *nix.Shell
	CribConfigsDir string
}

func (d *DeployCribJdInput) Validate() error {
	if d.JDInput == nil {
		return errors.New("jd input not set")
	}
	if d.NixShell == nil {
		return errors.New("nix shell not set")
	}
	if d.CribConfigsDir == "" {
		return errors.New("crib configs dir not set")
	}
	return nil
}

type DeployCribBlockchainInput struct {
	BlockchainInput *blockchain.Input
	NixShell        *nix.Shell
	CribConfigsDir  string
}

func (d *DeployCribBlockchainInput) Validate() error {
	if d.BlockchainInput == nil {
		return errors.New("blockchain input not set")
	}
	if d.NixShell == nil {
		return errors.New("nix shell not set")
	}
	if d.CribConfigsDir == "" {
		return errors.New("crib configs dir not set")
	}
	return nil
}

type StartNixShellInput struct {
	InfraInput     *types.InfraInput
	CribConfigsDir string
	ExtraEnvVars   map[string]string
}

func (s *StartNixShellInput) Validate() error {
	if s.InfraInput == nil {
		return errors.New("infra input not set")
	}
	if s.CribConfigsDir == "" {
		return errors.New("crib configs dir not set")
	}
	return nil
}
