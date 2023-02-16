package fees

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	feePolling = 5 * time.Second // TODO: make configurable
)

var _ Estimator = &recentFeeEstimator{}

type recentFeeEstimator struct {
	starter utils.StartStopOnce
	chStop  chan struct{}
	done    sync.WaitGroup

	cfg config.Config

	price uint64
	lock  sync.RWMutex
}

func NewRecentFeeEstimator(cfg config.Config) (Estimator, error) {
	return &recentFeeEstimator{
		chStop: make(chan struct{}),
	}, fmt.Errorf("estimator not available - RPC method not released") // TODO: implement when RPC method available
}

func (est *recentFeeEstimator) Start(ctx context.Context) error {
	return est.starter.StartOnce("solana_recentFeeEstimator", func() error {
		est.done.Add(1)
		go est.run()
		return nil
	})
}

func (est *recentFeeEstimator) run() {
	defer est.done.Done()

	tick := time.After(0)
	for {
		select {
		case <-est.chStop:
			return
		case <-tick:
			// TODO: query endpoint - not available yet

			est.lock.Lock()
			est.price = 0
			est.lock.Unlock()
		}

		tick = time.After(utils.WithJitter(feePolling))
	}
}

func (est *recentFeeEstimator) Close() error {
	close(est.chStop)
	est.done.Wait()
	return nil
}

func (est *recentFeeEstimator) BaseComputeUnitPrice() uint64 {
	est.lock.RLock()
	defer est.lock.RUnlock()

	if est.price >= est.cfg.ComputeUnitPriceMin() && est.price <= est.cfg.ComputeUnitPriceMax() {
		return est.price
	}

	if est.price < est.cfg.ComputeUnitPriceMin() {
		return est.cfg.ComputeUnitPriceMin()
	}

	return est.cfg.ComputeUnitPriceMax()
}
