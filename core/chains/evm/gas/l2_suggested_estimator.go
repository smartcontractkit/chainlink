package gas

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	_ EvmEstimator = &l2SuggestedPriceEstimator{}
)

//go:generate mockery --quiet --name rpcClient --output ./mocks/ --case=underscore --structname RPCClient
type rpcClient interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

// l2SuggestedPriceEstimator is an Estimator which uses the L2 suggested gas price from eth_gasPrice.
type l2SuggestedPriceEstimator struct {
	utils.StartStopOnce

	client     rpcClient
	pollPeriod time.Duration
	logger     logger.Logger

	gasPriceMu sync.RWMutex
	l2GasPrice *assets.Wei

	chForceRefetch chan (chan struct{})
	chInitialised  chan struct{}
	chStop         utils.StopChan
	chDone         chan struct{}
}

// NewL2SuggestedPriceEstimator returns a new Estimator which uses the L2 suggested gas price.
func NewL2SuggestedPriceEstimator(lggr logger.Logger, client rpcClient) EvmEstimator {
	return &l2SuggestedPriceEstimator{
		client:         client,
		pollPeriod:     10 * time.Second,
		logger:         lggr.Named("L2SuggestedEstimator"),
		chForceRefetch: make(chan (chan struct{})),
		chInitialised:  make(chan struct{}),
		chStop:         make(chan struct{}),
		chDone:         make(chan struct{}),
	}
}

func (o *l2SuggestedPriceEstimator) Name() string {
	return o.logger.Name()
}

func (o *l2SuggestedPriceEstimator) Start(context.Context) error {
	return o.StartOnce("L2SuggestedEstimator", func() error {
		go o.run()
		<-o.chInitialised
		return nil
	})
}
func (o *l2SuggestedPriceEstimator) Close() error {
	return o.StopOnce("L2SuggestedEstimator", func() error {
		close(o.chStop)
		<-o.chDone
		return nil
	})
}

func (o *l2SuggestedPriceEstimator) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.StartStopOnce.Healthy()}
}

func (o *l2SuggestedPriceEstimator) run() {
	defer close(o.chDone)

	t := o.refreshPrice()
	close(o.chInitialised)

	for {
		select {
		case <-o.chStop:
			return
		case ch := <-o.chForceRefetch:
			t.Stop()
			t = o.refreshPrice()
			close(ch)
		case <-t.C:
			t = o.refreshPrice()
		}
	}
}

func (o *l2SuggestedPriceEstimator) refreshPrice() (t *time.Timer) {
	t = time.NewTimer(utils.WithJitter(o.pollPeriod))

	var res hexutil.Big
	ctx, cancel := o.chStop.CtxCancel(evmclient.ContextWithDefaultTimeout())
	defer cancel()

	if err := o.client.CallContext(ctx, &res, "eth_gasPrice"); err != nil {
		o.logger.Warnf("Failed to refresh prices, got error: %s", err)
		return
	}
	bi := (*assets.Wei)(&res)

	o.logger.Debugw("refreshPrice", "l2GasPrice", bi)

	o.gasPriceMu.Lock()
	defer o.gasPriceMu.Unlock()
	o.l2GasPrice = bi
	return
}

func (o *l2SuggestedPriceEstimator) OnNewLongestChain(context.Context, *evmtypes.Head) {}

func (*l2SuggestedPriceEstimator) GetDynamicFee(_ context.Context, _ uint32, _ *assets.Wei) (fee DynamicFee, chainSpecificGasLimit uint32, err error) {
	err = errors.New("dynamic fees are not implemented for this layer 2")
	return
}

func (*l2SuggestedPriceEstimator) BumpDynamicFee(_ context.Context, _ DynamicFee, _ uint32, _ *assets.Wei, _ []EvmPriorAttempt) (bumped DynamicFee, chainSpecificGasLimit uint32, err error) {
	err = errors.New("dynamic fees are not implemented for this layer 2")
	return
}

func (o *l2SuggestedPriceEstimator) GetLegacyGas(ctx context.Context, _ []byte, l2GasLimit uint32, maxGasPriceWei *assets.Wei, opts ...txmgrtypes.Opt) (gasPrice *assets.Wei, chainSpecificGasLimit uint32, err error) {
	chainSpecificGasLimit = l2GasLimit

	ok := o.IfStarted(func() {
		if slices.Contains(opts, txmgrtypes.OptForceRefetch) {
			ch := make(chan struct{})
			select {
			case o.chForceRefetch <- ch:
			case <-o.chStop:
				err = errors.New("estimator stopped")
				return
			case <-ctx.Done():
				err = ctx.Err()
				return
			}
			select {
			case <-ch:
			case <-o.chStop:
				err = errors.New("estimator stopped")
				return
			case <-ctx.Done():
				err = ctx.Err()
				return
			}
		}
		if gasPrice = o.getGasPrice(); gasPrice == nil {
			err = errors.New("failed to estimate l2 gas; gas price not set")
			return
		}
		o.logger.Debugw("GetLegacyGas", "l2GasPrice", gasPrice, "l2GasLimit", l2GasLimit)
	})
	if !ok {
		return nil, 0, errors.New("estimator is not started")
	} else if err != nil {
		return
	}
	// For L2 chains (e.g. Optimism), submitting a transaction that is not priced high enough will cause the call to fail, so if the cap is lower than the RPC suggested gas price, this transaction cannot succeed
	if gasPrice != nil && gasPrice.Cmp(maxGasPriceWei) > 0 {
		return nil, 0, errors.Errorf("estimated gas price: %s is greater than the maximum gas price configured: %s", gasPrice.String(), maxGasPriceWei.String())
	}
	return
}

func (o *l2SuggestedPriceEstimator) BumpLegacyGas(_ context.Context, _ *assets.Wei, _ uint32, _ *assets.Wei, _ []EvmPriorAttempt) (bumpedGasPrice *assets.Wei, chainSpecificGasLimit uint32, err error) {
	return nil, 0, errors.New("bump gas is not supported for this l2")
}

func (o *l2SuggestedPriceEstimator) getGasPrice() (l2GasPrice *assets.Wei) {
	o.gasPriceMu.RLock()
	defer o.gasPriceMu.RUnlock()
	return o.l2GasPrice
}
