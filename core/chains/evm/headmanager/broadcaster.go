package headmanager

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/common/headmanager"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type broadcaster = headmanager.Broadcaster[*evmtypes.Head, common.Hash]

var _ commontypes.Broadcaster[*evmtypes.Head, common.Hash] = &broadcaster{}

func NewBroadcaster(
	lggr logger.Logger,
) *broadcaster {
	return headmanager.NewBroadcaster[*evmtypes.Head, common.Hash](lggr)
}
