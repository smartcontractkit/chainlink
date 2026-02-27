package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

// ShardInfo holds the topology and connection manager for a single shard.
type ShardInfo struct {
	DON     DON
	Members []config.NodeConfig
	F       int
}

// ShardRouter routes SendToNode calls to the correct shard's connection manager
// and exposes per-shard topology for shard-aware fan-out and aggregation.
type ShardRouter struct {
	shards        []ShardInfo
	nodeToConnMgr map[string]DON
	nodeToShard   map[string]int // lowercase node address -> shard index
}

// NewShardRouter creates a ShardRouter from sharded DON configs and their connection managers.
// All shards across all DONs are flattened into a single ordered list.
// shardsConnMgrs[donIdx][shardIdx] is the connection manager for DON donIdx, shard shardIdx.
func NewShardRouter(shardedDONs []config.ShardedDONConfig, shardsConnMgrs [][]DON) (*ShardRouter, error) {
	nodeToConnMgr := make(map[string]DON)
	nodeToShard := make(map[string]int)
	var shards []ShardInfo
	shardIndex := 0

	for donIdx, donCfg := range shardedDONs {
		if donIdx >= len(shardsConnMgrs) {
			return nil, fmt.Errorf("missing connection managers for DON %s", donCfg.DonName)
		}
		for shardIdx, shard := range donCfg.Shards {
			if shardIdx >= len(shardsConnMgrs[donIdx]) {
				return nil, fmt.Errorf("missing connection manager for DON %s shard %d", donCfg.DonName, shardIdx)
			}
			connMgr := shardsConnMgrs[donIdx][shardIdx]
			for _, node := range shard.Nodes {
				addr := strings.ToLower(node.Address)
				if _, exists := nodeToConnMgr[addr]; exists {
					return nil, fmt.Errorf("duplicate node address %s across shards", addr)
				}
				nodeToConnMgr[addr] = connMgr
				nodeToShard[addr] = shardIndex
			}
			shards = append(shards, ShardInfo{
				DON:     connMgr,
				Members: shard.Nodes,
				F:       donCfg.F,
			})
			shardIndex++
		}
	}
	return &ShardRouter{
		shards:        shards,
		nodeToConnMgr: nodeToConnMgr,
		nodeToShard:   nodeToShard,
	}, nil
}

func (r *ShardRouter) SendToNode(ctx context.Context, nodeAddress string, req *jsonrpc.Request[json.RawMessage]) error {
	connMgr, ok := r.nodeToConnMgr[strings.ToLower(nodeAddress)]
	if !ok {
		return fmt.Errorf("node %s not found in any shard", nodeAddress)
	}
	return connMgr.SendToNode(ctx, nodeAddress, req)
}

func (r *ShardRouter) Shard(idx int) ShardInfo {
	return r.shards[idx]
}

func (r *ShardRouter) NumShards() int {
	return len(r.shards)
}

// ShardIndexForNode returns the shard index that the given node belongs to.
func (r *ShardRouter) ShardIndexForNode(nodeAddr string) (int, bool) {
	idx, ok := r.nodeToShard[strings.ToLower(nodeAddr)]
	return idx, ok
}

// AllMembers returns all node configs across all shards (for metrics initialization).
func (r *ShardRouter) AllMembers() []config.NodeConfig {
	var members []config.NodeConfig
	for _, s := range r.shards {
		members = append(members, s.Members...)
	}
	return members
}

// DonID returns a summary ID for logging purposes.
func (r *ShardRouter) DonID() string {
	return fmt.Sprintf("sharded(%d)", len(r.shards))
}
