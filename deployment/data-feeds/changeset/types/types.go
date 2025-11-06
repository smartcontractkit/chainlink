package types

import (
	"embed"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/api"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/offchain"

	modulefeeds "github.com/smartcontractkit/chainlink-aptos/bindings/data_feeds"
	moduleplatform "github.com/smartcontractkit/chainlink-aptos/bindings/platform"
	moduleplatform_secondary "github.com/smartcontractkit/chainlink-aptos/bindings/platform_secondary"
	proxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/aggregator_proxy"
	bundleproxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/bundle_aggregator_proxy"
	cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
)

type MCMSConfig struct {
	MinDelay time.Duration `json:"minDelay" yaml:"minDelay"` // delay for timelock worker to execute the transfers.
}

type AddressType string

type DeployCacheResponse struct {
	Address  common.Address
	Tx       common.Hash
	Tv       cldf.TypeAndVersion
	Contract *cache.DataFeedsCache
}

type DeployConfig struct {
	ChainsToDeploy []uint64 `json:"chainsToDeploy" yaml:"chainsToDeploy"` // Chain Selectors
	Labels         []string `json:"labels" yaml:"labels"`                 // Labels for the cache, applies to all chains
	Qualifier      string   `json:"qualifier" yaml:"qualifier"`           // Qualifier for the contract, applies to all chains
}

type DeployAggregatorProxyConfig struct {
	ChainsToDeploy   []uint64         `json:"chainsToDeploy" yaml:"chainsToDeploy"`     // Chain Selectors
	AccessController []common.Address `json:"accessController" yaml:"accessController"` // AccessController addresses per chain
	Labels           []string         `json:"labels" yaml:"labels"`                     // Labels for the contract, applies to all chains
	Qualifier        string           `json:"qualifier" yaml:"qualifier"`               // Qualifier for the contract, applies to all chains
}

type DeployAggregatorProxyTronConfig struct {
	ChainsToDeploy   []uint64          // Chain Selectors
	AccessController []address.Address // AccessController address per chain
	Labels           []string          // Data Store labels for the deployed contracts, applies to all chains
	Qualifier        string            // Data Store qualifier for the deployed contracts, applies to all chains
	DeployOptions    *cldf_tron.DeployOptions
}

type DeployBundleAggregatorProxyConfig struct {
	ChainsToDeploy []uint64                  `json:"chainsToDeploy" yaml:"chainsToDeploy"` // Chain Selectors
	Owners         map[uint64]common.Address `json:"owners" yaml:"owners"`
	Labels         []string                  `json:"labels" yaml:"labels"`         // Labels for the BundleAggregatorProxy, applies to all chains
	CacheLabel     string                    `json:"cacheLabel" yaml:"cacheLabel"` // Label to find the DataFeedsCache contract address in addressbook
	Qualifier      string                    `json:"qualifier" yaml:"qualifier"`   // Qualifier for the contract, applies to all chains
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
	ChainSelector uint64         `json:"chainSelector" yaml:"chainSelector"`
	CacheAddress  common.Address `json:"cacheAddress" yaml:"cacheAddress"`
	AdminAddress  common.Address `json:"adminAddress" yaml:"adminAddress"`
	IsAdmin       bool           `json:"isAdmin" yaml:"isAdmin"`
	McmsConfig    *MCMSConfig    `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type SetFeedAdminTronConfig struct {
	ChainSelector  uint64
	CacheAddress   address.Address
	AdminAddress   address.Address
	IsAdmin        bool
	TriggerOptions *cldf_tron.TriggerOptions
}

type ProposeConfirmAggregatorConfig struct {
	ChainSelector        uint64         `json:"chainSelector" yaml:"chainSelector"`
	ProxyAddress         common.Address `json:"proxyAddress" yaml:"proxyAddress"`
	NewAggregatorAddress common.Address `json:"newAggregatorAddress" yaml:"newAggregatorAddress"`
	McmsConfig           *MCMSConfig    `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type SetFeedDecimalConfig struct {
	ChainSelector    uint64                                 `json:"chainSelector" yaml:"chainSelector"`
	CacheAddress     common.Address                         `json:"cacheAddress" yaml:"cacheAddress"`
	DataIDs          []string                               `json:"dataIDs" yaml:"dataIDs"`
	Descriptions     []string                               `json:"descriptions" yaml:"descriptions"`
	WorkflowMetadata []cache.DataFeedsCacheWorkflowMetadata `json:"workflowMetadata" yaml:"workflowMetadata"`
	McmsConfig       *MCMSConfig                            `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type DataFeedsCacheTronWorkflowMetadata struct {
	AllowedSender        address.Address
	AllowedWorkflowOwner address.Address
	AllowedWorkflowName  [10]byte
}

type SetFeedDecimalTronConfig struct {
	ChainSelector    uint64
	CacheAddress     address.Address
	DataIDs          []string
	Descriptions     []string
	WorkflowMetadata []DataFeedsCacheTronWorkflowMetadata
	TriggerOptions   *cldf_tron.TriggerOptions
}

type SetFeedBundleConfig struct {
	ChainSelector    uint64                                 `json:"chainSelector" yaml:"chainSelector"`
	CacheAddress     common.Address                         `json:"cacheAddress" yaml:"cacheAddress"`
	DataIDs          []string                               `json:"dataIDs" yaml:"dataIDs"`
	Descriptions     []string                               `json:"descriptions" yaml:"descriptions"`
	DecimalsMatrix   [][]uint8                              `json:"decimalsMatrix" yaml:"decimalsMatrix"`
	WorkflowMetadata []cache.DataFeedsCacheWorkflowMetadata `json:"workflowMetadata" yaml:"workflowMetadata"`
	McmsConfig       *MCMSConfig                            `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type RemoveFeedConfig struct {
	ChainSelector  uint64           `json:"chainSelector" yaml:"chainSelector"`
	CacheAddress   common.Address   `json:"cacheAddress" yaml:"cacheAddress"`
	ProxyAddresses []common.Address `json:"proxyAddresses" yaml:"proxyAddresses"`
	DataIDs        []string         `json:"dataIDs" yaml:"dataIDs"`
	McmsConfig     *MCMSConfig      `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type RemoveFeedConfigCSConfig struct {
	ChainSelector uint64         `json:"chainSelector" yaml:"chainSelector"`
	CacheAddress  common.Address `json:"cacheAddress" yaml:"cacheAddress"`
	DataIDs       []string       `json:"dataIDs" yaml:"dataIDs"`
	McmsConfig    *MCMSConfig    `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type UpdateDataIDProxyConfig struct {
	ChainSelector  uint64           `json:"chainSelector" yaml:"chainSelector"`
	CacheAddress   common.Address   `json:"cacheAddress" yaml:"cacheAddress"`
	ProxyAddresses []common.Address `json:"proxyAddresses" yaml:"proxyAddresses"`
	DataIDs        []string         `json:"dataIDs" yaml:"dataIDs"`
	McmsConfig     *MCMSConfig      `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type UpdateDataIDProxyTronConfig struct {
	ChainSelector  uint64
	CacheAddress   address.Address
	ProxyAddresses []address.Address
	DataIDs        []string
	TriggerOptions *cldf_tron.TriggerOptions
}

type RemoveFeedProxyConfig struct {
	ChainSelector  uint64           `json:"chainSelector" yaml:"chainSelector"`
	CacheAddress   common.Address   `json:"cacheAddress" yaml:"cacheAddress"`
	ProxyAddresses []common.Address `json:"proxyAddresses" yaml:"proxyAddresses"`
	McmsConfig     *MCMSConfig      `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type RemoveFeedProxyTronConfig struct {
	ChainSelector  uint64
	CacheAddress   address.Address
	ProxyAddresses []address.Address
	TriggerOptions *cldf_tron.TriggerOptions
}

type AddressSchema struct {
	Address   string                 `json:"address" yaml:"address"`
	Type      datastore.ContractType `json:"type" yaml:"type"`
	Version   string                 `json:"version" yaml:"version"`
	Qualifier string                 `json:"qualifier" yaml:"qualifier"`
	Labels    []string               `json:"labels" yaml:"labels"`
}

type ImportAddressesConfig struct {
	ChainSelector uint64           `json:"chainSelector" yaml:"chainSelector"`
	Addresses     []*AddressSchema `json:"addresses" yaml:"addresses"`
}

type MigrationSchema struct {
	Address        string              `json:"address" yaml:"address"`
	TypeAndVersion cldf.TypeAndVersion `json:"typeAndVersion" yaml:"typeAndVersion"`
	FeedID         string              `json:"feedId" yaml:"feedID"`
	Description    string              `json:"description" yaml:"description"`
}

type MigrationConfig struct {
	Proxies          []*MigrationSchema                     `json:"proxies" yaml:"proxies"`
	CacheAddress     common.Address                         `json:"cacheAddress" yaml:"cacheAddress"`
	ChainSelector    uint64                                 `json:"chainSelector" yaml:"chainSelector"`
	WorkflowMetadata []cache.DataFeedsCacheWorkflowMetadata `json:"workflowMetadata" yaml:"workflowMetadata"`
}

type AcceptOwnershipConfig struct {
	ContractAddresses []common.Address `json:"contractAddresses" yaml:"contractAddresses"`
	ChainSelector     uint64           `json:"chainSelector" yaml:"chainSelector"`
	McmsConfig        *MCMSConfig      `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type NewFeedWithProxyConfig struct {
	ChainSelector    uint64                                 `json:"chainSelector" yaml:"chainSelector"`
	AccessController common.Address                         `json:"accessController" yaml:"accessController"`
	Labels           []string                               `json:"labels" yaml:"labels"`         // labels for AggregatorProxy
	Qualifiers       []string                               `json:"qualifiers" yaml:"qualifiers"` // Qualifiers for AggregatorProxy
	DataIDs          []string                               `json:"dataIDs" yaml:"dataIDs"`
	Descriptions     []string                               `json:"descriptions" yaml:"descriptions"`
	WorkflowMetadata []cache.DataFeedsCacheWorkflowMetadata `json:"workflowMetadata" yaml:"workflowMetadata"`
	McmsConfig       *MCMSConfig                            `json:"mcmsConfig" yaml:"mcmsConfig"`
}

type NodeConfigSchema struct {
	ID           string          `json:"id" yaml:"id"`                       // node id
	Name         string          `json:"name" yaml:"name"`                   // new node name
	Labels       []*ptypes.Label `json:"labels" yaml:"labels"`               // new labels
	AppendLabels bool            `json:"append_labels" yaml:"append_labels"` // if true, append new labels to existing labels, otherwise replace
}

type UpdateNodeConfig struct {
	Nodes []*NodeConfigSchema `json:"nodes" yaml:"nodes"`
}

type MinimalNodeCfg struct {
	Name        string          `json:"name" yaml:"name"`
	CSAKey      string          `json:"csa_key" yaml:"csa_key"`
	IsBootstrap bool            `json:"is_bootstrap" yaml:"is_bootstrap"`
	Labels      []*ptypes.Label `json:"labels" yaml:"labels"`
}

type DONConfigSchema struct {
	ID    int              `json:"id" yaml:"id"`
	Name  string           `json:"name" yaml:"name"`
	Nodes []MinimalNodeCfg `json:"nodes" yaml:"nodes"`
}

type RegisterNodeConfig struct {
	DONs []*DONConfigSchema `json:"dons" yaml:"dons"` // list of DONs to register
}

type WorkflowSpecConfig struct {
	TargetContractEncoderType        string `json:"targetContractEncoderType" yaml:"targetContractEncoderType"`   // Required. "data-feeds_decimal", "aptos" or "ccip"
	ConsensusAggregationMethod       string `json:"consensusAggregationMethod" yaml:"consensusAggregationMethod"` // Required. "llo_streams" or "data_feeds"
	TriggerCapability                string `json:"triggerCapability" yaml:"triggerCapability"`                   // Required
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

// ProposeWFJobsConfig legacy type for legacy changet
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
	JobIDs       []string `json:"jobIDs" yaml:"jobIDs,omitempty"`             // Optional. If provided, all jobs with these IDs will be deleted.
	WorkflowName string   `json:"workflowName" yaml:"workflowName,omitempty"` // Optional. If provided, all jobs with this workflow name will be deleted.
	Environment  string   `json:"environment" yaml:"environment"`             // Optional. If provided, the jobs will be deleted only in this environment.
	Zone         string   `json:"zone" yaml:"zone"`                           // Optional. If provided, the jobs will be deleted only in this zone.
}

type SetRegistryWorkflowConfig struct {
	ChainSelector         uint64   `json:"chainSelector" yaml:"chainSelector"`
	AllowedWorkflowOwners []string `json:"allowedWorkflowOwners" yaml:"allowedWorkflowOwners"`
	AllowedWorkflowNames  []string `json:"allowedWorkflowNames" yaml:"allowedWorkflowNames"`
	CacheAddress          string   `json:"cacheAddress" yaml:"cacheAddress"`
}

type SetRegistryFeedConfig struct {
	ChainSelector uint64   `json:"chainSelector" yaml:"chainSelector"`
	DataIDs       []string `json:"dataIDs" yaml:"dataIDs"`
	Descriptions  []string `json:"descriptions" yaml:"descriptions"`
	CacheAddress  string   `json:"cacheAddress" yaml:"cacheAddress"`
}

type TransferDataFeedsAptosOwnershipConfig struct {
	ChainSelector    uint64 `json:"chainSelector" yaml:"chainSelector"`
	Address          string `json:"address" yaml:"address"`
	NewOwner         string `json:"NewOwner" yaml:"NewOwner"`
	TransferRegistry bool   `json:"transferRegistry" yaml:"transferRegistry"`
	TransferRouter   bool   `json:"transferRouter" yaml:"transferRouter"`
}

type AcceptDataFeedsAptosOwnershipConfig struct {
	ChainSelector  uint64 `json:"chainSelector" yaml:"chainSelector"`
	Address        string `json:"address" yaml:"address"`
	AcceptRegistry bool   `json:"acceptRegistry" yaml:"acceptRegistry"`
	AcceptRouter   bool   `json:"acceptRouter" yaml:"acceptRouter"`
}

type DeployDataFeedsResponse struct {
	Address  aptos.AccountAddress
	Tx       api.Hash
	Tv       cldf.TypeAndVersion
	Contract *modulefeeds.DataFeeds
}

type DeployPlatformResponse struct {
	Address  aptos.AccountAddress
	Tx       api.Hash
	Tv       cldf.TypeAndVersion
	Contract *moduleplatform.Platform
}

type DeployPlatformSecondaryResponse struct {
	Address  aptos.AccountAddress
	Tx       api.Hash
	Tv       cldf.TypeAndVersion
	Contract *moduleplatform_secondary.PlatformSecondary
}

type DeployAptosConfig struct {
	ChainsToDeploy           []uint64 `json:"chainsToDeploy" yaml:"chainsToDeploy"`                     // Chain Selectors
	Labels                   []string `json:"labels" yaml:"labels"`                                     // Data Store labels for the deployed contracts, applies to all chains
	Qualifier                string   `json:"qualifier" yaml:"qualifier"`                               // Data Store qualifier for the deployed contracts, applies to all chains
	OwnerAddress             string   `json:"ownerAddress" yaml:"ownerAddress"`                         // Owner of the deployed contracts
	PlatformAddress          string   `json:"platformAddress" yaml:"platformAddress"`                   // Address of the ChainLinkPlatform package
	SecondaryPlatformAddress string   `json:"secondaryPlatformAddress" yaml:"secondaryPlatformAddress"` // Secondary address of the ChainLinkPlatform package
}

type DeployTronResponse struct {
	Address address.Address
	Tx      string
	Tv      cldf.TypeAndVersion
}

type DeployTronConfig struct {
	ChainsToDeploy []uint64 // Chain Selectors
	Labels         []string // Data Store labels for the deployed contracts, applies to all chains
	Qualifier      string   // Data Store qualifier for the deployed contracts, applies to all chains
	DeployOptions  *cldf_tron.DeployOptions
}
