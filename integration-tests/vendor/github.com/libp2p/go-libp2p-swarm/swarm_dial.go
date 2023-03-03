package swarm

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"

	addrutil "github.com/libp2p/go-addr-util"
	lgbl "github.com/libp2p/go-libp2p-loggables"

	logging "github.com/ipfs/go-log"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

// Diagram of dial sync:
//
//   many callers of Dial()   synched w.  dials many addrs       results to callers
//  ----------------------\    dialsync    use earliest            /--------------
//  -----------------------\              |----------\           /----------------
//  ------------------------>------------<-------     >---------<-----------------
//  -----------------------|              \----x                 \----------------
//  ----------------------|                \-----x                \---------------
//                                         any may fail          if no addr at end
//                                                             retry dialAttempt x

var (
	// ErrDialBackoff is returned by the backoff code when a given peer has
	// been dialed too frequently
	ErrDialBackoff = errors.New("dial backoff")

	// ErrDialToSelf is returned if we attempt to dial our own peer
	ErrDialToSelf = errors.New("dial to self attempted")

	// ErrNoTransport is returned when we don't know a transport for the
	// given multiaddr.
	ErrNoTransport = errors.New("no transport for protocol")

	// ErrAllDialsFailed is returned when connecting to a peer has ultimately failed
	ErrAllDialsFailed = errors.New("all dials failed")

	// ErrNoAddresses is returned when we fail to find any addresses for a
	// peer we're trying to dial.
	ErrNoAddresses = errors.New("no addresses")

	// ErrNoGoodAddresses is returned when we find addresses for a peer but
	// can't use any of them.
	ErrNoGoodAddresses = errors.New("no good addresses")

	// ErrGaterDisallowedConnection is returned when the gater prevents us from
	// forming a connection with a peer.
	ErrGaterDisallowedConnection = errors.New("gater disallows connection to peer")
)

// DialAttempts governs how many times a goroutine will try to dial a given peer.
// Note: this is down to one, as we have _too many dials_ atm. To add back in,
// add loop back in Dial(.)
const DialAttempts = 1

// ConcurrentFdDials is the number of concurrent outbound dials over transports
// that consume file descriptors
const ConcurrentFdDials = 160

// DefaultPerPeerRateLimit is the number of concurrent outbound dials to make
// per peer
const DefaultPerPeerRateLimit = 8

// dialbackoff is a struct used to avoid over-dialing the same, dead peers.
// Whenever we totally time out on a peer (all three attempts), we add them
// to dialbackoff. Then, whenevers goroutines would _wait_ (dialsync), they
// check dialbackoff. If it's there, they don't wait and exit promptly with
// an error. (the single goroutine that is actually dialing continues to
// dial). If a dial is successful, the peer is removed from backoff.
// Example:
//
//  for {
//  	if ok, wait := dialsync.Lock(p); !ok {
//  		if backoff.Backoff(p) {
//  			return errDialFailed
//  		}
//  		<-wait
//  		continue
//  	}
//  	defer dialsync.Unlock(p)
//  	c, err := actuallyDial(p)
//  	if err != nil {
//  		dialbackoff.AddBackoff(p)
//  		continue
//  	}
//  	dialbackoff.Clear(p)
//  }
//

// DialBackoff is a type for tracking peer dial backoffs.
//
// * It's safe to use its zero value.
// * It's thread-safe.
// * It's *not* safe to move this type after using.
type DialBackoff struct {
	entries map[peer.ID]map[string]*backoffAddr
	lock    sync.RWMutex
}

type backoffAddr struct {
	tries int
	until time.Time
}

func (db *DialBackoff) init(ctx context.Context) {
	if db.entries == nil {
		db.entries = make(map[peer.ID]map[string]*backoffAddr)
	}
	go db.background(ctx)
}

func (db *DialBackoff) background(ctx context.Context) {
	ticker := time.NewTicker(BackoffMax)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			db.cleanup()
		}
	}
}

// Backoff returns whether the client should backoff from dialing
// peer p at address addr
func (db *DialBackoff) Backoff(p peer.ID, addr ma.Multiaddr) (backoff bool) {
	db.lock.Lock()
	defer db.lock.Unlock()

	ap, found := db.entries[p][string(addr.Bytes())]
	return found && time.Now().Before(ap.until)
}

// BackoffBase is the base amount of time to backoff (default: 5s).
var BackoffBase = time.Second * 5

// BackoffCoef is the backoff coefficient (default: 1s).
var BackoffCoef = time.Second

// BackoffMax is the maximum backoff time (default: 5m).
var BackoffMax = time.Minute * 5

// AddBackoff lets other nodes know that we've entered backoff with
// peer p, so dialers should not wait unnecessarily. We still will
// attempt to dial with one goroutine, in case we get through.
//
// Backoff is not exponential, it's quadratic and computed according to the
// following formula:
//
//     BackoffBase + BakoffCoef * PriorBackoffs^2
//
// Where PriorBackoffs is the number of previous backoffs.
func (db *DialBackoff) AddBackoff(p peer.ID, addr ma.Multiaddr) {
	saddr := string(addr.Bytes())
	db.lock.Lock()
	defer db.lock.Unlock()
	bp, ok := db.entries[p]
	if !ok {
		bp = make(map[string]*backoffAddr, 1)
		db.entries[p] = bp
	}
	ba, ok := bp[saddr]
	if !ok {
		bp[saddr] = &backoffAddr{
			tries: 1,
			until: time.Now().Add(BackoffBase),
		}
		return
	}

	backoffTime := BackoffBase + BackoffCoef*time.Duration(ba.tries*ba.tries)
	if backoffTime > BackoffMax {
		backoffTime = BackoffMax
	}
	ba.until = time.Now().Add(backoffTime)
	ba.tries++
}

// Clear removes a backoff record. Clients should call this after a
// successful Dial.
func (db *DialBackoff) Clear(p peer.ID) {
	db.lock.Lock()
	defer db.lock.Unlock()
	delete(db.entries, p)
}

func (db *DialBackoff) cleanup() {
	db.lock.Lock()
	defer db.lock.Unlock()
	now := time.Now()
	for p, e := range db.entries {
		good := false
		for _, backoff := range e {
			backoffTime := BackoffBase + BackoffCoef*time.Duration(backoff.tries*backoff.tries)
			if backoffTime > BackoffMax {
				backoffTime = BackoffMax
			}
			if now.Before(backoff.until.Add(backoffTime)) {
				good = true
				break
			}
		}
		if !good {
			delete(db.entries, p)
		}
	}
}

// DialPeer connects to a peer.
//
// The idea is that the client of Swarm does not need to know what network
// the connection will happen over. Swarm can use whichever it choses.
// This allows us to use various transport protocols, do NAT traversal/relay,
// etc. to achieve connection.
func (s *Swarm) DialPeer(ctx context.Context, p peer.ID) (network.Conn, error) {
	if s.gater != nil && !s.gater.InterceptPeerDial(p) {
		log.Debugf("gater disallowed outbound connection to peer %s", p.Pretty())
		return nil, &DialError{Peer: p, Cause: ErrGaterDisallowedConnection}
	}

	return s.dialPeer(ctx, p)
}

// internal dial method that returns an unwrapped conn
//
// It is gated by the swarm's dial synchronization systems: dialsync and
// dialbackoff.
func (s *Swarm) dialPeer(ctx context.Context, p peer.ID) (*Conn, error) {
	log.Debugf("[%s] swarm dialing peer [%s]", s.local, p)
	var logdial = lgbl.Dial("swarm", s.LocalPeer(), p, nil, nil)
	err := p.Validate()
	if err != nil {
		return nil, err
	}

	if p == s.local {
		log.Event(ctx, "swarmDialSelf", logdial)
		return nil, ErrDialToSelf
	}

	defer log.EventBegin(ctx, "swarmDialAttemptSync", p).Done()

	// check if we already have an open connection first
	conn := s.bestConnToPeer(p)
	if conn != nil {
		return conn, nil
	}

	// apply the DialPeer timeout
	ctx, cancel := context.WithTimeout(ctx, network.GetDialPeerTimeout(ctx))
	defer cancel()

	conn, err = s.dsync.DialLock(ctx, p)
	if err == nil {
		return conn, nil
	}

	log.Debugf("network for %s finished dialing %s", s.local, p)

	if ctx.Err() != nil {
		// Context error trumps any dial errors as it was likely the ultimate cause.
		return nil, ctx.Err()
	}

	if s.ctx.Err() != nil {
		// Ok, so the swarm is shutting down.
		return nil, ErrSwarmClosed
	}

	return nil, err
}

// doDial is an ugly shim method to retain all the logging and backoff logic
// of the old dialsync code
func (s *Swarm) doDial(ctx context.Context, p peer.ID) (*Conn, error) {
	// Short circuit.
	// By the time we take the dial lock, we may already *have* a connection
	// to the peer.
	c := s.bestConnToPeer(p)
	if c != nil {
		return c, nil
	}

	logdial := lgbl.Dial("swarm", s.LocalPeer(), p, nil, nil)

	// ok, we have been charged to dial! let's do it.
	// if it succeeds, dial will add the conn to the swarm itself.
	defer log.EventBegin(ctx, "swarmDialAttemptStart", logdial).Done()

	conn, err := s.dial(ctx, p)
	if err != nil {
		conn = s.bestConnToPeer(p)
		if conn != nil {
			// Hm? What error?
			// Could have canceled the dial because we received a
			// connection or some other random reason.
			// Just ignore the error and return the connection.
			log.Debugf("ignoring dial error because we have a connection: %s", err)
			return conn, nil
		}

		// ok, we failed.
		return nil, err
	}
	return conn, nil
}

func (s *Swarm) canDial(addr ma.Multiaddr) bool {
	t := s.TransportForDialing(addr)
	return t != nil && t.CanDial(addr)
}

// dial is the actual swarm's dial logic, gated by Dial.
func (s *Swarm) dial(ctx context.Context, p peer.ID) (*Conn, error) {
	var logdial = lgbl.Dial("swarm", s.LocalPeer(), p, nil, nil)
	if p == s.local {
		log.Event(ctx, "swarmDialDoDialSelf", logdial)
		return nil, ErrDialToSelf
	}
	defer log.EventBegin(ctx, "swarmDialDo", logdial).Done()
	logdial["dial"] = "failure" // start off with failure. set to "success" at the end.

	sk := s.peers.PrivKey(s.local)
	logdial["encrypted"] = sk != nil // log whether this will be an encrypted dial or not.
	if sk == nil {
		// fine for sk to be nil, just log.
		log.Debug("Dial not given PrivateKey, so WILL NOT SECURE conn.")
	}

	//////
	peerAddrs := s.peers.Addrs(p)
	if len(peerAddrs) == 0 {
		return nil, &DialError{Peer: p, Cause: ErrNoAddresses}
	}
	goodAddrs := s.filterKnownUndialables(p, peerAddrs)
	if len(goodAddrs) == 0 {
		return nil, &DialError{Peer: p, Cause: ErrNoGoodAddresses}
	}

	/////// Check backoff andnRank addresses
	var nonBackoff bool
	for _, a := range goodAddrs {
		// skip addresses in back-off
		if !s.backf.Backoff(p, a) {
			nonBackoff = true
		}
	}
	if !nonBackoff {
		return nil, ErrDialBackoff
	}

	// ranks addresses in descending order of preference for dialing
	// Private UDP > Public UDP > Private TCP > Public TCP > UDP Relay server > TCP Relay server
	rankAddrsFnc := func(addrs []ma.Multiaddr) []ma.Multiaddr {
		var localUdpAddrs []ma.Multiaddr // private udp
		var relayUdpAddrs []ma.Multiaddr // relay udp
		var othersUdp []ma.Multiaddr     // public udp

		var localFdAddrs []ma.Multiaddr // private fd consuming
		var relayFdAddrs []ma.Multiaddr //  relay fd consuming
		var othersFd []ma.Multiaddr     // public fd consuming

		for _, a := range addrs {
			if _, err := a.ValueForProtocol(ma.P_CIRCUIT); err == nil {
				if s.IsFdConsumingAddr(a) {
					relayFdAddrs = append(relayFdAddrs, a)
					continue
				}
				relayUdpAddrs = append(relayUdpAddrs, a)
			} else if manet.IsPrivateAddr(a) {
				if s.IsFdConsumingAddr(a) {
					localFdAddrs = append(localFdAddrs, a)
					continue
				}
				localUdpAddrs = append(localUdpAddrs, a)
			} else {
				if s.IsFdConsumingAddr(a) {
					othersFd = append(othersFd, a)
					continue
				}
				othersUdp = append(othersUdp, a)
			}
		}

		relays := append(relayUdpAddrs, relayFdAddrs...)
		fds := append(localFdAddrs, othersFd...)

		return append(append(append(localUdpAddrs, othersUdp...), fds...), relays...)
	}

	connC, dialErr := s.dialAddrs(ctx, p, rankAddrsFnc(goodAddrs))

	if dialErr != nil {
		logdial["error"] = dialErr.Cause.Error()
		switch dialErr.Cause {
		case context.Canceled, context.DeadlineExceeded:
			// Always prefer the context errors as we rely on being
			// able to check them.
			//
			// Removing this will BREAK backoff (causing us to
			// backoff when canceling dials).
			return nil, dialErr.Cause
		}
		return nil, dialErr
	}
	logdial["conn"] = logging.Metadata{
		"localAddr":  connC.LocalMultiaddr(),
		"remoteAddr": connC.RemoteMultiaddr(),
	}
	swarmC, err := s.addConn(connC, network.DirOutbound)
	if err != nil {
		logdial["error"] = err.Error()
		connC.Close() // close the connection. didn't work out :(
		return nil, &DialError{Peer: p, Cause: err}
	}

	logdial["dial"] = "success"
	return swarmC, nil
}

// filterKnownUndialables takes a list of multiaddrs, and removes those
// that we definitely don't want to dial: addresses configured to be blocked,
// IPv6 link-local addresses, addresses without a dial-capable transport,
// and addresses that we know to be our own.
// This is an optimization to avoid wasting time on dials that we know are going to fail.
func (s *Swarm) filterKnownUndialables(p peer.ID, addrs []ma.Multiaddr) []ma.Multiaddr {
	lisAddrs, _ := s.InterfaceListenAddresses()
	var ourAddrs []ma.Multiaddr
	for _, addr := range lisAddrs {
		protos := addr.Protocols()
		// we're only sure about filtering out /ip4 and /ip6 addresses, so far
		if len(protos) == 2 && (protos[0].Code == ma.P_IP4 || protos[0].Code == ma.P_IP6) {
			ourAddrs = append(ourAddrs, addr)
		}
	}

	return addrutil.FilterAddrs(addrs,
		addrutil.SubtractFilter(ourAddrs...),
		s.canDial,
		// TODO: Consider allowing link-local addresses
		addrutil.AddrOverNonLocalIP,
		func(addr ma.Multiaddr) bool {
			return s.gater == nil || s.gater.InterceptAddrDial(p, addr)
		},
	)
}

func (s *Swarm) dialAddrs(ctx context.Context, p peer.ID, remoteAddrs []ma.Multiaddr) (transport.CapableConn, *DialError) {
	/*
		This slice-to-chan code is temporary, the peerstore can currently provide
		a channel as an interface for receiving addresses, but more thought
		needs to be put into the execution. For now, this allows us to use
		the improved rate limiter, while maintaining the outward behaviour
		that we previously had (halting a dial when we run out of addrs)
	*/
	var remoteAddrChan chan ma.Multiaddr
	if len(remoteAddrs) > 0 {
		remoteAddrChan = make(chan ma.Multiaddr, len(remoteAddrs))
		for i := range remoteAddrs {
			remoteAddrChan <- remoteAddrs[i]
		}
		close(remoteAddrChan)
	}

	log.Debugf("%s swarm dialing %s", s.local, p)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // cancel work when we exit func

	// use a single response type instead of errs and conns, reduces complexity *a ton*
	respch := make(chan dialResult)
	err := &DialError{Peer: p}

	defer s.limiter.clearAllPeerDials(p)

	var active int
dialLoop:
	for remoteAddrChan != nil || active > 0 {
		// Check for context cancellations and/or responses first.
		select {
		case <-ctx.Done():
			break dialLoop
		case resp := <-respch:
			active--
			if resp.Err != nil {
				// Errors are normal, lots of dials will fail
				if resp.Err != context.Canceled {
					s.backf.AddBackoff(p, resp.Addr)
				}

				log.Infof("got error on dial: %s", resp.Err)
				err.recordErr(resp.Addr, resp.Err)
			} else if resp.Conn != nil {
				return resp.Conn, nil
			}

			// We got a result, try again from the top.
			continue
		default:
		}

		// Now, attempt to dial.
		select {
		case addr, ok := <-remoteAddrChan:
			if !ok {
				remoteAddrChan = nil
				continue
			}

			s.limitedDial(ctx, p, addr, respch)
			active++
		case <-ctx.Done():
			break dialLoop
		case resp := <-respch:
			active--
			if resp.Err != nil {
				// Errors are normal, lots of dials will fail
				if resp.Err != context.Canceled {
					s.backf.AddBackoff(p, resp.Addr)
				}

				log.Infof("got error on dial: %s", resp.Err)
				err.recordErr(resp.Addr, resp.Err)
			} else if resp.Conn != nil {
				return resp.Conn, nil
			}
		}
	}

	if ctxErr := ctx.Err(); ctxErr != nil {
		err.Cause = ctxErr
	} else if len(err.DialErrors) == 0 {
		err.Cause = network.ErrNoRemoteAddrs
	} else {
		err.Cause = ErrAllDialsFailed
	}
	return nil, err
}

// limitedDial will start a dial to the given peer when
// it is able, respecting the various different types of rate
// limiting that occur without using extra goroutines per addr
func (s *Swarm) limitedDial(ctx context.Context, p peer.ID, a ma.Multiaddr, resp chan dialResult) {
	s.limiter.AddDialJob(&dialJob{
		addr: a,
		peer: p,
		resp: resp,
		ctx:  ctx,
	})
}

func (s *Swarm) dialAddr(ctx context.Context, p peer.ID, addr ma.Multiaddr) (transport.CapableConn, error) {
	// Just to double check. Costs nothing.
	if s.local == p {
		return nil, ErrDialToSelf
	}
	log.Debugf("%s swarm dialing %s %s", s.local, p, addr)

	tpt := s.TransportForDialing(addr)
	if tpt == nil {
		return nil, ErrNoTransport
	}

	connC, err := tpt.Dial(ctx, addr, p)
	if err != nil {
		return nil, err
	}

	// Trust the transport? Yeah... right.
	if connC.RemotePeer() != p {
		connC.Close()
		err = fmt.Errorf("BUG in transport %T: tried to dial %s, dialed %s", p, connC.RemotePeer(), tpt)
		log.Error(err)
		return nil, err
	}

	// success! we got one!
	return connC, nil
}

// TODO We should have a `IsFdConsuming() bool` method on the `Transport` interface in go-libp2p-core/transport.
// This function checks if any of the transport protocols in the address requires a file descriptor.
// For now:
// A Non-circuit address which has the TCP/UNIX protocol is deemed FD consuming.
// For a circuit-relay address, we look at the address of the relay server/proxy
// and use the same logic as above to decide.
func (s *Swarm) IsFdConsumingAddr(addr ma.Multiaddr) bool {
	first, _ := ma.SplitFunc(addr, func(c ma.Component) bool {
		return c.Protocol().Code == ma.P_CIRCUIT
	})

	// for safety
	if first == nil {
		return true
	}

	_, err1 := first.ValueForProtocol(ma.P_TCP)
	_, err2 := first.ValueForProtocol(ma.P_UNIX)
	return err1 == nil || err2 == nil
}
