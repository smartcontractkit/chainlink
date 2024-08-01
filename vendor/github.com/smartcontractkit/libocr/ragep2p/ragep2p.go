package ragep2p

import (
	"context"
	"crypto/ed25519"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/ragep2p/internal/knock"
	"github.com/smartcontractkit/libocr/ragep2p/internal/msgbuf"
	"github.com/smartcontractkit/libocr/ragep2p/internal/mtls"
	"github.com/smartcontractkit/libocr/ragep2p/internal/ratelimit"
	"github.com/smartcontractkit/libocr/ragep2p/internal/ratelimitedconn"
	"github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

// Maximum number of streams with another peer that can be opened on a host
const MaxStreamsPerPeer = 2_000

// Maximum stream name length
const MaxStreamNameLength = 256

// Maximum length of messages sent with ragep2p
const MaxMessageLength = 1024 * 1024 * 1024 // 1 GiB. This must be smaller than INT32_MAX

const newConnTokens = MaxStreamsPerPeer * (frameHeaderEncodedSize + MaxStreamNameLength)

// assumes we re-open every stream every ten minutes during regular operation
const controlRate = MaxStreamsPerPeer / (10.0 * 60) * (frameHeaderEncodedSize + MaxStreamNameLength)

// The 5 second value is cribbed from go standard library's tls package as of version 1.16 and later
// https://cs.opensource.google/go/go/+/master:src/crypto/tls/conn.go;drc=059a9eedf45f4909db6a24242c106be15fb27193;l=1454
const netTimeout = 5 * time.Second

type hostState uint8

const (
	_ = iota
	hostStatePending
	hostStateOpen
	hostStateClosed
)

type streamID [32]byte

var _ fmt.Stringer = streamID{}

func (s streamID) String() string {
	return hex.EncodeToString(s[:])
}

type peerStreamOpenRequest struct {
	streamID           streamID
	streamName         string
	incomingBufferSize int
	maxMessageLength   int
	messagesLimit      TokenBucketParams
	bytesLimit         TokenBucketParams
}

type peerStreamOpenResponse struct {
	chSendOnOff <-chan bool
	demux       *demuxer
	err         error
}

type peerStreamCloseRequest struct {
	streamID streamID
}

type peerStreamCloseResponse struct {
	peerHasNoStreams bool
	err              error
}

type newConnNotification struct {
	chConnTerminated <-chan struct{}
}

type streamStateNotification struct {
	streamID   streamID
	streamName string // Used for sanity check, populated only on stream open and empty on stream close
	open       bool
}

type streamIDAndData struct {
	StreamID streamID
	Data     []byte
}

type peerConnLifeCycle struct {
	connCancel       context.CancelFunc
	connSubs         subprocesses.Subprocesses
	chConnTerminated <-chan struct{}
}

type peer struct {
	chDone <-chan struct{}

	other  types.PeerID
	logger loghelper.LoggerWithContext

	metrics *peerMetrics

	incomingConnsLimiterMu sync.Mutex
	incomingConnsLimiter   *ratelimit.TokenBucket

	connRateLimiter *connRateLimiter

	connLifeCycleMu sync.Mutex
	connLifeCycle   peerConnLifeCycle

	chStreamToConn chan streamIDAndData
	demuxer        *demuxer

	chNewConnNotification chan<- newConnNotification

	chOtherStreamStateNotification chan<- streamStateNotification
	chSelfStreamStateNotification  <-chan streamStateNotification

	chStreamOpenRequest  chan<- peerStreamOpenRequest
	chStreamOpenResponse <-chan peerStreamOpenResponse

	chStreamCloseRequest  chan<- peerStreamCloseRequest
	chStreamCloseResponse <-chan peerStreamCloseResponse
}

type HostConfig struct {
	// DurationBetweenDials is the minimum duration between two dials. It is
	// not the exact duration because of jitter.
	DurationBetweenDials time.Duration
}

// A Host allows users to establish Streams with other peers identified by their
// PeerID. The host will transparently handle peer discovery, secure connection
// (re)establishment, multiplexing streams over the connection and rate
// limiting.
type Host struct {
	// Constructor args
	config            HostConfig
	secretKey         ed25519.PrivateKey
	listenAddresses   []string
	discoverer        Discoverer
	logger            loghelper.LoggerWithContext
	metricsRegisterer prometheus.Registerer

	hostMetrics *hostMetrics

	// Derived from secretKey
	id      types.PeerID
	tlsCert tls.Certificate

	// Host state
	stateMu sync.Mutex
	state   hostState

	// Manage various subprocesses of host
	subprocesses subprocesses.Subprocesses
	ctx          context.Context
	cancel       context.CancelFunc

	// Peers
	peersMu sync.Mutex
	peers   map[types.PeerID]*peer
}

// NewHost creates a new Host with the provided config, Ed25519 secret key,
// network listen address. A Discoverer is also provided to NewHost for
// discovering addresses of peers.
func NewHost(
	config HostConfig,
	secretKey ed25519.PrivateKey,
	listenAddresses []string,
	discoverer Discoverer,
	logger commontypes.Logger,
	metricsRegisterer prometheus.Registerer,
) (*Host, error) {
	if len(listenAddresses) == 0 {
		return nil, fmt.Errorf("no listen addresses provided")
	}

	id, err := mtls.StaticallySizedEd25519PublicKey(secretKey.Public())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Host{
		config,
		secretKey,
		listenAddresses,
		discoverer,
		// peerID might already be set to the same value if we are managed, but we don't take any chances
		loghelper.MakeRootLoggerWithContext(logger).MakeChild(commontypes.LogFields{"id": "ragep2p", "peerID": types.PeerID(id)}),
		metricsRegisterer,

		newHostMetrics(metricsRegisterer, logger, types.PeerID(id)),

		id,
		mtls.NewMinimalX509CertFromPrivateKey(secretKey),

		sync.Mutex{},
		hostStatePending,

		subprocesses.Subprocesses{},
		ctx,
		cancel,

		sync.Mutex{},
		map[types.PeerID]*peer{},
	}, nil
}

// Start listening on the network interfaces and dialling peers.
func (ho *Host) Start() error {
	succeeded := false
	defer func() {
		if !succeeded {
			ho.Close()
		}
	}()
	ho.logger.Trace("ragep2p Start()", commontypes.LogFields{"listenAddresses": ho.listenAddresses})
	ho.stateMu.Lock()
	defer ho.stateMu.Unlock()

	if ho.state != hostStatePending {
		return fmt.Errorf("cannot Start() host that has already been started")
	}
	ho.state = hostStateOpen

	ho.subprocesses.Go(func() {
		ho.dialLoop()
	})
	for _, addr := range ho.listenAddresses {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("net.Listen(%q) failed: %w", addr, err)
		}
		ho.subprocesses.Go(func() {
			ho.listenLoop(ln)
		})
	}

	err := ho.discoverer.Start(ho, ho.secretKey, ho.logger)
	if err != nil {
		return fmt.Errorf("failed to start discoverer: %w", err)
	}

	succeeded = true
	return nil
}

func remotePeerIDField(other types.PeerID) commontypes.LogFields {
	return commontypes.LogFields{
		"remotePeerID": other,
	}
}

// Caller should hold peersMu.
func (ho *Host) findOrCreatePeer(other types.PeerID) *peer {
	if _, ok := ho.peers[other]; !ok {
		logger := ho.logger.MakeChild(remotePeerIDField(other))

		metrics := newPeerMetrics(ho.metricsRegisterer, logger, ho.id, other)

		chDone := make(chan struct{})

		chConnTerminated := make(chan struct{})
		// close so that we re-dial and establish a connection
		close(chConnTerminated)

		demuxer := newDemuxer()

		chNewConnNotification := make(chan newConnNotification)

		chOtherStreamStateNotification := make(chan streamStateNotification)
		chSelfStreamStateNotification := make(chan streamStateNotification)

		chStreamOpenRequest := make(chan peerStreamOpenRequest)
		chStreamOpenResponse := make(chan peerStreamOpenResponse)

		chStreamCloseRequest := make(chan peerStreamCloseRequest)
		chStreamCloseResponse := make(chan peerStreamCloseResponse)

		incomingConnsLimiter := ratelimit.NewTokenBucket(incomingConnsRateLimit(ho.config.DurationBetweenDials), 4, true)

		connRateLimiter := newConnRateLimiter(logger)
		connRateLimiter.AddStream(TokenBucketParams{}, TokenBucketParams{controlRate, newConnTokens})
		metrics.SetConnRateLimit(connRateLimiter.TokenBucketParams())

		p := peer{
			chDone,

			other,
			logger,

			metrics,

			sync.Mutex{},
			incomingConnsLimiter,

			connRateLimiter,

			sync.Mutex{},
			peerConnLifeCycle{
				func() {},
				subprocesses.Subprocesses{},
				chConnTerminated,
			},

			make(chan streamIDAndData),
			demuxer,

			chNewConnNotification,

			chOtherStreamStateNotification,
			chSelfStreamStateNotification,

			chStreamOpenRequest,
			chStreamOpenResponse,

			chStreamCloseRequest,
			chStreamCloseResponse,
		}
		ho.peers[other] = &p

		ho.subprocesses.Go(func() {
			peerLoop(
				ho.ctx,
				chDone,
				p.connRateLimiter,
				chNewConnNotification,
				chOtherStreamStateNotification,
				chSelfStreamStateNotification,
				demuxer,
				chStreamOpenRequest,
				chStreamOpenResponse,
				chStreamCloseRequest,
				chStreamCloseResponse,
				logger,
				metrics,
			)
		})
	}
	return ho.peers[other]
}

func peerLoop(
	ctx context.Context,
	chDone chan<- struct{},
	connRateLimiter *connRateLimiter,
	chNewConnNotification <-chan newConnNotification,
	chOtherStreamStateNotification <-chan streamStateNotification,
	chSelfStreamStateNotification chan<- streamStateNotification,
	demux *demuxer,
	chStreamOpenRequest <-chan peerStreamOpenRequest,
	chStreamOpenResponse chan<- peerStreamOpenResponse,
	chStreamCloseRequest <-chan peerStreamCloseRequest,
	chStreamCloseResponse chan<- peerStreamCloseResponse,
	logger loghelper.LoggerWithContext,
	metrics *peerMetrics,
) {
	defer close(chDone)
	defer logger.Info("peerLoop exiting", nil)

	defer metrics.Close()

	type stream struct {
		name                      string
		chOnOff                   chan<- bool
		messagesLimit, bytesLimit TokenBucketParams
	}
	streams := map[streamID]stream{}
	otherStreams := map[streamID]struct{}{}

	var chConnTerminated <-chan struct{}

	pendingSelfStreamStateNotifications := map[streamID]bool{}
	var selfStreamStateNotification streamStateNotification
	var chSelfStreamStateNotificationOrNil chan<- streamStateNotification

	for {
		chSelfStreamStateNotificationOrNil = nil
		// fake loop, we only perform zero or one iteration of this
		for streamID, state := range pendingSelfStreamStateNotifications {
			chSelfStreamStateNotificationOrNil = chSelfStreamStateNotification
			selfStreamStateNotification = streamStateNotification{
				streamID,
				streams[streamID].name,
				state,
			}
			break
		}

		select {
		case chSelfStreamStateNotificationOrNil <- selfStreamStateNotification:

			delete(pendingSelfStreamStateNotifications, selfStreamStateNotification.streamID)

			// if the stream has been opened by the other end already, switch it on right away
			if _, other := otherStreams[selfStreamStateNotification.streamID]; other && selfStreamStateNotification.open {
				select {
				case streams[selfStreamStateNotification.streamID].chOnOff <- true:
				case <-ctx.Done():
				}
			}

		case notification := <-chNewConnNotification:
			logger.Trace("New connection, creating pending notifications of all streams", nil)

			connRateLimiter.AddTokens(newConnTokens)
			metrics.SetConnRateLimit(connRateLimiter.TokenBucketParams())

			chConnTerminated = notification.chConnTerminated
			for streamID := range streams {
				pendingSelfStreamStateNotifications[streamID] = true
			}

		case <-chConnTerminated:
			chConnTerminated = nil
			logger.Trace("Connection terminated, pausing all streams", nil)

			// Clear pending notifications
			pendingSelfStreamStateNotifications = map[streamID]bool{}

			// Reset streams on other side
			otherStreams = map[streamID]struct{}{}

			// Pause all streams on our side
			for _, stream := range streams {
				select {
				case stream.chOnOff <- false:
				case <-ctx.Done():
				}
			}

			logger.Trace("Connection terminated, paused all streams", nil)

		case notification := <-chOtherStreamStateNotification:
			logger.Trace("Received stream state notification", commontypes.LogFields{
				"notification": notification,
			})

			_, other := otherStreams[notification.streamID]
			if other == notification.open {
				break
			}
			if notification.open {
				otherStreams[notification.streamID] = struct{}{}
			} else {
				delete(otherStreams, notification.streamID)
			}
			if s, ok := streams[notification.streamID]; ok {
				selfStreamName := streams[notification.streamID].name
				if notification.open && selfStreamName != notification.streamName {
					logger.Warn("Name mismatch between self and other stream", commontypes.LogFields{
						"localStreamName":  selfStreamName,
						"remoteStreamName": notification.streamName,
					})
				}
				select {
				case s.chOnOff <- notification.open:
				case <-ctx.Done():
				}
			}

		case req := <-chStreamOpenRequest:
			if _, ok := streams[req.streamID]; ok {
				chStreamOpenResponse <- peerStreamOpenResponse{
					nil,
					nil,
					fmt.Errorf("stream already exists"),
				}
			} else if len(streams) >= MaxStreamsPerPeer {
				chStreamOpenResponse <- peerStreamOpenResponse{
					nil,
					nil,
					fmt.Errorf("too many streams, expected at most %d", MaxStreamsPerPeer),
				}
			} else {
				connRateLimiter.AddStream(req.messagesLimit, req.bytesLimit)
				metrics.SetConnRateLimit(connRateLimiter.TokenBucketParams())
				if !demux.AddStream(req.streamID, req.incomingBufferSize, req.maxMessageLength, req.messagesLimit, req.bytesLimit) {
					logger.Warn("Assumption violation. Failed to add already existing stream to demuxer", commontypes.LogFields{
						"streamOpenRequest": req,
					})
					// let's try to fix the problem by removing and adding the stream again
					demux.RemoveStream(req.streamID)
					demux.AddStream(req.streamID, req.incomingBufferSize, req.maxMessageLength, req.messagesLimit, req.bytesLimit)
				}
				chOnOff := make(chan bool)
				streams[req.streamID] = stream{
					req.streamName,
					chOnOff,
					req.messagesLimit,
					req.bytesLimit,
				}
				if chConnTerminated != nil {
					pendingSelfStreamStateNotifications[req.streamID] = true
				}
				chStreamOpenResponse <- peerStreamOpenResponse{
					chOnOff,
					demux,
					nil,
				}
			}

		case req := <-chStreamCloseRequest:
			if s, ok := streams[req.streamID]; ok {
				connRateLimiter.RemoveStream(s.messagesLimit, s.bytesLimit)
				metrics.SetConnRateLimit(connRateLimiter.TokenBucketParams())
				demux.RemoveStream(req.streamID)
				delete(streams, req.streamID)
				if chConnTerminated != nil {
					pendingSelfStreamStateNotifications[req.streamID] = false
				}
				chStreamCloseResponse <- peerStreamCloseResponse{
					len(streams) == 0,
					nil,
				}

				if len(streams) == 0 {
					return
				}
			} else {
				chStreamCloseResponse <- peerStreamCloseResponse{
					false,
					fmt.Errorf("stream not found"),
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

// Close stops listening on the network interface(s) and closes all active
// streams.
func (ho *Host) Close() error {
	ho.stateMu.Lock()
	defer ho.stateMu.Unlock()

	if ho.state != hostStateOpen {
		return fmt.Errorf("cannot Close() host that isn't open")
	}
	ho.logger.Info("Host closing discoverer", nil)
	err := ho.discoverer.Close()
	ho.logger.Info("Host winding down", nil)
	ho.state = hostStateClosed
	ho.cancel()
	ho.subprocesses.Wait()
	ho.hostMetrics.Close()
	ho.logger.Info("Host exiting", nil)
	if err != nil {
		return fmt.Errorf("failed to close discoverer: %w", err)
	}
	return nil
}

func (ho *Host) ID() types.PeerID {
	return ho.id
}

func (ho *Host) dialLoop() {
	type dialState struct {
		next uint
	}
	dialStates := make(map[types.PeerID]*dialState)
	for {
		var dialProcesses subprocesses.Subprocesses
		ho.peersMu.Lock()
		peers := make([]*peer, 0, len(ho.peers))
		for pid, p := range ho.peers {
			peers = append(peers, p)
			if dialStates[pid] == nil {
				dialStates[pid] = &dialState{0}
			}
		}
		// Some peers may have been discarded, garbage collect dial states
		for pid := range dialStates {
			if ho.peers[pid] == nil {
				delete(dialStates, pid)
			}
		}
		ho.peersMu.Unlock()
		for _, p := range peers {
			p := p // copy for goroutine
			ds := dialStates[p.other]
			dialProcesses.Go(func() {
				p.connLifeCycleMu.Lock()
				chConnTerminated := p.connLifeCycle.chConnTerminated
				p.connLifeCycleMu.Unlock()
				select {
				case <-chConnTerminated:
					p.logger.Debug("Dialing", nil)
				default:
					p.logger.Trace("Dial skip", nil)
					return
				}

				addresses, err := ho.discoverer.FindPeer(p.other)
				if err != nil {
					p.logger.Warn("Discoverer error", commontypes.LogFields{"error": err})
					return
				}
				if len(addresses) == 0 {
					p.logger.Warn("Discoverer found no addresses", nil)
					return
				}

				address := string(addresses[ds.next%uint(len(addresses))])

				// We used to increment this only on dial error but a connection might fail after the Dial itself has
				// succeeded (eg. this happens with self-dials where the connection is reset after the incorrect knock
				// is received). Tracking an error so far down the stack is much harder so increment every time to give
				// a fair chance to every address.
				ds.next++

				logger := p.logger.MakeChild(commontypes.LogFields{"direction": "out", "remoteAddr": address})

				dialer := net.Dialer{
					Timeout: ho.config.DurationBetweenDials,
				}
				conn, err := dialer.DialContext(ho.ctx, "tcp", address)
				if err != nil {
					logger.Warn("Dial error", commontypes.LogFields{"error": err})
					return
				}

				logger.Trace("Dial succeeded", nil)
				ho.subprocesses.Go(func() {
					ho.handleOutgoingConnection(conn, p.other, logger)
				})
			})

		}
		dialProcesses.Wait()

		select {
		//case <-time.After(5 * time.Second): // good for testing simultaneous dials, real version is on next line
		case <-time.After(ho.config.DurationBetweenDials + time.Duration(rand.Float32()*float32(ho.config.DurationBetweenDials))):
		case <-ho.ctx.Done():
			ho.logger.Trace("Host.dialLoop exiting", nil)
			return
		}
	}
}

func (ho *Host) listenLoop(ln net.Listener) {
	ho.subprocesses.Go(func() {
		<-ho.ctx.Done()
		if err := ln.Close(); err != nil {
			ho.logger.Warn("Failed to close listener", commontypes.LogFields{"error": err})
		}
	})

	for {
		conn, err := ln.Accept()
		ho.hostMetrics.inboundDialsTotal.Inc()
		if err != nil {
			ho.logger.Info("Exiting Host.listenLoop due to error while Accepting", commontypes.LogFields{"error": err})
			return
		}
		ho.subprocesses.Go(func() {
			ho.handleIncomingConnection(conn)
		})
	}
}

func (ho *Host) handleOutgoingConnection(conn net.Conn, other types.PeerID, logger loghelper.LoggerWithContext) {
	shouldClose := true
	defer func() {
		if shouldClose {
			if err := safeClose(conn); err != nil {
				logger.Warn("Failed to close outgoing connection", commontypes.LogFields{"error": err})
			}
		}
	}()

	knck := knock.BuildKnock(other, ho.id, ho.secretKey)
	if err := conn.SetWriteDeadline(time.Now().Add(netTimeout)); err != nil {
		logger.Warn("Closing connection, error during SetWriteDeadline", commontypes.LogFields{"error": err})
		return
	}
	if _, err := conn.Write(knck); err != nil {
		logger.Warn("Error while sending knock", commontypes.LogFields{"error": err})
		return
	}

	ho.peersMu.Lock()
	peer, ok := ho.peers[other]
	ho.peersMu.Unlock()
	if !ok {
		// peer must have been deleted in the time between the dial being
		// started and now
		return
	}

	shouldClose = false

	rlConn := ratelimitedconn.NewRateLimitedConn(conn, peer.connRateLimiter, logger, peer.metrics.rawconnReadBytesTotal, peer.metrics.rawconnWrittenBytesTotal)

	tlsConfig := newTLSConfig(
		ho.tlsCert,
		mtls.VerifyCertMatchesPubKey(other),
	)
	tlsConn := tls.Client(rlConn, tlsConfig)
	ho.handleConnection(false, rlConn, tlsConn, peer, logger)
}

func (ho *Host) handleIncomingConnection(conn net.Conn) {
	remoteAddrLogFields := commontypes.LogFields{"direction": "in", "remoteAddr": conn.RemoteAddr()}
	logger := ho.logger.MakeChild(remoteAddrLogFields)
	shouldClose := true
	defer func() {
		if shouldClose {
			if err := safeClose(conn); err != nil {
				logger.Warn("Failed to close incoming connection", commontypes.LogFields{"error": err})
			}
		}
	}()

	knck := make([]byte, knock.KnockSize)
	if err := conn.SetReadDeadline(time.Now().Add(netTimeout)); err != nil {
		logger.Warn("Closing connection, error during SetReadDeadline", commontypes.LogFields{"error": err})
		return
	}
	n, err := conn.Read(knck)
	if err != nil {
		logger.Warn("Error while reading knock", commontypes.LogFields{"error": err})
		return
	}
	if n != knock.KnockSize {
		logger.Warn("Knock too short", nil)
		return
	}

	other, err := knock.VerifyKnock(ho.id, knck)
	if err != nil {
		if errors.Is(err, knock.ErrFromSelfDial) {
			logger.Info("Self-dial knock, dropping connection. Someone has likely misconfigured their announce addresses.", nil)
		} else {
			logger.Warn("Invalid knock", commontypes.LogFields{"error": err})
		}
		return
	}

	ho.peersMu.Lock()
	peer, ok := ho.peers[*other]
	ho.peersMu.Unlock()
	if !ok {
		logger.Warn("Received incoming connection from an unknown peer, closing", remotePeerIDField(*other))
		return
	}
	logger = peer.logger.MakeChild(remoteAddrLogFields) // introduce remotePeerID in our logs since we now know it
	rl := peer.connRateLimiter
	rlConn := ratelimitedconn.NewRateLimitedConn(conn, rl, logger, peer.metrics.rawconnReadBytesTotal, peer.metrics.rawconnWrittenBytesTotal)

	shouldClose = false

	tlsConfig := newTLSConfig(
		ho.tlsCert,
		mtls.VerifyCertMatchesPubKey(*other),
	)
	tlsConn := tls.Server(rlConn, tlsConfig)
	ho.handleConnection(true, rlConn, tlsConn, peer, logger)
}

func (ho *Host) handleConnection(incoming bool, rlConn *ratelimitedconn.RateLimitedConn, tlsConn *tls.Conn, peer *peer, logger loghelper.LoggerWithContext) {
	shouldClose := true
	defer func() {
		if shouldClose {
			if err := safeClose(tlsConn); err != nil {
				logger.Warn("Failed to close connection", commontypes.LogFields{"error": err})
			}
		}
	}()

	// Handshake reads and write to the connection. Set a deadline to prevent tarpitting
	if err := tlsConn.SetDeadline(time.Now().Add(netTimeout)); err != nil {
		logger.Warn("Closing connection, error during SetDeadline", commontypes.LogFields{"error": err})
		return
	}
	// Perform handshake so that we know the public key
	if err := tlsConn.Handshake(); err != nil {
		logger.Warn("Closing connection, error during Handshake", commontypes.LogFields{"error": err})
		return
	}
	// Disable deadline. Whoever uses the connection next will have to set their own timeouts.
	if err := tlsConn.SetDeadline(time.Time{}); err != nil {
		logger.Warn("Closing connection, error during SetDeadline", commontypes.LogFields{"error": err})
		return
	}

	// get public key
	pubKey, err := mtls.PubKeyFromCert(tlsConn.ConnectionState().PeerCertificates[0])
	if err != nil {
		logger.Warn("Closing connection, error getting public key", commontypes.LogFields{"error": err})
		return
	}
	if peer.other != pubKey {
		logger.Warn("TLS handshake PeerID mismatch", commontypes.LogFields{
			"expected": peer.other,
			"actual":   types.PeerID(pubKey),
		})
		return
	}

	if incoming {
		peer.incomingConnsLimiterMu.Lock()
		allowed := peer.incomingConnsLimiter.RemoveTokens(1)
		peer.incomingConnsLimiterMu.Unlock()
		if !allowed {
			logger.Warn("Incoming connection rate limited", nil)
			return
		}
	}

	rlConn.EnableRateLimiting()

	logger.Info("Connection established", nil)
	peer.metrics.connEstablishedTotal.Inc()
	if incoming {
		peer.metrics.connEstablishedInboundTotal.Inc()
	}

	// the lock here ensures there is at most one active connection at any time.
	// it also prevents races on connLifeCycle.connSubs.
	peer.connLifeCycleMu.Lock()
	peer.connLifeCycle.connCancel()
	peer.connLifeCycle.connSubs.Wait()
	connCtx, connCancel := context.WithCancel(ho.ctx)
	chConnTerminated := make(chan struct{})
	peer.connLifeCycle.connCancel = connCancel
	peer.connLifeCycle.chConnTerminated = chConnTerminated
	peer.connLifeCycle.connSubs.Go(func() {
		defer connCancel()
		authenticatedConnectionLoop(
			connCtx,
			tlsConn,
			peer.chOtherStreamStateNotification,
			peer.chSelfStreamStateNotification,
			peer.demuxer,
			peer.chStreamToConn,
			chConnTerminated,
			logger,
			peer.metrics,
		)
	})
	peer.connLifeCycleMu.Unlock()

	select {
	case peer.chNewConnNotification <- newConnNotification{chConnTerminated}:
		// keep the connection
		shouldClose = false
	case <-peer.chDone:
	case <-ho.ctx.Done():
	}
}

// TokenBucketParams contains the two parameters for a token bucket rate
// limiter.
type TokenBucketParams struct {
	Rate     float64
	Capacity uint32
}

// NewStream creates a new bidirectional stream with peer other for streamName.
// It is parameterized with a maxMessageLength, the maximum size of a message in
// bytes and two parameters for rate limiting.
func (ho *Host) NewStream(
	other types.PeerID,
	streamName string,
	outgoingBufferSize int, // number of messages that fit in the outgoing buffer
	incomingBufferSize int, // number of messages that fit in the incoming buffer
	maxMessageLength int,
	messagesLimit TokenBucketParams, // rate limit for incoming messages
	bytesLimit TokenBucketParams, // rate limit for incoming messages
) (*Stream, error) {
	if other == ho.id {
		return nil, fmt.Errorf("stream with self is forbidden")
	}

	if MaxStreamNameLength < len(streamName) {
		return nil, fmt.Errorf("streamName '%v' is longer than maximum length %v", streamName, MaxStreamNameLength)
	}

	if MaxMessageLength < maxMessageLength {
		return nil, fmt.Errorf("maxMessageLength %v is greater than global MaxMessageLength %v", maxMessageLength, MaxMessageLength)
	}

	ho.peersMu.Lock()
	defer ho.peersMu.Unlock()
	p := ho.findOrCreatePeer(other)

	sid := getStreamID(ho.id, other, streamName)

	var response peerStreamOpenResponse
	select {
	// it's important that we hold peersMu here. otherwise the peer could have
	// shut down and we'd block on the send until the host is shut down
	case p.chStreamOpenRequest <- peerStreamOpenRequest{
		sid,
		streamName,
		incomingBufferSize,
		maxMessageLength,
		messagesLimit,
		bytesLimit,
	}:
		response = <-p.chStreamOpenResponse
		if response.err != nil {
			return nil, response.err
		}
	case <-ho.ctx.Done():
		return nil, fmt.Errorf("host shut down")
	}

	ctx, cancel := context.WithCancel(ho.ctx)
	streamID := getStreamID(ho.id, other, streamName)
	streamLogger := loghelper.MakeRootLoggerWithContext(p.logger).MakeChild(commontypes.LogFields{
		"streamID":   streamID,
		"streamName": streamName,
	})
	s := Stream{
		sync.Mutex{},
		false,

		streamName,
		other,
		streamID,

		outgoingBufferSize,
		ho,

		subprocesses.Subprocesses{},
		ctx,
		cancel,
		streamLogger,
		make(chan []byte),
		make(chan []byte, 5),

		p.chStreamToConn,
		response.demux,
		response.chSendOnOff,

		p.chStreamCloseRequest,
		p.chStreamCloseResponse,
	}

	s.subprocesses.Go(func() {
		s.receiveLoop()
	})
	s.subprocesses.Go(func() {
		s.sendLoop()
	})

	streamLogger.Info("NewStream succeeded", commontypes.LogFields{
		"incomingBufferSize": incomingBufferSize,
		"maxMessageLength":   maxMessageLength,
		"messagesLimit":      messagesLimit,
		"bytesLimit":         bytesLimit,
	})

	return &s, nil
}

// Stream is an over-the-network channel between two peers. Two peers may share
// multiple disjoint streams with different names. Streams are persistent and
// agnostic to the state of the connection. They completely abstract the
// underlying connection. Messages are delivered on a best effort basis.
type Stream struct {
	closedMu sync.Mutex
	closed   bool

	name     string
	other    types.PeerID
	streamID streamID

	outgoingBufferSize int

	host *Host

	subprocesses subprocesses.Subprocesses
	ctx          context.Context
	cancel       context.CancelFunc
	logger       loghelper.LoggerWithContext
	chSend       chan []byte
	chReceive    chan []byte

	chStreamToConn chan<- streamIDAndData
	demux          *demuxer
	chStreamOnOff  <-chan bool

	chStreamCloseRequest  chan<- peerStreamCloseRequest
	chStreamCloseResponse <-chan peerStreamCloseResponse
}

// Other returns the peer ID of the stream counterparty.
func (st *Stream) Other() types.PeerID {
	return st.other
}

// Name returns the name of the stream.
func (st *Stream) Name() string {
	return st.name
}

// Best effort sending of messages. May fail without returning an error.
func (st *Stream) SendMessage(data []byte) {
	select {
	case st.chSend <- data:
	case <-st.ctx.Done():
	}
}

// Best effort receiving of messages. The returned channel will be closed when
// the stream is closed. Note that this function may return the same channel
// across invocations.
func (st *Stream) ReceiveMessages() <-chan []byte {
	return st.chReceive
}

// Close the stream. This closes any channel returned by ReceiveMessages earlier.
// After close the stream cannot be reopened. If the stream is needed in the
// future it should be created again through NewStream.
// After close, any messages passed to SendMessage will be dropped.
func (st *Stream) Close() error {
	st.closedMu.Lock()
	defer st.closedMu.Unlock()
	host := st.host

	if st.closed {
		return fmt.Errorf("already closed stream")
	}

	st.logger.Info("Stream winding down", nil)

	err := func() error {
		// Grab peersMu in case the peer has no streams left and we need to
		// delete it
		host.peersMu.Lock()
		defer host.peersMu.Unlock()

		select {
		case st.chStreamCloseRequest <- peerStreamCloseRequest{st.streamID}:
			resp := <-st.chStreamCloseResponse
			if resp.err != nil {
				st.logger.Error("Unexpected error during stream Close()", commontypes.LogFields{
					"error": resp.err,
				})
				return resp.err
			}
			if resp.peerHasNoStreams {
				st.logger.Trace("Garbage collecting peer", nil)
				peer := host.peers[st.other]
				host.subprocesses.Go(func() {
					peer.connLifeCycleMu.Lock()
					defer peer.connLifeCycleMu.Unlock()
					peer.connLifeCycle.connCancel()
					peer.connLifeCycle.connSubs.Wait()
				})
				delete(host.peers, st.other)
			}
		case <-st.ctx.Done():
		}
		return nil
	}()
	if err != nil {
		return err
	}

	st.closed = true
	st.cancel()
	st.subprocesses.Wait()
	close(st.chReceive)
	st.logger.Info("Stream exiting", nil)
	return nil
}

func (st *Stream) receiveLoop() {
	chSignalMaybePending := st.demux.SignalMaybePending(st.streamID)
	chDone := st.ctx.Done()
	for {
		select {
		case <-chSignalMaybePending:
			msg, popResult := st.demux.PopMessage(st.streamID)
			switch popResult {
			case popResultEmpty:
				st.logger.Debug("Demuxer buffer is empty", nil)
			case popResultUnknownStream:
				// Closing of streams does not happen in a single step, and so
				// it could be that in the process of closing, the stream has
				// been removed from demuxer, but receiveLoop has not stopped
				// yet (but should stop soon).
				st.logger.Info("Demuxer does not know of the stream, it is likely we are in the process of closing the stream", nil)
			case popResultSuccess:
				if msg != nil {
					select {
					case st.chReceive <- msg:
					case <-chDone:
					}
				} else {
					st.logger.Error("Demuxer indicated success but we received nil msg, this should not happen", nil)
				}
			}
		case <-chDone:
			return
		}
	}
}

func (st *Stream) sendLoop() {
	var chStreamToPeerOrNil chan<- streamIDAndData
	var pending streamIDAndData
	var onOff bool
	pendingFilled := false

	ringBuffer := msgbuf.NewMessageBuffer(st.outgoingBufferSize)

	for {
		select {
		case onOff = <-st.chStreamOnOff:
			if onOff {
				if pendingFilled {
					chStreamToPeerOrNil = st.chStreamToConn
				}
				st.logger.Info("Turned on stream", nil)
			} else {
				chStreamToPeerOrNil = nil
				st.logger.Info("Turned off stream", nil)
			}

		case msg := <-st.chSend:
			if ringBuffer.Push(msg) != nil || !pendingFilled {
				pending = streamIDAndData{st.streamID, ringBuffer.Peek()}
				pendingFilled = true
				if onOff {
					chStreamToPeerOrNil = st.chStreamToConn
				}
			}

		case chStreamToPeerOrNil <- pending:
			ringBuffer.Pop()
			if p := ringBuffer.Peek(); p != nil {
				pending = streamIDAndData{st.streamID, p}
			} else {
				pendingFilled = false
				chStreamToPeerOrNil = nil
			}

		case <-st.ctx.Done():
			return
		}
	}
}

/////////////////////////////////////////////
// authenticated connection handling
//////////////////////////////////////////////

func authenticatedConnectionLoop(
	ctx context.Context,
	conn net.Conn,
	chOtherStreamStateNotification chan<- streamStateNotification,
	chSelfStreamStateNotification <-chan streamStateNotification,
	demux *demuxer,
	chWriteData <-chan streamIDAndData,
	chTerminated chan<- struct{},
	logger loghelper.LoggerWithContext,
	metrics *peerMetrics,
) {
	defer func() {
		close(chTerminated)
		logger.Info("authenticatedConnectionLoop: exited", nil)
	}()

	var subs subprocesses.Subprocesses
	defer subs.Wait()

	defer func() {
		if err := safeClose(conn); err != nil {
			logger.Warn("Failed to close connection", commontypes.LogFields{"error": err})
		}
	}()

	childCtx, childCancel := context.WithCancel(ctx)
	defer childCancel()

	chReadTerminated := make(chan struct{})
	subs.Go(func() {
		authenticatedConnectionReadLoop(
			childCtx,
			conn,
			chOtherStreamStateNotification,
			demux,
			chReadTerminated,
			logger,
			metrics,
		)
	})

	chWriteTerminated := make(chan struct{})
	subs.Go(func() {
		authenticatedConnectionWriteLoop(
			childCtx,
			conn,
			chSelfStreamStateNotification,
			chWriteData,
			chWriteTerminated,
			logger,
			metrics,
		)
	})

	select {
	case <-ctx.Done():
	case <-chReadTerminated:
	case <-chWriteTerminated:
	}

	logger.Info("authenticatedConnectionLoop: winding down", nil)
}

func authenticatedConnectionReadLoop(
	ctx context.Context,
	conn net.Conn,
	chOtherStreamStateNotification chan<- streamStateNotification,
	demux *demuxer,
	chReadTerminated chan<- struct{},
	logger loghelper.LoggerWithContext,
	metrics *peerMetrics,
) {
	defer close(chReadTerminated)

	readInternal := func(buf []byte) bool {
		_, err := io.ReadFull(conn, buf)
		if err != nil {
			logger.Warn("Error reading from connection", commontypes.LogFields{"error": err})
			return false
		}
		metrics.connReadProcessedBytesTotal.Add(float64(len(buf)))
		return true
	}

	skipInternal := func(n uint32) bool {
		r, err := io.Copy(io.Discard, io.LimitReader(conn, int64(n)))
		if err != nil || r != int64(n) {
			logger.Warn("Error reading from connection", commontypes.LogFields{"error": err})
			return false
		}
		metrics.connReadSkippedBytesTotal.Add(float64(n))
		return true
	}

	// We taper some logs to prevent an adversary from spamming our logs
	limitsExceededTaper := loghelper.LogarithmicTaper{}
	// Note that we never reset this taper. There shouldn't be many messages
	// with unknown stream id.
	unknownStreamIDTaper := loghelper.LogarithmicTaper{}

	// We keep track of stream names for logging.
	// Note that entries in this map are not checked for truthfulness, the remote
	// could lie about the stream name.
	remoteStreamNameByID := make(map[streamID]string)

	logWithHeader := func(header frameHeader) commontypes.Logger {
		return logger.MakeChild(commontypes.LogFields{
			"payloadLength":    header.PayloadLength,
			"streamID":         header.StreamID,
			"remoteStreamName": remoteStreamNameByID[header.StreamID],
		})
	}

	// We keep track of the number of open & close frames that we have received.
	openCloseFramesReceived := 0
	const maxOpenCloseFramesReceived = 2 * MaxStreamsPerPeer

	rawHeader := make([]byte, frameHeaderEncodedSize)

	for {
		if !readInternal(rawHeader) {
			return
		}

		header, err := decodeFrameHeader(rawHeader)
		if err != nil {
			logger.Warn("Error decoding header", commontypes.LogFields{"error": err})
			return
		}

		switch header.Type {
		case frameTypeOpen:
			openCloseFramesReceived++
			if header.PayloadLength == 0 || header.PayloadLength > MaxStreamNameLength {
				return
			}
			streamName := make([]byte, header.PayloadLength)
			if !readInternal(streamName) {
				return
			}
			remoteStreamNameByID[header.StreamID] = string(streamName)
			select {
			case chOtherStreamStateNotification <- streamStateNotification{
				header.StreamID,
				string(streamName),
				true,
			}:
			case <-ctx.Done():
				return
			}
		case frameTypeClose:
			openCloseFramesReceived++
			if header.PayloadLength != 0 {
				logWithHeader(header).Warn("Frame close payload length is not zero", nil)
				return
			}
			delete(remoteStreamNameByID, header.StreamID)
			select {
			case chOtherStreamStateNotification <- streamStateNotification{
				header.StreamID,
				"",
				false,
			}:
			case <-ctx.Done():
				return
			}
		case frameTypeData:
			if MaxMessageLength < header.PayloadLength {
				logWithHeader(header).Warn("authenticatedConnectionReadLoop: message exceeds ragep2p message length limit, closing connection", commontypes.LogFields{
					"payloadLength":           header.PayloadLength,
					"ragep2pMaxMessageLength": MaxMessageLength,
				})
				return
			}
			// Cast to int is safe since header.PayloadLength <= MaxMessageLength <= INT_MAX
			switch demux.ShouldPush(header.StreamID, int(header.PayloadLength)) {
			case shouldPushResultMessageTooBig:
				logWithHeader(header).Warn("authenticatedConnectionReadLoop: message too big, closing connection", commontypes.LogFields{
					"payloadLength": header.PayloadLength,
				})
				return
			case shouldPushResultMessagesLimitExceeded:
				limitsExceededTaper.Trigger(func(count uint64) {
					logWithHeader(header).Warn("authenticatedConnectionReadLoop: message limit exceeded, dropping message", commontypes.LogFields{
						"limitsExceededDroppedCount": count,
					})
				})
				if !skipInternal(header.PayloadLength) {
					return
				}
			case shouldPushResultBytesLimitExceeded:
				limitsExceededTaper.Trigger(func(count uint64) {
					logWithHeader(header).Warn("authenticatedConnectionReadLoop: bytes limit exceeded, dropping message", commontypes.LogFields{
						"limitsExceededDroppedCount": count,
					})
				})
				if !skipInternal(header.PayloadLength) {
					return
				}
			case shouldPushResultUnknownStream:
				unknownStreamIDTaper.Trigger(func(count uint64) {
					logWithHeader(header).Warn("authenticatedConnectionReadLoop: unknown stream id, dropping message", commontypes.LogFields{
						"unknownStreamIDDroppedCount": count,
					})
				})
				if !skipInternal(header.PayloadLength) {
					return
				}
			case shouldPushResultYes:
				limitsExceededTaper.Reset(func(oldCount uint64) {
					logWithHeader(header).Info("authenticatedConnectionReadLoop: limits are no longer being exceeded", commontypes.LogFields{
						"droppedCount": oldCount,
					})
				})
				data := make([]byte, header.PayloadLength)
				if !readInternal(data) {
					return
				}
				switch demux.PushMessage(header.StreamID, data) {
				case pushResultSuccess:
				case pushResultDropped:
					logWithHeader(header).Trace("authenticatedConnectionReadLoop: demuxer is overflowing for stream, dropping oldest message", nil)
				case pushResultUnknownStream:
					unknownStreamIDTaper.Trigger(func(count uint64) {
						logWithHeader(header).Warn("authenticatedConnectionReadLoop: unknown stream id, dropping message", commontypes.LogFields{
							"unknownStreamIDDroppedCount": count,
						})
					})
				}

			}
		}

		if openCloseFramesReceived > maxOpenCloseFramesReceived {
			logWithHeader(header).Warn("authenticatedConnectionReadLoop: peer received too many open/close frames, closing connection", commontypes.LogFields{
				"maxOpenCloseFramesReceived": maxOpenCloseFramesReceived,
			})
			return
		}
	}
}

func authenticatedConnectionWriteLoop(
	ctx context.Context,
	conn net.Conn,
	chSelfStreamStateNotification <-chan streamStateNotification,
	chWriteData <-chan streamIDAndData,
	chWriteTerminated chan<- struct{},
	logger loghelper.LoggerWithContext,
	metrics *peerMetrics,
) {
	writeInternal := func(buf []byte) bool {
		_, err := conn.Write(buf)
		if err != nil {
			logger.Warn("Error writing to connection", commontypes.LogFields{"error": err})
			// shut everything down
			if err := safeClose(conn); err != nil {
				logger.Warn("Failed to close connection", commontypes.LogFields{"error": err})
			}
			close(chWriteTerminated)
			return false
		}
		metrics.connWrittenBytesTotal.Add(float64(len(buf)))
		return true
	}

	for {
		select {
		case data := <-chWriteData:
			if err := conn.SetWriteDeadline(time.Now().Add(netTimeout)); err != nil {
				logger.Warn("Closing connection, error during SetWriteDeadline", commontypes.LogFields{"error": err})
				return
			}
			header := frameHeader{
				frameTypeData,
				data.StreamID,
				uint32(len(data.Data)),
			}
			if !writeInternal(header.Encode()) {
				return
			}
			if !writeInternal(data.Data) {
				return
			}
			metrics.messageBytes.Observe(float64(len(data.Data)))
		case notification := <-chSelfStreamStateNotification:
			if err := conn.SetWriteDeadline(time.Now().Add(netTimeout)); err != nil {
				logger.Warn("Closing connection, error during SetWriteDeadline", commontypes.LogFields{"error": err})
				return
			}
			var header frameHeader
			streamName := []byte(notification.streamName)
			if notification.open {
				header = frameHeader{
					frameTypeOpen,
					notification.streamID,
					uint32(len(streamName)),
				}
			} else {
				header = frameHeader{
					frameTypeClose,
					notification.streamID,
					uint32(0),
				}
			}
			if !writeInternal(header.Encode()) {
				return
			}
			if notification.open && !writeInternal(streamName) {
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// gotta be careful about closing tls connections to make sure we don't get
// tarpitted
func safeClose(conn net.Conn) error {
	// This isn't needed in more recent versions of go, but better safe than sorry!
	errDeadline := conn.SetWriteDeadline(time.Now().Add(netTimeout))
	errClose := conn.Close()
	if errClose != nil {
		return errClose
	}
	if errDeadline != nil {
		return errDeadline
	}
	return nil
}

func incomingConnsRateLimit(durationBetweenDials time.Duration) ratelimit.MillitokensPerSecond {
	// 2 dials per DurationBetweenDials are okay
	result := ratelimit.MillitokensPerSecond(2.0 / durationBetweenDials.Seconds() * 1000.0)
	// dialing once every two seconds is always okay
	if result < 500 {
		result = 500
	}
	return result
}

// Discoverer is responsible for discovering the addresses of peers on the network.
type Discoverer interface {
	Start(host *Host, privKey ed25519.PrivateKey, logger loghelper.LoggerWithContext) error
	Close() error
	FindPeer(peer types.PeerID) ([]types.Address, error)
}
