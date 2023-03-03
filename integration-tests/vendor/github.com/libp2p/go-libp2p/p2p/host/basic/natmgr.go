package basichost

import (
	"net"
	"strconv"
	"sync"

	goprocess "github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"
	"github.com/libp2p/go-libp2p-core/network"
	inat "github.com/libp2p/go-libp2p-nat"
	ma "github.com/multiformats/go-multiaddr"
)

// A simple interface to manage NAT devices.
type NATManager interface {

	// Get the NAT device managed by the NAT manager.
	NAT() *inat.NAT

	// Receive a notification when the NAT device is ready for use.
	Ready() <-chan struct{}

	// Close all resources associated with a NAT manager.
	Close() error
}

// Create a NAT manager.
func NewNATManager(net network.Network) NATManager {
	return newNatManager(net)
}

// natManager takes care of adding + removing port mappings to the nat.
// Initialized with the host if it has a NATPortMap option enabled.
// natManager receives signals from the network, and check on nat mappings:
//  * natManager listens to the network and adds or closes port mappings
//    as the network signals Listen() or ListenClose().
//  * closing the natManager closes the nat and its mappings.
type natManager struct {
	net   network.Network
	natmu sync.RWMutex
	nat   *inat.NAT

	ready    chan struct{} // closed once the nat is ready to process port mappings
	syncFlag chan struct{}

	proc goprocess.Process // natManager has a process + children. can be closed.
}

func newNatManager(net network.Network) *natManager {
	nmgr := &natManager{
		net:      net,
		ready:    make(chan struct{}),
		syncFlag: make(chan struct{}, 1),
	}

	nmgr.proc = goprocess.WithParent(goprocess.Background())

	nmgr.start()
	return nmgr
}

// Close closes the natManager, closing the underlying nat
// and unregistering from network events.
func (nmgr *natManager) Close() error {
	return nmgr.proc.Close()
}

// Ready returns a channel which will be closed when the NAT has been found
// and is ready to be used, or the search process is done.
func (nmgr *natManager) Ready() <-chan struct{} {
	return nmgr.ready
}

func (nmgr *natManager) start() {
	nmgr.proc.Go(func(worker goprocess.Process) {
		// inat.DiscoverNAT blocks until the nat is found or a timeout
		// is reached. we unfortunately cannot specify timeouts-- the
		// library we're using just blocks.
		//
		// Note: on early shutdown, there may be a case where we're trying
		// to close before DiscoverNAT() returns. Since we cant cancel it
		// (library) we can choose to (1) drop the result and return early,
		// or (2) wait until it times out to exit. For now we choose (2),
		// to avoid leaking resources in a non-obvious way. the only case
		// this affects is when the daemon is being started up and _immediately_
		// asked to close. other services are also starting up, so ok to wait.

		natInstance, err := inat.DiscoverNAT(goprocessctx.OnClosingContext(worker))
		if err != nil {
			log.Info("DiscoverNAT error:", err)
			close(nmgr.ready)
			return
		}

		nmgr.natmu.Lock()
		nmgr.nat = natInstance
		nmgr.natmu.Unlock()
		close(nmgr.ready)

		// wire up the nat to close when nmgr closes.
		// nmgr.proc is our parent, and waiting for us.
		nmgr.proc.AddChild(nmgr.nat.Process())

		// sign natManager up for network notifications
		// we need to sign up here to avoid missing some notifs
		// before the NAT has been found.
		nmgr.net.Notify((*nmgrNetNotifiee)(nmgr))
		defer nmgr.net.StopNotify((*nmgrNetNotifiee)(nmgr))

		nmgr.doSync() // sync one first.
		for {
			select {
			case <-nmgr.syncFlag:
				nmgr.doSync() // sync when our listen addresses chnage.
			case <-worker.Closing():
				return
			}
		}
	})
}

func (nmgr *natManager) sync() {
	select {
	case nmgr.syncFlag <- struct{}{}:
	default:
	}
}

// doSync syncs the current NAT mappings, removing any outdated mappings and adding any
// new mappings.
func (nmgr *natManager) doSync() {
	ports := map[string]map[int]bool{
		"tcp": map[int]bool{},
		"udp": map[int]bool{},
	}
	for _, maddr := range nmgr.net.ListenAddresses() {
		// Strip the IP
		maIP, rest := ma.SplitFirst(maddr)
		if maIP == nil || rest == nil {
			continue
		}

		switch maIP.Protocol().Code {
		case ma.P_IP6, ma.P_IP4:
		default:
			continue
		}

		// Only bother if we're listening on a
		// unicast/unspecified IP.
		ip := net.IP(maIP.RawValue())
		if !(ip.IsGlobalUnicast() || ip.IsUnspecified()) {
			continue
		}

		// Extract the port/protocol
		proto, _ := ma.SplitFirst(rest)
		if proto == nil {
			continue
		}

		var protocol string
		switch proto.Protocol().Code {
		case ma.P_TCP:
			protocol = "tcp"
		case ma.P_UDP:
			protocol = "udp"
		default:
			continue
		}

		port, err := strconv.ParseUint(proto.Value(), 10, 16)
		if err != nil {
			// bug in multiaddr
			panic(err)
		}
		ports[protocol][int(port)] = false
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	// Close old mappings
	for _, m := range nmgr.nat.Mappings() {
		mappedPort := m.InternalPort()
		if _, ok := ports[m.Protocol()][mappedPort]; !ok {
			// No longer need this mapping.
			wg.Add(1)
			go func(m inat.Mapping) {
				defer wg.Done()
				m.Close()
			}(m)
		} else {
			// already mapped
			ports[m.Protocol()][mappedPort] = true
		}
	}

	// Create new mappings.
	for proto, pports := range ports {
		for port, mapped := range pports {
			if mapped {
				continue
			}
			wg.Add(1)
			go func(proto string, port int) {
				defer wg.Done()
				_, err := nmgr.nat.NewMapping(proto, port)
				if err != nil {
					log.Errorf("failed to port-map %s port %d: %s", proto, port, err)
				}
			}(proto, port)
		}
	}
}

// NAT returns the natManager's nat object. this may be nil, if
// (a) the search process is still ongoing, or (b) the search process
// found no nat. Clients must check whether the return value is nil.
func (nmgr *natManager) NAT() *inat.NAT {
	nmgr.natmu.Lock()
	defer nmgr.natmu.Unlock()
	return nmgr.nat
}

type nmgrNetNotifiee natManager

func (nn *nmgrNetNotifiee) natManager() *natManager {
	return (*natManager)(nn)
}

func (nn *nmgrNetNotifiee) Listen(n network.Network, addr ma.Multiaddr) {
	nn.natManager().sync()
}

func (nn *nmgrNetNotifiee) ListenClose(n network.Network, addr ma.Multiaddr) {
	nn.natManager().sync()
}

func (nn *nmgrNetNotifiee) Connected(network.Network, network.Conn)      {}
func (nn *nmgrNetNotifiee) Disconnected(network.Network, network.Conn)   {}
func (nn *nmgrNetNotifiee) OpenedStream(network.Network, network.Stream) {}
func (nn *nmgrNetNotifiee) ClosedStream(network.Network, network.Stream) {}
