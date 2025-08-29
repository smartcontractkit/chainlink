package cre

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ks_sol "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/system-tests/lib/nix"

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
}

type ContractVersionsProvider interface {
	// GetContractVersions returns a map of contract name to semver
	GetContractVersions() map[string]string
}

type contractVersionsProvider struct {
	contracts map[string]string
}

func (cvp *contractVersionsProvider) GetContractVersions() map[string]string {
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
	for k, v := range overrides {
		cvp.contracts[k] = v
	}
	return cvp
}

type CapabilityFlagsProvider interface {
	SupportedCapabilityFlags() []CapabilityFlag
	GlobalCapabilityFlags() []CapabilityFlag
	ChainSpecificCapabilityFlags() []CapabilityFlag
}

func NewEnvironmentDependencies(cfp CapabilityFlagsProvider, cvp ContractVersionsProvider) *envionmentDependencies {
	return &envionmentDependencies{
		flagsProvider:       cfp,
		contractSetProvider: cvp,
	}
}

type envionmentDependencies struct {
	flagsProvider       CapabilityFlagsProvider
	contractSetProvider ContractVersionsProvider
}

func (e *envionmentDependencies) GetContractVersions() map[string]string {
	return e.contractSetProvider.GetContractVersions()
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
	InfraInput       *infra.Input
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
	NodeSets      []*CapabilitiesAwareNodeSet

	OCR3Config  keystone_changeset.OracleConfig
	OCR3Address *common.Address

	DONTimeConfig  keystone_changeset.OracleConfig
	DONTimeAddress *common.Address

	VaultOCR3Config  keystone_changeset.OracleConfig
	VaultOCR3Address *common.Address

	EVMOCR3Config    keystone_changeset.OracleConfig
	EVMOCR3Addresses *map[uint64]common.Address

	ConsensusV2OCR3Config  keystone_changeset.OracleConfig
	ConsensusV2OCR3Address *common.Address

	CapabilitiesRegistryAddress *common.Address
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
	if len(c.NodeSets) != len(c.Topology.DonsMetadata) {
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

const VaultGatewayDonID = "vault"

type GatewayConnectorDons struct {
	MembersEthAddresses []string `toml:"members_eth_addresses" json:"members_eth_addresses"`
	ID                  string   `toml:"id" json:"id"`
	Handlers            map[string]string
}
type GatewayConnectorOutput struct {
	Configurations []*GatewayConfiguration `toml:"configurations" json:"configurations"`
}

type GatewayConfiguration struct {
	Dons          []GatewayConnectorDons `toml:"dons" json:"dons"` // do not set, it will be set dynamically
	Outgoing      Outgoing               `toml:"outgoing" json:"outgoing"`
	Incoming      Incoming               `toml:"incoming" json:"incoming"`
	AuthGatewayID string                 `toml:"auth_gateway_id" json:"auth_gateway_id"`
}

type Outgoing struct {
	Host string `toml:"host" json:"host"` // do not set, it will be set dynamically
	Path string `toml:"path" json:"path"`
	Port int    `toml:"port" json:"port"`
}

type Incoming struct {
	Protocol     string `toml:"protocol" json:"protocol"` // do not set, it will be set dynamically
	Host         string `toml:"host" json:"host"`         // do not set, it will be set dynamically
	Path         string `toml:"path" json:"path"`
	InternalPort int    `toml:"internal_port" json:"internal_port"`
	ExternalPort int    `toml:"external_port" json:"external_port"`
}

type NodeConfigFn = func(input GenerateConfigsInput) (NodeIndexToConfigOverride, error)

type (
	HandlerTypeToConfig    = map[string]string
	GatewayHandlerConfigFn = func(donMetadata *DonMetadata) (HandlerTypeToConfig, error)
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
	return nil
}

type ToplogyInput struct {
	NodeSetInput    []*CapabilitiesAwareNodeSet
	DonToEthAddress map[uint32][]common.Address
}

type DonWithMetadata struct {
	DON *devenv.DON `toml:"-" json:"-"`
	*DonMetadata
}

type DonMetadata struct {
	NodesMetadata   []*NodeMetadata `toml:"nodes_metadata" json:"nodes_metadata"`
	Flags           []string        `toml:"flags" json:"flags"`
	ID              uint64          `toml:"id" json:"id"`
	Name            string          `toml:"name" json:"name"`
	SupportedChains []uint64        `toml:"supported_chains" json:"supported_chains"` // chain IDs that the DON supports, empty means all chains
}

type Label struct {
	Key   string `toml:"key" json:"key"`
	Value string `toml:"value" json:"value"`
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
	Labels []*Label `toml:"labels" json:"labels"`
}

type Topology struct {
	WorkflowDONID           uint64                  `toml:"workflow_don_id" json:"workflow_don_id"`
	HomeChainSelector       uint64                  `toml:"home_chain_selector" json:"home_chain_selector"`
	DonsMetadata            []*DonMetadata          `toml:"dons_metadata" json:"dons_metadata"`
	CapabilitiesPeeringData CapabilitiesPeeringData `toml:"capabilities_peering_data" json:"capabilities_peering_data"`
	OCRPeeringData          OCRPeeringData          `toml:"ocr_peering_data" json:"ocr_peering_data"`
	GatewayConnectorOutput  *GatewayConnectorOutput `toml:"gateway_connector_output" json:"gateway_connector_output"`
}

type DonTopology struct {
	WorkflowDonID           uint64                  `toml:"workflow_don_id" json:"workflow_don_id"`
	HomeChainSelector       uint64                  `toml:"home_chain_selector" json:"home_chain_selector"`
	CapabilitiesPeeringData CapabilitiesPeeringData `toml:"capabilities_peering_data" json:"capabilities_peering_data"`
	OCRPeeringData          OCRPeeringData          `toml:"ocr_peering_data" json:"ocr_peering_data"`
	DonsWithMetadata        []*DonWithMetadata      `toml:"dons_with_metadata" json:"dons_with_metadata"`
	GatewayConnectorOutput  *GatewayConnectorOutput `toml:"gateway_connector_output" json:"gateway_connector_output"`
}

type CapabilitiesAwareNodeSet struct {
	*ns.Input
	Capabilities         []string          `toml:"capabilities"` // global capabilities that have no chain-specific configuration (like cron, web-api-target, web-api-trigger, etc.)
	DONTypes             []string          `toml:"don_types"`
	SupportedChains      []uint64          `toml:"supported_chains"`     // chain IDs that the DON supports, empty means all chains
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

type GenerateKeysInput struct {
	GenerateEVMKeysForChainIDs []int
	GenerateSolKeysForChainIDs []string
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
type ChainIDToEVMKeys = map[int]*crypto.EVMKeys

// chainID -> SolKeys
type ChainIDToSolKeys = map[string]*crypto.SolKeys

// donID -> chainID -> EVMKeys
type DonsToEVMKeys = map[uint64]ChainIDToEVMKeys

// donID -> chainID -> SolKeys
type DonsToSolKeys = map[uint64]ChainIDToSolKeys

// donID -> P2PKeys
type DonsToP2PKeys = map[uint64]*crypto.P2PKeys

type GenerateKeysOutput struct {
	EVMKeys DonsToEVMKeys
	SolKeys DonsToSolKeys
	P2PKeys DonsToP2PKeys
}

type GenerateSecretsInput struct {
	DonMetadata *DonMetadata
	EVMKeys     ChainIDToEVMKeys
	SolKeys     ChainIDToSolKeys
	P2PKeys     *crypto.P2PKeys
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
	BlockchainOutputs map[uint64]*WrappedBlockchainOutput
	SethClients       map[uint64]*seth.Client
	SolClients        map[uint64]*solrpc.Client
	NodeSetOutput     []*WrappedNodeOutput
	ExistingAddresses cldf.AddressBook
	Datastore         datastore.DataStore
	Topology          *Topology
	OperationsBundle  operations.Bundle
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

	var expectedSeth, expectedSols int
	for _, chain := range f.BlockchainOutputs {
		if chain.SolChain != nil {
			expectedSols++
			continue
		}
		expectedSeth++
	}

	if expectedSeth != len(f.SethClients) {
		return errors.Errorf("expected '%d' got '%d' unexpected number of seth clients", expectedSeth, len(f.SethClients))
	}
	if expectedSols != len(f.SolClients) {
		return errors.Errorf("expected '%d' got '%d' unexpected number of sol clients", expectedSols, len(f.SolClients))
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
	Topology      *Topology
	NodeSetInputs []*CapabilitiesAwareNodeSet
	// todo cleanup this
	NixShell       *nix.Shell
	CribConfigsDir string
	Namespace      string
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
	JDInput *jd.Input
	// todo:  cleanup this
	NixShell       *nix.Shell
	CribConfigsDir string
	Namespace      string
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
	// todo:  cleanup this
	NixShell       *nix.Shell
	CribConfigsDir string
	Namespace      string
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
	if d.Namespace == "" {
		return errors.New("namespace not set")
	}
	return nil
}

type StartNixShellInput struct {
	InfraInput     *infra.Input
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

type (
	CapabilityRegistryConfigFn = func(donFlags []CapabilityFlag, nodeSetInput *CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error)
	JobSpecFn                  = func(input *JobSpecInput) (DonsToJobSpecs, error)
)

type JobSpecInput struct {
	CldEnvironment            *cldf.Environment
	BlockchainOutput          *blockchain.Output
	DonTopology               *DonTopology
	InfraInput                *infra.Input
	CapabilityConfigs         map[string]CapabilityConfig
	Capabilities              []InstallableCapability
	CapabilitiesAwareNodeSets []*CapabilitiesAwareNodeSet
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

	// NodeConfigFn returns a function to generate node-level configuration,
	// or nil if no node-specific config is needed. Most capabilities don't need this.
	NodeConfigFn() NodeConfigFn

	// GatewayJobHandlerConfigFn returns a function to configure gateway handlers in the gateway jobspec,
	// or nil if no gateway handler configuration is required for this capability. Only capabilities
	// that need to connect to external resources might need this.
	GatewayJobHandlerConfigFn() GatewayHandlerConfigFn

	// CapabilityRegistryV1ConfigFn returns a function to generate capability registry
	// configuration for the v1 registry format
	CapabilityRegistryV1ConfigFn() CapabilityRegistryConfigFn
}
