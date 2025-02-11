package headtracker

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-framework/chains/heads"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
)

type headBroadcaster = heads.Broadcaster[*evmtypes.Head, common.Hash]

func NewHeadBroadcaster(
	lggr logger.Logger,
) headBroadcaster {
	return heads.NewBroadcaster[*evmtypes.Head, common.Hash](lggr)
}
