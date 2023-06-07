package headtracker

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/common/headtracker"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type evmHeadBroadcaster = headtracker.HeadBroadcaster[*evmtypes.Head, common.Hash]

var _ commontypes.HeadBroadcaster[*evmtypes.Head, common.Hash] = &evmHeadBroadcaster{}

func NewEvmHeadBroadcaster(
	lggr logger.Logger,
) *evmHeadBroadcaster {
	return headtracker.NewHeadBroadcaster[*evmtypes.Head, common.Hash](lggr)
}
