package client

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// RPCClient includes all the necessary generalized RPC methods used by Node to perform health checks
type RPCClient[
	CHAIN_ID types.ID,
	HEAD Head,
] interface {
	// ChainID - fetches ChainID from the RPC to verify that it matches config
	ChainID(ctx context.Context) (CHAIN_ID, error)
	// Dial - prepares the RPC for usage. Can be called on fresh or closed RPC
	Dial(ctx context.Context) error
	// SubscribeToHeads - returns channel and subscription for new heads.
	SubscribeToHeads(ctx context.Context) (<-chan HEAD, types.Subscription, error)
	// SubscribeToFinalizedHeads - returns channel and subscription for finalized heads.
	SubscribeToFinalizedHeads(ctx context.Context) (<-chan HEAD, types.Subscription, error)
	// Ping - returns error if RPC is not reachable
	Ping(context.Context) error
	// IsSyncing - returns true if the RPC is in Syncing state and can not process calls
	IsSyncing(ctx context.Context) (bool, error)
	// UnsubscribeAllExcept - close all subscriptions except `subs`
	UnsubscribeAllExcept(subs ...types.Subscription)
	// Close - closes all subscriptions and aborts all RPC calls
	Close()
	// GetInterceptedChainInfo - returns latest and highest observed by application layer ChainInfo.
	// latest ChainInfo is the most recent value received within a NodeClient's current lifecycle between Dial and DisconnectAll.
	// highestUserObservations ChainInfo is the highest ChainInfo observed excluding health checks calls.
	// Its values must not be reset.
	// The results of corresponding calls, to get the most recent head and the latest finalized head, must be
	// intercepted and reflected in ChainInfo before being returned to a caller. Otherwise, MultiNode is not able to
	// provide repeatable read guarantee.
	// DisconnectAll must reset latest ChainInfo to default value.
	// Ensure implementation does not have a race condition when values are reset before request completion and as
	// a result latest ChainInfo contains information from the previous cycle.
	GetInterceptedChainInfo() (latest, highestUserObservations ChainInfo)
}

// Head is the interface required by the NodeClient
type Head interface {
	BlockNumber() int64
	BlockDifficulty() *big.Int
	IsValid() bool
}

// PoolChainInfoProvider - provides aggregation of nodes pool ChainInfo
type PoolChainInfoProvider interface {
	// LatestChainInfo - returns number of live nodes available in the pool, so we can prevent the last alive node in a pool from being
	// moved to out-of-sync state. It is better to have one out-of-sync node than no nodes at all.
	// Returns highest latest ChainInfo within the alive nodes. E.g. most recent block number and highest block number
	// observed by Node A are 10 and 15; Node B - 12 and 14. This method will return 12.
	LatestChainInfo() (int, ChainInfo)
	// HighestUserObservations - returns highest ChainInfo ever observed by any user of MultiNode.
	HighestUserObservations() ChainInfo
}

// ChainInfo - defines RPC's or MultiNode's view on the chain
type ChainInfo struct {
	BlockNumber          int64
	FinalizedBlockNumber int64
	TotalDifficulty      *big.Int
}

func MaxTotalDifficulty(a, b *big.Int) *big.Int {
	if a == nil {
		if b == nil {
			return nil
		}

		return big.NewInt(0).Set(b)
	}

	if b == nil || a.Cmp(b) >= 0 {
		return big.NewInt(0).Set(a)
	}

	return big.NewInt(0).Set(b)
}
