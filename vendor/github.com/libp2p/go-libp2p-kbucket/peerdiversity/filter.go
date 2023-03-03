package peerdiversity

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"sync"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/libp2p/go-cidranger"
	asnutil "github.com/libp2p/go-libp2p-asn-util"

	logging "github.com/ipfs/go-log"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

var dfLog = logging.Logger("diversityFilter")

type asnStore interface {
	AsnForIPv6(ip net.IP) (string, error)
}

// PeerIPGroupKey is a unique key that represents ONE of the IP Groups the peer belongs to.
// A peer has one PeerIPGroupKey per address. Thus, a peer can belong to MULTIPLE Groups if it has
// multiple addresses.
// For now, given a peer address, our grouping mechanism is as follows:
// 1. For IPv6 addresses, we group by the ASN of the IP address.
// 2. For IPv4 addresses, all addresses that belong to same legacy (Class A)/8 allocations
//    OR share the same /16 prefix are in the same group.
type PeerIPGroupKey string

// https://en.wikipedia.org/wiki/List_of_assigned_/8_IPv4_address_blocks
var legacyClassA = []string{"12.0.0.0/8", "17.0.0.0/8", "19.0.0.0/8", "38.0.0.0/8", "48.0.0.0/8", "56.0.0.0/8", "73.0.0.0/8", "53.0.0.0/8"}

// PeerGroupInfo represents the grouping info for a Peer.
type PeerGroupInfo struct {
	Id         peer.ID
	Cpl        int
	IPGroupKey PeerIPGroupKey
}

// PeerIPGroupFilter is the interface that must be implemented by callers who want to
// instantiate a `peerdiversity.Filter`. This interface provides the function hooks
// that are used/called by the `peerdiversity.Filter`.
type PeerIPGroupFilter interface {
	// Allow is called by the Filter to test if a peer with the given
	// grouping info should be allowed/rejected by the Filter. This will be called ONLY
	// AFTER the peer has successfully passed all of the Filter's internal checks.
	// Note: If the peer is whitelisted on the Filter, the peer will be allowed by the Filter without calling this function.
	Allow(PeerGroupInfo) (allow bool)

	// Increment is called by the Filter when a peer with the given Grouping Info.
	// is added to the Filter state. This will happen after the peer has passed
	// all of the Filter's internal checks and the Allow function defined above for all of it's Groups.
	Increment(PeerGroupInfo)

	// Decrement is called by the Filter when a peer with the given
	// Grouping Info is removed from the Filter. This will happen when the caller/user of the Filter
	// no longer wants the peer and the IP groups it belongs to to count towards the Filter state.
	Decrement(PeerGroupInfo)

	// PeerAddresses is called by the Filter to determine the addresses of the given peer
	// it should use to determine the IP groups it belongs to.
	PeerAddresses(peer.ID) []ma.Multiaddr
}

// Filter is a peer diversity filter that accepts or rejects peers based on the whitelisting rules configured
// AND the diversity policies defined by the implementation of the PeerIPGroupFilter interface
// passed to it.
type Filter struct {
	mu sync.Mutex
	// An implementation of the `PeerIPGroupFilter` interface defined above.
	pgm        PeerIPGroupFilter
	peerGroups map[peer.ID][]PeerGroupInfo

	// whitelisted peers
	wlpeers map[peer.ID]struct{}

	// legacy IPv4 Class A networks.
	legacyCidrs cidranger.Ranger

	logKey string

	cplFnc func(peer.ID) int

	cplPeerGroups map[int]map[peer.ID][]PeerIPGroupKey

	asnStore asnStore
}

// NewFilter creates a Filter for Peer Diversity.
func NewFilter(pgm PeerIPGroupFilter, logKey string, cplFnc func(peer.ID) int) (*Filter, error) {
	if pgm == nil {
		return nil, errors.New("peergroup implementation can not be nil")
	}

	// Crate a Trie for legacy Class N networks
	legacyCidrs := cidranger.NewPCTrieRanger()
	for _, cidr := range legacyClassA {
		_, nn, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		if err := legacyCidrs.Insert(cidranger.NewBasicRangerEntry(*nn)); err != nil {
			return nil, err
		}
	}

	return &Filter{
		pgm:           pgm,
		peerGroups:    make(map[peer.ID][]PeerGroupInfo),
		wlpeers:       make(map[peer.ID]struct{}),
		legacyCidrs:   legacyCidrs,
		logKey:        logKey,
		cplFnc:        cplFnc,
		cplPeerGroups: make(map[int]map[peer.ID][]PeerIPGroupKey),
		asnStore:      asnutil.Store,
	}, nil
}

func (f *Filter) Remove(p peer.ID) {
	f.mu.Lock()
	defer f.mu.Unlock()

	cpl := f.cplFnc(p)

	for _, info := range f.peerGroups[p] {
		f.pgm.Decrement(info)
	}
	f.peerGroups[p] = nil
	delete(f.peerGroups, p)
	delete(f.cplPeerGroups[cpl], p)

	if len(f.cplPeerGroups[cpl]) == 0 {
		delete(f.cplPeerGroups, cpl)
	}
}

// TryAdd attempts to add the peer to the Filter state and returns true if it's successful, false otherwise.
func (f *Filter) TryAdd(p peer.ID) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.wlpeers[p]; ok {
		return true
	}

	cpl := f.cplFnc(p)

	// don't allow peers for which we can't determine addresses.
	addrs := f.pgm.PeerAddresses(p)
	if len(addrs) == 0 {
		dfLog.Debugw("no addresses found for peer", "appKey", f.logKey, "peer", p.Pretty())
		return false
	}

	peerGroups := make([]PeerGroupInfo, 0, len(addrs))
	for _, a := range addrs {
		ip, err := manet.ToIP(a)
		if err != nil {
			dfLog.Errorw("failed to parse IP from multiaddr", "appKey", f.logKey,
				"multiaddr", a.String(), "err", err)
			return false
		}

		// reject the peer if we can't determine a grouping for one of it's address.
		key, err := f.ipGroupKey(ip)
		if err != nil {
			dfLog.Errorw("failed to find Group Key", "appKey", f.logKey, "ip", ip.String(), "peer", p,
				"err", err)
			return false
		}
		if len(key) == 0 {
			dfLog.Errorw("group key is empty", "appKey", f.logKey, "ip", ip.String(), "peer", p)
			return false
		}
		group := PeerGroupInfo{Id: p, Cpl: cpl, IPGroupKey: key}

		if !f.pgm.Allow(group) {
			return false
		}

		peerGroups = append(peerGroups, group)
	}

	if _, ok := f.cplPeerGroups[cpl]; !ok {
		f.cplPeerGroups[cpl] = make(map[peer.ID][]PeerIPGroupKey)
	}

	for _, g := range peerGroups {
		f.pgm.Increment(g)

		f.peerGroups[p] = append(f.peerGroups[p], g)
		f.cplPeerGroups[cpl][p] = append(f.cplPeerGroups[cpl][p], g.IPGroupKey)
	}

	return true
}

// WhitelistPeers will always allow the given peers.
func (f *Filter) WhitelistPeers(peers ...peer.ID) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, p := range peers {
		f.wlpeers[p] = struct{}{}
	}
}

// returns the PeerIPGroupKey to which the given IP belongs.
func (f *Filter) ipGroupKey(ip net.IP) (PeerIPGroupKey, error) {
	switch bz := ip.To4(); bz {
	case nil:
		// TODO Clean up the ASN codebase
		// ipv6 Address -> get ASN
		s, err := f.asnStore.AsnForIPv6(ip)
		if err != nil {
			return "", fmt.Errorf("failed to fetch ASN for IPv6 addr %s: %w", ip.String(), err)
		}

		// if no ASN found then fallback on using the /32 prefix
		if len(s) == 0 {
			dfLog.Debugw("ASN not known", "appKey", f.logKey, "ip", ip)
			s = fmt.Sprintf("unknown ASN: %s", net.CIDRMask(32, 128).String())
		}

		return PeerIPGroupKey(s), nil
	default:
		// If it belongs to a legacy Class 8, we return the /8 prefix as the key
		rs, _ := f.legacyCidrs.ContainingNetworks(ip)
		if len(rs) != 0 {
			key := ip.Mask(net.IPv4Mask(255, 0, 0, 0)).String()
			return PeerIPGroupKey(key), nil
		}

		// otherwise -> /16 prefix
		key := ip.Mask(net.IPv4Mask(255, 255, 0, 0)).String()
		return PeerIPGroupKey(key), nil
	}
}

// CplDiversityStats contains the peer diversity stats for a Cpl.
type CplDiversityStats struct {
	Cpl   int
	Peers map[peer.ID][]PeerIPGroupKey
}

// GetDiversityStats returns the diversity stats for each CPL and is sorted by the CPL.
func (f *Filter) GetDiversityStats() []CplDiversityStats {
	f.mu.Lock()
	defer f.mu.Unlock()

	stats := make([]CplDiversityStats, 0, len(f.cplPeerGroups))

	var sortedCpls []int
	for cpl := range f.cplPeerGroups {
		sortedCpls = append(sortedCpls, cpl)
	}
	sort.Ints(sortedCpls)

	for _, cpl := range sortedCpls {
		ps := make(map[peer.ID][]PeerIPGroupKey, len(f.cplPeerGroups[cpl]))
		cd := CplDiversityStats{cpl, ps}

		for p, groups := range f.cplPeerGroups[cpl] {
			ps[p] = groups
		}
		stats = append(stats, cd)
	}

	return stats
}
