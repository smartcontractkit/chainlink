package dhtrouter

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
	"github.com/smartcontractkit/chainlink/libocr/subprocesses"

	"github.com/libp2p/go-libp2p-core/peer"
	p2pprotocol "github.com/libp2p/go-libp2p-core/protocol"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	ma "github.com/multiformats/go-multiaddr"

	kbucket "github.com/libp2p/go-libp2p-kbucket"
)

var InvalidSignature = errors.New("invalid signature")
var CannotGetPrivateKey = errors.New("cannot get peer private key")
var InvalidDhtKey = errors.New("invalid dht key")
var InvalidDhtValue = errors.New("invalid dht msg")

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
		aclHost: aclHost,
		dht:     kadDht,
		config:  config,

		processes: &subprocesses.Subprocesses{},
		ctx:       newCtx,
		ctxCancel: cancel,
	}

	return router, nil
}

func (router DHTRouter) Start() {
	router.processes.Go(func() {
		router.logger().Debug("DHT initial bootstrap starting", nil)

		err := tryConnectToBootstrappers(router.ctx, router.aclHost, router.config.bootstrapNodes)
		if err != nil {
			router.logger().Error("DHT initial bootstrap connect failed", types.LogFields{
				"err": err.Error(),
			})
		} else {
			err = <-router.dht.ForceRefresh()
			if err != nil {
				router.logger().Warn("Initial DHT table refresh failed", types.LogFields{
					"err": err.Error(),
				})
			}
		}

		router.logger().Info("DHT initial bootstrap complete", nil)

		router.processes.RepeatWithCancel("bootstrap", 10*time.Second, router.ctx, func() {
			toConnect := false
			for _, peer := range router.config.bootstrapNodes {
				if router.aclHost.Network().Connectedness(peer.ID) != network.Connected {
					toConnect = true
					break
				}
			}

			if toConnect {
				if err := tryConnectToBootstrappers(router.ctx, router.aclHost, router.config.bootstrapNodes); err != nil {
					router.logger().Warn("DHT has no connection to any bootstrappers", types.LogFields{
						"err": err.Error(),
					})
				}
			}
		})
		router.startAnnounceInBackground()

		if router.config.extendedDHTLogging {
			router.printRoutingTableAndAclTable(10 * time.Minute)
		}
	})
}

func (router DHTRouter) ProtocolID() p2pprotocol.ID {
	return router.config.ProtocolID()
}

func (router DHTRouter) logger() types.Logger {
	return router.config.logger
}

func (router DHTRouter) buildSignedAnnouncement() (ann Announcement, err error) {
	addresses := router.aclHost.Addrs()
	sk := router.aclHost.Peerstore().PrivKey(router.aclHost.ID())

	ann.Addrs = append([]ma.Multiaddr{}, addresses...)
	ann.Pk = sk.GetPublic()
	ann.timestamp = time.Now().Unix()

	err = ann.SelfSign(sk)
	if err != nil {
		return ann, err
	}

	return ann, nil
}

func peerIdToDhtKey(id peer.ID) string {
	return fmt.Sprintf("/%s/%s", ValidatorNamespace, id.String())
}

func dhtKeyToPeerId(s string) (id peer.ID, err error) {
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

	router.logger().Debug("DHT: Lookup peer", types.LogFields{
		"peerID": peerId,
		"key":    key,
	})

	marshaled, err := router.dht.GetValue(ctx, key)
	if err != nil {
		return addr, err
	}

	var ann Announcement
	err = ann.UnmarshalJSON(marshaled)
	if err != nil {
		return addr, err
	}

	ok, err := ann.SelfVerify()
	if err != nil {
		return addr, err
	} else if !ok {
		return addr, InvalidSignature
	}

	addr.ID = peerId
	addr.Addrs = ann.Addrs

	router.logger().Debug("DHT: Found peer", types.LogFields{
		"peerID":        peerId,
		"key":           key,
		"foundPeerAddr": addr,
		"addrVersion":   ann.timestamp,
	})
	return addr, nil
}

func (router DHTRouter) publishHostAddr(ctx context.Context) error {
	ann, err := router.buildSignedAnnouncement()
	if err != nil {
		return err
	}

	marshaled, err := ann.MarshalJSON()
	if err != nil {
		return err
	}

	key := peerIdToDhtKey(router.aclHost.ID())
	router.logger().Debug("DHT: Put value", types.LogFields{
		"key":   key,
		"value": string(marshaled),
	})

	err = router.dht.PutValue(ctx, key, marshaled)
	if err != nil {
		return errors.Wrap(err, "could publish address")
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
					router.logger().Error("DHT: Error publishing address", types.LogFields{
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

func (router DHTRouter) printRoutingTableAndAclTable(interval time.Duration) {
	if !router.config.extendedDHTLogging {
		return
	}

	router.processes.RepeatWithCancel("periodical report", interval, router.ctx, func() {
		var rtString []string

		myKadId := kbucket.ConvertPeerID(router.dht.PeerID())
		rt := router.dht.RoutingTable()

		for _, p := range rt.ListPeers() {
			rtString = append(rtString, fmt.Sprintf("peerId=%s, kadId=%s,allowed=%t, cpl=%d",
				p.Pretty(),
				hex.EncodeToString(kbucket.ConvertPeerID(p)),
				router.aclHost.GetACL().IsAllowed(p, router.config.ProtocolID()),
				kbucket.CommonPrefixLen(kbucket.ConvertPeerID(p), myKadId)))
		}

		router.logger().Debug("DHT periodical report", types.LogFields{
			"protocolID": router.config.ProtocolID(),
			"acl":        router.aclHost.GetACL(),
			"rt":         rtString,
		})
	})
}

func (router DHTRouter) Close() error {
	router.ctxCancel()
	router.processes.Wait()

	return router.dht.Close()
}
