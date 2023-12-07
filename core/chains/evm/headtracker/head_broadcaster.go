package headtracker

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/common/headtracker"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type headBroadcaster = headtracker.HeadBroadcaster[*evmtypes.Head, common.Hash]

var _ commontypes.HeadBroadcaster[*evmtypes.Head, common.Hash] = &headBroadcaster{}

func NewHeadBroadcaster(
	lggr logger.Logger,
) *headBroadcaster {
	return headtracker.NewHeadBroadcaster[*evmtypes.Head, common.Hash](lggr)
}
