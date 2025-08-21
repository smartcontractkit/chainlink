package pkg

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"sort"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func PeerIDsToBytes(p2pIDs []p2pkey.PeerID) [][32]byte {
	out := make([][32]byte, len(p2pIDs))
	for i, p2pID := range p2pIDs {
		out[i] = p2pID
	}
	return out
}

func SortedHash(p2pids [][32]byte) string {
	sha256Hash := sha256.New()
	sort.Slice(p2pids, func(i, j int) bool {
		return bytes.Compare(p2pids[i][:], p2pids[j][:]) < 0
	})
	for _, id := range p2pids {
		sha256Hash.Write(id[:])
	}
	return hex.EncodeToString(sha256Hash.Sum(nil))
}

func BytesToPeerIDs(p2pIDs [][32]byte) []p2pkey.PeerID {
	out := make([]p2pkey.PeerID, len(p2pIDs))
	for i, p2pID := range p2pIDs {
		out[i] = p2pID
	}
	return out
}
