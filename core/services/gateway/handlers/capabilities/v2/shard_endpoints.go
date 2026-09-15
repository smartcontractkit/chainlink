package v2

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

type shardEndpoint struct {
	// donID is the shard-specific DON identifier, derived via config.ShardDONID.
	// Shard 0 uses the bare DON name; shard N>0 uses "donName_N".
	donID    string
	donName  string
	shardIdx int
	// connMgr is the connection manager for this shard. It only knows how to reach
	// nodes that are members of this shard.
	connMgr handlers.DON
	members []config.NodeConfig
	f       int
}

// buildShardEndpoints expands the full DON×shard matrix into a flat list of
// shardEndpoints, and builds a nodeAddr -> shardEndpoint lookup. The outer
// dimension of shardsConnMgrs indexes DONs (parallel to shardedDONs); the inner
// dimension indexes shards within that DON (parallel to shardedDONs[i].Shards).
// Memberships must be disjoint.
func buildShardEndpoints(shardedDONs []config.ShardedDONConfig, shardsConnMgrs [][]handlers.DON) ([]*shardEndpoint, map[string]*shardEndpoint, error) {
	if len(shardedDONs) == 0 || len(shardsConnMgrs) == 0 {
		return nil, nil, errors.New("at least one DON and connection manager required")
	}
	if len(shardedDONs) != len(shardsConnMgrs) {
		return nil, nil, fmt.Errorf("shardedDONs length %d does not match shardsConnMgrs length %d", len(shardedDONs), len(shardsConnMgrs))
	}

	var endpoints []*shardEndpoint
	addrToShard := make(map[string]*shardEndpoint)

	for i, sd := range shardedDONs {
		connMgrs := shardsConnMgrs[i]
		if len(sd.Shards) != len(connMgrs) {
			return nil, nil, fmt.Errorf("DON %s: shards length %d does not match connection managers length %d", sd.DonName, len(sd.Shards), len(connMgrs))
		}
		for shardIdx, shard := range sd.Shards {
			if connMgrs[shardIdx] == nil {
				return nil, nil, fmt.Errorf("DON %s shard %d: nil connection manager", sd.DonName, shardIdx)
			}
			ep := &shardEndpoint{
				donID:    config.ShardDONID(sd.DonName, shardIdx),
				donName:  sd.DonName,
				shardIdx: shardIdx,
				connMgr:  connMgrs[shardIdx],
				members:  shard.Nodes,
				f:        sd.F,
			}
			endpoints = append(endpoints, ep)
			for _, member := range shard.Nodes {
				if existing, ok := addrToShard[member.Address]; ok {
					return nil, nil, fmt.Errorf("node address %s appears in both %s and %s; shard memberships must be disjoint", member.Address, existing.donID, ep.donID)
				}
				addrToShard[member.Address] = ep
			}
		}
	}

	return endpoints, addrToShard, nil
}

// allMembers returns the union of all shard members across the matrix.
func allMembers(endpoints []*shardEndpoint) []config.NodeConfig {
	var out []config.NodeConfig
	for _, ep := range endpoints {
		out = append(out, ep.members...)
	}
	return out
}
