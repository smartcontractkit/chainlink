package nat

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	logging "github.com/ipfs/go-log"
	goprocess "github.com/jbenet/goprocess"
	periodic "github.com/jbenet/goprocess/periodic"
	nat "github.com/libp2p/go-nat"
)

var (
	// ErrNoMapping signals no mapping exists for an address
	ErrNoMapping = errors.New("mapping not established")
)

var log = logging.Logger("nat")

// MappingDuration is a default port mapping duration.
// Port mappings are renewed every (MappingDuration / 3)
const MappingDuration = time.Second * 60

// CacheTime is the time a mapping will cache an external address for
const CacheTime = time.Second * 15

// DiscoverNAT looks for a NAT device in the network and
// returns an object that can manage port mappings.
func DiscoverNAT(ctx context.Context) (*NAT, error) {
	var (
		natInstance nat.NAT
		err         error
	)

	done := make(chan struct{})
	go func() {
		defer close(done)
		// This will abort in 10 seconds anyways.
		natInstance, err = nat.DiscoverGateway()
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if err != nil {
		return nil, err
	}

	// Log the device addr.
	addr, err := natInstance.GetDeviceAddress()
	if err != nil {
		log.Debug("DiscoverGateway address error:", err)
	} else {
		log.Debug("DiscoverGateway address:", addr)
	}

	return newNAT(natInstance), nil
}

// NAT is an object that manages address port mappings in
// NATs (Network Address Translators). It is a long-running
// service that will periodically renew port mappings,
// and keep an up-to-date list of all the external addresses.
type NAT struct {
	natmu sync.Mutex
	nat   nat.NAT
	proc  goprocess.Process

	mappingmu sync.RWMutex // guards mappings
	mappings  map[*mapping]struct{}
}

func newNAT(realNAT nat.NAT) *NAT {
	return &NAT{
		nat:      realNAT,
		proc:     goprocess.WithParent(goprocess.Background()),
		mappings: make(map[*mapping]struct{}),
	}
}

// Close shuts down all port mappings. NAT can no longer be used.
func (nat *NAT) Close() error {
	return nat.proc.Close()
}

// Process returns the nat's life-cycle manager, for making it listen
// to close signals.
func (nat *NAT) Process() goprocess.Process {
	return nat.proc
}

// Mappings returns a slice of all NAT mappings
func (nat *NAT) Mappings() []Mapping {
	nat.mappingmu.Lock()
	maps2 := make([]Mapping, 0, len(nat.mappings))
	for m := range nat.mappings {
		maps2 = append(maps2, m)
	}
	nat.mappingmu.Unlock()
	return maps2
}

func (nat *NAT) addMapping(m *mapping) {
	// make mapping automatically close when nat is closed.
	nat.proc.AddChild(m.proc)

	nat.mappingmu.Lock()
	nat.mappings[m] = struct{}{}
	nat.mappingmu.Unlock()
}

func (nat *NAT) rmMapping(m *mapping) {
	nat.mappingmu.Lock()
	delete(nat.mappings, m)
	nat.mappingmu.Unlock()
}

// NewMapping attempts to construct a mapping on protocol and internal port
// It will also periodically renew the mapping until the returned Mapping
// -- or its parent NAT -- is Closed.
//
// May not succeed, and mappings may change over time;
// NAT devices may not respect our port requests, and even lie.
// Clients should not store the mapped results, but rather always
// poll our object for the latest mappings.
func (nat *NAT) NewMapping(protocol string, port int) (Mapping, error) {
	if nat == nil {
		return nil, fmt.Errorf("no nat available")
	}

	switch protocol {
	case "tcp", "udp":
	default:
		return nil, fmt.Errorf("invalid protocol: %s", protocol)
	}

	m := &mapping{
		intport: port,
		nat:     nat,
		proto:   protocol,
	}

	m.proc = goprocess.WithTeardown(func() error {
		nat.rmMapping(m)
		nat.natmu.Lock()
		defer nat.natmu.Unlock()
		nat.nat.DeletePortMapping(m.Protocol(), m.InternalPort())
		return nil
	})

	nat.addMapping(m)

	m.proc.AddChild(periodic.Every(MappingDuration/3, func(worker goprocess.Process) {
		nat.establishMapping(m)
	}))

	// do it once synchronously, so first mapping is done right away, and before exiting,
	// allowing users -- in the optimistic case -- to use results right after.
	nat.establishMapping(m)
	return m, nil
}

func (nat *NAT) establishMapping(m *mapping) {
	oldport := m.ExternalPort()

	log.Debugf("Attempting port map: %s/%d", m.Protocol(), m.InternalPort())
	comment := "libp2p"

	nat.natmu.Lock()
	newport, err := nat.nat.AddPortMapping(m.Protocol(), m.InternalPort(), comment, MappingDuration)
	if err != nil {
		// Some hardware does not support mappings with timeout, so try that
		newport, err = nat.nat.AddPortMapping(m.Protocol(), m.InternalPort(), comment, 0)
	}
	nat.natmu.Unlock()

	if err != nil || newport == 0 {
		m.setExternalPort(0) // clear mapping
		// TODO: log.Event
		log.Warnf("failed to establish port mapping: %s", err)
		// we do not close if the mapping failed,
		// because it may work again next time.
		return
	}

	m.setExternalPort(newport)
	log.Debugf("NAT Mapping: %d --> %d (%s)", m.ExternalPort(), m.InternalPort(), m.Protocol())
	if oldport != 0 && newport != oldport {
		log.Debugf("failed to renew same port mapping: ch %d -> %d", oldport, newport)
	}
}
