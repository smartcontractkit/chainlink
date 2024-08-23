package gas

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups"
)

type ArbConfig interface {
	LimitMax() uint64
	BumpPercent() uint16
	BumpMin() *assets.Wei
}

// arbitrumEstimator is an Estimator which extends SuggestedPriceEstimator to use getPricesInArbGas() for gas limit estimation.
type arbitrumEstimator struct {
	services.StateMachine
	cfg ArbConfig

	EvmEstimator // *SuggestedPriceEstimator

	pollPeriod time.Duration
	logger     logger.Logger

	getPricesInArbGasMu sync.RWMutex
	perL2Tx             uint32
	perL1CalldataUnit   uint32

	chForceRefetch chan (chan struct{})
	chInitialised  chan struct{}
	chStop         services.StopChan
	chDone         chan struct{}

	l1Oracle rollups.ArbL1GasOracle
}

func NewArbitrumEstimator(lggr logger.Logger, cfg ArbConfig, ethClient feeEstimatorClient, l1Oracle rollups.ArbL1GasOracle) EvmEstimator {
	lggr = logger.Named(lggr, "ArbitrumEstimator")

	return &arbitrumEstimator{
		cfg:            cfg,
		EvmEstimator:   NewSuggestedPriceEstimator(lggr, ethClient, cfg, l1Oracle),
		pollPeriod:     10 * time.Second,
		logger:         lggr,
		chForceRefetch: make(chan (chan struct{})),
		chInitialised:  make(chan struct{}),
		chStop:         make(chan struct{}),
		chDone:         make(chan struct{}),
		l1Oracle:       l1Oracle,
	}
}

func (a *arbitrumEstimator) Name() string {
	return a.logger.Name()
}

func (a *arbitrumEstimator) Start(ctx context.Context) error {
	return a.StartOnce("ArbitrumEstimator", func() error {
		if err := a.EvmEstimator.Start(ctx); err != nil {
			return pkgerrors.Wrap(err, "failed to start gas price estimator")
		}
		go a.run()
		<-a.chInitialised
		return nil
	})
}
func (a *arbitrumEstimator) Close() error {
	return a.StopOnce("ArbitrumEstimator", func() (err error) {
		close(a.chStop)
		err = pkgerrors.Wrap(a.EvmEstimator.Close(), "failed to stop gas price estimator")
		<-a.chDone
		return
	})
}

func (a *arbitrumEstimator) Ready() error { return a.StateMachine.Ready() }

func (a *arbitrumEstimator) HealthReport() map[string]error {
	hp := map[string]error{a.Name(): a.Healthy()}
	services.CopyHealth(hp, a.EvmEstimator.HealthReport())
	return hp
}

// GetLegacyGas estimates both the gas price and the gas limit.
//   - Price is delegated to the embedded SuggestedPriceEstimator.
//   - Limit is computed from the dynamic values perL2Tx and perL1CalldataUnit, provided by the getPricesInArbGas() method
//     of the precompilie contract at ArbGasInfoAddress. perL2Tx is a constant amount of gas, and perL1CalldataUnit is
//     multiplied by the length of the tx calldata. The sum of these two values plus the original l2GasLimit is returned.
func (a *arbitrumEstimator) GetLegacyGas(ctx context.Context, calldata []byte, l2GasLimit uint64, maxGasPriceWei *assets.Wei, opts ...feetypes.Opt) (gasPrice *assets.Wei, chainSpecificGasLimit uint64, err error) {
	gasPrice, _, err = a.EvmEstimator.GetLegacyGas(ctx, calldata, l2GasLimit, maxGasPriceWei, opts...)
	if err != nil {
		return
	}
	gasPrice = a.gasPriceWithBuffer(gasPrice, maxGasPriceWei)
	ok := a.IfStarted(func() {
		if slices.Contains(opts, feetypes.OptForceRefetch) {
			ch := make(chan struct{})
			select {
			case a.chForceRefetch <- ch:
			case <-a.chStop:
				err = pkgerrors.New("estimator stopped")
				return
			case <-ctx.Done():
				err = ctx.Err()
				return
			}
			select {
			case <-ch:
			case <-a.chStop:
				err = pkgerrors.New("estimator stopped")
				return
			case <-ctx.Done():
				err = ctx.Err()
				return
			}
		}
		perL2Tx, perL1CalldataUnit := a.getPricesInArbGas()
		chainSpecificGasLimit = l2GasLimit + uint64(perL2Tx) + uint64(len(calldata))*uint64(perL1CalldataUnit)
		a.logger.Debugw("GetLegacyGas", "l2GasLimit", l2GasLimit, "calldataLen", len(calldata), "perL2Tx", perL2Tx,
			"perL1CalldataUnit", perL1CalldataUnit, "chainSpecificGasLimit", chainSpecificGasLimit)
	})
	if !ok {
		return nil, 0, pkgerrors.New("estimator is not started")
	} else if err != nil {
		return
	}
	if max := a.cfg.LimitMax(); chainSpecificGasLimit > max {
		err = fmt.Errorf("estimated gas limit: %d is greater than the maximum gas limit configured: %d", chainSpecificGasLimit, max)
		return
	}
	return
}

// During network congestion Arbitrum's suggested gas price can be extremely volatile, making gas estimations less accurate. For any transaction, Arbitrum will only charge
// the block's base fee. If the base fee increases rapidly there is a chance the suggested gas price will fall under that value, resulting in a fee too low error.
// We use gasPriceWithBuffer to increase the estimated gas price by some percentage to avoid fee too low errors. Eventually, only the base fee will be paid, regardless of the price.
func (a *arbitrumEstimator) gasPriceWithBuffer(gasPrice *assets.Wei, maxGasPriceWei *assets.Wei) *assets.Wei {
	const gasPriceBufferPercentage = 50

	gasPrice = gasPrice.AddPercentage(gasPriceBufferPercentage)
	if gasPrice.Cmp(maxGasPriceWei) > 0 {
		a.logger.Warnw("Updated gasPrice with buffer is higher than the max gas price limit. Falling back to max gas price", "gasPriceWithBuffer", gasPrice, "maxGasPriceWei", maxGasPriceWei)
		gasPrice = maxGasPriceWei
	}
	a.logger.Debugw("gasPriceWithBuffer", "updatedGasPrice", gasPrice)
	return gasPrice
}

func (a *arbitrumEstimator) getPricesInArbGas() (perL2Tx uint32, perL1CalldataUnit uint32) {
	a.getPricesInArbGasMu.RLock()
	perL2Tx, perL1CalldataUnit = a.perL2Tx, a.perL1CalldataUnit
	a.getPricesInArbGasMu.RUnlock()
	return
}

func (a *arbitrumEstimator) run() {
	defer close(a.chDone)

	a.refreshPricesInArbGas()
	close(a.chInitialised)

	t := services.TickerConfig{
		Initial:   a.pollPeriod,
		JitterPct: services.DefaultJitter,
	}.NewTicker(a.pollPeriod)
	defer t.Stop()

	for {
		select {
		case <-a.chStop:
			return
		case ch := <-a.chForceRefetch:
			a.refreshPricesInArbGas()
			t.Reset()
			close(ch)
		case <-t.C:
			a.refreshPricesInArbGas()
		}
	}
}

// refreshPricesInArbGas calls getPricesInArbGas() and caches the refreshed prices.
func (a *arbitrumEstimator) refreshPricesInArbGas() {
	perL2Tx, perL1CalldataUnit, err := a.l1Oracle.GetPricesInArbGas()
	if err != nil {
		a.logger.Warnw("Failed to refresh prices", "err", err)
		return
	}

	a.logger.Debugw("refreshPricesInArbGas", "perL2Tx", perL2Tx, "perL2CalldataUnit", perL1CalldataUnit)

	a.getPricesInArbGasMu.Lock()
	a.perL2Tx = perL2Tx
	a.perL1CalldataUnit = perL1CalldataUnit
	a.getPricesInArbGasMu.Unlock()
}
