package kbucket

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"

	mh "github.com/multiformats/go-multihash"
)

// maxCplForRefresh is the maximum cpl we support for refresh.
// This limit exists because we can only generate 'maxCplForRefresh' bit prefixes for now.
const maxCplForRefresh uint = 15

// GetTrackedCplsForRefresh returns the Cpl's we are tracking for refresh.
// Caller is free to modify the returned slice as it is a defensive copy.
func (rt *RoutingTable) GetTrackedCplsForRefresh() []time.Time {
	maxCommonPrefix := rt.maxCommonPrefix()
	if maxCommonPrefix > maxCplForRefresh {
		maxCommonPrefix = maxCplForRefresh
	}

	rt.cplRefreshLk.RLock()
	defer rt.cplRefreshLk.RUnlock()

	cpls := make([]time.Time, maxCommonPrefix+1)
	for i := uint(0); i <= maxCommonPrefix; i++ {
		// defaults to the zero value if we haven't refreshed it yet.
		cpls[i] = rt.cplRefreshedAt[i]
	}
	return cpls
}

func randUint16() (uint16, error) {
	// Read a random prefix.
	var prefixBytes [2]byte
	_, err := rand.Read(prefixBytes[:])
	return binary.BigEndian.Uint16(prefixBytes[:]), err
}

// GenRandPeerID generates a random peerID for a given Cpl
func (rt *RoutingTable) GenRandPeerID(targetCpl uint) (peer.ID, error) {
	if targetCpl > maxCplForRefresh {
		return "", fmt.Errorf("cannot generate peer ID for Cpl greater than %d", maxCplForRefresh)
	}

	localPrefix := binary.BigEndian.Uint16(rt.local)

	// For host with ID `L`, an ID `K` belongs to a bucket with ID `B` ONLY IF CommonPrefixLen(L,K) is EXACTLY B.
	// Hence, to achieve a targetPrefix `T`, we must toggle the (T+1)th bit in L & then copy (T+1) bits from L
	// to our randomly generated prefix.
	toggledLocalPrefix := localPrefix ^ (uint16(0x8000) >> targetCpl)
	randPrefix, err := randUint16()
	if err != nil {
		return "", err
	}

	// Combine the toggled local prefix and the random bits at the correct offset
	// such that ONLY the first `targetCpl` bits match the local ID.
	mask := (^uint16(0)) << (16 - (targetCpl + 1))
	targetPrefix := (toggledLocalPrefix & mask) | (randPrefix & ^mask)

	// Convert to a known peer ID.
	key := keyPrefixMap[targetPrefix]
	id := [34]byte{mh.SHA2_256, 32}
	binary.BigEndian.PutUint32(id[2:], key)
	return peer.ID(id[:]), nil
}

// ResetCplRefreshedAtForID resets the refresh time for the Cpl of the given ID.
func (rt *RoutingTable) ResetCplRefreshedAtForID(id ID, newTime time.Time) {
	cpl := CommonPrefixLen(id, rt.local)
	if uint(cpl) > maxCplForRefresh {
		return
	}

	rt.cplRefreshLk.Lock()
	defer rt.cplRefreshLk.Unlock()

	rt.cplRefreshedAt[uint(cpl)] = newTime
}
