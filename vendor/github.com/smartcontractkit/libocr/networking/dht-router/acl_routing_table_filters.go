package dhtrouter

import (
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/smartcontractkit/libocr/commontypes"
)

// ACLQueryFilter is a filter applied when considering peers to dial when querying
func ACLQueryFilter(acl ACL, protocol protocol.ID, logger commontypes.Logger) dht.QueryFilterFunc {
	if acl.IsACLEnforced(protocol) {
		return func(dht *dht.IpfsDHT, ai peer.AddrInfo) bool {
			b := acl.IsAllowed(ai.ID, protocol)
			if !b {
				logger.Warn("QueryFilter: denied", commontypes.LogFields{
					"remotePeerID":     ai.ID,
					"remoteMultiaddrs": ai,
					"id":               "DHT_ACL",
					"protocolID":       protocol,
				})
			}
			return b
		}
	}
	logger.Warn("QueryFilter: ACL disabled for this protocol", commontypes.LogFields{
		"id":         "DHT_ACL",
		"protocolID": protocol,
	})
	return func(dht *dht.IpfsDHT, ai peer.AddrInfo) bool {
		return false
	}
}

// ACLRoutingTableFilter is a filter applied when considering connections to keep in
// the local route table.
func ACLRoutingTableFilter(acl ACL, protocol protocol.ID, logger commontypes.Logger) dht.RouteTableFilterFunc {
	if acl.IsACLEnforced(protocol) {
		return func(dht *dht.IpfsDHT, conns []network.Conn) bool {
			for _, conn := range conns {
				b := acl.IsAllowed(conn.RemotePeer(), protocol)
				if !b {
					logger.Warn("RoutingTableFilter: denied", commontypes.LogFields{
						"remotePeerID":    conn.RemotePeer(),
						"remoteMultiaddr": conn.RemoteMultiaddr(),
						"id":              "DHT_ACL",
						"protocolID":      protocol,
					})
					return false
				}
			}
			return true
		}
	}
	// if ACL is not enforced for this protocol, always deny.
	logger.Warn("RoutingTableFilter: ACL disabled for this protocol", commontypes.LogFields{
		"id":         "DHT_ACL",
		"protocolID": protocol,
	})
	return func(dht *dht.IpfsDHT, conns []network.Conn) bool {
		return false
	}
}
