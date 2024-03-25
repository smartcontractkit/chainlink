package vrfv2plus

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func WaitForRequestCountEqualToFulfilmentCount(
	ctx context.Context,
	consumer LoadTestConsumer,
	timeout time.Duration,
	wg *sync.WaitGroup,
) (*big.Int, *big.Int, error) {
	metricsChannel := make(chan *contracts.VRFV2PlusLoadTestMetrics)
	metricsErrorChannel := make(chan error)

	testContext, testCancel := context.WithTimeout(ctx, timeout)
	defer testCancel()

	ticker := time.NewTicker(time.Second * 1)
	var metrics *contracts.VRFV2PlusLoadTestMetrics
	for {
		select {
		case <-testContext.Done():
			ticker.Stop()
			wg.Done()
			return metrics.RequestCount, metrics.FulfilmentCount,
				fmt.Errorf("timeout waiting for rand request and fulfilments to be equal AFTER performance test was executed. Request Count: %d, Fulfilment Count: %d",
					metrics.RequestCount.Uint64(), metrics.FulfilmentCount.Uint64())
		case <-ticker.C:
			go retrieveLoadTestMetrics(ctx, consumer, metricsChannel, metricsErrorChannel)
		case metrics = <-metricsChannel:
			if metrics.RequestCount.Cmp(metrics.FulfilmentCount) == 0 {
				ticker.Stop()
				wg.Done()
				return metrics.RequestCount, metrics.FulfilmentCount, nil
			}
		case err := <-metricsErrorChannel:
			ticker.Stop()
			wg.Done()
			return nil, nil, err
		}
	}
}

func retrieveLoadTestMetrics(
	ctx context.Context,
	consumer LoadTestConsumer,
	metricsChannel chan *contracts.VRFV2PlusLoadTestMetrics,
	metricsErrorChannel chan error,
) {
	metrics, err := consumer.GetLoadTestMetrics(ctx)
	if err != nil {
		metricsErrorChannel <- err
	}
	metricsChannel <- metrics
}
