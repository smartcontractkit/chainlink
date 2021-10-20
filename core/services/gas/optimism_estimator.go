package gas

import (
	"context"
	"math/big"
	"sync"
	"time"

	optimismfees "github.com/ethereum-optimism/go-optimistic-ethereum-utils/fees"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/tidwall/gjson"
	"go.uber.org/multierr"
)

// It's always 0.015 GWei
// See: https://www.notion.so/How-to-pay-Fees-in-Optimistic-Ethereum-f706f4e5b13e460fa5671af48ce9a695
const optimisml1GasPrice = 15000000

var _ Estimator = &optimismEstimator{}

//go:generate mockery --name optimismRPCClient --output ./mocks/ --case=underscore --structname OptimismRPCClient
type optimismRPCClient interface {
	Call(result interface{}, method string, args ...interface{}) error
}

type optimismEstimator struct {
	utils.StartStopOnce

	config     Config
	client     optimismRPCClient
	pollPeriod time.Duration

	gasPriceMu sync.RWMutex
	l1GasPrice *big.Int
	l2GasPrice *big.Int

	chForceRefetch chan (chan struct{})
	chInitialised  chan struct{}
	chStop         chan struct{}
	chDone         chan struct{}
}

// NewOptimismEstimator returns a new optimism estimator
func NewOptimismEstimator(config Config, client optimismRPCClient) Estimator {
	return &optimismEstimator{
		utils.StartStopOnce{},
		config,
		client,
		10 * time.Second,
		sync.RWMutex{},
		nil,
		nil,
		make(chan (chan struct{})),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
	}
}

func (o *optimismEstimator) Start() error {
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
		logger.Warnf("OptimismEstimator: Failed to refresh prices, got error: %s", err)
		return
	}

	logger.Debugw("OptimismEstimator#refreshPrices", "l1GasPrice", res.L1GasPrice, "l2GasPrice", res.L2GasPrice)

	o.gasPriceMu.Lock()
	defer o.gasPriceMu.Unlock()
	o.l1GasPrice, o.l2GasPrice = res.L1GasPrice, res.L2GasPrice
	return
}

func (o *optimismEstimator) EstimateGas(calldata []byte, gasLimit uint64, opts ...Opt) (gasPrice *big.Int, chainSpecificGasLimit uint64, err error) {
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

func (o *optimismEstimator) BumpGas(originalGasPrice *big.Int, originalGasLimit uint64) (gasPrice *big.Int, gasLimit uint64, err error) {
	return nil, 0, errors.New("bump gas is not supported for optimism")
}

func (o *optimismEstimator) OnNewLongestChain(_ context.Context, _ models.Head) {}

func (o *optimismEstimator) calcGas(calldata []byte, l2GasLimit uint64) (chainSpecificGasPrice *big.Int, chainSpecificGasLimit uint64, err error) {
	l1GasPrice, l2GasPrice := o.getGasPrices()
	if l1GasPrice == nil || l2GasPrice == nil {
		return nil, 0, errors.New("failed to estimate optimism gas; gas prices not set")
	}

	optimismGasLimitBig := optimismfees.EncodeTxGasLimit(calldata, l1GasPrice, big.NewInt(int64(l2GasLimit)), l2GasPrice)
	if !optimismGasLimitBig.IsInt64() {
		logger.Errorw(
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

	logger.Debugw("OptimismEstimator#EstimateGas", "l1GasPrice", l1GasPrice, "l2GasPrice", l2GasPrice, "l2GasLimit", l2GasLimit, "chainSpecificGasLimit", chainSpecificGasLimit, "optimisml1GasPrice", optimisml1GasPrice)
	return big.NewInt(optimisml1GasPrice), chainSpecificGasLimit, nil
}

func (o *optimismEstimator) getGasPrices() (l1GasPrice, l2GasPrice *big.Int) {
	o.gasPriceMu.RLock()
	defer o.gasPriceMu.RUnlock()
	return o.l1GasPrice, o.l2GasPrice
}
