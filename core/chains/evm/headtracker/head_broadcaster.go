package headtracker

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/common/headtracker"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type headBroadcaster = headtracker.HeadBroadcaster[*evmtypes.Head, common.Hash]

func NewHeadBroadcaster(
	lggr logger.Logger,
) headBroadcaster {
	return headtracker.NewHeadBroadcaster[*evmtypes.Head, common.Hash](lggr)
}
