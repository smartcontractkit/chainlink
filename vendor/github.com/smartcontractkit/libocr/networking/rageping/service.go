package rageping

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	unsafeRand "math/rand"
)

type latencyMetricsService struct {
	host              *ragep2p.Host
	metricsRegisterer prometheus.Registerer
	logger            loghelper.LoggerWithContext
	peerStates        map[ragetypes.PeerID]*latencyMetricsPeerState
	config            *LatencyMetricsServiceConfig
	streamConfig      *latencyMetricsServiceStreamLimits

	// Mutex ensuring that the interface functions are thread safe.
	mu sync.Mutex
}

// A smaller wrapper which runs a LatencyMetricsService instead for each passed configuration.
// Calls are are simply forwarded to the internal instances.
type latencyMetricsServiceGroup struct {
	instances []LatencyMetricsService
}

var _ LatencyMetricsService = &latencyMetricsService{}
var _ LatencyMetricsService = &latencyMetricsServiceGroup{}

// Request format:
//   - message type: 1 => request/ping; 4 bytes big-endian
//   - payload: random fill bytes to achieve the desired request/response size, must contain at least 128 bits of
//     cryptographically secure randomness.
//
// Response format:
//   - message type: 2 => response/pong; 4 bytes big-endian
//   - request hash: SHA2-256 hash of the corresponding request; 32 bytes

const (
	msgTypePing uint32 = 1
	msgTypePong uint32 = 2
	minPingSize int    = 20
	pongSize    int    = 4 + 32
)

// Internal struct holding the state for each remote peer.
type latencyMetricsPeerState struct {
	// Exposed prometheus metrics; cleaned up when refCount reaches zero.
	metrics *latencyMetrics

	// Main stream used for sending/receiving PING/PONG messages.
	stream *ragep2p.Stream

	// A reference counter for the number of times this particular peer has been registered. Only a single state is
	// kept per peer (and config). Registering the same peer multiple times does not create a new ping/pong protocol
	// instance but rather only increases the reference count. Whenever it reached zero, the underlying resources
	// (metrics, and stream) are cleaned up, and this state is removed from the service's state map.
	refCount int

	// Channel indicating when the main ping/pong protocol has terminated (after a shutdown was requested).
	chDone chan struct{}
}

func (s *latencyMetricsPeerState) Done() chan struct{} {
	return s.chDone
}

// Struct holding all the configuration parameters required to set up the stream for the underlying service
// configuration. Correct values are computed from latencyMetricsServiceConfig during service initialization.
type latencyMetricsServiceStreamLimits struct {
	outgoingBufferSize int
	incomingBufferSize int
	maxMessageLength   int
	messagesLimit      ragep2p.TokenBucketParams
	bytesLimit         ragep2p.TokenBucketParams
}

func (c *LatencyMetricsServiceConfig) getStreamLimits() *latencyMetricsServiceStreamLimits {
	// We are sending and receiving ping and pong messages, so the outgoing and incoming buffers must be able to hold
	// the larger of the two message types.
	maxMessageLength := c.PingSize
	if pongSize > c.PingSize {
		maxMessageLength = pongSize
	}

	// Buffer sizes are specified in number of messages (and not in bytes). A buffer size of 2 messages is required
	// because (1) outgoing: sending a PING and sending a PONG as response to a previously received PING may happen
	// concurrently, and (2) incoming: receiving a new PING and receiving a PONG as response to a previously sent PING
	// may happen concurrently.
	outgoingBufferSize := 2
	incomingBufferSize := 2

	// There is at most one ping and one pong message received per c.minPeriod. (Only the inbound messages are
	// considered for the rate limits.)
	msgsCapacity := uint32(2 + 1 /* margin of error */)
	msgsRate := 2.0 / c.MinPeriod.Seconds()
	msgsLimit := ragep2p.TokenBucketParams{msgsRate, msgsCapacity}
	bytesCapacity := uint32((c.PingSize + pongSize) * 2)
	bytesRate := float64(bytesCapacity) / c.MinPeriod.Seconds()
	bytesLimit := ragep2p.TokenBucketParams{bytesRate, bytesCapacity}

	return &latencyMetricsServiceStreamLimits{
		outgoingBufferSize,
		incomingBufferSize,
		maxMessageLength,
		msgsLimit,
		bytesLimit,
	}
}

// Register a list of (new) peers and executes the ping-pong protocol between this host and each peer. If a peer was
// already added by a prior call to this function, it is not added again - only its reference count is incremented.
func (s *latencyMetricsService) RegisterPeers(peerIDs []ragetypes.PeerID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Trace("latencyMetricsService.RegisterPeers", commontypes.LogFields{"remotePeerIDs": peerIDs})

	for _, peerID := range peerIDs {
		// Ignore registration request for the host itself.
		if s.host.ID() == peerID {
			continue
		}

		if peerState, keyExists := s.peerStates[peerID]; keyExists {
			// A ping/pong instance is already running for this peer, therefore only increment reference count.
			peerState.refCount += 1
			continue
		}

		// At this point in the code, we know that no ping/pong instance is running for the given peer.
		// Initialize and start a new instance below.
		s.logger.Info("initializing rageping instance", commontypes.LogFields{"remotePeerID": peerID})

		stream, err := s.initStream(peerID)
		if err != nil {
			s.logger.Error(
				"initializing rageping instance failed (initStream call failed)",
				commontypes.LogFields{"error": err, "remotePeerID": peerID},
			)
			continue
		}

		metrics := newLatencyMetrics(s.metricsRegisterer, s.logger, s.host.ID(), peerID, s.config)
		refCount := 1
		peerState := &latencyMetricsPeerState{metrics, stream, refCount, make(chan struct{})}
		s.peerStates[peerID] = peerState

		go s.run(peerID, peerState)
		s.logger.Info("initializing rageping instance completed", commontypes.LogFields{"remotePeerID": peerID})
	}
}

// Unregister a list of peers. If we only have a single registration for a particular peer, the underlying
// resources are freed. Otherwise, only the reference count is decremented, but the ping/pong protocol continues to
// execute until UnregisterPeers is called once for each call to RegisterPeers.
func (s *latencyMetricsService) UnregisterPeers(peerIDs []ragetypes.PeerID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Trace("latencyMetricsService.UnregisterPeers", commontypes.LogFields{"remotePeerIDs": peerIDs})

	for _, peerID := range peerIDs {
		// Ignore unregistration request for the host itself, registration is ignored already, so there is no
		// possibility that the host itself was registered, and therefore it cannot (and does not need to be)
		// unregistered.
		if s.host.ID() == peerID {
			continue
		}

		peerState, keyExists := s.peerStates[peerID]
		if !keyExists {
			s.logger.Error(
				"failed to unregister peer from latency metrics service (not registered)",
				commontypes.LogFields{
					"remotePeerID": peerID.String(),
				},
			)
			continue
		}

		// Decrement refCount and check if we need to cleanup resources or if other registrations prevent that.
		peerState.refCount -= 1
		if peerState.refCount > 0 {
			continue
		}

		// Here, refCount reached zero. Therefore, all resources are cleaned up below.

		// Close the underlying stream, this will cause the primary service loop to exit.
		if err := peerState.stream.Close(); err != nil {
			s.logger.Warn("failed to close stream", commontypes.LogFields{"error": err})
		}

		// Wait for the primary loop to be shutdown.
		<-peerState.Done()

		peerState.metrics.Close()

		delete(s.peerStates, peerID)
	}
}

func (s *latencyMetricsService) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, peerState := range s.peerStates {
		// Close the underlying stream, this will cause the primary service loop to exit.
		if err := peerState.stream.Close(); err != nil {
			s.logger.Warn("failed to close stream", commontypes.LogFields{"error": err})
		}
		<-peerState.Done()
		peerState.metrics.Close()
	}

	// Clear all peerStates.
	s.peerStates = make(map[ragetypes.PeerID]*latencyMetricsPeerState)
}

// Forward the RegisterPeers to each underlying service instance.
func (sg *latencyMetricsServiceGroup) RegisterPeers(peerIDs []ragetypes.PeerID) {
	for _, instance := range sg.instances {
		instance.RegisterPeers(peerIDs)
	}
}

// Forward the UnregisterPeers call to each underlying service instance.
func (sg *latencyMetricsServiceGroup) UnregisterPeers(peerIDs []ragetypes.PeerID) {
	for _, instance := range sg.instances {
		instance.UnregisterPeers(peerIDs)
	}
}

// Forward the Close call to each underlying service instance.
func (sg *latencyMetricsServiceGroup) Close() {
	for _, instance := range sg.instances {
		instance.Close()
	}
}

func (s *latencyMetricsService) initStream(peerID ragetypes.PeerID) (*ragep2p.Stream, error) {
	// Get a unique stream name for each configuration.
	streamName := fmt.Sprintf(
		"ping-pong-(%v|%v|%v|%v)", s.config.PingSize, s.config.MinPeriod, s.config.MaxPeriod, s.config.Timeout,
	)

	return s.host.NewStream(
		peerID,
		streamName,
		s.streamConfig.outgoingBufferSize,
		s.streamConfig.incomingBufferSize,
		s.streamConfig.maxMessageLength,
		s.streamConfig.messagesLimit,
		s.streamConfig.bytesLimit,
	)
}

// Return a uniformly-at-random selected duration from the interval [s.config.MinPeriod, s.config.MaxPeriod].
func (s *latencyMetricsService) getNextDelay() time.Duration {
	// Get a random value in the range 0.0 to 1.0.
	r := unsafeRand.Float64()

	// Scale the value of r to size of the interval [0, maxPeriod-minPeriod].
	r *= float64(s.config.MaxPeriod - s.config.MinPeriod)

	// Shift the value of r to interval [minPeriod, maxPeriod].
	return time.Duration(r) + s.config.MinPeriod

	// The above uniform distribution should be simple and sufficient for our purposes.
	// If we really want to get fancy, something like the following would be possible:
	//
	// def get_random_delay(MIN, AVG, MAX):
	//     # Calculate alpha and beta based on the desired average
	//         alpha = AVG - MIN
	//         beta = MAX - AVG
	//
	//     # Scaling factor to adjust alpha and beta for better shape handling
	//     scale = 0.5
	//     alpha *= scale
	//     beta *= scale
	//
	//     return MIN + (MAX - MIN) * random.betavariate(alpha, beta)
}

func (s *latencyMetricsService) preparePingMessage() ([]byte, error) {
	// Initialize the message buffer and set the pingMsg type to PING.
	pingMsg := make([]byte, s.config.PingSize)
	binary.BigEndian.PutUint32(pingMsg, msgTypePing)

	// Randomly generate a tag and extra fill bytes.
	if _, err := rand.Read(pingMsg[4:]); err != nil {
		s.logger.Error(
			"internal call to preparePingMessage() failed unexpectedly",
			commontypes.LogFields{"error": err},
		)
		return nil, err
	}

	return pingMsg, nil
}

func (s *latencyMetricsService) preparePongMessage(pingMsg []byte) ([]byte, error) {
	// Initialize the message buffer and set the pingMessage type to PONG.
	pongMsg := make([]byte, 4, pongSize)
	binary.BigEndian.PutUint32(pongMsg, msgTypePong)

	// Compute the response value as the hash over the request data and append it to the response.
	hasher := sha256.New()
	_, err := hasher.Write(pingMsg)
	if err != nil {
		s.logger.Error(
			"internal call to preparePongMessage() failed unexpectedly",
			commontypes.LogFields{"error": err},
		)
		return nil, err
	}
	pongMsg = hasher.Sum(pongMsg)

	// Here response holds:
	//  - [0: 4]  the value 1 (message type: response/pong)
	//  - [4:36]  the hash of request
	return pongMsg, nil
}

// Core ping-pong protocol between the host and the given remote peer. After an initial startup delay, the protocol,
// in some regular (but somewhat randomized) interval, sends out a ping messages to the remote peer and measures the
// round-trip-time until a response is received . When a ping is received from the remote peer, it responds with the
// corresponding pong.
func (s *latencyMetricsService) run(remotePeerID ragetypes.PeerID, peerState *latencyMetricsPeerState) {
	s.logger.Info("starting rageping instance", commontypes.LogFields{"remotePeerID": remotePeerID})
	defer func() {
		s.logger.Info("stopping rageping instance", commontypes.LogFields{"remotePeerID": remotePeerID})
		close(peerState.chDone)
	}()

	stream := peerState.stream
	metrics := peerState.metrics

	ticker := time.NewTicker(s.config.StartupDelay + s.getNextDelay())
	defer ticker.Stop()

	var lastPingSentAt time.Time
	var expectedPongMsg []byte

	for {
		select {
		case <-ticker.C:
			// Check if we are currently waiting for a PONG message.
			if expectedPongMsg == nil {
				// The ticker event triggered, but we are not waiting for a PONG message, therefore:
				//  1. Send a PING message.
				//  2. Configure the ticker such that the next tick event is triggered after the configured timeout for
				//     receiving the corresponding PONG message.
				lastPingSentAt, expectedPongMsg = s.sendPing(remotePeerID, stream, metrics)
				ticker.Reset(s.config.Timeout)
			} else {
				// The ticker event triggered, and we are currently awaiting a PONG message. So no PONG message was
				// received before the configured timeout, therefore:
				//  1. Stop waiting for the expected PONG message.
				//  2. Log this timeout and update metrics accordingly.
				//  3. Reschedule the ticker for sending a new PING message.
				expectedPongMsg = nil
				s.processTimedOutPing(remotePeerID, metrics)
				ticker.Reset(s.getNextDelay())
			}

		case msg, ok := <-stream.ReceiveMessages():
			if !ok {
				// Stream was closed, so we are shutting down.
				return
			}

			// Some message was received from the remote peer.
			//  - For an incoming (valid) PING message: respond with the corresponding PONG message.
			//  - For an incoming (valid) PONG message:
			//      1. Stop waiting for the PONG message.
			//      2. Measure latency and update metrics.
			//      3. Reschedule the ticker for sending a new PING message.
			//  - Log invalid messages.
			if len(msg) >= 4 {
				msgType := binary.BigEndian.Uint32(msg)
				if msgType == msgTypePing && len(msg) == s.config.PingSize {
					s.processIncomingPingMessage(msg, remotePeerID, stream, metrics)
					break
				}
				if msgType == msgTypePong && len(msg) == pongSize {
					if s.processIncomingPongMessage(msg, expectedPongMsg, lastPingSentAt, remotePeerID, metrics) {
						expectedPongMsg = nil
						ticker.Reset(s.getNextDelay())
					}
					break
				}
			}

			// Truncate long messages before logging them. Using minPingSize here is just some suitable value,
			// other small values are equally good.
			msgPrefix := msg
			if len(msg) > minPingSize {
				msgPrefix = msg[:minPingSize]
			}

			s.logger.Warn(
				"invalid message received",
				commontypes.LogFields{"remotePeerID": remotePeerID, "msgPrefix": msgPrefix, "msgLen": len(msg)},
			)
			metrics.invalidMessagesReceivedTotal.Inc()
		}
	}
}

func (s *latencyMetricsService) sendPing(
	remotePeerID ragetypes.PeerID, stream *ragep2p.Stream, metrics *latencyMetrics,
) (lastPingSentAt time.Time, expectedPongMsg []byte) {
	// Generate a new random PING message to be sent to the remote peer.
	pingMsg, err := s.preparePingMessage()
	if err != nil {
		return
	}

	// For the above PING message, compute the PONG message we expect the remote peer to respond with.
	expectedPongMsg, err = s.preparePongMessage(pingMsg)
	if err != nil {
		return
	}

	// Actually send the PING message and keep track of the current time to compute the latency when we
	// receive corresponding PONG message.
	lastPingSentAt = time.Now()
	stream.SendMessage(pingMsg)
	metrics.sentRequestsTotal.Inc()
	s.logger.Trace(
		"sending PING",
		commontypes.LogFields{
			"remotePeerID": remotePeerID, "msgPrefix": pingMsg[:minPingSize], "msgLen": len(pingMsg),
		},
	)
	return
}

func (s *latencyMetricsService) processTimedOutPing(remotePeerID ragetypes.PeerID, metrics *latencyMetrics) {
	// expectedPongMessage != nil
	// No PONG message for was received before the configured timeout.
	s.logger.Debug(
		"peer failed to respond to last PING request in time",
		commontypes.LogFields{"remotePeerID": remotePeerID},
	)
	metrics.timedOutRequestsTotal.Inc()
}

func (s *latencyMetricsService) processIncomingPingMessage(
	pingMsg []byte,
	remotePeerID ragetypes.PeerID,
	stream *ragep2p.Stream,
	metrics *latencyMetrics,
) {
	// Some valid PING message was received from the remote peer.
	s.logger.Trace(
		"PING received",
		commontypes.LogFields{
			"remotePeerID": remotePeerID, "msgPrefix": pingMsg[:minPingSize], "msgLen": len(pingMsg),
		},
	)
	metrics.receivedRequestsTotal.Inc()

	// Respond with the corresponding PONG message.
	pongMsg, err := s.preparePongMessage(pingMsg)
	if err != nil {
		return
	}

	s.logger.Trace("sending PONG", commontypes.LogFields{"remotePeerID": remotePeerID, "msg": pongMsg})
	stream.SendMessage(pongMsg)
}

func (s *latencyMetricsService) processIncomingPongMessage(
	pongMsg []byte,
	expectedPongMsg []byte,
	lastPingSentAt time.Time,
	remotePeerID ragetypes.PeerID,
	metrics *latencyMetrics,
) bool {
	// Some (valid or invalid) PONG message was received from the remote peer.
	if bytes.Equal(pongMsg, expectedPongMsg) {
		// The value matches the expected one, so the PONG message is valid and we compute the latency
		// and update the metric.
		latency := time.Since(lastPingSentAt)
		s.logger.Trace(
			"PONG received",
			commontypes.LogFields{
				"remotePeerID": remotePeerID, "msg": pongMsg, "latency_seconds": latency.Seconds(),
			},
		)
		metrics.roundTripLatencySeconds.Observe(latency.Seconds())
		return true
	} else {
		if expectedPongMsg != nil {
			s.logger.Debug("invalid (conflicting) PONG received. The typical cause are restarts of the underlying network connection.", commontypes.LogFields{
				"remotePeerID": remotePeerID, "msg": pongMsg, "expectedMsg": expectedPongMsg,
			})
		} else {
			s.logger.Debug("invalid (unexpected) PONG received. The typical cause are restarts of the underlying network connection.", commontypes.LogFields{
				"remotePeerID": remotePeerID, "msg": pongMsg,
			})
		}
		metrics.invalidMessagesReceivedTotal.Inc()
		return false
	}
}
