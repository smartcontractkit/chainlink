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
type L1Oracle interface {
	services.Service

	GasPrice(ctx context.Context) (*assets.Wei, error)
	GetGasCost(ctx context.Context, tx *types.Transaction, blockNum *big.Int) (*assets.Wei, error)
	ChainType(ctx context.Context) chaintype.ChainType
}

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

var supportedChainTypes = []chaintype.ChainType{chaintype.ChainArbitrum, chaintype.ChainOptimismBedrock, chaintype.ChainKroma, chaintype.ChainScroll, chaintype.ChainZkSync, chaintype.ChainMantle}

func IsRollupWithL1Support(chainType chaintype.ChainType) bool {
	return slices.Contains(supportedChainTypes, chainType)
}

func NewL1GasOracle(lggr logger.Logger, ethClient l1OracleClient, chainType chaintype.ChainType) (L1Oracle, error) {
	if !IsRollupWithL1Support(chainType) {
		return nil, nil
	}
	var l1Oracle L1Oracle
	var err error
	switch chainType {
	case chaintype.ChainOptimismBedrock, chaintype.ChainKroma, chaintype.ChainScroll, chaintype.ChainMantle, chaintype.ChainZircuit:
		l1Oracle, err = NewOpStackL1GasOracle(lggr, ethClient, chainType)
	case chaintype.ChainArbitrum:
		l1Oracle, err = NewArbitrumL1GasOracle(lggr, ethClient)
	case chaintype.ChainZkSync:
		l1Oracle = NewZkSyncL1GasOracle(lggr, ethClient)
	default:
		return nil, fmt.Errorf("received unsupported chaintype %s", chainType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initialize L1 oracle for chaintype %s: %w", chainType, err)
	}
	return l1Oracle, nil
}
