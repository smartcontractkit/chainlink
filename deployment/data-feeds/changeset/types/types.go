package types

import (
	"embed"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	proxy "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/data-feeds/generated/aggregator_proxy"
	cache "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/data-feeds/generated/data_feeds_cache"
)

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
	UseMCMS       bool
}

type ProposeConfirmAggregatorConfig struct {
	ChainSelector uint64
	Proxy         common.Address
	NewAggregator common.Address
}

type SetFeedDecimalConfig struct {
	ChainSelector    uint64
	CacheAddress     common.Address
	DataIDs          [][16]byte // without the 0x prefix
	Descriptions     []string
	WorkflowMetadata []cache.DataFeedsCacheWorkflowMetadata
	UseMCMS          bool
}

type UpdateDataIDProxyConfig struct {
	ChainSelector uint64
	CacheAddress  common.Address
	Proxies       []common.Address
	DataIDs       [][16]byte
	UseMCMS       bool
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
