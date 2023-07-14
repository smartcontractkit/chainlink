package headmanager

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/common/headmanager"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type listener = headmanager.Listener[*evmtypes.Head, ethereum.Subscription, *big.Int, common.Hash]

var _ commontypes.Listener[*evmtypes.Head, common.Hash] = (*listener)(nil)

func NewListener(
	lggr logger.Logger,
	ethClient evmclient.Client,
	config Config, chStop chan struct{},
) *listener {
	return headmanager.NewListener[
		*evmtypes.Head,
		ethereum.Subscription, *big.Int, common.Hash,
	](lggr, ethClient, config, chStop)
}
