package types

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

type NodeType = string

const (
	BootstrapNode NodeType = "bootstrap"
	WorkerNode    NodeType = "worker"
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

type NodeIndexToConfigOverrides = map[int]string
type DonsToConfigOverrides = map[uint32]NodeIndexToConfigOverrides

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

type ConfigureDonInput struct {
	CldEnv               *deployment.Environment
	BlockchainOutput     *blockchain.Output
	DonTopology          *DonTopology
	JdOutput             *jd.Output
	DonToJobSpecs        DonsToJobSpecs
	DonToConfigOverrides DonsToConfigOverrides
}

func (c *ConfigureDonInput) Validate() error {
	if c.CldEnv == nil {
		return errors.New("chainlink deployment env not set")
	}
	if c.BlockchainOutput == nil {
		return errors.New("blockchain output not set")
	}
	if c.DonTopology == nil {
		return errors.New("don topology not set")
	}
	if len(c.DonTopology.MetaDons) == 0 {
		return errors.New("meta dons not set")
	}
	if c.JdOutput == nil {
		return errors.New("jd output not set")
	}
	if len(c.DonToJobSpecs) == 0 {
		return errors.New("don to job specs not set")
	}
	if len(c.DonToConfigOverrides) == 0 {
		return errors.New("don to config overrides not set")
	}

	return nil
}

type ConfigureDonOutput struct {
	JdOutput *deployment.OffchainClient
}

type DebugInput struct {
	DonTopology      *DonTopology
	BlockchainOutput *blockchain.Output
}

func (d *DebugInput) Validate() error {
	if d.DonTopology == nil {
		return errors.New("don topology not set")
	}
	if len(d.DonTopology.MetaDons) == 0 {
		return errors.New("meta dons not set")
	}
	if d.BlockchainOutput == nil {
		return errors.New("blockchain output not set")
	}

	return nil
}

type ConfigureKeystoneInput struct {
	ChainSelector uint64
	DonTopology   *DonTopology
	CldEnv        *deployment.Environment
}

func (c *ConfigureKeystoneInput) Validate() error {
	if c.ChainSelector == 0 {
		return errors.New("chain selector not set")
	}
	if c.DonTopology == nil {
		return errors.New("don topology not set")
	}
	if len(c.DonTopology.MetaDons) == 0 {
		return errors.New("meta dons not set")
	}
	if c.CldEnv == nil {
		return errors.New("chainlink deployment env not set")
	}

	return nil
}

type GatewayConnectorOutput struct {
	Host string // do not set, it will be set dynamically
	Path string
	Port int
}

type GeneratePoRJobSpecsInput struct {
	CldEnv                 *deployment.Environment
	Don                    *devenv.DON
	NodeOutput             *WrappedNodeOutput
	BlockchainOutput       *blockchain.Output
	DonID                  uint32
	Flags                  []string
	OCR3CapabilityAddress  common.Address
	ExtraAllowedPorts      []int
	ExtraAllowedIPs        []string
	CronCapBinName         string
	GatewayConnectorOutput GatewayConnectorOutput
}

func (g *GeneratePoRJobSpecsInput) Validate() error {
	if g.CldEnv == nil {
		return errors.New("chainlink deployment env not set")
	}
	if g.Don == nil {
		return errors.New("don not set")
	}
	if len(g.Don.Nodes) == 0 {
		return errors.New("don nodes not set")
	}
	if g.NodeOutput == nil {
		return errors.New("node output not set")
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
	if g.OCR3CapabilityAddress == (common.Address{}) {
		return errors.New("ocr3 capability address not set")
	}
	if g.CronCapBinName == "" {
		return errors.New("cron cap bin name not set")
	}
	if g.GatewayConnectorOutput == (GatewayConnectorOutput{}) {
		return errors.New("gateway connector output not set")
	}

	return nil
}

type GeneratePoRConfigsInput struct {
	Don                         *devenv.DON
	NodeInput                   *CapabilitiesAwareNodeSet
	BlockchainOutput            *blockchain.Output
	DonID                       uint32
	Flags                       []string
	PeeringData                 PeeringData
	CapabilitiesRegistryAddress common.Address
	WorkflowRegistryAddress     common.Address
	ForwarderAddress            common.Address
	GatewayConnectorOutput      *GatewayConnectorOutput
}

func (g *GeneratePoRConfigsInput) Validate() error {
	if g.Don == nil {
		return errors.New("don not set")
	}
	if len(g.Don.Nodes) == 0 {
		return errors.New("don nodes not set")
	}
	if g.NodeInput == nil {
		return errors.New("node input not set")
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
	if g.PeeringData == (PeeringData{}) {
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

// DonWithMetadata is a struct that holds the DON references and various metadata
type DonWithMetadata struct {
	DON        *devenv.DON
	NodeInput  *CapabilitiesAwareNodeSet
	NodeOutput *WrappedNodeOutput
	ID         uint32
	Flags      []string
}

type DonTopology struct {
	WorkflowDONID uint32
	MetaDons      []*DonWithMetadata
}

type CapabilitiesAwareNodeSet struct {
	*ns.Input
	Capabilities []string `toml:"capabilities"`
	DONType      string   `toml:"don_type"`
}

type PeeringData struct {
	GlobalBootstraperPeerID string
	GlobalBootstraperHost   string
	Port                    int
}

type OCRPeeringData struct {
	OCRBootstraperPeerID string
	OCRBootstraperHost   string
	Port                 int
}
