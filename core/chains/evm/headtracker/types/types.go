package types

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/common/headtracker"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// HeadSaver maintains chains persisted in DB. All methods are thread-safe.
type HeadSaver interface {
	headtracker.HeadSaver[*evmtypes.Head, common.Hash]
	// LatestHeadFromDB returns the highest seen head from DB.
	LatestHeadFromDB(ctx context.Context) (*evmtypes.Head, error)
}

// Type Alias for EVM Head Tracker Components
type (
	HeadTracker     = headtracker.HeadTracker[*evmtypes.Head, common.Hash]
	HeadTrackable   = headtracker.HeadTrackable[*evmtypes.Head, common.Hash]
	HeadListener    = headtracker.HeadListener[*evmtypes.Head, common.Hash]
	HeadBroadcaster = headtracker.HeadBroadcaster[*evmtypes.Head, common.Hash]
)
