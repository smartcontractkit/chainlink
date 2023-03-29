// arbgas takes a single URL argument and prints the result of three GetLegacyGas calls to the Arbitrum gas estimator.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/zap/zapcore"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func main() {
	if l := len(os.Args); l != 2 {
		log.Fatal("Expected one URL argument but got", l-1)
	}
	url := os.Args[1]
	lggr, sync := logger.NewLogger()
	defer func() { _ = sync() }()
	lggr.SetLogLevel(zapcore.DebugLevel)

	ctx := context.Background()
	withEstimator(ctx, logger.Sugared(lggr), url, func(e gas.EvmEstimator) {
		printGetLegacyGas(ctx, e, make([]byte, 10), 500_000, assets.GWei(1))
		printGetLegacyGas(ctx, e, make([]byte, 10), 500_000, assets.GWei(1), txmgrtypes.OptForceRefetch)
		printGetLegacyGas(ctx, e, make([]byte, 10), max, assets.GWei(1))
	})
}

func printGetLegacyGas(ctx context.Context, e gas.EvmEstimator, calldata []byte, l2GasLimit uint32, maxGasPrice *assets.Wei, opts ...txmgrtypes.Opt) {
	price, limit, err := e.GetLegacyGas(ctx, calldata, l2GasLimit, maxGasPrice, opts...)
	if err != nil {
		log.Println("failed to get legacy gas:", err)
		return
	}
	fmt.Println("Price:", price)
	fmt.Println("Limit:", limit)
}

const max = 50_000_000

func withEstimator(ctx context.Context, lggr logger.SugaredLogger, url string, f func(e gas.EvmEstimator)) {
	rc, err := rpc.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	ec := ethclient.NewClient(rc)
	e := gas.NewArbitrumEstimator(lggr, &config{max: max}, rc, ec)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	err = e.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer lggr.ErrorIfFn(e.Close, "Error closing ArbitrumEstimator")

	f(e)
}

var _ gas.ArbConfig = &config{}

type config struct {
	max uint32
}

func (c *config) EvmGasLimitMax() uint32 {
	return c.max
}
