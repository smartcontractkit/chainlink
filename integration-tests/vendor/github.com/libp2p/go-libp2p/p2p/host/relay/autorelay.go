package relay

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/event"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"

	circuit "github.com/libp2p/go-libp2p-circuit"
	discovery "github.com/libp2p/go-libp2p-discovery"
	basic "github.com/libp2p/go-libp2p/p2p/host/basic"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

const (
	RelayRendezvous = "/libp2p/relay"
)

var (
	DesiredRelays = 1

	BootDelay = 20 * time.Second
)

// These are the known PL-operated relays
var DefaultRelays = []string{
	"/ip4/147.75.80.110/tcp/4001/p2p/QmbFgm5zan8P6eWWmeyfncR5feYEMPbht5b1FW1C37aQ7y",
	"/ip4/147.75.195.153/tcp/4001/p2p/QmW9m57aiBDHAkKj9nmFSEn7ZqrcF1fZS4bipsTCHburei",
	"/ip4/147.75.70.221/tcp/4001/p2p/Qme8g49gm3q4Acp7xWBKg3nAa9fxZ1YmyDJdyGgoG6LsXh",
}

// AutoRelay is a Host that uses relays for connectivity when a NAT is detected.
type AutoRelay struct {
	host     *basic.BasicHost
	discover discovery.Discoverer
	router   routing.PeerRouting
	addrsF   basic.AddrsFactory

	static []peer.AddrInfo

	disconnect chan struct{}

	mx     sync.Mutex
	relays map[peer.ID]struct{}
	status network.Reachability

	cachedAddrs       []ma.Multiaddr
	cachedAddrsExpiry time.Time
}

func NewAutoRelay(ctx context.Context, bhost *basic.BasicHost, discover discovery.Discoverer, router routing.PeerRouting, static []peer.AddrInfo) *AutoRelay {
	ar := &AutoRelay{
		host:       bhost,
		discover:   discover,
		router:     router,
		addrsF:     bhost.AddrsFactory,
		static:     static,
		relays:     make(map[peer.ID]struct{}),
		disconnect: make(chan struct{}, 1),
		status:     network.ReachabilityUnknown,
	}
	bhost.AddrsFactory = ar.hostAddrs
	bhost.Network().Notify(ar)
	go ar.background(ctx)
	return ar
}

func (ar *AutoRelay) hostAddrs(addrs []ma.Multiaddr) []ma.Multiaddr {
	return ar.relayAddrs(ar.addrsF(addrs))
}

func (ar *AutoRelay) background(ctx context.Context) {
	subReachability, _ := ar.host.EventBus().Subscribe(new(event.EvtLocalReachabilityChanged))
	defer subReachability.Close()

	// when true, we need to identify push
	push := false

	for {
		select {
		case ev, ok := <-subReachability.Out():
			if !ok {
				return
			}
			evt, ok := ev.(event.EvtLocalReachabilityChanged)
			if !ok {
				return
			}

			var update bool
			if evt.Reachability == network.ReachabilityPrivate {
				// TODO: this is a long-lived (2.5min task) that should get spun up in a separate thread
				// and canceled if the relay learns the nat is now public.
				update = ar.findRelays(ctx)
			}

			ar.mx.Lock()
			if update || (ar.status != evt.Reachability && evt.Reachability != network.ReachabilityUnknown) {
				push = true
			}
			ar.status = evt.Reachability
			ar.mx.Unlock()
		case <-ar.disconnect:
			push = true
		case <-ctx.Done():
			return
		}

		if push {
			ar.mx.Lock()
			ar.cachedAddrs = nil
			ar.mx.Unlock()
			push = false
			ar.host.SignalAddressChange()
		}
	}
}

func (ar *AutoRelay) findRelays(ctx context.Context) bool {
	if ar.numRelays() >= DesiredRelays {
		return false
	}

	update := false
	for retry := 0; retry < 5; retry++ {
		if retry > 0 {
			log.Debug("no relays connected; retrying in 30s")
			select {
			case <-time.After(30 * time.Second):
			case <-ctx.Done():
				return update
			}
		}

		update = ar.findRelaysOnce(ctx) || update
		if ar.numRelays() > 0 {
			return update
		}
	}
	return update
}

func (ar *AutoRelay) findRelaysOnce(ctx context.Context) bool {
	pis, err := ar.discoverRelays(ctx)
	if err != nil {
		log.Debugf("error discovering relays: %s", err)
		return false
	}
	log.Debugf("discovered %d relays", len(pis))
	pis = ar.selectRelays(ctx, pis)
	log.Debugf("selected %d relays", len(pis))

	update := false
	for _, pi := range pis {
		update = ar.tryRelay(ctx, pi) || update
		if ar.numRelays() >= DesiredRelays {
			break
		}
	}
	return update
}

func (ar *AutoRelay) numRelays() int {
	ar.mx.Lock()
	defer ar.mx.Unlock()
	return len(ar.relays)
}

// usingRelay returns if we're currently using the given relay.
func (ar *AutoRelay) usingRelay(p peer.ID) bool {
	ar.mx.Lock()
	defer ar.mx.Unlock()
	_, ok := ar.relays[p]
	return ok
}

// addRelay adds the given relay to our set of relays.
// returns true when we add a new relay
func (ar *AutoRelay) tryRelay(ctx context.Context, pi peer.AddrInfo) bool {
	if ar.usingRelay(pi.ID) {
		return false
	}

	if !ar.connect(ctx, pi) {
		return false
	}

	ok, err := circuit.CanHop(ctx, ar.host, pi.ID)
	if err != nil {
		log.Debugf("error querying relay: %s", err.Error())
		return false
	}

	if !ok {
		// not a hop relay
		return false
	}

	ar.mx.Lock()
	defer ar.mx.Unlock()

	// make sure we're still connected.
	if ar.host.Network().Connectedness(pi.ID) != network.Connected {
		return false
	}
	ar.relays[pi.ID] = struct{}{}

	return true
}

func (ar *AutoRelay) connect(ctx context.Context, pi peer.AddrInfo) bool {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	if len(pi.Addrs) == 0 {
		var err error
		pi, err = ar.router.FindPeer(ctx, pi.ID)
		if err != nil {
			log.Debugf("error finding relay peer %s: %s", pi.ID, err.Error())
			return false
		}
	}

	err := ar.host.Connect(ctx, pi)
	if err != nil {
		log.Debugf("error connecting to relay %s: %s", pi.ID, err.Error())
		return false
	}

	// tag the connection as very important
	ar.host.ConnManager().TagPeer(pi.ID, "relay", 42)
	return true
}

func (ar *AutoRelay) discoverRelays(ctx context.Context) ([]peer.AddrInfo, error) {
	if len(ar.static) > 0 {
		return ar.static, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return discovery.FindPeers(ctx, ar.discover, RelayRendezvous, discovery.Limit(1000))
}

func (ar *AutoRelay) selectRelays(ctx context.Context, pis []peer.AddrInfo) []peer.AddrInfo {
	// TODO better relay selection strategy; this just selects random relays
	//      but we should probably use ping latency as the selection metric

	shuffleRelays(pis)
	return pis
}

// This function is computes the NATed relay addrs when our status is private:
// - The public addrs are removed from the address set.
// - The non-public addrs are included verbatim so that peers behind the same NAT/firewall
//   can still dial us directly.
// - On top of those, we add the relay-specific addrs for the relays to which we are
//   connected. For each non-private relay addr, we encapsulate the p2p-circuit addr
//   through which we can be dialed.
func (ar *AutoRelay) relayAddrs(addrs []ma.Multiaddr) []ma.Multiaddr {
	ar.mx.Lock()
	defer ar.mx.Unlock()

	if ar.status != network.ReachabilityPrivate {
		return addrs
	}

	if ar.cachedAddrs != nil && time.Now().Before(ar.cachedAddrsExpiry) {
		return ar.cachedAddrs
	}

	raddrs := make([]ma.Multiaddr, 0, 4*len(ar.relays)+4)

	// only keep private addrs from the original addr set
	for _, addr := range addrs {
		if manet.IsPrivateAddr(addr) {
			raddrs = append(raddrs, addr)
		}
	}

	// add relay specific addrs to the list
	for p := range ar.relays {
		addrs := cleanupAddressSet(ar.host.Peerstore().Addrs(p))

		circuit, err := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s/p2p-circuit", p.Pretty()))
		if err != nil {
			panic(err)
		}

		for _, addr := range addrs {
			pub := addr.Encapsulate(circuit)
			raddrs = append(raddrs, pub)
		}
	}

	ar.cachedAddrs = raddrs
	ar.cachedAddrsExpiry = time.Now().Add(30 * time.Second)

	return raddrs
}

func shuffleRelays(pis []peer.AddrInfo) {
	for i := range pis {
		j := rand.Intn(i + 1)
		pis[i], pis[j] = pis[j], pis[i]
	}
}

// Notifee
func (ar *AutoRelay) Listen(network.Network, ma.Multiaddr)      {}
func (ar *AutoRelay) ListenClose(network.Network, ma.Multiaddr) {}
func (ar *AutoRelay) Connected(network.Network, network.Conn)   {}

func (ar *AutoRelay) Disconnected(net network.Network, c network.Conn) {
	p := c.RemotePeer()

	ar.mx.Lock()
	defer ar.mx.Unlock()

	if ar.host.Network().Connectedness(p) == network.Connected {
		// We have a second connection.
		return
	}

	if _, ok := ar.relays[p]; ok {
		delete(ar.relays, p)
		select {
		case ar.disconnect <- struct{}{}:
		default:
		}
	}
}

func (ar *AutoRelay) OpenedStream(network.Network, network.Stream) {}
func (ar *AutoRelay) ClosedStream(network.Network, network.Stream) {}
