package rageping

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

type LatencyMetricsService interface {
	// Adds the given list of peers to the service. Registering the same peers multiple times is supported, upon each
	// registration, an internal reference count is incremented.
	RegisterPeers(peerIDs []ragetypes.PeerID)

	// Decremented the internal reference count for each of the given peers. Only if a reference count reaches zero,
	// the execution of the core ping/pong protocol is stopped.
	UnregisterPeers(peerIDs []ragetypes.PeerID)

	// Unregisters all peers (if any) and releases all resources.
	Close()
}

type LatencyMetricsServiceConfig struct {
	// The size of the PING message to be sent in bytes.
	// The minimal allowed value is 20 (4 bytes for the message type, 16 bytes for a random tag).
	// The size of the corresponding PONG message is constant (36 bytes: 4 bytes for message type, 32 bytes hash).
	PingSize int

	// The minimal, and maximal configured delay between two consecutive requests.
	// The actual delay used is computed uniformly at random from the interval [MinPeriod, maxPeriod].
	MinPeriod time.Duration
	MaxPeriod time.Duration

	// The maximal time to wait for a PONG message in response to a sent PING message, before considering the request to
	// be failed. Note: As long no PONG message was received for an active PING message and no Timeout was reached, the
	// protocol does not send out a new PING message.
	Timeout time.Duration

	// Extra time to wait until the first PING messages are sent out.
	// Useful, e.g., to avoid failing PING messages during testing when all nodes are started roughly at the same time.
	StartupDelay time.Duration

	// The bucket values for the prometheus.Histogram metric capturing the round-trip latencies, to be specified in
	// fractional seconds.
	//
	// Example: The bucket values [0.05, 0.10, 0.5, 1.0, 5.0] capture latencies in the following ranges:
	//  -   0 ms <= x <=  50 ms
	//  -  50 ms <  x <= 100 ms
	//  - 100 ms <  x <= 500 ms
	//  - 500 ms <  x <=   1 s
	//  -   1 s  <  x <=   5 s
	//  -   5 s  <  x <= infinity
	//
	// The value `nil` may be specified to use a set of pre-configured default bucket values.
	// See DefaultLatencyBuckets() for the default values.
	Buckets []float64
}

// Default latency histogram bucket boundaries, denoted in seconds
func DefaultLatencyBuckets() []float64 {
	return []float64{
		0.025, 0.050, 0.075, 0.100,
		0.150, 0.200, 0.250, 0.300,
		0.400, 0.500,
		0.750, 1.000,
		2.500, 5.000,
		10.000,
	}
}

// Initializes a new instance for collecting latency metrics. Metrics are collected for each passed configuration
// (PING request size, periods, ...). The passed configurations must be pairwise distinct, i.e., do not pass the same
// configuration multiple times as parameter. (This minor restriction is a result of how the underlying streams are
// initialized, and may be lifted if needed.)
func NewLatencyMetricsService(
	host *ragep2p.Host,
	registerer prometheus.Registerer,
	logger loghelper.LoggerWithContext,
	configs []*LatencyMetricsServiceConfig,
) LatencyMetricsService {
	// Create child logger to make finding rageping-related logs easier.
	logger = logger.MakeChild(commontypes.LogFields{"in": "rageping"})

	if len(configs) == 0 {
		logger.Warn("latency metrics service not starting, no configs provided", commontypes.LogFields{
			"hostPeerID": host.ID(),
		})
	}

	// Create a latencyMetricsService instance per configuration and manage all of them using a
	// latencyMetricsServiceGroup, i.e., a rapper which forwards all calls to the individual instances.
	serviceGroup := latencyMetricsServiceGroup{make([]LatencyMetricsService, 0, len(configs))}
	for _, config := range configs {
		if config.PingSize < minPingSize {
			logger.Error(
				"invalid ping size, ignoring configuration",
				commontypes.LogFields{"pingSize": config.PingSize, "minPingSize": minPingSize},
			)
			continue
		}

		serviceInstance := latencyMetricsService{
			host,
			registerer,
			logger,
			make(map[ragetypes.PeerID]*latencyMetricsPeerState),
			config,
			config.getStreamLimits(),
			sync.Mutex{},
		}
		serviceGroup.instances = append(serviceGroup.instances, &serviceInstance)
	}

	return &serviceGroup
}

// The default configurations run the protocol for two different ping sizes:
//   - a small 20 B ping, every 10-20s (10 second timeout)
//   - a larger 200 KiB ping, every 2-3min (30 second timeout)
//
// During testing, significant latency differences between the two message sizes have been found. So it makes sense to
// run the protocol for different ping sizes. However, there is no particular reason to use those exact values. Timeouts
// and startup delay are set quite conservatively. There is no need to make those tighter.
func DefaultConfigs() []*LatencyMetricsServiceConfig {
	return []*LatencyMetricsServiceConfig{
		{
			20,               // PingSize, smallest allowed size: 20 bytes
			10 * time.Second, // MinPeriod
			20 * time.Second, // MaxPeriod
			10 * time.Second, // Timeout
			30 * time.Second, // StartupDelay
			nil,              // use default latency buckets
		},
		{
			200 * 1024,       // PingSize, 200 KiB
			2 * time.Minute,  // MinPeriod
			3 * time.Minute,  // MaxPeriod
			30 * time.Second, // Timeout
			30 * time.Second, // StartupDelay
			nil,              // use default latency buckets
		},
	}
}
