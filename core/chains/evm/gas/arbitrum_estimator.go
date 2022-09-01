package gas

import (
	"context"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type ethClient interface {
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}

// arbitrumEstimator is an Estimator which extends l2SuggestedPriceEstimator to use getPricesInArbGas() for gas limit estimation.
type arbitrumEstimator struct {
	utils.StartStopOnce

	Estimator // *l2SuggestedPriceEstimator

	client     ethClient
	pollPeriod time.Duration
	logger     logger.Logger

	getPricesInArbGasMu sync.RWMutex
	perL2Tx             uint32
	perL1CalldataUnit   uint32

	chForceRefetch chan (chan struct{})
	chInitialised  chan struct{}
	chStop         chan struct{}
	chDone         chan struct{}
}

func NewArbitrumEstimator(lggr logger.Logger, rpcClient rpcClient, ethClient ethClient) Estimator {
	lggr = lggr.Named("ArbitrumEstimator")
	return &arbitrumEstimator{
		Estimator:      NewL2SuggestedPriceEstimator(lggr, rpcClient),
		client:         ethClient,
		pollPeriod:     10 * time.Second,
		logger:         lggr,
		chForceRefetch: make(chan (chan struct{})),
		chInitialised:  make(chan struct{}),
		chStop:         make(chan struct{}),
		chDone:         make(chan struct{}),
	}
}

func (a *arbitrumEstimator) Start(ctx context.Context) error {
	return a.StartOnce("ArbitrumEstimator", func() error {
		if err := a.Estimator.Start(ctx); err != nil {
			return errors.Wrap(err, "failed to start gas price estimator")
		}
		go a.run()
		<-a.chInitialised
		return nil
	})
}
func (a *arbitrumEstimator) Close() error {
	return a.StopOnce("ArbitrumEstimator", func() (err error) {
		close(a.chStop)
		err = errors.Wrap(a.Estimator.Close(), "failed to stop gas price estimator")
		<-a.chDone
		return
	})
}

// GetLegacyGas estimates both the gas price and the gas limit.
func (a *arbitrumEstimator) GetLegacyGas(calldata []byte, l2GasLimit uint32, maxGasPriceWei *big.Int, opts ...Opt) (gasPrice *big.Int, chainSpecificGasLimit uint32, err error) {
	gasPrice, _, err = a.Estimator.GetLegacyGas(calldata, l2GasLimit, maxGasPriceWei, opts...)
	if err != nil {
		return
	}
	ok := a.IfStarted(func() {
		if slices.Contains(opts, OptForceRefetch) {
			ch := make(chan struct{})
			select {
			case a.chForceRefetch <- ch:
			case <-a.chStop:
				err = errors.New("estimator stopped")
				return
			}
			select {
			case <-ch:
			case <-a.chStop:
				err = errors.New("estimator stopped")
				return
			}
		}
		perL2Tx, perL1CalldataUnit := a.getPricesInArbGas()
		// TODO is it ok to return lsGasLimit if unitialized?
		chainSpecificGasLimit = l2GasLimit + perL2Tx + uint32(len(calldata))*perL1CalldataUnit
		a.logger.Debugw("GetLegacyGas", "l2GasLimit", l2GasLimit, "calldataLen", len(calldata), "perL2Tx", perL2Tx,
			"perL1CalldataUnit", perL1CalldataUnit, "chainSpecificGasLimit", chainSpecificGasLimit)
	})
	if !ok {
		return nil, 0, errors.New("estimator is not started")
	} else if err != nil {
		return
	}
	//TODO enforce a maximum? (limit, or overall fee? txm could limit overall fee instead....)
	return
}

func (a *arbitrumEstimator) getPricesInArbGas() (perL2Tx uint32, perL1CalldataUnit uint32) {
	a.getPricesInArbGasMu.RLock()
	perL2Tx, perL1CalldataUnit = a.perL2Tx, a.perL1CalldataUnit
	a.getPricesInArbGasMu.RUnlock()
	return
}

func (a *arbitrumEstimator) run() {
	defer close(a.chDone)

	t := a.refreshPricesInArbGas()
	close(a.chInitialised)

	for {
		select {
		case <-a.chStop:
			return
		case ch := <-a.chForceRefetch:
			t.Stop()
			t = a.refreshPricesInArbGas()
			close(ch)
		case <-t.C:
			t = a.refreshPricesInArbGas()
		}
	}
}

// refreshPricesInArbGas calls getPricesInArbGas() on the precompile contract 0x000000000000000000000000000000000000006c.
func (a *arbitrumEstimator) refreshPricesInArbGas() (t *time.Timer) {
	t = time.NewTimer(utils.WithJitter(a.pollPeriod))

	ctx, cancel := evmclient.ContextWithDefaultTimeoutFromChan(a.chStop)
	defer cancel()

	// @return (per L2 tx, per L1 calldata unit, per storage allocation)
	// function getPricesInArbGas() external view returns (uint256, uint256, uint256);
	//
	// https://github.com/OffchainLabs/nitro/blob/f7645453cfc77bf3e3644ea1ac031eff629df325/contracts/src/precompiles/ArbGasInfo.sol#L69
	precompile := common.HexToAddress("0x000000000000000000000000000000000000006c")
	b, err := a.client.CallContract(ctx, ethereum.CallMsg{
		To:   &precompile,
		Data: common.Hex2Bytes("02199f34"),
	}, big.NewInt(-1))
	if err != nil {
		a.logger.Warnf("Failed to refresh prices, got error: %s", err)
		return
	}

	if len(b) != 3*32 { // returns (uint256, uint256, uint256);
		a.logger.Errorf("Failed to refresh prices, return data length (%d) different than expected (%d)", len(b), 3*32)
		return
	}
	bPerL2Tx := new(big.Int).SetBytes(b[:32])
	bPerL1CalldataUnit := new(big.Int).SetBytes(b[32:64])
	// ignore perStorageAllocation
	if !bPerL2Tx.IsUint64() || !bPerL1CalldataUnit.IsUint64() {
		a.logger.Errorf("Failed to refresh prices, returned integers are not uint64", "perL2Tx", bPerL2Tx.String(),
			"perL1CalldataUnit", bPerL1CalldataUnit.String())
	}

	perL2Tx := bPerL2Tx.Uint64()
	perL1CalldataUnit := bPerL1CalldataUnit.Uint64()
	if perL2Tx > math.MaxUint32 || perL1CalldataUnit > math.MaxUint32 {
		a.logger.Errorf("Failed to refresh prices, returned integers are not uint32", "perL2Tx", perL2Tx,
			"perL1CalldataUnit", perL1CalldataUnit)
	}

	a.logger.Debugw("refreshPricesInArbGas", "perL2Tx", perL2Tx, "perL2CalldataUnit", perL1CalldataUnit)

	a.getPricesInArbGasMu.Lock()
	defer a.getPricesInArbGasMu.Unlock()
	a.perL2Tx = uint32(perL2Tx)
	a.perL1CalldataUnit = uint32(perL1CalldataUnit)
	return
}
