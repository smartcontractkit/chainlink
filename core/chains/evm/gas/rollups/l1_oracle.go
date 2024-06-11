package rollups

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
)

// L1Oracle provides interface for fetching L1-specific fee components if the chain is an L2.
// For example, on Optimistic Rollups, this oracle can return rollup-specific l1BaseFee
//
//go:generate mockery --quiet --name L1Oracle --output ./mocks/ --case=underscore
type L1Oracle interface {
	services.Service

	GasPrice(ctx context.Context) (*assets.Wei, error)
	GetGasCost(ctx context.Context, tx *types.Transaction, blockNum *big.Int) (*assets.Wei, error)
}

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

var supportedChainTypes = []chaintype.ChainType{chaintype.ChainArbitrum, chaintype.ChainOptimismBedrock, chaintype.ChainKroma, chaintype.ChainScroll}

func IsRollupWithL1Support(chainType chaintype.ChainType) bool {
	return slices.Contains(supportedChainTypes, chainType)
}

func NewL1GasOracle(lggr logger.Logger, ethClient l1OracleClient, chainType chaintype.ChainType) L1Oracle {
	if !IsRollupWithL1Support(chainType) {
		return nil
	}
	var l1Oracle L1Oracle
	switch chainType {
	case chaintype.ChainOptimismBedrock, chaintype.ChainKroma, chaintype.ChainScroll:
		l1Oracle = NewOpStackL1GasOracle(lggr, ethClient, chainType)
	case chaintype.ChainArbitrum:
		l1Oracle = NewArbitrumL1GasOracle(lggr, ethClient)
	default:
		panic(fmt.Sprintf("Received unspported chaintype %s", chainType))
	}
	return l1Oracle
}
