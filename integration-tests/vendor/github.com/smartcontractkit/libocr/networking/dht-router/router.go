package dhtrouter

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p-core/peerstore"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/subprocesses"

	"github.com/libp2p/go-libp2p-core/peer"
	p2pprotocol "github.com/libp2p/go-libp2p-core/protocol"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	kbucket "github.com/libp2p/go-libp2p-kbucket"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
)

var InvalidDhtKey = errors.New("invalid dht key")

// PeerDiscoveryRouter is a router (rhost.Routing) with resource management capabilities via Start() and Close()
type PeerDiscoveryRouter interface {
	rhost.Routing
	Start()
	Close() error
	ProtocolID() p2pprotocol.ID
}

type DHTRouter struct {
	aclHost ACLHost
	dht     *dht.IpfsDHT
	config  DHTNodeConfig

	processes *subprocesses.Subprocesses
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewDHTRouter(ctx context.Context, config DHTNodeConfig, aclHost ACLHost) (PeerDiscoveryRouter, error) {
	kadDht, err := newDHT(ctx, config, aclHost)
	if err != nil {
		return nil, err
	}

	newCtx, cancel := context.WithCancel(ctx)

	router := DHTRouter{
		aclHost,
		kadDht,
		config,
		&subprocesses.Subprocesses{},
		newCtx,
		cancel,
	}

	return router, nil
}

// Start ought to be called before the router is used
//
// Start optimistically tries to connect to bootstrappers and populate the
// initial table, but if that fails it puts the connect into
// a background thread and exits
//
// Start runs asynchronously in the background and can be cancelled early at
// any time by calling Close() (which does run synchronously and waits for
// Start to exit)
//

func (router DHTRouter) Start() {
	// Connect to all the bootstrap nodes
	router.processes.Go(func() {
		router.logger().Debug("DHT initial bootstrap starting", nil)

		err := tryConnectToBootstrappers(router.ctx, router.aclHost, router.config.bootstrapNodes)
		if err != nil {
			router.logger().Error("DHT initial bootstrap connect failed", commontypes.LogFields{
				"err": err.Error(),
			})
		} else {
			// Make a best effort to populate initial routing table
			err = <-router.dht.ForceRefresh()
			if err != nil {
				router.logger().Warn("Initial DHT table refresh failed", commontypes.LogFields{
					"err": err.Error(),
				})
			}
		}

		router.logger().Info("DHT initial bootstrap complete", commontypes.LogFields{
			"bnodes": router.config.bootstrapNodes,
		})

		// start a thread to re-connect to all bootstrap nodes every router.config.bootstrapCheckInterval
		router.processes.RepeatWithCancel("bootstrap", router.config.bootstrapCheckInterval, router.ctx, func() {
			toConnect := false
			for _, p := range router.config.bootstrapNodes {
				// reconnect if any connection is lost.
				if router.aclHost.Network().Connectedness(p.ID) != network.Connected {
					toConnect = true
					break
				}
			}

			if toConnect {
				router.logger().Debug("connect to bootstrap nodes", commontypes.LogFields{
					"bnodes": router.config.bootstrapNodes,
				})
				if err := tryConnectToBootstrappers(router.ctx, router.aclHost, router.config.bootstrapNodes); err != nil {
					router.logger().Warn("DHT has no connection to any bootstrappers", commontypes.LogFields{
						"err": err.Error(),
					})
				}
			}
		})
		router.startAnnounceInBackground()

		if router.config.extendedDHTLogging {
			// default RT refresh time in libp2p is 10 minutes
			router.printPeriodicReport(10 * time.Minute)
		}
	})
}

func (router DHTRouter) ProtocolID() p2pprotocol.ID {
	return router.config.ProtocolID()
}

func (router DHTRouter) logger() commontypes.Logger {
	return router.config.logger
}

func peerIdToDhtKey(id peer.ID) string {
	return fmt.Sprintf("/%s/%s", ValidatorNamespace, id.String())
}

func dhtKeyToPeerId(s string) (id peer.ID, err error) {
	// of format /ns/key
	ss := strings.Split(s, "/")

	if len(ss) != 3 {
		return id, InvalidDhtKey
	}

	ns := ss[1]
	key := ss[2]

	if ns != ValidatorNamespace {
		return id, InvalidDhtKey
	}

	id, err = peer.Decode(key)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (router DHTRouter) FindPeer(ctx context.Context, peerId peer.ID) (addr peer.AddrInfo, error error) {
	key := peerIdToDhtKey(peerId)

	router.logger().Debug("DHT: Lookup peer", commontypes.LogFields{
		"peerID": peerId,
		"key":    key,
	})

	marshaled, err := router.dht.GetValue(ctx, key)
	if err != nil {
		return addr, err
	}

	// unmarshal
	ann, err := deserializeSignedAnnouncement(marshaled)
	if err != nil {
		return addr, err
	}

	// at this point ann has been verified by AnnouncementValidator (called in GetValue)
	addr.ID = peerId
	addr.Addrs = ann.Addrs

	router.logger().Debug("DHT: Found peer", commontypes.LogFields{
		"peerID":              peerId,
		"key":                 key,
		"foundPeerAddr":       addr,
		"announcementCounter": ann.Counter,
	})
	return addr, nil
}

const counterKeyName = "counter"

func (router DHTRouter) publishHostAddr(ctx context.Context) error {
	var addrs = router.aclHost.Addrs()

	// cap at maxAddrInAnnouncements addresses
	if len(router.aclHost.Addrs()) > maxAddrInAnnouncements {
		router.logger().Warn("trying to publish many addresses. capped.", commontypes.LogFields{
			"nAddrs":            len(router.aclHost.Addrs()),
			"addrsNotAnnounced": addrs[:maxAddrInAnnouncements],
		})

		addrs = addrs[:maxAddrInAnnouncements]
	}

	// try to retrieve last counter
	counter, err := router.aclHost.Peerstore().Get(router.aclHost.ID(), counterKeyName)
	// IMPORTANT: 🚨🚨🚨 make sure our own peer store returns peerstore.ErrNotFound too! 🚨🚨🚨
	if errors.Is(err, peerstore.ErrNotFound) {
		// if the db is empty, counter starts with zero
		counter = uint64(0)
	} else if err != nil {
		return err
	}

	c, ok := counter.(uint64)
	if !ok {
		return errors.New("cannot convert counter to uint64")
	}

	// advance the counter but dont persist it yet
	newCounter := c + 1
	if newCounter < c {
		return errors.New("DHT Announcement counter overflowed")
	}

	// persist the new counter before making an announcement
	err = router.aclHost.Peerstore().Put(router.aclHost.ID(), counterKeyName, newCounter)
	if err != nil {
		return err
	}

	// retrieve the private keys
	sk := router.aclHost.Peerstore().PrivKey(router.aclHost.ID())

	ann := announcement{
		addrs,
		announcementCounter{
			router.config.announcementUserPrefix,
			newCounter,
		},
	}

	signedAnn, err := ann.sign(sk)
	if err != nil {
		return err
	}

	marshaled, err := signedAnn.serialize()
	if err != nil {
		return err
	}

	key := peerIdToDhtKey(router.aclHost.ID())
	router.logger().Debug("DHT: Put value", commontypes.LogFields{
		"key":         key,
		"addrs":       ann.Addrs,
		"counter":     ann.Counter,
		"value (hex)": hex.EncodeToString(marshaled),
	})

	err = router.dht.PutValue(ctx, key, marshaled)
	if err != nil {
		return errors.Wrap(err, "could not publish address")
	} else {
		return nil
	}
}

func (router DHTRouter) startAnnounceInBackground() {
	router.processes.Go(func() {
		const retryInterval = 10 * time.Second

		ticker := time.NewTicker(retryInterval)
		defer ticker.Stop()

		for {
			select {
			case <-router.ctx.Done():
				return
			case <-ticker.C:
				err := router.publishHostAddr(router.ctx)
				if err != nil {
					router.logger().Warn("DHT: Error publishing address", commontypes.LogFields{
						"err":     err.Error(),
						"retryIn": retryInterval,
					})
				} else {
					router.logger().Info("DHT: Published address", nil)
					return
				}
			}
		}
	})
}

// This is used by bootstrap nodes to periodically report their states. (Otherwise their logging is pretty sparse.)
func (router DHTRouter) printPeriodicReport(interval time.Duration) {
	if !router.config.extendedDHTLogging {
		return
	}

	router.processes.RepeatWithCancel("periodic report", interval, router.ctx, func() {
		rt := router.dht.RoutingTable()

		rtString := []string{fmt.Sprintf("RT has %d entries,", len(rt.ListPeers()))}
		for _, p := range rt.ListPeers() {
			rtString = append(rtString, fmt.Sprintf("peerId=%s, kadId=%s, allowed=%t",
				p.Pretty(),
				hex.EncodeToString(kbucket.ConvertPeerID(p)),
				router.aclHost.GetACL().IsAllowed(p, router.config.ProtocolID())))
		}

		pstore := router.aclHost.Peerstore()
		peerAddrs := make(map[string][]string)
		for _, myPeer := range pstore.Peers() {
			var addrs []string
			for _, addr := range pstore.Addrs(myPeer) {
				addrs = append(addrs, addr.String())
			}
			peerAddrs[myPeer.Pretty()] = addrs
		}

		router.logger().Debug("DHT periodical report", commontypes.LogFields{
			"protocolID": router.config.ProtocolID(),
			"acl":        router.aclHost.GetACL(),
			"rt":         rtString,
			"peerstore":  peerAddrs,
		})
	})
}

func (router DHTRouter) Close() error {
	router.ctxCancel()
	router.processes.Wait()

	return router.dht.Close()
}
