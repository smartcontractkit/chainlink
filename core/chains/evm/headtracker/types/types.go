package types

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// HeadSaver maintains chains persisted in DB. All methods are thread-safe.
type HeadSaver interface {
	commontypes.HeadSaver[*evmtypes.Head, common.Hash]
	// LatestHeadFromDB returns the highest seen head from DB.
	LatestHeadFromDB(ctx context.Context) (*evmtypes.Head, error)
}

// Type Alias for EVM Head Tracker Components
type (
	HeadBroadcasterRegistry = commontypes.HeadBroadcasterRegistry[*evmtypes.Head, common.Hash]
	HeadTracker             = commontypes.HeadTracker[*evmtypes.Head, common.Hash]
	HeadTrackable           = commontypes.HeadTrackable[*evmtypes.Head, common.Hash]
	HeadListener            = commontypes.HeadListener[*evmtypes.Head, common.Hash]
	HeadBroadcaster         = commontypes.HeadBroadcaster[*evmtypes.Head, common.Hash]
)
