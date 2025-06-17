package test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

// P2PIDs is a slice of p2pkey.PeerID with convenient transform methods.
type P2PIDs []p2pkey.PeerID

// Strings returns the string representation of the p2p IDs.
func (ps P2PIDs) Strings() []string {
	out := make([]string, len(ps))
	for i, p := range ps {
		out[i] = p.String()
	}
	return out
}

// Bytes32 returns the byte representation of the p2p IDs.
func (ps P2PIDs) Bytes32() [][32]byte {
	out := make([][32]byte, len(ps))
	for i, p := range ps {
		out[i] = p
	}
	return out
}

// Unique returns a new slice with duplicate p2p IDs removed.
func (ps P2PIDs) Unique() P2PIDs {
	dedup := make(map[p2pkey.PeerID]struct{})
	var out []p2pkey.PeerID
	for _, p := range ps {
		if _, exists := dedup[p]; !exists {
			out = append(out, p)

			dedup[p] = struct{}{}
		}
	}
	return out
}

func capIDs(t *testing.T, cfgs []kcr.CapabilitiesRegistryCapabilityConfiguration) [][32]byte {
	var out [][32]byte
	for _, cfg := range cfgs {
		out = append(out, cfg.CapabilityId)
	}
	return out
}

func expectedHashedCapabilities(t *testing.T, registry *kcr.CapabilitiesRegistry, don internal.DonCapabilities) [][32]byte {
	out := make([][32]byte, len(don.Capabilities))
	var err error
	for i, capWithCfg := range don.Capabilities {
		out[i], err = registry.GetHashedCapabilityId(nil, capWithCfg.Capability.LabelledName, capWithCfg.Capability.Version)
		require.NoError(t, err)
	}
	return out
}

func sortedHash(p2pids [][32]byte) string {
	sha256Hash := sha256.New()
	sort.Slice(p2pids, func(i, j int) bool {
		return bytes.Compare(p2pids[i][:], p2pids[j][:]) < 0
	})
	for _, id := range p2pids {
		sha256Hash.Write(id[:])
	}
	return hex.EncodeToString(sha256Hash.Sum(nil))
}

func p2p32Bytes(t *testing.T, p2pIDs []p2pkey.PeerID) [][32]byte {
	bs := make([][32]byte, len(p2pIDs))
	for i, p := range p2pIDs {
		bs[i] = p
	}
	return bs
}
