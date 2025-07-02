package types

import (
	"embed"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/api"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/offchain"

	modulefeeds "github.com/smartcontractkit/chainlink-aptos/bindings/data_feeds"
	proxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/aggregator_proxy"
	bundleproxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/bundle_aggregator_proxy"
	cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
)

type MCMSConfig struct {
	MinDelay time.Duration // delay for timelock worker to execute the transfers.
}

type AddressType string

type DeployCacheResponse struct {
	Address  common.Address
	Tx       common.Hash
	Tv       cldf.TypeAndVersion
	Contract *cache.DataFeedsCache
}

type DeployConfig struct {
	ChainsToDeploy []uint64 // Chain Selectors
	Labels         []string // Labels for the cache, applies to all chains
}

type DeployAggregatorProxyConfig struct {
	ChainsToDeploy   []uint64         // Chain Selectors
	AccessController []common.Address // AccessController addresses per chain
	Labels           []string         // Labels for the cache, applies to all chains
}

type DeployBundleAggregatorProxyConfig struct {
	ChainsToDeploy []uint64 // Chain Selectors
	Owners         map[uint64]common.Address
	Labels         []string // Labels for the BundleAggregatorProxy, applies to all chains
	CacheLabel     string   // Label to find the DataFeedsCache contract address in addressbook
}

type DeployBundleAggregatorProxyResponse struct {
	Address  common.Address
	Tx       common.Hash
	Tv       cldf.TypeAndVersion
	Contract *bundleproxy.BundleAggregatorProxy
}

type DeployProxyResponse struct {
	Address  common.Address
	Tx       common.Hash
	Tv       cldf.TypeAndVersion
	Contract *proxy.AggregatorProxy
}

type SetFeedAdminConfig struct {
	ChainSelector uint64
	CacheAddress  common.Address
	AdminAddress  common.Address
	IsAdmin       bool
	McmsConfig    *MCMSConfig
}

type ProposeConfirmAggregatorConfig struct {
	ChainSelector        uint64
	ProxyAddress         common.Address
	NewAggregatorAddress common.Address
	McmsConfig           *MCMSConfig
}

type SetFeedDecimalConfig struct {
	ChainSelector    uint64
	CacheAddress     common.Address
	DataIDs          []string
	Descriptions     []string
	WorkflowMetadata []cache.DataFeedsCacheWorkflowMetadata
	McmsConfig       *MCMSConfig
}

type SetFeedBundleConfig struct {
	ChainSelector    uint64
	CacheAddress     common.Address
	DataIDs          []string
	Descriptions     []string
	DecimalsMatrix   [][]uint8
	WorkflowMetadata []cache.DataFeedsCacheWorkflowMetadata
	McmsConfig       *MCMSConfig
}

type RemoveFeedConfig struct {
	ChainSelector  uint64
	CacheAddress   common.Address
	ProxyAddresses []common.Address
	DataIDs        []string
	McmsConfig     *MCMSConfig
}

type RemoveFeedConfigCSConfig struct {
	ChainSelector uint64
	CacheAddress  common.Address
	DataIDs       []string
	McmsConfig    *MCMSConfig
}

type UpdateDataIDProxyConfig struct {
	ChainSelector  uint64
	CacheAddress   common.Address
	ProxyAddresses []common.Address
	DataIDs        []string
	McmsConfig     *MCMSConfig
}

type RemoveFeedProxyConfig struct {
	ChainSelector  uint64
	CacheAddress   common.Address
	ProxyAddresses []common.Address
	McmsConfig     *MCMSConfig
}

type ImportAddressesConfig struct {
	InputFileName string
	ChainSelector uint64
	InputFS       embed.FS
}

type MigrationConfig struct {
	InputFileName    string
	CacheAddress     common.Address
	ChainSelector    uint64
	InputFS          embed.FS
	WorkflowMetadata []cache.DataFeedsCacheWorkflowMetadata
}

type AcceptOwnershipConfig struct {
	ContractAddresses []common.Address
	ChainSelector     uint64
	McmsConfig        *MCMSConfig
}

type NewFeedWithProxyConfig struct {
	ChainSelector    uint64
	AccessController common.Address
	Labels           []string // labels for AggregatorProxy
	DataIDs          []string
	Descriptions     []string
	WorkflowMetadata []cache.DataFeedsCacheWorkflowMetadata
	McmsConfig       *MCMSConfig
}

type NodeConfig struct {
	InputFileName string
	InputFS       embed.FS
}

type WorkflowSpecConfig struct {
	TargetContractEncoderType        string `json:"targetContractEncoderType" yaml:"targetContractEncoderType"`   // Required. "data-feeds_decimal", "aptos" or "ccip"
	ConsensusAggregationMethod       string `json:"consensusAggregationMethod" yaml:"consensusAggregationMethod"` // Required. "llo_streams" or "data_feeds"
	WorkflowName                     string `json:"workflowName" yaml:"workflowName"`                             // Required
	ConsensusReportID                string `json:"consensusReportID" yaml:"consensusReportID"`                   // Required
	WriteTargetTrigger               string `json:"writeTargetTrigger" yaml:"writeTargetTrigger"`                 // Required
	ConsensusRef                     string `json:"consensusRef" yaml:"consensusRef"`                             // Default "data-feeds"
	ConsensusConfigKeyID             string `json:"consensusConfigKeyID" yaml:"consensusConfigKeyID"`             // Default "evm"
	ConsensusAllowedPartialStaleness string `json:"consensusAllowedPartialStaleness,omitempty" yaml:"consensusAllowedPartialStaleness,omitempty"`
	DeltaStageSec                    *int   `json:"deltaStageSec,omitempty" yaml:"deltaStageSec,omitempty"`     // Default 45
	TargetsSchedule                  string `json:"targetsSchedule,omitempty" yaml:"targetsSchedule,omitempty"` // Default "oneAtATime"
	TargetProcessor                  string `json:"targetProcessor,omitempty" yaml:"targetProcessor,omitempty"`
	TriggersMaxFrequencyMs           *int   `json:"triggersMaxFrequencyMs,omitempty" yaml:"triggersMaxFrequencyMs,omitempty"` // Default 5000
	CREStepTimeout                   int64  `json:"creStepTimeout,omitempty" yaml:"creStepTimeout,omitempty"`
}

type ProposeWFJobsConfig struct {
	ChainSelector      uint64
	CacheLabel         string   // Label for the DataFeedsCache contract in AB
	MigrationName      string   // Name of the migration in CLD
	InputFS            embed.FS // filesystem to read the feeds json mapping
	WorkflowJobName    string   // Required
	WorkflowSpecConfig WorkflowSpecConfig
	NodeFilter         *offchain.NodesFilter // Required. Node filter to select the nodes to send the jobs to.
}

type ProposeWFJobsV2Config struct {
	ChainSelector      uint64                `json:"chainSelector" yaml:"chainSelector"`
	CacheLabel         string                `json:"cacheLabel" yaml:"cacheLabel"`           // Label for the DataFeedsCache contract in AB, or qualifier in DataStore
	Domain             string                `json:"domain" yaml:"domain"`                   // default to data-feeds
	WorkflowJobName    string                `json:"workflowJobName" yaml:"workflowJobName"` // Required
	WorkflowSpecConfig WorkflowSpecConfig    `json:"workflowSpecConfig" yaml:"workflowSpecConfig"`
	NodeFilter         *offchain.NodesFilter `json:"nodeFilter" yaml:"nodeFilter"` // Required. Node filter to select the nodes to send the jobs to.
}

type ProposeBtJobsConfig struct {
	ChainSelector    uint64
	BootstrapJobName string
	Contract         string
	NodeFilter       *offchain.NodesFilter // Node filter to select the nodes to send the jobs to.
}

type DeleteJobsConfig struct {
	JobIDs       []string `json:"jobIDs"`       // Optional. If provided, all jobs with these IDs will be deleted.
	WorkflowName string   `json:"workflowName"` // Optional. If provided, all jobs with this workflow name will be deleted.
	Environment  string   `json:"environment"`  // Optional. If provided, the jobs will be deleted only in this environment.
}

type SetRegistryWorkflowConfig struct {
	ChainSelector         uint64
	AllowedWorkflowOwners []string
	AllowedWorkflowNames  []string
	CacheAddress          string
}

type SetRegistryFeedConfig struct {
	ChainSelector uint64
	DataIDs       []string
	Descriptions  []string
	CacheAddress  string
}

type DeployDataFeedsResponse struct {
	Address  aptos.AccountAddress
	Tx       api.Hash
	Tv       cldf.TypeAndVersion
	Contract *modulefeeds.DataFeeds
}

type DeployAptosConfig struct {
	ChainsToDeploy           []uint64 // Chain Selectors
	Labels                   []string // Data Store labels for the deployed contracts, applies to all chains
	Qualifier                string   // Data Store qualifier for the deployed contracts, applies to all chains
	OwnerAddress             string   // Owner of the deployed contracts
	PlatformAddress          string   // Address of the ChainLinkPlatform package
	SecondaryPlatformAddress string   // Secondary address of the ChainLinkPlatform package
}
