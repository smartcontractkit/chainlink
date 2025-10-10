package cre

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	ks_sol "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"

	cldf_jd "github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
)

type CapabilityFlag = string

// DON types
const (
	WorkflowDON     CapabilityFlag = "workflow"
	CapabilitiesDON CapabilityFlag = "capabilities"
	GatewayDON      CapabilityFlag = "gateway"
)

// Capabilities
const (
	ConsensusCapability     CapabilityFlag = "ocr3"
	ConsensusCapabilityV2   CapabilityFlag = "consensus" // v2
	CronCapability          CapabilityFlag = "cron"
	EVMCapability           CapabilityFlag = "evm"
	CustomComputeCapability CapabilityFlag = "custom-compute"
	WriteEVMCapability      CapabilityFlag = "write-evm"
	WriteSolanaCapability   CapabilityFlag = "write-solana"
	ReadContractCapability  CapabilityFlag = "read-contract"
	LogTriggerCapability    CapabilityFlag = "log-event-trigger"
	WebAPITargetCapability  CapabilityFlag = "web-api-target"
	WebAPITriggerCapability CapabilityFlag = "web-api-trigger"
	MockCapability          CapabilityFlag = "mock"
	VaultCapability         CapabilityFlag = "vault"
	HTTPTriggerCapability   CapabilityFlag = "http-trigger"
	HTTPActionCapability    CapabilityFlag = "http-action"
	// Add more capabilities as needed
)

type CLIEnvironmentDependencies interface {
	CapabilityFlagsProvider
	ContractVersionsProvider
	CLIFlagsProvider
}

// CLIFlagsProvider provides access to select command line flags passed to the
// start command of the environment script.
type CLIFlagsProvider interface {
	// If true, then use V2 Capability and Workflow Registries.
	WithV2Registries() bool
}

func NewCLIFlagsProvider(withV2Registries bool) *cliFlagsProvider {
	return &cliFlagsProvider{
		withV2Registries: withV2Registries,
	}
}

type cliFlagsProvider struct {
	withV2Registries bool
}

func (cfp *cliFlagsProvider) WithV2Registries() bool {
	return cfp.withV2Registries
}

type ContractVersionsProvider interface {
	// ContractVersions returns a map of contract name to semver
	ContractVersions() map[string]string
}

type contractVersionsProvider struct {
	contracts map[string]string
}

func (cvp *contractVersionsProvider) ContractVersions() map[string]string {
	cv := make(map[string]string, 0)
	maps.Copy(cv, cvp.contracts)
	return cv
}

func NewContractVersionsProvider(overrides map[string]string) *contractVersionsProvider {
	cvp := &contractVersionsProvider{
		contracts: map[string]string{
			keystone_changeset.OCR3Capability.String():       "1.0.0",
			keystone_changeset.WorkflowRegistry.String():     "1.0.0",
			keystone_changeset.CapabilitiesRegistry.String(): "1.1.0",
			keystone_changeset.KeystoneForwarder.String():    "1.0.0",
			ks_sol.ForwarderContract.String():                "1.0.0",
			ks_sol.ForwarderState.String():                   "1.0.0",
		},
	}
	maps.Copy(cvp.contracts, overrides)
	return cvp
}

type CapabilityFlagsProvider interface {
	SupportedCapabilityFlags() []CapabilityFlag
	GlobalCapabilityFlags() []CapabilityFlag
	ChainSpecificCapabilityFlags() []CapabilityFlag
}

func NewEnvironmentDependencies(
	cfp CapabilityFlagsProvider,
	cvp ContractVersionsProvider,
	cliFlagsProvider CLIFlagsProvider,
) *envionmentDependencies {
	return &envionmentDependencies{
		flagsProvider:       cfp,
		contractSetProvider: cvp,
		cliFlagsProvider:    cliFlagsProvider,
	}
}

type envionmentDependencies struct {
	flagsProvider       CapabilityFlagsProvider
	contractSetProvider ContractVersionsProvider
	cliFlagsProvider    CLIFlagsProvider
}

func (e *envionmentDependencies) WithV2Registries() bool {
	return e.cliFlagsProvider.WithV2Registries()
}

func (e *envionmentDependencies) ContractVersions() map[string]string {
	return e.contractSetProvider.ContractVersions()
}

func (e *envionmentDependencies) SupportedCapabilityFlags() []CapabilityFlag {
	return e.flagsProvider.SupportedCapabilityFlags()
}

func (e *envionmentDependencies) GlobalCapabilityFlags() []CapabilityFlag {
	return e.flagsProvider.GlobalCapabilityFlags()
}

func (e *envionmentDependencies) ChainSpecificCapabilityFlags() []CapabilityFlag {
	return e.flagsProvider.ChainSpecificCapabilityFlags()
}

type NodeType = string

const (
	BootstrapNode NodeType = "bootstrap"
	GatewayNode   NodeType = "gateway"

	// WorkerNode The value here is `plugin` to match the filtering performed by JD to get non-bootstrap nodes.
	// See: https://github.com/smartcontractkit/chainlink/blob/develop/deployment/data-feeds/offchain/jd.go#L57
	WorkerNode NodeType = "plugin"
)

type (
	DonJobs        = []*jobv1.ProposeJobRequest
	DonsToJobSpecs = map[uint64]DonJobs
)

const (
	CapabilityLabelKey = "capability"
)

// TODO stop using index to identify nodes, use some unique ID instead
type (
	NodeIndexToConfigOverride  = map[int]string
	NodeIndexToSecretsOverride = map[int]string
)

type CapabilityConfigs = map[string]CapabilityConfig

type CapabilityConfig struct {
	BinaryPath   string         `toml:"binary_path"`
	Config       map[string]any `toml:"config"`
	Chains       []string       `toml:"chains"`
	ChainConfigs map[string]any `toml:"chain_configs"`
}

type WorkflowRegistryInput struct {
	ContractAddress common.Address          `toml:"_"`
	ContractVersion cldf.TypeAndVersion     `toml:"_"`
	ChainSelector   uint64                  `toml:"-"`
	CldEnv          *cldf.Environment       `toml:"-"`
	AllowedDonIDs   []uint64                `toml:"-"`
	WorkflowOwners  []common.Address        `toml:"-"`
	Out             *WorkflowRegistryOutput `toml:"out"`
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

func (c *WorkflowRegistryOutput) Store(absPath string) error {
	framework.L.Info().Msgf("Storing Workflow Registry state file: %s", absPath)
	return storeLocalArtifact(c, absPath)
}

func (c WorkflowRegistryOutput) WorkflowOwnersStrings() []string {
	owners := make([]string, len(c.WorkflowOwners))
	for idx, owner := range c.WorkflowOwners {
		owners[idx] = owner.String()
	}

	return owners
}

func storeLocalArtifact(artifact any, absPath string) error {
	dErr := os.MkdirAll(filepath.Dir(absPath), 0755)
	if dErr != nil {
		return errors.Wrap(dErr, "failed to create directory for the environment artifact")
	}

	d, mErr := toml.Marshal(artifact)
	if mErr != nil {
		return errors.Wrap(mErr, "failed to marshal environment artifact to TOML")
	}

	return os.WriteFile(absPath, d, 0600)
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

type WrappedBlockchainOutput struct {
	ChainSelector      uint64
	ChainID            uint64
	BlockchainOutput   *blockchain.Output
	SethClient         *seth.Client
	SolClient          *solrpc.Client
	DeployerPrivateKey string
	SolChain           *SolChain
}

type SolChain struct {
	ChainSelector uint64
	ChainID       string
	ChainName     string
	PrivateKey    solana.PrivateKey
	ArtifactsDir  string
}

type ConfigureKeystoneInput struct {
	ChainSelector               uint64
	Topology                    *Topology
	CldEnv                      *cldf.Environment
	NodeSets                    []*CapabilitiesAwareNodeSet
	CapabilityRegistryConfigFns []CapabilityRegistryConfigFn
	BlockchainOutputs           []*WrappedBlockchainOutput

	OCR3Config  keystone_changeset.OracleConfig
	OCR3Address *common.Address // v1 consensus contract address

	DONTimeConfig  keystone_changeset.OracleConfig
	DONTimeAddress *common.Address

	VaultOCR3Config  keystone_changeset.OracleConfig
	VaultOCR3Address *common.Address

	DKGReportingPluginConfig *dkgocrtypes.ReportingPluginConfig
	DKGOCR3Config            keystone_changeset.OracleConfig
	DKGOCR3Address           *common.Address

	EVMOCR3Config    keystone_changeset.OracleConfig
	EVMOCR3Addresses map[uint64]common.Address // chain selector to address map

	ConsensusV2OCR3Config  keystone_changeset.OracleConfig // v2 consensus contract config
	ConsensusV2OCR3Address *common.Address

	CapabilitiesRegistryAddress *common.Address

	WithV2Registries bool
}

func (c *ConfigureKeystoneInput) Validate() error {
	if c.ChainSelector == 0 {
		return errors.New("chain selector not set")
	}
	if c.Topology == nil {
		return errors.New("don topology not set")
	}
	if len(c.Topology.DonsMetadata.List()) == 0 {
		return errors.New("meta dons not set")
	}
	if len(c.NodeSets) != len(c.Topology.DonsMetadata.List()) {
		return errors.New("node sets and don metadata must have the same length")
	}
	if c.CldEnv == nil {
		return errors.New("chainlink deployment env not set")
	}
	if c.OCR3Address == nil || c.CapabilitiesRegistryAddress == nil {
		return errors.New("OCR3Address and CapabilitiesRegistryAddress must be set")
	}

	return nil
}

type GatewayConnectorDons struct {
	MembersEthAddresses []string `toml:"members_eth_addresses" json:"members_eth_addresses"`
	ID                  string   `toml:"id" json:"id"`
	Handlers            map[string]string
}
type GatewayConnectorOutput struct {
	Configurations []*DonGatewayConfiguration `toml:"configurations" json:"configurations"`
}

func NewGatewayConnectorOutput() *GatewayConnectorOutput {
	return &GatewayConnectorOutput{
		Configurations: make([]*DonGatewayConfiguration, 0),
	}
}

type DonGatewayConfiguration struct {
	Dons []GatewayConnectorDons `toml:"dons" json:"dons"` // do not set, it will be set dynamically
	*GatewayConfiguration
}

type NodeConfigTransformerFn = func(input GenerateConfigsInput, existingConfigs NodeIndexToConfigOverride) (NodeIndexToConfigOverride, error)

type (
	HandlerTypeToConfig    = map[string]string
	GatewayHandlerConfigFn = func(don *DON) (HandlerTypeToConfig, error)
)

type GenerateConfigsInput struct {
	Datastore               datastore.DataStore
	DonMetadata             *DonMetadata
	BlockchainOutput        map[uint64]*WrappedBlockchainOutput
	HomeChainSelector       uint64
	Flags                   []string
	CapabilitiesPeeringData CapabilitiesPeeringData
	OCRPeeringData          OCRPeeringData
	AddressBook             cldf.AddressBook
	NodeSet                 *CapabilitiesAwareNodeSet
	CapabilityConfigs       CapabilityConfigs
	GatewayConnectorOutput  *GatewayConnectorOutput // optional, automatically set if some DON in the topology has the GatewayDON flag
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
	if g.CapabilitiesPeeringData == (CapabilitiesPeeringData{}) {
		return errors.New("peering data not set")
	}
	if g.OCRPeeringData == (OCRPeeringData{}) {
		return errors.New("ocr peering data not set")
	}
	_, addrErr := g.AddressBook.AddressesForChain(g.HomeChainSelector)
	if addrErr != nil {
		return fmt.Errorf("failed to get addresses for chain %d: %w", g.HomeChainSelector, addrErr)
	}
	_, dsErr := g.Datastore.Addresses().Fetch()
	if dsErr != nil {
		return fmt.Errorf("failed to get addresses from datastore: %w", dsErr)
	}
	h := g.Datastore.Addresses().Filter(datastore.AddressRefByChainSelector(g.HomeChainSelector))
	if len(h) == 0 {
		return fmt.Errorf("no addresses found for home chain %d in datastore", g.HomeChainSelector)
	}
	// TODO check for required registry contracts by type and version
	return nil
}

type DonMetadata struct {
	NodesMetadata []*NodeMetadata `toml:"nodes_metadata" json:"nodes_metadata"`
	Flags         []string        `toml:"flags" json:"flags"`
	ID            uint64          `toml:"id" json:"id"`
	Name          string          `toml:"name" json:"name"`

	ns CapabilitiesAwareNodeSet // computed field, not serialized
	gh GatewayHelper
}

func NewDonMetadata(c *CapabilitiesAwareNodeSet, id uint64, provider infra.Provider) (*DonMetadata, error) {
	cfgs := make([]NodeMetadataConfig, len(c.NodeSpecs))
	for i, nodeSpec := range c.NodeSpecs {
		nodeType := WorkerNode
		if c.BootstrapNodeIndex != -1 && i == c.BootstrapNodeIndex {
			nodeType = BootstrapNode
		}

		cfg := NodeMetadataConfig{
			Keys: NodeKeyInput{
				EVMChainIDs:     c.EVMChains(),
				SolanaChainIDs:  c.SupportedSolChains,
				Password:        "dev-password",
				ImportedSecrets: nodeSpec.Node.TestSecretsOverrides,
			},
			Host:  provider.InternalHost(i, nodeType == BootstrapNode, c.Name),
			Roles: []string{nodeType},
			Index: i,
		}

		if slices.Contains(c.DONTypes, GatewayDON) && c.GatewayNodeIndex != -1 && i == c.GatewayNodeIndex {
			cfg.Roles = append(cfg.Roles, GatewayNode)
		}

		cfgs[i] = cfg
	}

	nodes, err := newNodes(cfgs)
	if err != nil {
		return nil, fmt.Errorf("failed to create nodes metadata: %w", err)
	}
	out := &DonMetadata{
		ID:            id,
		Flags:         c.Flags(),
		NodesMetadata: nodes,
		Name:          c.Name,
		ns:            *c,
	}

	return out, nil
}

func (m *DonMetadata) GatewayConfig(p infra.Provider) (*DonGatewayConfiguration, error) {
	gatewayNode, hasGateway := m.Gateway()
	if !hasGateway {
		return nil, errors.New("don does not have a gateway node")
	}

	return &DonGatewayConfiguration{
		Dons:                 make([]GatewayConnectorDons, 0),
		GatewayConfiguration: NewGatewayConfig(p, gatewayNode.Index, gatewayNode.HasRole(BootstrapNode), m.Name),
	}, nil
}

func (m *DonMetadata) Workers() ([]*NodeMetadata, error) {
	workers := make([]*NodeMetadata, 0)
	for _, node := range m.NodesMetadata {
		if slices.Contains(node.Roles, WorkerNode) {
			workers = append(workers, node)
		}
	}

	if len(workers) == 0 {
		return nil, errors.New("don does not contain any worker nodes")
	}

	return workers, nil
}

// Currently only one bootstrap node is supported.
func (m *DonMetadata) Bootstrap() (*NodeMetadata, bool) {
	for _, node := range m.NodesMetadata {
		if slices.Contains(node.Roles, BootstrapNode) {
			return node, true
		}
	}

	return nil, false
}

// For now we support only one gateway node per DON
func (m *DonMetadata) Gateway() (*NodeMetadata, bool) {
	for _, node := range m.NodesMetadata {
		if slices.Contains(node.Roles, GatewayNode) {
			return node, true
		}
	}
	return nil, false
}

func (m *DonMetadata) HasFlag(flag CapabilityFlag) bool {
	return HasFlag(m.Flags, flag)
}

func (m *DonMetadata) CapabilitiesAwareNodeSet() *CapabilitiesAwareNodeSet {
	return &m.ns
}

func (m *DonMetadata) EVMChains() []uint64 {
	return m.ns.EVMChains()
}

func (m *DonMetadata) RequiresOCR() bool {
	return slices.Contains(m.Flags, ConsensusCapability) || slices.Contains(m.Flags, ConsensusCapabilityV2) ||
		slices.Contains(m.Flags, VaultCapability) || slices.Contains(m.Flags, EVMCapability)
}

func (m *DonMetadata) RequiresGateway() bool {
	return m.gh.RequiresGateway(m.Flags)
}

func (m *DonMetadata) RequiresWebAPI() bool {
	return m.gh.RequiresWebAPI(m.Flags)
}

func (m *DonMetadata) IsWorkflowDON() bool {
	// is there a case where flags are not set yet?
	if len(m.Flags) == 0 && len(m.ns.DONTypes) != 0 {
		return slices.Contains(m.ns.DONTypes, WorkflowDON)
	}

	return slices.Contains(m.Flags, WorkflowDON)
}

type Dons struct {
	Dons []*DON `toml:"dons" json:"dons"`
}

func (d *Dons) List() []*DON {
	return d.Dons
}

func NewDons(dons []*DON) *Dons {
	return &Dons{
		Dons: dons,
	}
}

type DonsMetadata struct {
	dons  []*DonMetadata
	infra infra.Provider
}

func NewDonsMetadata(dons []*DonMetadata, infra infra.Provider) (*DonsMetadata, error) {
	if dons == nil {
		dons = make([]*DonMetadata, 0)
	}
	out := &DonsMetadata{
		dons:  dons,
		infra: infra,
	}
	return out, out.validate()
}

func (m DonsMetadata) DonCount() int {
	return len(m.dons)
}

func (m DonsMetadata) List() []*DonMetadata {
	return m.dons
}

func (m DonsMetadata) validate() error {
	if len(m.dons) == 0 {
		return errors.New("at least one don is required")
	}

	if m.BootstrapCount() == 0 {
		return errors.New("at least one nodeSet must have a bootstrap node")
	}

	wfDon, err := m.WorkflowDON()
	if err != nil {
		return fmt.Errorf("failed to get workflow DON: %w", err)
	}

	if _, isBootstrap := wfDon.Bootstrap(); !isBootstrap {
		return errors.New("due to the limitations of our implementation, workflow DON must always have a bootstrap node")
	}

	if m.RequiresGateway() && !m.GatewayEnabled() {
		return errors.New("at least one DON requires gateway due to its capabilities, but no DON is configured with gateway")
	}

	return nil
}

func (m DonsMetadata) BootstrapCount() int {
	count := 0
	for _, don := range m.dons {
		if _, isBootstrap := don.Bootstrap(); isBootstrap {
			count++
		}
	}
	return count
}

func (m DonsMetadata) Bootstrap() (*NodeMetadata, bool) {
	for _, don := range m.dons {
		if node, isBootstrap := don.Bootstrap(); isBootstrap {
			return node, true
		}
	}
	return nil, false
}

// WorkflowDON returns the DON with the WorkflowDON flag. Returns an error if
// there is not exactly one such DON. Currently, the WorkflowDON flag is required on exactly one DON.
func (m DonsMetadata) WorkflowDON() (*DonMetadata, error) {
	// don't use flag b/c may not be set
	for _, don := range m.dons {
		if don.IsWorkflowDON() {
			return don, nil
		}
	}
	return nil, fmt.Errorf("no dons with flag %s found", WorkflowDON)
}

func (m DonsMetadata) GatewayEnabled() bool {
	for _, don := range m.dons {
		if _, hasGateway := don.Gateway(); hasGateway {
			return true
		}
	}
	return false
}

func (m DonsMetadata) GetGatewayDON() (*DonMetadata, error) {
	for _, don := range m.dons {
		if _, hasGateway := don.Gateway(); hasGateway {
			return don, nil
		}
	}
	return nil, fmt.Errorf("no dons with flag %s found", GatewayDON)
}

func (m DonsMetadata) RequiresGateway() bool {
	for _, don := range m.dons {
		if don.RequiresGateway() {
			return true
		}
	}
	return false
}

type NodeMetadata struct {
	Keys  *secrets.NodeKeys `toml:"keys" json:"keys"`
	Host  string            `toml:"host" json:"host"`
	Roles []string          `toml:"roles" json:"roles"`
	Index int               `toml:"index" json:"index"` // hopefully we can remove it later, but for now we need it to construct urls in CRIB
}

func (n *NodeMetadata) HasRole(role string) bool {
	return slices.Contains(n.Roles, role)
}

func (n *NodeMetadata) GetHost() string {
	return n.Host
}

func (n *NodeMetadata) PeerID() string {
	return strings.TrimPrefix(n.Keys.PeerID(), "p2p_")
}

type NodeMetadataConfig struct {
	Keys  NodeKeyInput
	Host  string
	Roles []string
	Index int
}

func NewNodeMetadata(c NodeMetadataConfig) (*NodeMetadata, error) {
	keys, err := NewNodeKeys(c.Keys)
	if err != nil {
		return nil, err
	}

	return &NodeMetadata{
		Keys:  keys,
		Host:  c.Host,
		Roles: c.Roles,
		Index: c.Index,
	}, nil
}

func newNodes(cfgs []NodeMetadataConfig) ([]*NodeMetadata, error) {
	nodes := make([]*NodeMetadata, len(cfgs))

	for i := range nodes {
		node, err := NewNodeMetadata(cfgs[i])
		if err != nil {
			return nil, fmt.Errorf("failed to create node (index: %d): %w", i, err)
		}
		nodes[i] = node
	}

	return nodes, nil
}

func NewDonTopology(registryChainSelector uint64, topology *Topology, dons *Dons) *DonTopology {
	return &DonTopology{
		WorkflowDonID:          topology.WorkflowDONID,
		HomeChainSelector:      registryChainSelector,
		Dons:                   dons,
		GatewayConnectorOutput: topology.GatewayConnectorOutput,
	}
}

type DonTopology struct {
	WorkflowDonID          uint64                  `toml:"workflow_don_id" json:"workflow_don_id"`
	HomeChainSelector      uint64                  `toml:"home_chain_selector" json:"home_chain_selector"`
	Dons                   *Dons                   `toml:"dons" json:"dons"`
	GatewayConnectorOutput *GatewayConnectorOutput `toml:"gateway_connector_output" json:"gateway_connector_output"`
}

// BootstrapNode returns the the bootstrap node that should be used as the bootstrap node for P2P peering
// Currently only one bootstrap is supported.
func (t *DonTopology) Bootstrap() (*Node, bool) {
	for _, don := range t.Dons.List() {
		if node, isBootstrap := don.Bootstrap(); isBootstrap {
			return node, true
		}
	}

	return nil, false
}

func (t *DonTopology) Gateway() (*Node, bool) {
	for _, don := range t.Dons.List() {
		if node, hasGateway := don.Gateway(); hasGateway {
			return node, true
		}
	}

	return nil, false
}

func (t *DonTopology) AnyDonHasCapability(capability CapabilityFlag) bool {
	for _, don := range t.Dons.List() {
		if don.HasFlag(capability) {
			return true
		}
	}

	return false
}

// CapabilitiesAwareNodeSet is the serialized form that declares nodesets in a topology.
type CapabilitiesAwareNodeSet struct {
	*ns.Input
	Capabilities []string `toml:"capabilities"` // global capabilities that have no chain-specific configuration (like cron, web-api-target, web-api-trigger, etc.)
	DONTypes     []string `toml:"don_types"`
	// SupportedEVMChains is filter. Use EVMChains() to get the actual list of chains supported by the nodeset.
	SupportedEVMChains []uint64 `toml:"supported_evm_chains"` // chain IDs that the DON supports, empty means all chains
	// TODO separate out bootstrap as a concept rather than index
	BootstrapNodeIndex   int               `toml:"bootstrap_node_index"` // -1 -> no bootstrap, only used if the DON doesn't hae the GatewayDON flag
	GatewayNodeIndex     int               `toml:"gateway_node_index"`   // -1 -> no gateway, only used if the DON has the GatewayDON flag
	EnvVars              map[string]string `toml:"env_vars"`             // additional environment variables to be set on each node
	RawChainCapabilities any               `toml:"chain_capabilities"`
	// ChainCapabilities allows enabling capabilities per chain with optional per-chain overrides.
	// Example syntaxes accepted per capability key:
	//   evm = ["1337", "2337"]
	//   evm = { enabled_chains = ["1337", "2337"], chain_overrides = { "1337" = { ReceiverGasMinimum = 1000 } } }
	ChainCapabilities map[string]*ChainCapabilityConfig `toml:"-"`

	// CapabilityOverrides allows overriding global capability configuration per DON.
	// Example: [nodesets.capability_overrides.web-api-target] GlobalRPS = 2000.0
	CapabilityOverrides map[string]map[string]any `toml:"capability_overrides"`

	SupportedSolChains []string `toml:"supported_sol_chains"` // sol chain IDs that the DON supports
	// Merged list of global and chain-specific capabilities. The latter ones are transformed to the format "capability-chainID", e.g. "evm-1337" for the evm capability on chain 1337.
	ComputedCapabilities []string `toml:"computed_capabilities"`
}

func (c *CapabilitiesAwareNodeSet) Flags() []string {
	var stringCaps []string

	return append(stringCaps, append(c.ComputedCapabilities, c.DONTypes...)...)
}

// EVMChains returns the list of EVM chain IDs that the nodeset supports. If SupportedChains is set, it is returned directly.
// Otherwise, the chain IDs are computed from the ChainCapabilities map by collecting all EnabledChains from each capability.
// The returned list is deduplicated and sorted.
func (c *CapabilitiesAwareNodeSet) EVMChains() []uint64 {
	if len(c.SupportedEVMChains) != 0 {
		return c.SupportedEVMChains
	}

	t := make(map[uint64]struct{})
	for _, cc := range c.ChainCapabilities {
		if cc != nil {
			for _, chainID := range cc.EnabledChains {
				t[chainID] = struct{}{}
			}
		}
	}

	// deduplicate
	out := []uint64{}
	for chainID := range t {
		out = append(out, chainID)
	}
	slices.Sort(out)

	return out
}

type CapabilitiesPeeringData struct {
	GlobalBootstraperPeerID string `toml:"global_bootstraper_peer_id" json:"global_bootstraper_peer_id"`
	GlobalBootstraperHost   string `toml:"global_bootstraper_host" json:"global_bootstraper_host"`
	Port                    int    `toml:"port" json:"port"`
}

type OCRPeeringData struct {
	OCRBootstraperPeerID string `toml:"ocr_bootstraper_peer_id" json:"ocr_bootstraper_peer_id"`
	OCRBootstraperHost   string `toml:"ocr_bootstraper_host" json:"ocr_bootstraper_host"`
	Port                 int    `toml:"port" json:"port"`
}

// ChainCapabilityConfig is a universal, static envelope for per-capability configuration.
// It supports both simple and complex TOML syntaxes via UnmarshalTOML:
// - capability = ["1337", "2337"]
// - capability = { enabled_chains=["1337","2337"], chain_overrides={"1337"={ ... }} }
type ChainCapabilityConfig struct {
	EnabledChains  []uint64                  `toml:"-"`
	ChainOverrides map[uint64]map[string]any `toml:"-"`
}

// ParseChainCapabilities parses chain_capabilities from raw TOML data and sets it on the CapabilitiesAwareNodeSet.
// This allows us to handle the flexible chain_capabilities syntax without a complex custom unmarshaler.
func (c *CapabilitiesAwareNodeSet) ParseChainCapabilities() error {
	c.ChainCapabilities = make(map[string]*ChainCapabilityConfig)
	c.ComputedCapabilities = append(c.ComputedCapabilities, c.Capabilities...)

	if c.RawChainCapabilities == nil {
		return nil
	}

	capMap, ok := c.RawChainCapabilities.(map[string]any)
	if !ok {
		return fmt.Errorf("chain_capabilities must be a map, but got %T", c.RawChainCapabilities)
	}

	parseChainID := func(v any) (uint64, error) {
		var chainID uint64
		var err error

		switch t := v.(type) {
		case string:
			trimmed := strings.TrimSpace(t)
			if trimmed == "" {
				return 0, errors.New("chain id cannot be empty")
			}
			chainID, err = strconv.ParseUint(trimmed, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid chain id string '%s': %w", trimmed, err)
			}
		case int64:
			if t < 0 {
				return 0, fmt.Errorf("chain id cannot be negative: %d", t)
			}
			chainID = uint64(t)
		case int:
			if t < 0 {
				return 0, fmt.Errorf("chain id cannot be negative: %d", t)
			}
			chainID = uint64(t)
		case uint64:
			chainID = t
		default:
			return 0, fmt.Errorf("invalid chain id type: %T. Supported types are string, int64, int, uint64", v)
		}

		if chainID == 0 {
			return 0, errors.New("chain id cannot be zero")
		}

		return chainID, nil
	}

	for capName, capValue := range capMap {
		config := &ChainCapabilityConfig{}
		computedCapabilities := []string{}

		switch v := capValue.(type) {
		case []any:
			// Handle array syntax: capability = ["1337", "2337"]
			for _, chainIDVal := range v {
				chainID, err := parseChainID(chainIDVal)
				if err != nil {
					return errors.Wrapf(err, "invalid chain ID in %s", capName)
				}
				config.EnabledChains = append(config.EnabledChains, chainID)
				computedCapabilities = append(computedCapabilities, capName+"-"+strconv.FormatUint(chainID, 10))
			}
		case map[string]any:
			// Handle map syntax: capability = { enabled_chains = [...], chain_overrides = {...} }
			if enabledChainsVal, ok := v["enabled_chains"]; ok {
				enabledChains, ok := enabledChainsVal.([]any)
				if !ok {
					return fmt.Errorf("enabled_chains must be an array in %s", capName)
				}
				for _, chainIDVal := range enabledChains {
					chainID, err := parseChainID(chainIDVal)
					if err != nil {
						return errors.Wrapf(err, "invalid chain ID in %s.enabled_chains", capName)
					}
					config.EnabledChains = append(config.EnabledChains, chainID)
					computedCapabilities = append(computedCapabilities, capName+"-"+strconv.FormatUint(chainID, 10))
				}
			}

			if chainOverridesVal, ok := v["chain_overrides"]; ok {
				chainOverrides, ok := chainOverridesVal.(map[string]any)
				if !ok {
					return errors.Errorf("chain_overrides must be a map in %s", capName)
				}
				config.ChainOverrides = make(map[uint64]map[string]any)
				for chainIDStr, overrides := range chainOverrides {
					chainID, err := strconv.ParseUint(chainIDStr, 10, 64)
					if err != nil {
						return errors.Wrapf(err, "invalid chain ID key %s in %s.chain_overrides", chainIDStr, capName)
					}

					if _, ok := overrides.(map[string]any); !ok {
						return errors.Errorf("chain override for %d in %s must be a map", chainID, capName)
					}
					config.ChainOverrides[chainID] = overrides.(map[string]any)
				}
			}
		default:
			return fmt.Errorf("unsupported chain capability format for %s: %T", capName, capValue)
		}

		c.ChainCapabilities[capName] = config
		c.ComputedCapabilities = append(c.ComputedCapabilities, computedCapabilities...)
	}

	return nil
}

func (c *CapabilitiesAwareNodeSet) ValidateChainCapabilities(bcInput []blockchain.Input) error {
	knownChains := []uint64{}
	for _, bc := range bcInput {
		if bc.Type == blockchain.FamilySolana {
			continue
		}
		chainIDUint64, convErr := strconv.ParseUint(bc.ChainID, 10, 64)
		if convErr != nil {
			return errors.Wrapf(convErr, "failed to convert chain ID %s to uint64", bc.ChainID)
		}
		knownChains = append(knownChains, chainIDUint64)
	}

	for capName, chain := range c.ChainCapabilities {
		for _, chainID := range chain.EnabledChains {
			if !slices.Contains(knownChains, chainID) {
				return fmt.Errorf("capability %s is enabled for chain %d, but chain %d is not present in the environment. Make sure you have added it to '[[blockchains]] table'", capName, chainID, chainID)
			}
		}
	}

	return nil
}

// MaxFaultyNodes returns the maximum number of faulty (Byzantine) nodes
// that a network of `n` total nodes can tolerate while still maintaining
// consensus safety under the standard BFT assumption (n >= 3f + 1).
//
// For example, with 4 nodes, at most 1 can be faulty.
// With 7 nodes, at most 2 can be faulty.
func (c *CapabilitiesAwareNodeSet) MaxFaultyNodes() (uint32, error) {
	if c.Nodes <= 0 {
		return 0, fmt.Errorf("total nodes must be greater than 0, got %d", c.Nodes)
	}
	return uint32((c.Nodes - 1) / 3), nil //nolint:gosec // disable G115
}

type NodeKeyInput struct {
	EVMChainIDs    []uint64
	SolanaChainIDs []string
	Password       string

	ImportedSecrets string // raw JSON string of secrets to import (usually from a previous run)
}

func NewNodeKeys(input NodeKeyInput) (*secrets.NodeKeys, error) {
	out := &secrets.NodeKeys{
		EVM:    make(map[uint64]*crypto.EVMKey),
		Solana: make(map[string]*crypto.SolKey),
	}

	if input.ImportedSecrets != "" {
		importedKeys, err := secrets.ImportNodeKeys(input.ImportedSecrets)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse imported secrets")
		}

		return importedKeys, nil
	}

	p2pKey, err := crypto.NewP2PKey(input.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate P2P keys")
	}
	out.P2PKey = p2pKey

	dkgKey, err := crypto.NewDKGRecipientKey(input.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate DKG recipient keys")
	}
	out.DKGKey = dkgKey

	if len(input.EVMChainIDs) > 0 {
		for _, chainID := range input.EVMChainIDs {
			k, err := crypto.NewEVMKey(input.Password, chainID)
			if err != nil {
				return nil, fmt.Errorf("failed to generate EVM keys: %w", err)
			}
			out.EVM[chainID] = k
		}
	}

	for _, chainID := range input.SolanaChainIDs {
		k, err := crypto.NewSolKey(input.Password, chainID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate Sol keys: %w", err)
		}
		out.Solana[chainID] = k
	}
	return out, nil
}

type LinkDonsToJDInput struct {
	JDClient          *cldf_jd.JobDistributor
	BlockchainOutputs []*WrappedBlockchainOutput
	DONs              []*DON
	Topology          *Topology
	CldfEnvironment   *cldf.Environment
}

func (f *LinkDonsToJDInput) Validate() error {
	if f.JDClient == nil {
		return errors.New("jd client not set")
	}
	if len(f.BlockchainOutputs) == 0 {
		return errors.New("blockchain output not set")
	}

	var expectedSeth, expectedSols int
	for _, chain := range f.BlockchainOutputs {
		if chain.SolChain != nil {
			expectedSols++
			continue
		}
		expectedSeth++
	}
	if len(f.DONs) == 0 {
		return errors.New("DONS not set")
	}
	if f.Topology == nil {
		return errors.New("topology not set")
	}
	if len(f.Topology.DonsMetadata.List()) == 0 {
		return errors.New("metadata not set")
	}
	if f.CldfEnvironment == nil {
		return errors.New("cldf environment not set")
	}

	return nil
}

type Environment struct {
	CldfEnvironment *cldf.Environment
	DonTopology     *DonTopology
}

type DeployCribDonsInput struct {
	Topology       *Topology
	NodeSetInputs  []*CapabilitiesAwareNodeSet
	CribConfigsDir string
	Namespace      string
}

func (d *DeployCribDonsInput) Validate() error {
	if d.Topology == nil {
		return errors.New("topology not set")
	}
	if len(d.Topology.DonsMetadata.List()) == 0 {
		return errors.New("metadata not set")
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
	JDInput        jd.Input
	CribConfigsDir string
	Namespace      string
}

func (d *DeployCribJdInput) Validate() error {
	if d.CribConfigsDir == "" {
		return errors.New("crib configs dir not set")
	}
	return nil
}

type DeployCribBlockchainInput struct {
	BlockchainInput *blockchain.Input
	CribConfigsDir  string
	Namespace       string
}

func (d *DeployCribBlockchainInput) Validate() error {
	if d.BlockchainInput == nil {
		return errors.New("blockchain input not set")
	}
	if d.CribConfigsDir == "" {
		return errors.New("crib configs dir not set")
	}
	if d.Namespace == "" {
		return errors.New("namespace not set")
	}
	return nil
}

type (
	CapabilityRegistryConfigFn = func(donFlags []CapabilityFlag, nodeSetInput *CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error)
	JobSpecFn                  = func(input *JobSpecInput) (DonsToJobSpecs, error)
)

type JobSpecInput struct {
	CldEnvironment            *cldf.Environment
	BlockchainOutput          *blockchain.Output
	DonTopology               *DonTopology
	InfraInput                infra.Provider
	CapabilityConfigs         map[string]CapabilityConfig
	Capabilities              []InstallableCapability
	CapabilitiesAwareNodeSets []*CapabilitiesAwareNodeSet
}

// InstallableCapability defines the interface for capabilities that can be dynamically
// registered and deployed across DONs. This interface enables plug-and-play capability
// extension without modifying core infrastructure code.
type InstallableCapability interface {
	// Flag returns the unique identifier used in TOML configurations and internal references
	Flag() CapabilityFlag

	// JobSpecFn returns a function that generates job specifications for this capability
	// based on the provided input configuration and topology. Most capabilities need this.
	// Exceptions include capabilities that are configured via the node config, like write-evm, aptos, tron or solana.
	JobSpecFn() JobSpecFn

	// NodeConfigTransformerFn returns a function to modify node-level configuration,
	// or nil if node config modification is not needed. Most capabilities don't need this.
	NodeConfigTransformerFn() NodeConfigTransformerFn

	// GatewayJobHandlerConfigFn returns a function to configure gateway handlers in the gateway jobspec,
	// or nil if no gateway handler configuration is required for this capability. Only capabilities
	// that need to connect to external resources might need this.
	GatewayJobHandlerConfigFn() GatewayHandlerConfigFn

	// CapabilityRegistryV1ConfigFn returns a function to generate capability registry
	// configuration for the v1 registry format
	CapabilityRegistryV1ConfigFn() CapabilityRegistryConfigFn

	// CapabilityRegistryV2ConfigFn returns a function to generate capability registry
	// configuration for the v2 registry format
	CapabilityRegistryV2ConfigFn() CapabilityRegistryConfigFn
}

type PersistentConfig interface {
	Load(absPath string) error
	Store(absPath string) error
}
