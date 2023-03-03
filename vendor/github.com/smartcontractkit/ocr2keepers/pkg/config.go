package ocr2keepers

import (
	"runtime"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	ktypes "github.com/smartcontractkit/ocr2keepers/pkg/types"
)

const (
	// DefaultCacheExpiration is the default amount of time a key can remain
	// in the cache before being eligible to be cleared
	DefaultCacheExpiration = 20 * time.Minute
	// DefaultCacheClearInterval is the default setting for the interval at
	// which the cache attempts to evict expired keys
	DefaultCacheClearInterval = 30 * time.Second
	// DefaultServiceQueueLength is the default buffer size for the RPC worker
	// queue.
	DefaultServiceQueueLength = 1000
)

var (
	// DefaultMaxServiceWorkers is the max number of workers allowed to make
	// simultaneous RPC calls. The default is based on the number of CPUs
	// available to the current process.
	DefaultMaxServiceWorkers = 10 * runtime.GOMAXPROCS(0)
)

// DelegateConfig provides a single configuration struct for all options
// to be passed to the oracle, oracle factory, and underlying plugin/services.
type DelegateConfig struct {
	BinaryNetworkEndpointFactory types.BinaryNetworkEndpointFactory
	V2Bootstrappers              []commontypes.BootstrapperLocator
	ContractConfigTracker        types.ContractConfigTracker
	ContractTransmitter          types.ContractTransmitter
	KeepersDatabase              types.Database
	Logger                       commontypes.Logger
	MonitoringEndpoint           commontypes.MonitoringEndpoint
	OffchainConfigDigester       types.OffchainConfigDigester
	OffchainKeyring              types.OffchainKeyring
	OnchainKeyring               types.OnchainKeyring
	LocalConfig                  types.LocalConfig

	// EVMClient is an EVM head subscriber
	HeadSubscriber ktypes.HeadSubscriber
	// Registry is an abstract plugin registry; can be evm based or anything else
	Registry ktypes.Registry
	// PerformLogProvider is an abstract provider of logs where upkeep performs
	// occur. This interface provides subscribe and unsubscribe methods.
	PerformLogProvider ktypes.PerformLogProvider
	// ReportEncoder is an abstract encoder for encoding reports destined for
	// transmission; can be evm based or anything else.
	ReportEncoder ktypes.ReportEncoder
	// CacheExpiration is the duration of time a cached key is available. Use
	// this value to balance memory usage and RPC calls. A new set of keys is
	// generated with every block so a good setting might come from block time
	// times number of blocks of history to support not replaying reports.
	CacheExpiration time.Duration
	// CacheEvictionInterval is a parameter for how often the cache attempts to
	// evict expired keys. This value should be short enough to ensure key
	// eviction doesn't block for too long, and long enough that it doesn't
	// cause frequent blocking.
	CacheEvictionInterval time.Duration
	// MaxServiceWorkers is the total number of go-routines allowed to make RPC
	// simultaneous calls on behalf of the sampling operation. This parameter
	// is 10x the number of available CPUs by default. The RPC calls are memory
	// heavy as opposed to CPU heavy as most of the work involves waiting on
	// network responses.
	MaxServiceWorkers int
	// ServiceQueueLength is the buffer size for the RPC service queue. Fewer
	// workers or slower RPC responses will cause this queue to build up.
	// Adding new items to the queue will block if the queue becomes full.
	ServiceQueueLength int
}
