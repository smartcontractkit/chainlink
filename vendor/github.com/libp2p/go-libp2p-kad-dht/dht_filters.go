package dht

import (
	"bytes"
	"net"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/google/gopacket/routing"
	netroute "github.com/libp2p/go-netroute"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

// QueryFilterFunc is a filter applied when considering peers to dial when querying
type QueryFilterFunc func(dht *IpfsDHT, ai peer.AddrInfo) bool

// RouteTableFilterFunc is a filter applied when considering connections to keep in
// the local route table.
type RouteTableFilterFunc func(dht *IpfsDHT, conns []network.Conn) bool

var publicCIDR6 = "2000::/3"
var public6 *net.IPNet

func init() {
	_, public6, _ = net.ParseCIDR(publicCIDR6)
}

// isPublicAddr follows the logic of manet.IsPublicAddr, except it uses
// a stricter definition of "public" for ipv6: namely "is it in 2000::/3"?
func isPublicAddr(a ma.Multiaddr) bool {
	ip, err := manet.ToIP(a)
	if err != nil {
		return false
	}
	if ip.To4() != nil {
		return !inAddrRange(ip, manet.Private4) && !inAddrRange(ip, manet.Unroutable4)
	}

	return public6.Contains(ip)
}

// isPrivateAddr follows the logic of manet.IsPrivateAddr, except that
// it uses a stricter definition of "public" for ipv6
func isPrivateAddr(a ma.Multiaddr) bool {
	ip, err := manet.ToIP(a)
	if err != nil {
		return false
	}
	if ip.To4() != nil {
		return inAddrRange(ip, manet.Private4)
	}

	return !public6.Contains(ip) && !inAddrRange(ip, manet.Unroutable6)
}

// PublicQueryFilter returns true if the peer is suspected of being publicly accessible
func PublicQueryFilter(_ *IpfsDHT, ai peer.AddrInfo) bool {
	if len(ai.Addrs) == 0 {
		return false
	}

	var hasPublicAddr bool
	for _, a := range ai.Addrs {
		if !isRelayAddr(a) && isPublicAddr(a) {
			hasPublicAddr = true
		}
	}
	return hasPublicAddr
}

var _ QueryFilterFunc = PublicQueryFilter

// PublicRoutingTableFilter allows a peer to be added to the routing table if the connections to that peer indicate
// that it is on a public network
func PublicRoutingTableFilter(dht *IpfsDHT, conns []network.Conn) bool {
	if len(conns) == 0 {
		return false
	}

	// Do we have a public address for this peer?
	id := conns[0].RemotePeer()
	known := dht.peerstore.PeerInfo(id)
	for _, a := range known.Addrs {
		if !isRelayAddr(a) && isPublicAddr(a) {
			return true
		}
	}

	return false
}

var _ RouteTableFilterFunc = PublicRoutingTableFilter

// PrivateQueryFilter doens't currently restrict which peers we are willing to query from the local DHT.
func PrivateQueryFilter(dht *IpfsDHT, ai peer.AddrInfo) bool {
	return len(ai.Addrs) > 0
}

var _ QueryFilterFunc = PrivateQueryFilter

// We call this very frequently but routes can technically change at runtime.
// Cache it for two minutes.
const routerCacheTime = 2 * time.Minute

var routerCache struct {
	sync.RWMutex
	router  routing.Router
	expires time.Time
}

func getCachedRouter() routing.Router {
	routerCache.RLock()
	router := routerCache.router
	expires := routerCache.expires
	routerCache.RUnlock()

	if time.Now().Before(expires) {
		return router
	}

	routerCache.Lock()
	defer routerCache.Unlock()

	now := time.Now()
	if now.Before(routerCache.expires) {
		return router
	}
	routerCache.router, _ = netroute.New()
	routerCache.expires = now.Add(routerCacheTime)
	return router
}

// PrivateRoutingTableFilter allows a peer to be added to the routing table if the connections to that peer indicate
// that it is on a private network
func PrivateRoutingTableFilter(dht *IpfsDHT, conns []network.Conn) bool {
	router := getCachedRouter()
	myAdvertisedIPs := make([]net.IP, 0)
	for _, a := range dht.Host().Addrs() {
		if isPublicAddr(a) && !isRelayAddr(a) {
			ip, err := manet.ToIP(a)
			if err != nil {
				continue
			}
			myAdvertisedIPs = append(myAdvertisedIPs, ip)
		}
	}

	for _, c := range conns {
		ra := c.RemoteMultiaddr()
		if isPrivateAddr(ra) && !isRelayAddr(ra) {
			return true
		}

		if isPublicAddr(ra) {
			ip, err := manet.ToIP(ra)
			if err != nil {
				continue
			}

			// if the ip is the same as one of the local host's public advertised IPs - then consider it local
			for _, i := range myAdvertisedIPs {
				if i.Equal(ip) {
					return true
				}
				if ip.To4() == nil {
					if i.To4() == nil && isEUI(ip) && sameV6Net(i, ip) {
						return true
					}
				}
			}

			// if there's no gateway - a direct host in the OS routing table - then consider it local
			// This is relevant in particular to ipv6 networks where the addresses may all be public,
			// but the nodes are aware of direct links between each other.
			if router != nil {
				_, gw, _, err := router.Route(ip)
				if gw == nil && err == nil {
					return true
				}
			}
		}
	}

	return false
}

var _ RouteTableFilterFunc = PrivateRoutingTableFilter

func isEUI(ip net.IP) bool {
	// per rfc 2373
	return len(ip) == net.IPv6len && ip[11] == 0xff && ip[12] == 0xfe
}

func sameV6Net(a, b net.IP) bool {
	return len(a) == net.IPv6len && len(b) == net.IPv6len && bytes.Equal(a[0:8], b[0:8]) //nolint
}

func isRelayAddr(a ma.Multiaddr) bool {
	found := false
	ma.ForEach(a, func(c ma.Component) bool {
		found = c.Protocol().Code == ma.P_CIRCUIT
		return !found
	})
	return found
}

func inAddrRange(ip net.IP, ipnets []*net.IPNet) bool {
	for _, ipnet := range ipnets {
		if ipnet.Contains(ip) {
			return true
		}
	}

	return false
}
