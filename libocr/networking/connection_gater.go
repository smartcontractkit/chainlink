package networking

import (
	"sync"

	p2pcontrol "github.com/libp2p/go-libp2p-core/control"
	p2pnetwork "github.com/libp2p/go-libp2p-core/network"
	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/loghelper"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

type allower interface {
	isAllowed(p2ppeer.ID) bool
	allowlist() []p2ppeer.ID
}

type connectionGater struct {
	allowers    map[allower]struct{}
	allowersMtx sync.RWMutex
	logger      types.Logger
}

func newConnectionGater(logger types.Logger) (*connectionGater, error) {
	allowers := make(map[allower]struct{})

	logger = loghelper.MakeLoggerWithContext(logger, types.LogFields{
		"id": "ConnectionGater",
	})

	return &connectionGater{
		allowers:    allowers,
		allowersMtx: sync.RWMutex{},
		logger:      logger,
	}, nil
}

func (c *connectionGater) add(g allower) {
	c.allowersMtx.Lock()
	defer c.allowersMtx.Unlock()
	if _, exists := c.allowers[g]; exists {
		panic("allower has already been added")
	}
	c.allowers[g] = struct{}{}
}

func (c *connectionGater) remove(g allower) {
	c.allowersMtx.Lock()
	defer c.allowersMtx.Unlock()
	if _, exists := c.allowers[g]; !exists {
		panic("allower is not in list")
	}
	delete(c.allowers, g)
}

func (c *connectionGater) isAllowed(id p2ppeer.ID) bool {
	c.allowersMtx.RLock()
	defer c.allowersMtx.RUnlock()
	for g := range c.allowers {
		if g.isAllowed(id) {
			return true
		}
	}
	c.logger.Warn("ConnectionGater: denied access", types.LogFields{
		"remotePeerID": id,
	})
	return false
}

func (c *connectionGater) allowlist() (allowlist []p2ppeer.ID) {
	c.allowersMtx.RLock()
	defer c.allowersMtx.RUnlock()
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

func (c *connectionGater) InterceptAccept(p2pnetwork.ConnMultiaddrs) (allow bool) {
	return true
}

func (c *connectionGater) InterceptAddrDial(id p2ppeer.ID, _ ma.Multiaddr) (allow bool) {
	return c.isAllowed(id)
}

func (c *connectionGater) InterceptPeerDial(id p2ppeer.ID) (allow bool) {
	return c.isAllowed(id)
}

func (c *connectionGater) InterceptSecured(_ p2pnetwork.Direction, id p2ppeer.ID, _ p2pnetwork.ConnMultiaddrs) (allow bool) {
	return c.isAllowed(id)
}

func (c *connectionGater) InterceptUpgraded(p2pnetwork.Conn) (allow bool, reason p2pcontrol.DisconnectReason) {
	return true, p2pcontrol.DisconnectReason(0)
}
