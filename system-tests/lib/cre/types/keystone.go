package types

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/nix"
	"github.com/smartcontractkit/chainlink/system-tests/lib/types"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

type NodeType = string

const (
	BootstrapNode NodeType = "bootstrap"
	WorkerNode    NodeType = "worker"
	GatewayNode   NodeType = "gateway"
)

type ConfigDescription struct {
	Flag     CapabilityFlag
	NodeType string
}

type DonJobs = []*jobv1.ProposeJobRequest
type DonsToJobSpecs = map[uint32]DonJobs

type NodeIndexToConfigOverride = map[int]string
type NodeIndexToSecretsOverride = map[int]string

type WorkflowRegistryInput struct {
	ChainSelector  uint64                  `toml:"-"`
	CldEnv         *cldf.Environment       `toml:"-"`
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

type ConfigureDataFeedsCacheOutput struct {
	UseCache              bool             `toml:"use_cache"`
	DataFeedsCacheAddress common.Address   `toml:"data_feeds_cache_address"`
	FeedIDs               []string         `toml:"feed_is"`
	Descriptions          []string         `toml:"descriptions"`
	AdminAddress          common.Address   `toml:"admin_address"`
	AllowedSenders        []common.Address `toml:"allowed_senders"`
	AllowedWorkflowOwners []common.Address `toml:"allowed_workflow_owners"`
	AllowedWorkflowNames  []string         `toml:"allowed_workflow_names"`
}

type ConfigureDataFeedsCacheInput struct {
	CldEnv                *cldf.Environment              `toml:"-"`
	ChainSelector         uint64                         `toml:"-"`
	FeedIDs               []string                       `toml:"-"`
	Descriptions          []string                       `toml:"-"`
	DataFeedsCacheAddress common.Address                 `toml:"-"`
	AdminAddress          common.Address                 `toml:"-"`
	AllowedSenders        []common.Address               `toml:"-"`
	AllowedWorkflowOwners []common.Address               `toml:"-"`
	AllowedWorkflowNames  []string                       `toml:"-"`
	Out                   *ConfigureDataFeedsCacheOutput `toml:"out"`
}

func (c *ConfigureDataFeedsCacheInput) Validate() error {
	if c.CldEnv == nil {
		return errors.New("chainlink deployment env not set")
	}
	if len(c.FeedIDs) == 0 {
		return errors.New("feed ids not set")
	}
	if len(c.Descriptions) == 0 {
		return errors.New("descriptions not set")
	}
	if c.ChainSelector == 0 {
		return errors.New("chain selector not set")
	}
	if c.DataFeedsCacheAddress == (common.Address{}) {
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

	if (len(c.AllowedWorkflowNames) != len(c.AllowedWorkflowOwners)) || (len(c.AllowedWorkflowNames) != len(c.AllowedSenders)) {
		return errors.New("allowed workflow names, owners and senders must have the same length")
	}

	return nil
}

type WrappedNodeOutput struct {
	*ns.Output
	NodeSetName  string
	Capabilities []string
}

type CreateJobsInput struct {
	CldEnv        *cldf.Environment
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
	CldEnv        *cldf.Environment
	OCR3Config    keystone_changeset.OracleConfig
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

type ConfigFactoryFn = func(input GenerateConfigsInput) (NodeIndexToConfigOverride, error)

type GenerateConfigsInput struct {
	DonMetadata            *DonMetadata
	BlockchainOutput       map[uint64]*blockchain.Output
	HomeChainSelector      uint64
	Flags                  []string
	PeeringData            CapabilitiesPeeringData
	AddressBook            cldf.AddressBook
	GatewayConnectorOutput *GatewayConnectorOutput // optional, automatically set if some DON in the topology has the GatewayDON flag
}

func (g *GenerateConfigsInput) Validate() error {
	if len(g.DonMetadata.NodesMetadata) == 0 {
		return errors.New("don nodes not set")
	}
	if len(g.BlockchainOutput) == 0 {
		return errors.New("blockchain output not set")
	}
	if g.HomeChainSelector == 0 {
		return errors.New("home chain selector not set")
	}
	if len(g.Flags) == 0 {
		return errors.New("flags not set")
	}
	if g.PeeringData == (CapabilitiesPeeringData{}) {
		return errors.New("peering data not set")
	}
	_, addrErr := g.AddressBook.AddressesForChain(g.HomeChainSelector)
	if addrErr != nil {
		return errors.Wrapf(addrErr, "failed to get addresses for chain %d", g.HomeChainSelector)
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
	HomeChainSelector      uint64
	DonsMetadata           []*DonMetadata
	GatewayConnectorOutput *GatewayConnectorOutput
}

type DonTopology struct {
	WorkflowDonID          uint32
	HomeChainSelector      uint64
	DonsWithMetadata       []*DonWithMetadata
	GatewayConnectorOutput *GatewayConnectorOutput
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
	Out                        *GenerateKeysOutput
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

// chainID -> EVMKeys
type ChainIDToEVMKeys = map[int]*types.EVMKeys

// donID -> chainID -> EVMKeys
type DonsToEVMKeys = map[uint32]ChainIDToEVMKeys

// donID -> P2PKeys
type DonsToP2PKeys = map[uint32]*types.P2PKeys

type GenerateKeysOutput struct {
	EVMKeys DonsToEVMKeys
	P2PKeys DonsToP2PKeys
}

type GenerateSecretsInput struct {
	DonMetadata *DonMetadata
	EVMKeys     ChainIDToEVMKeys
	P2PKeys     *types.P2PKeys
}

func (g *GenerateSecretsInput) Validate() error {
	if g.DonMetadata == nil {
		return errors.New("don metadata not set")
	}
	if g.EVMKeys != nil {
		if len(g.EVMKeys) == 0 {
			return errors.New("chain ids not set")
		}
		for chainID, evmKeys := range g.EVMKeys {
			if len(evmKeys.EncryptedJSONs) == 0 {
				return errors.New("encrypted jsons not set")
			}
			if len(evmKeys.PublicAddresses) == 0 {
				return errors.New("public addresses not set")
			}
			if len(evmKeys.EncryptedJSONs) != len(evmKeys.PublicAddresses) {
				return errors.New("encrypted jsons and public addresses must have the same length")
			}
			if chainID == 0 {
				return errors.New("chain id 0 not allowed")
			}
		}
	}
	if g.P2PKeys != nil {
		if len(g.P2PKeys.EncryptedJSONs) == 0 {
			return errors.New("encrypted jsons not set")
		}
		if len(g.P2PKeys.PeerIDs) == 0 {
			return errors.New("peer ids not set")
		}
		if len(g.P2PKeys.EncryptedJSONs) != len(g.P2PKeys.PeerIDs) {
			return errors.New("encrypted jsons and peer ids must have the same length")
		}
	}

	return nil
}

type FullCLDEnvironmentInput struct {
	JdOutput          *jd.Output
	BlockchainOutputs map[uint64]*blockchain.Output
	SethClients       map[uint64]*seth.Client
	NodeSetOutput     []*WrappedNodeOutput
	ExistingAddresses cldf.AddressBook
	Topology          *Topology
}

func (f *FullCLDEnvironmentInput) Validate() error {
	if f.JdOutput == nil {
		return errors.New("jd output not set")
	}
	if len(f.BlockchainOutputs) == 0 {
		return errors.New("blockchain output not set")
	}
	if len(f.SethClients) == 0 {
		return errors.New("seth clients are not set")
	}
	if len(f.BlockchainOutputs) != len(f.SethClients) {
		return errors.New("blockchain outputs and seth clients must have the same length")
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
	Environment *cldf.Environment
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
	PurgeNamespace bool
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

type DONCapabilityWithConfigFactoryFn = func(donFlags []CapabilityFlag) []keystone_changeset.DONCapabilityWithConfig
type CapabilitiesBinaryPathFactoryFn = func(donMetadata *DonMetadata) ([]string, error)
type JobSpecFactoryFn = func(input *JobSpecFactoryInput) (DonsToJobSpecs, error)

type JobSpecFactoryInput struct {
	CldEnvironment   *cldf.Environment
	BlockchainOutput *blockchain.Output
	DonTopology      *DonTopology
	AddressBook      cldf.AddressBook
}

type ManageWorkflowWithCRECLIInput struct {
	DoNotUseCRECLI           bool
	ShouldCompileNewWorkflow bool
	ChainSelector            uint64
	WorkflowName             string
	WorkflowDonID            uint32
	WorkflowOwnerAddress     common.Address
	CRECLIPrivateKey         string
	CRECLIAbsPath            string
	CRESettingsFile          *os.File
	NewWorkflow              *NewWorkflow
	ExistingWorkflow         *ExistingWorkflow
	CRECLIProfile            string
}

type NewWorkflow struct {
	WorkflowFileName string
	FolderLocation   string
	ConfigFilePath   *string
	SecretsFilePath  *string
	Secrets          map[string]string
}

type ExistingWorkflow struct {
	BinaryURL  string
	ConfigURL  *string
	SecretsURL *string
}

func (w *ManageWorkflowWithCRECLIInput) Validate() error {
	if w.ChainSelector == 0 {
		return errors.New("ChainSelector is required")
	}
	if w.WorkflowName == "" {
		return errors.New("WorkflowName is required")
	}
	if w.WorkflowDonID == 0 {
		return errors.New("WorkflowDonID is required")
	}
	if w.CRECLIPrivateKey == "" {
		return errors.New("CRECLIPrivateKey is required")
	}
	if w.CRESettingsFile == nil {
		return errors.New("CRESettingsFile is required")
	}
	if w.NewWorkflow != nil && w.ExistingWorkflow != nil {
		return errors.New("only one of NewWorkflow or ExistingWorkflow can be provided")
	}

	return nil
}
