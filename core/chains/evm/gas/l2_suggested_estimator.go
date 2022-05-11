package gas

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	_ Estimator = &l2SuggestedEstimator{}
)

//go:generate mockery --name rpcClient --output ./mocks/ --case=underscore --structname OptimismRPCClient
type rpcClient interface {
	Call(result interface{}, method string, args ...interface{}) error
}

type l2SuggestedEstimator struct {
	utils.StartStopOnce

	config     Config
	client     rpcClient
	pollPeriod time.Duration
	logger     logger.Logger

	gasPriceMu sync.RWMutex
	l2GasPrice *big.Int

	chForceRefetch chan (chan struct{})
	chInitialised  chan struct{}
	chStop         chan struct{}
	chDone         chan struct{}
}

// NewL2SuggestedEstimator returns a new optimism 2.0 estimator
func NewL2SuggestedEstimator(lggr logger.Logger, config Config, client rpcClient) Estimator {
	return &l2SuggestedEstimator{
		utils.StartStopOnce{},
		config,
		client,
		10 * time.Second,
		lggr.Named("L2SuggestedEstimator"),
		sync.RWMutex{},
		nil,
		make(chan (chan struct{})),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
	}
}

func (o *l2SuggestedEstimator) Start(context.Context) error {
	return o.StartOnce("L2SuggestedEstimator", func() error {
		go o.run()
		<-o.chInitialised
		return nil
	})
}
func (o *l2SuggestedEstimator) Close() error {
	return o.StopOnce("L2SuggestedEstimator", func() error {
		close(o.chStop)
		<-o.chDone
		return nil
	})
}

func (o *l2SuggestedEstimator) run() {
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

func (o *l2SuggestedEstimator) refreshPrice() (t *time.Timer) {
	t = time.NewTimer(utils.WithJitter(o.pollPeriod))

	var res hexutil.Big
	if err := o.client.Call(&res, "eth_gasPrice"); err != nil {
		o.logger.Warnf("Failed to refresh prices, got error: %s", err)
		return
	}
	bi := (*big.Int)(&res)

	o.logger.Debugw("refreshPrice", "l2GasPrice", bi)

	o.gasPriceMu.Lock()
	defer o.gasPriceMu.Unlock()
	o.l2GasPrice = bi
	return
}

func (o *l2SuggestedEstimator) OnNewLongestChain(_ context.Context, _ *evmtypes.Head) {}

func (*l2SuggestedEstimator) GetDynamicFee(_ uint64) (fee DynamicFee, chainSpecificGasLimit uint64, err error) {
	err = errors.New("dynamic fees are not implemented for this layer 2")
	return
}

func (*l2SuggestedEstimator) BumpDynamicFee(_ DynamicFee, _ uint64) (bumped DynamicFee, chainSpecificGasLimit uint64, err error) {
	err = errors.New("dynamic fees are not implemented for this layer 2")
	return
}

func (o *l2SuggestedEstimator) GetLegacyGas(_ []byte, l2GasLimit uint64, opts ...Opt) (gasPrice *big.Int, chainSpecificGasLimit uint64, err error) {
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
			err = errors.New("failed to estimate l2 gas; gas price not set")
			return
		}
		o.logger.Debugw("EstimateGas", "l2GasPrice", gasPrice, "l2GasLimit", l2GasLimit)
	})
	if !ok {
		return nil, 0, errors.New("estimator is not started")
	}
	return
}

func (o *l2SuggestedEstimator) BumpLegacyGas(_ *big.Int, _ uint64) (bumpedGasPrice *big.Int, chainSpecificGasLimit uint64, err error) {
	return nil, 0, errors.New("bump gas is not supported for this l2")
}

func (o *l2SuggestedEstimator) getGasPrice() (l2GasPrice *big.Int) {
	o.gasPriceMu.RLock()
	defer o.gasPriceMu.RUnlock()
	return o.l2GasPrice
}
