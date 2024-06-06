package headtracker

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/common/headtracker"

	htrktypes "github.com/smartcontractkit/chainlink/v2/common/headtracker/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type headListener = headtracker.HeadListener[*evmtypes.Head, common.Hash]

func NewHeadListener(
	lggr logger.Logger,
	ethClient evmclient.Client,
	config htrktypes.Config, chStop chan struct{},
) headListener {
	return headtracker.NewHeadListener[
		*evmtypes.Head,
		ethereum.Subscription, *big.Int, common.Hash,
	](lggr, ethClient, config, chStop)
}
