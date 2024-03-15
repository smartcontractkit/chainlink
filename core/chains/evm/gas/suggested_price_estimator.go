package gas

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

var (
	_ EvmEstimator = &SuggestedPriceEstimator{}
)

//go:generate mockery --quiet --name rpcClient --output ./mocks/ --case=underscore --structname RPCClient
type rpcClient interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

// SuggestedPriceEstimator is an Estimator which uses the suggested gas price from eth_gasPrice.
type SuggestedPriceEstimator struct {
	services.StateMachine

	client     rpcClient
	pollPeriod time.Duration
	logger     logger.Logger

	gasPriceMu sync.RWMutex
	GasPrice   *assets.Wei

	chForceRefetch chan (chan struct{})
	chInitialised  chan struct{}
	chStop         services.StopChan
	chDone         chan struct{}
}

// NewSuggestedPriceEstimator returns a new Estimator which uses the suggested gas price.
func NewSuggestedPriceEstimator(lggr logger.Logger, client rpcClient) EvmEstimator {
	return &SuggestedPriceEstimator{
		client:         client,
		pollPeriod:     10 * time.Second,
		logger:         logger.Named(lggr, "SuggestedPriceEstimator"),
		chForceRefetch: make(chan (chan struct{})),
		chInitialised:  make(chan struct{}),
		chStop:         make(chan struct{}),
		chDone:         make(chan struct{}),
	}
}

func (o *SuggestedPriceEstimator) Name() string {
	return o.logger.Name()
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

func (o *SuggestedPriceEstimator) refreshPrice() (t *time.Timer) {
	t = time.NewTimer(utils.WithJitter(o.pollPeriod))

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
	return
}

func (o *SuggestedPriceEstimator) OnNewLongestChain(context.Context, *evmtypes.Head) {}

func (*SuggestedPriceEstimator) GetDynamicFee(_ context.Context, _ uint32, _ *assets.Wei) (fee DynamicFee, chainSpecificGasLimit uint32, err error) {
	err = errors.New("dynamic fees are not implemented for this layer 2")
	return
}

func (*SuggestedPriceEstimator) BumpDynamicFee(_ context.Context, _ DynamicFee, _ uint32, _ *assets.Wei, _ []EvmPriorAttempt) (bumped DynamicFee, chainSpecificGasLimit uint32, err error) {
	err = errors.New("dynamic fees are not implemented for this layer 2")
	return
}

func (o *SuggestedPriceEstimator) GetLegacyGas(ctx context.Context, _ []byte, GasLimit uint32, maxGasPriceWei *assets.Wei, opts ...feetypes.Opt) (gasPrice *assets.Wei, chainSpecificGasLimit uint32, err error) {
	chainSpecificGasLimit = GasLimit

	ok := o.IfStarted(func() {
		if slices.Contains(opts, feetypes.OptForceRefetch) {
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
			err = errors.New("failed to estimate gas; gas price not set")
			return
		}
		o.logger.Debugw("GetLegacyGas", "GasPrice", gasPrice, "GasLimit", GasLimit)
	})
	if !ok {
		return nil, 0, errors.New("estimator is not started")
	} else if err != nil {
		return
	}
	// For L2 chains, submitting a transaction that is not priced high enough will cause the call to fail, so if the cap is lower than the RPC suggested gas price, this transaction cannot succeed
	if gasPrice != nil && gasPrice.Cmp(maxGasPriceWei) > 0 {
		return nil, 0, errors.Errorf("estimated gas price: %s is greater than the maximum gas price configured: %s", gasPrice.String(), maxGasPriceWei.String())
	}
	return
}

func (o *SuggestedPriceEstimator) BumpLegacyGas(_ context.Context, _ *assets.Wei, _ uint32, _ *assets.Wei, _ []EvmPriorAttempt) (bumpedGasPrice *assets.Wei, chainSpecificGasLimit uint32, err error) {
	return nil, 0, errors.New("bump gas is not supported for this chain")
}

func (o *SuggestedPriceEstimator) getGasPrice() (GasPrice *assets.Wei) {
	o.gasPriceMu.RLock()
	defer o.gasPriceMu.RUnlock()
	return o.GasPrice
}
