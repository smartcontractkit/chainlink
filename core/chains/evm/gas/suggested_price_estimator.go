package gas

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	bigmath "github.com/smartcontractkit/chainlink-common/pkg/utils/big_math"

	"github.com/smartcontractkit/chainlink/v2/common/fee"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

var (
	_ EvmEstimator = &SuggestedPriceEstimator{}
)

type suggestedPriceConfig interface {
	BumpPercent() uint16
	BumpMin() *assets.Wei
}

type suggestedPriceEstimatorClient interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

// SuggestedPriceEstimator is an Estimator which uses the suggested gas price from eth_gasPrice.
type SuggestedPriceEstimator struct {
	services.StateMachine

	cfg        suggestedPriceConfig
	client     suggestedPriceEstimatorClient
	pollPeriod time.Duration
	logger     logger.Logger

	gasPriceMu sync.RWMutex
	GasPrice   *assets.Wei

	chForceRefetch chan (chan struct{})
	chInitialised  chan struct{}
	chStop         services.StopChan
	chDone         chan struct{}

	l1Oracle rollups.L1Oracle
}

// NewSuggestedPriceEstimator returns a new Estimator which uses the suggested gas price.
func NewSuggestedPriceEstimator(lggr logger.Logger, client feeEstimatorClient, cfg suggestedPriceConfig, l1Oracle rollups.L1Oracle) EvmEstimator {
	return &SuggestedPriceEstimator{
		client:         client,
		pollPeriod:     10 * time.Second,
		logger:         logger.Named(lggr, "SuggestedPriceEstimator"),
		cfg:            cfg,
		chForceRefetch: make(chan (chan struct{})),
		chInitialised:  make(chan struct{}),
		chStop:         make(chan struct{}),
		chDone:         make(chan struct{}),
		l1Oracle:       l1Oracle,
	}
}

func (o *SuggestedPriceEstimator) Name() string {
	return o.logger.Name()
}

func (o *SuggestedPriceEstimator) L1Oracle() rollups.L1Oracle {
	return o.l1Oracle
}

func (o *SuggestedPriceEstimator) Start(context.Context) error {
	return o.StartOnce("SuggestedPriceEstimator", func() error {
		go o.run()
		<-o.chInitialised
		return nil
	})
}
func (o *SuggestedPriceEstimator) Close() error {
	return o.StopOnce("SuggestedPriceEstimator", func() error {
		close(o.chStop)
		<-o.chDone
		return nil
	})
}

func (o *SuggestedPriceEstimator) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.Healthy()}
}

func (o *SuggestedPriceEstimator) run() {
	defer close(o.chDone)

	o.refreshPrice()
	close(o.chInitialised)

	t := services.TickerConfig{
		Initial:   o.pollPeriod,
		JitterPct: services.DefaultJitter,
	}.NewTicker(o.pollPeriod)

	for {
		select {
		case <-o.chStop:
			return
		case ch := <-o.chForceRefetch:
			o.refreshPrice()
			t.Reset()
			close(ch)
		case <-t.C:
			o.refreshPrice()
		}
	}
}

func (o *SuggestedPriceEstimator) refreshPrice() {
	var res hexutil.Big
	ctx, cancel := o.chStop.CtxCancel(evmclient.ContextWithDefaultTimeout())
	defer cancel()

	if err := o.client.CallContext(ctx, &res, "eth_gasPrice"); err != nil {
		o.logger.Warnf("Failed to refresh prices, got error: %s", err)
		return
	}
	bi := (*assets.Wei)(&res)

	o.logger.Debugw("refreshPrice", "GasPrice", bi)

	o.gasPriceMu.Lock()
	defer o.gasPriceMu.Unlock()
	o.GasPrice = bi
}

// Uses the force refetch chan to trigger a price update and blocks until complete
func (o *SuggestedPriceEstimator) forceRefresh(ctx context.Context) (err error) {
	ch := make(chan struct{})
	select {
	case o.chForceRefetch <- ch:
	case <-o.chStop:
		return pkgerrors.New("estimator stopped")
	case <-ctx.Done():
		return ctx.Err()
	}
	select {
	case <-ch:
	case <-o.chStop:
		return pkgerrors.New("estimator stopped")
	case <-ctx.Done():
		return ctx.Err()
	}
	return
}

func (o *SuggestedPriceEstimator) OnNewLongestChain(context.Context, *evmtypes.Head) {}

func (*SuggestedPriceEstimator) GetDynamicFee(_ context.Context, _ *assets.Wei) (fee DynamicFee, err error) {
	err = pkgerrors.New("dynamic fees are not implemented for this estimator")
	return
}

func (*SuggestedPriceEstimator) BumpDynamicFee(_ context.Context, _ DynamicFee, _ *assets.Wei, _ []EvmPriorAttempt) (bumped DynamicFee, err error) {
	err = pkgerrors.New("dynamic fees are not implemented for this estimator")
	return
}

func (o *SuggestedPriceEstimator) GetLegacyGas(ctx context.Context, _ []byte, GasLimit uint64, maxGasPriceWei *assets.Wei, opts ...feetypes.Opt) (gasPrice *assets.Wei, chainSpecificGasLimit uint64, err error) {
	chainSpecificGasLimit = GasLimit
	ok := o.IfStarted(func() {
		if slices.Contains(opts, feetypes.OptForceRefetch) {
			err = o.forceRefresh(ctx)
		}
		if gasPrice = o.getGasPrice(); gasPrice == nil {
			err = pkgerrors.New("failed to estimate gas; gas price not set")
			return
		}
		o.logger.Debugw("GetLegacyGas", "GasPrice", gasPrice, "GasLimit", GasLimit)
	})
	if !ok {
		return nil, 0, pkgerrors.New("estimator is not started")
	} else if err != nil {
		return
	}
	// For L2 chains, submitting a transaction that is not priced high enough will cause the call to fail, so if the cap is lower than the RPC suggested gas price, this transaction cannot succeed
	if gasPrice != nil && gasPrice.Cmp(maxGasPriceWei) > 0 {
		return nil, 0, pkgerrors.Errorf("estimated gas price: %s is greater than the maximum gas price configured: %s", gasPrice.String(), maxGasPriceWei.String())
	}
	return
}

// Refreshes the gas price by making a call to the RPC in case the current one has gone stale.
// Adds the larger of BumpPercent and BumpMin configs as a buffer on top of the price returned from the RPC.
// The only reason bumping logic would be called on the SuggestedPriceEstimator is if there was a significant price spike
// between the last price update and when the tx was submitted. Refreshing the price helps ensure the latest market changes are accounted for.
func (o *SuggestedPriceEstimator) BumpLegacyGas(ctx context.Context, originalFee *assets.Wei, feeLimit uint64, maxGasPriceWei *assets.Wei, _ []EvmPriorAttempt) (newGasPrice *assets.Wei, chainSpecificGasLimit uint64, err error) {
	chainSpecificGasLimit = feeLimit
	ok := o.IfStarted(func() {
		// Immediately return error if original fee is greater than or equal to the max gas price
		// Prevents a loop of resubmitting the attempt with the max gas price
		if originalFee.Cmp(maxGasPriceWei) >= 0 {
			err = fmt.Errorf("original fee (%s) greater than or equal to max gas price (%s) so cannot be bumped further", originalFee.String(), maxGasPriceWei.String())
			return
		}
		err = o.forceRefresh(ctx)
		if newGasPrice = o.getGasPrice(); newGasPrice == nil {
			err = pkgerrors.New("failed to refresh and return gas; gas price not set")
			return
		}
		o.logger.Debugw("BumpLegacyGas", "GasPrice", newGasPrice, "GasLimit", feeLimit)
	})
	if !ok {
		return nil, 0, pkgerrors.New("estimator is not started")
	} else if err != nil {
		return
	}
	if newGasPrice != nil && newGasPrice.Cmp(maxGasPriceWei) > 0 {
		return nil, 0, pkgerrors.Errorf("estimated gas price: %s is greater than the maximum gas price configured: %s", newGasPrice.String(), maxGasPriceWei.String())
	}
	// Add a buffer on top of the gas price returned by the RPC.
	// Bump logic when using the suggested gas price from an RPC is realistically only needed when there is increased volatility in gas price.
	// This buffer is a precaution to increase the chance of getting this tx on chain
	bufferedPrice := fee.MaxBumpedFee(newGasPrice.ToInt(), o.cfg.BumpPercent(), o.cfg.BumpMin().ToInt())
	// If the new suggested price is less than or equal to the max and the buffer puts the new price over the max, return the max price instead
	// The buffer is added on top of the suggested price during bumping as just a precaution. It is better to resubmit the transaction with the max gas price instead of erroring.
	newGasPrice = assets.NewWei(bigmath.Min(bufferedPrice, maxGasPriceWei.ToInt()))
	// Return the original price if the refreshed price with the buffer is lower to ensure the bumped gas price is always equal or higher to the previous attempt
	if originalFee != nil && originalFee.Cmp(newGasPrice) > 0 {
		return originalFee, chainSpecificGasLimit, nil
	}
	return
}

func (o *SuggestedPriceEstimator) getGasPrice() (GasPrice *assets.Wei) {
	o.gasPriceMu.RLock()
	defer o.gasPriceMu.RUnlock()
	return o.GasPrice
}
