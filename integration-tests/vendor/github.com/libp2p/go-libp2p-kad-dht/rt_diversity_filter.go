package dht

import (
	"sync"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/libp2p/go-libp2p-kbucket/peerdiversity"

	ma "github.com/multiformats/go-multiaddr"
)

var _ peerdiversity.PeerIPGroupFilter = (*rtPeerIPGroupFilter)(nil)

type rtPeerIPGroupFilter struct {
	mu sync.RWMutex
	h  host.Host

	maxPerCpl   int
	maxForTable int

	cplIpGroupCount   map[int]map[peerdiversity.PeerIPGroupKey]int
	tableIpGroupCount map[peerdiversity.PeerIPGroupKey]int
}

// NewRTPeerDiversityFilter constructs the `PeerIPGroupFilter` that will be used to configure
// the diversity filter for the Routing Table.
// Please see the docs for `peerdiversity.PeerIPGroupFilter` AND `peerdiversity.Filter` for more details.
func NewRTPeerDiversityFilter(h host.Host, maxPerCpl, maxForTable int) *rtPeerIPGroupFilter {
	return &rtPeerIPGroupFilter{
		h: h,

		maxPerCpl:   maxPerCpl,
		maxForTable: maxForTable,

		cplIpGroupCount:   make(map[int]map[peerdiversity.PeerIPGroupKey]int),
		tableIpGroupCount: make(map[peerdiversity.PeerIPGroupKey]int),
	}

}

func (r *rtPeerIPGroupFilter) Allow(g peerdiversity.PeerGroupInfo) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := g.IPGroupKey
	cpl := g.Cpl

	if r.tableIpGroupCount[key] >= r.maxForTable {

		return false
	}

	c, ok := r.cplIpGroupCount[cpl]
	allow := !ok || c[key] < r.maxPerCpl
	return allow
}

func (r *rtPeerIPGroupFilter) Increment(g peerdiversity.PeerGroupInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := g.IPGroupKey
	cpl := g.Cpl

	r.tableIpGroupCount[key] = r.tableIpGroupCount[key] + 1
	if _, ok := r.cplIpGroupCount[cpl]; !ok {
		r.cplIpGroupCount[cpl] = make(map[peerdiversity.PeerIPGroupKey]int)
	}

	r.cplIpGroupCount[cpl][key] = r.cplIpGroupCount[cpl][key] + 1
}

func (r *rtPeerIPGroupFilter) Decrement(g peerdiversity.PeerGroupInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := g.IPGroupKey
	cpl := g.Cpl

	r.tableIpGroupCount[key] = r.tableIpGroupCount[key] - 1
	if r.tableIpGroupCount[key] == 0 {
		delete(r.tableIpGroupCount, key)
	}

	r.cplIpGroupCount[cpl][key] = r.cplIpGroupCount[cpl][key] - 1
	if r.cplIpGroupCount[cpl][key] == 0 {
		delete(r.cplIpGroupCount[cpl], key)
	}
	if len(r.cplIpGroupCount[cpl]) == 0 {
		delete(r.cplIpGroupCount, cpl)
	}
}

func (r *rtPeerIPGroupFilter) PeerAddresses(p peer.ID) []ma.Multiaddr {
	cs := r.h.Network().ConnsToPeer(p)
	addr := make([]ma.Multiaddr, 0, len(cs))
	for _, c := range cs {
		addr = append(addr, c.RemoteMultiaddr())
	}
	return addr
}
