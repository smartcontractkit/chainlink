// package kbucket implements a kademlia 'k-bucket' routing table.
package kbucket

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"

	"github.com/libp2p/go-libp2p-kbucket/peerdiversity"

	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("table")

var ErrPeerRejectedHighLatency = errors.New("peer rejected; latency too high")
var ErrPeerRejectedNoCapacity = errors.New("peer rejected; insufficient capacity")

// RoutingTable defines the routing table.
type RoutingTable struct {
	// the routing table context
	ctx context.Context
	// function to cancel the RT context
	ctxCancel context.CancelFunc

	// ID of the local peer
	local ID

	// Blanket lock, refine later for better performance
	tabLock sync.RWMutex

	// latency metrics
	metrics peerstore.Metrics

	// Maximum acceptable latency for peers in this cluster
	maxLatency time.Duration

	// kBuckets define all the fingers to other nodes.
	buckets    []*bucket
	bucketsize int

	cplRefreshLk   sync.RWMutex
	cplRefreshedAt map[uint]time.Time

	// notification functions
	PeerRemoved func(peer.ID)
	PeerAdded   func(peer.ID)

	// usefulnessGracePeriod is the maximum grace period we will give to a
	// peer in the bucket to be useful to us, failing which, we will evict
	// it to make place for a new peer if the bucket is full
	usefulnessGracePeriod time.Duration

	df *peerdiversity.Filter
}

// NewRoutingTable creates a new routing table with a given bucketsize, local ID, and latency tolerance.
func NewRoutingTable(bucketsize int, localID ID, latency time.Duration, m peerstore.Metrics, usefulnessGracePeriod time.Duration,
	df *peerdiversity.Filter) (*RoutingTable, error) {
	rt := &RoutingTable{
		buckets:    []*bucket{newBucket()},
		bucketsize: bucketsize,
		local:      localID,

		maxLatency: latency,
		metrics:    m,

		cplRefreshedAt: make(map[uint]time.Time),

		PeerRemoved: func(peer.ID) {},
		PeerAdded:   func(peer.ID) {},

		usefulnessGracePeriod: usefulnessGracePeriod,

		df: df,
	}

	rt.ctx, rt.ctxCancel = context.WithCancel(context.Background())

	return rt, nil
}

// Close shuts down the Routing Table & all associated processes.
// It is safe to call this multiple times.
func (rt *RoutingTable) Close() error {
	rt.ctxCancel()
	return nil
}

// NPeersForCPL returns the number of peers we have for a given Cpl
func (rt *RoutingTable) NPeersForCpl(cpl uint) int {
	rt.tabLock.RLock()
	defer rt.tabLock.RUnlock()

	// it's in the last bucket
	if int(cpl) >= len(rt.buckets)-1 {
		count := 0
		b := rt.buckets[len(rt.buckets)-1]
		for _, p := range b.peers() {
			if CommonPrefixLen(rt.local, p.dhtId) == int(cpl) {
				count++
			}
		}
		return count
	} else {
		return rt.buckets[cpl].len()
	}
}

// TryAddPeer tries to add a peer to the Routing table.
// If the peer ALREADY exists in the Routing Table and has been queried before, this call is a no-op.
// If the peer ALREADY exists in the Routing Table but hasn't been queried before, we set it's LastUsefulAt value to
// the current time. This needs to done because we don't mark peers as "Useful"(by setting the LastUsefulAt value)
// when we first connect to them.
//
// If the peer is a queryPeer i.e. we queried it or it queried us, we set the LastSuccessfulOutboundQuery to the current time.
// If the peer is just a peer that we connect to/it connected to us without any DHT query, we consider it as having
// no LastSuccessfulOutboundQuery.
//
//
// If the logical bucket to which the peer belongs is full and it's not the last bucket, we try to replace an existing peer
// whose LastSuccessfulOutboundQuery is above the maximum allowed threshold in that bucket with the new peer.
// If no such peer exists in that bucket, we do NOT add the peer to the Routing Table and return error "ErrPeerRejectedNoCapacity".

// It returns a boolean value set to true if the peer was newly added to the Routing Table, false otherwise.
// It also returns any error that occurred while adding the peer to the Routing Table. If the error is not nil,
// the boolean value will ALWAYS be false i.e. the peer wont be added to the Routing Table it it's not already there.
//
// A return value of false with error=nil indicates that the peer ALREADY exists in the Routing Table.
func (rt *RoutingTable) TryAddPeer(p peer.ID, queryPeer bool, isReplaceable bool) (bool, error) {
	rt.tabLock.Lock()
	defer rt.tabLock.Unlock()

	return rt.addPeer(p, queryPeer, isReplaceable)
}

// locking is the responsibility of the caller
func (rt *RoutingTable) addPeer(p peer.ID, queryPeer bool, isReplaceable bool) (bool, error) {
	bucketID := rt.bucketIdForPeer(p)
	bucket := rt.buckets[bucketID]

	now := time.Now()
	var lastUsefulAt time.Time
	if queryPeer {
		lastUsefulAt = now
	}

	// peer already exists in the Routing Table.
	if peer := bucket.getPeer(p); peer != nil {
		// if we're querying the peer first time after adding it, let's give it a
		// usefulness bump. This will ONLY happen once.
		if peer.LastUsefulAt.IsZero() && queryPeer {
			peer.LastUsefulAt = lastUsefulAt
		}
		return false, nil
	}

	// peer's latency threshold is NOT acceptable
	if rt.metrics.LatencyEWMA(p) > rt.maxLatency {
		// Connection doesnt meet requirements, skip!
		return false, ErrPeerRejectedHighLatency
	}

	// add it to the diversity filter for now.
	// if we aren't able to find a place for the peer in the table,
	// we will simply remove it from the Filter later.
	if rt.df != nil {
		if !rt.df.TryAdd(p) {
			return false, errors.New("peer rejected by the diversity filter")
		}
	}

	// We have enough space in the bucket (whether spawned or grouped).
	if bucket.len() < rt.bucketsize {
		bucket.pushFront(&PeerInfo{
			Id:                            p,
			LastUsefulAt:                  lastUsefulAt,
			LastSuccessfulOutboundQueryAt: now,
			AddedAt:                       now,
			dhtId:                         ConvertPeerID(p),
			replaceable:                   isReplaceable,
		})
		rt.PeerAdded(p)
		return true, nil
	}

	if bucketID == len(rt.buckets)-1 {
		// if the bucket is too large and this is the last bucket (i.e. wildcard), unfold it.
		rt.nextBucket()
		// the structure of the table has changed, so let's recheck if the peer now has a dedicated bucket.
		bucketID = rt.bucketIdForPeer(p)
		bucket = rt.buckets[bucketID]

		// push the peer only if the bucket isn't overflowing after slitting
		if bucket.len() < rt.bucketsize {
			bucket.pushFront(&PeerInfo{
				Id:                            p,
				LastUsefulAt:                  lastUsefulAt,
				LastSuccessfulOutboundQueryAt: now,
				AddedAt:                       now,
				dhtId:                         ConvertPeerID(p),
				replaceable:                   isReplaceable,
			})
			rt.PeerAdded(p)
			return true, nil
		}
	}

	// the bucket to which the peer belongs is full. Let's try to find a peer
	// in that bucket which is replaceable.
	// we don't really need a stable sort here as it dosen't matter which peer we evict
	// as long as it's a replaceable peer.
	replaceablePeer := bucket.min(func(p1 *PeerInfo, p2 *PeerInfo) bool {
		return p1.replaceable
	})

	if replaceablePeer != nil && replaceablePeer.replaceable {
		// let's evict it and add the new peer
		if rt.removePeer(replaceablePeer.Id) {
			bucket.pushFront(&PeerInfo{
				Id:                            p,
				LastUsefulAt:                  lastUsefulAt,
				LastSuccessfulOutboundQueryAt: now,
				AddedAt:                       now,
				dhtId:                         ConvertPeerID(p),
				replaceable:                   isReplaceable,
			})
			rt.PeerAdded(p)
			return true, nil
		}
	}

	// we weren't able to find place for the peer, remove it from the filter state.
	if rt.df != nil {
		rt.df.Remove(p)
	}
	return false, ErrPeerRejectedNoCapacity
}

// MarkAllPeersIrreplaceable marks all peers in the routing table as irreplaceable
// This means that we will never replace an existing peer in the table to make space for a new peer.
// However, they can still be removed by calling the `RemovePeer` API.
func (rt *RoutingTable) MarkAllPeersIrreplaceable() {
	rt.tabLock.Lock()
	defer rt.tabLock.Unlock()

	for i := range rt.buckets {
		b := rt.buckets[i]
		b.updateAllWith(func(p *PeerInfo) {
			p.replaceable = false
		})
	}
}

// GetPeerInfos returns the peer information that we've stored in the buckets
func (rt *RoutingTable) GetPeerInfos() []PeerInfo {
	rt.tabLock.RLock()
	defer rt.tabLock.RUnlock()

	var pis []PeerInfo
	for _, b := range rt.buckets {
		for _, p := range b.peers() {
			pis = append(pis, p)
		}
	}
	return pis
}

// UpdateLastSuccessfulOutboundQuery updates the LastSuccessfulOutboundQueryAt time of the peer.
// Returns true if the update was successful, false otherwise.
func (rt *RoutingTable) UpdateLastSuccessfulOutboundQueryAt(p peer.ID, t time.Time) bool {
	rt.tabLock.Lock()
	defer rt.tabLock.Unlock()

	bucketID := rt.bucketIdForPeer(p)
	bucket := rt.buckets[bucketID]

	if pc := bucket.getPeer(p); pc != nil {
		pc.LastSuccessfulOutboundQueryAt = t
		return true
	}
	return false
}

// UpdateLastUsefulAt updates the LastUsefulAt time of the peer.
// Returns true if the update was successful, false otherwise.
func (rt *RoutingTable) UpdateLastUsefulAt(p peer.ID, t time.Time) bool {
	rt.tabLock.Lock()
	defer rt.tabLock.Unlock()

	bucketID := rt.bucketIdForPeer(p)
	bucket := rt.buckets[bucketID]

	if pc := bucket.getPeer(p); pc != nil {
		pc.LastUsefulAt = t
		return true
	}
	return false
}

// RemovePeer should be called when the caller is sure that a peer is not useful for queries.
// For eg: the peer could have stopped supporting the DHT protocol.
// It evicts the peer from the Routing Table.
func (rt *RoutingTable) RemovePeer(p peer.ID) {
	rt.tabLock.Lock()
	defer rt.tabLock.Unlock()
	rt.removePeer(p)
}

// locking is the responsibility of the caller
func (rt *RoutingTable) removePeer(p peer.ID) bool {
	bucketID := rt.bucketIdForPeer(p)
	bucket := rt.buckets[bucketID]
	if bucket.remove(p) {
		if rt.df != nil {
			rt.df.Remove(p)
		}
		for {
			lastBucketIndex := len(rt.buckets) - 1

			// remove the last bucket if it's empty and it isn't the only bucket we have
			if len(rt.buckets) > 1 && rt.buckets[lastBucketIndex].len() == 0 {
				rt.buckets[lastBucketIndex] = nil
				rt.buckets = rt.buckets[:lastBucketIndex]
			} else if len(rt.buckets) >= 2 && rt.buckets[lastBucketIndex-1].len() == 0 {
				// if the second last bucket just became empty, remove and replace it with the last bucket.
				rt.buckets[lastBucketIndex-1] = rt.buckets[lastBucketIndex]
				rt.buckets[lastBucketIndex] = nil
				rt.buckets = rt.buckets[:lastBucketIndex]
			} else {
				break
			}
		}

		// peer removed callback
		rt.PeerRemoved(p)
		return true
	}
	return false
}

func (rt *RoutingTable) nextBucket() {
	// This is the last bucket, which allegedly is a mixed bag containing peers not belonging in dedicated (unfolded) buckets.
	// _allegedly_ is used here to denote that *all* peers in the last bucket might feasibly belong to another bucket.
	// This could happen if e.g. we've unfolded 4 buckets, and all peers in folded bucket 5 really belong in bucket 8.
	bucket := rt.buckets[len(rt.buckets)-1]
	newBucket := bucket.split(len(rt.buckets)-1, rt.local)
	rt.buckets = append(rt.buckets, newBucket)

	// The newly formed bucket still contains too many peers. We probably just unfolded a empty bucket.
	if newBucket.len() >= rt.bucketsize {
		// Keep unfolding the table until the last bucket is not overflowing.
		rt.nextBucket()
	}
}

// Find a specific peer by ID or return nil
func (rt *RoutingTable) Find(id peer.ID) peer.ID {
	srch := rt.NearestPeers(ConvertPeerID(id), 1)
	if len(srch) == 0 || srch[0] != id {
		return ""
	}
	return srch[0]
}

// NearestPeer returns a single peer that is nearest to the given ID
func (rt *RoutingTable) NearestPeer(id ID) peer.ID {
	peers := rt.NearestPeers(id, 1)
	if len(peers) > 0 {
		return peers[0]
	}

	log.Debugf("NearestPeer: Returning nil, table size = %d", rt.Size())
	return ""
}

// NearestPeers returns a list of the 'count' closest peers to the given ID
func (rt *RoutingTable) NearestPeers(id ID, count int) []peer.ID {
	// This is the number of bits _we_ share with the key. All peers in this
	// bucket share cpl bits with us and will therefore share at least cpl+1
	// bits with the given key. +1 because both the target and all peers in
	// this bucket differ from us in the cpl bit.
	cpl := CommonPrefixLen(id, rt.local)

	// It's assumed that this also protects the buckets.
	rt.tabLock.RLock()

	// Get bucket index or last bucket
	if cpl >= len(rt.buckets) {
		cpl = len(rt.buckets) - 1
	}

	pds := peerDistanceSorter{
		peers:  make([]peerDistance, 0, count+rt.bucketsize),
		target: id,
	}

	// Add peers from the target bucket (cpl+1 shared bits).
	pds.appendPeersFromList(rt.buckets[cpl].list)

	// If we're short, add peers from all buckets to the right. All buckets
	// to the right share exactly cpl bits (as opposed to the cpl+1 bits
	// shared by the peers in the cpl bucket).
	//
	// This is, unfortunately, less efficient than we'd like. We will switch
	// to a trie implementation eventually which will allow us to find the
	// closest N peers to any target key.

	if pds.Len() < count {
		for i := cpl + 1; i < len(rt.buckets); i++ {
			pds.appendPeersFromList(rt.buckets[i].list)
		}
	}

	// If we're still short, add in buckets that share _fewer_ bits. We can
	// do this bucket by bucket because each bucket will share 1 fewer bit
	// than the last.
	//
	// * bucket cpl-1: cpl-1 shared bits.
	// * bucket cpl-2: cpl-2 shared bits.
	// ...
	for i := cpl - 1; i >= 0 && pds.Len() < count; i-- {
		pds.appendPeersFromList(rt.buckets[i].list)
	}
	rt.tabLock.RUnlock()

	// Sort by distance to local peer
	pds.sort()

	if count < pds.Len() {
		pds.peers = pds.peers[:count]
	}

	out := make([]peer.ID, 0, pds.Len())
	for _, p := range pds.peers {
		out = append(out, p.p)
	}

	return out
}

// Size returns the total number of peers in the routing table
func (rt *RoutingTable) Size() int {
	var tot int
	rt.tabLock.RLock()
	for _, buck := range rt.buckets {
		tot += buck.len()
	}
	rt.tabLock.RUnlock()
	return tot
}

// ListPeers takes a RoutingTable and returns a list of all peers from all buckets in the table.
func (rt *RoutingTable) ListPeers() []peer.ID {
	rt.tabLock.RLock()
	defer rt.tabLock.RUnlock()

	var peers []peer.ID
	for _, buck := range rt.buckets {
		peers = append(peers, buck.peerIds()...)
	}
	return peers
}

// Print prints a descriptive statement about the provided RoutingTable
func (rt *RoutingTable) Print() {
	fmt.Printf("Routing Table, bs = %d, Max latency = %d\n", rt.bucketsize, rt.maxLatency)
	rt.tabLock.RLock()

	for i, b := range rt.buckets {
		fmt.Printf("\tbucket: %d\n", i)

		for e := b.list.Front(); e != nil; e = e.Next() {
			p := e.Value.(*PeerInfo).Id
			fmt.Printf("\t\t- %s %s\n", p.Pretty(), rt.metrics.LatencyEWMA(p).String())
		}
	}
	rt.tabLock.RUnlock()
}

// GetDiversityStats returns the diversity stats for the Routing Table if a diversity Filter
// is configured.
func (rt *RoutingTable) GetDiversityStats() []peerdiversity.CplDiversityStats {
	if rt.df != nil {
		return rt.df.GetDiversityStats()
	}
	return nil
}

// the caller is responsible for the locking
func (rt *RoutingTable) bucketIdForPeer(p peer.ID) int {
	peerID := ConvertPeerID(p)
	cpl := CommonPrefixLen(peerID, rt.local)
	bucketID := cpl
	if bucketID >= len(rt.buckets) {
		bucketID = len(rt.buckets) - 1
	}
	return bucketID
}

// maxCommonPrefix returns the maximum common prefix length between any peer in
// the table and the current peer.
func (rt *RoutingTable) maxCommonPrefix() uint {
	rt.tabLock.RLock()
	defer rt.tabLock.RUnlock()

	for i := len(rt.buckets) - 1; i >= 0; i-- {
		if rt.buckets[i].len() > 0 {
			return rt.buckets[i].maxCommonPrefix(rt.local)
		}
	}
	return 0
}
