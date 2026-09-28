package handlers

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

// ShardedDONs pairs a DON×shard config matrix with the connection managers that
// serve it. The outer dimension of ConnMgrs indexes DONs (parallel to DONs); the
// inner dimension indexes shards within that DON (parallel to DONs[i].Shards).
// Use NewShardedDONs to build one: it checks the two stay in lockstep once, so
// callers can pass a single value around instead of two slices that must always
// agree in length.
type ShardedDONs struct {
	DONs     []config.ShardedDONConfig
	ConnMgrs [][]DON
}

// NewShardedDONs validates that shardedDONs and shardsConnMgrs describe the same
// DON×shard matrix (equal outer length, equal per-DON shard counts, no nil
// connection managers) and bundles them together.
func NewShardedDONs(shardedDONs []config.ShardedDONConfig, shardsConnMgrs [][]DON) (*ShardedDONs, error) {
	if len(shardedDONs) == 0 || len(shardsConnMgrs) == 0 {
		return nil, errors.New("at least one DON and connection manager required")
	}
	if len(shardedDONs) != len(shardsConnMgrs) {
		return nil, fmt.Errorf("shardedDONs length %d does not match shardsConnMgrs length %d", len(shardedDONs), len(shardsConnMgrs))
	}
	for i, sd := range shardedDONs {
		connMgrs := shardsConnMgrs[i]
		if len(sd.Shards) != len(connMgrs) {
			return nil, fmt.Errorf("DON %s: shards length %d does not match connection managers length %d", sd.DonName, len(sd.Shards), len(connMgrs))
		}
		for shardIdx, connMgr := range connMgrs {
			if connMgr == nil {
				return nil, fmt.Errorf("DON %s shard %d: nil connection manager", sd.DonName, shardIdx)
			}
		}
	}
	return &ShardedDONs{DONs: shardedDONs, ConnMgrs: shardsConnMgrs}, nil
}

// ShardEndpoint is one shard of one DON: its identifier, the connection manager
// that reaches its members, and the config (membership, Byzantine fault
// tolerance threshold F) that applies to it.
type ShardEndpoint struct {
	// DonID is the shard-specific DON identifier, derived via config.ShardDONID.
	// Shard 0 uses the bare DON name; shard N>0 uses "donName_N".
	DonID    string
	DonName  string
	ShardIdx int
	// ConnMgr is the connection manager for this shard. It only knows how to reach
	// nodes that are members of this shard.
	ConnMgr DON
	Members []config.NodeConfig
	F       int
}

// BuildShardEndpoints expands the DON×shard matrix into a flat list of
// ShardEndpoints, and builds a nodeAddr -> ShardEndpoint lookup. Node addresses
// must be unique (disjoint) across the whole matrix.
func (s *ShardedDONs) BuildShardEndpoints() ([]*ShardEndpoint, map[string]*ShardEndpoint, error) {
	var endpoints []*ShardEndpoint
	addrToShard := make(map[string]*ShardEndpoint)

	for i, sd := range s.DONs {
		connMgrs := s.ConnMgrs[i]
		for shardIdx, shard := range sd.Shards {
			ep := &ShardEndpoint{
				DonID:    config.ShardDONID(sd.DonName, shardIdx),
				DonName:  sd.DonName,
				ShardIdx: shardIdx,
				ConnMgr:  connMgrs[shardIdx],
				Members:  shard.Nodes,
				F:        sd.F,
			}
			endpoints = append(endpoints, ep)
			for _, member := range shard.Nodes {
				if existing, ok := addrToShard[member.Address]; ok {
					return nil, nil, fmt.Errorf("node address %s appears in both %s and %s; shard memberships must be disjoint", member.Address, existing.DonID, ep.DonID)
				}
				addrToShard[member.Address] = ep
			}
		}
	}

	return endpoints, addrToShard, nil
}

// AllMembers returns the union of all shard members across the matrix.
func AllMembers(endpoints []*ShardEndpoint) []config.NodeConfig {
	var out []config.NodeConfig
	for _, ep := range endpoints {
		out = append(out, ep.Members...)
	}
	return out
}
