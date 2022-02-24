package gas

import (
	"context"
	"math/big"
	"sync"
	"time"

	optimismfees "github.com/ethereum-optimism/go-optimistic-ethereum-utils/fees"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"go.uber.org/multierr"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	_ Estimator = &optimismEstimator{}
	_ Estimator = &optimism2Estimator{}
)

//go:generate mockery --name optimismRPCClient --output ./mocks/ --case=underscore --structname OptimismRPCClient
type optimismRPCClient interface {
	Call(result interface{}, method string, args ...interface{}) error
}

type optimismEstimator struct {
	utils.StartStopOnce

	config     Config
	client     optimismRPCClient
	pollPeriod time.Duration
	logger     logger.Logger

	gasPriceMu sync.RWMutex
	l1GasPrice *big.Int
	l2GasPrice *big.Int

	chForceRefetch chan (chan struct{})
	chInitialised  chan struct{}
	chStop         chan struct{}
	chDone         chan struct{}
}

// NewOptimismEstimator returns a new optimism estimator
func NewOptimismEstimator(lggr logger.Logger, config Config, client optimismRPCClient) Estimator {
	return &optimismEstimator{
		utils.StartStopOnce{},
		config,
		client,
		10 * time.Second,
		lggr.Named("OptimismEstimator"),
		sync.RWMutex{},
		nil,
		nil,
		make(chan (chan struct{})),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
	}
}

func (o *optimismEstimator) Start(context.Context) error {
	return o.StartOnce("OptimismEstimator", func() error {
		go o.run()
		<-o.chInitialised
		return nil
	})
}
func (o *optimismEstimator) Close() error {
	return o.StopOnce("OptimismEstimator", func() error {
		close(o.chStop)
		<-o.chDone
		return nil
	})
}

func (o *optimismEstimator) run() {
	defer close(o.chDone)

	t := o.refreshPrices()
	close(o.chInitialised)

	for {
		select {
		case <-o.chStop:
			return
		case ch := <-o.chForceRefetch:
			t.Stop()
			t = o.refreshPrices()
			close(ch)
		case <-t.C:
			t = o.refreshPrices()
		}
	}
}

// OptimismGasPricesResponse is the shape of the response when calling rollup_gasPrices
type OptimismGasPricesResponse struct {
	L1GasPrice *big.Int
	L2GasPrice *big.Int
}

func (g *OptimismGasPricesResponse) UnmarshalJSON(b []byte) error {
	var l1Hex string = gjson.GetBytes(b, "l1GasPrice").Str
	var l2Hex string = gjson.GetBytes(b, "l2GasPrice").Str
	l1 := new(hexutil.Big)
	l2 := new(hexutil.Big)
	if err := multierr.Combine(l1.UnmarshalText([]byte(l1Hex)), l2.UnmarshalText([]byte(l2Hex))); err != nil {
		return err
	}
	g.L1GasPrice = l1.ToInt()
	g.L2GasPrice = l2.ToInt()
	return nil
}

func (o *optimismEstimator) refreshPrices() (t *time.Timer) {
	var res OptimismGasPricesResponse
	t = time.NewTimer(utils.WithJitter(o.pollPeriod))

	if err := o.client.Call(&res, "rollup_gasPrices"); err != nil {
		o.logger.Warnf("Failed to refresh prices, got error: %s", err)
		return
	}

	o.logger.Debugw("OptimismEstimator#refreshPrices", "l1GasPrice", res.L1GasPrice, "l2GasPrice", res.L2GasPrice)

	o.gasPriceMu.Lock()
	defer o.gasPriceMu.Unlock()
	o.l1GasPrice, o.l2GasPrice = res.L1GasPrice, res.L2GasPrice
	return
}

func (o *optimismEstimator) GetLegacyGas(calldata []byte, gasLimit uint64, opts ...Opt) (gasPrice *big.Int, chainSpecificGasLimit uint64, err error) {
	ok := o.IfStarted(func() {
		var forceRefetch bool
		for _, opt := range opts {
			if opt == OptForceRefetch {
				forceRefetch = true
			}
		}
		if forceRefetch {
			ch := make(chan struct{})
			o.chForceRefetch <- ch
			select {
			case <-ch:
			case <-o.chStop:
				err = errors.New("estimator stopped")
				return
			}
		}
		gasPrice, chainSpecificGasLimit, err = o.calcGas(calldata, gasLimit)
	})
	if !ok {
		return nil, 0, errors.New("estimator is not started")
	}
	return
}

func (o *optimismEstimator) BumpLegacyGas(originalGasPrice *big.Int, originalGasLimit uint64) (gasPrice *big.Int, gasLimit uint64, err error) {
	return nil, 0, errors.New("bump gas is not supported for optimism")
}

func (o *optimismEstimator) OnNewLongestChain(_ context.Context, _ *evmtypes.Head) {}

func (*optimismEstimator) GetDynamicFee(gasLimit uint64) (fee DynamicFee, chainSpecificGasLimit uint64, err error) {
	err = errors.New("dynamic fees are not implemented for Optimism")
	return
}

func (o *optimismEstimator) BumpDynamicFee(original DynamicFee, gasLimit uint64) (bumped DynamicFee, chainSpecificGasLimit uint64, err error) {
	err = errors.New("dynamic fees are not implemented for Optimism")
	return
}

func (o *optimismEstimator) calcGas(calldata []byte, l2GasLimit uint64) (chainSpecificGasPrice *big.Int, chainSpecificGasLimit uint64, err error) {
	l1GasPrice, l2GasPrice := o.getGasPrices()
	if l1GasPrice == nil || l2GasPrice == nil {
		return nil, 0, errors.New("failed to estimate optimism gas; gas prices not set")
	}

	optimismGasLimitBig := optimismfees.EncodeTxGasLimit(calldata, l1GasPrice, big.NewInt(int64(l2GasLimit)), l2GasPrice)
	if !optimismGasLimitBig.IsInt64() {
		o.logger.Errorw(
			"Optimism: unable to represent gas limit as Int64, this is an unexpected error and should be reported to the Chainlink team",
			"calldata", calldata,
			"l2GasLimit", l2GasLimit,
			"l1GasPrice", l1GasPrice,
			"l2GasPrice", l2GasPrice,
			"optimismGasLimitBig", optimismGasLimitBig.String(),
		)
		return nil, 0, errors.New("gas limit overflows int64")
	}
	chainSpecificGasLimit = uint64(optimismGasLimitBig.Int64())

	// It's always 0.015 GWei
	// See: https://www.notion.so/How-to-pay-Fees-in-Optimistic-Ethereum-f706f4e5b13e460fa5671af48ce9a695
	const optimisml1GasPrice = 15000000

	o.logger.Debugw("OptimismEstimator#EstimateGas", "l1GasPrice", l1GasPrice, "l2GasPrice", l2GasPrice, "l2GasLimit", l2GasLimit, "chainSpecificGasLimit", chainSpecificGasLimit, "optimisml1GasPrice", optimisml1GasPrice)
	return big.NewInt(optimisml1GasPrice), chainSpecificGasLimit, nil
}

func (o *optimismEstimator) getGasPrices() (l1GasPrice, l2GasPrice *big.Int) {
	o.gasPriceMu.RLock()
	defer o.gasPriceMu.RUnlock()
	return o.l1GasPrice, o.l2GasPrice
}

type optimism2Estimator struct {
	utils.StartStopOnce

	config     Config
	client     optimismRPCClient
	pollPeriod time.Duration
	logger     logger.Logger

	gasPriceMu sync.RWMutex
	l2GasPrice *big.Int

	chForceRefetch chan (chan struct{})
	chInitialised  chan struct{}
	chStop         chan struct{}
	chDone         chan struct{}
}

// NewOptimism2Estimator returns a new optimism 2.0 estimator
func NewOptimism2Estimator(lggr logger.Logger, config Config, client optimismRPCClient) Estimator {
	return &optimism2Estimator{
		utils.StartStopOnce{},
		config,
		client,
		10 * time.Second,
		lggr.Named("Optimism2Estimator"),
		sync.RWMutex{},
		nil,
		make(chan (chan struct{})),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
	}
}

func (o *optimism2Estimator) Start(context.Context) error {
	return o.StartOnce("Optimism2Estimator", func() error {
		go o.run()
		<-o.chInitialised
		return nil
	})
}
func (o *optimism2Estimator) Close() error {
	return o.StopOnce("Optimism2Estimator", func() error {
		close(o.chStop)
		<-o.chDone
		return nil
	})
}

func (o *optimism2Estimator) run() {
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

func (o *optimism2Estimator) refreshPrice() (t *time.Timer) {
	t = time.NewTimer(utils.WithJitter(o.pollPeriod))

	var res hexutil.Big
	if err := o.client.Call(&res, "eth_gasPrice"); err != nil {
		o.logger.Warnf("Optimism2Estimator: Failed to refresh prices, got error: %s", err)
		return
	}
	bi := (*big.Int)(&res)

	o.logger.Debugw("Optimism2Estimator#refreshPrice", "l2GasPrice", bi)

	o.gasPriceMu.Lock()
	defer o.gasPriceMu.Unlock()
	o.l2GasPrice = bi
	return
}

func (o *optimism2Estimator) OnNewLongestChain(_ context.Context, _ *evmtypes.Head) {}

func (*optimism2Estimator) GetDynamicFee(_ uint64) (fee DynamicFee, chainSpecificGasLimit uint64, err error) {
	err = errors.New("dynamic fees are not implemented for Optimism")
	return
}

func (*optimism2Estimator) BumpDynamicFee(_ DynamicFee, _ uint64) (bumped DynamicFee, chainSpecificGasLimit uint64, err error) {
	err = errors.New("dynamic fees are not implemented for Optimism")
	return
}

func (o *optimism2Estimator) GetLegacyGas(_ []byte, l2GasLimit uint64, opts ...Opt) (gasPrice *big.Int, chainSpecificGasLimit uint64, err error) {
	chainSpecificGasLimit = l2GasLimit
	ok := o.IfStarted(func() {
		var forceRefetch bool
		for _, opt := range opts {
			if opt == OptForceRefetch {
				forceRefetch = true
			}
		}
		if forceRefetch {
			ch := make(chan struct{})
			o.chForceRefetch <- ch
			select {
			case <-ch:
			case <-o.chStop:
				err = errors.New("estimator stopped")
				return
			}
		}
		if gasPrice = o.getGasPrice(); gasPrice == nil {
			err = errors.New("failed to estimate optimism gas; gas price not set")
			return
		}
		o.logger.Debugw("Optimism2Estimator#EstimateGas", "l2GasPrice", gasPrice, "l2GasLimit", l2GasLimit)
	})
	if !ok {
		return nil, 0, errors.New("estimator is not started")
	}
	return
}

func (o *optimism2Estimator) BumpLegacyGas(_ *big.Int, _ uint64) (bumpedGasPrice *big.Int, chainSpecificGasLimit uint64, err error) {
	return nil, 0, errors.New("bump gas is not supported for optimism")
}

func (o *optimism2Estimator) getGasPrice() (l2GasPrice *big.Int) {
	o.gasPriceMu.RLock()
	defer o.gasPriceMu.RUnlock()
	return o.l2GasPrice
}
