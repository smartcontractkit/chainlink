package types

import (
	"embed"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	proxy "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/data-feeds/generated/aggregator_proxy"
	cache "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/data-feeds/generated/data_feeds_cache"
)

type MCMSConfig struct {
	MinDelay time.Duration // delay for timelock worker to execute the transfers.
}

type AddressType string

type DeployCacheResponse struct {
	Address  common.Address
	Tx       common.Hash
	Tv       deployment.TypeAndVersion
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
	ChainsToDeploy    []uint64 // Chain Selectors
	MCMSAddressesPath string   // Path to the MCMS addresses JSON file, per chain
	InputFS           embed.FS // Filesystem to read MCMS addresses JSON file
}

type DeployProxyResponse struct {
	Address  common.Address
	Tx       common.Hash
	Tv       deployment.TypeAndVersion
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
	DataIDs          [][16]byte // without the 0x prefix
	Descriptions     []string
	WorkflowMetadata []cache.DataFeedsCacheWorkflowMetadata
	McmsConfig       *MCMSConfig
}

type RemoveFeedConfig struct {
	ChainSelector  uint64
	CacheAddress   common.Address
	ProxyAddresses []common.Address
	DataIDs        [][16]byte // without the 0x prefix
	McmsConfig     *MCMSConfig
}

type RemoveFeedConfigCSConfig struct {
	ChainSelector uint64
	CacheAddress  common.Address
	DataIDs       [][16]byte // without the 0x prefix
	McmsConfig    *MCMSConfig
}

type UpdateDataIDProxyConfig struct {
	ChainSelector  uint64
	CacheAddress   common.Address
	ProxyAddresses []common.Address
	DataIDs        [][16]byte
	McmsConfig     *MCMSConfig
}

type RemoveFeedProxyConfig struct {
	ChainSelector  uint64
	CacheAddress   common.Address
	ProxyAddresses []common.Address
	McmsConfig     *MCMSConfig
}

type ImportToAddressbookConfig struct {
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
	ContractAddress common.Address
	ChainSelector   uint64
	McmsConfig      *MCMSConfig
}

type NewFeedWithProxyConfig struct {
	ChainSelector    uint64
	AccessController common.Address
	Labels           []string // labels for AggregatorProxy
	DataID           [16]byte // without the 0x prefix
	Description      string
	WorkflowMetadata []cache.DataFeedsCacheWorkflowMetadata
	McmsConfig       *MCMSConfig
}
