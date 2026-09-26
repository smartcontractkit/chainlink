package handlers

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

// NodeRouter maps node addresses to the connection manager that owns them, across
// every DON and shard a handler serves. It lets a handler reach (or route a
// response back to) any node in the DON×shard matrix without hardcoding a single
// connection manager, which breaks once a DON has more than one shard.
type NodeRouter struct {
	addrToDON map[string]DON
	members   []config.NodeConfig
}

// BuildNodeRouter expands shardedDONs/shardsConnMgrs into a single node address ->
// connection manager lookup. The outer dimension of shardsConnMgrs indexes DONs
// (parallel to shardedDONs); the inner dimension indexes shards within that DON
// (parallel to shardedDONs[i].Shards). Node addresses must be unique across the
// whole matrix.
func BuildNodeRouter(shardedDONs []config.ShardedDONConfig, shardsConnMgrs [][]DON) (*NodeRouter, error) {
	if len(shardedDONs) == 0 || len(shardsConnMgrs) == 0 {
		return nil, errors.New("at least one DON and connection manager required")
	}
	if len(shardedDONs) != len(shardsConnMgrs) {
		return nil, fmt.Errorf("shardedDONs length %d does not match shardsConnMgrs length %d", len(shardedDONs), len(shardsConnMgrs))
	}

	router := &NodeRouter{addrToDON: make(map[string]DON)}
	for i, sd := range shardedDONs {
		connMgrs := shardsConnMgrs[i]
		if len(sd.Shards) != len(connMgrs) {
			return nil, fmt.Errorf("DON %s: shards length %d does not match connection managers length %d", sd.DonName, len(sd.Shards), len(connMgrs))
		}
		for shardIdx, shard := range sd.Shards {
			if connMgrs[shardIdx] == nil {
				return nil, fmt.Errorf("DON %s shard %d: nil connection manager", sd.DonName, shardIdx)
			}
			for _, member := range shard.Nodes {
				if _, exists := router.addrToDON[member.Address]; exists {
					return nil, fmt.Errorf("node address %s appears in more than one shard; shard memberships must be disjoint", member.Address)
				}
				router.addrToDON[member.Address] = connMgrs[shardIdx]
				router.members = append(router.members, member)
			}
		}
	}
	return router, nil
}

// DONFor returns the connection manager that owns nodeAddr, if any.
func (r *NodeRouter) DONFor(nodeAddr string) (DON, bool) {
	don, ok := r.addrToDON[nodeAddr]
	return don, ok
}

// Members returns every node across the full DON×shard matrix.
func (r *NodeRouter) Members() []config.NodeConfig {
	return r.members
}
