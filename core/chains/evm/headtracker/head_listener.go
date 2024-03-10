package headtracker

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/common/headtracker"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type headListener = headtracker.HeadListener[*evmtypes.Head, ethereum.Subscription, *big.Int, common.Hash]

var _ commontypes.HeadListener[*evmtypes.Head, common.Hash] = (*headListener)(nil)

func NewHeadListener(
	lggr logger.Logger,
	ethClient evmclient.Client,
	config Config, chStop chan struct{},
) *headListener {
	return headtracker.NewHeadListener[
		*evmtypes.Head,
		ethereum.Subscription, *big.Int, common.Hash,
	](lggr, ethClient, config, chStop)
}
