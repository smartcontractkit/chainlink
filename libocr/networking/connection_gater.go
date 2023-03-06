package networking

import (
	"sync"

	p2pcontrol "github.com/libp2p/go-libp2p-core/control"
	p2pnetwork "github.com/libp2p/go-libp2p-core/network"
	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"golang.org/x/time/rate"
)

// Rate limits for inbound connection establishment. The limits are per-peer.
// Note that we only rateLimit inbound connections *after* authentication
// has taken place. This is to prevent an attacker from spoofing their identity
// and exhausting an honest node's rate limit.
const (
	maxConnectionsPerSecond = .5
	maxConnectionsBurst     = 2
)

// allower controls which peers are allowed
type allower interface {
	isAllowed(p2ppeer.ID) bool
	allowlist() []p2ppeer.ID
}

// Allowers are OR'd together. As long as the remote peer is allowed by one of
// the allowers, the connection will be allowed.
type connectionGater struct {
	connLimiters map[p2ppeer.ID]*rate.Limiter
	allowers     map[allower]struct{}
	mutex        sync.RWMutex
	logger       loghelper.LoggerWithContext
}

func newConnectionGater(logger loghelper.LoggerWithContext) (*connectionGater, error) {
	allowers := make(map[allower]struct{})

	logger = logger.MakeChild(commontypes.LogFields{
		"id": "ConnectionGater",
	})

	return &connectionGater{
		connLimiters: map[p2ppeer.ID]*rate.Limiter{},
		allowers:     allowers,
		mutex:        sync.RWMutex{},
		logger:       logger,
	}, nil
}

func (c *connectionGater) add(g allower) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, exists := c.allowers[g]; exists {
		panic("allower has already been added")
	}
	c.allowers[g] = struct{}{}
}

func (c *connectionGater) remove(g allower) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, exists := c.allowers[g]; !exists {
		panic("allower is not in list")
	}
	delete(c.allowers, g)
}

func (c *connectionGater) isAllowed(id p2ppeer.ID, checkRateLimit bool) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	oneAllowerPasses := false
	for g := range c.allowers {
		if g.isAllowed(id) {
			oneAllowerPasses = true
			break
		}
	}
	if !oneAllowerPasses {
		c.logger.Warn("ConnectionGater: denied access", commontypes.LogFields{
			"remotePeerID": id,
		})
		return false
	}
	if !checkRateLimit {
		return true
	}
	// instantiate new limiter if needed
	_, found := c.connLimiters[id]
	if !found {
		c.connLimiters[id] = rate.NewLimiter(maxConnectionsPerSecond, maxConnectionsBurst)
	}
	return c.connLimiters[id].Allow()
}

// Returns the full set of peers allowlisted by the current allowers
func (c *connectionGater) allowlist() (allowlist []p2ppeer.ID) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	set := make(map[p2ppeer.ID]struct{})

	for a := range c.allowers {
		for _, pid := range a.allowlist() {
			set[pid] = struct{}{}
		}
	}

	for pid := range set {
		allowlist = append(allowlist, pid)
	}
	return
}

// Interface methods for ConnectionGater
// See: https://github.com/libp2p/go-libp2p-core/blob/909c77480f732b9e9e0aa6857220950665b1b64b/connmgr/gater.go

// InterceptAccept tests whether an incipient inbound connection is allowed.
//
// This is called by the upgrader, or by the transport directly (e.g. QUIC,
// Bluetooth), straight after it has accepted a connection from its socket.
//
// Implementation accepts connections from all sources
func (c *connectionGater) InterceptAccept(p2pnetwork.ConnMultiaddrs) (allow bool) {
	return true
}

// InterceptAddrDial tests whether we're permitted to dial the specified
// multiaddr for the given peer.
//
// This is called by the network.Network implementation after it has
// resolved the peer's addrs, and prior to dialling each.
//
// Implementation restricts incoming peer IDs to those present in our oracle mappings.
// We don't apply bandwidth limiting at this point because we want the peer to authenticate first.
func (c *connectionGater) InterceptAddrDial(id p2ppeer.ID, _ ma.Multiaddr) (allow bool) {
	return c.isAllowed(id, false)
}

// InterceptPeerDial tests whether we're permitted to Dial the specified peer.
//
// This is called by the network.Network implementation when dialling a peer.
//
// Implementation prevents dialling to any peer not present in our oracle mappings.
// We don't apply bandwidth limiting at this point because we want the peer to authenticate first.
func (c *connectionGater) InterceptPeerDial(id p2ppeer.ID) (allow bool) {
	return c.isAllowed(id, false)
}

// InterceptSecured tests whether a given connection, now authenticated,
// is allowed.
//
// This is called by the upgrader, after it has performed the security
// handshake, and before it negotiates the muxer, or by the directly by the
// transport, at the exact same checkpoint.
//
// Implementation restricts incoming peer IDs to those present in our oracle mappings.
// This applies to both incoming and outgoing connections. Note that this is the only
// function that can prevent unknown incoming peers because any peer can lie about its
// peer ID until after a secure connection has been established.
//
// Bandwidth rate limiting kicks in here for inbound traffic coming from any other peer.
func (c *connectionGater) InterceptSecured(d p2pnetwork.Direction, id p2ppeer.ID, _ p2pnetwork.ConnMultiaddrs) (allow bool) {
	return c.isAllowed(id, d == p2pnetwork.DirInbound)
}

// InterceptUpgraded tests whether a fully capable connection is allowed.
//
// At this point, the connection a multiplexer has been selected.
// When rejecting a connection, the gater can return a DisconnectReason.
// Refer to the godoc on the ConnectionGater type for more information.
//
// NOTE: the go-libp2p implementation currently IGNORES the disconnect reason.
//
// Implementation allows everything.
func (c *connectionGater) InterceptUpgraded(p2pnetwork.Conn) (allow bool, reason p2pcontrol.DisconnectReason) {
	return true, p2pcontrol.DisconnectReason(0)
}
