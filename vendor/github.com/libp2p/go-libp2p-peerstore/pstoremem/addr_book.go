package pstoremem

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/peer"
	pstore "github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/record"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/libp2p/go-libp2p-peerstore/addr"
)

var log = logging.Logger("peerstore")

type expiringAddr struct {
	Addr    ma.Multiaddr
	TTL     time.Duration
	Expires time.Time
}

func (e *expiringAddr) ExpiredBy(t time.Time) bool {
	return !t.Before(e.Expires)
}

type peerRecordState struct {
	Envelope *record.Envelope
	Seq      uint64
}

type addrSegments [256]*addrSegment

type addrSegment struct {
	sync.RWMutex

	// Use pointers to save memory. Maps always leave some fraction of their
	// space unused. storing the *values* directly in the map will
	// drastically increase the space waste. In our case, by 6x.
	addrs map[peer.ID]map[string]*expiringAddr

	signedPeerRecords map[peer.ID]*peerRecordState
}

func (segments *addrSegments) get(p peer.ID) *addrSegment {
	return segments[byte(p[len(p)-1])]
}

// memoryAddrBook manages addresses.
type memoryAddrBook struct {
	segments addrSegments

	ctx    context.Context
	cancel func()

	subManager *AddrSubManager
}

var _ pstore.AddrBook = (*memoryAddrBook)(nil)
var _ pstore.CertifiedAddrBook = (*memoryAddrBook)(nil)

func NewAddrBook() *memoryAddrBook {
	ctx, cancel := context.WithCancel(context.Background())

	ab := &memoryAddrBook{
		segments: func() (ret addrSegments) {
			for i, _ := range ret {
				ret[i] = &addrSegment{
					addrs:             make(map[peer.ID]map[string]*expiringAddr),
					signedPeerRecords: make(map[peer.ID]*peerRecordState)}
			}
			return ret
		}(),
		subManager: NewAddrSubManager(),
		ctx:        ctx,
		cancel:     cancel,
	}

	go ab.background()
	return ab
}

// background periodically schedules a gc
func (mab *memoryAddrBook) background() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mab.gc()

		case <-mab.ctx.Done():
			return
		}
	}
}

func (mab *memoryAddrBook) Close() error {
	mab.cancel()
	return nil
}

// gc garbage collects the in-memory address book.
func (mab *memoryAddrBook) gc() {
	now := time.Now()
	for _, s := range mab.segments {
		s.Lock()
		var collectedPeers []peer.ID
		for p, amap := range s.addrs {
			for k, addr := range amap {
				if addr.ExpiredBy(now) {
					delete(amap, k)
				}
			}
			if len(amap) == 0 {
				delete(s.addrs, p)
				collectedPeers = append(collectedPeers, p)
			}
		}
		// remove signed records for peers whose signed addrs have all been removed
		for _, p := range collectedPeers {
			delete(s.signedPeerRecords, p)
		}
		s.Unlock()
	}
}

func (mab *memoryAddrBook) PeersWithAddrs() peer.IDSlice {
	// deduplicate, since the same peer could have both signed & unsigned addrs
	pidSet := peer.NewSet()
	for _, s := range mab.segments {
		s.RLock()
		for pid, amap := range s.addrs {
			if amap != nil && len(amap) > 0 {
				pidSet.Add(pid)
			}
		}
		s.RUnlock()
	}
	return pidSet.Peers()
}

// AddAddr calls AddAddrs(p, []ma.Multiaddr{addr}, ttl)
func (mab *memoryAddrBook) AddAddr(p peer.ID, addr ma.Multiaddr, ttl time.Duration) {
	mab.AddAddrs(p, []ma.Multiaddr{addr}, ttl)
}

// AddAddrs gives memoryAddrBook addresses to use, with a given ttl
// (time-to-live), after which the address is no longer valid.
// This function never reduces the TTL or expiration of an address.
func (mab *memoryAddrBook) AddAddrs(p peer.ID, addrs []ma.Multiaddr, ttl time.Duration) {
	// if we have a valid peer record, ignore unsigned addrs
	// peerRec := mab.GetPeerRecord(p)
	// if peerRec != nil {
	// 	return
	// }
	mab.addAddrs(p, addrs, ttl)
}

// ConsumePeerRecord adds addresses from a signed peer.PeerRecord (contained in
// a record.Envelope), which will expire after the given TTL.
// See https://godoc.org/github.com/libp2p/go-libp2p-core/peerstore#CertifiedAddrBook for more details.
func (mab *memoryAddrBook) ConsumePeerRecord(recordEnvelope *record.Envelope, ttl time.Duration) (bool, error) {
	r, err := recordEnvelope.Record()
	if err != nil {
		return false, err
	}
	rec, ok := r.(*peer.PeerRecord)
	if !ok {
		return false, fmt.Errorf("unable to process envelope: not a PeerRecord")
	}
	if !rec.PeerID.MatchesPublicKey(recordEnvelope.PublicKey) {
		return false, fmt.Errorf("signing key does not match PeerID in PeerRecord")
	}

	// ensure seq is greater than, or equal to, the last received
	s := mab.segments.get(rec.PeerID)
	s.Lock()
	defer s.Unlock()
	lastState, found := s.signedPeerRecords[rec.PeerID]
	if found && lastState.Seq > rec.Seq {
		return false, nil
	}
	s.signedPeerRecords[rec.PeerID] = &peerRecordState{
		Envelope: recordEnvelope,
		Seq:      rec.Seq,
	}
	mab.addAddrsUnlocked(s, rec.PeerID, rec.Addrs, ttl, true)
	return true, nil
}

func (mab *memoryAddrBook) addAddrs(p peer.ID, addrs []ma.Multiaddr, ttl time.Duration) {
	if err := p.Validate(); err != nil {
		log.Warningf("tried to set addrs for invalid peer ID %s: %s", p, err)
		return
	}

	s := mab.segments.get(p)
	s.Lock()
	defer s.Unlock()

	mab.addAddrsUnlocked(s, p, addrs, ttl, false)
}

func (mab *memoryAddrBook) addAddrsUnlocked(s *addrSegment, p peer.ID, addrs []ma.Multiaddr, ttl time.Duration, signed bool) {
	// if ttl is zero, exit. nothing to do.
	if ttl <= 0 {
		return
	}

	amap, ok := s.addrs[p]
	if !ok {
		amap = make(map[string]*expiringAddr)
		s.addrs[p] = amap
	}

	exp := time.Now().Add(ttl)
	addrSet := make(map[string]struct{}, len(addrs))
	for _, addr := range addrs {
		if addr == nil {
			log.Warnf("was passed nil multiaddr for %s", p)
			continue
		}
		k := string(addr.Bytes())
		addrSet[k] = struct{}{}

		// find the highest TTL and Expiry time between
		// existing records and function args
		a, found := amap[k] // won't allocate.

		if !found {
			// not found, announce it.
			entry := &expiringAddr{Addr: addr, Expires: exp, TTL: ttl}
			amap[k] = entry
			mab.subManager.BroadcastAddr(p, addr)
		} else {
			// update ttl & exp to whichever is greater between new and existing entry
			if ttl > a.TTL {
				a.TTL = ttl
			}
			if exp.After(a.Expires) {
				a.Expires = exp
			}
		}
	}
}

// SetAddr calls mgr.SetAddrs(p, addr, ttl)
func (mab *memoryAddrBook) SetAddr(p peer.ID, addr ma.Multiaddr, ttl time.Duration) {
	if err := p.Validate(); err != nil {
		log.Warningf("tried to set addrs for invalid peer ID %s: %s", p, err)
		return
	}

	mab.SetAddrs(p, []ma.Multiaddr{addr}, ttl)
}

// SetAddrs sets the ttl on addresses. This clears any TTL there previously.
// This is used when we receive the best estimate of the validity of an address.
func (mab *memoryAddrBook) SetAddrs(p peer.ID, addrs []ma.Multiaddr, ttl time.Duration) {
	if err := p.Validate(); err != nil {
		log.Warningf("tried to set addrs for invalid peer ID %s: %s", p, err)
		return
	}

	s := mab.segments.get(p)
	s.Lock()
	defer s.Unlock()

	amap, ok := s.addrs[p]
	if !ok {
		amap = make(map[string]*expiringAddr)
		s.addrs[p] = amap
	}

	exp := time.Now().Add(ttl)
	for _, addr := range addrs {
		if addr == nil {
			log.Warnf("was passed nil multiaddr for %s", p)
			continue
		}
		aBytes := addr.Bytes()
		key := string(aBytes)

		// re-set all of them for new ttl.
		if ttl > 0 {
			amap[key] = &expiringAddr{Addr: addr, Expires: exp, TTL: ttl}
			mab.subManager.BroadcastAddr(p, addr)
		} else {
			delete(amap, key)
		}
	}

	// if we've expired all the signed addresses for a peer, remove their signed routing state record
	if len(amap) == 0 {
		delete(s.signedPeerRecords, p)
	}
}

// UpdateAddrs updates the addresses associated with the given peer that have
// the given oldTTL to have the given newTTL.
func (mab *memoryAddrBook) UpdateAddrs(p peer.ID, oldTTL time.Duration, newTTL time.Duration) {
	if err := p.Validate(); err != nil {
		log.Warningf("tried to set addrs for invalid peer ID %s: %s", p, err)
		return
	}

	s := mab.segments.get(p)
	s.Lock()
	defer s.Unlock()
	exp := time.Now().Add(newTTL)
	amap, found := s.addrs[p]
	if !found {
		return
	}

	for k, a := range amap {
		if oldTTL == a.TTL {
			if newTTL == 0 {
				delete(amap, k)
			} else {
				a.TTL = newTTL
				a.Expires = exp
				amap[k] = a
			}
		}
	}

	// if we've expired all the signed addresses for a peer, remove their signed routing state record
	if len(amap) == 0 {
		delete(s.signedPeerRecords, p)
	}
}

// Addrs returns all known (and valid) addresses for a given peer
func (mab *memoryAddrBook) Addrs(p peer.ID) []ma.Multiaddr {
	if err := p.Validate(); err != nil {
		// invalid peer ID = no addrs
		return nil
	}

	s := mab.segments.get(p)
	s.RLock()
	defer s.RUnlock()

	return validAddrs(s.addrs[p])
}

func validAddrs(amap map[string]*expiringAddr) []ma.Multiaddr {
	now := time.Now()
	good := make([]ma.Multiaddr, 0, len(amap))
	if amap == nil {
		return good
	}
	for _, m := range amap {
		if !m.ExpiredBy(now) {
			good = append(good, m.Addr)
		}
	}

	return good
}

// GetPeerRecord returns a Envelope containing a PeerRecord for the
// given peer id, if one exists.
// Returns nil if no signed PeerRecord exists for the peer.
func (mab *memoryAddrBook) GetPeerRecord(p peer.ID) *record.Envelope {
	if err := p.Validate(); err != nil {
		// invalid peer ID = no addrs
		return nil
	}

	s := mab.segments.get(p)
	s.RLock()
	defer s.RUnlock()

	// although the signed record gets garbage collected when all addrs inside it are expired,
	// we may be in between the expiration time and the GC interval
	// so, we check to see if we have any valid signed addrs before returning the record
	if len(validAddrs(s.addrs[p])) == 0 {
		return nil
	}

	state := s.signedPeerRecords[p]
	if state == nil {
		return nil
	}
	return state.Envelope
}

// ClearAddrs removes all previously stored addresses
func (mab *memoryAddrBook) ClearAddrs(p peer.ID) {
	if err := p.Validate(); err != nil {
		// nothing to clear
		return
	}

	s := mab.segments.get(p)
	s.Lock()
	defer s.Unlock()

	delete(s.addrs, p)
	delete(s.signedPeerRecords, p)
}

// AddrStream returns a channel on which all new addresses discovered for a
// given peer ID will be published.
func (mab *memoryAddrBook) AddrStream(ctx context.Context, p peer.ID) <-chan ma.Multiaddr {
	if err := p.Validate(); err != nil {
		log.Warningf("tried to get addrs for invalid peer ID %s: %s", p, err)
		ch := make(chan ma.Multiaddr)
		close(ch)
		return ch
	}

	s := mab.segments.get(p)
	s.RLock()
	defer s.RUnlock()

	baseaddrslice := s.addrs[p]
	initial := make([]ma.Multiaddr, 0, len(baseaddrslice))
	for _, a := range baseaddrslice {
		initial = append(initial, a.Addr)
	}

	return mab.subManager.AddrStream(ctx, p, initial)
}

type addrSub struct {
	pubch  chan ma.Multiaddr
	lk     sync.Mutex
	buffer []ma.Multiaddr
	ctx    context.Context
}

func (s *addrSub) pubAddr(a ma.Multiaddr) {
	select {
	case s.pubch <- a:
	case <-s.ctx.Done():
	}
}

// An abstracted, pub-sub manager for address streams. Extracted from
// memoryAddrBook in order to support additional implementations.
type AddrSubManager struct {
	mu   sync.RWMutex
	subs map[peer.ID][]*addrSub
}

// NewAddrSubManager initializes an AddrSubManager.
func NewAddrSubManager() *AddrSubManager {
	return &AddrSubManager{
		subs: make(map[peer.ID][]*addrSub),
	}
}

// Used internally by the address stream coroutine to remove a subscription
// from the manager.
func (mgr *AddrSubManager) removeSub(p peer.ID, s *addrSub) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	subs := mgr.subs[p]
	if len(subs) == 1 {
		if subs[0] != s {
			return
		}
		delete(mgr.subs, p)
		return
	}

	for i, v := range subs {
		if v == s {
			subs[i] = subs[len(subs)-1]
			subs[len(subs)-1] = nil
			mgr.subs[p] = subs[:len(subs)-1]
			return
		}
	}
}

// BroadcastAddr broadcasts a new address to all subscribed streams.
func (mgr *AddrSubManager) BroadcastAddr(p peer.ID, addr ma.Multiaddr) {
	mgr.mu.RLock()
	defer mgr.mu.RUnlock()

	if subs, ok := mgr.subs[p]; ok {
		for _, sub := range subs {
			sub.pubAddr(addr)
		}
	}
}

// AddrStream creates a new subscription for a given peer ID, pre-populating the
// channel with any addresses we might already have on file.
func (mgr *AddrSubManager) AddrStream(ctx context.Context, p peer.ID, initial []ma.Multiaddr) <-chan ma.Multiaddr {
	sub := &addrSub{pubch: make(chan ma.Multiaddr), ctx: ctx}
	out := make(chan ma.Multiaddr)

	mgr.mu.Lock()
	if _, ok := mgr.subs[p]; ok {
		mgr.subs[p] = append(mgr.subs[p], sub)
	} else {
		mgr.subs[p] = []*addrSub{sub}
	}
	mgr.mu.Unlock()

	sort.Sort(addr.AddrList(initial))

	go func(buffer []ma.Multiaddr) {
		defer close(out)

		sent := make(map[string]bool, len(buffer))
		var outch chan ma.Multiaddr

		for _, a := range buffer {
			sent[string(a.Bytes())] = true
		}

		var next ma.Multiaddr
		if len(buffer) > 0 {
			next = buffer[0]
			buffer = buffer[1:]
			outch = out
		}

		for {
			select {
			case outch <- next:
				if len(buffer) > 0 {
					next = buffer[0]
					buffer = buffer[1:]
				} else {
					outch = nil
					next = nil
				}
			case naddr := <-sub.pubch:
				if sent[string(naddr.Bytes())] {
					continue
				}

				sent[string(naddr.Bytes())] = true
				if next == nil {
					next = naddr
					outch = out
				} else {
					buffer = append(buffer, naddr)
				}
			case <-ctx.Done():
				mgr.removeSub(p, sub)
				return
			}
		}

	}(initial)

	return out
}
