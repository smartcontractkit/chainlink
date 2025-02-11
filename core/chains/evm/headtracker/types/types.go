package types

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-framework/chains/heads"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
)

// HeadSaver maintains chains persisted in DB. All methods are thread-safe.
type HeadSaver interface {
	heads.Saver[*evmtypes.Head, common.Hash]
	// LatestHeadFromDB returns the highest seen head from DB.
	LatestHeadFromDB(ctx context.Context) (*evmtypes.Head, error)
}

// Type Alias for EVM Head Tracker Components
type (
	HeadTracker     = heads.Tracker[*evmtypes.Head, common.Hash]
	HeadTrackable   = heads.Trackable[*evmtypes.Head, common.Hash]
	HeadListener    = heads.Listener[*evmtypes.Head, common.Hash]
	HeadBroadcaster = heads.Broadcaster[*evmtypes.Head, common.Hash]
	Client          = heads.Client[*evmtypes.Head, ethereum.Subscription, *big.Int, common.Hash]
)
