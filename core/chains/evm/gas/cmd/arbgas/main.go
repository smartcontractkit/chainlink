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

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/core/logger"
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
	withEstimator(ctx, logger.Sugared(lggr), url, func(e gas.Estimator) {
		printGetLegacyGas(ctx, e, make([]byte, 10), 500_000, assets.GWei(1))
		printGetLegacyGas(ctx, e, make([]byte, 10), 500_000, assets.GWei(1), gas.OptForceRefetch)
		printGetLegacyGas(ctx, e, make([]byte, 10), max, assets.GWei(1))
	})
}

func printGetLegacyGas(ctx context.Context, e gas.Estimator, calldata []byte, l2GasLimit uint32, maxGasPrice *assets.Wei, opts ...gas.Opt) {
	price, limit, err := e.GetLegacyGas(ctx, calldata, l2GasLimit, maxGasPrice, opts...)
	if err != nil {
		log.Println("failed to get legacy gas:", err)
		return
	}
	fmt.Println("Price:", price)
	fmt.Println("Limit:", limit)
}

const max = 50_000_000

func withEstimator(ctx context.Context, lggr logger.SugaredLogger, url string, f func(e gas.Estimator)) {
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
