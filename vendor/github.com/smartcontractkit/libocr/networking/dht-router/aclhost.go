package dhtrouter

import (
	"context"

	"github.com/libp2p/go-libp2p-core/connmgr"
	"github.com/libp2p/go-libp2p-core/event"
	"github.com/libp2p/go-libp2p-core/host"
	p2pnetwork "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/protocol"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
)

type ACL interface {
	// if an ID in the white list of a protocol
	IsAllowed(id peer.ID, protocol protocol.ID) bool

	// if ACL is enforced for a protocol
	IsACLEnforced(protocol protocol.ID) bool

	String() string
}

type ACLHost interface {
	host.Host

	SetACL(acl ACL)
	GetACL() ACL
}

type BasicACLHost struct {
	host   host.Host
	acl    ACL
	logger loghelper.LoggerWithContext
}

func WrapACL(h host.Host, acl ACL, logger loghelper.LoggerWithContext) ACLHost {
	return &BasicACLHost{host: h, acl: acl, logger: logger}
}

func (aclHost *BasicACLHost) SetACL(acl ACL) {
	aclHost.acl = acl
}

func (aclHost BasicACLHost) GetACL() ACL {
	return aclHost.acl
}

func (aclHost BasicACLHost) ID() peer.ID {
	return aclHost.host.ID()
}

func (aclHost BasicACLHost) Peerstore() peerstore.Peerstore {
	return aclHost.host.Peerstore()
}

func (aclHost BasicACLHost) Addrs() []ma.Multiaddr {
	return aclHost.host.Addrs()
}

func (aclHost BasicACLHost) Network() p2pnetwork.Network {
	return aclHost.host.Network()
}

func (aclHost BasicACLHost) Mux() protocol.Switch {
	// code should never reach here
	panic("direct access to Mux is not allowed")
}

func (aclHost BasicACLHost) Connect(ctx context.Context, pi peer.AddrInfo) error {
	return aclHost.host.Connect(ctx, pi)
}

func (aclHost BasicACLHost) SetStreamHandler(protocol protocol.ID, handler p2pnetwork.StreamHandler) {
	aclHost.logger.Debug("ACLHost: setting stream handler", commontypes.LogFields{
		"id":         "DHT_ACL",
		"protocolID": protocol,
	})

	wrapped := func(stream p2pnetwork.Stream) {
		if !aclHost.acl.IsAllowed(stream.Conn().RemotePeer(), protocol) {
			aclHost.logger.Warn("ACLHost: denied stream", commontypes.LogFields{
				"id":              "DHT_ACL",
				"protocolID":      protocol,
				"remotePeerID":    stream.Conn().RemotePeer(),
				"remoteMultiaddr": stream.Conn().RemoteMultiaddr(),
			})
			if err := stream.Reset(); err != nil {
				aclHost.logger.Error("ACLHost: Could not reset stream", commontypes.LogFields{
					"id":              "DHT_ACL",
					"protocolID":      protocol,
					"remotePeerID":    stream.Conn().RemotePeer(),
					"remoteMultiaddr": stream.Conn().RemoteMultiaddr(),
					"err":             err.Error(),
				})
			}

			return
		}
		handler(stream)

	}

	if aclHost.acl.IsACLEnforced(protocol) {
		aclHost.logger.Debug("ACLHost: Wrapping ACL", commontypes.LogFields{
			"id":         "DHT_ACL",
			"protocolID": protocol,
		})
		aclHost.host.SetStreamHandler(protocol, wrapped)
	} else {
		aclHost.logger.Debug("ACLHost: ACL not enforced for this protocol", commontypes.LogFields{
			"id":         "DHT_ACL",
			"protocolID": protocol,
		})
		aclHost.host.SetStreamHandler(protocol, handler)
	}
}

func (aclHost BasicACLHost) SetStreamHandlerMatch(id protocol.ID, f func(string) bool, handler p2pnetwork.StreamHandler) {
	// code should never reach here
	panic("SetStreamHandlerMatch not allowed")
}

func (aclHost BasicACLHost) RemoveStreamHandler(pid protocol.ID) {
	aclHost.host.RemoveStreamHandler(pid)
}

func (aclHost BasicACLHost) NewStream(ctx context.Context, p peer.ID, pids ...protocol.ID) (p2pnetwork.Stream, error) {
	var allowdPids []protocol.ID
	for _, pid := range pids {
		if aclHost.acl.IsAllowed(p, pid) {
			allowdPids = append(allowdPids, pid)
		} else {
			aclHost.logger.Warn("ACLHost: Denying NewStream", commontypes.LogFields{
				"id":           "DHT_ACL",
				"protocolID":   pid,
				"remotePeerID": p,
			})
		}
	}
	return aclHost.host.NewStream(ctx, p, allowdPids...)
}

func (aclHost BasicACLHost) Close() error {
	return aclHost.host.Close()
}

func (aclHost BasicACLHost) ConnManager() connmgr.ConnManager {
	return aclHost.host.ConnManager()
}

func (aclHost BasicACLHost) EventBus() event.Bus {
	return aclHost.host.EventBus()
}
