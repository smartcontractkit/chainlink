package swarm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/libp2p/go-libp2p-core/connmgr"
	"github.com/libp2p/go-libp2p-core/metrics"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/transport"

	logging "github.com/ipfs/go-log"
	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"

	ma "github.com/multiformats/go-multiaddr"
)

// DialTimeoutLocal is the maximum duration a Dial to local network address
// is allowed to take.
// This includes the time between dialing the raw network connection,
// protocol selection as well the handshake, if applicable.
var DialTimeoutLocal = 5 * time.Second

var log = logging.Logger("swarm2")

// ErrSwarmClosed is returned when one attempts to operate on a closed swarm.
var ErrSwarmClosed = errors.New("swarm closed")

// ErrAddrFiltered is returned when trying to register a connection to a
// filtered address. You shouldn't see this error unless some underlying
// transport is misbehaving.
var ErrAddrFiltered = errors.New("address filtered")

// ErrDialTimeout is returned when one a dial times out due to the global timeout
var ErrDialTimeout = errors.New("dial timed out")

// Swarm is a connection muxer, allowing connections to other peers to
// be opened and closed, while still using the same Chan for all
// communication. The Chan sends/receives Messages, which note the
// destination or source Peer.
type Swarm struct {
	// Close refcount. This allows us to fully wait for the swarm to be torn
	// down before continuing.
	refs sync.WaitGroup

	local peer.ID
	peers peerstore.Peerstore

	nextConnID   uint32 // guarded by atomic
	nextStreamID uint32 // guarded by atomic

	conns struct {
		sync.RWMutex
		m map[peer.ID][]*Conn
	}

	listeners struct {
		sync.RWMutex

		ifaceListenAddres []ma.Multiaddr
		cacheEOL          time.Time

		m map[transport.Listener]struct{}
	}

	notifs struct {
		sync.RWMutex
		m map[network.Notifiee]struct{}
	}

	transports struct {
		sync.RWMutex
		m map[int]transport.Transport
	}

	// new connection and stream handlers
	connh   atomic.Value
	streamh atomic.Value

	// dialing helpers
	dsync   *DialSync
	backf   DialBackoff
	limiter *dialLimiter
	gater   connmgr.ConnectionGater

	proc goprocess.Process
	ctx  context.Context
	bwc  metrics.Reporter
}

// NewSwarm constructs a Swarm.
//
// NOTE: go-libp2p will be moving to dependency injection soon. The variadic
// `extra` interface{} parameter facilitates the future migration. Supported
// elements are:
//  - connmgr.ConnectionGater
func NewSwarm(ctx context.Context, local peer.ID, peers peerstore.Peerstore, bwc metrics.Reporter, extra ...interface{}) *Swarm {
	s := &Swarm{
		local: local,
		peers: peers,
		bwc:   bwc,
	}

	s.conns.m = make(map[peer.ID][]*Conn)
	s.listeners.m = make(map[transport.Listener]struct{})
	s.transports.m = make(map[int]transport.Transport)
	s.notifs.m = make(map[network.Notifiee]struct{})

	for _, i := range extra {
		switch v := i.(type) {
		case connmgr.ConnectionGater:
			s.gater = v
		}
	}

	s.dsync = NewDialSync(s.doDial)
	s.limiter = newDialLimiter(s.dialAddr, s.IsFdConsumingAddr)
	s.proc = goprocessctx.WithContext(ctx)
	s.ctx = goprocessctx.OnClosingContext(s.proc)
	s.backf.init(s.ctx)

	// Set teardown after setting the context/process so we don't start the
	// teardown process early.
	s.proc.SetTeardown(s.teardown)

	return s
}

func (s *Swarm) teardown() error {
	// Wait for the context to be canceled.
	// This allows other parts of the swarm to detect that we're shutting
	// down.
	<-s.ctx.Done()

	// Prevents new connections and/or listeners from being added to the swarm.

	s.listeners.Lock()
	listeners := s.listeners.m
	s.listeners.m = nil
	s.listeners.Unlock()

	s.conns.Lock()
	conns := s.conns.m
	s.conns.m = nil
	s.conns.Unlock()

	// Lots of goroutines but we might as well do this in parallel. We want to shut down as fast as
	// possible.

	for l := range listeners {
		go func(l transport.Listener) {
			if err := l.Close(); err != nil {
				log.Errorf("error when shutting down listener: %s", err)
			}
		}(l)
	}

	for _, cs := range conns {
		for _, c := range cs {
			go func(c *Conn) {
				if err := c.Close(); err != nil {
					log.Errorf("error when shutting down connection: %s", err)
				}
			}(c)
		}
	}

	// Wait for everything to finish.
	s.refs.Wait()

	return nil
}

// Process returns the Process of the swarm
func (s *Swarm) Process() goprocess.Process {
	return s.proc
}

func (s *Swarm) addConn(tc transport.CapableConn, dir network.Direction) (*Conn, error) {
	var (
		p    = tc.RemotePeer()
		addr = tc.RemoteMultiaddr()
	)

	if s.gater != nil {
		if allow := s.gater.InterceptAddrDial(p, addr); !allow {
			err := tc.Close()
			if err != nil {
				log.Warnf("failed to close connection with peer %s and addr %s; err: %s", p.Pretty(), addr, err)
			}
			return nil, ErrAddrFiltered
		}
	}

	// Wrap and register the connection.
	stat := network.Stat{Direction: dir, Opened: time.Now()}
	c := &Conn{
		conn:  tc,
		swarm: s,
		stat:  stat,
		id:    atomic.AddUint32(&s.nextConnID, 1),
	}

	// we ONLY check upgraded connections here so we can send them a Disconnect message.
	// If we do this in the Upgrader, we will not be able to do this.
	if s.gater != nil {
		if allow, _ := s.gater.InterceptUpgraded(c); !allow {
			// TODO Send disconnect with reason here
			err := tc.Close()
			if err != nil {
				log.Warnf("failed to close connection with peer %s and addr %s; err: %s", p.Pretty(), addr, err)
			}
			return nil, ErrGaterDisallowedConnection
		}
	}

	// Add the public key.
	if pk := tc.RemotePublicKey(); pk != nil {
		s.peers.AddPubKey(p, pk)
	}

	// Clear any backoffs
	s.backf.Clear(p)

	// Finally, add the peer.
	s.conns.Lock()
	// Check if we're still online
	if s.conns.m == nil {
		s.conns.Unlock()
		tc.Close()
		return nil, ErrSwarmClosed
	}

	c.streams.m = make(map[*Stream]struct{})
	s.conns.m[p] = append(s.conns.m[p], c)

	// Add two swarm refs:
	// * One will be decremented after the close notifications fire in Conn.doClose
	// * The other will be decremented when Conn.start exits.
	s.refs.Add(2)

	// Take the notification lock before releasing the conns lock to block
	// Disconnect notifications until after the Connect notifications done.
	c.notifyLk.Lock()
	s.conns.Unlock()

	// We have a connection now. Cancel all other in-progress dials.
	// This should be fast, no reason to wait till later.
	if dir == network.DirOutbound {
		s.dsync.CancelDial(p)
	}

	s.notifyAll(func(f network.Notifiee) {
		f.Connected(s, c)
	})
	c.notifyLk.Unlock()

	c.start()

	// TODO: Get rid of this. We use it for identify but that happen much
	// earlier (really, inside the transport and, if not then, during the
	// notifications).
	if h := s.ConnHandler(); h != nil {
		go h(c)
	}

	return c, nil
}

// Peerstore returns this swarms internal Peerstore.
func (s *Swarm) Peerstore() peerstore.Peerstore {
	return s.peers
}

// Context returns the context of the swarm
func (s *Swarm) Context() context.Context {
	return s.ctx
}

// Close stops the Swarm.
func (s *Swarm) Close() error {
	return s.proc.Close()
}

// TODO: We probably don't need the conn handlers.

// SetConnHandler assigns the handler for new connections.
// You will rarely use this. See SetStreamHandler
func (s *Swarm) SetConnHandler(handler network.ConnHandler) {
	s.connh.Store(handler)
}

// ConnHandler gets the handler for new connections.
func (s *Swarm) ConnHandler() network.ConnHandler {
	handler, _ := s.connh.Load().(network.ConnHandler)
	return handler
}

// SetStreamHandler assigns the handler for new streams.
func (s *Swarm) SetStreamHandler(handler network.StreamHandler) {
	s.streamh.Store(handler)
}

// StreamHandler gets the handler for new streams.
func (s *Swarm) StreamHandler() network.StreamHandler {
	handler, _ := s.streamh.Load().(network.StreamHandler)
	return handler
}

// NewStream creates a new stream on any available connection to peer, dialing
// if necessary.
func (s *Swarm) NewStream(ctx context.Context, p peer.ID) (network.Stream, error) {
	log.Debugf("[%s] opening stream to peer [%s]", s.local, p)

	// Algorithm:
	// 1. Find the best connection, otherwise, dial.
	// 2. Try opening a stream.
	// 3. If the underlying connection is, in fact, closed, close the outer
	//    connection and try again. We do this in case we have a closed
	//    connection but don't notice it until we actually try to open a
	//    stream.
	//
	// Note: We only dial once.
	//
	// TODO: Try all connections even if we get an error opening a stream on
	// a non-closed connection.
	dials := 0
	for {
		c := s.bestConnToPeer(p)
		if c == nil {
			if nodial, _ := network.GetNoDial(ctx); nodial {
				return nil, network.ErrNoConn
			}

			if dials >= DialAttempts {
				return nil, errors.New("max dial attempts exceeded")
			}
			dials++

			var err error
			c, err = s.dialPeer(ctx, p)
			if err != nil {
				return nil, err
			}
		}
		s, err := c.NewStream(ctx)
		if err != nil {
			if c.conn.IsClosed() {
				continue
			}
			return nil, err
		}
		return s, nil
	}
}

// ConnsToPeer returns all the live connections to peer.
func (s *Swarm) ConnsToPeer(p peer.ID) []network.Conn {
	// TODO: Consider sorting the connection list best to worst. Currently,
	// it's sorted oldest to newest.
	s.conns.RLock()
	defer s.conns.RUnlock()
	conns := s.conns.m[p]
	output := make([]network.Conn, len(conns))
	for i, c := range conns {
		output[i] = c
	}
	return output
}

// bestConnToPeer returns the best connection to peer.
func (s *Swarm) bestConnToPeer(p peer.ID) *Conn {
	// Selects the best connection we have to the peer.
	// TODO: Prefer some transports over others. Currently, we just select
	// the newest non-closed connection with the most streams.
	s.conns.RLock()
	defer s.conns.RUnlock()

	var best *Conn
	bestLen := 0
	for _, c := range s.conns.m[p] {
		if c.conn.IsClosed() {
			// We *will* garbage collect this soon anyways.
			continue
		}
		c.streams.Lock()
		cLen := len(c.streams.m)
		c.streams.Unlock()

		if cLen >= bestLen {
			best = c
			bestLen = cLen
		}

	}
	return best
}

// Connectedness returns our "connectedness" state with the given peer.
//
// To check if we have an open connection, use `s.Connectedness(p) ==
// network.Connected`.
func (s *Swarm) Connectedness(p peer.ID) network.Connectedness {
	if s.bestConnToPeer(p) != nil {
		return network.Connected
	}
	return network.NotConnected
}

// Conns returns a slice of all connections.
func (s *Swarm) Conns() []network.Conn {
	s.conns.RLock()
	defer s.conns.RUnlock()

	conns := make([]network.Conn, 0, len(s.conns.m))
	for _, cs := range s.conns.m {
		for _, c := range cs {
			conns = append(conns, c)
		}
	}
	return conns
}

// ClosePeer closes all connections to the given peer.
func (s *Swarm) ClosePeer(p peer.ID) error {
	conns := s.ConnsToPeer(p)
	switch len(conns) {
	case 0:
		return nil
	case 1:
		return conns[0].Close()
	default:
		errCh := make(chan error)
		for _, c := range conns {
			go func(c network.Conn) {
				errCh <- c.Close()
			}(c)
		}

		var errs []string
		for _ = range conns {
			err := <-errCh
			if err != nil {
				errs = append(errs, err.Error())
			}
		}
		if len(errs) > 0 {
			return fmt.Errorf("when disconnecting from peer %s: %s", p, strings.Join(errs, ", "))
		}
		return nil
	}
}

// Peers returns a copy of the set of peers swarm is connected to.
func (s *Swarm) Peers() []peer.ID {
	s.conns.RLock()
	defer s.conns.RUnlock()
	peers := make([]peer.ID, 0, len(s.conns.m))
	for p := range s.conns.m {
		peers = append(peers, p)
	}

	return peers
}

// LocalPeer returns the local peer swarm is associated to.
func (s *Swarm) LocalPeer() peer.ID {
	return s.local
}

// Backoff returns the DialBackoff object for this swarm.
func (s *Swarm) Backoff() *DialBackoff {
	return &s.backf
}

// notifyAll sends a signal to all Notifiees
func (s *Swarm) notifyAll(notify func(network.Notifiee)) {
	var wg sync.WaitGroup

	s.notifs.RLock()
	wg.Add(len(s.notifs.m))
	for f := range s.notifs.m {
		go func(f network.Notifiee) {
			defer wg.Done()
			notify(f)
		}(f)
	}

	wg.Wait()
	s.notifs.RUnlock()
}

// Notify signs up Notifiee to receive signals when events happen
func (s *Swarm) Notify(f network.Notifiee) {
	s.notifs.Lock()
	s.notifs.m[f] = struct{}{}
	s.notifs.Unlock()
}

// StopNotify unregisters Notifiee fromr receiving signals
func (s *Swarm) StopNotify(f network.Notifiee) {
	s.notifs.Lock()
	delete(s.notifs.m, f)
	s.notifs.Unlock()
}

func (s *Swarm) removeConn(c *Conn) {
	p := c.RemotePeer()

	s.conns.Lock()
	defer s.conns.Unlock()
	cs := s.conns.m[p]
	for i, ci := range cs {
		if ci == c {
			if len(cs) == 1 {
				delete(s.conns.m, p)
			} else {
				// NOTE: We're intentionally preserving order.
				// This way, connections to a peer are always
				// sorted oldest to newest.
				copy(cs[i:], cs[i+1:])
				cs[len(cs)-1] = nil
				s.conns.m[p] = cs[:len(cs)-1]
			}
			return
		}
	}
}

// String returns a string representation of Network.
func (s *Swarm) String() string {
	return fmt.Sprintf("<Swarm %s>", s.LocalPeer())
}

// Swarm is a Network.
var _ network.Network = (*Swarm)(nil)
var _ transport.TransportNetwork = (*Swarm)(nil)
