package rollups

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/common/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
)

//go:generate mockery --quiet --name l1OracleClient --output ./mocks/ --case=underscore --structname L1OracleClient
type l1OracleClient interface {
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
}

type priceEntry struct {
	price     *assets.Wei
	timestamp time.Time
}

const (
	// Interval at which to poll for L1BaseFee. A good starting point is the L1 block time.
	PollPeriod = 6 * time.Second
)

var supportedChainTypes = []config.ChainType{config.ChainArbitrum, config.ChainOptimismBedrock, config.ChainKroma, config.ChainScroll}

func IsRollupWithL1Support(chainType config.ChainType) bool {
	return slices.Contains(supportedChainTypes, chainType)
}

func NewL1GasOracle(lggr logger.Logger, ethClient l1OracleClient, chainType config.ChainType) L1Oracle {
	var l1Oracle L1Oracle
	switch chainType {
	case config.ChainOptimismBedrock, config.ChainKroma:
		l1Oracle = NewOpStackL1GasOracle(lggr, ethClient, chainType)
	case config.ChainArbitrum:
		l1Oracle = NewArbitrumL1GasOracle(lggr, ethClient)
	case config.ChainScroll:
		l1Oracle = NewScrollL1GasOracle(lggr, ethClient)
	default:
		panic(fmt.Sprintf("Received unspported chaintype %s", chainType))
	}
	return l1Oracle
}
